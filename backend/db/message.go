package db

import (
	"context"
	"encoding/json"
	"errors"

	ollama "github.com/ollama/ollama/api"
	_ "modernc.org/sqlite"
)

var getChatbyIdQuery = `SELECT messages FROM messages WHERE id = ? AND user_id = ?`
var createMessageQuery = `INSERT INTO messages (user_id, messages, description) VALUES (?, ?, ?)`
var getDescriptionsQuery = `SELECT id, description FROM messages WHERE user_id = ? ORDER BY created_at LIMIT ? OFFSET ?`
var deleteMessageQuery = `DELETE FROM messages WHERE id = ? AND user_id = ?`


type Message_db struct {
    Id int64
    Messages []ollama.Message
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

func (r *sqliteRepo) GetMessage(ctx context.Context, userId, messageId int64) ([]ollama.Message, error) {
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

func (r *sqliteRepo) GetMessages() () {}
func (r *sqliteRepo) GetMessageCount() () {}
func (r *sqliteRepo) UpdateMessage() () {}

// TODO: test this out
func (r *sqliteRepo) CreateMessage(ctx context.Context, userId int64, messages []ollama.Message) (int64, error) {
	/*
	1. have to generate the description
	2. insert into db
	3. messages to string
	4. return message ID
	*/

	// TEMP: grab first 10 chars of the first user message
    desc := generateDesc(messages)

	messageString, err := json.Marshal(messages)

	if err != nil {
		return 0, errors.New("Couldn't marshal messages")
	}

	result, err := r.db.Exec(createMessageQuery, userId, string(messageString), desc)

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

func (r *sqliteRepo) DeleteMessage(ctx context.Context, id int64, user_id int64) (bool, error) {
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
func (r *sqliteRepo) AddMessage(ctx context.Context, messageId, userId int64, messages []ollama.Message) (bool, error) {
    /*
     **kind of a terrible impl**
    */

    messageString, err := json.Marshal(messages)

    if err != nil {
        return false, err
    }

    query := `UPDATE messages SET messages = ? WHERE id = ? AND user_id = ?`

    result, err := r.db.Exec(query, string(messageString), messageId, userId)

    if err != nil {
        return false, err
    }

    i, err := result.RowsAffected()

    if err != nil {
        return false, err
    }

    if i == 0 { return false, nil }

	return true, nil
}

func generateDesc(message []ollama.Message) string {
	var lastQ string
	for _, m := range message {
		if m.Role == "user" {
			lastQ = m.Content
			break
		}
	}

    var desc string
    if len(lastQ) >= 10 {
        desc = lastQ[:10]
    } else {
        desc = lastQ
    }

    return desc
}
