package pg

import (
	"database/sql"

	"github.com/Fe4p3b/gophermart/internal/model"
	"github.com/Fe4p3b/gophermart/internal/storage"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var _ storage.AuthStorer = (*pg)(nil)

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
	return nil
}
