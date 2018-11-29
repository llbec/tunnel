package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/tunnel/tbrurl"
)

func main() {
	log.Print("Server start ...")
	http.HandleFunc("/", hello) //设置访问的路由
	http.HandleFunc("/tbr/", tbrurl.TransHandle)

	err := http.ListenAndServe(":9090", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func hello(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	fmt.Fprintf(w, "Welcome to 666!")
}
