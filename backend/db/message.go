package db

import (
    "errors"

    _ "modernc.org/sqlite"
)

type Message_db struct {
}

type MessageErrorResponse struct {
}

var (
    ErrMessageNotFound = errors.New("user not found")
    ErrInvalidMessage = errors.New("invalid user ID")
)

// TODO:
func (r *sqliteRepo) GetMessage() () 
func (r *sqliteRepo) GetMessages() ()
func (r *sqliteRepo) GetMessageCount() ()
func (r *sqliteRepo) CreateMessage() ()
func (r *sqliteRepo) UpdateMessage() ()
func (r *sqliteRepo) GetDescriptions() ()
