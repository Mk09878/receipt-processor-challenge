package main

import (
	"log"
	"os"
	"processor/receipt-processor-challenge/controller"
	"processor/receipt-processor-challenge/middleware"
	"processor/receipt-processor-challenge/repository"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func getEnvVar(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}

func main() {
	pointRepository := repository.GetPointRepository()
	router := gin.Default()

	// Get PORT number from .env file
	port := getEnvVar("PORT")

	// Middleware
	router.Use(middleware.LoggingMiddleware)
	router.Use(middleware.CorsMiddleware())

	router.GET("/receipts/:id/points", func(c *gin.Context) {
		controller.GetPointsById(c, pointRepository)
	})
	router.POST("/receipts/process", func(c *gin.Context) {
		controller.AddReceipt(c, pointRepository)
	})

	router.Run(":" + port)
}
