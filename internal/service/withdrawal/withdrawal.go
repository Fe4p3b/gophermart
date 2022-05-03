package withdrawal

import (
	"errors"

	"github.com/Fe4p3b/gophermart/internal/model"
	"github.com/Fe4p3b/gophermart/internal/service/balance"
	"github.com/Fe4p3b/gophermart/internal/storage"
	"go.uber.org/zap"
)

var _ WithdrawalService = (*withdrawalService)(nil)

var (
	ErrNoOrder             error = errors.New("order with such number doesn't exist")
	ErrInsufficientBalance error = errors.New("insufficient balance")
)

type WithdrawalService interface {
	Withdraw(string, string, float64) error
	GetWithdrawals(string) ([]model.Withdrawal, error)
}

type withdrawalService struct {
	l *zap.SugaredLogger
	w storage.WithdrawalRepository
	b balance.BalanceService
}

func New(l *zap.SugaredLogger, w storage.WithdrawalRepository, b balance.BalanceService) *withdrawalService {
	return &withdrawalService{l: l, w: w, b: b}
}

func (ws *withdrawalService) Withdraw(userID, orderNumber string, sum float64) error {
	ub, err := ws.b.Get(userID)
	if err != nil {
		return err
	}

	s := uint64(sum * 100)
	if s > ub.Current {
		return ErrInsufficientBalance
	}

	if err := ws.w.AddWithdrawal(model.Withdrawal{OrderNumber: orderNumber, Sum: s, UserID: userID}); err != nil {
		return err
	}
	return nil
}

func (ws *withdrawalService) GetWithdrawals(userID string) ([]model.Withdrawal, error) {
	w, err := ws.w.GetWithdrawalsForUser(userID)
	if err != nil {
		return nil, err
	}

	return w, nil
}
