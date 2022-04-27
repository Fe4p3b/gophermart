package auth

import (
	"errors"

	"github.com/Fe4p3b/gophermart/internal/model"
	"github.com/Fe4p3b/gophermart/internal/storage"
	"go.uber.org/zap"
)

var _ Auth = (*auth)(nil)

var ErrUserExists error = errors.New("login already exists")

type Auth interface {
	Register(string, string) error
	Login(string, string) error
}

type auth struct {
	l *zap.SugaredLogger
	s storage.AuthStorer
}

func NewAuth(l *zap.SugaredLogger) *auth {
	return &auth{l: l}
}

func (a *auth) Register(l string, p string) error {
	a.s.AddUser(&model.User{Login: l, Passord: p})
	return ErrUserExists
}

func (a *auth) Login(l string, p string) error {
	return nil
}
