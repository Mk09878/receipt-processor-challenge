package service

import (
	"log"
	"math"
	"processor/receipt-processor-challenge/models"
	"processor/receipt-processor-challenge/utils"
	"strconv"
	"strings"
	"time"
)

// ProcessReceipt calculates the points associated with the given receipt.
func ProcessReceipt(receipt *models.Receipt) (int, error) {
	receiptDateTime, err := utils.CreateTimeFromDateAndTime(receipt.PurchaseDate, receipt.PurchaseTime)
	if err != nil {
		log.Println("Error creating time:", err)
		return 0, err
	}

	points := 0

	// One point for every alphanumeric character in the retailer name
	points += utils.CountAlphanumericChars(receipt.Retailer)

	// Check if the total is a round dollar amount with no cents
	receiptTotal, err := strconv.ParseFloat(receipt.Total, 64)
	if err != nil {
		log.Println("Error parsing total:", err)
		return 0, err
	}

	points += calculatePointsForTotal(receiptTotal)

	points += calculatePointsForItems(receipt.Items)

	points += calculatePointsForOddDay(receiptDateTime)

	points += calculatePointsForTime(receiptDateTime)

	return points, nil
}

func calculatePointsForTotal(receiptTotal float64) int {
	points := 0
	totalCents := receiptTotal * 100

	if totalCents == 0 {
		return 0
	}

	// Check if the total is a rounded amount
	if math.Mod(totalCents, 100) == 0 {
		points += 50
	}

	// Check if the total is a multiple of 0.25
	if math.Mod(totalCents, 25) == 0 {
		points += 25
	}

	return points
}

func calculatePointsForItems(items []models.Item) int {
	points := 0

	// 5 points for every two items on the receipt
	points += (len(items) / 2) * 5

	for _, item := range items {
		// If the trimmed length of the item description is a multiple of 3,
		// multiply the price by 0.2 and round up to the nearest integer.
		// The result is the number of points earned.
		trimmedString := strings.TrimSpace(item.ShortDescription)
		trimmedLength := len(trimmedString)
		if trimmedLength%3 == 0 {
			price, err := strconv.ParseFloat(item.Price, 64)
			if err != nil {
				log.Println("Error parsing price:", err)
			}
			points += int(math.Ceil(price * 0.2))
		}
	}

	return points
}

// 6 points if the day in the purchase date is odd
func calculatePointsForOddDay(receiptDateTime time.Time) int {
	if receiptDateTime.Day()%2 != 0 {
		return 6
	} else {
		return 0
	}
}

// 10 points if the time of purchase is after 2:00pm and before 4:00pm
func calculatePointsForTime(receiptDateTime time.Time) int {
	if receiptDateTime.Hour() >= 14 && receiptDateTime.Minute() > 0 && receiptDateTime.Hour() < 16 {
		return 10
	} else {
		return 0
	}
}
