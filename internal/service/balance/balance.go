package balance

import (
	"context"

	"github.com/Fe4p3b/gophermart/internal/model"
	"github.com/Fe4p3b/gophermart/internal/storage"
	"go.uber.org/zap"
)

var _ BalanceService = (*balanceService)(nil)

type BalanceService interface {
	Get(context.Context, string) (*model.Balance, error)
}

type balanceService struct {
	l *zap.SugaredLogger
	b storage.BalanceRepository
}

func New(l *zap.SugaredLogger, b storage.BalanceRepository) *balanceService {
	return &balanceService{l: l, b: b}
}

func (b *balanceService) Get(ctx context.Context, userID string) (*model.Balance, error) {
	ub, err := b.b.GetBalanceForUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	return ub, nil
}
