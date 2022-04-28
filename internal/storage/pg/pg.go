package pg

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Fe4p3b/gophermart/internal/model"
	"github.com/Fe4p3b/gophermart/internal/service/auth"
	"github.com/Fe4p3b/gophermart/internal/storage"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var (
	_ storage.UserRepository       = (*pg)(nil)
	_ storage.OrderRepository      = (*pg)(nil)
	_ storage.BalanceRepository    = (*pg)(nil)
	_ storage.WithdrawalRepository = (*pg)(nil)
)

type pg struct {
	db *sql.DB
}

func New(dsn string) (*pg, error) {
	conn, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	return &pg{db: conn}, nil
}

func (p *pg) AddUser(u *model.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sql := `INSERT INTO gophermart.users(login, password) VALUES($1, $2)`

	_, err := p.db.ExecContext(ctx, sql, u.Login, u.Passord)
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
		return auth.ErrUserExists
	}

	return err
}

func (p *pg) GetUserByLogin(l string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sql := `SELECT id, login, password FROM gophermart.users WHERE login = $1`

	var u model.User
	row := p.db.QueryRowContext(ctx, sql, l)
	if err := row.Scan(&u.ID, &u.Login, &u.Passord); err != nil {
		return nil, err
	}

	return &u, nil
}

func (p *pg) GetOrdersForUser(u string) ([]model.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sql := `SELECT id, user_id, number, status, accrual, upload_date FROM gophermart.orders WHERE user_id = $1`

	rows, err := p.db.QueryContext(ctx, sql, u)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var o model.Order
		if err := rows.Scan(&o.ID, &o.UserID, &o.Number, &o.Status, &o.Accrual, &o.UploadDate); err != nil {
			return nil, err
		}

		orders = append(orders, o)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (p *pg) AddAccrual(u, o string, s uint32) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	sql := `UPDATE gophermart.orders SET accrual = $1 WHERE id = $2`
	if _, err := tx.ExecContext(ctx, sql, o, s); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	sql = `UPDATE gophermart.balance SET current = current + $1 WHERE user_id = $2`
	if _, err := tx.ExecContext(ctx, sql, s, u); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	return nil
}

func (p *pg) GetForUser(u string) (*model.Balance, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sql := `SELECT id, user_id, current WHERE user_id = $1`
	row := p.db.QueryRowContext(ctx, sql, u)

	var balance model.Balance
	if err := row.Scan(&balance.ID, &balance.UserID, &balance.Current); err != nil {
		return nil, err
	}

	return &balance, nil
}

func (p *pg) AddWithdrawal(u string, w model.Withdrawal) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	sql := `INSERT INTO gophermart.withdrawals(id, order_id, sum, date) VALUES($1, $2, $3, $4)`

	if _, err := tx.ExecContext(ctx, sql, w.ID, w.OrderID, w.Sum, w.Date); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	sql = `UPDATE gophermart.balance SET current=current-$1 WHERE user_id=$2`

	if _, err := tx.ExecContext(ctx, sql, w.Sum, u); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	if err = tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	return nil
}

func (p *pg) GetWithdrawalsForUser(u string) ([]model.Withdrawal, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sql := `SELECT id, order_id, sum, date FROM gophermart.withdrawals as w, gophermart.orders as o WHERE w.order_id=o.order_id and o.user_id = $1`

	rows, err := p.db.QueryContext(ctx, sql, u)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var withdrawals []model.Withdrawal
	for rows.Next() {
		var w model.Withdrawal
		if err := rows.Scan(&w.ID, &w.OrderID, &w.Sum, &w.Date); err != nil {
			return nil, err
		}
		withdrawals = append(withdrawals, w)
	}

	return withdrawals, nil
}
