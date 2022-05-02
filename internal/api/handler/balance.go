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
	user, ok := r.Context().Value(middleware.Key).(string)
	if !ok {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	balance, err := b.s.Get(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bb := model.Balance{Current: float64(balance.Current) / 100, Withdrawn: balance.Withdrawn}

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
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	var withdrawal model.Withdrawal
	if err := json.NewDecoder(r.Body).Decode(&withdrawal); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := b.s.Withdraw(user, withdrawal.Order, withdrawal.Sum); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (b *balance) getWithdrawals(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.Key).(string)
	if !ok {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	withdrawals, err := b.s.GetWithdrawals(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bb, err := json.Marshal(model.ToWithdrawals(withdrawals))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bb)
}
