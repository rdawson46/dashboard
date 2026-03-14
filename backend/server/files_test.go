package server

import (
	"bytes"
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/charmbracelet/log"
	"github.com/rdawson46/dashboard/db"
)

func TestFileHandlers(t *testing.T) {
	mockRepo := &mockRepository{}
	s := &Server{
		db:     mockRepo,
		logger: log.Default(),
	}
	user := &User_jwt{ID: "user123", Username: "testuser"}

	t.Run("UploadFile_Success", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("file", "test.txt")
		part.Write([]byte("file content"))
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req = req.WithContext(contextWithUser(req.Context(), user))
		w := httptest.NewRecorder()

		mockRepo.SaveFileFunc = func(ctx context.Context, userId, fileName, contentType, content string) (string, error) {
			return "new-file-id", nil
		}

		s.uploadFile(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status OK, got %d", w.Code)
		}
		var resp map[string]string
		json.NewDecoder(w.Body).Decode(&resp)
		if resp["fileID"] != "new-file-id" {
			t.Errorf("expected fileID new-file-id, got %s", resp["fileID"])
		}
	})

	t.Run("GetFile_Success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/file?fileId=123", nil)
		req = req.WithContext(contextWithUser(req.Context(), user))
		w := httptest.NewRecorder()

		mockRepo.GetFileFunc = func(ctx context.Context, fileId, userId string) (*db.File, error) {
			return &db.File{ID: "123", FileName: "test.txt", Content: "hello", ContentType: "text/plain"}, nil
		}

		s.getFile(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status OK, got %d", w.Code)
		}
		if w.Body.String() != "hello" {
			t.Errorf("expected body 'hello', got %s", w.Body.String())
		}
	})

	t.Run("GetFileList_Success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/files", nil)
		req = req.WithContext(contextWithUser(req.Context(), user))
		w := httptest.NewRecorder()

		mockRepo.GetFilesFunc = func(ctx context.Context, userId string, limit, offset int) ([]*db.File, error) {
			return []*db.File{{ID: "1", FileName: "f1.txt"}}, nil
		}
		mockRepo.GetFilesCountFunc = func(ctx context.Context, userId string) (int, error) {
			return 1, nil
		}

		s.getFileList(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status OK, got %d", w.Code)
		}
		var resp map[string]any
		json.NewDecoder(w.Body).Decode(&resp)
		if resp["totalItems"].(float64) != 1 {
			t.Errorf("expected totalItems 1, got %v", resp["totalItems"])
		}
	})

	t.Run("DeleteFile_Success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/delete-file?fileId=123", nil)
		req = req.WithContext(contextWithUser(req.Context(), user))
		w := httptest.NewRecorder()

		mockRepo.DeleteFileFunc = func(ctx context.Context, fileId, userId string) error {
			return nil
		}

		s.deleteFile(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status OK, got %d", w.Code)
		}
	})
}
