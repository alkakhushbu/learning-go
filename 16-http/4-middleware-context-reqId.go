package main

import (
	"context"
	"github.com/google/uuid"
	"log"
	"net/http"
	"time"
)

type ContextKey string

const ReqIdKey ContextKey = "reqId"

func main() {

}

// fetch the requestId and log on terminal reqId: hello username
// return hello username to the client

func HelloHandler(w http.ResponseWriter, r *http.Request) {

}
func ReqIdMid(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id := uuid.NewString()
		ctx := r.Context()                         // fetching the ctx object from the request
		ctx = context.WithValue(ctx, ReqIdKey, id) // creating an updated ctx with a traceId store in it
		r = r.WithContext(ctx)                     // putting context inside the request object
		next(w, r)                                 // calling next thing in the chain

	}
}
func LogMid(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		ctx := r.Context()            // fetching the context from the request
		val := ctx.Value(ReqIdKey)    // fetch the request id for the key
		reqId, ok := val.(ContextKey) // checking if the values exist and of correct type
		if !ok {
			reqId = "unknown"
		}
		log.Printf("reqId, method, url: %s, %s, %s\n", reqId, r.Method, r.URL)
		defer log.Printf("reqId,  duration: %s, %s\n", reqId, time.Since(t))
		next(w, r) // calling next thing in the chain

	}

}
