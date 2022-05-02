package middleware

import (
	"context"
	"net/http"

	"github.com/Fe4p3b/gophermart/internal/service/auth"
	"go.uber.org/zap"
)

type Middleware interface {
	Middleware(http.HandlerFunc) http.HandlerFunc
}

type authMiddleware struct {
	l    *zap.SugaredLogger
	auth auth.AuthService
}
type ContextKey string

var Key ContextKey = "user"

func NewAuthMiddleware(as auth.AuthService) *authMiddleware {
	return &authMiddleware{}
}

func (a *authMiddleware) Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("token")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		uuid := token.Value

		ctx := context.WithValue(r.Context(), Key, uuid)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
