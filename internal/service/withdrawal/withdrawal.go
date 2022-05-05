package withdrawal

import (
	"context"
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
	Withdraw(context.Context, string, string, float64) error
	GetWithdrawals(context.Context, string) ([]model.Withdrawal, error)
}

type withdrawalService struct {
	l *zap.SugaredLogger
	w storage.WithdrawalRepository
	b balance.BalanceService
}

func New(l *zap.SugaredLogger, w storage.WithdrawalRepository, b balance.BalanceService) *withdrawalService {
	return &withdrawalService{l: l, w: w, b: b}
}

func (ws *withdrawalService) Withdraw(ctx context.Context, userID, orderNumber string, sum float64) error {
	ub, err := ws.b.Get(ctx, userID)
	if err != nil {
		return err
	}

	s := uint64(sum * 100)
	if s > ub.Current {
		return ErrInsufficientBalance
	}

	if err := ws.w.AddWithdrawal(ctx, model.Withdrawal{OrderNumber: orderNumber, Sum: s, UserID: userID}); err != nil {
		return err
	}
	return nil
}

func (ws *withdrawalService) GetWithdrawals(ctx context.Context, userID string) ([]model.Withdrawal, error) {
	w, err := ws.w.GetWithdrawalsForUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	return w, nil
}
