package handlers

import (
	"encoding/json"
	"habit-tracker/internal/models"
	"habit-tracker/internal/storage"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func getUserID(r *http.Request) int {
	idStr := r.Header.Get("X-User-ID")
	if idStr == "" {
		return 1
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 1
	}
	return id
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, models.ErrorResponse{Error: message})
}

func CreateHabit(w http.ResponseWriter, r *http.Request) {
	var req models.CreateHabitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Title == "" {
		respondError(w, http.StatusBadRequest, "title is required")
		return
	}

	if req.GoalPerDay <= 0 {
		req.GoalPerDay = 1
	}

	userID := getUserID(r)
	habit, err := storage.CreateHabit(userID, &req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to create habit")
		return
	}

	respondJSON(w, http.StatusCreated, habit)
}

func GetHabit(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid habit id")
		return
	}

	userID := getUserID(r)
	habit, err := storage.GetHabitByID(id, userID)
	if err != nil {
		respondError(w, http.StatusNotFound, "habit not found")
		return
	}

	respondJSON(w, http.StatusOK, habit)
}

func ListHabits(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	habits, err := storage.GetAllHabits(userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to fetch habits")
		return
	}

	if habits == nil {
		habits = []models.Habit{}
	}

	respondJSON(w, http.StatusOK, habits)
}

func UpdateHabit(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid habit id")
		return
	}

	var req models.UpdateHabitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	userID := getUserID(r)
	habit, err := storage.UpdateHabit(id, userID, &req)
	if err != nil {
		respondError(w, http.StatusNotFound, "habit not found")
		return
	}

	respondJSON(w, http.StatusOK, habit)
}

func DeleteHabit(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid habit id")
		return
	}

	userID := getUserID(r)
	if err := storage.DeleteHabit(id, userID); err != nil {
		respondError(w, http.StatusNotFound, "habit not found")
		return
	}

	respondJSON(w, http.StatusOK, models.SuccessResponse{Message: "habit deleted"})
}

func LogHabit(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid habit id")
		return
	}

	var req models.LogHabitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		req.Count = 1
	}

	if req.Count <= 0 {
		req.Count = 1
	}

	userID := getUserID(r)
	if err := storage.LogHabitCompletion(id, userID, req.Count); err != nil {
		respondError(w, http.StatusNotFound, "habit not found")
		return
	}

	respondJSON(w, http.StatusOK, models.SuccessResponse{Message: "habit logged"})
}

func GetHabitLogs(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid habit id")
		return
	}

	userID := getUserID(r)
	logs, err := storage.GetHabitLogs(id, userID)
	if err != nil {
		respondError(w, http.StatusNotFound, "habit not found")
		return
	}

	if logs == nil {
		logs = []models.HabitLog{}
	}

	respondJSON(w, http.StatusOK, logs)
}

func GetHabitStats(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid habit id")
		return
	}

	userID := getUserID(r)
	stats, err := storage.GetHabitStats(id, userID)
	if err != nil {
		respondError(w, http.StatusNotFound, "habit not found")
		return
	}

	respondJSON(w, http.StatusOK, stats)
}
