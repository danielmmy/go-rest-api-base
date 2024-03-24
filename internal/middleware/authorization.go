package middleware

import (
	"net/http"

	"gorest/internal/tools"
)

func Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token != "123456" {
			tools.ErrorLogger.Printf("unauthorized access attempt from %v\n", r.RemoteAddr)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
