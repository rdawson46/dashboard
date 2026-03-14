package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/charmbracelet/log"
	"github.com/rdawson46/dashboard/jobs"
)

func TestJobHandlers(t *testing.T) {
	mockRepo := &mockRepository{}
	s := &Server{
		db:     mockRepo,
		logger: log.Default(),
	}
	user := &User_jwt{ID: "user123", Username: "testuser"}

	t.Run("CreateJob_Success", func(t *testing.T) {
		form := url.Values{}
		form.Add("name", "Test Job")
		form.Add("type", "LLM")
		form.Add("freq", "daily")
		form.Add("query", "hello")

		req := httptest.NewRequest(http.MethodPost, "/create-job", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req = req.WithContext(contextWithUser(req.Context(), user))
		w := httptest.NewRecorder()

		mockRepo.CreateJobFunc = func(ctx context.Context, userId string, job *jobs.Job) (*jobs.Job, error) {
			job.Id = "job-123"
			return job, nil
		}

		s.createJob(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status OK, got %d. Body: %s", w.Code, w.Body.String())
		}
		var resp jobs.Job
		json.NewDecoder(w.Body).Decode(&resp)
		if resp.Id != "job-123" {
			t.Errorf("expected job Id job-123, got %s", resp.Id)
		}
	})

	t.Run("DeleteJob_Success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/delete-job?jobId=job-123", nil)
		req = req.WithContext(contextWithUser(req.Context(), user))
		w := httptest.NewRecorder()

		mockRepo.DeleteJobFunc = func(ctx context.Context, jobId string, userId string) error {
			return nil
		}

		s.deleteJob(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status OK, got %d", w.Code)
		}
	})

	t.Run("ViewJob_Success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/view-job?jobId=job-123", nil)
		req = req.WithContext(contextWithUser(req.Context(), user))
		w := httptest.NewRecorder()

		mockRepo.GetJobFunc = func(ctx context.Context, jobId string, userId string) (*jobs.Job, error) {
			return &jobs.Job{Id: "job-123", Name: "Test Job"}, nil
		}

		s.viewJob(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status OK, got %d", w.Code)
		}
		var resp map[string]any
		json.NewDecoder(w.Body).Decode(&resp)
		if resp["status"] != "ok" {
			t.Errorf("expected status ok, got %v", resp["status"])
		}
	})

	t.Run("ViewAllJobs_Success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/view-all-jobs", nil)
		req = req.WithContext(contextWithUser(req.Context(), user))
		w := httptest.NewRecorder()

		mockRepo.GetJobsFunc = func(ctx context.Context, userId string, limit, offset int) (jobs.Jobs, error) {
			return jobs.Jobs{{Id: "1", Name: "j1"}}, nil
		}

		s.viewAllJobs(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status OK, got %d", w.Code)
		}
	})
}
