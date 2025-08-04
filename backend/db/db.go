package db

import (
    "database/sql"
    "errors"
	"os"

    _ "modernc.org/sqlite"
)

var ErrDatabaseError = errors.New("database error")

type sqliteRepo struct {
    db *sql.DB
}

func NewSqliteRepository() (Repository, error) {
	db_url := os.Getenv("DB_URL")

	if db_url == "" {
		return nil, errors.New("No DB URL found in env")
	}

	db, err := sql.Open("sqlite", db_url)

	if err != nil {
		return nil, err
	}

    return &sqliteRepo{db: db}, nil
}

func (r *sqliteRepo) Close() error {
	return r.db.Close()
}
