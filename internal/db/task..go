package db

import (
	"database/sql"
	"strings"
	"time"

	"github.com/Minkaill/planner-service.git/internal/models"
)

func AddTask(db *sql.DB, t *models.Task) (int64, error) {
	var id int64

	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)`
	res, err := db.Exec(query, sql.Named("date", t.Date), sql.Named("title", t.Title), sql.Named("comment", t.Comment), sql.Named("repeat", t.Repeat))
	if err == nil {
		id, err = res.LastInsertId()
	}

	return id, err
}

func GetTasks(db *sql.DB, search string) ([]models.Task, error) {
	var (
		rows *sql.Rows
		err  error
	)

	if search == "" {
		rows, err = db.Query(`SELECT id, date, title, comment, repeat from scheduler`)
		if err != nil {
			return nil, err
		}

		defer rows.Close()
		return scanTasks(rows)
	}

	if t, parseErr := time.Parse("02.01.2006", search); parseErr == nil {
		search = t.Format("20060102")
		rows, err = db.Query(`
			SELECT id, date, title, comment, repeat
			FROM scheduler
			WHERE date = :date
		`, sql.Named("date", search))
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		return scanTasks(rows)
	}

	search = strings.TrimSpace(search)
	rows, err = db.Query(`
		SELECT id, date, title, comment, repeat
		FROM scheduler
		WHERE title LIKE '%' || :search || '%' COLLATE NOCASE
		   OR comment LIKE '%' || :search || '%' COLLATE NOCASE
	`, sql.Named("search", search))
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	return scanTasks(rows)
}

func scanTasks(rows *sql.Rows) ([]models.Task, error) {
	var tasks []models.Task

	for rows.Next() {
		var t models.Task

		if err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}
