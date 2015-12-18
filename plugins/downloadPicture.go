package plugins

import (
	"github.com/Sirupsen/logrus"
	"go_distributed_spider/config"
	"go_distributed_spider/util"
	"regexp"
	"strings"
)

type Picture struct {
	Config *config.Config
}

func NewPicture(conf *config.Config) *Picture {
	return &Picture{Config: conf}
}

func (p *Picture) Init() {
	method := "Picture init"
	//picture path check.
	err := util.MakedirPath(p.Config.Controller.PicturePath)
	if err != nil {
		logrus.Errorln(method, err)
	}
}

func (p *Picture) Done(url string, content []byte) {
	method := "Picture Done"
	//if save the picture
	pictureNum, err := p.MatchPictureOnce(url, string(content))
	if err != nil {
		logrus.Errorln(method, err)
		return
	}
	logrus.Infof("picture num: %d", pictureNum)
}

func (p *Picture) MatchPictureOnce(parentUrl, body string) (int, error) {
	method := "MatchPictureOnce"
	_ = method
	//<\s*(a|img)\s*(href|src)="(http://)?[0-9a-zA-Z/_\.#-]*[^/"]
	//find all have link part
	urlMatch := regexp.MustCompile(`<\s*(a|img)\s*(src)="(http://|https://)?[0-9a-zA-Z/_\.#-]*[^/"]`)
	urls := urlMatch.FindAllString(body, -1)

	//parse parent url.use for getting like "/abc".
	parseParentUrl := regexp.MustCompile(`(http://|https://)?[0-9a-zA-Z_\.#-]*[^/]`)
	parentUrlHead := parseParentUrl.FindAllString(parentUrl, 1)

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

	for _, val := range urls {
		logrus.Infoln("pictur url:", val)
		if !strings.HasPrefix(val, "http") {
			continue
		}

		// name := path.Base(val)
		name := val

		/*r := rand.New(rand.NewSource(time.Now().UnixNano()))
		prefix := r.Intn(10000)*/

		err := util.Download(p.Config.Controller.PicturePath /*strconv.Itoa(prefix)+"-"+*/, name, val)
		if err != nil {
			logrus.Errorln(method, err)
			continue
		}
	}
	return len(urls), nil
}
