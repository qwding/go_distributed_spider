package plugins

import (
	"go_distributed_spider/config"
)

func GetPlugins(conf *config.Config) []Plugins {
	arr := []Plugins{}
	if conf.Controller.IfPicture {
		picture := NewPicture(conf)
		arr = append(arr, picture)
	}
	return arr
}
