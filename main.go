package main

import (
	"processor/receipt-processor-challenge/controller"
	"processor/receipt-processor-challenge/repository"

	"github.com/gin-gonic/gin"
)

func main() {
	pointRepository := repository.GetPointRepository()
	router := gin.Default()
	router.GET("/receipts/:id/points", func(c *gin.Context) {
		controller.GetPointsById(c, pointRepository)
	})
	router.POST("/receipts/process", func(c *gin.Context) {
		controller.AddReceipt(c, pointRepository)
	})

	router.Run("localhost:8080")
}
