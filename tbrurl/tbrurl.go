package tbrurl

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/tunnel/urlget"
)

//Get public func, return target url
func Get() (string, error) {
	var name string
	var slaiceItems []tItem

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

	jsonparser.ArrayEach([]byte(body), func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		target, _ := jsonparser.GetString(value, "body")
		reg, _ := regexp.Compile(`tumblr_([0-9a-zA-Z]{17}).mp4`)
		url := reg.FindString(target)
		if url == "" {
			target, _ = jsonparser.GetString(value, "video_url")
			reg, _ = regexp.Compile(`tumblr_([0-9a-zA-Z]{17}).mp4`)
			url = reg.FindString(target)
		}
		var summary string
		summarys, _ := jsonparser.GetString(value, "summary")
		reg, _ = regexp.Compile(`\n`)
		if reg.MatchString(summarys) == true {
			titles := strings.Split(summarys, "\n")
			for i, title := range titles {
				if title != "" {
					summary = title
					break
				}
				if i == len(titles) {
					fmt.Print("[ERROR] no title\n")
				}
			}
		} else {
			summary = summarys
		}
		date, _ := jsonparser.GetString(value, "date")

		slaiceItems = append(slaiceItems, tItem{summary, date, url})
	}, "response", "posts")

	for i, obj := range slaiceItems {
		fmt.Printf("%d. %s\t%s\n", i, func(o tItem) string {
			if o.summary == "" {
				return o.date
			}
			return o.summary
		}(obj), obj.item)
	}
	var sIndex int
	fmt.Print("Select a number and enter: ")
	fmt.Scanln(&sIndex)

	return itemPrefix + slaiceItems[sIndex].item, nil
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

//TransHandle response to /tbr/ request
func TransHandle(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	args := strings.Split(req.URL.Path, "/")

	log.Print(func(list []string) string {
		res := "tbrurl ask: "
		for i, s := range list {
			res += fmt.Sprintf("[%d]%s", i, s)
		}
		return res
	}(args))

	if len(args) > 2 {
		rs, err := GetItems(args[2])
		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}
		fmt.Fprintf(w, rs)
		return
	}
	http.NotFound(w, req)
}

//DownLoadHandle reponse to /tbrget/ request
func DownLoadHandle(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	args := strings.Split(req.URL.Path, "/")

	log.Print(func(list []string) string {
		res := "tbrurl download: "
		for i, s := range list {
			res += fmt.Sprintf("[%d]%s", i, s)
		}
		return res
	}(args))

	if req.Method == "Get" {
		fmt.Print("GET:")
		if len(args) > 2 {
			url := itemPrefix + args[2]
			newTask := urlget.NewTask(url)
			newTask.Relay(w)
			return
		}
	} else if req.Method == "POST" {
		fmt.Print("POST:")
		body, _ := ioutil.ReadAll(req.Body)
		req.Body.Close()
		fmt.Print("Post: ", string(body))
		var items []string
		jsonparser.ArrayEach([]byte(body), func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			s, _ := jsonparser.GetString(value, "item")
			items = append(items, s)
		}, "items")
		fmt.Print(func() string {
			res := "Items is:"
			for i, s := range items {
				s += fmt.Sprintf("\n\t[%2d]%s", i, s)
			}
			return res
		}())
		index, _ := jsonparser.GetInt(body, "selected")
		if index < 0 || index > int64(len(items)) {
			fmt.Fprintf(w, "[ERROR] Invalid index")
		}
		newTask := urlget.NewTask(items[index])
		newTask.Relay(w)
		return
	}
	http.NotFound(w, req)
}

//GetItems arg usrname, return tables of item
func GetItems(usrname string) (string, error) {
	var slaiceItems []tItem
	var result string

	tGeter = newGeter(usrname)
	resp, err := http.Get(tGeter.url())
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}

	jsonparser.ArrayEach([]byte(body), func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		target, _ := jsonparser.GetString(value, "body")
		reg, _ := regexp.Compile(`tumblr_([0-9a-zA-Z]{17}).mp4`)
		url := reg.FindString(target)
		if url == "" {
			target, _ = jsonparser.GetString(value, "video_url")
			reg, _ = regexp.Compile(`tumblr_([0-9a-zA-Z]{17}).mp4`)
			url = reg.FindString(target)
		}
		var summary string
		summarys, _ := jsonparser.GetString(value, "summary")
		reg, _ = regexp.Compile(`\n`)
		if reg.MatchString(summarys) == true {
			titles := strings.Split(summarys, "\n")
			for i, title := range titles {
				if title != "" {
					summary = title
					break
				}
				if i == len(titles) {
					fmt.Print("[ERROR] no title\n")
				}
			}
		} else {
			summary = summarys
		}
		date, _ := jsonparser.GetString(value, "date")

		slaiceItems = append(slaiceItems, tItem{summary, date, url})
	}, "response", "posts")

	result = "{\"items\":["
	for i, obj := range slaiceItems {
		result += fmt.Sprintf("{\"index\":%2d,\"item\":\"%s\",\"summary\":\"%s\"}", i, obj.item, func() string {
			if obj.summary == "" {
				return obj.date
			}
			return obj.summary
		}())
		if i != len(slaiceItems)-1 {
			result += ","
		}
	}
	result += "]}"
	return result, nil
}

//private
type tItem struct {
	summary string
	date    string
	item    string
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
	geter.key = "?api_key=takRkZUgF7x3h5Dh296ZDZt3jkaFdILFsBLYBLG9M1pwSArUOe" //secret ByrLhtHqjefHa0T7pqvsowTkMcsVwmOTTdGXIiOuxOeuf11nNM //"?api_key=5iYundnZV0CW2fIdBafMhShEWx0mOep8SFVKXmCUi8oEAqABSZ"
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
