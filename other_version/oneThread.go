package main

/**
 * @author carlding
 * @date : 2015-9-25
 */

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"time"
)

/*
*单线程跑真的超级慢。所以干脆不看了。
 */

type Address struct {
	Url string
	Has bool
}

var logger *log.Logger
var Visit map[string]Address
var target []string = []string{"golang", "docker"}
var count int64 = 0

//flag
var (
	help        *bool   = flag.Bool("help", false, "see help detail.")
	saveFile    *string = flag.String("file", "ret.txt", "save result into file.")
	hasFile     *string = flag.String("has", "has.txt", "save has result into file.")
	timeout     *int    = flag.Int("timeout", 2, "dail tcp connect timeout.default is 2s.")
	giveUrls    *string = flag.String("urls", "http://hedengcheng.com", "give spider urls,use ',' split the urls")
	giveTargets *string = flag.String("targets", "golang,docker", "give spider targets,use ',' split the target")
)

func init() {
	// logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	logger = log.New(os.Stdout, "", log.Ldate|log.Ltime)
}

func main() {
	flag.Parse()

	if *help {
		helpFunc()
		return
	}

	urls := strings.Split(*giveUrls, ",")
	target = strings.Split(*giveTargets, ",")
	if len(urls) <= 0 || len(target) <= 0 {
		logger.Println("please give right urls.")
		return
	}

	//listen ctrl + c.and save the urls to file before exit.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			_ = sig
			logger.Println("Stoping by Ctrl + C ...... ")
			OutputRes()
			os.Exit(0)
		}
	}()

	Visit = make(map[string]Address, 0)
	AppendAddressArr(urls[0], false)

	//go to recursion
	Spider(urls)

	//output result and save into file before exit.
	defer OutputRes()
}

func OutputRes() {
	//save all result
	allFile, err := os.OpenFile(*saveFile, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		logger.Println("Open file error,", err)
	}
	defer allFile.Close()

	//save has result
	hasFile, err := os.OpenFile(*hasFile, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		logger.Println("Open file error,", err)
	}
	defer hasFile.Close()

	var num int64 = 0
	var has int64 = 0
	for k, address := range Visit {
		str := fmt.Sprintf("%-4d has: %-5t,  url:  %s\n", num, address.Has, k)
		// logger.Printf("%-4d  has: %t,url: %-40s\n", num, address.Has, k)
		logger.Printf(str)
		allFile.WriteString(str)
		if address.Has {
			hasFile.WriteString(fmt.Sprintf("%-4d has: %-5t,  url:  %s\n", has, address.Has, k))
		}
		num++
		has++
	}
}

func Spider(urls []string) {
	method := "Spider"
	if len(urls) <= 0 {
		return
	}

	logger.Printf("%-4d Spidering ... ... ... url is:   %s\n", count, urls[0])
	count++
	// logger.Println("debug urls is                 ", urls)
	// logger.Println("visit is ******", Visit)

	url := urls[0]
	//添加时候已经判断是不是访问过了
	/*if len(Visit) != 1 && judgeVisit(url) {
		logger.Println(method, "come into already visit.and url is !!!!!!!!!!!!!!!!!!!!!", url)
		Spider(urls[1:])
		return
	}*/

	client := &http.Client{

		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				deadline := time.Now().Add(time.Second * time.Duration(*timeout))
				c, err := net.DialTimeout(netw, addr, time.Second*time.Duration(*timeout))
				if err != nil {
					// fmt.Println("dail timeout", err)
					return nil, err
				}
				c.SetDeadline(deadline)
				return c, nil
			},
		},
	}

	/*req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Println(method, err)
	}*/
	// resp, err := client.Do(req)
	resp, err := client.Get(url)
	if err != nil {
		AppendAddressArr(url, false)
		logger.Println(method, err)
		Spider(urls[1:])
		return
	}

	body := resp.Body
	defer body.Close()
	bodyByte, err := ioutil.ReadAll(body)
	if err != nil {
		logger.Println(method, err)
	}

	// logger.Println("resp is ", string(bodyByte))
	has, matchUrls := MatchTargetUrl(url, string(bodyByte), urls)
	AppendAddressArr(url, has)

	// logger.Println(method, "debug @@@@@@@@@@@ this add array is ", len(matchUrls))
	urls = append(urls, matchUrls...)
	// DebugArray(urls)

	resp.Body.Close()
	Spider(urls[1:])
}

