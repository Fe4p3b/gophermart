package auth

import (
	"errors"

	"github.com/Fe4p3b/gophermart/internal/model"
	"github.com/Fe4p3b/gophermart/internal/storage"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var _ AuthService = (*authService)(nil)

var (
	ErrUserExists       error = errors.New("login already exists")
	ErrWrongCredentials error = errors.New("wrong credentials")
)

const hashCost = 14

type AuthService interface {
	Register(string, string) error
	Login(string, string) error
}

type AuthServiceConfiguration func(as *AuthService) error

type authService struct {
	l *zap.SugaredLogger
	s storage.UserRepository
}

func NewAuth(l *zap.SugaredLogger, s storage.UserRepository) *authService {
	return &authService{l: l, s: s}
}

func WithPostgresRepository() {}

func (a *authService) Register(l string, p string) error {
	hash, err := a.hashPassword(p)
	if err != nil {
		return err
	}

	if err := a.s.AddUser(&model.User{Login: l, Passord: string(hash)}); err != nil {
		return err
	}

	return nil
}

func (a *authService) Login(l string, p string) error {
	u, err := a.s.GetUserByLogin(l)
	if err != nil {
		return err
	}
	a.l.Info(u)
	if err := a.checkPasswordHash(p, u.Passord); err != nil {
		return ErrWrongCredentials
	}

	return nil
}

func (a *authService) hashPassword(p string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(p), hashCost)
	return string(bytes), err
}

func (a *authService) checkPasswordHash(p string, h string) error {
	return bcrypt.CompareHashAndPassword([]byte(h), []byte(p))
}
