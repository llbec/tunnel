package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	filename string
)

func init() {
	flag.StringVar(&filename, "filename", "", "Specify the file to save result html")
	flag.StringVar(&filename, "f", "", "Specify the file to save result html")
}

func main() {
	flag.Parse()
	if len(flag.Args()) < 1 {
		return
	}
	resp, err := http.Get(func(url string) string {
		prefix := "http://"
		url = strings.Replace(url, "https://", prefix, 1)
		if len(url) < len(prefix) || url[0:len(prefix)] != prefix {
			return prefix + url
		}
		return url
	}(flag.Args()[1]))
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	if filename != "" {
		if ioutil.WriteFile(filename+".html", body, 0644) == nil {
			fmt.Printf("Save to file %s.html success\n", filename)
			return
		}
	}
	fmt.Print(string(body))
	return
}
