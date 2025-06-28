package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// ========== ROUTING STRUCTS ==========
type HealthResponse struct {
    Status string `json:"status"`
}

type TestResponse struct {
    Message string `json:"status"`
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

func chat(w http.ResponseWriter, r *http.Request) {
    response := TestResponse{Message: "chat"}
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)
}

func search(w http.ResponseWriter, r *http.Request) {
    response := TestResponse{Message: "search"}
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)
}

// ========== JWT (temp) FUNCTIONS ==========

type User struct {
    Username string `json:"username"`
    ID int `json:"id"`
}


type Claims struct {
    Username string `json:"username"`
    jwt.RegisteredClaims
}

type JWTManager struct {
    secretKey []byte
    tokenDuration time.Duration
}

func NewJWTManager(secretKey string, tokenDuration time.Duration) *JWTManager {
    return &JWTManager{
        secretKey: []byte(secretKey),
        tokenDuration: tokenDuration,
    }
}

func (manager *JWTManager) GenerateToken(user *User) (string, error) {
    claims := &Claims{
        Username: user.Username,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(manager.tokenDuration)),
            IssuedAt: jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
            Issuer: "server",
            Subject: user.Username,
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    return token.SignedString(manager.secretKey)
}

func (manager *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }

        return manager.secretKey, nil
    })

    if err != nil {
        return nil, err
    }

    claims, ok := token.Claims.(*Claims)

    if !ok || !token.Valid {
        return nil, fmt.Errorf("invalid token")
    }

    return claims, nil
}

func (manager  *JWTManager) RefreshToken(tokenString string) (string, error) {
    claims, err := manager.ValidateToken(tokenString)

    if err != nil {
        return "", err
    }

    newClaims := &Claims{
        Username: claims.Username,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(manager.tokenDuration)),
            IssuedAt: jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
            Issuer: "server",
            Subject: claims.Username,
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
    return token.SignedString(manager.secretKey)
}

func (manager *JWTManager) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")

        if authHeader == "" {
            // TODO: split between api and page routing handlers
            http.Error(w, "Authorization required", http.StatusUnauthorized)
            http.Redirect(w, r, "/login", http.StatusUnauthorized)
            return
        }

        tokenParts := strings.Split(authHeader, " ")

        if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
            http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
            http.Redirect(w, r, "/login", http.StatusUnauthorized)
            return
        }

        tokenString := tokenParts[1]
        claims, err := manager.ValidateToken(tokenString)

        if err != nil {
            http.Error(w, "Invalid Token", http.StatusUnauthorized)
            return
        }

        r = r.WithContext(contextWithUser(r.Context(), &User{
            Username: claims.Username,
        }))

        next(w, r)
    }
}

func (manager *JWTManager) SetTokenCookie(w http.ResponseWriter, token string) {
    cookie := &http.Cookie{
        Name: "auth_token",
        Value: token,
        Expires: time.Now().Add(manager.tokenDuration),
        HttpOnly: true,
        Secure: true,
        SameSite: http.SameSiteStrictMode,
        Path: "/",
    }

    http.SetCookie(w, cookie)
}

func (manager *JWTManager) GetTokenFromCookie(r *http.Request) (string, error) {
    cookie, err := r.Cookie("auth_token")

    if err != nil {
        return "", err
    }
    return cookie.Value, nil
}

func (manager *JWTManager) ClearTokenCookie(w http.ResponseWriter) {
    cookie := &http.Cookie{
        Name: "auth_token",
        Value: "",
        Expires: time.Now().Add(manager.tokenDuration),
        HttpOnly: true,
        Secure: true,
        SameSite: http.SameSiteStrictMode,
        Path: "/",
    }

    http.SetCookie(w, cookie)
}

type contextKey string

const userContextKey contextKey = "user"

func contextWithUser(ctx context.Context, user *User) context.Context {
    return context.WithValue(ctx, userContextKey, user)
}

func userFromContext(ctx context.Context) (*User, bool) {
    user, ok := ctx.Value(userContextKey).(*User)
    return user, ok
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
            fmt.Printf("username: %s, password: %s\n", loginReq.Username, loginReq.Password)
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

func refreshHandler(jwtManager *JWTManager) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Authorization header required", http.StatusUnauthorized)
            return
        }

        tokeParts := strings.Split(authHeader, " ")
        if len(tokeParts) != 2 || tokeParts[0] != "Bearer" {
            http.Error(w, "Invalid header format", http.StatusUnauthorized)
            return
        }

        newToken, err := jwtManager.RefreshToken(tokeParts[1])
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

// TODO: will require jwt set up and the db
// set up basic UI and chats
// need to setup the jwt middleware
    // either in this file or server
func register(w http.ResponseWriter, r *http.Request) {}
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
    // TODO: api handler and middleware
    h.HandleFunc("/dashboard", jwt_manager.AuthMiddleware(
        s.rateLimitMiddleware(dashboard),
    ))
    /*
    h.HandleFunc("api/search", jwt_manager.AuthMiddleware(
        s.rateLimitMiddleware(search),
    ))

    h.HandleFunc("api/chat", jwt_manager.AuthMiddleware(
        s.rateLimitMiddleware(chat),
    ))
    */
}

// =======================================
