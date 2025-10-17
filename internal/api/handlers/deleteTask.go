package handlers

import (
	"errors"
	"net/http"

	"github.com/Minkaill/planner-service.git/internal/db"
	"github.com/gin-gonic/gin"
)

func DeleteTaskHandler(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "отсутствует параметр id"})
		return
	}

	err := db.DeleteTask(db.DB, id)
	if err != nil {
		if errors.Is(err, db.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
