package pg

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"os"
	"time"

	"github.com/Fe4p3b/gophermart/internal/model"
	"github.com/Fe4p3b/gophermart/internal/service/auth"
	"github.com/Fe4p3b/gophermart/internal/service/balance"
	"github.com/Fe4p3b/gophermart/internal/service/order"
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

func (p *pg) InitialMigration() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	sql, err := os.ReadFile("./migrations/001_init.sql")
	if err != nil {
		return err
	}

	_, err = p.db.ExecContext(ctx, string(sql))
	if err != nil {
		return err
	}

	return nil
}

func (p *pg) AddUser(u *model.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := p.db.BeginTx(ctx, nil)
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

func (p *pg) VerifyUser(u string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sql := `SELECT id FROM gophermart.users WHERE id=$1`

	row := p.db.QueryRowContext(ctx, sql, u)
	var uuid string

	if err := row.Scan(&uuid); err != nil {
		return err
	}

	return nil
}

func (p *pg) GetOrdersForUser(u string) ([]model.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sql := `SELECT number, user_id, number, status, accrual, upload_date FROM gophermart.orders WHERE user_id = $1`

	rows, err := p.db.QueryContext(ctx, sql, u)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var o model.Order
		var s string
		if err := rows.Scan(&o.Number, &o.UserID, &o.Number, &s, &o.Accrual, &o.UploadDate); err != nil {
			return nil, err
		}

		status, err := model.ToOrderStatus(s)
		if err != nil {
			return nil, err
		}
		o.Status = status
		orders = append(orders, o)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (p *pg) AddOrder(o *model.Order) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sql := `INSERT INTO gophermart.orders(user_id, number, status, accrual, upload_date) VALUES ($1, $2, $3, $4, $5)`
	if _, err := p.db.ExecContext(ctx, sql, o.UserID, o.Number, o.Status.String(), o.Accrual, o.UploadDate); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			sql := `SELECT user_id FROM gophermart.orders WHERE number=$1`
			row := p.db.QueryRowContext(ctx, sql, o.Number)

			var userID string
			if err = row.Scan(&userID); err != nil {
				return err
			}

			if o.UserID != userID {
				return order.ErrOrderExistsForAnotherUser
			}
			return order.ErrOrderForUserExists
		}

		return err
	}

	return nil
}

func (p *pg) UpdateOrder(o *model.Order) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sql := `UPDATE gophermart.orders SET status = $1, accrual = $2 WHERE number = $3`

	if _, err := p.db.ExecContext(ctx, sql, o.Status, o.Accrual, o.Number); err != nil {
		return err
	}
	return nil
}

func (p *pg) UpdateBalanceForProcessedOrder(o *model.Order) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sql := `UPDATE gophermart.balances SET current = current + $1 WHERE user_id = $2`
	if _, err := p.db.ExecContext(ctx, sql, o.Accrual, o.UserID); err != nil {
		return err
	}

	return nil
}

func (p *pg) GetBalanceForUser(u string) (*model.Balance, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sql := `SELECT b.id, b.user_id, b.current, COALESCE(SUM(w.sum),0) as withdrawn
FROM gophermart.balances as b
LEFT JOIN gophermart.withdrawals as w
ON b.user_id = w.user_id
WHERE b.user_id = $1
GROUP BY b.id`
	row := p.db.QueryRowContext(ctx, sql, u)

	var balance model.Balance
	if err := row.Scan(&balance.ID, &balance.UserID, &balance.Current, &balance.Withdrawn); err != nil {
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
	defer tx.Rollback()

	sql := `INSERT INTO gophermart.withdrawals(user_id, order_number, sum, date) VALUES($1, $2, $3, $4)`

	if _, err := tx.ExecContext(ctx, sql, u, w.OrderNumber, w.Sum, w.Date); err != nil {
		log.Printf("AddWithdrawal - %s", err)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.ForeignKeyViolation {
			return balance.ErrNoOrder
		}

		return err
	}

	sql = `UPDATE gophermart.balances SET current=current-$1 WHERE user_id=$2`

	if _, err := tx.ExecContext(ctx, sql, w.Sum, u); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (p *pg) GetWithdrawalsForUser(u string) ([]model.Withdrawal, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sql := `SELECT id, order_number, sum, date FROM gophermart.withdrawals WHERE user_id = $1`

	rows, err := p.db.QueryContext(ctx, sql, u)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var withdrawals []model.Withdrawal
	for rows.Next() {
		var w model.Withdrawal
		if err := rows.Scan(&w.ID, &w.OrderNumber, &w.Sum, &w.Date); err != nil {
			return nil, err
		}
		withdrawals = append(withdrawals, w)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return withdrawals, nil
}
