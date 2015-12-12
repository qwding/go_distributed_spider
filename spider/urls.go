package spider

import (
// "github.com/Sirupsen/logrus"
)

type urls []string

func newUrls(start []string) urls {
	return urls(start)
}

func (u urls) inList(url string) bool {
	// for _, val := range []string(u) {
	for _, val := range u {
		if val == url {
			return true
		}
	}
	return false
}

func (u *urls) add(url string) {
	if !u.inList(url) {
		*u = append(*u, url)
		// u = urls(tmp)
	}
}

func (u *urls) addList(urls []string) {
	for _, url := range urls {
		u.add(url)
	}
}

func (u urls) len() int {

	return len([]string(u))
}

func (u urls) GiveUrls(index int) (int, []string) {
	length := u.len()
	end := 0
	if (length-index-1)/Share <= 0 {
		end = index + 1 //this may out of range
	} else if (length-index-1)/Share > 100 {
		end = index + 100
	} else {
		end = (length-index-1)/Share + index
	}
	tmp := []string(u)

	return end, tmp[index:end]
}
