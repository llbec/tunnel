package main

import (
	"fmt"

	"github.com/buger/jsonparser"
	"github.com/tunnel/tbrurl"
)

func main() {
	var name string

	fmt.Println("Enter the usrname:")
	fmt.Scanln(&name)
	//geter := tbrget.NewGeter(name)
	//fmt.Println(tbrurl.Get(name))
	s, err := tbrurl.Get(name)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Print(s)
	jsonparser.ArrayEach([]byte(s), func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		fmt.Println(jsonparser.GetString(value, "date"))
	}, "response", "posts")
}
