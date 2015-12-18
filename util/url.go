package util

import (
	"strings"
)

func GetDaemon(url string) string {
	if url == "" {
		return ""
	}
	url = strings.TrimPrefix(url, "http://")
	url = strings.TrimPrefix(url, "https://")

	arr := strings.Split(url, "/")
	url = arr[0]
	dae := strings.Split(url, ".")
	if len(dae) <= 1 {
		return ""
	} else if len(dae) == 2 {
		return dae[0] + "." + dae[1]
	} else {
		return dae[len(dae)-2] + "." + dae[len(dae)-1]
	}
}
