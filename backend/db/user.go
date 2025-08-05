package db

import (
    "context"
    "database/sql"
    "errors"
    "fmt"
    "time"

	"golang.org/x/crypto/bcrypt"
    _ "modernc.org/sqlite"
)

// TODO: add logging and make theses tables

type User_db struct {
    ID int64 `json:"id"`
    Name string `json:"name"`
    CreatedAt time.Time `json:"createdAt"`
}

// HACK: can most likely remove this, can keep for logging
type UserErrorResponse struct {
    Error string `json:"error"`
    Code string `json:"code,omitempty"`
    Details string `json:"details,omitempty"`
}

var (
    ErrUserNotFound = errors.New("user not found")
    ErrInvalidUserId = errors.New("invalid user ID")
)

func (r *sqliteRepo) GetUser(ctx context.Context, id int64) (*User_db, error) {
    query := `SELECT id, name, created_at FROM users WHERE id = ?`

    row := r.db.QueryRowContext(ctx, query, id)

    var user User_db
    err := row.Scan(&user.ID, &user.Name, &user.CreatedAt)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, ErrUserNotFound
        } 
        return nil, fmt.Errorf("failed to scan user: %w", err)
    }

    return &user, nil
}


func (r *sqliteRepo) GetUsers(ctx context.Context, limit, offset int64) ([]*User_db, error){
    query := `SELECT id, name, created_at FROM users ORDER BY created_at DESC LIMIT ? OFFSET ?`

    rows, err := r.db.QueryContext(ctx, query, limit, offset)

    if err != nil {
        return nil, fmt.Errorf("failed to query users: %w", err)
    }
    defer rows.Close()

    var users []*User_db
    for rows.Next() {
        var user User_db
        err := rows.Scan(&user.ID, &user.Name, &user.CreatedAt)

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


func (r *sqliteRepo) GetUserCount(ctx context.Context) (int64, error) {
    query := `SELECT COUNT(*) FROM users`

    var count int64
    err := r.db.QueryRowContext(ctx, query).Scan(&count)

    if err != nil {
        return 0, fmt.Errorf("failed to count users: %w", err)
    }

    return count, nil
}


func (r *sqliteRepo) CreateUser(ctx context.Context, username, password string) (*User_db, error) {
	if !checkPassword(password) {
		return nil, errors.New("Invalid password")
	}
	
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	query := `INSERT INTO users (username, password) VALUES (?, ?)`

	result, err := r.db.Exec(query, username, hashedPass)

	if err != nil {
		return nil, err
	}

	insertedId, err := result.LastInsertId()

	if err != nil {
		return nil, err
	}

	var insertedUser User_db
	lastQuery := `SELECT id, username, created_at FROM users WHERE id = ?`
	err = r.db.QueryRowContext(ctx, lastQuery, insertedId).Scan(&insertedUser.ID, &insertedUser.Name, &insertedUser.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &insertedUser, nil
}


// TODO:
func (r *sqliteRepo) UpdateUser() () {}


func (r *sqliteRepo) SignInUser(ctx context.Context, username, enteredPassword string) (*User_db, error) {
	return nil, nil
}
