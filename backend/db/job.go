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
	createJobQuery = `INSERT INTO jobs (id, user_id, tasks, model, status, time, freq, result) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	getJobQuery = `SELECT id, user_id, tasks, model, status, time, freq, result FROM jobs WHERE id = ? AND user_id = ?`
	getJobsQuery = `SELECT id, user_id, tasks, model, status, time, freq, result FROM jobs WHERE user_id = ? LIMIT ? OFFSET ?`
	updateJobQuery = `UPDATE jobs SET tasks = ?, model = ?, status = ?, time = ?, freq = ?, result = ? WHERE id = ?`
	deleteJobQuery = `DELETE FROM jobs WHERE id = ? AND user_id = ?`
	peekJobQuery = `SELECT id, user_id, tasks, model, status, time, freq, result FROM jobs WHERE status = 'pending' AND datetime(time) <= CURRENT_TIMESTAMP ORDER BY created_at ASC LIMIT 1`
	updateJobStatusQuery = `UPDATE jobs SET status = ? WHERE id = ?`
)

func (r *sqliteRepo) CreateJob(ctx context.Context, userId string, job jobs.Job) (*jobs.Job, error) {
	jobId := uuid.New().String()
	job.Id = jobId

	if job.Status == "" {
		job.Status = "pending"
	}

	tasks, err := json.Marshal(job.Tasks)
	if err != nil {
		return nil, err
	}

	results, err := json.Marshal(job.Result)
	if err != nil {
		return nil, err
	}

	_, err = r.db.ExecContext(ctx, createJobQuery, job.Id, userId, string(tasks), job.Model, job.Status, job.Time, job.Freq, string(results))
	if err != nil {
		return nil, err
	}

	return &job, nil
}

func (r *sqliteRepo) GetJob(ctx context.Context, jobId string, userId string) (*jobs.Job, error) {
	row := r.db.QueryRowContext(ctx, getJobQuery, jobId, userId)

	var job jobs.Job
	var tasks, result string

	err := row.Scan(&job.Id, &userId, &tasks, &job.Model, &job.Status, &job.Time, &job.Freq, &result)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("job not found")
		}
		return nil, err
	}

	err = json.Unmarshal([]byte(tasks), &job.Tasks)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(result), &job.Result)
	if err != nil {
		return nil, err
	}

	return &job, nil
}

func (r *sqliteRepo) GetJobs(ctx context.Context, userId string, limit, offset int) ([]*jobs.Job, error) {
	rows, err := r.db.QueryContext(ctx, getJobsQuery, userId, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobList []*jobs.Job
	for rows.Next() {
		var job jobs.Job
		var tasks, result string

		err := rows.Scan(&job.Id, &userId, &tasks, &job.Model, &job.Status, &job.Time, &job.Freq, &result)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal([]byte(tasks), &job.Tasks)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal([]byte(result), &job.Result)
		if err != nil {
			return nil, err
		}

		jobList = append(jobList, &job)
	}

	return jobList, nil
}

func (r *sqliteRepo) UpdateJob(ctx context.Context, job jobs.Job) (*jobs.Job, error) {
	tasks, err := json.Marshal(job.Tasks)
	if err != nil {
		return nil, err
	}

	results, err := json.Marshal(job.Result)
	if err != nil {
		return nil, err
	}

	_, err = r.db.ExecContext(ctx, updateJobQuery, string(tasks), job.Model, job.Status, job.Time, job.Freq, string(results), job.Id)
	if err != nil {
		return nil, err
	}

	return &job, nil
}

func (r *sqliteRepo) DeleteJob(ctx context.Context, jobId string, userId string) error {
	_, err := r.db.ExecContext(ctx, deleteJobQuery, jobId, userId)
	return err
}

func (r *sqliteRepo) Peek(ctx context.Context) (*jobs.Job, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback() // Rollback on any error

	row := tx.QueryRowContext(ctx, peekJobQuery)

	var job jobs.Job
	var tasks, result, userId string

	err = row.Scan(&job.Id, &userId, &tasks, &job.Model, &job.Status, &job.Time, &job.Freq, &result)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // No job available
		}
		return nil, err
	}

	_, err = tx.ExecContext(ctx, updateJobStatusQuery, "running", job.Id)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(tasks), &job.Tasks)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(result), &job.Result)
	if err != nil {
		return nil, err
	}

	return &job, nil
}
