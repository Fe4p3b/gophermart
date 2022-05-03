package handler

import (
	"encoding/json"
	"errors"
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
	get(w http.ResponseWriter, r *http.Request)
	withdraw(w http.ResponseWriter, r *http.Request)
	getWithdrawals(w http.ResponseWriter, r *http.Request)
}

type balance struct {
	l *zap.SugaredLogger
	s service.BalanceService
}

func NewBalance(l *zap.SugaredLogger, s service.BalanceService) *balance {
	return &balance{l: l, s: s}
}

func (b *balance) SetupRouting(r *chi.Mux, m middleware.Middleware) {
	r.Get("/api/user/balance", m.Middleware(b.get))
	r.Post("/api/user/balance/withdraw", m.Middleware(b.withdraw))
	r.Get("/api/user/balance/withdrawals", m.Middleware(b.getWithdrawals))
}

func (b *balance) get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user, ok := r.Context().Value(middleware.Key).(string)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	balance, err := b.s.Get(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	bb := model.Balance{Current: float64(balance.Current) / 100, Withdrawn: float64(balance.Withdrawn) / 100}

	resp, err := json.Marshal(bb)
	if err != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func (b *balance) withdraw(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.Key).(string)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var withdrawal model.Withdrawal
	if err := json.NewDecoder(r.Body).Decode(&withdrawal); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := b.s.Withdraw(user, withdrawal.Order, withdrawal.Sum); err != nil {
		if errors.Is(err, service.ErrNoOrder) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		if errors.Is(err, service.ErrInsufficientBalance) {
			w.WriteHeader(http.StatusPaymentRequired)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (b *balance) getWithdrawals(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user, ok := r.Context().Value(middleware.Key).(string)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	withdrawals, err := b.s.GetWithdrawals(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(withdrawals) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	bb, err := json.Marshal(model.ToAPIWithdrawals(withdrawals))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bb)
}
