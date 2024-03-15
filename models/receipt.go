package models

import (
	"processor/receipt-processor-challenge/utils"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Receipt struct {
	Retailer     string `json:"retailer" binding:"required,retailerValidator"`
	PurchaseDate string `json:"purchaseDate" binding:"required,purchaseDateValidator"`
	PurchaseTime string `json:"purchaseTime" binding:"required,purchaseTimeValidator"`
	Items        []Item `json:"items" binding:"required"`
	Total        string `json:"total" binding:"required,totalValidator"`
}

func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("retailerValidator", utils.ValidateRetailer)
		_ = v.RegisterValidation("purchaseDateValidator", utils.ValidatePurchaseDate)
		_ = v.RegisterValidation("purchaseTimeValidator", utils.ValidatePurchaseTime)
		_ = v.RegisterValidation("totalValidator", utils.ValidatePrice)
	}
}
