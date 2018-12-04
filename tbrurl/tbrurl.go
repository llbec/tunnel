package tbrurl

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/tunnel/urlget"
)

//Get public func, return target url
func Get() (string, error) {
	var (
		name        string
		slaiceItems []tItem
		nOffset     int64
		nIndex      int
	)

	fmt.Println("Enter the usrname:")
	fmt.Scanln(&name)

	posts, err := getPosts(name)
	if err != nil {
		log.Printf("%s:%s", name, err.Error())
		return "", err
	}
	log.Printf("%s total %d posts", name, posts)

	showItems := func() {
		slaiceItems = getItemList(name, nOffset)
		fmt.Printf("%s: Total %d medias\n", name, posts)
		for i, obj := range slaiceItems {
			fmt.Printf("%d. %s\t%s\n", int64(i)+nOffset, func(o tItem) string {
				if o.summary == "" {
					return o.date
				}
				return o.summary
			}(obj), obj.item)
		}
	}

	showItems()
	for cmd := ""; cmd != "q"; {
		fmt.Print("q to quit; n to next; p to prev; number to download.")
		fmt.Scanln(&cmd)
		if cmd == "n" {
			if nOffset+20 >= posts {
				fmt.Print("The End!\n")
			} else {
				nOffset += 20
				showItems()
			}
		} else if cmd == "p" {
			if nOffset-20 < 0 {
				fmt.Print("The Begin!\n")
			} else {
				nOffset -= 20
				showItems()
			}
		} else if true == func() bool {
			nIndex, err = strconv.Atoi(cmd)
			if err == nil && int64(nIndex)-nOffset >= 0 && int(int64(nIndex)-nOffset) < len(slaiceItems) {
				return true
			}
			return false
		}() {
			return itemPrefix + slaiceItems[int(int64(nIndex)-nOffset)].item, nil
		}
	}
	return "", nil
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
		res := "tbrurl download: " + req.Method + " "
		for i, s := range list {
			res += fmt.Sprintf("[%d]%s", i, s)
		}
		return res
	}(args))

	if req.Method == "GET" {
		if len(args) > 2 {
			url := itemPrefix + args[2]
			newTask := urlget.NewTask(url)
			newTask.Relay(w)
			return
		}
	} else if req.Method == "POST" {
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
	slaiceItems := getItemList(usrname, 0)
	var result string

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

//TbrDownLoader run a downloader to get all items of a user
func TbrDownLoader(name string) error {
	var nOffset int64

	posts, err := getPosts(name)
	if err != nil {
		log.Printf("%s:%s", name, err.Error())
		return nil
	}
	log.Printf("%s total %d posts", name, posts)

	for nOffset < posts {
		download := func() error {
			n, err := downloadPage(name, nOffset)
			if err != nil {
				return err
			}
			if n != 20 {
				log.Printf("%s offset %d download %d videos", name, nOffset, n)
			}
			return nil
		}
		if download() != nil {
			if download() != nil {
				if err := download(); err != nil {
					log.Printf("%s:%s", name, err.Error())
					return err
				}
			}
		}
		nOffset += 20
		log.Printf("%s finish %d", name, nOffset)
	}

	return nil
}

func getPosts(name string) (int64, error) {
	tGeter := newGeter(name)
	resp, err := http.Get(tGeter.url())
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	return jsonparser.GetInt(body, "response", "blog", "total_posts")
}

func getItemList(name string, offset int64) []tItem {
	var (
		slaiceItems []tItem
	)
	tGeter := newGeter(name)
	resp, err := http.Get(tGeter.pageurl(offset))
	if err != nil {
		log.Printf("%s:%s", name, err.Error())
		return nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("%s:%s", name, err.Error())
		return nil
	}
	jsonparser.ArrayEach([]byte(body), func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		target, _ := jsonparser.GetString(value, "body")
		reg, _ := regexp.Compile(`tumblr_([0-9a-zA-Z]{17}).mp4`)
		url := reg.FindString(target)
		if url == "" {
			target, _ = jsonparser.GetString(value, "video_url")
			url = reg.FindString(target)
			if url == "" {
				url = reg.FindString(string(value))
			}
		}
		var summary string
		summarys, _ := jsonparser.GetString(value, "summary")
		reg, _ = regexp.Compile(`\n`)
		if reg.MatchString(summarys) == true {
			titles := strings.Split(summarys, "\n")
			for _, title := range titles {
				if title != "" {
					summary = title
					break
				}
			}
		} else {
			summary = summarys
		}
		date, _ := jsonparser.GetString(value, "date")

		slaiceItems = append(slaiceItems, tItem{summary, date, url})
	}, "response", "posts")

	return slaiceItems
}

func downloadPage(name string, offset int64) (int, error) {
	var (
		slaiceItems []string
		nCount      int
	)
	tGeter := newGeter(name)
	resp, err := http.Get(tGeter.pageurl(offset))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	jsonparser.ArrayEach([]byte(body), func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		target, _ := jsonparser.GetString(value, "body")
		reg, _ := regexp.Compile(`tumblr_([0-9a-zA-Z]{17}).mp4`)
		url := reg.FindString(target)
		if url == "" {
			target, _ = jsonparser.GetString(value, "video_url")
			//reg, _ = regexp.Compile(`tumblr_([0-9a-zA-Z]{17}).mp4`)
			url = reg.FindString(target)
			if url == "" {
				url = reg.FindString(string(value))
			}
		}
		slaiceItems = append(slaiceItems, url)
	}, "response", "posts")

	for _, item := range slaiceItems {
		newTask := urlget.CreateTask(itemPrefix+item, name)
		if newTask == nil {
			continue
		}
		newTask.Run()
		nCount++
	}
	return nCount, nil
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

//https://api.tumblr.com/v2/blog/username/posts/video?api_key=takRkZUgF7x3h5Dh296ZDZt3jkaFdILFsBLYBLG9M1pwSArUOe
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

//&offset=${(page-1)*20}
func (geter *tbrGet) pageurl(offset int64) string {
	return geter.url() + fmt.Sprintf("&offset=%d", offset)
}
