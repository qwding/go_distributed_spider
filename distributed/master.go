package distributed

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"go_distributed_spider/config"
	"go_distributed_spider/spider"
	"go_distributed_spider/util"
	"net/http"
)

type Master struct {
	Mux    *http.ServeMux
	Spider *spider.MasterSpider
	Config *config.Config
}

func (m *Master) Init(config *config.Config) {
	method := "Master Init"
	logrus.Infoln(method)
	m.Config = config
	m.Spider = spider.NewMasterSpider(config)
	err := util.RemoveFile(m.Spider.AllF)
	if err != nil {
		logrus.Errorln(method, err)
	}
	err = util.RemoveFile(m.Spider.MatchF)
	if err != nil {
		logrus.Errorln(method, err)
	}
	m.Mux = http.NewServeMux()
	for k, v := range m.Spider.Maps {
		m.Mux.HandleFunc(k, v)
	}
}

func (m *Master) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.Mux.ServeHTTP(w, r)
}

func (m *Master) Done(str string) {
	method := "Master Done"
	logrus.Infoln("server start on ", m.Config.Port)
	fmt.Printf("Done %#v\n", m)
	err := http.ListenAndServe(":"+m.Config.Port, m)
	if err != nil {
		logrus.Errorln(method, err)
	}
}
