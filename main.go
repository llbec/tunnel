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
	n := 1
	jsonparser.ArrayEach([]byte(s), func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		summary, err := jsonparser.GetString(data, "summary")
		if err != nil {
			fmt.Println(err)
			continue
		}
		date, err := jsonparser.GetString(data, "date")
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Printf("%d. %s\t%s\n", n, summary, date)
		n++
	}, "response", "posts")
}
