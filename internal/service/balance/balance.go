package balance

import (
	"github.com/Fe4p3b/gophermart/internal/model"
	"go.uber.org/zap"
)

type Balance interface {
	Get(userId string) error
}

type balance struct {
	l *zap.SugaredLogger
}

func New(l *zap.SugaredLogger) *balance {
	return &balance{l: l}
}

func (b *balance) Get(userId string) error {
	return nil
}

func (b *balance) Withdraw(orderId string, sum uint) error {
	return nil
}

func (b *balance) getWithdrawals(userId string) ([]model.Bonus, error) {
	return []model.Bonus{}, nil
}
