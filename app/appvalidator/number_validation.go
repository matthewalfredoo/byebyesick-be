package appvalidator

import (
	"github.com/go-playground/validator/v10"
	"strconv"
)

func StringNumberGreaterThanValidation(fl validator.FieldLevel) bool {
	str := fl.Field().Interface().(string)
	gt := fl.Param()

	number, _ := strconv.Atoi(str)
	numbergt, err := strconv.Atoi(gt)
	if err != nil {
		panic("numbergt value must be a valid number")
	}

	return number > numbergt
}
