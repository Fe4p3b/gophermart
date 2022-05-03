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
	AddOrder(string, string) error
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

func (o *orderService) AddOrder(userID, number string) error {
	order := &model.Order{UserID: userID, Number: number, Status: model.StatusNew, UploadDate: time.Now()}

	if err := o.s.AddOrder(order); err != nil {
		return err
	}

	go func(l *zap.SugaredLogger, order *model.Order) {
		for {
			time.Sleep(20 * time.Second)
			n, err := o.a.GetAccrual(order)
			if err != nil {
				if errors.Is(err, accrual.ErrTooManyRequests) {
					time.Sleep(time.Duration(n) * time.Second)
					continue
				}
				o.l.Errorw("error getting accrual", "order", order, "error", err)
				continue
			}

			o.l.Infow("GetAccrual", "order", order)

			if err := o.s.UpdateOrder(order); err != nil {
				o.l.Errorw("error adding accrual", "order", order, "error", err)
				return
			}

			switch order.Status {
			case model.StatusInvalid:
				return
			case model.StatusProcessed:
				if err := o.s.UpdateBalanceForProcessedOrder(order); err != nil {
					o.l.Errorw("error updating balance", "order", order, "error", err)
					return
				}
			default:
				continue
			}
			return
		}
	}(o.l, order)

	return nil
}
