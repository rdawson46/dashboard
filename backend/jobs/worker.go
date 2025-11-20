package jobs

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	"github.com/charmbracelet/log"
)

const (
	updateStatusQuery = `UPDATE jobs SET status = ? WHERE id = ?`
	saveResultQuery = `UPDATE jobs SET result = ? WHERE id = ?`
)

func updateJobResult(jobId string, result string, db *sql.DB, logger *log.Logger) {
	res, err := db.Exec(saveResultQuery, result, jobId)

	if err != nil {
		logger.Info("Error updating job result", "id", jobId, "err", err.Error())
		return
	}

	i, err := res.RowsAffected()

	if err != nil {
		logger.Info("Error getting job updated count", "id", jobId, "err", err.Error())
		return
	}

	logger.Info("Updated job result", "id", jobId, "rows_affected", i)
}

func updateJobStatus(jobId string, status Status, db *sql.DB, logger *log.Logger) {
	if !status.isValidStatus() {
		return
	}

	res, err := db.Exec(updateStatusQuery, status, jobId)

	if err != nil {
		logger.Info("Error updating job status", "id", jobId, "err", err.Error())
		return
	}

	i, err := res.RowsAffected()

	if err != nil {
		logger.Info("Error getting job updated count", "id", jobId, "err", err.Error())
		return
	}

	logger.Info("Updated job stauts", "id", jobId, "rows_affected", i)
}

func StartWorkerPool(ctx context.Context, jobs <-chan *Job, count int, logger *log.Logger, db *sql.DB) {
	var wg sync.WaitGroup
	for i := range count {
		wg.Add(1)
        l := logger.WithPrefix(fmt.Sprintf("Worker (%d)", i))

		go func(id int, logger *log.Logger) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case job := <-jobs:
					l.Info("Job recieved", "id", job.Id)
					taskResult := job.Run()

					updateJobResult(job.Id, job.Result, db, l)

					// TODO: update status message to completed or failed from job
					updateJobStatus(job.Id, taskResult.status, db, l)
				}
			}
        }(i, l)
	}
	wg.Wait()
}
