package main

import (
	"fmt"
	"net/http"
	"os"
)

type Config struct {
	PORT string
}

// factory function to generate Config
func loadConfig() Config {
	cfg := Config{PORT: "8080"}

	if port := os.Getenv("BACKEND_PORT"); port != "" {
		cfg.PORT = port
	}

	return cfg
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "PONG")
}

func main() {

	var cfg Config = loadConfig()
	fmt.Printf("Using port: %s", cfg.PORT)

	http.HandleFunc("/ping", healthHandler)
	http.ListenAndServe(":"+cfg.PORT, nil) // NOTE: Interesting. Go handles "" as a string and '' as a Rune

}
