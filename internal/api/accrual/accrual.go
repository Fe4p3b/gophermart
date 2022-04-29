package accrual

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Fe4p3b/gophermart/internal/model"
	"go.uber.org/zap"
)

var (
	_ AccrualAcquirer = (*accrual)(nil)
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
	a.l.Infof("accrual %#v", o)

	resp, err := http.Get(fmt.Sprintf("%s/%s", a.baseURL, o.Number))
	if err != nil {
		return err
	}

	if err := json.NewDecoder(resp.Body).Decode(&o); err != nil {
		return err
	}
	defer resp.Body.Close()

	a.l.Infof("accrual %#v", o)

	return nil
}
