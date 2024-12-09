package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
)

type response struct {
	url  string
	resp *http.Response
	err  error
}

func main() {
	urls := []string{
		"https://www.google.com/",
		"https://stackoverflow.com/",
		"https://go.dev/",
		"https://jhgffsda.com/",
		"https://www.facebook.com/"}
	doGetRequest(urls)
	fmt.Println("Done")
}

func doGetRequest(urls []string) {
	var wg sync.WaitGroup
	httpChannel := make(chan response)

	wg.Add(1)
	go func() {
		defer wg.Done()
		var workerGroup sync.WaitGroup
		for _, url := range urls {
			//fan out pattern. Create goroutine for each of the http network call
			workerGroup.Add(1)
			go func() {
				defer workerGroup.Done()
				resp, err := http.Get(url)
				httpResponse := response{url: url, resp: resp, err: err}
				httpChannel <- httpResponse
			}()
		}
		workerGroup.Wait()
		close(httpChannel)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for v := range httpChannel {
			resp, err := v.resp, v.err
			if err != nil {
				fmt.Println("Error found:", err)
				continue
			}
			_, err = io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Error found:", err)
				continue
			}
			if resp.StatusCode > 299 {
				fmt.Println("Error in response code: ", resp.StatusCode)
				continue
			}
			fmt.Println("Successful response: ", v.url, v.err)
		}
	}()

	wg.Wait()
}
