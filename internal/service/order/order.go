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

	if err := o.s.AddAccrual(order); err != nil {
		return err
	}

	go func(l *zap.SugaredLogger, order *model.Order) {
		for {
			n, err := o.a.GetAccrual(order)
			if err != nil {
				if errors.Is(err, accrual.ErrTooManyRequests) {
					time.Sleep(time.Duration(n) * time.Second)
					continue
				}
				o.l.Errorf("error getting accrual - %v, error - %s", order, err)
				return
			}

			switch order.Status {
			case model.StatusInvalid:
				return
			case model.StatusProcessed:
			default:
				continue
			}

			if err := o.s.AddAccrual(order); err != nil {
				o.l.Errorf("error adding accrual - %v, error - %s", order, err)
				return
			}
			return
		}
	}(o.l, order)

	return nil
}
