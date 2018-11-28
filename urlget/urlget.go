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

//TTask is a dscription about download task
type TTask struct {
	name   string //filename
	url    string
	pieces []tPiece
	state  int64
	file   *os.File
}

//NewTask is TTask's constructor
func NewTask(url string, name string) (task TTask) {
	task.url = url
	task.name = name
	task.state = -1

	len, err := probe(task.url)
	if err != nil {
		return
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

	task.file, err = os.Create(task.name)
	if err != nil {
		log.Fatal(err)
		return
	}

	task.state = n
	return
}

//Run start the task
func (task *TTask) Run() {
	if len(task.pieces) == 0 {
		_, err := task.direcDownload()
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	//var thchannel chan int = make(chan int, gThreadNum)
	for i := 0; i < len(task.pieces); i++ {
		//thchannel <- i
		go task.partialDownload(i)
	}
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
	task.pieces[pos].status = 1
	task.pieces[pos].data = string(data[:])
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
