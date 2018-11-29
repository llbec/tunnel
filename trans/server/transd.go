package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/tunnel/tbrurl"
)

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()                     //解析参数，默认是不会解析的
	fmt.Println("******\n\t", r.Form) //这些信息是输出到服务器端的打印信息
	fmt.Println("\tpath", r.URL.Path)
	fmt.Println("\tscheme", r.URL.Scheme)
	fmt.Println("\tform url_long", r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("\t\tform key:", k)
		fmt.Println("\t\tform val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Hello astaxie!") //这个写入到w的是输出到客户端的
}

func main() {
	log.Print("Server start ...")
	http.HandleFunc("/", transHandle)        //设置访问的路由
	err := http.ListenAndServe(":9090", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func transHandle(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	args := strings.Split(req.URL.Path, "/")

	log.Print(func(list []string) string {
		var res string
		for i, s := range list {
			res += fmt.Sprintf("[%d]%s", i, s)
		}
		return res
	}(args))

	if len(args) > 2 {
		if args[1] == "tbr" {
			rs, err := tbrurl.GetItems(args[2])
			if err != nil {
				fmt.Fprintf(w, err.Error())
				return
			}
			fmt.Fprintf(w, rs)
			return
		}
		http.NotFound(w, req)
	}
}
