package server

import (
	"encoding/json"
	"net/http"
)

type HealthResponse struct {
    Status string `json:"status"`
}

func index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Page Not Found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, "/health", http.StatusTemporaryRedirect)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
    response := HealthResponse{Status: "ok"}

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)
}


func addRoutes(h *http.ServeMux, s *Server) {
    // basic routes
    h.HandleFunc("/", index)
    h.HandleFunc("/health", healthCheck)

    // user status routes
    h.HandleFunc("/api/login", s.loginHandler)
    h.HandleFunc("/api/logout", s.logoutHandler)
    h.HandleFunc("/api/refresh", s.refreshHandler)
    h.HandleFunc("/api/register", s.registerHandler)

	// chat and model routes
    h.Handle("/api/chat", s.jwt_manager.AuthMiddleware(
		s.rateLimitMiddleware(http.HandlerFunc(chatHandler)),
	))

    h.Handle("/api/stream", s.jwt_manager.AuthMiddleware(s.streamHandler))
    h.Handle("/api/modelList", s.jwt_manager.AuthMiddleware(s.modelListHandler))
    h.Handle("/api/modelInfo", s.jwt_manager.AuthMiddleware(modelShowHandler))

	h.Handle("/api/chatDescription", s.jwt_manager.AuthMiddleware(s.chatDescriptionHandler))

	h.Handle("/api/me", s.jwt_manager.AuthMiddleware(s.verifyHandler))

	h.Handle("/api/messages", s.jwt_manager.AuthMiddleware(s.getChatHandler))
	h.Handle("/api/deleteMessages", s.jwt_manager.AuthMiddleware(s.deleteChatHandler))

	// job routes
	h.Handle("/api/jobList", s.jwt_manager.AuthMiddleware(s.viewAllJobs))
	h.Handle("/api/createJob", s.jwt_manager.AuthMiddleware(s.createJob))
}

// =======================================
