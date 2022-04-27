package model

import "time"

type User struct {
	ID      string
	Login   string
	Passord string
}

type Bonus struct {
	ID      string
	OrderID string
	Sum     uint
	Date    time.Time
}
