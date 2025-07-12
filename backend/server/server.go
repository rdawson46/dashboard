package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"golang.org/x/time/rate"
)

type ServerConfig struct {
    Port int
    ShutdownTimeout time.Duration
    RateLimitReq float64
    RateLimitBurst int
}


func NewConfig(port, rateLimitBurst int, rateLimitReq float64, shutdownTimeout time.Duration) ServerConfig {
    return ServerConfig {
        Port: port,
        ShutdownTimeout: shutdownTimeout,
        RateLimitBurst: rateLimitBurst,
        RateLimitReq: rateLimitReq,
    }
}

type Server struct {
    config ServerConfig
    httpServer *http.Server
    rateLimiter *rate.Limiter
    logger *log.Logger
}

func NewServer(config ServerConfig) *Server {
    logger := log.NewWithOptions(os.Stderr, log.Options{
        ReportCaller: true,
    })

    return &Server{
        config: config,
        logger: logger,
        rateLimiter: rate.NewLimiter(
            rate.Limit(config.RateLimitReq),
            config.RateLimitBurst,
        ),
    }
}


// ========== MIDDLEWARE FUNCTIONS ==========
func (s *Server) rateLimitMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if !s.rateLimiter.Allow() {
            http.Error(w, "Rate limit Exceeded", http.StatusTooManyRequests)
            return
        }

        next.ServeHTTP(w, r)
    }
}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        s.logger.Info("Incoming request",
            "method", r.Method,
            "path", r.URL.Path,
            "remote", r.RemoteAddr,
        )

        next.ServeHTTP(w, r)

        duration := time.Since(start)
        s.logger.Info("Request handled",
            "method", r.Method,
            "path", r.URL.Path,
            "duration", duration,
        )
    })
}
// ==========================================



// ========== SERVER FUNCTIONS ==========

func (s *Server) Start() error {
    s.httpServer = &http.Server{
        Addr: fmt.Sprintf(":%d", s.config.Port),
        ReadTimeout: 5 * time.Second,
        // disabled for longer responses
        // WriteTimeout: 10 * time.Second,
    }

    mux := http.NewServeMux()
    addRoutes(mux, s)
    loggedMux := s.loggingMiddleware(mux)

    s.httpServer.Handler = loggedMux

    go func() {
        s.logger.Infof("Started listening on port %d", s.config.Port)

        if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
            s.logger.Errorf("Hit error: %s", err)
        }
    }()

    return nil
}

func (s *Server) Shutdown() error {
    ctx, cancel := context.WithTimeout(
        context.Background(),
        s.config.ShutdownTimeout,
    )

    defer cancel()

    if err := s.httpServer.Shutdown(ctx); err != nil {
        s.logger.Errorf("Error shutting down: %s", err)
        return err
    }

    s.logger.Info("Server shutdown gracefully")
    return nil
}
// ======================================
