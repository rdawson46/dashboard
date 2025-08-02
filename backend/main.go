package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/rdawson46/dashboard/server"
)

// TODO: need to run migrations on startup
func run() error {
	db, err := server.NewSqliteRepository()

	if err != nil {
		return err
	}
	
	defer db.Close()

    config := server.NewConfig(8000, 300, 100, 100)
    s := server.NewServer(config, db)

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
