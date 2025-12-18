package utils

import (
	"errors"
	"strings"
	"time"
)

// ParseEntryDate
// This function will probably need a lot of tweaking in the future to handle
// different requirements and functionalities
func ParseEntryDate(input string) (string, error) {
	loc, err := time.LoadLocation("Asia/Singapore")
	if err != nil {
		loc = time.Local
	}

	now := time.Now().In(loc)
	input = strings.TrimSpace(strings.ToLower(input))

	if input == "" {
		return now.Format("2006-01-02"), nil
	}

	if input == "yesterday" {
		d := now.AddDate(0, 0, -1)
		return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location()).Format("2006-01-02"), nil
	}

	if t, err := time.Parse("2/1/2006", input); err == nil {
		return t.Format("2006-01-02"), nil
	}

	if t, err := time.Parse("2/1", input); err == nil {
		return t.AddDate(now.Year(), 0, 0).Format("2006-01-02"), nil
	}

	return time.Time{}.Format("2006-01-02"), errors.New("invalid date format: use 'dd/mm/yyyy', 'dd/mm', or 'Yesterday'")
}
