package service

import "net/http"

func TraceId(next http.Handler) http.Handler {
	panic("Not implemented")
}

func Logging(http.Handler) http.Handler {
	panic("Not implemented")
}
