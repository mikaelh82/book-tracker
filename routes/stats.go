package routes

import (
	"net/http"

	"book-tracker/handlers"
)

func SetupStatsRoutes(mux *http.ServeMux, handler *handlers.StatsHandler) {
	// Note:
	// Handle GET /api/v1/stats.
	mux.HandleFunc("/api/v1/stats", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" { // This endpoint is only related to GET i.e. GET all stats
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handler.GetStats(w, r)
	})
}
