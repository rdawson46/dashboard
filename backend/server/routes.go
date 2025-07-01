package server

import (
	"encoding/json"
	"net/http"
	"time"

	api "github.com/rdawson46/dashboard/api"
)

// ========== ROUTING STRUCTS ==========
type HealthResponse struct {
    Status string `json:"status"`
}

// =====================================


// ========== ROUTING FUNCTIONS ==========
func index(w http.ResponseWriter, r *http.Request) {
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

    // HACK: placed here temp
    url := "http://10.0.2.2:11434"
    ctx := r.Context()

    if ctx == nil {
        http.Error(w, "context is nil", http.StatusInternalServerError)
        return
    }

    oc, err := api.NewOllamaClient(url, ctx)

    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return 
    }

    res, err := oc.Chat(chatReq.Query)

    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return 
    }

    response := map[string]string{"respones": res}
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)
}

func search(w http.ResponseWriter, r *http.Request) {
    response := map[string]string{"message": "chat"}
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)
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

1. get an ollama client running
2. set up the api/chat endpoint
3. better logging ...

* better logging
* set up a generic UI for testing 
    * send file and then routing through Vue
    * or through html (not sold on Vue routing)
* move jwt stuff to a new file
* create the dashboard routing
* set up the db
* create db methods
* set up Ollama client
* api auth handler

*/

// WARN: requires the db to be set up
// and better user struct
func register(w http.ResponseWriter, r *http.Request) {}


// FIX: will have to think about what this does and how to handle
func dashboard(w http.ResponseWriter, r *http.Request) {}

// =======================================



// ========== HANDLER FUNCTIONS ==========
// TODO: set up and api handler

func addRoutes(h *http.ServeMux, s *Server) {
    jwt_manager := NewJWTManager("test_key", time.Minute*10)


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
    h.HandleFunc("api/search", jwt_manager.AuthMiddleware(
        s.rateLimitMiddleware(search),
    ))

    h.HandleFunc("/api/chat", chatHandler)
    /*
    h.HandleFunc("api/chat", jwt_manager.AuthMiddleware(
        s.rateLimitMiddleware(chat),
    ))
    */
}

// =======================================
