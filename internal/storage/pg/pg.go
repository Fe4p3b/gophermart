package pg

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"os"
	"time"

	"github.com/Fe4p3b/gophermart/internal/model"
	"github.com/Fe4p3b/gophermart/internal/service/balance"
	"github.com/Fe4p3b/gophermart/internal/storage"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var (
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

func (p *pg) AddWithdrawal(w model.Withdrawal) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	sql := `INSERT INTO gophermart.withdrawals(user_id, order_number, sum, date) VALUES($1, $2, $3, $4)`

	if _, err := tx.ExecContext(ctx, sql, w.UserID, w.OrderNumber, w.Sum, w.Date); err != nil {
		log.Printf("AddWithdrawal - %v, %v", err, w)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.ForeignKeyViolation {
			return balance.ErrNoOrder
		}

		return err
	}

	sql = `UPDATE gophermart.balances SET current=current-$1 WHERE user_id=$2`

	if _, err := tx.ExecContext(ctx, sql, w.Sum, w.UserID); err != nil {
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
