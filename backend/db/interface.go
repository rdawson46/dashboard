package db

import (
    "context"
	ollama "github.com/ollama/ollama/api"
)

type Repository interface {
    MessageDB 
    UserDB
    JobDB
	Close() error
}

type MessageDB interface {
    GetMessage(ctx context.Context, userId, messageId string) ([]ollama.Message, error)
    DeleteMessage(ctx context.Context, id string, user_id string) (bool, error)
    GetMessages() ()
    GetMessageCount() ()

    CreateMessage(ctx context.Context, userId string, message []ollama.Message) (string, error)
    UpdateMessage() ()
	AddMessage(ctx context.Context, messageId, userId string, message []ollama.Message) (bool, error)

    GetDescriptions(ctx context.Context, userId string, limit, offset int) (Descriptions, error)
}

type UserDB interface {
    GetUser(ctx context.Context, id string) (*User_db, error)
    GetUsers(ctx context.Context, limit, offset int64) ([]*User_db, error)
    GetUserCount(ctx context.Context) (int64, error)

    CreateUser(ctx context.Context, username, password string) (*User_db, error)
    SignInUser(ctx context.Context, username, password string) (*User_db, error)
    UpdateUser() ()

	GetPerferredModel(ctx context.Context, userId string) (string, error)
	SetPerferredModel(ctx context.Context, userId, model string) (error)
}

type JobDB interface {
	CreateJob(ctx context.Context) ()
	UpdateJob(ctx context.Context) ()
	GetJob(ctx context.Context) ()

	Peek(ctx context.Context) ()
	Run(ctx context.Context) ()
}
