package main

import (
	"fmt"
	"log"
	"os"

	"github.com/tunnel/tbrurl"
)

/*func main() {
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
		if len(os.Args) > 2 {
			return os.Args[2]
		}
		return "video"
	}()
	log.Printf("%s\t%s download start", username, media)

	err := tbrurl.TbrDownLoader(username)
	if err != nil {
		log.Printf("%s:%s", username, err.Error())
	}
}*/

func main() {
	if len(os.Args) < 2 {
		fmt.Print("Description:will download all the videos of the specify user of tumblr\nExample: tbr username1 ...\nNotice:winows set proxy by \"set http_proxy=127.0.0.1:1080\" and \"set https_proxy=127.0.0.1:1080\"\n")
		return
	}
	log.Printf("Total %d tasks", len(os.Args)-1)
	for i, name := range os.Args {
		if 0 == i {
			continue
		}
		log.Printf("[%s] download start", name)
		err := tbrurl.TbrDownLoader(name)
		if err != nil {
			log.Printf("%s:%s", name, err.Error())
		}
	}
}
