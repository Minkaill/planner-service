package db

import (
	"database/sql"

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
