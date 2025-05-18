package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type Config struct {
	PORT string
}

func loadConfig() Config {
	cfg := Config{PORT: "8080"}

	if port := os.Getenv("BACKEND_PORT"); port != "" {
		cfg.PORT = port
	}

	return cfg
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "OK")
}

func main() {

	var cfg Config = loadConfig()
	log.Printf("Loaded config with PORT: %s", cfg.PORT)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/health", healthHandler)

	srv := http.Server{
		Addr:         ":" + cfg.PORT,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("Server error: %s", err)
		}
	}()

	select {} // TODO: Lets block the main thread for now. later graceful shutdown

}
