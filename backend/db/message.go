package db

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	ollama "github.com/ollama/ollama/api"
	_ "modernc.org/sqlite"
)

var (
	getChatbyIdQuery = `SELECT messages FROM messages WHERE id = ? AND user_id = ?`
	createMessageQuery = `INSERT INTO messages (id, user_id, messages, description, model) VALUES (?, ?, ?, ?, ?)`
	getDescriptionsQuery = `SELECT id, description FROM messages WHERE user_id = ? ORDER BY updated_at DESC LIMIT ? OFFSET ?`
	deleteMessageQuery = `DELETE FROM messages WHERE id = ? AND user_id = ?`
	updateMessageQuery = `UPDATE messages SET messages = ?, model = ? WHERE id = ? AND user_id = ?`
)

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
	r.logger.Info("Fetching message", "messageId", messageId, "userId", userId)
	err := r.Db.QueryRowContext(ctx, getChatbyIdQuery, messageId, userId).Scan(&messagesStr)

	if err != nil {
		r.logger.Error("Error fetching message", "messageId", messageId, "userId", userId)
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

func (r *sqliteRepo) CreateMessage(ctx context.Context, userId, model string, messages []ollama.Message) (string, error) {
	r.logger.Info("Creating message", "userId", userId)

	desc := generateDesc(messages)

	messageString, err := json.Marshal(messages)

	if err != nil {
		r.logger.Error("failed to marshal messages", "userId", userId)
		return "", errors.New("Couldn't marshal messages")
	}

	messageID := uuid.New().String()

	_, err = r.Db.Exec(createMessageQuery, messageID, userId, string(messageString), desc, model)

	if err != nil {
		r.logger.Error("failed to create message", "messageId", messageID, "userId", userId)
		return "", err
	}

	return messageID, nil
}

func (r *sqliteRepo) GetDescriptions(ctx context.Context, userId string, limit, offset int) (Descriptions, error) {
	r.logger.Info("Fetching descriptions", "userId", userId, "limit", limit, "offset", offset)
	rows, err := r.Db.QueryContext(ctx, getDescriptionsQuery, userId, limit, offset)

	if err != nil {
		r.logger.Error("failed to query descriptions", "userId", userId, "limit", limit, "offset", offset)
		return nil, err
	}

	defer rows.Close()

	var descs Descriptions
	for rows.Next() {
		var d ChatDesc

		err := rows.Scan(&d.Id, &d.Description)

		if err != nil {
			r.logger.Error("failed to scan description", "userId", userId)
			return nil, err
		}

		descs = append(descs, &d)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("error iterating rows", "userId", userId)
		return nil, err
	}

	return descs, nil
}

func (r *sqliteRepo) DeleteMessage(ctx context.Context, id string, user_id string) (bool, error) {
	r.logger.Info("Deleting message", "messageId", id, "userId", user_id)
	result, err := r.Db.Exec(deleteMessageQuery, id, user_id)

	if err != nil {
		r.logger.Error("failed to delete message", "messageId", id, "userId", user_id)
		return false, err
	}

	i, err := result.RowsAffected()

	if i <= 0 {
		return false, nil
	}

	return true, nil
}

func (r *sqliteRepo) AddMessage(ctx context.Context, messageId, userId, model string, messages []ollama.Message) (bool, error) {
	r.logger.Info("Adding message", "messageId", messageId, "userId", userId)
	/*
	 **kind of a terrible impl**
	*/

	messageString, err := json.Marshal(messages)

	if err != nil {
		r.logger.Error("failed to marshal messages", "messageId", messageId, "userId", userId)
		return false, err
	}

	result, err := r.Db.Exec(updateMessageQuery, string(messageString), model, messageId, userId)

	if err != nil {
		r.logger.Error("failed to update message", "messageId", messageId, "userId", userId)
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
