package storage

import "github.com/Fe4p3b/gophermart/internal/model"

type UserRepository interface {
	AddUser(*model.User) error
	GetUserByLogin(string) (*model.User, error)
	VerifyUser(string) error
}

type OrderRepository interface {
	GetOrdersForUser(string) ([]model.Order, error)
	AddOrder(*model.Order) error
	UpdateOrder(*model.Order) error
	UpdateBalanceForProcessedOrder(*model.Order) error
}

type BalanceRepository interface {
	GetBalanceForUser(string) (*model.Balance, error)
}

type WithdrawalRepository interface {
	AddWithdrawal(string, model.Withdrawal) error
	GetWithdrawalsForUser(string) ([]model.Withdrawal, error)
}
