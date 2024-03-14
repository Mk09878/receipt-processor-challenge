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
func ProcessReceipt(receipt *models.Receipt) int {
	receiptDateTime, err := utils.CreateTimeFromDateAndTime(receipt.PurchaseDate, receipt.PurchaseTime)
	points := 0

	// One point for every alphanumeric character in the retailer name
	points += utils.CountAlphanumericChars(receipt.Retailer)

	log.Println("After alphanumeric count", points)

	// Check if the total is a round dollar amount with no cents
	receiptTotal, err := strconv.ParseFloat(receipt.Total, 64)
	if err != nil {
		// Handle error (e.g., invalid total format)
		// For now, let's just log the error and continue
		log.Println("Error parsing total:", err)
		// return or continue with default value
	}

	points += calculatePointsForTotal(receiptTotal)

	points += calculatePointsForItems(receipt.Items)

	points += calculatePointsForOddDay(receiptDateTime)

	log.Println("After day count", points)

	points += calculatePointsForTime(receiptDateTime)

	log.Println("After time count", points)

	return points
}

func calculatePointsForTotal(receiptTotal float64) int {
	points := 0
	totalCents := receiptTotal * 100

	// Check if the total is a rounded amount
	if math.Mod(totalCents, 100) == 0 {
		points += 50
	}

	log.Println("After rounded count", points)

	// Check if the total is a multiple of 0.25
	if math.Mod(totalCents, 25) == 0 {
		points += 25
	}

	log.Println("After 0.25 count", points)

	return points
}

func calculatePointsForItems(items []models.Item) int {
	points := 0

	// 5 points for every two items on the receipt
	points += (len(items) / 2) * 5

	log.Println("After item length", items)

	log.Println("After item length", points)

	// Check item descriptions and prices
	for _, item := range items {
		// If the trimmed length of the item description is a multiple of 3,
		// multiply the price by 0.2 and round up to the nearest integer.
		// The result is the number of points earned.
		log.Println("Original item name:", item.ShortDescription)
		trimmedString := strings.TrimSpace(item.ShortDescription)
		log.Println("Trimmed item name:", trimmedString)
		trimmedLength := len(trimmedString)
		log.Println("trimmedLength", trimmedLength)
		if trimmedLength%3 == 0 {
			price, err := strconv.ParseFloat(item.Price, 64)
			if err != nil {
				// Handle error (e.g., invalid price format)
				// For now, let's just log the error and continue
				log.Println("Error parsing price:", err)
				// return or continue with default value
			}
			log.Println("trimmed BS", int(math.Ceil(price*0.2)))
			points += int(math.Ceil(price * 0.2))
		}
	}

	log.Println("After item description", points)

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
