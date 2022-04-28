package order

import (
	"github.com/Fe4p3b/gophermart/internal/model"
	"github.com/Fe4p3b/gophermart/internal/storage"
	"go.uber.org/zap"
)

var _ OrderService = (*orderService)(nil)

type OrderService interface {
	List(string) ([]model.Order, error)
	AddAccrual(string, string) error
}

type orderService struct {
	l *zap.SugaredLogger
	s storage.OrderRepository
}

func New(l *zap.SugaredLogger, s storage.OrderRepository) *orderService {
	return &orderService{l: l, s: s}
}

func (o *orderService) List(userId string) ([]model.Order, error) {
	orders, err := o.s.GetOrdersForUser(userId)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (o *orderService) AddAccrual(userId, id string) error {
	sum := uint32(999)
	if err := o.s.AddAccrual("asd", id, sum); err != nil {
		return err
	}

	return nil
}
