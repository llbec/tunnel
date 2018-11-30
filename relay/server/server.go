package server

import (
	"log"
	"os"
	"strconv"
)

func run() {
	port := 9090
	if len(os.Args) > 1 {
		port, _ = strconv.Atoi(os.Args[1])
	}
	log.Print(port)
}
