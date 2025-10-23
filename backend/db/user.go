package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

var (
	getUserQuery = `SELECT id, username, created_at FROM users WHERE id = ?`
	getUsersQuery = `SELECT id, username, created_at FROM users ORDER BY created_at DESC LIMIT ? OFFSET ?`
	getUserCountQuery = `SELECT COUNT(*) FROM users`
	createUserQuery = `INSERT INTO users (id, username, password) VALUES (?, ?, ?)`
	signInUserQuery = `SELECT id, username, created_at, password FROM users WHERE username = ?`
	getModelQuery = `SELECT model FROM users WHERE id = ?`
)

type User_db struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"createdAt"`
}

type UserErrorResponse struct {
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
}

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrInvalidUserId = errors.New("invalid user ID")
)

func (r *sqliteRepo) GetUser(ctx context.Context, id string) (*User_db, error) {
	r.logger.Info("Fetching user", "userId", id)
	row := r.db.QueryRowContext(ctx, getUserQuery, id)

	var user User_db
	err := row.Scan(&user.ID, &user.Username, &user.CreatedAt)
	if err != nil {
		r.logger.Error("failed to get user", "userId", id)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to scan user: %w", err)
	}

	return &user, nil
}

func (r *sqliteRepo) GetUsers(ctx context.Context, limit, offset int64) ([]*User_db, error) {
	r.logger.Info("Fetching users", "limit", limit, "offset", offset)
	rows, err := r.db.QueryContext(ctx, getUsersQuery, limit, offset)

	if err != nil {
		r.logger.Error("failed to query users", "limit", limit, "offset", offset)
		return nil, err
	}
	defer rows.Close()

	var users []*User_db
	for rows.Next() {
		var user User_db
		err := rows.Scan(&user.ID, &user.Username, &user.CreatedAt)

		if err != nil {
			r.logger.Error("failed to scan user")
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}

		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("error iterating rows")
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return users, nil
}

func (r *sqliteRepo) GetUserCount(ctx context.Context) (int64, error) {
	r.logger.Info("Getting user count")
	var count int64
	err := r.db.QueryRowContext(ctx, getUserCountQuery).Scan(&count)

	if err != nil {
		r.logger.Error("failed to count users")
		return 0, fmt.Errorf("failed to count users: %w", err)
	}

	return count, nil
}

func (r *sqliteRepo) CreateUser(ctx context.Context, username, password string) (*User_db, error) {
	r.logger.Info("Creating user", "username", username)
	if !checkPassword(password) {
		r.logger.Error("invalid password")
		return nil, errors.New("Invalid password")
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		r.logger.Error("failed to hash password")
		return nil, err
	}

	userID := uuid.New().String()

	_, err = r.db.Exec(createUserQuery, userID, username, hashedPass)

	if err != nil {
		r.logger.Error("failed to create user", "userId", userID, "username", username)
		return nil, err
	}

	var insertedUser User_db
	err = r.db.QueryRowContext(ctx, getUserQuery, userID).Scan(&insertedUser.ID, &insertedUser.Username, &insertedUser.CreatedAt)

	if err != nil {
		r.logger.Error("failed to fetch created user", "userId", userID, "username", username)
		return nil, err
	}

	return &insertedUser, nil
}

func (r *sqliteRepo) UpdateUser() () {}

func (r *sqliteRepo) SignInUser(ctx context.Context, username, enteredPassword string) (*User_db, error) {
	r.logger.Info("Signing in user", "username", username)
	row := r.db.QueryRowContext(ctx, signInUserQuery, username)

	var user User_db
	var hashedPass string
	err := row.Scan(&user.ID, &user.Username, &user.CreatedAt, &hashedPass)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.Error("user not found", "username", username)
			return nil, ErrUserNotFound
		}
		r.logger.Error("failed to scan user", "username", username)
		return nil, fmt.Errorf("failed to scan user: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(enteredPassword))

	if err != nil {
		r.logger.Error("incorrect password", "username", username)
		return nil, fmt.Errorf("Incorrect Password")
	}

	return &user, nil
}


func (r *sqliteRepo) GetPerferredModel(ctx context.Context, userId string) (string, error) {
	r.logger.Info("Getting preferred model", "userId", userId)
	var modelStr string
	err := r.db.QueryRowContext(ctx, getModelQuery, userId).Scan(&modelStr)

	if err != nil {
		r.logger.Error("failed to get preferred model", "userId", userId)
		return "", err
	}

	return modelStr, nil
}

func (r *sqliteRepo) SetPerferredModel(ctx context.Context, userId, model string) error {
	/*
	r.logger.Info("Setting preferred model", "userId", userId, "model", model)
	result, err := r.db.Exec(setModelQuery, model, userId)

	if err != nil {
		r.logger.Error("failed to set preferred model", "userId", userId, "model", model)
		return nil
	}

	i, err := result.RowsAffected()

	if err != nil {
		r.logger.Error("failed to get rows affected", "userId", userId, "model", model)
		return err
	}

	if i == 0 {
		return nil
	}

	return nil
	*/
	return nil
}
