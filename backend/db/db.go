package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/charmbracelet/log"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "modernc.org/sqlite"
)

var ErrDatabaseError = errors.New("database error")

type sqliteRepo struct {
    Db *sql.DB
	logger *log.Logger
	mu sync.Mutex
}

func NewSqliteRepository(log *log.Logger) (*sqliteRepo, error) {
	l := log.WithPrefix("")
	db_url := os.Getenv("DB_URL")

	if db_url == "" {
		return nil, errors.New("No DB URL found in env")
	}

	db, err := sql.Open("sqlite", db_url)

	if err != nil {
		return nil, err
	}

    // migrations
    m, err := migrate.New(
        "file://db/migrations", // TODO: move migrations dir and name files
        fmt.Sprintf("sqlite://%s", db_url),
    )

	if err != nil {
		return nil, err
	}

    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, err
    }

	l.Info("DB initiated")

    return &sqliteRepo{
		Db: db,
		logger: l,
	}, nil
}

func (r *sqliteRepo) Close() error {
	r.logger.Info("Closing DB")
	return r.Db.Close()
}
