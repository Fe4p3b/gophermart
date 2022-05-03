package order

import (
	"errors"
	"time"

	"github.com/Fe4p3b/gophermart/internal/api/accrual"
	"github.com/Fe4p3b/gophermart/internal/model"
	"github.com/Fe4p3b/gophermart/internal/storage"
	"go.uber.org/zap"
)

var _ OrderService = (*orderService)(nil)

var (
	ErrOrderForUserExists        error = errors.New("order for user already exists")
	ErrOrderExistsForAnotherUser error = errors.New("order already exists for another user")
)

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

func (o *orderService) List(userID string) ([]model.Order, error) {
	orders, err := o.s.GetOrdersForUser(userID)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (o *orderService) AddAccrual(userID, number string) error {
	order := &model.Order{UserID: userID, Number: number, Status: model.StatusProcessing, UploadDate: time.Now()}

	err := o.a.GetAccrual(order)
	if err != nil {
		return err
	}

	if err := o.s.AddAccrual(order); err != nil {
		return err
	}

	return nil
}
