package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
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
	if len(os.Args) < 2 {
		return
	}
	resp, err := http.Get(strings.Replace(os.Args[1], "https://", "http://", 1))
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
	fmt.Print(string(body))
	return
}
