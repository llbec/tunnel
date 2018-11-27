package tbrurl

import (
	"io/ioutil"
	"net/http"
)

//Get public func, return target url
func Get(name string) (string, error) {
	tGeter = newGeter(name)
	resp, err := http.Get(tGeter.url())
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

//private
type tbrGet struct {
	prefix  string
	region  string
	usrname string
	method  string
	media   string
	key     string
}

var tGeter tbrGet

func (geter *tbrGet) init() {
	geter.prefix = "https://api.tumblr.com/v2/"
	geter.region = "blog/"
	geter.method = "posts/"
	geter.media = "video"
	geter.key = "?api_key=5iYundnZV0CW2fIdBafMhShEWx0mOep8SFVKXmCUi8oEAqABSZ"
}

func newGeter(name string) tbrGet {
	var geter tbrGet
	geter.init()
	geter.usrname = name + "/"
	return geter
}

func (geter *tbrGet) url() string {
	var url string
	url += geter.prefix
	url += geter.region
	url += geter.usrname
	url += geter.method
	url += geter.media
	url += geter.key
	return url
}
