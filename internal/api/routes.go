package api

import (
	"github.com/Minkaill/planner-service.git/internal/api/handlers"
	"github.com/Minkaill/planner-service.git/internal/api/middleware"
	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.Engine) {
	r.POST("/api/signin", handlers.SignInHandler)

	protected := r.Group("/api")

	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/nextdate", handlers.NextDateHandler)
		protected.GET("/tasks", handlers.GetTasksHandler)
		protected.GET("/task", handlers.GetTaskHandler)
		protected.POST("/task", handlers.AddTaskHandler)
		protected.POST("/task/done", handlers.TaskDoneHandler)
		protected.PUT("/task", handlers.UpdateTaskHandler)
		protected.DELETE("/task", handlers.DeleteTaskHandler)
	}
}
