package handler

import (
	"net/http"

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

func (b *balance) SetupRouting(r *chi.Mux) {
	r.Get("/api/user/balance", b.get)
	r.Post("/api/user/balance/withdraw", b.withdraw)
	r.Get("/api/user/balance/withdrawals", b.getWithdrawals)
}

func (b *balance) get(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("get"))
}

func (b *balance) withdraw(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("withdraw"))
}

func (b *balance) getWithdrawals(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("getWithdrawals"))
}
