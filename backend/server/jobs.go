package server

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func (s *Server) createJob(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
    }

    // create in the db

    w.WriteHeader(http.StatusNotImplemented)
    json.NewEncoder(w).Encode(map[string]string{
        "error": "Currently unavailable",
    })
}


func updateJob(w http.ResponseWriter, r *http.Request) {}
func viewJob(w http.ResponseWriter, r *http.Request) {}


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

    w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}


func runJob(w http.ResponseWriter, r *http.Request) {}
