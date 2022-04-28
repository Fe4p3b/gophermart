package model

import (
	"fmt"
	"time"
)

type User struct {
	ID      string
	Login   string
	Passord string
}

type Balance struct {
	ID      string
	UserID  string
	Current uint32
}

type Withdrawal struct {
	ID      string
	OrderID string
	Sum     uint32
	Date    time.Time
}

type OrderStatus int8

const (
	StatusNew OrderStatus = iota + 1
	StatusProcessing
	StatusInvalid
	StatusProcessed
)

func (o OrderStatus) String() string {
	switch o {
	case StatusNew:
		return "NEW"
	case StatusProcessing:
		return "PROCESSING"
	case StatusProcessed:
		return "PROCESSED"
	case StatusInvalid:
		return "INVALID"
	default:
		return fmt.Sprintf("status unknown - %d", o)
	}
}

type Order struct {
	ID         string
	UserID     string
	Number     string
	Status     OrderStatus
	Accrual    uint32
	UploadDate time.Time
}
