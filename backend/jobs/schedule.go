package jobs

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

/*
Options:
 - 15mins
 - 30mins
 - 1hour
*/

func FindReadyJobs(db *sql.DB) (Jobs, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback() // Rollback on any error.

	// TODO: will have to rework so that jobs will repeat
	query := fmt.Sprintf(`SELECT id, name, user_id, task, model, status, time, freq, result FROM jobs 
	WHERE status = '%s' AND (
		(freq = '15mins' AND datetime(time, '+15 minutes') <= CURRENT_TIMESTAMP) OR
		(freq = '30mins' AND datetime(time, '+30 minutes') <= CURRENT_TIMESTAMP) OR
		(freq = '1hour' AND datetime(time, '+1 hour') <= CURRENT_TIMESTAMP)
	) ORDER BY created_at ASC`, StatusPending)
	rows, err := tx.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobList Jobs
	var jobIds []string
	for rows.Next() {
		var job Job
		var task, userId string

		err := rows.Scan(&job.Id, &job.Name, &userId, &task, &job.Model, &job.Status, &job.Time, &job.Freq, &job.Result)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal([]byte(task), &job.Task)
		if err != nil {
			return nil, err
		}

		jobList = append(jobList, &job)
		jobIds = append(jobIds, job.Id)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// if we don't have any jobs, we can commit and return
	if len(jobIds) == 0 {
		return nil, tx.Commit()
	}

	updateQuery := "UPDATE jobs SET status = ? WHERE id IN (?" + strings.Repeat(",?", len(jobIds)-1) + ")"
	args := make([]any, len(jobIds)+1)
	args[0] = StatusRunning
	for i, id := range jobIds {
		args[i+1] = id
	}

	_, err = tx.Exec(updateQuery, args...)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return jobList, nil
}


func (jp *JobPipeline) StartScheduler(ctx context.Context, jobs chan<- *Job) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <- ticker.C:
			readyJobs, err := FindReadyJobs(jp.db)
			if err != nil {
				jp.logger.Error(
					"Error searching for ready jobs",
					"err", err.Error(),
				)

				continue
			}

			count := len(readyJobs)

			if count > 0 {
				jp.logger.Info("Found ready jobs", "count", len(readyJobs))
			} else {
				jp.logger.Info("No ready jobs found")
			}

			for _, job := range readyJobs {
				jobs <- job
			}
		}
	}
}
