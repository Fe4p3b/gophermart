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
	Get(userId string) (*model.Balance, error)
	Withdraw(string, string, uint64) error
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

func (b *balanceService) Get(userId string) (*model.Balance, error) {
	ub, err := b.b.GetBalanceForUser(userId)
	if err != nil {
		return nil, err
	}

	return ub, nil
}

func (b *balanceService) Withdraw(userId, orderNumber string, sum uint64) error {
	ub, err := b.b.GetBalanceForUser(userId)
	if err != nil {
		return err
	}

	if sum > uint64(ub.Current) {
		return ErrInsufficientBalance
	}

	if err := b.w.AddWithdrawal(userId, model.Withdrawal{OrderNumber: orderNumber, Sum: sum}); err != nil {
		return err
	}
	return nil
}

func (b *balanceService) GetWithdrawals(userId string) ([]model.Withdrawal, error) {
	w, err := b.w.GetWithdrawalsForUser(userId)
	if err != nil {
		return nil, err
	}

	return w, nil
}
