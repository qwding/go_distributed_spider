package spider

import (
	"encoding/json"
	// "fmt"
	"github.com/Sirupsen/logrus"
	"go_distributed_spider/config"
	"io/ioutil"
	"net/http"
	"sync"
)

type MasterSpider struct {
	idxLock *sync.RWMutex
	Config  *config.Config                                          //config file
	Urls    urls                                                    //record all the urls
	Index   int                                                     //record in which the urls have been distributed.
	Maps    map[string]func(w http.ResponseWriter, r *http.Request) //just for look beauty.no use.
}

func NewMasterSpider(config *config.Config) *MasterSpider {
	spider := &MasterSpider{Config: config, Urls: newUrls(config.Start), Index: 0, idxLock: new(sync.RWMutex)}
	spider.Maps = map[string]func(w http.ResponseWriter, r *http.Request){
		Url:      spider.Base,
		"/hello": spider.Hello,
	}
	return spider
}

func (s *MasterSpider) Base(w http.ResponseWriter, r *http.Request) {
	method := "MasterSpider Base"
	logrus.Infof("%-15s|%-15s|%-15s|", Url, r.Method, r.RemoteAddr)

	//read request
	body := r.Body
	defer body.Close()
	bodyBytes, err := ioutil.ReadAll(body)
	if err != nil {
		logrus.Errorln(method, err)
		w.Write([]byte(""))
		return
	}
	var slaveToMaster SlaveToMaster
	err = json.Unmarshal(bodyBytes, &slaveToMaster)
	if err != nil {
		logrus.Errorln(method, err)
		w.Write([]byte(""))
		return
	}

	s.Urls.addList(slaveToMaster.NewUrls)

	var masterToSlave MasterToSlave
	if s.Index >= len(s.Urls) {
		masterToSlave.Task = []string{}
	} else {
		s.idxLock.Lock()
		s.Index, masterToSlave.Task = s.Urls.GiveUrls(s.Index)
		s.idxLock.Unlock()
	}

	logrus.Debugln("all urls length:", len(s.Urls), "index:", s.Index, "send task:", len(masterToSlave.Task), "slaveToMaster:", len(slaveToMaster.NewUrls))
	// logrus.Debugln("index", s.Index, "send urls:", masterToSlave.Task)
	// logrus.Debugln("slave give urls:", slaveToMaster.NewUrls)

	tmparr, err := json.Marshal(masterToSlave)
	if err != nil {
		logrus.Errorln(method, err)
	}
	w.Write(tmparr)

}

func (s *MasterSpider) Hello(w http.ResponseWriter, r *http.Request) {
	method := "Hello"
	w.Write([]byte(method + " spider"))
}
