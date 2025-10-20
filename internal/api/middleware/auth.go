package middleware

import (
	"net/http"
	"os"

	"github.com/Minkaill/planner-service.git/pkg/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		pass := os.Getenv("TODO_PASSWORD")

		if pass == "" {
			c.Next()
			return
		}

		cookie, err := c.Cookie("token")
		if err != nil || !utils.ValidateJWT(cookie) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentification required"})
			c.Abort()
			return
		}

		c.Next()
	}
}
