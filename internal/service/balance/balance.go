package balance

import (
	"github.com/Fe4p3b/gophermart/internal/model"
	"github.com/Fe4p3b/gophermart/internal/storage"
	"go.uber.org/zap"
)

type BalanceService interface {
	Get(userId string) (*model.Balance, error)
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
	ub, err := b.b.GetForUser(userId)
	if err != nil {
		return nil, err
	}

	return ub, nil
}

func (b *balanceService) Withdraw(userId, orderId string, sum uint32) error {
	if err := b.w.AddWithdrawal(userId, model.Withdrawal{OrderID: orderId, Sum: sum}); err != nil {
		return err
	}
	return nil
}

func (b *balanceService) getWithdrawals(userId string) ([]model.Withdrawal, error) {
	w, err := b.w.GetWithdrawalsForUser(userId)
	if err != nil {
		return nil, err
	}

	return w, nil
}
