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
	Date    time.Time `json:"date"`
	Count   int       `json:"count"`
}
