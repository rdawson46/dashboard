package jobs

import (
	"context"
	"database/sql"
)

// TODO: don't use background context, use one with cancel
func StartJobCheck(ctx context.Context, db *sql.DB, workers int) {
    j := make(chan Job, 10)
    go StartWorkerPool(ctx, j, workers)
    StartScheduler(ctx, db, j)
}
