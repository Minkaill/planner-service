package server

import (
	"fmt"
	"os"

	"github.com/Minkaill/planner-service.git/internal/api"
	"github.com/Minkaill/planner-service.git/internal/db"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func Run() error {
	if err := godotenv.Load(); err != nil {
		fmt.Println(".env не найден, используются значения по умолчанию")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "7540"
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "scheduler.db"
	}

	if err := db.Init(dbPath); err != nil {
		return fmt.Errorf("ошибка инициализации БД: %w", err)
	}

	r := gin.Default()

	api.InitRoutes(r)

	r.Static("/web", "./web")

	addr := fmt.Sprintf(":%s", port)
	fmt.Printf("Сервер запущен на http://localhost%s\n", addr)
	return r.Run(addr)
}
