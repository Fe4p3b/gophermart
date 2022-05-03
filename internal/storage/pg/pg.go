package pg

import (
	"context"
	"database/sql"
	"os"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
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
