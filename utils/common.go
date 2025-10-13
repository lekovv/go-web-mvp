package utils

import (
	"errors"
	"time"
)

func ParseDate(date string) (time.Time, error) {
	result, err := time.Parse("2006-01-02", date)
	if err != nil {
		return time.Time{}, errors.New("invalid date format, expected YYYY-MM-DD")
	}

	return result, nil
}
