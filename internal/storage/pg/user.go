package pg

import (
	"context"
	sqlErr "database/sql"
	"errors"

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
func (us *UserStorage) AddUser(ctx context.Context, u *model.User) error {
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

func (us *UserStorage) GetUserByLogin(ctx context.Context, l string) (*model.User, error) {
	sql := `SELECT id, login, password FROM gophermart.users WHERE login = $1`

	var u model.User
	row := us.pg.db.QueryRowContext(ctx, sql, l)
	if err := row.Scan(&u.ID, &u.Login, &u.Passord); err != nil {
		if errors.Is(err, sqlErr.ErrNoRows) {
			return nil, auth.ErrWrongCredentials
		}
		return nil, err
	}

	return &u, nil
}

func (us *UserStorage) VerifyUser(ctx context.Context, u string) error {
	sql := `SELECT id FROM gophermart.users WHERE id=$1`

	row := us.pg.db.QueryRowContext(ctx, sql, u)
	var uuid string

	if err := row.Scan(&uuid); err != nil {
		if errors.Is(err, sqlErr.ErrNoRows) {
			return auth.ErrWrongCredentials
		}
		return err
	}

	return nil
}
