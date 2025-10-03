package db

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	ollama "github.com/ollama/ollama/api"
	_ "modernc.org/sqlite"
)

var getChatbyIdQuery = `SELECT messages FROM messages WHERE id = ? AND user_id = ?`
var createMessageQuery = `INSERT INTO messages (id, user_id, messages, description) VALUES (?, ?, ?, ?)`
var getDescriptionsQuery = `SELECT id, description FROM messages WHERE user_id = ? ORDER BY created_at LIMIT ? OFFSET ?`
var deleteMessageQuery = `DELETE FROM messages WHERE id = ? AND user_id = ?`
var updateMessageQuery = `UPDATE messages SET messages = ? WHERE id = ? AND user_id = ?`

type Message_db struct {
	Id       string
	Messages []ollama.Message
}

// TODO:
type MessageErrorResponse struct {
}

type ChatDesc struct {
	Id          string `json:"id"`
	Description string `json:"description"`
}

type Descriptions []*ChatDesc

var (
	ErrMessageNotFound = errors.New("user not found")
	ErrInvalidMessage  = errors.New("invalid user ID")
)

func (r *sqliteRepo) GetMessage(ctx context.Context, userId, messageId string) ([]ollama.Message, error) {
	var messagesStr string
	err := r.db.QueryRowContext(ctx, getChatbyIdQuery, messageId, userId).Scan(&messagesStr)

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

func (r *sqliteRepo) GetMessages()    {}
func (r *sqliteRepo) GetMessageCount() {}
func (r *sqliteRepo) UpdateMessage()  {}

// TODO: test this out
func (r *sqliteRepo) CreateMessage(ctx context.Context, userId string, messages []ollama.Message) (string, error) {
	/*
		1. have to generate the description
		2. insert into db
		3. messages to string
		4. return message ID
	*/

	desc := generateDesc(messages)

	messageString, err := json.Marshal(messages)

	if err != nil {
		return "", errors.New("Couldn't marshal messages")
	}

	messageID := uuid.New().String()

	_, err = r.db.Exec(createMessageQuery, messageID, userId, string(messageString), desc)

	if err != nil {
		return "", err
	}

	return messageID, nil
}

func (r *sqliteRepo) GetDescriptions(ctx context.Context, userId string, limit, offset int) (Descriptions, error) {
	rows, err := r.db.QueryContext(ctx, getDescriptionsQuery, userId, limit, offset)

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

func (r *sqliteRepo) DeleteMessage(ctx context.Context, id string, user_id string) (bool, error) {
	result, err := r.db.Exec(deleteMessageQuery, id, user_id)

	if err != nil {
		return false, err
	}

	i, err := result.RowsAffected()

	if i <= 0 {
		return false, nil
	}

	return true, nil
}

// HACK: returning bool as a temp placeholder
func (r *sqliteRepo) AddMessage(ctx context.Context, messageId, userId string, messages []ollama.Message) (bool, error) {
	/*
	 **kind of a terrible impl**
	*/

	messageString, err := json.Marshal(messages)

	if err != nil {
		return false, err
	}

	result, err := r.db.Exec(updateMessageQuery, string(messageString), messageId, userId)

	if err != nil {
		return false, err
	}

	i, err := result.RowsAffected()

	if err != nil {
		return false, err
	}

	if i == 0 {
		return false, nil
	}

	return true, nil
}
