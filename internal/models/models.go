package models

import "time"

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

type Habit struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	GoalPerDay  int       `json:"goal_per_day"`
	CreatedAt   time.Time `json:"created_at"`
	Streak      int       `json:"streak"`
}

type HabitLog struct {
	ID      int       `json:"id"`
	HabitID int       `json:"habit_id"`
	Date    string    `json:"date"`
	Count   int       `json:"count"`
}

type CreateHabitRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	GoalPerDay  int    `json:"goal_per_day"`
}

type UpdateHabitRequest struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	GoalPerDay  *int    `json:"goal_per_day,omitempty"`
}

type LogHabitRequest struct {
	Count int `json:"count"`
}

type StatsResponse struct {
	HabitID       int     `json:"habit_id"`
	TotalLogs     int     `json:"total_logs"`
	CurrentStreak int     `json:"current_streak"`
	BestStreak    int     `json:"best_streak"`
	CompletionRate float64 `json:"completion_rate"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}
