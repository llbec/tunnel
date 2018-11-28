package tbrurl

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/buger/jsonparser"
)

//Get public func, return target url
func Get() (string, error) {
	var name string
	mapURL := make(map[string]string)

	fmt.Println("Enter the usrname:")
	fmt.Scanln(&name)

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

	n := 1
	jsonparser.ArrayEach([]byte(body), func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		target, _ := jsonparser.GetString(value, "body")
		reg, _ := regexp.Compile(`tumblr_([0-9a-zA-Z]{17}).mp4`)
		url := reg.FindString(target)
		if url == "" {
			target, _ = jsonparser.GetString(value, "video_url")
			reg, _ = regexp.Compile(`tumblr_([0-9a-zA-Z]{17}).mp4`)
			url = reg.FindString(target)
		}
		summary, _ := jsonparser.GetString(value, "summary")
		if summary == "" {
			date, _ := jsonparser.GetString(value, "date")
			//fmt.Printf("%d. %s\n", n, date)
			mapURL[date] = url
		} else {
			//fmt.Printf("%d. %s\n", n, summary)
			titles := strings.Split(summary, "\n")
			for i, title := range titles {
				if title != "" {
					mapURL[summary] = url
					break
				}
				if i == len(titles) {
					fmt.Print("[ERROR] no title\n")
				}
			}
		}
	}, "response", "posts")

	for title := range mapURL {
		fmt.Printf("%d. %s\n\t%s\n", n, title, mapURL[title])
		n++
	}

	var sIndex int
	fmt.Print("Select a number and enter: ")
	fmt.Scanln(&sIndex)

	return string(body), nil
}

//GetFile get url file
func GetFile() (string, error) {
	var name string

	fmt.Println("Enter the usrname:")
	fmt.Scanln(&name)

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
type tItem struct {
	title string
	url   string
}

type tbrGet struct {
	prefix  string
	region  string
	usrname string
	method  string
	media   string
	key     string
}

var itemPrefix = "http://ve.media.tumblr.com/"

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
