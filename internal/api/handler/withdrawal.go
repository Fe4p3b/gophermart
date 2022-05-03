package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Fe4p3b/gophermart/internal/api/middleware"
	"github.com/Fe4p3b/gophermart/internal/api/model"
	service "github.com/Fe4p3b/gophermart/internal/service/withdrawal"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

var _ Handler = (*withdrawal)(nil)
var _ Withdrawal = (*withdrawal)(nil)

type Withdrawal interface {
	Withdraw(w http.ResponseWriter, r *http.Request)
	GetWithdrawals(w http.ResponseWriter, r *http.Request)
}

type withdrawal struct {
	l *zap.SugaredLogger
	s service.WithdrawalService
}

func NewWithdrawal(l *zap.SugaredLogger, s service.WithdrawalService) *withdrawal {
	return &withdrawal{l: l, s: s}
}

func (wh *withdrawal) SetupRouting(r *chi.Mux, m middleware.Middleware) {
	r.Get("/api/user/balance/withdrawals", m.Middleware(wh.GetWithdrawals))
	r.Post("/api/user/balance/withdraw", m.Middleware(wh.Withdraw))
}

func (wh *withdrawal) Withdraw(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.Key).(string)
	if !ok {
		wh.l.Error("error getting user uuid from context")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var withdrawal model.Withdrawal
	if err := json.NewDecoder(r.Body).Decode(&withdrawal); err != nil {
		wh.l.Errorw("error decoding body", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := wh.s.Withdraw(user, withdrawal.Order, withdrawal.Sum); err != nil {
		if errors.Is(err, service.ErrNoOrder) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		if errors.Is(err, service.ErrInsufficientBalance) {
			w.WriteHeader(http.StatusPaymentRequired)
			return
		}

		wh.l.Errorw("error withdrawing", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (wh *withdrawal) GetWithdrawals(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user, ok := r.Context().Value(middleware.Key).(string)
	if !ok {
		wh.l.Error("error getting user uuid from context")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	withdrawals, err := wh.s.GetWithdrawals(user)
	if err != nil {
		wh.l.Errorw("error getting withdrawals", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(withdrawals) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	bb, err := json.Marshal(model.ToAPIWithdrawals(withdrawals))
	if err != nil {
		wh.l.Errorw("error marshallin response withdrawals", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bb)
}
