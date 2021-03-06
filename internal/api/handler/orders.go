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
	GetOrders(w http.ResponseWriter, r *http.Request)
	AddOrder(w http.ResponseWriter, r *http.Request)
}

type order struct {
	l *zap.SugaredLogger
	s service.OrderService
}

func NewOrder(l *zap.SugaredLogger, s service.OrderService) *order {
	return &order{l: l, s: s}
}

func (o *order) SetupRouting(r *chi.Mux, m middleware.Middleware) {
	r.Get("/api/user/orders", m.Middleware(o.GetOrders))
	r.Post("/api/user/orders", m.Middleware(o.AddOrder))
}

func (o *order) GetOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user, ok := r.Context().Value(middleware.Key).(string)
	if !ok {
		o.l.Error("error getting user uuid from context")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	orders, err := o.s.List(r.Context(), user)
	if err != nil {
		o.l.Errorw("error listing orders", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(orders) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	jsonOrders := make([]model.Order, 0)

	for _, v := range orders {
		jsonOrders = append(jsonOrders, model.Order{Number: v.Number, Status: v.Status.String(), Accrual: float64(v.Accrual) / 100, UploadDate: v.UploadDate.Format(time.RFC3339)})
	}

	b, err := json.Marshal(jsonOrders)
	if err != nil {
		o.l.Errorw("error marshalling orders", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (o *order) AddOrder(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.Key).(string)
	if !ok {
		o.l.Error("error getting user uuid from context")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		o.l.Errorw("error reading body", "error", err)
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

	if err := o.s.AddOrder(r.Context(), user, string(b)); err != nil {
		if errors.Is(err, service.ErrOrderForUserExists) {
			w.WriteHeader(http.StatusOK)
			return
		} else if errors.Is(err, service.ErrOrderExistsForAnotherUser) {
			w.WriteHeader(http.StatusConflict)
			return
		}

		o.l.Errorw("error adding order", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
