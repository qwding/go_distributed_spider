package main

/**
 * @author carlding
 * @date : 2015-12-8
 */
import (
	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/pkg/mflag"
	"go_distributed_spider/config"
	"go_distributed_spider/distributed"
	"os"
)

var (
	master, help  bool
	env, conffile string
)

func init() {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
	mflag.BoolVar(&master, []string{"m", "-master", "#mt"}, false, "if the master server.")
	mflag.BoolVar(&help, []string{"h", "-help", "#hhhhelp"}, false, "show help.")
	mflag.StringVar(&env, []string{"e", "-env"}, "test", "which config env to run.")
	// mflag.BoolVar(&picture, []string{"i", "-image"}, false, "download the picture that go throuth.")
	mflag.StringVar(&conffile, []string{"f", "-config"}, "", "point a config file.default is config/config.json.")
	mflag.Parse()

}
func main() {
	if help {
		mflag.PrintDefaults()
		return
	}

	config, err := config.ReadConfig(conffile, env)
	if err != nil {
		logrus.Errorln("main", err)
		return
	}
	var distribute distributed.Distributed
	if master {
		distribute = &distributed.Master{}
	} else {
		distribute = &distributed.Slave{}
	}
	distribute.Init(config)
	distribute.Done("str")
}
