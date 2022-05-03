package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Fe4p3b/gophermart/internal/api/middleware"
	"github.com/Fe4p3b/gophermart/internal/api/model"
	service "github.com/Fe4p3b/gophermart/internal/service/balance"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

var _ Handler = (*balance)(nil)
var _ Balance = (*balance)(nil)

type Balance interface {
	Get(w http.ResponseWriter, r *http.Request)
}

type balance struct {
	l *zap.SugaredLogger
	s service.BalanceService
}

func NewBalance(l *zap.SugaredLogger, s service.BalanceService) *balance {
	return &balance{l: l, s: s}
}

func (b *balance) SetupRouting(r *chi.Mux, m middleware.Middleware) {
	r.Get("/api/user/balance", m.Middleware(b.Get))
}

func (b *balance) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user, ok := r.Context().Value(middleware.Key).(string)
	if !ok {
		b.l.Error("error getting user uuid from context")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	balance, err := b.s.Get(user)
	if err != nil {
		b.l.Errorw("error getting user balance", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	bb := model.Balance{Current: float64(balance.Current) / 100, Withdrawn: float64(balance.Withdrawn) / 100}

	resp, err := json.Marshal(bb)
	if err != nil {
		b.l.Errorw("error marshalling response", "error", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
