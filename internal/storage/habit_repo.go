package storage

import (
	"habit-tracker/internal/models"
	"time"
)

func CreateHabit(habit *models.Habit) error {
	query := `INSERT INTO habits (user_id, title, description, goal_per_day) VALUES (?, ?, ?, ?)`

	res, err := DB.Exec(query, habit.UserID, habit.Title, habit.Description, habit.GoalPerDay)
	if err != nil {
		return err
	}

	id, _ := res.LastInsertId()
	habit.ID = int(id)
	return nil
}

func GetAllHabits(userID int) ([]models.Habit, error) {
	rows, err := DB.Query(`SELECT id, user_id, title, description, goal_per_day, created_at, streak FROM habits WHERE user_id = ?`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var habits []models.Habit
	for rows.Next() {
		var h models.Habit
		err := rows.Scan(&h.ID, &h.UserID, &h.Title, &h.Description, &h.GoalPerDay, &h.CreatedAt, &h.Streak)
		if err != nil {
			return nil, err
		}
		habits = append(habits, h)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return habits, nil
}

func LogHabit(habitID int, count int) error {
	today := time.Now().Format("2006-01-02")
	query := `INSERT INTO habit_logs (habit_id, date, count) VALUES (?, ?, ?)
	ON CONFLICT(habit_id, date) DO UPDATE SET count = count + ?`

	_, err := DB.Exec(query, habitID, today, count, count)
	if err != nil {
		return err
	}

	_, err = DB.Exec(`UPDATE habits SET streak = streak + 1 WHERE id = ?`, habitID)

	return err
}
