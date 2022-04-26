package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

var _ Handlers = &balance{}
var _ Balance = &balance{}

type Balance interface {
	get(w http.ResponseWriter, r *http.Request)
	withdraw(w http.ResponseWriter, r *http.Request)
	getWithdrawals(w http.ResponseWriter, r *http.Request)
}

type balance struct {
	l *zap.SugaredLogger
}

func NewBalance(l *zap.SugaredLogger) *balance {
	return &balance{l: l}
}

func (b *balance) SetupRouting(r *chi.Mux) {
	r.Get("/api/user/balance", b.get)
	r.Post("/api/user/balance/withdraw", b.withdraw)
	r.Get("/api/user/balance/withdrawals", b.getWithdrawals)
}

func (b *balance) get(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("get"))
}

func (b *balance) withdraw(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("withdraw"))
}

func (b *balance) getWithdrawals(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("getWithdrawals"))
}
