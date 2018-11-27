package main

import (
	"fmt"

	"github.com/tunnel/tbrurl"
)

func main() {
	var name string
	fmt.Println("Enter the usrname:")
	fmt.Scanln(&name)
	//geter := tbrget.NewGeter(name)
	fmt.Println(tbrurl.Get(name))
}
