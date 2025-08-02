package server

import (
    "context"
    "database/sql"
    "errors"
    "fmt"
	"os"
    "time"

	// "github.com/charmbracelet/log"
    _ "modernc.org/sqlite"
)

// TODO: change the logging and make theses tables

type User_db struct {
    ID int64 `json:"id"`
    Name string `json:"name"`
    Email string `json:"email"`
    CreatedAt time.Time `json:"createdAt"`
}

type UserErrorResponse struct {
    Error string `json:"error"`
    Code string `json:"code,omitempty"`
    Details string `json:"details,omitempty"`
}

var (
    ErrUserNotFound = errors.New("user not found")
    ErrInvalidUserId = errors.New("invalid user ID")
    ErrDatabaseError = errors.New("database error")
)

type UserRepository interface {
    GetUser(ctx context.Context, id int64) (*User_db, error)
    GetUsers(ctx context.Context, limit, offset int64) ([]*User_db, error)
    GetUserCount(ctx context.Context) (int64, error)
	Close() error
}

type sqliteUserRepo struct {
    db *sql.DB
}

func NewSqliteRepository() (UserRepository, error) {
	db_url := os.Getenv("DB_URL")

	if db_url == "" {
		return nil, errors.New("No DB URL found in env")
	}

	db, err := sql.Open("sqlite", db_url)

	if err != nil {
		return nil, err
	}

    return &sqliteUserRepo{db: db}, nil
}

func (r *sqliteUserRepo) GetUser(ctx context.Context, id int64) (*User_db, error) {
    query := `SELECT id, name, email, created_at FROM users WHERE id = ?`

    row := r.db.QueryRowContext(ctx, query, id)

    var user User_db
    err := row.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, ErrUserNotFound
        } 
        return nil, fmt.Errorf("failed to scan user: %w", err)
    }

    return &user, nil
}

func (r *sqliteUserRepo) GetUsers(ctx context.Context, limit, offset int64) ([]*User_db, error){
    query := `SELECT id, name, email, created_at FROM users ORDER BY created_at DESC LIMIT ? OFFSET ?`

    rows, err := r.db.QueryContext(ctx, query, limit, offset)

    if err != nil {
        return nil, fmt.Errorf("failed to query users: %w", err)
    }
    defer rows.Close()

    var users []*User_db
    for rows.Next() {
        var user User_db
        err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)

        if err != nil {
            return nil, fmt.Errorf("failed to scan user: %w", err)
        } 

        users = append(users, &user)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("error iterating rows: %w", err)
    }

    return users, nil
}

func (r *sqliteUserRepo) GetUserCount(ctx context.Context) (int64, error) {
    query := `SELECT COUNT(*) FROM users`

    var count int64
    err := r.db.QueryRowContext(ctx, query).Scan(&count)

    if err != nil {
        return 0, fmt.Errorf("failed to count users: %w", err)
    }

    return count, nil
}

func (r *sqliteUserRepo) Close() error {
	return r.db.Close()
}
