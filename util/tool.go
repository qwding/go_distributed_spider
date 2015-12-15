package util

import (
	"bytes"
	"github.com/Sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func DownPicture(path, name string, url string) error {
	method := "DownPicture"
	resp, err := http.Get(url)
	if err != nil {
		logrus.Errorln(method, "http.Get", err)
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorln(method, "ioutil.ReadAll", err)
		return err
	}

	//windows not support the file name.
	replacer := strings.NewReplacer(":", "-", "/", "-", "?", "-", "\\", "-", "*", "-", "\"", "-", ">", "-", "<", "-", "|", "-")
	name = replacer.Replace(name)

	out, err := os.Create(path + "/" + name)
	if err != nil {
		logrus.Errorln(method, "os.Create", err)
		return err
	}
	io.Copy(out, bytes.NewReader(body))
	return nil
}
