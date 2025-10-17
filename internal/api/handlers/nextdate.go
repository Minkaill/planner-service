package handlers

import (
	"net/http"
	"time"

	"github.com/Minkaill/planner-service.git/pkg/utils"
	"github.com/gin-gonic/gin"
)

func NextDateHandler(c *gin.Context) {
	nowStr := c.Query("now")
	dateStr := c.Query("date")
	repeat := c.Query("repeat")

	if dateStr == "" || repeat == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Параметры 'date' и 'repeat' обязательны"})
		return
	}

	var now time.Time
	var err error
	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(utils.DateFormat, nowStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Недействительный формат даты 'now'"})
			return
		}
	}

	next, err := utils.NextDate(now, dateStr, repeat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Ошибка": err.Error()})
		return
	}

	c.String(http.StatusOK, next)
}
