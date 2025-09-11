package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)


// ========== JWT FUNCTIONS ==========

// TODO: add logger to jwt
type User_jwt struct {
    Username string `json:"username"`
    ID int64 `json:"id"`
}


type Claims struct {
    Username string `json:"username"`
	ID int64 `json:"id"`
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

func (manager *JWTManager) GenerateToken(user *User_jwt) (string, error) {
    claims := &Claims{
        Username: user.Username,
		ID: user.ID,
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
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
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
		ID: claims.ID,
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

// TODO: split between api and page routing handlers
func (manager *JWTManager) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        token, err := manager.GetTokenFromCookie(r)

        if err != nil {
            if err == http.ErrNoCookie {
                http.Error(w, "Unauthorized: No session cookie found", http.StatusUnauthorized)
                return
            }
            http.Error(w, "Bad request", http.StatusUnauthorized)
            return
        }

        claims, err := manager.ValidateToken(token)

        if err != nil {
            http.Error(w, "Invalid Token", http.StatusUnauthorized)
            return
        }

        r = r.WithContext(contextWithUser(r.Context(), &User_jwt{
            Username: claims.Username,
			ID: claims.ID,
        }))

        next(w, r)
    }
}

// HACK: will have to decide when building the UI if the 2 auth funcs need seperation
func (manager *JWTManager) AuthApiMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        token, err := manager.GetTokenFromCookie(r)

        if err != nil {
            if err == http.ErrNoCookie {
                http.Error(w, "Unauthorized: No session cookie found", http.StatusUnauthorized)
                return
            }

            http.Error(w, "Bad request", http.StatusUnauthorized)
            return
        }

        claims, err := manager.ValidateToken(token)

        if err != nil {
            http.Error(w, "Invalid Token", http.StatusUnauthorized)
            return
        }

        r = r.WithContext(contextWithUser(r.Context(), &User_jwt{
            Username: claims.Username,
			ID: claims.ID,
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

