package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

/*
q1. Create a handler which accepts any kind of json and prints it

	check if client is still connected or not
	if client is connected then return json processed otherwise just move on
*/
func main() {
	http.HandleFunc("/print", print)
	// http.HandlerFunc
	http.ListenAndServe(":8084", nil)
}

func print(w http.ResponseWriter, r *http.Request) {
	//process the incoming the request here
	//handle the error in the request
	s, err := decode(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//prepare response for sending back to the client
	data, err := json.Marshal(s)
	// err = json.NewEncoder(w).Encode(s)
	if err != nil {
		http.Error(w, "Error encoding data", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(data)
	if err != nil {
		http.Error(w, "Error writing data", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func decode(r *http.Request) (map[string]any, error) {
	var s map[string]any
	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		return nil, fmt.Errorf("error decoding data: %w", err)
	}
	return s, nil
}
