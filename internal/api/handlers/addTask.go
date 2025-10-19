package handlers

import (
	"net/http"

	"github.com/Minkaill/planner-service.git/internal/db"
	"github.com/Minkaill/planner-service.git/internal/models"
	"github.com/Minkaill/planner-service.git/pkg/utils"
	"github.com/gin-gonic/gin"
)

func AddTaskHandler(c *gin.Context) {
	var t models.Task

	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат JSON"})
		return
	}

	if t.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Заголовок не может быть пустым"})
		return
	}

	if err := utils.CheckDate(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := db.AddTask(&t)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при добавлении задачи"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id": id,
	})
}
