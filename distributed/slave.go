package distributed

import (
	"github.com/Sirupsen/logrus"
	"go_distributed_spider/config"

	"go_distributed_spider/spider"
)

type Slave struct {
	Spider *spider.SlaveSpider
	Config *config.Config
}

func (m *Slave) Init(conf *config.Config) {
	method := "Slave Init"
	logrus.Infoln(method)
	m.Config = conf
	m.Spider = spider.NewSlaveSpider(conf)

	//init plugins.
	for _, plugin := range m.Spider.Plugins {
		plugin.Init()
	}
}

func (m *Slave) Done(str string) {
	method := "Slave Done"
	logrus.Infoln(method)
	m.Spider.Base()
}
