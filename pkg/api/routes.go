package api

import (
	"github.com/Minkaill/planner-service.git/pkg/api/handlers"
	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.Engine) {
	r.GET("/api/nextdate", handlers.NextDateHandler)
}