func judgeVisit(url string) bool {
	if _, ok := Visit[url]; ok {
		return true
	}
	return false
}

func MatchTargetUrl(parentUrl, body string, parentUrls []string) (bool, []string) {
	method := "MatchTargetUrl"
	_ = method
	//<\s*(a|img)\s*(href|src)="(http://)?[0-9a-zA-Z/_\.#-]*[^/"]
	urlMatch := regexp.MustCompile(`<\s*(a)\s*(href)="(http://|https://)?[0-9a-zA-Z/_\.#-]*[^/"]`)
	urls := urlMatch.FindAllString(body, -1)

	geturlMatch := regexp.MustCompile(`href="(http://|https://)?[0-9a-zA-Z/_\.#-]*`)
	parseParentUrl := regexp.MustCompile(`(http://|https://)?[0-9a-zA-Z_\.#-]*[^/]`)
	parentUrlHead := parseParentUrl.FindAllString(parentUrl, 1)
	lenParentUrl := len(parentUrlHead)

	for i := 0; i < len(urls); i++ {
		href := geturlMatch.FindAllString(urls[i], 1)
		if len(href) == 0 || len(href[0]) <= 7 {
			// logger.Println(method, "len(href) is ", href, "and url is ", urls[i])
			urls[i] = ""
			continue
		}
		url := href[0][6:]
		if strings.HasPrefix(url, "/") {
			if lenParentUrl <= 0 {
				url = ""
			} else {
				url = parentUrlHead[0] + url
			}
		}
		if strings.HasPrefix(url, "#") {
			url = ""
		}
		url = strings.TrimSpace(url)
		urls[i] = strings.TrimSuffix(url, "/")
	}

	trim := TrimSameSpaceHave(urls, parentUrls)

	// logger.Println(method, "parse urls is ", urls)

	targetReg := ""
	for _, str := range target {
		targetReg += str + "|"
	}
	targetReg = targetReg[:len(targetReg)-1]
	targetMatch := regexp.MustCompile(targetReg)
	target := targetMatch.FindAllString(body, 1)

	if len(target) <= 0 {
		return false, trim
	}
	return true, trim
}

func AppendAddressArr(url string, has bool) {
	Visit[url] = Address{Url: url, Has: has}
}

func TrimSameSpaceHave(urls, parentUrls []string) []string {
	method := "TrimSameSpaceHave"
	_ = method
	ret := make([]string, 0)
	for i := 0; i < len(urls); i++ {
		if urls[i] == "" {
			continue
		}
		have := false
		for _, r := range ret {
			if r == urls[i] {
				have = true
				break
			}
		}
		// logger.Println(method, "judge before url :", urls[i])
		// logger.Println(method, " url is :", urls[i])
		// fmt.Println(method, "visit is ", Visit)
		if have == false && judgeVisit(urls[i]) {
			// logger.Println(method, "judeg in have :")
			have = true
		}
		if have == false {
			for _, r := range parentUrls {
				if r == urls[i] {
					have = true
					break
				}
			}
		}

		if have {
			continue
		} else {
			ret = append(ret, urls[i])
		}
	}
	return ret
}

func DebugArray(arr []string) {
	method := "DebugArray"
	for i, val := range arr {
		logger.Println(method, i, val)
	}
}

func helpFunc() {
	fmt.Printf(`使用说明:\n -urls 和-targets输入要爬的url和关键词，
		目前只支持这个，整个过程可以随时ctrl+c掉，并且结果回存到当
		前目录下的has.txt和ret.txt，主要是保存符合条件的和所有爬多的url。`)

	helpStr := `
-help      ==>   see help detail.
-file      ==>   save result into file.eg: -file=test.txt
-timeout   ==>   dail tcp connect timeout.default is 2s. eg:-timeout=3
-urls      ==>   give spider urls,use ',' split the urls. eg:-urls=www.baidu.com,www.xyz.com
-targets   ==>   give spider targets,use ',' split the target. eg:-targets=golang,docker `

	fmt.Println(helpStr)
}
