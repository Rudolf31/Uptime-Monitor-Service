package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"uptime-monitor/internal/handlers"
	"uptime-monitor/internal/storage"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,

		AddSource: true,
	}))
	slog.SetDefault(logger)

	store := storage.NewStorage()
	monitorHandler := handlers.NewMonitorHandler(store)

	mux := http.NewServeMux()

	mux.Handle("POST /monitors", http.HandlerFunc(monitorHandler.Create))
	mux.Handle("GET /monitors", http.HandlerFunc(monitorHandler.List))
	mux.Handle("GET /monitors/{id}", http.HandlerFunc(monitorHandler.Get))
	mux.Handle("DELETE /monitors/{id}", http.HandlerFunc(monitorHandler.Delete))

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  69 * time.Second,
	}

	go func() {
		logger.Info("server started",
			"url", "http://localhost:"+port,
			"port", port,
		)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Ошибка запуска сервера",
				"Error", err,
			)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	logger.Info("Начинаем graceful shutdown",
		"signal", sig,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Ошибка при shutdown",
			"error", err,
			"action", "force_close",
		)

		server.Close()
	}

	logger.Info("Сервер остановлен")
}
