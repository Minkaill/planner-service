package api

import (
	"github.com/Minkaill/planner-service.git/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.Engine) {
	r.GET("/api/nextdate", handlers.NextDateHandler)
	r.GET("/api/tasks", handlers.GetTasksHandler)
	r.POST("/api/task", handlers.AddTaskHandler)
}
