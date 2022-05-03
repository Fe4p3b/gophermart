package model

import (
	"time"

	"github.com/Fe4p3b/gophermart/internal/model"
)

type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Order struct {
	Number     string  `json:"number"`
	Status     string  `json:"status"`
	Accrual    float64 `json:"accrual"`
	UploadDate string  `json:"uploaded_at"`
}

type Balance struct {
	Current   float64 `json:"current"`
	Withdrawn uint64  `json:"withdrawn"`
}

type Withdrawal struct {
	Order       string  `json:"order"`
	Sum         float64 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
}

func ToAPIWithdrawal(w model.Withdrawal) Withdrawal {
	return Withdrawal{Order: w.OrderNumber, Sum: float64(w.Sum) / 100, ProcessedAt: w.Date.Format(time.RFC3339)}
}

func ToAPIWithdrawals(withdrawals []model.Withdrawal) []Withdrawal {
	w := make([]Withdrawal, 0)
	for _, v := range withdrawals {
		w = append(w, ToAPIWithdrawal(v))
	}
	return w
}
