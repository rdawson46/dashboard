package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/log"
	"github.com/rdawson46/dashboard/db"
)

func TestAuthHandlers(t *testing.T) {
	mockRepo := &mockRepository{}
	jwtManager := NewJWTManager("secret", time.Hour)
	s := &Server{
		db:          mockRepo,
		jwt_manager: jwtManager,
		logger:      log.Default(),
	}

	t.Run("LoginHandler_Success", func(t *testing.T) {
		form := url.Values{}
		form.Add("username", "testuser")
		form.Add("password", "password")

		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()

		mockRepo.SignInUserFunc = func(ctx context.Context, username, password string) (*db.User_db, error) {
			return &db.User_db{ID: "123", Username: "testuser"}, nil
		}

		s.loginHandler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status OK, got %d", w.Code)
		}

		var resp User_jwt
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if resp.Username != "testuser" {
			t.Errorf("expected username testuser, got %s", resp.Username)
		}

		// Check for cookie
		cookies := w.Result().Cookies()
		found := false
		for _, c := range cookies {
			if c.Name == "auth_token" {
				found = true
				break
			}
		}
		if !found {
			t.Error("expected auth_token cookie not found")
		}
	})

	t.Run("LoginHandler_Failure", func(t *testing.T) {
		form := url.Values{}
		form.Add("username", "baduser")
		form.Add("password", "badpass")

		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()

		mockRepo.SignInUserFunc = func(ctx context.Context, username, password string) (*db.User_db, error) {
			return nil, errors.New("auth failed")
		}

		s.loginHandler(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected status Unauthorized, got %d", w.Code)
		}
	})

	t.Run("RegisterHandler_Success", func(t *testing.T) {
		form := url.Values{}
		form.Add("username", "newuser")
		form.Add("password", "password")
		form.Add("confirm", "password")

		req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()

		mockRepo.CreateUserFunc = func(ctx context.Context, username, password string) (*db.User_db, error) {
			return &db.User_db{ID: "456", Username: "newuser"}, nil
		}

		s.registerHandler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status OK, got %d", w.Code)
		}
	})

	t.Run("RegisterHandler_PasswordMismatch", func(t *testing.T) {
		form := url.Values{}
		form.Add("username", "newuser")
		form.Add("password", "password")
		form.Add("confirm", "mismatch")

		req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()

		s.registerHandler(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected status Unauthorized, got %d", w.Code)
		}
	})
}
