package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

var _ Handlers = &orders{}
var _ Orders = &orders{}

type Orders interface {
	getOrders(w http.ResponseWriter, r *http.Request)
	createBonus(w http.ResponseWriter, r *http.Request)
}

type orders struct {
	l *zap.SugaredLogger
}

func NewOrders(l *zap.SugaredLogger) *orders {
	return &orders{l: l}
}

func (o *orders) SetupRouting(r *chi.Mux) {
	r.Get("/api/user/orders", o.getOrders)
	r.Post("/api/user/orders", o.createBonus)
}

func (o *orders) getOrders(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("get"))
}

func (o *orders) createBonus(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("load"))
}
