package pg

import (
	"context"
	"errors"

	"github.com/Fe4p3b/gophermart/internal/model"
	"github.com/Fe4p3b/gophermart/internal/service/order"
	"github.com/Fe4p3b/gophermart/internal/storage"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
)

var (
	_ storage.OrderRepository = (*OrderStorage)(nil)
)

type OrderStorage struct {
	pg *pg
}

func NewOrderStorage(pg *pg) *OrderStorage {
	return &OrderStorage{pg: pg}
}

func (os *OrderStorage) GetOrdersForUser(ctx context.Context, u string) ([]model.Order, error) {
	sql := `SELECT number, user_id, number, status, accrual, upload_date FROM gophermart.orders WHERE user_id = $1`

	rows, err := os.pg.db.QueryContext(ctx, sql, u)
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

func (os *OrderStorage) AddOrder(ctx context.Context, o *model.Order) error {
	sql := `INSERT INTO gophermart.orders(user_id, number, status, accrual, upload_date) VALUES ($1, $2, $3, $4, $5)`
	if _, err := os.pg.db.ExecContext(ctx, sql, o.UserID, o.Number, o.Status.String(), o.Accrual, o.UploadDate); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			sql := `SELECT user_id FROM gophermart.orders WHERE number=$1`
			row := os.pg.db.QueryRowContext(ctx, sql, o.Number)

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

func (os *OrderStorage) UpdateOrder(ctx context.Context, o *model.Order) error {
	sql := `UPDATE gophermart.orders SET status = $1, accrual = $2 WHERE number = $3`

	if _, err := os.pg.db.ExecContext(ctx, sql, o.Status, o.Accrual, o.Number); err != nil {
		return err
	}
	return nil
}

func (os *OrderStorage) UpdateBalanceForProcessedOrder(ctx context.Context, o *model.Order) error {
	sql := `UPDATE gophermart.balances SET current = current + $1 WHERE user_id = $2`
	if _, err := os.pg.db.ExecContext(ctx, sql, o.Accrual, o.UserID); err != nil {
		return err
	}

	return nil
}
