package dateparse

import (
	"fmt"
	"time"
)

func ParseDate(input string) (time.Time, error) {
	now := time.Now()

	switch input {
	case "today":
		return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()), nil

	case "tomorrow":
		tomorrow := now.AddDate(0, 0, 1)
		return time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, tomorrow.Location()), nil

	case "end-of-week":
		daysUntilFriday := (5 - int(now.Weekday()) + 7) % 7
		if daysUntilFriday == 0 {
			daysUntilFriday = 7
		}
		friday := now.AddDate(0, 0, daysUntilFriday)
		return time.Date(friday.Year(), friday.Month(), friday.Day(), 0, 0, 0, 0, friday.Location()), nil

	case "end-of-month":
		firstOfNextMonth := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location())
		lastOfMonth := firstOfNextMonth.AddDate(0, 0, -1)
		return lastOfMonth, nil

	case "next-week":
		nextWeek := now.AddDate(0, 0, 7)
		return time.Date(nextWeek.Year(), nextWeek.Month(), nextWeek.Day(), 0, 0, 0, 0, nextWeek.Location()), nil

	case "next-month":
		nextMonth := now.AddDate(0, 1, 0)
		return time.Date(nextMonth.Year(), nextMonth.Month(), nextMonth.Day(), 0, 0, 0, 0, nextMonth.Location()), nil

	default:
		t, err := time.Parse("2006-01-02", input)
		if err != nil {
			return time.Time{}, fmt.Errorf("Invalid date '%s'. Use YYYY-MM-DD or: today, tomorrow, end-of-week, end-of-month, next-week, next-month", input)
		}
		return t, nil
	}
}
