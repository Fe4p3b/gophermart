package handler

import (
	"net/http"

	service "github.com/Fe4p3b/gophermart/internal/service/order"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

var _ Handler = (*order)(nil)
var _ Order = (*order)(nil)

type Order interface {
	getOrders(w http.ResponseWriter, r *http.Request)
	addBonus(w http.ResponseWriter, r *http.Request)
}

type order struct {
	l *zap.SugaredLogger
	s service.Order
}

func NewOrder(l *zap.SugaredLogger, s service.Order) *order {
	return &order{l: l, s: s}
}

func (o *order) SetupRouting(r *chi.Mux) {
	r.Get("/api/user/orders", o.getOrders)
	r.Post("/api/user/orders", o.addBonus)
}

func (o *order) getOrders(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("get"))
}

func (o *order) addBonus(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("load"))
}
