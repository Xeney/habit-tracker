package handlers

import (
	"habit-tracker/internal/models"
	"habit-tracker/internal/storage"
	"net/http"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	if storage.DB != nil {
		respondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	} else {
		respondJSON(w, http.StatusServiceUnavailable, models.ErrorResponse{Error: "database not connected"})
	}
}
