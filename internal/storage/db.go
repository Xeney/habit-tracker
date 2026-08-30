package storage

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB(dbPath string) error {
	var err error
	DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}

	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL
	);

	CREATE TABLE IF NOT EXISTS habits (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		title TEXT NOT NULL,
		description TEXT DEFAULT '',
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

	_, err = DB.Exec(schema)
	if err != nil {
		return err
	}

	_, err = DB.Exec(`INSERT OR IGNORE INTO users (id, username) VALUES (1, 'default_user')`)
	return err
}

func Close() {
	if DB != nil {
		DB.Close()
	}
}
