package storage

import "github.com/Fe4p3b/gophermart/internal/model"

type AuthStorer interface {
	AddUser(*model.User) error
}
