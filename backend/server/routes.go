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

	addAuthRoute := func(path string, handler http.HandlerFunc) {
		h.Handle(path, s.jwt_manager.AuthMiddleware(handler))
	}

	addAuthAndLimiting := func(path string, handler http.HandlerFunc) {
		h.Handle(path, s.rateLimitMiddleware(
			s.jwt_manager.AuthMiddleware(handler),
		))
	}

	// chat and model routes
	addAuthAndLimiting("/api/chat", http.HandlerFunc(chatHandler))
	addAuthAndLimiting("/api/stream", s.streamHandler)
	addAuthRoute("/api/modelList", s.modelListHandler)
	addAuthRoute("/api/modelInfo", modelShowHandler)
	addAuthRoute("/api/chatDescription", s.chatDescriptionHandler)
	addAuthRoute("/api/me", s.verifyHandler)
	addAuthRoute("/api/messages", s.getChatHandler)
	addAuthRoute("/api/deleteMessages", s.deleteChatHandler)

	// job routes
	addAuthRoute("/api/jobList", s.viewAllJobs)
	addAuthRoute("/api/createJob", s.createJob)
	addAuthRoute("/api/deleteJob", s.deleteJob)
	addAuthRoute("/api/updateJob", s.updateJob)
	addAuthRoute("/api/viewJob", s.viewJob)

	// file routes
	addAuthRoute("/api/uploadFile", s.uploadFile)
	addAuthRoute("/api/getFileList", s.getFileList)
	addAuthRoute("/api/getFile", s.getFile)
	addAuthRoute("/api/deleteFile", s.deleteFile)
}

// =======================================
