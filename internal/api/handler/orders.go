package handler

import (
	"encoding/json"
	"io"
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
	s service.OrderService
}

func NewOrder(l *zap.SugaredLogger, s service.OrderService) *order {
	return &order{l: l, s: s}
}

func (o *order) SetupRouting(r *chi.Mux) {
	r.Get("/api/user/orders", o.getOrders)
	r.Post("/api/user/orders", o.addBonus)
}

func (o *order) getOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := o.s.List("9deb06e4-59b2-496e-9af4-17809f317e59")
	if err != nil {
		o.l.Errorw("order handler", "error", err)
	}

	b, err := json.Marshal(orders)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		// http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (o *order) addBonus(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := o.s.AddAccrual("9deb06e4-59b2-496e-9af4-17809f317e59", string(b)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
