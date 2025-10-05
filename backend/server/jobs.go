package server

import (
	"encoding/json"
	"net/http"
)

func createJob(w http.ResponseWriter, r *http.Request) {
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
func viewAllJobs(w http.ResponseWriter, r *http.Request) {}
func runJob(w http.ResponseWriter, r *http.Request) {}
