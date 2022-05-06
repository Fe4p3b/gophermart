package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Fe4p3b/gophermart/internal/api/middleware"
	"github.com/Fe4p3b/gophermart/internal/api/model"
	service "github.com/Fe4p3b/gophermart/internal/service/auth"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

var _ Handler = (*auth)(nil)
var _ Auth = (*auth)(nil)

type Auth interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
}

type auth struct {
	l *zap.SugaredLogger
	s service.AuthService
}

func NewAuth(l *zap.SugaredLogger, s service.AuthService) *auth {
	return &auth{l: l, s: s}
}

func (a *auth) SetupRouting(r *chi.Mux, _ middleware.Middleware) {
	r.Post("/api/user/register", a.Register)
	r.Post("/api/user/login", a.Login)
}

func (a *auth) Register(w http.ResponseWriter, r *http.Request) {
	var cred model.Credentials
	if err := json.NewDecoder(r.Body).Decode(&cred); err != nil {
		a.l.Errorw("error decoding body", "error", err)
		w.WriteHeader(http.StatusBadRequest)
	}

	token, err := a.s.Register(r.Context(), cred.Login, cred.Password)
	if err != nil {
		if errors.Is(err, service.ErrUserExists) {
			w.WriteHeader(http.StatusConflict)
			return
		}

		a.l.Errorw("error registering user", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{Name: "token", Value: token})
	w.WriteHeader(http.StatusOK)
}

func (a *auth) Login(w http.ResponseWriter, r *http.Request) {
	var cred model.Credentials
	if err := json.NewDecoder(r.Body).Decode(&cred); err != nil {
		a.l.Errorw("error decoding body", "error", err)
		w.WriteHeader(http.StatusBadRequest)
	}

	token, err := a.s.Login(r.Context(), cred.Login, cred.Password)
	if err != nil {
		if errors.Is(err, service.ErrWrongCredentials) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		a.l.Errorw("error logging in user", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{Name: "token", Value: token})
	w.WriteHeader(http.StatusOK)
}
