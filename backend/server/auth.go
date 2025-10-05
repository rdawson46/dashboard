package server

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/rdawson46/dashboard/db"
)


func (s *Server) loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		return
	}

	username := strings.TrimSpace(r.FormValue("username"))
	password := strings.TrimSpace(r.FormValue("password"))

	if username == "" || password == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Improper values"})
		s.logger.Error(
			"Login failed: field left blank",
			"path", r.URL.Path,
		)
		return
	}

	ctx := r.Context()
	user_db, err := s.db.SignInUser(ctx, username, password)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Couldn't login"})
		s.logger.Error(
			"Login failed: No user found",
			"error", err.Error(),
			"path", r.URL.Path,
		)
		return
	}

	s.setUserCookieResponse(*user_db, r, w)
}

func (s *Server) logoutHandler(w http.ResponseWriter, r *http.Request) {
	s.jwt_manager.ClearTokenCookie(w)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}

// TODO: modify for cookies and figure where to trigger
func (s *Server) refreshHandler(w http.ResponseWriter, r *http.Request) {
	token, err := s.jwt_manager.GetTokenFromCookie(r)

	if err != nil {
		if err == http.ErrNoCookie {
			http.Error(w, "Unauthorized: No session cookie found", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Bad request", http.StatusUnauthorized)
		return
	}

	newToken, err := s.jwt_manager.RefreshToken(token)

	if err != nil {
		http.Error(w, "Failed to refresh token: "+err.Error(), http.StatusUnauthorized)
		return
	}

	s.jwt_manager.SetTokenCookie(w, newToken)
	w.WriteHeader(http.StatusOK)
}

func (s *Server) verifyHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := userFromContext(r.Context())

	if !ok {
		http.Error(w, "Invalid user", http.StatusUnauthorized)
		return
	}

    w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(user)
} 

func (s *Server) registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid Method"})
		return
	}

	username := strings.TrimSpace(r.FormValue("username"))
	password := strings.TrimSpace(r.FormValue("password"))
	confirm := strings.TrimSpace(r.FormValue("confirm"))

	if username == "" || password == "" || confirm == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Improper values"})
		s.logger.Error(
			"User registration failed: field left blank",
			"path", r.URL.Path,
		)
		return
	}

	if password != confirm {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "password and confirmation do not match"})
		s.logger.Error(
			"User registration failed: password != confirmation",
			"path", r.URL.Path,
		)
		return
	}

	s.logger.Info("User information pulled",
		"path", r.URL.Path,
	)

	user_db, err := s.db.CreateUser(r.Context(), username, password)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create user"})
		s.logger.Error(
			"User registration failed: Failed to create in db",
			"error", err.Error(),
			"path", r.URL.Path,
		)
		return
	}

	s.setUserCookieResponse(*user_db, r, w)
}


func (s *Server)setUserCookieResponse(user_db db.User_db, r *http.Request, w http.ResponseWriter) {
	var user User_jwt
	user.Username = user_db.Username
	user.ID = user_db.ID

	token, err := s.jwt_manager.GenerateToken(&user)

	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		s.logger.Error(
			"Failed to generate user token",
			"error", err.Error(),
			"path", r.URL.Path,
		)
		return
	}

	s.jwt_manager.SetTokenCookie(w, token)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)

}
