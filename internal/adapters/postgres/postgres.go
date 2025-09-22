package postgres

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type Config struct {
	DSN string
}

type Pool struct {
	*sql.DB
}

func NewPostgresPool(c Config) (*Pool, error) {
	db, err := sql.Open("postgres", c.DSN)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &Pool{db}, nil
}
