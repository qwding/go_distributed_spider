package util

import (
	"bytes"
	"github.com/Sirupsen/logrus"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

// : -> $  ,  / -> (   x   \ -> )  ,   ? -> #   ,   * -> &  ,  " -> !   ,    > ->  @  ,  < -> %   ,   | -> +,
var symple = []string{":", "$$", "/", "$(", "?", "$#", "\\", "$)", "*", "$&", "\"", "$!", ">", "$@", "<", "$%", "|", "$+"}
var DailTimeout = time.Second * time.Duration(2)

func Download(path, name string, url string) error {
	method := "Down"

	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				deadline := time.Now().Add(DailTimeout)
				c, err := net.DialTimeout(netw, addr, DailTimeout)
				if err != nil {
					// fmt.Println("dail timeout", err)
					return nil, err
				}
				c.SetDeadline(deadline)
				return c, nil
			},
		},
	}

	resp, err := client.Get(url)
	if err != nil {
		logrus.Errorln(method, "http.Get", err)
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorln(method, "ioutil.ReadAll", err)
		return err
	}

	err = SaveFile(path, name, body)
	if err != nil {
		logrus.Errorln(method, err)
		return err
	}

	return nil
}

func ReConvertName(name string) (url string) {
	arr := ReConverSympleArr()
	replacer := strings.NewReplacer(arr...)
	name = replacer.Replace(name)
	return name
}

func ReConverSympleArr() []string {
	res := make([]string, len(symple))
	for i := 0; i < len(symple); i = i + 2 {
		if i+1 >= len(symple) {
			logrus.Errorln("symple arr error")
			break
		}
		res[i] = symple[i+1]
		res[i+1] = symple[i]
	}
	return res
}

func SaveFile(path, name string, body []byte) error {
	method := "util SaveFile"
	//windows not support the file name.
	// : -> $  ,  / -> (   x   \ -> )  ,   ? -> #   ,   * -> &  ,  " -> !   ,    > ->  @  ,  < -> %   ,   | -> +,

	replacer := strings.NewReplacer(symple...)
	name = replacer.Replace(name)

	out, err := os.Create(path + "/" + name)
	if err != nil {
		logrus.Errorln(method, "os.Create", err)
		return err
	}
	io.Copy(out, bytes.NewReader(body))
	return nil
}
