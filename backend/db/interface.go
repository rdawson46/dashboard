package db

import (
    "context"
	ollama "github.com/ollama/ollama/api"
)

type Repository interface {
    // User 
    GetUser(ctx context.Context, id int64) (*User_db, error)
    GetUsers(ctx context.Context, limit, offset int64) ([]*User_db, error)
    GetUserCount(ctx context.Context) (int64, error)

    CreateUser(ctx context.Context, username, password string) (*User_db, error)
    SignInUser(ctx context.Context, username, password string) (*User_db, error)
    UpdateUser() ()


    // Messages
    GetMessage(ctx context.Context, id int64) ([]ollama.Message, error)
    DeleteMessage(ctx context.Context, id int64, user_id int64) (bool, error)
    GetMessages() ()
    GetMessageCount() ()

    CreateMessage(ctx context.Context, userId int64, message []ollama.Message) (int64, error)
    UpdateMessage() ()
	AddMessage(ctx context.Context, messageId, userId int64, message []ollama.Message) (bool, error)

    GetDescriptions(ctx context.Context, userId int64, limit, offset int) (Descriptions, error)
	Close() error
}

