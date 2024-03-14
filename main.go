package main

import (
	"net/http"
	"processor/receipt-processor-challenge/models"
	"processor/receipt-processor-challenge/repository"
	"processor/receipt-processor-challenge/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var pointRepository *repository.PointRepository

func main() {
	router := gin.Default()
	pointRepository = repository.GetPointRepository()

	router.GET("/receipts/:id/points", getPointsById)
	router.POST("/receipts/process", addReceipt)

	router.Run("localhost:8080")
}

func addReceipt(c *gin.Context) {
	var receipt models.Receipt

	// Bind the JSON data from the request body to the Receipt struct
	if err := c.ShouldBindJSON(&receipt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate a random UUID
	id := uuid.New()

	// Call receiptProcessor to get points for a receipt
	points := service.ProcessReceipt(&receipt)

	// Store id -> point pair
	pointRepository.Put(id.String(), points)

	c.IndentedJSON(http.StatusCreated, id)
}

func getPointsById(c *gin.Context) {
	id := c.Param("id")
	points, exists := repository.GetPointRepository().Get(id)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "ID not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"points": points})
}
