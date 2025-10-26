package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/rdawson46/dashboard/jobs"
	_ "modernc.org/sqlite"
)

var (
	createJobQuery = `INSERT INTO jobs (id, name, user_id, task, model, status, time, freq, result) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	getJobQuery = `SELECT id, user_id, task, model, status, time, freq, result FROM jobs WHERE id = ? AND user_id = ?`
	getJobsQuery = `SELECT id, name, user_id, task, model, status, time, freq, result FROM jobs WHERE user_id = ? LIMIT ? OFFSET ?`
	updateJobQuery = `UPDATE jobs SET task = ?, model = ?, status = ?, time = ?, freq = ?, result = ? WHERE id = ?`
	deleteJobQuery = `DELETE FROM jobs WHERE id = ? AND user_id = ?`
	peekJobQuery = `SELECT id, user_id, task, model, status, time, freq, result FROM jobs WHERE status = 'pending' AND datetime(time) <= CURRENT_TIMESTAMP ORDER BY created_at ASC LIMIT 1`
	updateJobStatusQuery = `UPDATE jobs SET status = ? WHERE id = ?`
)

func (r *sqliteRepo) CreateJob(ctx context.Context, userId string, job *jobs.Job) (*jobs.Job, error) {
	jobId := uuid.New().String()
	job.Id = jobId
	r.logger.Info("Creating job", "userId", userId, "jobId", jobId)

	if job.Status == "" {
		job.Status = "pending"
	}

	task, err := json.Marshal(job.Task)
	if err != nil {
		r.logger.Error("failed to marshal task", "jobId", jobId, "userId", userId)
		return nil, err
	}

	results, err := json.Marshal(job.Result)
	if err != nil {
		r.logger.Error("failed to marshal results", "jobId", jobId, "userId", userId)
		return nil, err
	}

	_, err = r.db.ExecContext(ctx, createJobQuery, job.Id, job.Name, userId, string(task), job.Model, job.Status, job.Time, job.Freq, string(results))
	if err != nil {
		r.logger.Error("failed to create job", "jobId", jobId, "userId", userId)
		return nil, err
	}

	return job, nil
}

func (r *sqliteRepo) GetJob(ctx context.Context, jobId string, userId string) (*jobs.Job, error) {
	r.logger.Info("Fetching job", "jobId", jobId, "userId", userId)
	row := r.db.QueryRowContext(ctx, getJobQuery, jobId, userId)

	var job jobs.Job
	var task, result string

	err := row.Scan(&job.Id, &userId, &task, &job.Model, &job.Status, &job.Time, &job.Freq, &result)
	if err != nil {
		r.logger.Error("failed to scan job", "jobId", jobId, "userId", userId)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("job not found")
		}
		return nil, err
	}

	err = json.Unmarshal([]byte(task), &job.Task)
	if err != nil {
		r.logger.Error("failed to unmarshal task", "jobId", jobId, "userId", userId)
		return nil, err
	}

	err = json.Unmarshal([]byte(result), &job.Result)
	if err != nil {
		r.logger.Error("failed to unmarshal result", "jobId", jobId, "userId", userId)
		return nil, err
	}

	return &job, nil
}

func (r *sqliteRepo) GetJobs(ctx context.Context, userId string, limit, offset int) (jobs.Jobs, error) {
	r.logger.Info("Fetching jobs", "userId", userId, "limit", limit, "offset", offset)
	rows, err := r.db.QueryContext(ctx, getJobsQuery, userId, limit, offset)
	if err != nil {
		r.logger.Error("failed to query jobs", "userId", userId, "limit", limit, "offset", offset)
		return nil, err
	}
	defer rows.Close()

	var jobList jobs.Jobs
	for rows.Next() {
		var job jobs.Job
		var task, result string

		err := rows.Scan(&job.Id, &job.Name, &userId, &task, &job.Model, &job.Status, &job.Time, &job.Freq, &result)
		if err != nil {
			r.logger.Error("failed to scan job", "userId", userId)
			return nil, err
		}

		err = json.Unmarshal([]byte(task), &job.Task)
		if err != nil {
			r.logger.Error("failed to unmarshal task", "userId", userId)
			return nil, err
		}

		err = json.Unmarshal([]byte(result), &job.Result)
		if err != nil {
			r.logger.Error("failed to unmarshal result", "userId", userId)
			return nil, err
		}

		jobList = append(jobList, &job)
	}

	return jobList, nil
}

func (r *sqliteRepo) UpdateJob(ctx context.Context, job jobs.Job) (*jobs.Job, error) {
	r.logger.Info("Updating job", "jobId", job.Id)
	task, err := json.Marshal(job.Task)
	if err != nil {
		r.logger.Error("failed to marshal task", "jobId", job.Id)
		return nil, err
	}

	results, err := json.Marshal(job.Result)
	if err != nil {
		r.logger.Error("failed to marshal results", "jobId", job.Id)
		return nil, err
	}

	_, err = r.db.ExecContext(ctx, updateJobQuery, string(task), job.Model, job.Status, job.Time, job.Freq, string(results), job.Id)
	if err != nil {
		r.logger.Error("failed to update job", "jobId", job.Id)
		return nil, err
	}

	return &job, nil
}

func (r *sqliteRepo) DeleteJob(ctx context.Context, jobId string, userId string) error {
	r.logger.Info("Deleting job", "jobId", jobId, "userId", userId)
	_, err := r.db.ExecContext(ctx, deleteJobQuery, jobId, userId)
	if err != nil {
		r.logger.Error("failed to delete job", "jobId", jobId, "userId", userId)
	}
	return err
}

func (r *sqliteRepo) Peek(ctx context.Context) (*jobs.Job, error) {
	r.logger.Info("Peeking for a job")
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.logger.Error("failed to begin transaction")
		return nil, err
	}
	defer tx.Rollback() // Rollback on any error

	row := tx.QueryRowContext(ctx, peekJobQuery)

	var job jobs.Job
	var task, result, userId string

	err = row.Scan(&job.Id, &userId, &task, &job.Model, &job.Status, &job.Time, &job.Freq, &result)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // No job available
		}
		r.logger.Error("failed to scan job")
		return nil, err
	}

	_, err = tx.ExecContext(ctx, updateJobStatusQuery, "running", job.Id)
	if err != nil {
		r.logger.Error("failed to update job status", "jobId", job.Id)
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		r.logger.Error("failed to commit transaction", "jobId", job.Id)
		return nil, err
	}

	err = json.Unmarshal([]byte(task), &job.Task)
	if err != nil {
		r.logger.Error("failed to unmarshal task", "jobId", job.Id)
		return nil, err
	}

	err = json.Unmarshal([]byte(result), &job.Result)
	if err != nil {
		r.logger.Error("failed to unmarshal result", "jobId", job.Id)
		return nil, err
	}

	return &job, nil
}
