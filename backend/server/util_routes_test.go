package server

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/charmbracelet/log"
	ollama "github.com/ollama/ollama/api"
)

func TestUtilRoutes(t *testing.T) {
	mockRepo := &mockRepository{}
	s := &Server{
		db:     mockRepo,
		logger: log.Default(),
	}
	user := &User_jwt{ID: "user123", Username: "testuser"}

	t.Run("DeleteChat_Success", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{
			"chatId": "chat123",
			"userId": "user123",
		})
		req := httptest.NewRequest(http.MethodDelete, "/delete-chat", bytes.NewReader(body))
		req = req.WithContext(contextWithUser(req.Context(), user))
		w := httptest.NewRecorder()

		mockRepo.DeleteMessageFunc = func(ctx context.Context, id, userId string) (bool, error) {
			return true, nil
		}

		s.deleteChatHandler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status OK, got %d. Body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("GetChat_Success", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{
			"chatId": "chat123",
			"userId": "user123",
		})
		req := httptest.NewRequest(http.MethodPost, "/get-chat", bytes.NewReader(body))
		req = req.WithContext(contextWithUser(req.Context(), user))
		w := httptest.NewRecorder()

		mockRepo.GetMessageFunc = func(ctx context.Context, userId, id string) ([]ollama.Message, error) {
			return []ollama.Message{{Role: "user", Content: "hi"}}, nil
		}

		s.getChatHandler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status OK, got %d. Body: %s", w.Code, w.Body.String())
		}
		var resp []ollama.Message
		json.NewDecoder(w.Body).Decode(&resp)
		if len(resp) != 1 || resp[0].Content != "hi" {
			t.Errorf("unexpected response: %v", resp)
		}
	})
}
