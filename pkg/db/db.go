package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title VARCHAR(255) NOT NULL,
    date DATE NOT NULL,
    comment TEXT,
    repeat VARCHAR(100)
);

CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler(date);
`

var DB *sql.DB

func Init(dbFile string) error {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		fmt.Println("Создаётся новая база данных:", dbFile)
	}

	database, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("ошибка открытия БД: %w", err)
	}

	DB = database

	if _, err := DB.Exec(schema); err != nil {
		return fmt.Errorf("ошибка применения схемы: %w", err)
	}

	fmt.Println("База данных успешно инициализирована.")
	return nil
}
