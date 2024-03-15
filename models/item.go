package models

import (
	"processor/receipt-processor-challenge/utils"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Item struct {
	ShortDescription string `json:"shortDescription" binding:"required,shortDescriptionValidator"`
	Price            string `json:"price" binding:"required,priceValidator"`
}

func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("shortDescriptionValidator", utils.ValidateDescription)
		_ = v.RegisterValidation("priceValidator", utils.ValidatePrice)
	}
}
