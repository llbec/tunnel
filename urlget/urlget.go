package urlget

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const gRangeSize = 1024
const gThreadNum = 4

type tPiece struct {
	posStart int64
	posEnd   int64
	status   int
	data     string
	err      error
}

func (h *tPiece) String() string {
	return fmt.Sprintf("bytes=%d-%d", h.posStart, h.posEnd)
}

func (h *tPiece) Length() int {
	return int(h.posEnd - h.posStart + 1)
}

//TTask is a dscription about download task
type TTask struct {
	url    string
	pieces []tPiece
	state  int64
	file   *os.File
}

//NewTask is TTask's constructor
func NewTask(url string) *TTask {
	task := new(TTask)
	task.url = url
	task.state = -1

	len, err := probe(task.url)
	if err != nil {
		return task
	}
	var n int64
	if len != 0 {
		for i, j := 0, 1; j == 1; i++ {
			var pos int64
			pos, j = func(v int64) (int64, int) {
				if v+1 >= len {
					return len - 1, 0
				}
				return v, 1
			}(int64(i*gRangeSize + gRangeSize - 1))
			task.pieces = append(task.pieces, tPiece{int64(i * gRangeSize), pos, 0, "", nil})
			n++
		}
	}

	task.file, err = os.Create(parseFileName(task.url))
	if err != nil {
		log.Fatal(err)
		return task
	}

	task.state = n
	return task
}

//Run start the task
func (task *TTask) Run() {
	if task.state == -1 {
		log.Fatal("Task not ready!")
		return
	}
	if len(task.pieces) == 0 {
		_, err := task.direcDownload()
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	/*for i := 0; i < len(task.pieces); i++ {
		for task.pieces[i].status != 1 {
			task.partialDownload(i)
		}
		n, err := task.file.WriteAt([]byte(task.pieces[i].data), task.pieces[i].posStart)
		if n != task.pieces[i].Length() || err != nil {
			log.Fatal(n, "\t", err)
			return
		}
	}*/

	//multi threads
	thchannel := make(chan int, gThreadNum)
	for i := 0; i < len(task.pieces); i++ {
		if task.pieces[i].status != 0 {
			continue
		}
		go func(pos int) {
			task.pieces[pos].status = -1
			task.partialDownload(pos)
			thchannel <- pos
		}(i)
	}
	log.Printf("Total piece:%d started", len(task.pieces))
	for len(task.pieces) != 0 {
		pos := <-thchannel
		if task.pieces[pos].status == 1 {
			n, err := task.file.WriteAt([]byte(task.pieces[pos].data), task.pieces[pos].posStart)
			if n != task.pieces[pos].Length() || err != nil {
				log.Fatal(n, "\t", err)
				return
			}
			task.pieces = append(task.pieces[:pos], task.pieces[pos+1:]...)
			log.Printf("Piece number %d is complete", pos)
		}
		for i := 0; i < len(task.pieces); {
			if task.pieces[i].status == 1 {
				n, err := task.file.WriteAt([]byte(task.pieces[i].data), task.pieces[i].posStart)
				if n != task.pieces[i].Length() || err != nil {
					log.Fatal(n, "\t", err)
					close(thchannel)
					return
				}
				task.pieces = append(task.pieces[:i], task.pieces[i+1:]...)
				continue
			}
			if task.pieces[i].status != 0 {
				i++
				continue
			}
			go func(pos int) {
				task.pieces[pos].status = -1
				task.partialDownload(pos)
				thchannel <- pos
			}(i)
			i++
		}
	}
}

func parseFileName(url string) string {
	tokens := strings.Split(url, "/")
	fileName := tokens[len(tokens)-1]

	if fileName == "" {
		fileName = "index.html"
	}
	return fileName
}

func (task *TTask) direcDownload() (int64, error) {
	resp, err := http.Get(task.url)
	if err != nil {
		return 0, err
	}

	n, err := io.Copy(task.file, resp.Body)
	return n, err
}

func (task *TTask) partialDownload(pos int) {
	// make HTTP Range request to get file from server
	req, err := http.NewRequest(http.MethodGet, task.url, nil)
	if err != nil {
		task.pieces[pos].err = err
		task.pieces[pos].status = 0
		return
	}
	req.Header.Set("Range", task.pieces[pos].String())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		task.pieces[pos].err = err
		task.pieces[pos].status = 0
		//thchannel <- -1
		return
	}
	defer resp.Body.Close()

	// read data from response and write it
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		task.pieces[pos].err = err
		task.pieces[pos].status = 0
		//thchannel <- -1
		return
	}
	task.pieces[pos].data = string(data[:])
	task.pieces[pos].status = 1
	//thchannel <- pos
}

// probe makes am HTTP request to the site and return site infomation.
// If site is not reachable, return non-nil error.
// If site supports for range request, return the file length (should be greater than 0).
func probe(url string) (length int64, err error) {
	// Check whether site is reachable
	req, err := http.NewRequest(http.MethodHead, url, nil)
	if err != nil {
		log.Printf("Cannot create http request with the URL: %s, error: %v", url, err)
		return
	}

	// Do HTTP HEAD request with range header to this site
	client := &http.Client{
		Timeout: time.Second * 5,
	}
	req.Header.Set("Range", "bytes=0-")
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Remote site is not reachable: %s, error: %v", url, err)
		return
	}
	defer resp.Body.Close()

	// Collect site infomation
	switch resp.StatusCode {
	case http.StatusPartialContent:
		log.Println("Break-point is supported in this downloading task.")

		attr := resp.Header.Get("Content-Length")
		length, err = strconv.ParseInt(attr, 10, 0)
		if err != nil {
			log.Fatal(err)
		}
	case http.StatusOK, http.StatusRequestedRangeNotSatisfiable:
		log.Println(url, "does not support for range request.")
		// set length to N/A or unknown
		length = 0
	default:
		log.Fatal("Got unexpected status code", resp.StatusCode)
		err = errors.New("Unexpected error response when access site: " + url)
	}

	return
}
