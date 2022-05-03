package balance

import (
	"errors"

	"github.com/Fe4p3b/gophermart/internal/model"
	"github.com/Fe4p3b/gophermart/internal/storage"
	"go.uber.org/zap"
)

var _ BalanceService = (*balanceService)(nil)

var (
	ErrNoOrder             error = errors.New("order with such number doesn't exist")
	ErrInsufficientBalance error = errors.New("insufficient balance")
)

type BalanceService interface {
	Get(userID string) (*model.Balance, error)
	Withdraw(string, string, float64) error
	GetWithdrawals(string) ([]model.Withdrawal, error)
}

type balanceService struct {
	l *zap.SugaredLogger
	b storage.BalanceRepository
	w storage.WithdrawalRepository
}

func New(l *zap.SugaredLogger, b storage.BalanceRepository, w storage.WithdrawalRepository) *balanceService {
	return &balanceService{l: l, b: b, w: w}
}

func (b *balanceService) Get(userID string) (*model.Balance, error) {
	ub, err := b.b.GetBalanceForUser(userID)
	if err != nil {
		return nil, err
	}

	return ub, nil
}

func (b *balanceService) Withdraw(userID, orderNumber string, sum float64) error {
	ub, err := b.b.GetBalanceForUser(userID)
	if err != nil {
		return err
	}

	s := uint64(sum * 100)
	if s > ub.Current {
		return ErrInsufficientBalance
	}

	if err := b.w.AddWithdrawal(userID, model.Withdrawal{OrderNumber: orderNumber, Sum: s}); err != nil {
		return err
	}
	return nil
}

func (b *balanceService) GetWithdrawals(userID string) ([]model.Withdrawal, error) {
	w, err := b.w.GetWithdrawalsForUser(userID)
	if err != nil {
		return nil, err
	}

	return w, nil
}
