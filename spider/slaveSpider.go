package spider

import (
	"fmt"
	// "net/http"
	"encoding/json"
	"github.com/Sirupsen/logrus"
	"go_distributed_spider/config"
	"go_distributed_spider/util"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type SlaveSpider struct {
	// Maps map[string]func()
	Target      string
	Config      *config.Config
	DailTimeout time.Duration
	Urls        []string //tasks
	NewUrls     []string //send to master
	Match       []string //task matchs
}

func NewSlaveSpider(config *config.Config) *SlaveSpider {
	method := "NewSlaveSpider"
	spider := &SlaveSpider{Urls: []string{}, NewUrls: []string{}, Config: config}
	/*spider.Maps = map[string]func() {
		Url: spider.Base,
	}*/
	for i, val := range config.Target {
		spider.Target += val
		if i < len(config.Target)-1 {
			spider.Target += "|"
		}
	}
	logrus.Debugln(method, "target is:", spider.Target)
	if config.DailTimeout != 0 {
		spider.DailTimeout = time.Duration(config.DailTimeout) * time.Second
	} else {
		spider.DailTimeout = DefaultDailTimeout
	}
	fmt.Println("time out is :", spider.DailTimeout)
	return spider
}

func (s *SlaveSpider) Base() {
	method := "SlaveSpider Base"
	logrus.Println("SlaveSpider", method)

	for {
		if len(s.Urls) == 0 {
			// logrus.Infoln(method, "Send to Master New urls:", len(s.NewUrls) /*s.NewUrls*/)
			send := len(s.NewUrls)

			masterToSlave, err := s.Request()
			if err != nil {
				logrus.Errorln(method, err)
				continue
			}

			s.NewUrls = []string{}
			s.Match = []string{}

			logrus.Infoln(method, "slaveToMaster url:", send, "get Task:", len(masterToSlave.Task) /*, masterToSlave.Task*/)
			s.Urls = masterToSlave.Task
			time.Sleep(time.Second * 2)
		} else {
			err := s.MatchOnce()
			if err != nil {
				logrus.Errorln(method, err)
				continue
			}
		}

	}

	return
}

func (s *SlaveSpider) Request() (*MasterToSlave, error) {
	method := "SlaveSpider Request"
	/*arr := []string{}
	for i := 0; i < num; i++ {
		arr = append(arr, "give"+strconv.Itoa(num)+strconv.Itoa(i))
	}
	num--*/
	var slaveToMaster SlaveToMaster
	slaveToMaster.Match = s.Match
	slaveToMaster.NewUrls = s.NewUrls

	req, err := json.Marshal(slaveToMaster)
	if err != nil {
		logrus.Errorln(method, err)
		return nil, err
	}
	url := s.Config.GetUrl() + Url

	resp, err := http.Post(url, "application/json", strings.NewReader(string(req)))
	if err != nil {
		logrus.Errorln(method, "2", err)
		return nil, err
	}
	body := resp.Body
	defer body.Close()
	bodyBytes, err := ioutil.ReadAll(body)
	if err != nil {
		logrus.Errorln(method, "3", err)
		return nil, err
	}

	var masterToSlave *MasterToSlave

	err = json.Unmarshal(bodyBytes, &masterToSlave)
	if err != nil {
		logrus.Errorln(method, "4", err)
		return nil, err
	}
	return masterToSlave, nil
}

func (s *SlaveSpider) MatchOnce() error {
	method := "SlaveSpider MatchOnce"

	if len(s.Urls[0]) == 0 {
		return fmt.Errorf("none unMatch url")
	}
	url := s.Urls[0]
	s.Urls = s.Urls[1:]

	client := &http.Client{

		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				deadline := time.Now().Add(s.DailTimeout)
				c, err := net.DialTimeout(netw, addr, s.DailTimeout)
				if err != nil {
					// fmt.Println("dail timeout", err)
					return nil, err
				}
				c.SetDeadline(deadline)
				return c, nil
			},
		},
	}
	resp, err := client.Get(url)
	if err != nil {
		logrus.Infoln(method, "http get error", err)
		return err
	}

	body := resp.Body
	defer body.Close()
	bodyByte, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}

	// logrus.Println("resp is ", string(bodyByte))
	countAvailable, err := s.MatchUrlOnce(url, string(bodyByte))
	if err != nil {
		logrus.Errorln(method, err)
		return err
	}

	logrus.Infof("%s Search URL: %s. Get child: %d. ", method, url, countAvailable)
	if s.Config.IfPicture {
		pictureNum, err := s.MatchPictureOnce(url, string(bodyByte))
		if err != nil {
			logrus.Errorln(method, err)
			return err
		}
		logrus.Infof("picture num: %d\n", pictureNum)
	} else {
		logrus.Infof("\n")
	}

	resp.Body.Close()
	return nil
}

