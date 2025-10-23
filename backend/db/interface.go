package db

import (
    "context"
	jobs "github.com/rdawson46/dashboard/jobs"
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

    CreateMessage(ctx context.Context, userId, model string, message []ollama.Message) (string, error)
    UpdateMessage() ()
	AddMessage(ctx context.Context, messageId, userId, model string, message []ollama.Message) (bool, error)

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
    CreateJob(ctx context.Context, userId string, job jobs.Job) (*jobs.Job, error)
    GetJob(ctx context.Context, jobId string, userId string) (*jobs.Job, error)
    GetJobs(ctx context.Context, userId string, limit, offset int) (jobs.Jobs, error)
    UpdateJob(ctx context.Context, job jobs.Job) (*jobs.Job, error)
    DeleteJob(ctx context.Context, jobId string, userId string) error
    Peek(ctx context.Context) (*jobs.Job, error)
}
