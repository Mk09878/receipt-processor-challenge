package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCountAlphanumericChars(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"Target", 6},
		{"  !@#$%^&*() ", 0},
		{"M&M Corner Market", 14},
		{"abcABC123", 9},
	}

	for _, test := range tests {
		result := CountAlphanumericChars(test.input)
		assert.Equal(t, test.expected, result, "CountAlphanumericChars(%s) returned unexpected result", test.input)
	}
}

func TestCreateTimeFromDateAndTime(t *testing.T) {
	tests := []struct {
		name          string
		date          string
		time          string
		expected      time.Time
		expectedError bool
	}{
		{"ValidDateTime", "2024-03-14", "15:04", time.Date(2024, 3, 14, 15, 4, 0, 0, time.UTC), false},
		{"InvalidTimeFormat", "2024-03-14", "08:30:00", time.Time{}, true}, // incorrect time format
		{"InvalidDateFormat", "2024/03/14", "15:04", time.Time{}, true},    // incorrect date format
		{"InvalidHour", "2024-03-14", "25:00", time.Time{}, true},          // invalid hour
		{"InvalidMinute", "2024-03-14", "08:60", time.Time{}, true},        // invalid minute
		{"InvalidSecond", "2024-03-14", "08:30:61", time.Time{}, true},     // invalid second
		{"InvalidDateFebruary", "2024-02-30", "08:30", time.Time{}, true},  // invalid date (February 30th)
		{"InvalidDateMarch", "2024-03-32", "08:30", time.Time{}, true},     // invalid date (April 31st)
		{"ExtraCharacters", "2024-03-14", "08:30:00", time.Time{}, true},   // extra characters
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := CreateTimeFromDateAndTime(test.date, test.time)
			if test.expectedError {
				assert.Error(t, err, "Expected error for test case: %s", test.name)
			} else {
				assert.NoError(t, err, "Unexpected error for test case: %s", test.name)
				assert.Equal(t, test.expected, result, "Unexpected result for test case: %s", test.name)
			}
		})
	}
}
