package main

import (
	"fmt"
	"os"

	"github.com/tunnel/tbrurl"
	"github.com/tunnel/urlget"
)

func main() {
	var cmd int
	var s string

	if len(os.Args) < 2 {
		fmt.Print("Description:will show the videos of the specify user of tumblr,than choose one to download\nExample: tunnel username\nNotice:winows set proxy by \"set http_proxy=127.0.0.1:1080\" and \"set https_proxy=127.0.0.1:1080\"\n")
		return
	}

	fmt.Println("1. get file\n2. download one\nEnter the number:")
	fmt.Scanln(&cmd)

	if cmd == 1 {
		s, _ = tbrurl.GetFile()
	} else if cmd == 2 {
		s, _ = tbrurl.Get()
		if s == "" {
			fmt.Print("URL is NULL\n")
			return
		}
		newTask := urlget.NewTask(s)
		newTask.Run()
	} else {
		return
	}
	fmt.Print("Result is:\n", s, "\n")
}
