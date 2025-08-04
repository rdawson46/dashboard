package db

import (
    "context"
)

type Repository interface {
    // User 
    GetUser(ctx context.Context, id int64) (*User_db, error)
    GetUsers(ctx context.Context, limit, offset int64) ([]*User_db, error)
    GetUserCount(ctx context.Context) (int64, error)

    CreateUser() ()
    UpdateUser() ()


    // Messages
    GetMessage() ()
    GetMessages() ()
    GetMessageCount() ()

    CreateMessage() ()
    UpdateMessage() ()

    GetDescriptions() ()
	Close() error
}

