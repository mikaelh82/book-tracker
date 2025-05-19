package main

import (
	"book-tracker/handlers"
	"book-tracker/middleware"
	"book-tracker/routes"
	"book-tracker/services"
	"book-tracker/store"
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type Config struct {
	PORT          string
	Timeout       time.Duration
	DBPath        string
	AllowedOrigin []string
}

func loadConfig() Config {
	cfg := Config{
		PORT:          "8080",
		Timeout:       10 * time.Second,
		DBPath:        "books.db",
		AllowedOrigin: []string{"http://localhost:5173"},
	}

	if port := os.Getenv("BACKEND_PORT"); port != "" {
		cfg.PORT = port
	}

	if timeout := os.Getenv("HTTP_TIMEOUT"); timeout != "" {
		if d, err := time.ParseDuration(timeout); err == nil {
			cfg.Timeout = d
		}
	}

	if dbPath := os.Getenv("DB_PATH"); dbPath != "" {
		cfg.DBPath = dbPath
	}

	if origin := os.Getenv("ALLOWED_ORIGIN"); origin != "" {
		origins := strings.Split(origin, ",")
		for i, o := range origins {
			origins[i] = strings.TrimSpace(o)
		}
		if len(origins) > 0 && origins[0] != "" {
			cfg.AllowedOrigin = origins
		}
	}

	return cfg
}

// NOTE: Could be moved somewhere else but lets keep it here for now
func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "OK")
}

func main() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	var cfg Config = loadConfig()

	db, closeDB, err := store.NewDB(cfg.DBPath)
	if err != nil {
		logger.Error("Failed to initialize DuckDB", "error", err)
		os.Exit(1)
	}
	defer closeDB()

	bookStore := store.NewBookStore(db)
	statsStore := store.NewStatsStore(db)

	bookService := services.NewBookService(bookStore)
	statsService := services.NewStatsService(statsStore)

	bookHandler := handlers.NewBookHandler(bookService)
	statsHandler := handlers.NewStatsHandler(statsService)

	mux := http.NewServeMux()
	routes.SetupBooksRoutes(mux, bookHandler)
	routes.SetupStatsRoutes(mux, statsHandler)
	mux.HandleFunc("/api/v1/health", healthHandler)
	mux.Handle("/metrics", middleware.MetricsHandler())

	handler := middleware.NewChain(
		middleware.CORS(cfg.AllowedOrigin),
		middleware.Logging(logger),
		middleware.Metrics,
		middleware.Timeout(cfg.Timeout),
	).Then(mux)

	srv := &http.Server{
		Addr:         ":" + cfg.PORT,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed { // NOTE: Need to be added to work with Graceful shutdown
			log.Fatalf("Server error: %s", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM) // os.Interrupt: Control+c. SIGTERM process managers i.e. Kubernetes, Docker etc
	<-sigChan

	logger.Info("Shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Shutdown error", "error", err)
		os.Exit(1)
	}
	logger.Info("Server stopped")

}
