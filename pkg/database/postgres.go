package database

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type PostgresConfig struct {
	DSN string
}

type PostgresPool struct {
	*sql.DB
}

func NewPostgresPool(c PostgresConfig) (*PostgresPool, error) {
	db, err := sql.Open("postgres", c.DSN)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresPool{db}, nil
}
