package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"processor/receipt-processor-challenge/controller"
	"processor/receipt-processor-challenge/repository"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const receiptRequestBody1 = `
{
  "retailer": "Target",
  "purchaseDate": "2022-01-01",
  "purchaseTime": "13:01",
  "items": [
    {
      "shortDescription": "Mountain Dew 12PK",
      "price": "6.49"
    },{
      "shortDescription": "Emils Cheese Pizza",
      "price": "12.25"
    },{
      "shortDescription": "Knorr Creamy Chicken",
      "price": "1.26"
    },{
      "shortDescription": "Doritos Nacho Cheese",
      "price": "3.35"
    },{
      "shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ",
      "price": "12.00"
    }
  ],
  "total": "35.35"
}
`

const receiptRequestBody2 = `
{
	"retailer": "M&M Corner Market",
	"purchaseDate": "2022-03-20",
	"purchaseTime": "14:33",
	"items": [
	  {
		"shortDescription": "Gatorade",
		"price": "2.25"
	  },{
		"shortDescription": "Gatorade",
		"price": "2.25"
	  },{
		"shortDescription": "Gatorade",
		"price": "2.25"
	  },{
		"shortDescription": "Gatorade",
		"price": "2.25"
	  }
	],
	"total": "9.00"
  }
`

const receiptBadRequestBody = `
{
	"retailer": "M&M Corner Market",
	"purchaseDate": "2022-03-20",
	"purchaseTime": "14:33:00", // Bad purchase time
	"items": [
	  {
		"shortDescription": "Gatorade",
		"price": "2.25"
	  },{
		"shortDescription": "Gatorade",
		"price": "2.25"
	  },{
		"shortDescription": "Gatorade",
		"price": "2.25"
	  },{
		"shortDescription": "Gatorade",
		"price": "2.25"
	  }
	],
	"total": "9.00"
  }
`

func TestReceipt1(t *testing.T) {
	router := setupRouter()

	// Send a POST request to create a receipt
	postRecorder := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/receipts/process", bytes.NewBuffer([]byte(receiptRequestBody1)))
	router.ServeHTTP(postRecorder, req)

	// Check if the POST request was successful
	assert.Equal(t, http.StatusCreated, postRecorder.Code)

	// Extract the ID from the response body
	var responseBody map[string]interface{}
	err := json.Unmarshal(postRecorder.Body.Bytes(), &responseBody)
	if err != nil {
		t.Fatalf("failed to unmarshal response body: %v", err)
	}

	// Ensure that the response body contains the "id" field
	id, ok := responseBody["id"].(string)
	assert.True(t, ok, "response body does not contain 'id' field")

	// Ensure that the ID is not empty
	assert.NotEmpty(t, id)

	// Send a GET request to retrieve the receipt using the ID
	getRecorder := httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/receipts/"+id+"/points", nil)
	router.ServeHTTP(getRecorder, req)

	// Check if the GET request was successful
	assert.Equal(t, http.StatusOK, getRecorder.Code)

	// Validate the response body
	responseBody = nil
	err = json.Unmarshal(getRecorder.Body.Bytes(), &responseBody)
	if err != nil {
		t.Fatalf("failed to unmarshal response body: %v", err)
	}

	// Ensure that the response body contains the "points" field
	points, ok := responseBody["points"].(float64)
	assert.True(t, ok, "response body does not contain 'points' field")

	// Validate the points value
	assert.NotNil(t, points)
	assert.Equal(t, 28.0, points)
}

func TestReceipt2(t *testing.T) {
	router := setupRouter()

	// Send a POST request to create a receipt
	postRecorder := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/receipts/process", bytes.NewBuffer([]byte(receiptRequestBody2)))
	router.ServeHTTP(postRecorder, req)

	// Check if the POST request was successful
	assert.Equal(t, http.StatusCreated, postRecorder.Code)

	// Extract the ID from the response body
	var responseBody map[string]interface{}
	err := json.Unmarshal(postRecorder.Body.Bytes(), &responseBody)
	if err != nil {
		t.Fatalf("failed to unmarshal response body: %v", err)
	}

	// Ensure that the response body contains the "id" field
	id, ok := responseBody["id"].(string)
	assert.True(t, ok, "response body does not contain 'id' field")

	// Ensure that the ID is not empty
	assert.NotEmpty(t, id)

	// Send a GET request to retrieve the receipt using the ID
	getRecorder := httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/receipts/"+id+"/points", nil)
	router.ServeHTTP(getRecorder, req)

	// Check if the GET request was successful
	assert.Equal(t, http.StatusOK, getRecorder.Code)

	// Validate the response body
	responseBody = nil
	err = json.Unmarshal(getRecorder.Body.Bytes(), &responseBody)
	if err != nil {
		t.Fatalf("failed to unmarshal response body: %v", err)
	}

	// Ensure that the response body contains the "points" field
	points, ok := responseBody["points"].(float64)
	assert.True(t, ok, "response body does not contain 'points' field")

	// Validate the points value
	assert.NotNil(t, points)
	assert.Equal(t, 109.0, points)
}

func TestGetPointsById_NonExistantId(t *testing.T) {
	router := setupRouter()
	// Send a GET request to retrieve the receipt using the ID
	getRecorder := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/receipts/123/points", nil)
	router.ServeHTTP(getRecorder, req)

	// Check if the GET request returned not found code
	assert.Equal(t, http.StatusNotFound, getRecorder.Code)
}

func TestAddReceipt_BadRequest(t *testing.T) {
	router := setupRouter()

	// Send a POST request to create a receipt
	postRecorder := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/receipts/process", bytes.NewBuffer([]byte(receiptBadRequestBody)))
	router.ServeHTTP(postRecorder, req)

	// Check if the POST request returned bad request code
	assert.Equal(t, http.StatusBadRequest, postRecorder.Code)
}

func setupRouter() *gin.Engine {
	pointRepository := repository.GetPointRepository()
	router := gin.Default()

	router.GET("/receipts/:id/points", func(c *gin.Context) {
		controller.GetPointsById(c, pointRepository)
	})
	router.POST("/receipts/process", func(c *gin.Context) {
		controller.AddReceipt(c, pointRepository)
	})

	return router
}
