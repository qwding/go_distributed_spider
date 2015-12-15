package distributed

import (
	"github.com/Sirupsen/logrus"
	"go_distributed_spider/config"
	"go_distributed_spider/spider"
	"go_distributed_spider/util"
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

	//picture path check.
	if conf.IfPicture {
		err := util.MakedirPath(conf.PicturePath)
		if err != nil {
			logrus.Errorln(method, err)
		}
	}

}

func (m *Slave) Done(str string) {
	method := "Slave Done"
	logrus.Infoln(method)
	m.Spider.Base()
}
