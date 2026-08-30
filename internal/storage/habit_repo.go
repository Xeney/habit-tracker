package storage

import (
	"database/sql"
	"habit-tracker/internal/models"
	"time"
)

func CreateHabit(userID int, req *models.CreateHabitRequest) (*models.Habit, error) {
	goal := req.GoalPerDay
	if goal <= 0 {
		goal = 1
	}

	res, err := DB.Exec(
		`INSERT INTO habits (user_id, title, description, goal_per_day) VALUES (?, ?, ?, ?)`,
		userID, req.Title, req.Description, goal,
	)
	if err != nil {
		return nil, err
	}

	id, _ := res.LastInsertId()

	return GetHabitByID(int(id), userID)
}

func GetHabitByID(id, userID int) (*models.Habit, error) {
	var h models.Habit
	err := DB.QueryRow(
		`SELECT id, user_id, title, description, goal_per_day, created_at, streak FROM habits WHERE id = ? AND user_id = ?`,
		id, userID,
	).Scan(&h.ID, &h.UserID, &h.Title, &h.Description, &h.GoalPerDay, &h.CreatedAt, &h.Streak)
	if err != nil {
		return nil, err
	}
	return &h, nil
}

func GetAllHabits(userID int) ([]models.Habit, error) {
	rows, err := DB.Query(
		`SELECT id, user_id, title, description, goal_per_day, created_at, streak FROM habits WHERE user_id = ? ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var habits []models.Habit
	for rows.Next() {
		var h models.Habit
		if err := rows.Scan(&h.ID, &h.UserID, &h.Title, &h.Description, &h.GoalPerDay, &h.CreatedAt, &h.Streak); err != nil {
			return nil, err
		}
		habits = append(habits, h)
	}

	return habits, rows.Err()
}

func UpdateHabit(id, userID int, req *models.UpdateHabitRequest) (*models.Habit, error) {
	habit, err := GetHabitByID(id, userID)
	if err != nil {
		return nil, err
	}

	if req.Title != nil {
		habit.Title = *req.Title
	}
	if req.Description != nil {
		habit.Description = *req.Description
	}
	if req.GoalPerDay != nil && *req.GoalPerDay > 0 {
		habit.GoalPerDay = *req.GoalPerDay
	}

	_, err = DB.Exec(
		`UPDATE habits SET title = ?, description = ?, goal_per_day = ? WHERE id = ? AND user_id = ?`,
		habit.Title, habit.Description, habit.GoalPerDay, id, userID,
	)
	if err != nil {
		return nil, err
	}

	return GetHabitByID(id, userID)
}

func DeleteHabit(id, userID int) error {
	result, err := DB.Exec(`DELETE FROM habits WHERE id = ? AND user_id = ?`, id, userID)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func LogHabitCompletion(habitID, userID, count int) error {
	habit, err := GetHabitByID(habitID, userID)
	if err != nil {
		return err
	}

	today := time.Now().Format("2006-01-02")
	_, err = DB.Exec(
		`INSERT INTO habit_logs (habit_id, date, count) VALUES (?, ?, ?) ON CONFLICT(habit_id, date) DO UPDATE SET count = count + ?`,
		habitID, today, count, count,
	)
	if err != nil {
		return err
	}

	var todayCount int
	err = DB.QueryRow(`SELECT count FROM habit_logs WHERE habit_id = ? AND date = ?`, habitID, today).Scan(&todayCount)
	if err != nil {
		return err
	}

	if todayCount >= habit.GoalPerDay {
		_, err = DB.Exec(`UPDATE habits SET streak = streak + 1 WHERE id = ? AND user_id = ?`, habitID, userID)
	} else {
		_, err = DB.Exec(`UPDATE habits SET streak = 0 WHERE id = ? AND user_id = ?`, habitID, userID)
	}

	return err
}

func GetHabitLogs(habitID, userID int) ([]models.HabitLog, error) {
	_, err := GetHabitByID(habitID, userID)
	if err != nil {
		return nil, err
	}

	rows, err := DB.Query(
		`SELECT id, habit_id, date, count FROM habit_logs WHERE habit_id = ? ORDER BY date DESC`,
		habitID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []models.HabitLog
	for rows.Next() {
		var l models.HabitLog
		if err := rows.Scan(&l.ID, &l.HabitID, &l.Date, &l.Count); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}

	return logs, rows.Err()
}

func GetHabitStats(habitID, userID int) (*models.StatsResponse, error) {
	habit, err := GetHabitByID(habitID, userID)
	if err != nil {
		return nil, err
	}

	stats := &models.StatsResponse{
		HabitID:       habitID,
		CurrentStreak: habit.Streak,
	}

	DB.QueryRow(`SELECT COUNT(*) FROM habit_logs WHERE habit_id = ?`, habitID).Scan(&stats.TotalLogs)

	DB.QueryRow(`SELECT COALESCE(MAX(streak), 0) FROM habits WHERE id = ?`, habitID).Scan(&stats.BestStreak)

	var totalDays int
	DB.QueryRow(`SELECT COUNT(DISTINCT date) FROM habit_logs WHERE habit_id = ? AND count >= (SELECT goal_per_day FROM habits WHERE id = ?)`, habitID, habitID).Scan(&totalDays)

	var totalTrackedDays int
	DB.QueryRow(`SELECT COUNT(DISTINCT date) FROM habit_logs WHERE habit_id = ?`, habitID).Scan(&totalTrackedDays)

	if totalTrackedDays > 0 {
		stats.CompletionRate = float64(totalDays) / float64(totalTrackedDays) * 100
	}

	return stats, nil
}
