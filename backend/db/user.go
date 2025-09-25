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


var getUserQuery = `SELECT id, name, created_at FROM users WHERE id = ?`
var getUsersQuery = `SELECT id, name, created_at FROM users ORDER BY created_at DESC LIMIT ? OFFSET ?`
var getUserCountQuery = `SELECT COUNT(*) FROM users`
var createUserQuery = `INSERT INTO users (username, password) VALUES (?, ?)`
var signInUserQuery = `SELECT * FROM users WHERE username = ?`
var lastQuery = `SELECT id, username, created_at FROM users WHERE id = ?`


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
    row := r.db.QueryRowContext(ctx, getUserQuery, id)

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
    rows, err := r.db.QueryContext(ctx, getUsersQuery, limit, offset)

    if err != nil {
        return nil, err
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
    var count int64
    err := r.db.QueryRowContext(ctx, getUserCountQuery).Scan(&count)

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

	result, err := r.db.Exec(createMessageQuery, username, hashedPass)

	if err != nil {
		return nil, err
	}

	insertedId, err := result.LastInsertId()

	if err != nil {
		return nil, err
	}

	var insertedUser User_db
	err = r.db.QueryRowContext(ctx, lastQuery, insertedId).Scan(&insertedUser.ID, &insertedUser.Name, &insertedUser.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &insertedUser, nil
}


// TODO:
func (r *sqliteRepo) UpdateUser() () {}


func (r *sqliteRepo) SignInUser(ctx context.Context, username, enteredPassword string) (*User_db, error) {
    row := r.db.QueryRowContext(ctx, signInUserQuery, username)

    var user User_db
	var hashedPass string
    err := row.Scan(&user.ID, &user.Name, &user.CreatedAt, &hashedPass)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, ErrUserNotFound
        } 
        return nil, fmt.Errorf("failed to scan user: %w", err)
    }

	err = bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(enteredPassword))

	if err != nil {
		return nil, fmt.Errorf("Incorrect Password")
	}

	return &user, nil
}
