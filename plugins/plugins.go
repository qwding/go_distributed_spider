package plugins

import ()

type Plugins interface {
	Init()
	Done(string, []byte)
}
