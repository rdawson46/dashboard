package server

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/charmbracelet/log"
	ollama "github.com/ollama/ollama/api"
	"github.com/rdawson46/dashboard/api"
)

func TestSetSSE(t *testing.T) {
	w := httptest.NewRecorder()
	flusher, ok := SetSSE(w)

	if !ok {
		t.Error("SetSSE failed to return a flusher")
	}
	if flusher == nil {
		t.Error("SetSSE returned a nil flusher")
	}

	resp := w.Result()
	if resp.Header.Get("Content-Type") != "text/event-stream" {
		t.Errorf("expected Content-Type text/event-stream, got %s", resp.Header.Get("Content-Type"))
	}
}

func TestGetMessageID(t *testing.T) {
	s := &Server{
		db:     &mockRepository{},
		logger: log.Default(),
	}

	user := &User_jwt{ID: "user1", Username: "testuser"}
	ctx := context.Background()

	t.Run("NewMessage", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil).WithContext(ctx)
		chatReq := api.StreamRequest{MessageId: "", Model: "test-model"}

		msgCreation, err := GetMessageID(s, req, chatReq, user)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if msgCreation.Id != "new-message-id" {
			t.Errorf("expected new-message-id, got %s", msgCreation.Id)
		}
		if !msgCreation.New {
			t.Error("expected New to be true")
		}
	})

	t.Run("ExistingMessage", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil).WithContext(ctx)
		chatReq := api.StreamRequest{MessageId: "existing-id", Model: "test-model"}

		msgCreation, err := GetMessageID(s, req, chatReq, user)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if msgCreation.Id != "existing-id" {
			t.Errorf("expected existing-id, got %s", msgCreation.Id)
		}
		if msgCreation.New {
			t.Error("expected New to be false")
		}
	})
}

func TestSendSSEMessage(t *testing.T) {
	s := &Server{logger: log.Default()}
	w := httptest.NewRecorder()
	flusher := w

	data := map[string]string{"foo": "bar"}
	err := SendSSEMessage(s, flusher, w, data, "test-type")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	body := w.Body.String()
	expectedPrefix := "data: "
	if body[:len(expectedPrefix)] != expectedPrefix {
		t.Errorf("expected prefix %s, got %s", expectedPrefix, body)
	}

	var envelope map[string]any
	jsonErr := json.Unmarshal([]byte(body[len(expectedPrefix):]), &envelope)
	if jsonErr != nil {
		t.Fatalf("failed to unmarshal body: %v", jsonErr)
	}

	if envelope["type"] != "test-type" {
		t.Errorf("expected type test-type, got %s", envelope["type"])
	}
}

func TestSendSimpleSEE(t *testing.T) {
	s := &Server{logger: log.Default()}
	w := httptest.NewRecorder()
	flusher := w

	message := "simple message"
	SendSimpleSEE(s, flusher, w, message)

	body := w.Body.String()
	expected := "data: simple message\n\n"
	if body != expected {
		t.Errorf("expected %q, got %q", expected, body)
	}
}

func TestSaveMessage(t *testing.T) {
	s := &Server{
		db:     &mockRepository{},
		logger: log.Default(),
	}
	user := &User_jwt{ID: "user1", Username: "testuser"}
	req := httptest.NewRequest("POST", "/", nil)
	chatReq := &api.StreamRequest{Model: "test-model", Query: []ollama.Message{{Role: "user", Content: "hello"}}}

	// Just verifying it doesn't panic and handles success
	SaveMessage(s, req, "msg1", "hi back", chatReq, user)

	if len(chatReq.Query) != 2 {
		t.Errorf("expected 2 messages in query, got %d", len(chatReq.Query))
	}
	if chatReq.Query[1].Content != "hi back" {
		t.Errorf("expected last message content 'hi back', got %s", chatReq.Query[1].Content)
	}
}
