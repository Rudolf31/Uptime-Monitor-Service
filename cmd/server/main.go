package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"uptime-monitor/internal/db"
	"uptime-monitor/internal/handlers"
	"uptime-monitor/internal/middleware"
	"uptime-monitor/internal/storage"
	"uptime-monitor/internal/worker"

	"github.com/joho/godotenv"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	godotenv.Load()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	}))
	slog.SetDefault(logger)

	var store storage.Storage
	switch os.Getenv("STORAGE") {
	case "postgres":
		sqlDB, err := db.NewDBFromEnv()
		if err != nil {
			logger.Error("Failed to connect to Postgres", "error", err)
			return
		}
		store = storage.NewPostgresStorage(sqlDB)
	default:
		store = storage.NewMemoryStorage()
	}

	pool := worker.NewPool(store)
	scheduler := worker.NewScheduler(store, pool.Jobs())

	poolCtx, cancelPool := context.WithCancel(context.Background())
	schedulerCtx, cancelScheduler := context.WithCancel(context.Background())

	var schedulerWg sync.WaitGroup
	schedulerWg.Add(1)
	go func() {
		defer schedulerWg.Done()
		scheduler.Start(schedulerCtx)
	}()
	go pool.Start(poolCtx)

	monitorHandler := handlers.NewMonitorHandler(store)

	mux := http.NewServeMux()

	commonMiddleware := middleware.Chain(
		middleware.RecoveryMiddleware,
		middleware.LoggingMiddleware,
	)

	mux.Handle("POST /monitors", commonMiddleware(http.HandlerFunc(monitorHandler.Create)))
	mux.Handle("GET /monitors", commonMiddleware(http.HandlerFunc(monitorHandler.List)))
	mux.Handle("GET /monitors/{id}", commonMiddleware(http.HandlerFunc(monitorHandler.Get)))
	mux.Handle("DELETE /monitors/{id}", commonMiddleware(http.HandlerFunc(monitorHandler.Delete)))

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		storage := os.Getenv("STORAGE")
		if storage == "" {
			storage = "memory"
		}
		logger.Info("Server started",
			"url", "http://localhost:"+port,
			"port", port,
			"storage", storage,
		)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Failed to start server",
				"error", err,
			)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan

	logger.Info("Starting graceful shutdown",
		"signal", sig,
	)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("Failed to shutdown server",
			"error", err,
			"action", "force_close",
		)
		server.Close()
	}

	cancelScheduler()
	schedulerWg.Wait()

	cancelPool()
	pool.Stop()

	logger.Info("Server stopped")
}
