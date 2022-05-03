package pg

import (
	"context"
	"time"

	"github.com/Fe4p3b/gophermart/internal/model"
	"github.com/Fe4p3b/gophermart/internal/storage"
)

var (
	_ storage.BalanceRepository = (*BalanceStorage)(nil)
)

type BalanceStorage struct {
	pg *pg
}

func NewBalanceStorage(pg *pg) *BalanceStorage {
	return &BalanceStorage{pg: pg}
}

func (bs *BalanceStorage) GetBalanceForUser(u string) (*model.Balance, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sql := `SELECT b.id, b.user_id, b.current, COALESCE(SUM(w.sum),0) as withdrawn
FROM gophermart.balances as b
LEFT JOIN gophermart.withdrawals as w
ON b.user_id = w.user_id
WHERE b.user_id = $1
GROUP BY b.id`
	row := bs.pg.db.QueryRowContext(ctx, sql, u)

	var balance model.Balance
	if err := row.Scan(&balance.ID, &balance.UserID, &balance.Current, &balance.Withdrawn); err != nil {
		return nil, err
	}

	return &balance, nil
}
