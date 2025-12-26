package server

import (
	"encoding/json"
	"net/http"
	"os"

	ollama "github.com/ollama/ollama/api"
	api "github.com/rdawson46/dashboard/api"
	"github.com/rdawson46/dashboard/db"
)

func chatHandler(w http.ResponseWriter, r *http.Request) {
    type request struct {
        Query string `json:"query"`
    }

    var chatReq request
    err := json.NewDecoder(r.Body).Decode(&chatReq)

    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

	url := os.Getenv("OLLAMA_URL")
	
	if url == "" {
        http.Error(w, "Ollama URL not set", http.StatusInternalServerError)
        return
	}

    ctx := r.Context()

    oc, err := api.NewOllamaClient(url)

    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return 
    }

    res, err := oc.Chat(ctx, chatReq.Query)

    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return 
    }

    response := map[string]string{"respones": res}
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)
}

/*
TODO:
Frist Stage:
	* attached files list
	* rag flag
Second Stage:
	* send a `notification` for a new description
*/
func (s *Server) streamHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user, ok := userFromContext(r.Context())

	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"Error": "User not found"})
		return
	}

	var chatReq api.StreamRequest
	err := json.NewDecoder(r.Body).Decode(&chatReq)

	if err != nil {
		s.logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if user.Username != chatReq.Username {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"Error": "Invalid username"})
		return
	}

	// NOTE: this might make it harder to update history in UI in realtime
	messageId, err := GetMessageID(s, r, chatReq, user)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"Error": "Unable to Create Chat"})
		return
	}

	file_count := len(chatReq.FileIds)

	files := make([]*db.File, 0)
	if file_count > 0 {
		s.logger.Info("Fetching file content", "count", file_count)

		for _, f := range chatReq.FileIds {
			res, err := s.db.GetFile(r.Context(), f, user.ID)

			if err != nil {
				s.logger.Error("Error retrieving file", "FileId", f, "Error", err.Error())
				continue
			}

			files = append(files, res)
		}
	}

	flusher, ok := SetSSE(w)

	if !ok {
		s.logger.Error("Failed to get flusher")
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	s.logger.Info(
		"Setting chat Id",
		"User", user.Username,
		"Chat Id", messageId,
	)

	err = SendSSEMessage(s, flusher, w, messageId, "Message ID")
	if err != nil { return }

	s.logger.Info(
		"Chat Id set",
		"User", user.Username,
		"Chat Id", messageId,
	)

	originalMessages := make([]ollama.Message, len(chatReq.Query))
	copy(originalMessages, chatReq.Query)

	if len(files) > 0 {
		tempMessages, ok := addFileToMessages(r.Context(), chatReq.Model, chatReq.Query, files)

		if ok {
			chatReq.Query = tempMessages
		}
	}

	chatResponse, err := Streamer(s, w, flusher, r, chatReq, messageId, user)

	if len(originalMessages) > 0 {
		chatReq.Query = originalMessages
	}
	SaveMessage(s, r, messageId, chatResponse, &chatReq, user)
}

func (s *Server) chatDescriptionHandler(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "Invalid user", http.StatusUnauthorized)
		return
	}

	s.logger.Info(
		"Pulling user history",
		"userId", user.ID,
		"username", user.Username,
	)

	history, err := s.db.GetDescriptions(ctx, user.ID, 10, 0)

	if err != nil {
		s.logger.Error(
			"Failed to get descriptions",
			"error", err.Error(),
			"path", r.URL.Path,
		)

		http.Error(w, "Unable to get chat history", http.StatusInternalServerError)
		return
	}

	s.logger.Info(
		"History pulled",
		"length", len(history),
		"path", r.URL.Path,
	)

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(history)
}

