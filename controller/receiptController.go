package controller

import (
	"net/http"
	"processor/receipt-processor-challenge/models"
	"processor/receipt-processor-challenge/repository"
	"processor/receipt-processor-challenge/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AddReceipt(c *gin.Context, pointRepository *repository.PointRepository) {
	var receipt models.Receipt

	// Bind the JSON data from the request body to the Receipt struct
	if err := c.ShouldBindJSON(&receipt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate a random UUID
	id := uuid.New()

	// Call receiptProcessor to get points for a receipt
	points, err := service.ProcessReceipt(&receipt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Store id -> point pair
	pointRepository.Put(id.String(), points)

	c.IndentedJSON(http.StatusCreated, gin.H{"id": id.String()})
}

func GetPointsById(c *gin.Context, pointRepository *repository.PointRepository) {
	id := c.Param("id")
	points, exists := repository.GetPointRepository().Get(id)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "ID not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"points": points})
}
