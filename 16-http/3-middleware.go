package main

import (
	"fmt"
	"net/http"
)

// middleware that exec some pre-processing or the post-processing logic
// Req -> (Mid) -> HandleFunc -> Service -> User -> Database
//		<-		   <-			<-				 <-            <- return flow
// Middleware Examples: logging, Panic Recovery, Auth, Authorize, GenerateReqID, Fetching Headers

func main() {
	// handle functions handle the request and do request matching
	http.HandleFunc("/home", Mid(homeV2))
	http.ListenAndServe(":8083", nil)
}

// Mid is a middleware that accepts a handler func and returns a handler func
func Mid(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Mid layer started")
		fmt.Println("pre processing logic")
		f(w, r) //handler function would be called here that was passed to mid
		fmt.Println("Mid layer ended")
		fmt.Println("post processing logic")
	}
}

// homeV2 is a handler function.
// handler function is a type defined in standard lib
// type HandlerFunc func(ResponseWriter, *Request)
func homeV2(w http.ResponseWriter, r *http.Request) {
	fmt.Println("homeHandler layer started")
	Service()
	fmt.Println("homeHandler layer ended")
	fmt.Fprintln(w, "home page V2")
}

func Service() {
	fmt.Println("Service layer started")
	UserService()
	fmt.Println("Service layer ended")
}

func UserService() {
	fmt.Println("User layer started")
	fmt.Println("User layer ended")
}
