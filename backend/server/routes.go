package server

// TODO: add logging every where
import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	ollama "github.com/ollama/ollama/api"
	"github.com/charmbracelet/log"
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

    healthCheck(w, r)
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
func streamHandler(l *log.Logger) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
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
			l.Error("Failed to get flusher")
            http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
            return
        }

        var chatReq api.StreamRequest
        err := json.NewDecoder(r.Body).Decode(&chatReq)

        if err != nil {
			l.Error(err.Error())
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

		url := os.Getenv("OLLAMA_URL")

		if url == "" {
			l.Error("Ollama URL not set")
			http.Error(w, "Ollama URL not set", http.StatusInternalServerError)
			return
		}

        ctx := r.Context()

        oc, err := api.NewOllamaClient(url)

        if err != nil {
			l.Error(err.Error())
            http.Error(w, err.Error(), http.StatusBadRequest)
            return 
        }

        msgChan := make(chan ollama.ChatResponse)
        errChan := make(chan error)

        l.Info(
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
					l.Error(err.Error())
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
				l.Error(err.Error())
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
        }

        l.Info(
            "Streaming finished",
            "token_count", token_count,
            "remote", r.RemoteAddr,
        )

        fmt.Fprintf(w, "data: %s\n\n", `{"done": true}`)
        flusher.Flush()
    }
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

// TODO: get model name
func modelShowHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // TODO: get model name from the request body
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

// TODO: have to get the user ID from content
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


func loginHandler(jwtManager *JWTManager) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }

        var loginReq struct {
            Username string `json:"username"`
            Password string `json:"password"`
        }

        if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
            http.Error(w, "Invalid JSON", http.StatusBadRequest)
            return
        }

        // TODO: add db logic
        if loginReq.Username == "admin" && loginReq.Password == "password" {
            user := &User{
                Username: "admin", 
                ID: 1,
            }

            token, err := jwtManager.GenerateToken(user)
            if err != nil {
                http.Error(w, "Failed to generate token", http.StatusInternalServerError)
                return
            }

            jwtManager.SetTokenCookie(w, token)
            w.WriteHeader(http.StatusOK)
        } else {
            http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        }
    }
}

func logoutHandler(jwtManager *JWTManager) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        jwtManager.ClearTokenCookie(w)
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
    }
}

// TODO: modify for cookies and figure where to trigger
func refreshHandler(jwtManager *JWTManager) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        token, err := jwtManager.GetTokenFromCookie(r)

        if err != nil {
            if err == http.ErrNoCookie {
                http.Error(w, "Unauthorized: No session cookie found", http.StatusUnauthorized)
                return
            }
            http.Error(w, "Bad request", http.StatusUnauthorized)
            return
        }

        newToken, err := jwtManager.RefreshToken(token)

        if err != nil {
            http.Error(w, "Failed to refresh token: "+err.Error(), http.StatusUnauthorized)
            return
        }

        jwtManager.SetTokenCookie(w, newToken)
        w.WriteHeader(http.StatusOK)
    }
}

/*
func protectedHandler(w http.ResponseWriter, r *http.Request) {
    user, ok := userFromContext(r.Context())

    if !ok {
        http.Error(w, "User not found in context", http.StatusInternalServerError)
        return
    }

    response := map[string]any {
        "message": "this is a protected endopoint",
        "user": user,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}
*/

// =======================================

/*
TODO: 

3. better logging ...

* better logging
* set up a generic UI for testing 
    * send file and then routing through Vue
    * or through html (not sold on Vue routing)
* create the dashboard routing
* set up the db
* create db methods
* api auth handler

*/

// for returning the login page
func loginPageHandler(w http.ResponseWriter, r *http.Request) {}

// WARN: requires the db to be set up
// and better user struct
func register(w http.ResponseWriter, r *http.Request) {}


// FIX: will have to think about what this does and how to handle
func dashboard(w http.ResponseWriter, r *http.Request) {}

// =======================================



// ========== HANDLER FUNCTIONS ==========
// TODO: set up and api handler

func addRoutes(h *http.ServeMux, s *Server) {
    jwt_manager := NewJWTManager(
        "test_key",
        time.Minute*10,
    )

    // basic routes
    h.HandleFunc("/", index)
    h.HandleFunc("/health", healthCheck)

    // user status routes
    h.HandleFunc("/login", loginHandler(jwt_manager))
    h.HandleFunc("/logout", logoutHandler(jwt_manager))
    h.HandleFunc("/reresh", refreshHandler(jwt_manager))
    // TODO:
    h.HandleFunc("/register", register)


    // application routes
    h.HandleFunc("/dashboard", jwt_manager.AuthMiddleware(
        s.rateLimitMiddleware(dashboard),
    ))

    // api auth handler
    h.HandleFunc("/api/chat", jwt_manager.AuthApiMiddleware(
        s.rateLimitMiddleware(chatHandler),
    ))

    // HACK: need to add auth around this
    h.HandleFunc("/api/stream", streamHandler(s.logger))
    h.HandleFunc("/api/modelList", modelListHandler)
    h.HandleFunc("/api/modelInfo", modelShowHandler)

	h.HandleFunc("/api/chatDescription", chatDescriptionHandler)
}

// =======================================
