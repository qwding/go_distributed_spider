package main

/**
 * @author carlding
 * @date : 2015-10-10
 */

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	// "os/signal"
	"go_sample_spider/util"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"
)

/*
*写的还是比较混乱的。思路就是动态分配process，保证process和cpu核数相同，然后维护一个大的urlList，存储访问过和还没有访问的url。
*每个process会领去自己需要遍历的url 列表，遍历结束将获得的新url合并到urlList.
*如果urlList长度等于了index就退出程序，（一般是跑不完的。。。）
*需要改进：
*尽量少用全局变量，process主动在线程领任务改为分配任务。分配由main函数来分配，process只负责遍历，并返回新url list。这样写会清晰很多。
*因为全局变量多，导致用了不少锁，锁来锁去很容易死锁。虽然这么写对锁有了新认识，但是真的很耽误事。
*
 */

type Address struct {
	Url string
	Has bool
}

var logger *log.Logger

var idxLock *sync.RWMutex
var urlsLock *sync.RWMutex
var procLock *sync.RWMutex
var nowIndex int       //index for which the spider have running to.
var urlList []*Address //save all the urls.
var process int = 0    //number of running process.(max is number of cpu.)

var target []string = []string{"golang", "docker"} //
var findNum int = 0
var findHas int = 0

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
	numCpu := runtime.NumCPU()
	idxLock = new(sync.RWMutex)
	urlsLock = new(sync.RWMutex)
	procLock = new(sync.RWMutex)

	logger.Printf("cup number is %d\n", numCpu)
	runtime.GOMAXPROCS(runtime.NumCPU())
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
	urlList = make([]*Address, 0)
	for _, url := range urls {
		urlList = append(urlList, &Address{Url: url})
	}

	if util.IsExist(*saveFile) {
		err := os.Remove(*saveFile)
		if err != nil {
			logger.Println(" remove savefile error!")
		}
	}
	if util.IsExist(*hasFile) {
		err := os.Remove(*hasFile)
		if err != nil {
			logger.Println(" remove savefile error!")
		}
	}
	//every porcess save the file.
	//listen ctrl + c.and save the urls to file before exit.
	/*	c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		go func() {
			for sig := range c {
				_ = sig
				logger.Println("Stoping by Ctrl + C ...... ")
				OutputRes()
				os.Exit(0)
			}
		}()

		//output result and save into file before exit.
		defer OutputRes()*/

	//go to spider
	for i := 0; i < numCpu; i++ {
		procLock.Lock()
		process++
		procLock.Unlock()
		go Spider()
	}

	//control main over.
	for {
		if nowIndex >= len(urlList) && process == 0 {
			return
		} else {
			if process < numCpu {
				logger.Println("debug ", "now proc is ", process, "len(urlList):", len(urlList), "now index:", nowIndex)

				for i := 0; i < numCpu-process; i++ {
					procLock.Lock()
					process++
					procLock.Unlock()
					go Spider()
				}

			}
		}
		time.Sleep(time.Second * 5)
	}
}

func processSave(start, end int) {
	//save all result
	allFile, err := os.OpenFile(*saveFile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		logger.Println("Open file error,", err)
		return
	}
	defer allFile.Close()

	//save has result
	hasFile, err := os.OpenFile(*hasFile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		logger.Println("Open file error,", err)
		return
	}
	defer hasFile.Close()

	processHas := 0
	for ; start < end; start++ {
		str := fmt.Sprintf("has: %-5t,  url:  %s\n", urlList[start].Has, urlList[start].Url)
		allFile.WriteString(str)
		if urlList[start].Has {
			hasFile.WriteString(fmt.Sprintf("has: %-5t,  url:  %s\n", urlList[start].Has, urlList[start].Url))
			processHas++
		}
	}
	findNum += end - start
	findHas += processHas
}

func Spider() {
	method := "Spider"

	numCpu := runtime.NumCPU()
	var start int
	var end int //don't visit end.like slice

	idxLock.Lock()
	lenUrl := len(urlList)
	logger.Println(method, "in lock process num:", process, "len(urllist):", lenUrl, "nowindex", nowIndex)
	if lenUrl <= nowIndex {
		procLock.Lock()
		process--
		procLock.Unlock()
		idxLock.Unlock() //这里没有解锁，简直被坑到死啊，找了半年的问题
		logger.Println(method, "out lock process num:", process, "len(urllist):", lenUrl, "nowindex", nowIndex)
		return
	}

	//head(start) have,tail(end) not.
	start = nowIndex
	if (lenUrl-nowIndex-1)/numCpu <= 0 {
		end = nowIndex + 1
	} else if (lenUrl-nowIndex-1)/numCpu > 100 {
		end = nowIndex + 100
	} else {
		end = (lenUrl-nowIndex-1)/numCpu + nowIndex
	}
	nowIndex = end

	idxLock.Unlock()

	logger.Println(method, "out lock process num:", process, "len(urllist):", lenUrl, "nowindex", nowIndex)

	threadUrls := make([]string, 0)

	for idx := start; idx < end; idx++ {
		url := urlList[idx].Url

		//for debug. debug if the url has been visited.
		/*for i := 0; i < idx; i++ {
			if urlList[i].Url == url {
				logger.Println(method, " ######### url have visited. len(urlList):", len(urlList), url)
			}
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
		resp, err := client.Get(url)
		if err != nil {
			logger.Println(method, "http err process start:", start, "index:", idx, err)
			continue
		}

		body := resp.Body
		defer body.Close()
		bodyByte, err := ioutil.ReadAll(body)
		if err != nil {
			logger.Println(method, "read resp error process start:", start, "index:", idx, err)
		}

		// logger.Println("resp is ", string(bodyByte))
		has, matchUrls := MatchTargetUrl(url, string(bodyByte), threadUrls)
		urlList[idx].Has = has

		logger.Println(method, "res process start:", start, "index:", idx, "has:", has, "url:", url)
		// logger.Println(method, "debug @@@@@@@@@@@ this add array is ", len(matchUrls))

		threadUrls = append(threadUrls, matchUrls...)
		// DebugArray(urls)

		resp.Body.Close()
	}
	logger.Println(method, " process out for")
	urlsLock.Lock()
	mergeUrls(threadUrls)
	urlsLock.Unlock()

	processSave(start, end)
	procLock.Lock()
	process--
	procLock.Unlock()
	logger.Println(method, " this process over.")
}

func mergeUrls(threadUrls []string) {
	for _, url := range threadUrls {
		if !JudgeUrlListVisit(url) {
			urlList = append(urlList, &Address{Url: url})
		}
	}
}

func JudgeUrlListVisit(url string) bool {
	for _, urlStru := range urlList {
		if url == urlStru.Url {
			return true
		}
	}
	return false
}

func judgeVisit(url string, arrUrl []string) bool {
	for _, j := range arrUrl {
		if url == j {
			return true
		}
	}
	return false
}

func MatchTargetUrl(parentUrl, body string, threadUrls []string) (bool, []string) {
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

	trim := TrimSameSpaceHave(urls, threadUrls)

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

func TrimSameSpaceHave(urls, threadUrls []string) []string {
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
		if have == false && judgeVisit(urls[i], threadUrls) {
			// logger.Println(method, "judeg in have :")
			have = true
		}

		if have {
			continue
		} else {
			ret = append(ret, urls[i])
		}
	}
	return ret
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
