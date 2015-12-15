package spider

import (
	"encoding/json"
	// "fmt"
	"github.com/Sirupsen/logrus"
	"go_distributed_spider/config"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
)

type MasterSpider struct {
	idxLock *sync.RWMutex
	Config  *config.Config                                          //config file
	Urls    urls                                                    //record all the urls
	Index   int                                                     //record in which the urls have been distributed.
	Maps    map[string]func(w http.ResponseWriter, r *http.Request) //just for look beauty.no use.
	AllF    string
	MatchF  string
}

func NewMasterSpider(conf *config.Config) *MasterSpider {
	spider := &MasterSpider{Config: conf, Urls: newUrls(conf.Start), Index: 0, idxLock: new(sync.RWMutex)}
	spider.Maps = map[string]func(w http.ResponseWriter, r *http.Request){
		Url:      spider.Base,
		"/hello": spider.Hello,
	}
	if conf.AllF == "" {
		spider.AllF = AllFDefault
	}
	if conf.MatchFDefault == "" {
		spider.MatchF = MatchFDefault
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
	var slaveToMaster *SlaveToMaster
	err = json.Unmarshal(bodyBytes, &slaveToMaster)
	if err != nil {
		logrus.Errorln(method, err)
		w.Write([]byte(""))
		return
	}

	idxBefore := len(s.Urls)
	//add the urls which slave given to the all list.
	s.Urls.addList(slaveToMaster.NewUrls)
	idxAfter := len(s.Urls)

	//record all urls and match to file.
	err = s.SaveToFile(slaveToMaster, idxBefore, idxAfter)

	//master send task to slave
	var masterToSlave MasterToSlave
	if s.Index >= len(s.Urls) {
		masterToSlave.Task = []string{}
	} else {
		s.idxLock.Lock()
		s.Index, masterToSlave.Task = s.Urls.GiveUrls(s.Index)
		s.idxLock.Unlock()
	}

	logrus.Debugln("all urls length:", len(s.Urls), "index:", s.Index, "send task:", len(masterToSlave.Task), "slaveToMaster match:", len(slaveToMaster.Match), "slaveToMaster NewUrls:", len(slaveToMaster.NewUrls))
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

func (s *MasterSpider) SaveToFile(slaveToMaster *SlaveToMaster, idxBefore, idxAfter int) error {
	method := "MasterSpider SaveTdoFile"
	allf, err := os.OpenFile(s.AllF, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		logrus.Errorln(method, err)
		return err
	}
	defer allf.Close()

	matchf, err := os.OpenFile(s.MatchF, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		logrus.Errorln(method, err)
		return err
	}
	defer matchf.Close()

	all := ""
	for _, val := range s.Urls[idxBefore:idxAfter] {
		all += val + "\n"
	}
	logrus.Debugln(method, idxBefore, idxAfter)
	_, err = allf.WriteString(all)
	if err != nil {
		return err
	}

	match := ""
	for _, val := range slaveToMaster.Match {
		match += val + "\n"
	}

	_, err = matchf.WriteString(match)
	if err != nil {
		return err
	}
	return nil
}
