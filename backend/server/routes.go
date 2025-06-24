package server

import (
	"encoding/json"
	"net/http"
)

// ========== ROUTING STRUCTS ==========
type HealthResponse struct {
    Status string `json:"status"`
}

type TestResponse struct {
    Message string `json:"status"`
}

// =====================================


// ========== ROUTING FUNCTIONS ==========
func index(w http.ResponseWriter, r *http.Request) {
    healthCheck(w, r)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
    response := HealthResponse{Status: "ok"}

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)
}

func chat(w http.ResponseWriter, r *http.Request) {
    response := TestResponse{Message: "chat"}
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)
}

func search(w http.ResponseWriter, r *http.Request) {
    response := TestResponse{Message: "search"}
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)
}


// TODO: will require jwt set up and the db
// set up basic UI and chats
func register(w http.ResponseWriter, r *http.Request) {}
func login(w http.ResponseWriter, r *http.Request) {}
func logout(w http.ResponseWriter, r *http.Request) {}
func refresh(w http.ResponseWriter, r *http.Request) {}

// =======================================


// ========== HANDLER FUNCTIONS ==========

func addRoutes(h *http.ServeMux, s *Server) {
    h.HandleFunc("/", index)
    h.HandleFunc("/health", healthCheck)
    h.HandleFunc("/search", s.rateLimitMiddleware(search))
    h.HandleFunc("/chat", s.rateLimitMiddleware(chat))
}

// =======================================
