package main

import (
	"log"
	"net/http"
	"time"
)

/*
q2. Create a middleware which logs request url, http request method and time to complete the request
    Hint: time package can help
*/

func main() {
	http.HandleFunc("/home", Log(home))
	http.ListenAndServe(":8084", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	// this should be returned to the client
	// no ResponseWriter method should be called after the write
	_, err := w.Write([]byte("Response Received"))
	if err != nil {
		http.Error(w, "Cannot write response to response writer", http.StatusBadRequest)
		return
	}
	time.Sleep(time.Second * 1)
}

func Log(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t := time.Now().UTC()
		log.Println("Logging started", t)
		log.Println(r.URL.Path, r.Method)
		fn(w, r)
		diff := time.Since(t)
		log.Println("Logging ended", diff)
	}
}
