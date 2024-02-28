package appvalidator

import (
	"github.com/go-playground/validator/v10"
	"halodeksik-be/app/util"
)

func StringDateTimeValidation(fl validator.FieldLevel) bool {
	str := fl.Field().Interface().(string)
	format := fl.Param()

	_, err := util.ParseDateTime(str, format)
	if err != nil {
		return false
	}
	return true
}
