package storage

import (
	"database/sql"
	"log"
)

var DB *sql.DB

func InitDB(dbPath string) error {
	var err error
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	createdTables := `
	CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT UNIQUE NOT NULL
    );

    CREATE TABLE IF NOT EXISTS habits (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER NOT NULL,
        title TEXT NOT NULL,
        description TEXT,
        goal_per_day INTEGER DEFAULT 1,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        streak INTEGER DEFAULT 0,
        FOREIGN KEY (user_id) REFERENCES users(id)
    );

    CREATE TABLE IF NOT EXISTS habit_logs (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        habit_id INTEGER NOT NULL,
        date DATE DEFAULT CURRENT_DATE,
        count INTEGER DEFAULT 0,
        FOREIGN KEY (habit_id) REFERENCES habits(id),
        UNIQUE(habit_id, date)
    );
	`
	_, err = DB.Exec(createdTables)
	if err != nil {
		return err
	}

	err = CreatedDefaultUser()
	if err != nil {
		return err
	}

	log.Println("Database initialized succefully")
	return nil
}

// пользователь для тестов
func CreatedDefaultUser() error {
	_, err := DB.Exec(`INSERT OR IGNORE INTO users (id, username) VALUES (1, 'default_user')`)
	return err
}
