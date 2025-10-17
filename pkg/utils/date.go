package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const DateFormat = "20060102"

func AfterNow(date, now time.Time) bool {
	y1, m1, d1 := date.Date()
	y2, m2, d2 := now.Date()
	return time.Date(y1, m1, d1, 0, 0, 0, 0, time.UTC).
		After(time.Date(y2, m2, d2, 0, 0, 0, 0, time.UTC))
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	repeat = strings.TrimSpace(repeat)
	if repeat == "" {
		return "", fmt.Errorf("правило повтора пустое")
	}

	date, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("ошибка формата даты")
	}

	parts := strings.Split(repeat, " ")

	switch parts[0] {

	case "d":
		if len(parts) != 2 {
			return "", fmt.Errorf("ошибка чтения формата для 'd'")
		}

		interval, err := strconv.Atoi(parts[1])
		if err != nil || interval <= 0 || interval > 400 {
			return "", fmt.Errorf("недействительный интервал для 'd'")
		}

		for {
			date = date.AddDate(0, 0, interval)
			if AfterNow(date, now) {
				break
			}
		}

		return date.Format(DateFormat), nil

	case "w":
		if len(parts) != 2 {
			return "", fmt.Errorf("ошибка чтения формата для 'w'")
		}

		daysStr := strings.Split(parts[1], ",")
		var days [8]bool

		for _, ds := range daysStr {
			n, err := strconv.Atoi(strings.TrimSpace(ds))
			if err != nil || n < 1 || n > 7 {
				return "", fmt.Errorf("недействительный день недели: %s", ds)
			}
			days[n] = true
		}

		date = date.AddDate(0, 0, 1)

		for {
			weekday := int(date.Weekday())

			if weekday == 0 {
				weekday = 7
			}

			if days[weekday] && AfterNow(date, now) {
				break
			}

			date = date.AddDate(0, 0, 1)
		}

		return date.Format(DateFormat), nil

	case "m":
		if len(parts) < 2 {
			return "", fmt.Errorf("ошибка чтения формата для 'm'")
		}

		dayStr := strings.Split(parts[1], ",")
		var days [32]bool
		lastDay, beforeLast := false, false

		for _, ds := range dayStr {
			ds = strings.TrimSpace(ds)

			if ds == "-1" {
				lastDay = true
				continue
			}
			if ds == "-2" {
				beforeLast = true
				continue
			}
			n, err := strconv.Atoi(ds)
			if err != nil || n < 1 || n > 31 {
				return "", fmt.Errorf("недействительное число дня: %s", ds)
			}

			days[n] = true
		}

		var months [13]bool
		hasMonths := false
		if len(parts) == 3 {
			monthsStr := strings.Split(parts[2], ",")
			for _, ms := range monthsStr {
				n, err := strconv.Atoi(strings.TrimSpace(ms))
				if err != nil || n < 1 || n > 12 {
					return "", fmt.Errorf("недействительный номер месяца: %s", ms)
				}
				months[n] = true
				hasMonths = true
			}
		}

		date = date.AddDate(0, 0, 1)
		for {
			d := date.Day()
			m := int(date.Month())

			if hasMonths && !months[m] {
				date = date.AddDate(0, 0, 1)
				continue
			}

			ok := false
			if days[d] {
				ok = true
			}

			y, mon := date.Year(), date.Month()
			firstNextMonth := time.Date(y, mon+1, 1, 0, 0, 0, 0, time.UTC)
			lastDayOfMonth := firstNextMonth.AddDate(0, 0, -1).Day()

			if lastDay && d == lastDayOfMonth {
				ok = true
			}
			if beforeLast && d == lastDayOfMonth-1 {
				ok = true
			}

			if ok && AfterNow(date, now) {
				break
			}

			date = date.AddDate(0, 0, 1)
		}

		return date.Format(DateFormat), nil

	case "y":
		if len(parts) != 1 {
			return "", fmt.Errorf("ошибка чтения формата для 'y'")
		}

		for {
			date = date.AddDate(1, 0, 0)
			if AfterNow(date, now) {
				break
			}
		}

		return date.Format(DateFormat), nil

	default:
		return "", fmt.Errorf("неподдерживаемое правило повторения: %s", repeat)
	}
}
