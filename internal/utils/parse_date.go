package utils

import (
	"errors"
	"time"
)

func ParseDate(date string) (string, error) {
	result, err := time.Parse("20060102", date)
	if err != nil {
		return "", errors.New("invalid date format, expected YYYYMMDD")
	}
	return result.Format("20060102"), nil
}
