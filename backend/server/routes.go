package server

// TODO: add logging every where
import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	ollama "github.com/ollama/ollama/api"
	api "github.com/rdawson46/dashboard/api"
)

// ========== ROUTING STRUCTS ==========
type HealthResponse struct {
    Status string `json:"status"`
}

// =====================================


// ========== ROUTING FUNCTIONS ==========
func index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Page Not Found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, "/health", http.StatusFound)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
    response := HealthResponse{Status: "ok"}

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)
}

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

// will work on to replace the chatHandler
func (s *Server) streamHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control", "*")


	flusher, ok := w.(http.Flusher)

	if !ok {
		s.logger.Error("Failed to get flusher")
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	var chatReq api.StreamRequest
	err := json.NewDecoder(r.Body).Decode(&chatReq)

	if err != nil {
		s.logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	url := os.Getenv("OLLAMA_URL")

	if url == "" {
		s.logger.Error("Ollama URL not set")
		http.Error(w, "Ollama URL not set", http.StatusInternalServerError)
		return
	}

	ctx := r.Context()

	oc, err := api.NewOllamaClient(url)

	if err != nil {
		s.logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return 
	}

	msgChan := make(chan ollama.ChatResponse)
	errChan := make(chan error)

	s.logger.Info(
		"Starting stream",
		"model", chatReq.Model,
		"remote", r.RemoteAddr,
	)

	go oc.Stream(ctx, chatReq, chatReq.Model, msgChan, errChan)

	token_count := 0
	OuterLoop:
	for {
		select {
		case resp, ok := <-msgChan:
			if !ok {
				break OuterLoop
			}

			// encode resp
			b, err := json.Marshal(resp)

			if err != nil {
				s.logger.Error(err.Error())
				http.Error(w, "failed to marshal resp", http.StatusInternalServerError)
				return
			}

			fmt.Fprintf(w, "data: %s\n\n", b)
			flusher.Flush()
			token_count = resp.EvalCount;
		case err, ok := <-errChan:
			if !ok {
				break OuterLoop
			}

			// HACK: replace error message
			s.logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	s.logger.Info(
		"Streaming finished",
		"token_count", token_count,
		"remote", r.RemoteAddr,
	)

	fmt.Fprintf(w, "data: %s\n\n", `{"done": true}`)
	flusher.Flush()
}

func modelListHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
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

    resp, err := oc.GetModelList(ctx)

    if err != nil {
        http.Error(w, "Internal Error", http.StatusInternalServerError)
        return 
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(resp.Models)
}

func modelShowHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    type request struct {
        Model string `json:"model"`
    }

    var showReq request
    err := json.NewDecoder(r.Body).Decode(&showReq)

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

    resp, err := oc.GetShow(ctx, showReq.Model)

    if err != nil {
        http.Error(w, fmt.Sprintf("Error getting model details for: %s", showReq.Model), http.StatusInternalServerError)
        return 
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(resp)
}

// TODO: have to get the user ID from content, requires DB setup
func chatDescriptionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

    type chatDesc struct {
        Description string `json:"description"`
    }

	// TEMP: placeholder until db is set up
	history := []chatDesc{
		{
			Description: "questions",
		},
		{
			Description: "answers",
		},
		{
			Description: "coding",
		},
	}

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(history)
}


func (s *Server) loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username := strings.TrimSpace(r.FormValue("username"))
	password := strings.TrimSpace(r.FormValue("password"))

	if username == "" || password == "" {
		http.Error(w, "Improper values", http.StatusUnauthorized)
		s.logger.Error(
			"User registration failed: field left blank",
			"path", r.URL.Path,
		)
		return
	}

	ctx := r.Context()
	user_db, err := s.db.SignInUser(ctx, username, password)

	if err != nil {
		http.Error(w, "Couldn't login", http.StatusUnauthorized)
		s.logger.Error(
			"User registration failed: field left blank",
			"path", r.URL.Path,
		)
		return
	}

	// TODO: make and set token in jwt
	var user User

	user.Username = user_db.Name
	user.ID = int(user_db.ID)

	token, err := s.jwt_manager.GenerateToken(&user)

	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		s.logger.Errorf(
			"Failed to generate user token",
			"path", r.URL.Path,
		)
		return
	}

	s.jwt_manager.SetTokenCookie(w, token)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
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

// =======================================

func (s *Server) register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
		return
	}

	username := strings.TrimSpace(r.FormValue("username"))
	password := strings.TrimSpace(r.FormValue("password"))
	confirm := strings.TrimSpace(r.FormValue("confirm"))

	if username == "" || password == "" || confirm == "" {
		http.Error(w, "Improper values", http.StatusUnauthorized)
		s.logger.Error(
			"User registration failed: field left blank",
			"path", r.URL.Path,
		)
		return
	}

	if password != confirm {
		http.Error(w, "password and confirmation do not match", http.StatusUnauthorized)
		s.logger.Errorf(
			"User registration failed: password != confirmation",
			"path", r.URL.Path,
		)
		return
	}

	s.logger.Info("User information pulled",
		"path", r.URL.Path,
	)

	db_user, err := s.db.CreateUser(r.Context(), username, password)

	if err != nil {
		http.Error(w, "Failed to create user", http.StatusUnauthorized)
		s.logger.Errorf(
			"User registration failed: Failed to create in db",
			err.Error(),
			"path", r.URL.Path,
		)
		return
	}

	user := &User{} // username and ID
	user.Username = db_user.Name
	user.ID = int(db_user.ID)

	token, err := s.jwt_manager.GenerateToken(user)

	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		s.logger.Errorf(
			"Failed to generate user token",
			"path", r.URL.Path,
		)
		return
	}

	s.jwt_manager.SetTokenCookie(w, token)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}


// =======================================



// ========== HANDLER FUNCTIONS ==========
// TODO: set up and api handler

func addRoutes(h *http.ServeMux, s *Server) {
    // basic routes
    h.HandleFunc("/", index)
    h.HandleFunc("/health", healthCheck)

    // user status routes
	// TODO: will have to set up as an api and requires DB
    h.HandleFunc("/api/login", s.loginHandler)
    h.HandleFunc("/api/logout", s.logoutHandler)
    h.HandleFunc("/api/reresh", s.refreshHandler)

    // TODO: add user to db
    h.HandleFunc("/api/register", s.register)

    // api auth handler
    h.HandleFunc("/api/chat", s.jwt_manager.AuthApiMiddleware(
        s.rateLimitMiddleware(chatHandler),
    ))

    h.HandleFunc("/api/stream", s.jwt_manager.AuthApiMiddleware(s.streamHandler))
    h.HandleFunc("/api/modelList", s.jwt_manager.AuthApiMiddleware(modelListHandler))
    h.HandleFunc("/api/modelInfo", s.jwt_manager.AuthApiMiddleware(modelShowHandler))

	h.HandleFunc("/api/chatDescription", chatDescriptionHandler)

	h.HandleFunc("/api/me", s.jwt_manager.AuthApiMiddleware(s.verifyHandler))
}

// =======================================
