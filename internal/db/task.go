package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Minkaill/planner-service.git/internal/models"
	"github.com/Minkaill/planner-service.git/pkg/utils"
)

var ErrTaskNotFound = errors.New("задача не найдена")

func AddTask(t *models.Task) (int64, error) {
	var id int64

	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)`
	res, err := DB.Exec(query, sql.Named("date", t.Date), sql.Named("title", t.Title), sql.Named("comment", t.Comment), sql.Named("repeat", t.Repeat))
	if err == nil {
		id, err = res.LastInsertId()
	}

	return id, err
}

func GetTasks(search string) ([]models.Task, error) {
	var (
		rows *sql.Rows
		err  error
	)

	if search == "" {
		rows, err = DB.Query(`SELECT id, date, title, comment, repeat from scheduler`)
		if err != nil {
			return nil, err
		}

		defer rows.Close()
		return scanTasks(rows)
	}

	if t, parseErr := time.Parse("02.01.2006", search); parseErr == nil {
		search = t.Format("20060102")
		rows, err = DB.Query(`
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
	rows, err = DB.Query(`
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

func GetTask(id string) (*models.Task, error) {
	var t models.Task

	row := DB.QueryRow(`SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id`, sql.Named("id", id))
	err := row.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // если задачи с таким id нет
		}
	}

	return &t, nil
}

func UpdateTask(t *models.Task) error {
	if strings.TrimSpace(t.Title) == "" {
		return fmt.Errorf("заголовок не может быть пустым")
	}

	if err := utils.CheckDate(t); err != nil {
		return err
	}

	query := `
		UPDATE scheduler
		SET date = :date,
		    title = :title,
		    comment = :comment,
		    repeat = :repeat
		WHERE id = :id
	`

	res, err := DB.Exec(query,
		sql.Named("date", t.Date),
		sql.Named("title", t.Title),
		sql.Named("comment", t.Comment),
		sql.Named("repeat", t.Repeat),
		sql.Named("id", t.ID),
	)

	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrTaskNotFound
	}

	return nil
}

func DeleteTask(id string) error {
	query := `DELETE FROM scheduler WHERE id = :id`

	res, err := DB.Exec(query, sql.Named("id", id))
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrTaskNotFound
	}

	return nil
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

func UpdateDate(next string, id string) error {
	query := `UPDATE scheduler SET date = :date WHERE id = :id`

	res, err := DB.Exec(query, sql.Named("date", next), sql.Named("id", id))
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrTaskNotFound
	}

	return nil
}

func TaskDone(id string) error {
	task, err := GetTask(id)
	if err != nil {
		return err
	}

	if task.Repeat == "" {
		DeleteTask(id)
		return nil
	}

	now := time.Now()
	next, err := utils.NextDate(now, task.Date, task.Repeat)
	if err != nil {
		return err
	}

	if err := UpdateDate(next, id); err != nil {
		return err
	}

	return nil
}
