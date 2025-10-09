package jobs

import (
	"context"
	"time"
    "database/sql"
)


func FindReadyJobs(db *sql.DB) (Jobs, error) {
	return nil, nil
}


func StartScheduler(ctx context.Context, db *sql.DB, jobs chan<- Job) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <- ticker.C:
			readyJobs, err := FindReadyJobs(db)
			if err != nil {
				// TODO: add logger and log this error
			}

			for _, job := range readyJobs {
				jobs <- job
			}
		}
	}
}
