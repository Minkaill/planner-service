package utils

import (
	"errors"
	"time"

	"github.com/Minkaill/planner-service.git/internal/models"
)

func CheckDate(t *models.Task) error {
	now := time.Now()

	if t.Date == "" {
		t.Date = now.Format(DateFormat)
		return nil
	}

	time, err := time.Parse(DateFormat, t.Date)
	if err != nil {
		return errors.New("неверный формат даты (ожидается YYYYMMDD)")
	}

	if AfterNow(now, time) {
		if len(t.Repeat) == 0 {
			t.Date = now.Format(DateFormat)
		} else {
			next, err := NextDate(now, t.Date, t.Repeat)
			if err != nil {
				return err
			}
			t.Date = next
		}
	}

	return nil
}
