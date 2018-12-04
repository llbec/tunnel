package main

import (
	"fmt"
	"log"
	"os"

	"github.com/tunnel/tbrurl"
)

func main() {
	var (
		username string
		media    string
	)
	if len(os.Args) < 2 {
		fmt.Print("Example: tbr username [video|photo]")
		return
	}
	username = os.Args[1]
	media = func() string {
		/*if len(os.Args) > 2 {
			return os.Args[2]
		}*/
		return "video"
	}()
	log.Printf("%s\t%s download start", username, media)

	err := tbrurl.TbrDownLoader(username)
	if err != nil {
		log.Printf("%s:%s", username, err.Error())
	}
}
