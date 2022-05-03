package pg

import (
	"context"
	"errors"
	"time"

	"github.com/Fe4p3b/gophermart/internal/model"
	"github.com/Fe4p3b/gophermart/internal/service/auth"
	"github.com/Fe4p3b/gophermart/internal/storage"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
)

var (
	_ storage.UserRepository = (*UserStorage)(nil)
)

type UserStorage struct {
	pg *pg
}

func NewUserStorage(pg *pg) *UserStorage {
	return &UserStorage{pg: pg}
}
func (us *UserStorage) AddUser(u *model.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := us.pg.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	sql := `INSERT INTO gophermart.users(login, password) VALUES($1, $2) RETURNING id`

	row := tx.QueryRowContext(ctx, sql, u.Login, u.Passord)

	if err := row.Scan(&u.ID); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return auth.ErrUserExists
		}

		return err
	}

	sql = `INSERT INTO gophermart.balances(user_id) VALUES($1)`
	_, err = tx.ExecContext(ctx, sql, u.ID)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return err
}

func (us *UserStorage) GetUserByLogin(l string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sql := `SELECT id, login, password FROM gophermart.users WHERE login = $1`

	var u model.User
	row := us.pg.db.QueryRowContext(ctx, sql, l)
	if err := row.Scan(&u.ID, &u.Login, &u.Passord); err != nil {
		return nil, err
	}

	return &u, nil
}

func (us *UserStorage) VerifyUser(u string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sql := `SELECT id FROM gophermart.users WHERE id=$1`

	row := us.pg.db.QueryRowContext(ctx, sql, u)
	var uuid string

	if err := row.Scan(&uuid); err != nil {
		return err
	}

	return nil
}
