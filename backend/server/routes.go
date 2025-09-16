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


// TODO: set up and api handler
func addRoutes(h *http.ServeMux, s *Server) {
    // basic routes
    h.HandleFunc("/", index)
    h.HandleFunc("/health", healthCheck)

    // user status routes
    h.HandleFunc("/api/login", s.loginHandler)
    h.HandleFunc("/api/logout", s.logoutHandler)

	// TODO: will have to set up as an api and requires DB
    h.HandleFunc("/api/reresh", s.refreshHandler)

    h.HandleFunc("/api/register", s.register)

    h.HandleFunc("/api/chat", s.jwt_manager.AuthApiMiddleware(s.rateLimitMiddleware(chatHandler)))

    h.HandleFunc("/api/stream", s.jwt_manager.AuthApiMiddleware(s.streamHandler))
    h.HandleFunc("/api/modelList", s.jwt_manager.AuthApiMiddleware(modelListHandler))
    h.HandleFunc("/api/modelInfo", s.jwt_manager.AuthApiMiddleware(modelShowHandler))

	h.HandleFunc("/api/chatDescription", s.jwt_manager.AuthApiMiddleware(s.chatDescriptionHandler))

	h.HandleFunc("/api/me", s.jwt_manager.AuthApiMiddleware(s.verifyHandler))

	h.HandleFunc("/api/messages", s.jwt_manager.AuthApiMiddleware(s.getChatHandler))
	h.HandleFunc("/api/deleteMessages", s.jwt_manager.AuthApiMiddleware(s.deleteChatHandler))
}

// =======================================
