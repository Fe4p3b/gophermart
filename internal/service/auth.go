package service

import (
	"errors"

	"github.com/Fe4p3b/gophermart/internal/api/models"
	"go.uber.org/zap"
)

var _ Auth = &auth{}

var ErrUserExists error = errors.New("login already exists")

type Auth interface {
	Register(c *models.Credentials) error
}

type auth struct {
	l *zap.SugaredLogger
}

func NewAuth(l *zap.SugaredLogger) *auth {
	return &auth{l: l}
}

func (a *auth) Register(c *models.Credentials) error {
	return ErrUserExists
}
