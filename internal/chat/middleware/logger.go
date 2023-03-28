package middleware

import (
	"log"
	"net/http"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		log.Println("Request info", request.URL, request.Header, request.Method)
		next.ServeHTTP(writer, request)
	})
}
