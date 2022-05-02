package handler

import (
	"github.com/Fe4p3b/gophermart/internal/api/middleware"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Handler interface {
	SetupRouting(*chi.Mux, middleware.Middleware)
}

type handler struct {
	l *zap.SugaredLogger
}

func New(l *zap.SugaredLogger) *handler {
	return &handler{l: l}
}

func (h *handler) SetupRouting(r *chi.Mux, m middleware.Middleware, hh ...Handler) {
	for _, v := range hh {
		v.SetupRouting(r, m)
	}
}
