package api

import (
	"os"
	"testing"

	"github.com/charmbracelet/log"
	ollama "github.com/ollama/ollama/api"
)

func TestOllamaClient(t *testing.T) {
	t.Run("NewOllamaClient_Valid", func(t *testing.T) {
		client, err := NewOllamaClient("http://localhost:11434", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if client == nil {
			t.Fatal("expected client, got nil")
		}
	})

	t.Run("NewOllamaClient_InvalidURL", func(t *testing.T) {
		_, err := NewOllamaClient(":", nil)
		if err == nil {
			t.Error("expected error for invalid URL, got nil")
		}
	})

	t.Run("newRequest", func(t *testing.T) {
		client, _ := NewOllamaClient("http://localhost:11434", log.Default())
		
		os.Setenv("DEFAULT_MODEL", "test-model")
		defer os.Unsetenv("DEFAULT_MODEL")

		req := client.newRequest("hello", &[]bool{false}[0])
		if req == nil {
			t.Fatal("expected request, got nil")
		}
		if req.Model != "test-model" {
			t.Errorf("expected model test-model, got %s", req.Model)
		}
		if req.Messages[1].Content != "hello" {
			t.Errorf("expected message content hello, got %s", req.Messages[1].Content)
		}
	})

	t.Run("newRequestWithMessages", func(t *testing.T) {
		client, _ := NewOllamaClient("http://localhost:11434", log.Default())
		messages := []ollama.Message{{Role: "user", Content: "hi"}}
		req := client.newRequestWithMessages(messages, "model-x", true)
		
		if req.Model != "model-x" {
			t.Errorf("expected model model-x, got %s", req.Model)
		}
		if len(req.Messages) != 2 {
			t.Errorf("expected 2 messages (system + user), got %d", len(req.Messages))
		}
	})

	t.Run("Semaphore", func(t *testing.T) {
		InitOllamaSemaphore(1)
		// Should not panic or error
		if ollamaSemaphore == nil {
			t.Error("expected semaphore to be initialized")
		}
	})
}
