package handlers

import (
	"book-tracker/services"
	"encoding/json"
	"fmt"
	"net/http"
)

type StatsHandler struct {
	service services.StatsService
}

func NewStatsHandler(service services.StatsService) *StatsHandler {
	return &StatsHandler{service: service}
}

func (h *StatsHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.service.GetStats(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("Get stats error: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
