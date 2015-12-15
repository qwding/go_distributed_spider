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
	Master        string
	IfPicture     bool
	PicturePath   string
	Target        []string
	Port          string
	Start         []string
	Scheme        string
	DailTimeout   int
	AllF          string
	MatchFDefault string
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
	return config[env], nil
}
