package spider

import (
//"fmt"
)

/*type Body struct {
	Task  []string
	Match []string
}
*/

type SlaveToMaster struct {
	NewUrls []string
	Match   []string
}

type MasterToSlave struct {
	Task []string
}
