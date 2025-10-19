package handlers

import (
	"net/http"

	"github.com/Minkaill/planner-service.git/internal/db"
	"github.com/gin-gonic/gin"
)

func GetTaskHandler(c *gin.Context) {
	id := c.Query("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "отсутствует параметр id"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if task == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Задача не найдена"})
		return
	}

	c.JSON(http.StatusOK, task)
}
