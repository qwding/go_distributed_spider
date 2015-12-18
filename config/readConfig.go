package config

import (
	"encoding/json"
	"fmt"
	"go_distributed_spider/util"
	// "github.com/Sirupsen/logrus"
	"io/ioutil"
	"os"
)

type Config struct {
	Master        string //程序的master，slave需要
	Port          string //初始url port
	Scheme        string //初始url scheme
	DailTimeout   int    //连接url超时时间
	AllF          string //遍历过的url存放文件
	MatchFDefault string //遍历match的url存放文件
	Controller    controller
}

type controller struct {
	Start        []string //初始url集合
	Target       []string //目的关键字
	IfPicture    bool     //是否下载图片
	PicturePath  string   //下载图片路径
	IfDownHtml   bool     //是否下载html
	DownHtmlPath string   //下载html路径
	Filter       Filter   //过滤新url
}

type Filter struct {
	Only []string //唯一模式，只遍历only数组里的url开始的路径
	Cut  []string //cut数组里的url不做遍历，前两个不能同时用
}

const (
	defaultConfigPath  string = "config/config.json"
	defaultPicturePath string = "picture"
)

func (c *Config) GetUrl() string {
	return c.Scheme + "://" + c.Master + ":" + c.Port
}

func ReadConfig(conffile, env string) (*Config, error) {
	// method := "ReadConfig"
	if conffile == "" {
		conffile = defaultConfigPath
	}

	if !util.IsExist(conffile) {
		return nil, fmt.Errorf("config file not exit.")
	}
	file, err := os.Open(conffile)
	if err != nil {
		return nil, err
	}
	fileByte, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	var config map[string]*Config
	err = json.Unmarshal(fileByte, &config)
	if err != nil {
		return nil, err
	}
	fmt.Printf("config: %#v\n", config[env])

	return config[env], nil
}
