package storage

import (
	"context"

	"github.com/Fe4p3b/gophermart/internal/model"
)

type UserRepository interface {
	AddUser(context.Context, *model.User) error
	GetUserByLogin(context.Context, string) (*model.User, error)
	VerifyUser(context.Context, string) error
}

type OrderRepository interface {
	GetOrdersForUser(context.Context, string) ([]model.Order, error)
	AddOrder(context.Context, *model.Order) error
	UpdateOrder(context.Context, *model.Order) error
	UpdateBalanceForProcessedOrder(context.Context, *model.Order) error
}

type BalanceRepository interface {
	GetBalanceForUser(context.Context, string) (*model.Balance, error)
}

type WithdrawalRepository interface {
	AddWithdrawal(context.Context, model.Withdrawal) error
	GetWithdrawalsForUser(context.Context, string) ([]model.Withdrawal, error)
}
