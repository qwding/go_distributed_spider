package distributed

import (
	"go_distributed_spider/config"
)

type Distributed interface {
	Init(*config.Config)
	Done(string)
}
