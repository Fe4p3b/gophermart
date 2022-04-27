package handler

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Handler interface {
	SetupRouting(r *chi.Mux)
}

type handler struct {
	l *zap.SugaredLogger
}

func New(l *zap.SugaredLogger) *handler {
	return &handler{l: l}
}

func (h *handler) SetupRouting(r *chi.Mux, hh ...Handler) {
	for _, v := range hh {
		v.SetupRouting(r)
	}
}
