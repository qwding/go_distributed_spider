package util

import (
	"os"
	"strings"
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

func MakedirPath(path string) error {
	if !IsExist(path) {
		arr := strings.Split(path, "/")
		parent := ""
		if strings.HasPrefix(path, "/") {
			parent = "/"
		}
		for _, val := range arr {
			if !IsExist(parent + val) {
				err := MakeDir(parent + val)
				if err != nil {
					return err
				}
			}
			parent = parent + val + "/"
		}
	}
	return nil
}

func MakeDir(name string) error {
	err := os.Mkdir(name, 0666)
	if err != nil {
		return err
	}
	return nil
}
