package order

import (
	"time"

	"github.com/Fe4p3b/gophermart/internal/api/accrual"
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
	a accrual.AccrualAcquirer
}

func New(l *zap.SugaredLogger, s storage.OrderRepository, a accrual.AccrualAcquirer) *orderService {
	return &orderService{l: l, s: s, a: a}
}

func (o *orderService) List(userId string) ([]model.Order, error) {
	orders, err := o.s.GetOrdersForUser(userId)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (o *orderService) AddAccrual(userId, number string) error {
	order := &model.Order{UserID: userId, Number: number, Status: model.StatusProcessed, UploadDate: time.Now()}
	o.l.Infof("%#v", order)

	err := o.a.GetAccrual(order)
	if err != nil {
		return err
	}

	o.l.Infof("%#v", order)
	if err := o.s.AddAccrual(order); err != nil {
		return err
	}

	return nil
}
