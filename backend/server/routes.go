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


// =======================================

// TODO: will require jwt set up and the db
// set up basic UI and chats
// need to setup the jwt middleware
    // either in this file or server
func register(w http.ResponseWriter, r *http.Request) {}
func login(w http.ResponseWriter, r *http.Request) {}
func logout(w http.ResponseWriter, r *http.Request) {}
func refresh(w http.ResponseWriter, r *http.Request) {}
func dashboard(w http.ResponseWriter, r *http.Request) {}

// =======================================



// ========== HANDLER FUNCTIONS ==========
// TODO: set up and api handler

func addRoutes(h *http.ServeMux, s *Server) {
    h.HandleFunc("/", index)
    h.HandleFunc("/health", healthCheck)
    h.HandleFunc("/login", login)
    h.HandleFunc("/logout", logout)
    h.HandleFunc("/register", register)
    h.HandleFunc("/reresh", refresh)

    // protected routes
    jwt_manager := NewJWTManager("test_key", time.Minute*10)

    // limited routes: TODO: will need to wrap with jwt
    h.HandleFunc("/dashboard", jwt_manager.AuthMiddleware(
        s.rateLimitMiddleware(search),
    ))


    // TODO: api handler and middleware
    h.HandleFunc("api/search", jwt_manager.AuthMiddleware(
        s.rateLimitMiddleware(search),
    ))

    h.HandleFunc("api/chat", jwt_manager.AuthMiddleware(
        s.rateLimitMiddleware(chat),
    ))
}

// =======================================
