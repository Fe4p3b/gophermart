package accrual

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	apiModel "github.com/Fe4p3b/gophermart/internal/api/model"
	"github.com/Fe4p3b/gophermart/internal/model"
	"go.uber.org/zap"
)

var (
	_ AccrualAcquirer = (*accrual)(nil)

	ErrTooManyRequests error = errors.New("error too many requests")
	ErrGettingOrder    error = errors.New("error no order")
)

type AccrualAcquirer interface {
	GetAccrual(*model.Order) (int, error)
}

type accrual struct {
	l       *zap.SugaredLogger
	baseURL string
}

func New(l *zap.SugaredLogger, u string) *accrual {
	return &accrual{baseURL: u, l: l}
}

func (a *accrual) GetAccrual(o *model.Order) (int, error) {
	URL := fmt.Sprintf("%s/api/orders/%s", a.baseURL, o.Number)
	resp, err := http.Get(URL)
	if err != nil {
		return 0, err
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		s := resp.Header.Get("Retry-After")
		n, err := strconv.Atoi(s)
		if err != nil {
			return 0, err
		}
		return n, ErrTooManyRequests
	}

	if resp.StatusCode != http.StatusOK {
		return 0, ErrGettingOrder
	}

	var order apiModel.Order
	if err := json.NewDecoder(resp.Body).Decode(&order); err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	status, err := model.ToOrderStatus(order.Status)
	if err != nil {
		return 0, err
	}
	o.Status = status
	o.Accrual = uint32(order.Accrual * 100)

	return 0, nil
}
