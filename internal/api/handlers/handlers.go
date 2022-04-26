package handlers

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Handlers interface {
	SetupRouting(r *chi.Mux)
}

type handlers struct {
	l *zap.SugaredLogger
}

func New(l *zap.SugaredLogger) *handlers {
	return &handlers{l: l}
}

func (h *handlers) SetupRouting(r *chi.Mux, hh ...Handlers) {
	for _, v := range hh {
		v.SetupRouting(r)
	}
}
