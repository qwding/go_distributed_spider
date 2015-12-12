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

func RemoveFile(file string) error {
	if !IsExist(file) {
		return nil
	} else {
		err := os.Remove(file)
		return err
	}
}
