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
	return &authMiddleware{auth: as}
}

func (a *authMiddleware) Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("token")
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		uuid, err := a.auth.Decrypt(token.Value)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		if err := a.auth.VerifyUser(string(uuid)); err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), Key, string(uuid))
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
