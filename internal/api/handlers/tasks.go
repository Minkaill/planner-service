package handlers

import (
	"net/http"

	"github.com/Minkaill/planner-service.git/internal/db"
	"github.com/Minkaill/planner-service.git/internal/models"
	"github.com/gin-gonic/gin"
)

func GetTasksHandler(c *gin.Context) {
	search := c.Query("search")
	
	tasks, err := db.GetTasks(db.DB, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if tasks == nil {
		tasks = []models.Task{}
	}

	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}
