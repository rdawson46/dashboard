package server

import (
    "context"
    "database/sql"
    "encoding/json"
    "errors"
    "fmt"
    "net/http"
    "strconv"
    "time"

	"github.com/charmbracelet/log"
    _ "modernc.org/sqlite"
)

// TODO: change the logging and make theses tables

type User struct {
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
    GetUser(ctx context.Context, id int64) (*User, error)
    GetUsers(ctx context.Context, limit, offset int64) ([]*User, error)
    GetUserCount(ctx context.Context) (int64, error)
}

type sqliteUserRepo struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
    return &sqliteUserRepo{db: db}
}

func (r *sqliteUserRepo) GetUser(ctx context.Context, id int64) (*User, error) {
    query := `SELECT id, name, email, created_at FROM users WHERE id = ?`

    row := r.db.QueryRowContext(ctx, query, id)

    var user User
    err := row.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, ErrUserNotFound
        } 
        return nil, fmt.Errorf("failed to scan user: %w", err)
    }

    return &user, nil
}

func (r *sqliteUserRepo) GetUsers(ctx context.Context, limit, offset int64) ([]*User, error){
    query := `SELECT id, name, email, created_at FROM users ORDER BY created_at DESC LIMIT ? OFFSET ?`

    rows, err := r.db.QueryContext(ctx, query, limit, offset)

    if err != nil {
        return nil, fmt.Errorf("failed to query users: %w", err)
    }
    defer rows.Close()

    var users []*User
    for rows.Next() {
        var user User
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
