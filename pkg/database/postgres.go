package database

import "database/sql"

type PostgresConfig struct {
	DSN string
}

type PostgresPool struct {
	*sql.DB
}

func NewPostgres(c PostgresConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", c.DSN)
	if err != nil {
		return nil, err
	}
	return db, nil
}
