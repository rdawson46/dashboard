package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"github.com/rdawson46/dashboard/db"
	"github.com/rdawson46/dashboard/jobs"
	"github.com/rdawson46/dashboard/server"
)

// TODO: need to run migrations on startup
func run() error {
    logger := log.NewWithOptions(os.Stderr, log.Options{
        ReportCaller: true,
    })

	sqliteDb, err := db.NewSqliteRepository(logger)

	if err != nil {
		return err
	}

	jp := jobs.NewJobPipeline(logger, sqliteDb.Db)

	ctx := context.Background()
	go jp.StartJobCheck(ctx, 5)
	
	defer sqliteDb.Close()

    config := server.NewConfig(8000, 300, 100, 100)
    s := server.NewServer(config, db.Repository(sqliteDb), logger)

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

	server.InitialEnvCheck()

    if err := run(); err != nil {
        log.Fatal("Error running server:", err)
    }
}
