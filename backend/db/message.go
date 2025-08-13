package db

import (
	"context"
	"encoding/json"
	"errors"

	ollama "github.com/ollama/ollama/api"
	_ "modernc.org/sqlite"
)

// TODO:
type Message_db struct {
}

// TODO:
type MessageErrorResponse struct {
}

type ChatDesc struct {
	Id int64 `json:"id"`
	Description string `json:"description"`
}

type Descriptions []*ChatDesc

var (
    ErrMessageNotFound = errors.New("user not found")
    ErrInvalidMessage = errors.New("invalid user ID")
)

func (r *sqliteRepo) GetMessage(ctx context.Context, id int64) ([]ollama.Message, error) {
	query := `SELECT messages FROM messages WHERE id = ?`

	var messagesStr string
	err := r.db.QueryRowContext(ctx, query, id).Scan(&messagesStr)

	if err != nil {
		return nil, err
	}

	var messages []ollama.Message
	err = json.Unmarshal([]byte(messagesStr), &messages)

	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (r *sqliteRepo) GetMessages() () {}
func (r *sqliteRepo) GetMessageCount() () {}
func (r *sqliteRepo) UpdateMessage() () {}

// TODO: test this out
func (r *sqliteRepo) CreateMessage(ctx context.Context, userId int64, message []ollama.Message) (int64, error) {
	/*
	1. have to generate the description
	2. insert into db
	3. messages to string
	4. return message ID
	*/

	// TEMP: grab first 10 chars of the first user message
	var lastQ string
	for _, m := range message {
		if m.Role == "user" {
			lastQ = m.Content
			break
		}
	}

	desc := lastQ[:10]

	messageString, err := json.Marshal(message)

	if err != nil {
		return 0, errors.New("Couldn't marshal messages")
	}

	query := `INSERT INTO messages (user_id, messages, description) VALUES (?, ?, ?)`

	result, err := r.db.Exec(query, userId, string(messageString), desc)

	if err != nil {
		return 0, err
	}

	insertedId, err := result.LastInsertId()

	if err != nil {
		return 0, err
	}

	return insertedId, nil
}

func (r *sqliteRepo) GetDescriptions(ctx context.Context, userId int64, limit, offset int) (Descriptions, error) {
	query := `SELECT id, description FROM messages WHERE user_id = ? ORDER BY created_at LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, userId, limit, offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var descs Descriptions
	for rows.Next() {
		var d ChatDesc

		err := rows.Scan(&d.Id, &d.Description)

		if err != nil {
			return nil, err
		}
		
		descs = append(descs, &d)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return descs, nil
}

// HACK: returning bool as a temp placeholder
func (r *sqliteRepo) AddMessage(ctx context.Context, messageId, userId int64, message []ollama.Message) (bool, error) {
	return false, nil
}
