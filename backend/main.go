package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/rdawson46/dashboard/server"
)

func run() error {
    config := server.NewConfig(8000, 30, 10, 10)
    s := server.NewServer(config)

    if err := s.Start(); err != nil {
        log.Fatal(err)
    }

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

    <-quit
    return s.Shutdown()
}

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file:", err)
    }

    if err := run(); err != nil {
        log.Fatal("Error running server:", err)
    }
}
