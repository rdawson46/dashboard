package server

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/rdawson46/dashboard/jobs"
)


func (s *Server) createJob(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
		s.logger.Error("Invalid Method", "Method", r.Method)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		return
    }

	user, ok := userFromContext(r.Context())

	if !ok {
		s.logger.Error("User Not Set")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unable to authorize"})
		return
	}

	name := strings.TrimSpace(r.FormValue("name"))
	toolType := strings.TrimSpace(r.FormValue("type"))
	freq := strings.TrimSpace(r.FormValue("freq"))

	task := jobs.NewTask(toolType)

	if task == nil {
		s.logger.Error(
			"Unable to make task",
			"name", name,
			"toolType", toolType,
			"freq", freq,
		)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid task value"})
		return
	}

	// TODO: finish making the job with form values, add checks for valid name and freq
	j := jobs.NewJob(task, name, freq)

	err := j.FillIn(r.Form)

	if err != nil {
		s.logger.Error(
			"Unable to fill in job",
			"error", err.Error(),
		)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid task value"})
		return
	}

    // create in the db
	j, err = s.db.CreateJob(context.Background(), user.ID, j)

	if err != nil {
		s.logger.Error(
			"Unable to create job in DB",
			"error", err.Error(),
		)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error creating job"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(j)
}

func (s *Server) deleteJob(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
		s.logger.Error("Invalid Method", "Method", r.Method)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		return
    }

	ctx := r.Context()
	user, ok := userFromContext(ctx)

	if !ok {
		s.logger.Error(
			"User not in context",
			"path", r.URL.Path,
		)
		w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"Error": "User not found"})
        return
	}

	userId := user.ID
	jobId := r.URL.Query().Get("jobId")

	if jobId == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
        json.NewEncoder(w).Encode(map[string]string{"Error": "Missing Job Id"})
		return
	}

	background_ctx := context.Background()
	err := s.db.DeleteJob(background_ctx, jobId, userId)

	if err != nil {
		s.logger.Error(
			"Error when deleting job",
			"userId", userId,
			"jobId", jobId,
		)
		http.Error(w, "Unable to delete job", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (s *Server) updateJob(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
		s.logger.Error("Invalid Method", "Method", r.Method)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		return
    }

	user, ok := userFromContext(r.Context())

	if !ok {
		s.logger.Error("User Not Set")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unable to authorize"})
		return
	}

	name := strings.TrimSpace(r.FormValue("name"))
	toolType := strings.TrimSpace(r.FormValue("type"))
	freq := strings.TrimSpace(r.FormValue("freq"))
	jobId := strings.TrimSpace(r.FormValue("jobId"))

	task := jobs.NewTask(toolType)

	if task == nil {
		s.logger.Error(
			"Unable to make task",
			"name", name,
			"toolType", toolType,
			"freq", freq,
		)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid task value"})
		return
	}

	// TODO: finish making the job with form values, add checks for valid name and freq
	j := jobs.NewJob(task, name, freq)

	err := j.FillIn(r.Form)

	if err != nil {
		s.logger.Error(
			"Unable to fill in job",
			"error", err.Error(),
		)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid task value"})
		return
	}

	j.Id = jobId

	background_ctx := context.Background()
	updatedJob, err := s.db.UpdateJob(background_ctx, *j, user.ID)

	if err != nil {
		s.logger.Error(
			"Unable to create job in DB",
			"error", err.Error(),
		)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error creating job"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(updatedJob)
}


func (s *Server) viewJob(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.logger.Error("Invalid Method", "Method", r.Method)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		return
	}


	ctx := r.Context()
	user, ok := userFromContext(ctx)

	if !ok {
		s.logger.Error(
			"User not in context",
			"path", r.URL.Path,
		)
		w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"Error": "User not found"})
        return
	}

	userId := user.ID

	jobId := r.URL.Query().Get("jobId")

	if jobId == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
        json.NewEncoder(w).Encode(map[string]string{"Error": "Missing Job Id"})
		return
	}

	j, err := s.db.GetJob(ctx, jobId, userId)

	if err != nil {
		s.logger.Error(
			"Error when viewing job",
			"userId", userId,
			"jobId", jobId,
		)
		http.Error(w, "Unable to view job", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"status": "ok",
		"job": j,
	})
}


func (s *Server) viewAllJobs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	user, ok := userFromContext(ctx)

	if !ok {
		s.logger.Error(
			"User not in context",
			"path", r.URL.Path,
		)
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"Error": "User not found"})
        return
	}

	userId := user.ID

	queryParams := r.URL.Query()

	limitStr := queryParams.Get("limit")
	limit := 10
	if limitStr != "" {
		if val, err := strconv.Atoi(limitStr); err == nil {
			limit = val
		}
	}

	offsetStr := queryParams.Get("offset")
	offset := 0
	if offsetStr != "" {
		if val, err := strconv.Atoi(offsetStr); err == nil {
			offset = val
		}
	}

	s.logger.Info(
		"Pulling jobs form DB",
		"userId", userId,
		"limit", limit,
		"offset", offset,
	)
	jobs, err := s.db.GetJobs(ctx, userId, limit, offset)

	if err != nil {
		s.logger.Error(
			"Failed to get job list",
			"error", err.Error(),
			"path", r.URL.Path,
		)
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"Error": "User not found"})
        return
	}

	jobLen := len(jobs)

	s.logger.Info(
		"Jobs pulled",
		"userId", userId,
		"length", jobLen,
	)

	result := map[string]any{
		"jobs": jobs,
		"length": jobLen,
	}

	if offset == 0 {
		s.logger.Info(
			"Pulling job count",
			"userId", userId,
		)

		count, err := s.db.GetJobCount(ctx, userId)

		if err != nil {
			count = 0
		}

		s.logger.Info(
			"Job count pulled",
			"userId", userId,
			"count", count,
		)

		result["totalItems"] = count
	} 

    w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}


func runJob(w http.ResponseWriter, r *http.Request) {}
