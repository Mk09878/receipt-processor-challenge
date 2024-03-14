package utils

import (
	"fmt"
	"time"
	"unicode"
)

func CountAlphanumericChars(str string) int {
	count := 0
	for _, char := range str {
		if unicode.IsLetter(char) || unicode.IsNumber(char) {
			count++
		}
	}
	return count
}

func CreateTimeFromDateAndTime(dateString string, timeString string) (time.Time, error) {
	// Combine date and time strings into a single string
	dateTimeStr := fmt.Sprintf("%s %s", dateString, timeString)

	// Parse combined date and time string into a time.Time object
	dateTime, err := time.Parse("2006-01-02 15:04", dateTimeStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("error parsing date and time: %v", err)
	}

	return dateTime, nil
}
