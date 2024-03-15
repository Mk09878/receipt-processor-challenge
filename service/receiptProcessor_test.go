package service

import (
	"errors"
	"testing"
	"time"

	"processor/receipt-processor-challenge/models"

	"github.com/stretchr/testify/assert"
)

func TestCalculatePointsForTotal(t *testing.T) {
	tests := []struct {
		name         string
		receiptTotal float64
		expected     int
	}{
		{"RoundedAmountAndMultipleOf0.25", 100.0, 75}, // Rounded amount, should add 50 points
		{"MultipleOf0.25", 50.25, 25},                 // Not rounded, but multiple of 0.25, should add 25 points
		{"NotRounded", 99.99, 0},                      // Not rounded, not a multiple of 0.25, should add 0 points
		{"NotRounded", 0.0, 0},                        // No points for 0 dollars
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			points := calculatePointsForTotal(test.receiptTotal)
			assert.Equal(t, test.expected, points, "For receipt total %.2f, expected %d points", test.receiptTotal, test.expected)
		})
	}
}

func TestCalculatePointsForOddDay(t *testing.T) {
	tests := []struct {
		name        string
		receiptDate time.Time
		expected    int
	}{
		{"OddDay", time.Date(2024, time.March, 15, 0, 0, 0, 0, time.UTC), 6},  // Odd day, should add 6 points
		{"EvenDay", time.Date(2024, time.March, 14, 0, 0, 0, 0, time.UTC), 0}, // Even day, should add 0 points
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			points := calculatePointsForOddDay(test.receiptDate)
			assert.Equal(t, test.expected, points, "For receipt date %s, expected %d points", test.receiptDate.String(), test.expected)
		})
	}
}

func TestCalculatePointsForTime(t *testing.T) {
	tests := []struct {
		name        string
		receiptTime time.Time
		expected    int
	}{
		{"After2PM_Before4PM", time.Date(2024, time.March, 14, 15, 30, 0, 0, time.UTC), 10}, // Purchase time between 2:00pm and 4:00pm, should add 10 points
		{"Before2PM", time.Date(2024, time.March, 14, 13, 0, 0, 0, time.UTC), 0},            // Purchase time before 2:00pm, should add 0 points
		{"After4PM", time.Date(2024, time.March, 14, 16, 30, 0, 0, time.UTC), 0},            // Purchase time after 4:00pm, should add 0 points
		{"Exact2PM", time.Date(2024, time.March, 14, 14, 0, 0, 0, time.UTC), 0},             // Purchase time exactly at 2:00pm, should add 0 points
		{"Exact4PM", time.Date(2024, time.March, 14, 16, 0, 0, 0, time.UTC), 0},             // Purchase time exactly at 4:00pm, should add 0 points
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			points := calculatePointsForTime(test.receiptTime)
			assert.Equal(t, test.expected, points, "For receipt time %s, expected %d points", test.receiptTime.String(), test.expected)
		})
	}
}

func TestCalculatePointsForItems(t *testing.T) {
	tests := []struct {
		name     string
		items    []models.Item
		expected int
	}{
		{
			name:     "NoItems",
			items:    []models.Item{},
			expected: 0,
		},
		{
			name: "TwoItems",
			items: []models.Item{
				{ShortDescription: "Mountain Dew 12PK", Price: "6.49"},
				{ShortDescription: "Emils Cheese Pizza", Price: "12.25"},
				{ShortDescription: "Knorr Creamy Chicken", Price: "1.26"},
				{ShortDescription: "Doritos Nacho Cheese", Price: "3.35"},
				{ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ", Price: "12.00"},
			},
			expected: 16,
		},
		{
			name: "PriceNotFloat",
			items: []models.Item{
				{ShortDescription: "Item 1", Price: "invalid"},
			},
			expected: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			points := calculatePointsForItems(test.items)
			assert.Equal(t, test.expected, points, "For items %v, expected %d points", test.items, test.expected)
		})
	}
}

func TestProcessReceipt(t *testing.T) {
	tests := []struct {
		name           string
		receipt        *models.Receipt
		expectedPoints int
		expectedError  error
	}{
		{
			name: "ValidReceipt1",
			receipt: &models.Receipt{
				Retailer:     "Target",
				PurchaseDate: "2022-01-01",
				PurchaseTime: "13:01",
				Items: []models.Item{
					{ShortDescription: "Mountain Dew 12PK", Price: "6.49"},
					{ShortDescription: "Emils Cheese Pizza", Price: "12.25"},
					{ShortDescription: "Knorr Creamy Chicken", Price: "1.26"},
					{ShortDescription: "Doritos Nacho Cheese", Price: "3.35"},
					{ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ", Price: "12.00"},
				},
				Total: "35.35",
			},
			expectedPoints: 28,
			expectedError:  nil,
		},
		{
			name: "ValidReceipt2",
			receipt: &models.Receipt{
				Retailer:     "M&M Corner Market",
				PurchaseDate: "2022-03-20",
				PurchaseTime: "14:33",
				Items: []models.Item{
					{ShortDescription: "Gatorade", Price: "2.25"},
					{ShortDescription: "Gatorade", Price: "2.25"},
					{ShortDescription: "Gatorade", Price: "2.25"},
					{ShortDescription: "Gatorade", Price: "2.25"},
				},
				Total: "9.00",
			},
			expectedPoints: 109,
			expectedError:  nil,
		},
		{
			name: "InvalidTimeFormat",
			receipt: &models.Receipt{
				Retailer:     "Retailer",
				PurchaseDate: "2022-01-01",
				PurchaseTime: "invalid",
				Items: []models.Item{
					{ShortDescription: "Item 1", Price: "10.00"},
				},
				Total: "10.00",
			},
			expectedPoints: 0,
			expectedError:  errors.New("error parsing date and time: parsing time \"2022-01-01 invalid\" as \"2006-01-02 15:04\": cannot parse \"invalid\" as \"15\""),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			points, err := ProcessReceipt(test.receipt)
			assert.Equal(t, test.expectedPoints, points, "For test case %s, expected %d points", test.name, test.expectedPoints)

			if test.expectedError == nil {
				assert.NoError(t, err, "For test case %s, expected no error", test.name)
			} else {
				assert.EqualError(t, err, test.expectedError.Error(), "For test case %s, expected error: %v", test.name, test.expectedError)
			}
		})
	}
}
