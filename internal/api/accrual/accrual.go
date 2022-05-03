package accrual

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Fe4p3b/gophermart/internal/model"
	"go.uber.org/zap"
)

var (
	_ AccrualAcquirer = (*accrual)(nil)

	ErrTooManyRequests error = errors.New("too many requests")
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
	resp, err := http.Get(fmt.Sprintf("%s/api/orders/%s", a.baseURL, o.Number))
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

	a.l.Infof("status - %v", resp.StatusCode)
	a.l.Infof("before accrual - %v", o)
	if err := json.NewDecoder(resp.Body).Decode(&o); err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	a.l.Infof("after accrual - %v", o)

	return 0, nil
}
