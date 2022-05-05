package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func ContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("context ContextMiddleware")
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
