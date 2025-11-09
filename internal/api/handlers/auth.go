package handlers

import (
	"net/http"
	"os"

	"github.com/Minkaill/planner-service.git/pkg/utils"
	"github.com/gin-gonic/gin"
)

type SignInRequest struct {
	Password string `json:"password"`
}

// Тестовый коммит для проверки

func SignInHandler(c *gin.Context) {
	var req SignInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный запрос"})
		return
	}

	expected := os.Getenv("TODO_PASSWORD")
	if expected == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "TODO_PASSWORD не задан"})
		return
	}

	if req.Password != expected {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный пароль"})
		return
	}

	token, err := utils.GenerateJWT()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось создать токен"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
