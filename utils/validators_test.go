package utils

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

var priceItems = []struct {
	have string
	want bool
}{
	{"10.99", true},
	{"10", false},
	{"10.9a", false},
	{"10.123", false},
	{"-10.99", false},
}

func TestValidatePrice(t *testing.T) {
	validate := validator.New()
	validate.RegisterValidation("price", ValidatePrice)

	for _, item := range priceItems {
		err := validate.Var(item.have, "price")
		if item.want {
			assert.Nil(t, err)
		} else {
			assert.Error(t, err)
		}
	}
}

var purchaseDateItems = []struct {
	have string
	want bool
}{
	{"2024-03-16", true},
	{"2024-02-30", false},
	{"2024-03-32", false},
	{"2024-13-01", false},
	{"2024/03/16", false},
	{"notadate", false},
	{"", false},
}

func TestValidatePurchaseDate(t *testing.T) {
	validate := validator.New()
	validate.RegisterValidation("purchaseDate", ValidatePurchaseDate)

	for _, item := range purchaseDateItems {
		err := validate.Var(item.have, "purchaseDate")
		if item.want {
			assert.Nil(t, err, "Expected valid purchase date")
		} else {
			assert.Error(t, err, "Expected invalid purchase date")
		}
	}
}

var purchaseTimeItems = []struct {
	have string
	want bool
}{
	{"12:30", true},
	{"23:59", true},
	{"00:00", true},
	{"12:61", false},
	{"24:00", false},
	{"12-30", false},
	{"notatime", false},
	{"", false},
}

func TestValidatePurchaseTime(t *testing.T) {
	validate := validator.New()
	validate.RegisterValidation("purchaseTime", ValidatePurchaseTime)

	for _, item := range purchaseTimeItems {
		err := validate.Var(item.have, "purchaseTime")
		if item.want {
			assert.Nil(t, err)
		} else {
			assert.Error(t, err)
		}
	}
}
