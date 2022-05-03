package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/Fe4p3b/gophermart/internal/api/middleware"
	"github.com/Fe4p3b/gophermart/internal/api/model"
	service "github.com/Fe4p3b/gophermart/internal/service/order"
	"github.com/Fe4p3b/gophermart/pkg/luhn"
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

func (o *order) SetupRouting(r *chi.Mux, m middleware.Middleware) {
	r.Get("/api/user/orders", m.Middleware(o.getOrders))
	r.Post("/api/user/orders", m.Middleware(o.addBonus))
}

func (o *order) getOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user, ok := r.Context().Value(middleware.Key).(string)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	orders, err := o.s.List(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(orders) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	jsonOrders := make([]model.Order, 0)

	for _, v := range orders {
		jsonOrders = append(jsonOrders, model.Order{Number: v.Number, Status: v.Status.String(), Accrual: v.Accrual, UploadDate: v.UploadDate.Format(time.RFC3339)})
	}

	b, err := json.Marshal(jsonOrders)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (o *order) addBonus(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.Key).(string)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if isDigitsOnly := luhn.OnlyDigits(b); !isDigitsOnly {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if isLuhnValid := luhn.Luhn(b); !isLuhnValid {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	if err := o.s.AddAccrual(user, string(b)); err != nil {
		if errors.Is(err, service.ErrOrderForUserExists) {
			w.WriteHeader(http.StatusOK)
			return
		} else if errors.Is(err, service.ErrOrderExistsForAnotherUser) {
			w.WriteHeader(http.StatusConflict)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
