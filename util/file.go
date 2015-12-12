package util

import (
	"os"
)

func IsExist(file string) bool {
	_, err := os.Stat(file)
	if err != nil && os.IsExist(err) == false {
		return false
	}
	return true
}
