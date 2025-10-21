package jobs

import (
	"context"
	"database/sql"

	"github.com/charmbracelet/log"
)

type JobPipeline struct {
    logger *log.Logger
    db *sql.DB
}

func NewJobPipeline(logger *log.Logger, db *sql.DB) JobPipeline {
    return JobPipeline{
        logger: logger.WithPrefix("Scheduler"),
        db: db,
    }
}

// TODO: don't use background context, use one with cancel
func (jp *JobPipeline) StartJobCheck(ctx context.Context, workers int) {
    j := make(chan *Job, 10)
    worker_logger := jp.logger.WithPrefix("Worker")

    go StartWorkerPool(ctx, j, workers, worker_logger)
    jp.StartScheduler(ctx, j)
}
