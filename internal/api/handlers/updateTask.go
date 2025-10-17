package handlers

import (
	"database/sql"
	"net/http"

	"github.com/Minkaill/planner-service.git/internal/db"
	"github.com/Minkaill/planner-service.git/internal/models"
	"github.com/gin-gonic/gin"
)

func UpdateTaskHandler(c *gin.Context) {
	var t models.Task
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный JSON"})
		return
	}

	if t.ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "отсутствует id задачи"})
		return
	}

	err := db.UpdateTask(db.DB, &t)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "задача не найдена"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "задача обновлена"})
}
