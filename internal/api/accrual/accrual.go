package accrual

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Fe4p3b/gophermart/internal/model"
	"go.uber.org/zap"
)

var (
	_ AccrualAcquirer = (*accrual)(nil)

	ErrTooManyRequests error = errors.New("too many requests")
)

type AccrualAcquirer interface {
	GetAccrual(*model.Order) error
}

type accrual struct {
	l       *zap.SugaredLogger
	baseURL string
}

func New(l *zap.SugaredLogger, u string) *accrual {
	return &accrual{baseURL: u, l: l}
}

func (a *accrual) GetAccrual(o *model.Order) error {
	resp, err := http.Get(fmt.Sprintf("%s/%s", a.baseURL, o.Number))
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		return ErrTooManyRequests
	}

	a.l.Infof("before accrual - %v", o)
	if err := json.NewDecoder(resp.Body).Decode(&o); err != nil {
		return err
	}
	defer resp.Body.Close()
	a.l.Infof("after accrual - %v", o)

	return nil
}
