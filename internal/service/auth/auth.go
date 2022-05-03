package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"

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

type AuthService interface {
	Register(string, string) error
	Login(string, string) (string, error)
	VerifyUser(string) (string, error)
}

type AuthServiceConfiguration func(as *AuthService) error

type authService struct {
	l        *zap.SugaredLogger
	s        storage.UserRepository
	hashCost int
	key      [32]byte
	aesgcm   cipher.AEAD
}

func NewAuth(l *zap.SugaredLogger, s storage.UserRepository, c int, k []byte) (*authService, error) {
	authKey := sha256.Sum256(k)
	aesblock, err := aes.NewCipher(authKey[:])
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return nil, err
	}

	return &authService{l: l, s: s, hashCost: c, key: authKey, aesgcm: aesgcm}, nil
}

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

func (a *authService) Login(l string, p string) (string, error) {
	u, err := a.s.GetUserByLogin(l)
	if err != nil {
		return "", err
	}

	if err := a.checkPasswordHash(p, u.Passord); err != nil {
		return "", ErrWrongCredentials
	}

	token, err := a.encrypt(u.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (a *authService) VerifyUser(token string) (string, error) {
	uuid, err := a.decrypt(token)
	if err != nil {

		return "", err
	}

	if err := a.s.VerifyUser(string(uuid)); err != nil {
		return "", err
	}

	return string(uuid), nil
}

func (a *authService) hashPassword(p string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(p), a.hashCost)
	return string(bytes), err
}

func (a *authService) checkPasswordHash(p string, h string) error {
	return bcrypt.CompareHashAndPassword([]byte(h), []byte(p))
}

func (a *authService) encrypt(src string) (string, error) {
	nonce := a.key[len(a.key)-a.aesgcm.NonceSize():]
	dst := (a.aesgcm.Seal(nil, nonce, []byte(src), nil))

	return fmt.Sprintf("%x", dst), nil
}

func (a *authService) decrypt(src string) ([]byte, error) {
	nonce := a.key[len(a.key)-a.aesgcm.NonceSize():]
	encrypted, err := hex.DecodeString(src)
	if err != nil {
		return nil, err
	}

	dst, err := a.aesgcm.Open(nil, nonce, encrypted, nil)
	if err != nil {
		return nil, err
	}

	return dst, nil
}
