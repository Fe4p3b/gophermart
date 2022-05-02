package model

import (
	"errors"
	"fmt"
	"time"
)

var ErrUnknownStatus = errors.New("unknown status")

type User struct {
	ID      string
	Login   string
	Passord string
}

type Balance struct {
	ID        string
	UserID    string
	Current   uint32
	Withdrawn uint64
}

type Withdrawal struct {
	ID          string
	OrderNumber string
	UserID      string
	Sum         uint64
	Date        time.Time
}

type Order struct {
	Number     string
	UserID     string
	Status     OrderStatus
	Accrual    uint32
	UploadDate time.Time
}

type OrderStatus int8

const (
	StatusNew OrderStatus = iota + 1
	StatusProcessing
	StatusInvalid
	StatusProcessed
)

func ToOrderStatus(s string) (OrderStatus, error) {
	switch s {
	case "NEW":
		return StatusNew, nil
	case "PROCESSING":
		return StatusProcessing, nil
	case "PROCESSED":
		return StatusProcessed, nil
	case "INVALID":
		return StatusInvalid, nil
	}
	return 0, ErrUnknownStatus
}

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
