package server

import (
	"net/http/httptest"
	"testing"
	"time"
)

func TestJWTManager(t *testing.T) {
	secretKey := "secret"
	tokenDuration := time.Hour
	manager := NewJWTManager(secretKey, tokenDuration)

	user := &User_jwt{
		Username: "testuser",
		ID:       "123",
	}

	t.Run("GenerateAndValidateToken", func(t *testing.T) {
		token, err := manager.GenerateToken(user)
		if err != nil {
			t.Fatalf("failed to generate token: %v", err)
		}

		claims, err := manager.ValidateToken(token)
		if err != nil {
			t.Fatalf("failed to validate token: %v", err)
		}

		if claims.Username != user.Username {
			t.Errorf("expected username %s, got %s", user.Username, claims.Username)
		}
		if claims.ID != user.ID {
			t.Errorf("expected ID %s, got %s", user.ID, claims.ID)
		}
	})

	t.Run("ValidateInvalidToken", func(t *testing.T) {
		_, err := manager.ValidateToken("invalid-token")
		if err == nil {
			t.Error("expected error for invalid token, got nil")
		}
	})

	t.Run("RefreshToken", func(t *testing.T) {
		token, err := manager.GenerateToken(user)
		if err != nil {
			t.Fatalf("failed to generate token: %v", err)
		}

		// Wait for more than 1 second to ensure the new token has a different issued at time
		// NumericDate truncates to seconds.
		time.Sleep(time.Second + time.Millisecond*100)

		newToken, err := manager.RefreshToken(token)
		if err != nil {
			t.Fatalf("failed to refresh token: %v", err)
		}

		if newToken == token {
			t.Error("expected new token to be different from old token")
		}

		claims, err := manager.ValidateToken(newToken)
		if err != nil {
			t.Fatalf("failed to validate new token: %v", err)
		}

		if claims.Username != user.Username {
			t.Errorf("expected username %s, got %s", user.Username, claims.Username)
		}
	})

	t.Run("CookieOperations", func(t *testing.T) {
		token := "test-token"
		w := httptest.NewRecorder()

		manager.SetTokenCookie(w, token)

		resp := w.Result()
		cookies := resp.Cookies()
		found := false
		for _, c := range cookies {
			if c.Name == "auth_token" {
				found = true
				if c.Value != token {
					t.Errorf("expected cookie value %s, got %s", token, c.Value)
				}
				break
			}
		}
		if !found {
			t.Error("auth_token cookie not found in response")
		}

		// Test GetTokenFromCookie
		req := httptest.NewRequest("GET", "/", nil)
		for _, c := range cookies {
			req.AddCookie(c)
		}

		gotToken, err := manager.GetTokenFromCookie(req)
		if err != nil {
			t.Fatalf("failed to get token from cookie: %v", err)
		}
		if gotToken != token {
			t.Errorf("expected token %s, got %s", token, gotToken)
		}

		// Test ClearTokenCookie
		w2 := httptest.NewRecorder()
		manager.ClearTokenCookie(w2)
		resp2 := w2.Result()
		cookies2 := resp2.Cookies()
		found2 := false
		for _, c := range cookies2 {
			if c.Name == "auth_token" {
				found2 = true
				if c.Value != "" {
					t.Errorf("expected empty cookie value, got %s", c.Value)
				}
				if !c.Expires.Before(time.Now().Add(time.Second)) {
					t.Errorf("expected cookie to be expired, got expires %v", c.Expires)
				}
				break
			}
		}
		if !found2 {
			t.Error("auth_token cookie not found in clear response")
		}
	})
}
