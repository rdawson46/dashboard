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
	"github.com/rdawson46/dashboard/jobs"
)

func (r *sqliteRepo) CreateJob(ctx context.Context, job jobs.Job) (*jobs.Job, error) {
	jobId := uuid.New().String()

	return nil, nil
}

func (r *sqliteRepo) UpdateJob(ctx context.Context) ()

func (r *sqliteRepo) GetJob(ctx context.Context) ()

func (r *sqliteRepo) Peek(ctx context.Context) ()
func (r *sqliteRepo) Run(ctx context.Context) ()
