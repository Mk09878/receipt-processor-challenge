package utils

import (
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
)

func ValidatePrice(fl validator.FieldLevel) bool {
	return validateRegex(`^\d+\.\d{2}$`, fl.Field().String())
}

func ValidatePurchaseDate(fl validator.FieldLevel) bool {
	// Get the value of the PurchaseDate field
	dateStr := fl.Field().String()

	// Parse the date string into a time.Time object
	_, err := time.Parse("2006-01-02", dateStr)
	return err == nil
}

func ValidatePurchaseTime(fl validator.FieldLevel) bool {
	// Get the value of the PurchaseTime field
	timeStr := fl.Field().String()

	// Parse the time string into a time.Time object
	_, err := time.Parse("15:04", timeStr)
	return err == nil
}

func validateRegex(pattern string, input string) bool {
	// Compile the regex pattern
	regex := regexp.MustCompile(pattern)

	// Check if the value matches the regex pattern
	return regex.MatchString(input)
}
