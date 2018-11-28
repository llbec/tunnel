package main

import (
	"fmt"

	"github.com/tunnel/tbrurl"
)

func main() {
	var cmd int
	var s string

	fmt.Println("1. get file\n2. get url\nEnter the number:")
	fmt.Scanln(&cmd)

	if cmd == 1 {
		s, _ = tbrurl.GetFile()
		fmt.Print(s)
	} else if cmd == 2 {
		s, _ = tbrurl.Get()
	} else {
		return
	}
}
