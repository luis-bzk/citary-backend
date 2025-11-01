package middleware

import (
	"log"
	"net/http"
)

// Recovery middleware recovers from panics and returns a 500 error
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("PANIC RECOVERED: %v", err)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"success":false,"message":"Internal server error","error":{"code":500,"message":"An unexpected error occurred"}}`))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
