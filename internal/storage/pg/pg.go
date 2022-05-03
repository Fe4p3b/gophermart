package pg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type pg struct {
	db *sql.DB
	m  *migrate.Migrate
}

func New(dsn string, folder string) (*pg, error) {
	conn, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	pg := &pg{db: conn}
	if err := pg.InitialMigration(folder); err != nil {
		return nil, err
	}

	m, err := migrate.New(fmt.Sprintf("file://%s", folder), dsn)
	if err != nil {
		return nil, err
	}

	pg.m = m

	if err := pg.MigrationUp(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, err
	}

	return pg, nil
}

func (p *pg) InitialMigration(folder string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	file := fmt.Sprintf("./%s/001_init.sql", folder)
	log.Println(file)
	sql, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	_, err = p.db.ExecContext(ctx, string(sql))
	if err != nil {
		return err
	}

	return nil
}

func (p *pg) MigrationUp() error {
	if err := p.m.Up(); err != nil {
		return err
	}

	return nil
}

func (p *pg) MigrationDown() error {
	if err := p.m.Down(); err != nil {
		return err
	}

	return nil
}
