package main

import (
	"fmt"

	"github.com/tunnel/tbrurl"
	"github.com/tunnel/urlget"
)

func main() {
	var cmd int
	var s string

	fmt.Println("1. get file\n2. get url\nEnter the number:")
	fmt.Scanln(&cmd)

	if cmd == 1 {
		s, _ = tbrurl.GetFile()
	} else if cmd == 2 {
		s, _ = tbrurl.Get()
		newTask := urlget.NewTask(s, "test.mp4")
		newTask.Run()
	} else {
		return
	}
	fmt.Print("Result is:\n", s)
}
