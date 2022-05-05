package pg

import (
	"context"
	"errors"
	"log"

	"github.com/Fe4p3b/gophermart/internal/model"
	"github.com/Fe4p3b/gophermart/internal/service/withdrawal"
	"github.com/Fe4p3b/gophermart/internal/storage"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
)

var (
	_ storage.WithdrawalRepository = (*WithdrawalStorage)(nil)
)

type WithdrawalStorage struct {
	pg *pg
}

func NewWithdrawalStorage(pg *pg) *WithdrawalStorage {
	return &WithdrawalStorage{pg: pg}
}

func (ws *WithdrawalStorage) AddWithdrawal(ctx context.Context, w model.Withdrawal) error {
	tx, err := ws.pg.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	sql := `INSERT INTO gophermart.withdrawals(user_id, order_number, sum, date) VALUES($1, $2, $3, $4)`

	if _, err := tx.ExecContext(ctx, sql, w.UserID, w.OrderNumber, w.Sum, w.Date); err != nil {
		log.Printf("AddWithdrawal - %v, %v", err, w)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.ForeignKeyViolation {
			return withdrawal.ErrNoOrder
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

func (ws *WithdrawalStorage) GetWithdrawalsForUser(ctx context.Context, u string) ([]model.Withdrawal, error) {
	sql := `SELECT id, order_number, sum, date FROM gophermart.withdrawals WHERE user_id = $1`

	rows, err := ws.pg.db.QueryContext(ctx, sql, u)
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
