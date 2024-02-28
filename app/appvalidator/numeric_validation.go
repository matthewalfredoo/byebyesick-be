package appvalidator

import (
	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
)

func StringNumericGreaterThanValidation(fl validator.FieldLevel) bool {
	str := fl.Field().Interface().(string)
	gt := fl.Param()

	numeric, _ := decimal.NewFromString(str)
	numericgt, err := decimal.NewFromString(gt)
	if err != nil {
		panic("numericgt value must be a valid numeric")
	}

	if numeric.LessThanOrEqual(numericgt) {
		return false
	}
	return true
}
