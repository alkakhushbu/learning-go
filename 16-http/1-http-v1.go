package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// start up the server
	// this lines block forever
	// mux := http.NewServeMux()
	http.HandleFunc("/home", home)
	http.HandleFunc("/find", FetchUser)

	// run the server // it would run, until someone manually kills it
	err := http.ListenAndServe(":8084", nil)
	if err != nil {
		panic(err)
	}

	// mux // mux matches request to handler functions
	// http has a DefaultServeMux mux, which can match request to specific endpoints
	// in ListenAndServe if we pass the handler value as nil, by default it would use http.DefaultServeMux
}

func home(w http.ResponseWriter, r *http.Request) {
	//w http.ResponseWriter, is used to write resp to the client
	// http.Request// anything user send us would be in the request struct
	//w.Write([]byte("Hello World"))
	fmt.Fprintln(w, "Hello World")

}

func FetchUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// fields must be exported , if a struct is sent to be as json
	var user struct {
		Name         string `json:"first_name"` // field level tag
		Password     string `json:"-"`          // - is to ignore the value in json output
		PasswordHash string `json:"password_hash"`
		Marks        []int  `json:"marks"`
	}

	user.Name = "John"
	user.Password = "<PASSWORD>"
	user.PasswordHash = "<HASH>"
	user.Marks = []int{10, 20, 30}

	// NewEncoder can directly write JSON to the writer
	// Encode would convert struct/map to JSON

	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(user)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return // don't forget the return, program would move on
	}

}
