package main

import (
	"fmt"
	"os"

	"github.com/tunnel/tbrurl"
	"github.com/tunnel/urlget"
)

func main() {
	var cmd string
	var s string

	if len(os.Args) < 2 {
		fmt.Print("Description:will show the videos of the specify user of tumblr,than choose one to download\nExample: tunnel username\nNotice:winows set proxy by \"set http_proxy=127.0.0.1:1080\" and \"set https_proxy=127.0.0.1:1080\"\n")
		return
	}

	cmd = os.Args[1]

	if len(os.Args) == 3 && os.Args[2] == "json" {
		s, _ = tbrurl.GetFile(cmd)
	} else {
		s, _ = tbrurl.Get(cmd)
		if s == "" {
			fmt.Print("URL is NULL\n")
			return
		}
		newTask := urlget.NewTask(s)
		newTask.Run()
	}
	fmt.Print("Result is:\n", s, "\n")
}