func (s *SlaveSpider) MatchUrlOnce(url, body string) (int, error) {
	//match child urls in the page.
	matchUrls := MatchUrls(url, body)
	countAvailable := 0
	for _, url := range matchUrls {
		if url == "" {
			continue
		} else {
			s.NewUrls = append(s.NewUrls, url)
			countAvailable++
		}
	}

	//if this url has target.Add to the match.
	if IsTarget(body, s.Target) {
		s.Match = append(s.Match, url)
	}

	return countAvailable, nil
}

func (s *SlaveSpider) MatchPictureOnce(parentUrl, body string) (int, error) {
	method := "MatchPictureOnce"
	_ = method
	//<\s*(a|img)\s*(href|src)="(http://)?[0-9a-zA-Z/_\.#-]*[^/"]
	//find all have link part
	urlMatch := regexp.MustCompile(`<\s*(a|img)\s*(src)="(http://|https://)?[0-9a-zA-Z/_\.#-]*[^/"]`)
	urls := urlMatch.FindAllString(body, -1)

	//parse parent url.use for getting like "/abc".
	parseParentUrl := regexp.MustCompile(`(http://|https://)?[0-9a-zA-Z_\.#-]*[^/]`)
	parentUrlHead := parseParentUrl.FindAllString(parentUrl, 1)

	fmt.Println("debug .urls before:", urls)
	geturlMatch := regexp.MustCompile(`src="(http://|https://)?[0-9a-zA-Z/_\.#-]*`)
	for i := 0; i < len(urls); i++ {
		href := geturlMatch.FindAllString(urls[i], 1)
		if len(href) == 0 || len(href[0]) <= 6 {
			// logrus.Println(method, "len(href) is ", href, "and url is ", urls[i])
			urls[i] = ""
			continue
		}
		url := href[0][5:]
		if strings.HasPrefix(url, "/") {
			if len(parentUrlHead) <= 0 {
				url = ""
			} else {
				url = parentUrlHead[0] + url
			}
		}
		if strings.Contains(url, "#") {
			url = ""
		}
		url = strings.TrimSpace(url)
		urls[i] = strings.TrimSuffix(url, "/")
	}

	fmt.Println("debug .urls after:", urls)

	for _, val := range urls {
		logrus.Infoln("pictur url:", val)
		if !strings.HasPrefix(val, "http") {
			continue
		}

		name := path.Base(val)
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		prefix := r.Intn(10000)

		err := util.DownPicture(s.Config.PicturePath, strconv.Itoa(prefix)+"-"+name, val)
		if err != nil {
			logrus.Errorln(method, err)
			continue
		}
	}
	return len(urls), nil
}

func MatchUrls(parentUrl, body string) []string {
	method := "MatchTargetUrl"
	_ = method
	//<\s*(a|img)\s*(href|src)="(http://)?[0-9a-zA-Z/_\.#-]*[^/"]
	//find all have link part
	urlMatch := regexp.MustCompile(`<\s*(a)\s*(href)="(http://|https://)?[0-9a-zA-Z/_\.#-]*[^/"]`)
	urls := urlMatch.FindAllString(body, -1)

	//parse parent url.use for getting like "/abc".
	parseParentUrl := regexp.MustCompile(`(http://|https://)?[0-9a-zA-Z_\.#-]*[^/]`)
	parentUrlHead := parseParentUrl.FindAllString(parentUrl, 1)

	geturlMatch := regexp.MustCompile(`href="(http://|https://)?[0-9a-zA-Z/_\.#-]*`)
	for i := 0; i < len(urls); i++ {
		href := geturlMatch.FindAllString(urls[i], 1)
		if len(href) == 0 || len(href[0]) <= 7 {
			// logrus.Println(method, "len(href) is ", href, "and url is ", urls[i])
			urls[i] = ""
			continue
		}
		url := href[0][6:]
		if strings.HasPrefix(url, "/") {
			if len(parentUrlHead) <= 0 {
				url = ""
			} else {
				url = parentUrlHead[0] + url
			}
		}
		if strings.Contains(url, "#") {
			url = ""
		}
		url = strings.TrimSpace(url)
		urls[i] = strings.TrimSuffix(url, "/")
	}
	return urls
}

func IsTarget(body string, target string) bool {
	targetMatch := regexp.MustCompile(target)
	res := targetMatch.FindAllString(body, 1)

	if len(res) <= 0 {
		return false
	} else {
		return true
	}
}
