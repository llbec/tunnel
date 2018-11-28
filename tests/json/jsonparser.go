package main

import (
	"fmt"

	"github.com/buger/jsonparser"
)

func main() {
	data := []byte(`{
  "person": {
    "name":{
      "first": "Leonid",
      "last": "Bugaev",
      "fullName": "Leonid Bugaev"
    },
    "github": {
      "handle": "buger",
      "followers": 109
    },
    "avatars": [
			{ "url": "https://avatars1.githubusercontent.com/u/14009?v=3&s=460", "type": "thumbnail" },
			{ "url": "https://avatars2.githubusercontent.com/u/14009?v=4&s=560", "type": "testnail" }
    ]
  },
  "company": {
    "name": "Acme"
  }
}`)

	result, err := jsonparser.GetString(data, "person", "name", "fullName")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("result is: ", result)

	content, valueType, offset, err := jsonparser.Get(data, "person", "name", "fullName")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("content is: ", content, "valueType is: ", valueType, "offset is: ", offset)
	//jsonparser提供了解析bool、string、float64以及int64类型的方法，至于其他类型，我们可以通过valueType类型来自己进行转化
	result1, err := jsonparser.ParseString(content)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("content ParseString is: ", result1)

	fmt.Println("Parse person name:")
	err = jsonparser.ObjectEach(data, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		fmt.Printf("\tkey:%s\tvalue:%s\tType:%s\n", string(key), string(value), dataType)
		return nil
	}, "person", "name")

	fmt.Println("Parse person github:")
	err = jsonparser.ObjectEach(data, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		fmt.Printf("\tkey:%s\tvalue:%s\tType:%s\n", string(key), string(value), dataType)
		return nil
	}, "person", "github")

	fmt.Println("Parse person avatars:")
	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		//fmt.Println(jsonparser.Get(value, "url"))
		content, valueType, offset, err = jsonparser.Get(value, "url")
		if err != nil {
			fmt.Println(err)
		}
		result1, err := jsonparser.ParseString(content)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("url ParseString is: ", result1)
	}, "person", "avatars")
}
