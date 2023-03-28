package middleware

import (
	"net/http"
)

func Adapt(next http.Handler, adapters ...func(next http.Handler) http.Handler) http.Handler {
	for i := range adapters {
		next = adapters[i](next)
	}
	return next
}
