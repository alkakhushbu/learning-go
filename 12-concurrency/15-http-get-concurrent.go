package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
)

var urls = []string{
	`https://pkg.go.dev/`,
	`https://github.com/`,
	`abc.com/1234`,
}

type Response struct {
	url  string
	resp *http.Response
	err  error
}

func main() {
	doGetRequest(urls)
}

func doGetRequest(urls []string) {
	respChan := make(chan Response, len(urls))
	wgWorker := &sync.WaitGroup{}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, v := range urls {
			wgWorker.Add(1)
			//fanning out go routines // one task = one goroutine
			go func(url string) {
				defer wgWorker.Done()
				resp, err := http.Get(url)
				r := Response{
					url:  url,
					resp: resp,
					err:  err,
				}
				respChan <- r //sending the resp struct to respCh
			}(v)
		}

		//wait for go routines to finish the get request task
		wgWorker.Wait()
		close(respChan)
		// when channel is closed no more send can happen // only recv is possible

	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		//range would stop when channel is closed

		for r := range respChan {
			if r.err != nil {
				log.Println(r.err)
				continue
			}

			// reading the request body, ReadAll accepts a reader and body implements the reader interface
			bytes, err := io.ReadAll(r.resp.Body)
			if err != nil {
				log.Println(err)
				continue
			}

			// anything above 299 is a problem
			if r.resp.StatusCode > 299 {
				log.Printf("Response failed with status code: %d and\nbody: %s\n", r.resp.StatusCode, bytes)
				continue
			}

			// printing relevant results
			fmt.Println(r.url, r.resp.Status)
		}
	}()

	// blocking main goroutine
	wg.Wait()
}
