package storage

import "github.com/Fe4p3b/gophermart/internal/model"

type UserRepository interface {
	AddUser(*model.User) error
	GetUserByLogin(string) (*model.User, error)
}

type OrderRepository interface {
	GetOrdersForUser(string) ([]model.Order, error)
	AddAccrual(*model.Order) error
}

type BalanceRepository interface {
	GetForUser(string) (*model.Balance, error)
}

type WithdrawalRepository interface {
	AddWithdrawal(string, model.Withdrawal) error
	GetWithdrawalsForUser(string) ([]model.Withdrawal, error)
}
