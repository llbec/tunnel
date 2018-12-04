package main

import (
	"fmt"

	"github.com/tunnel/tbrurl"
	"github.com/tunnel/urlget"
)

func main() {
	var cmd int
	var s string

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
