package appvalidator

import (
	"github.com/go-playground/validator/v10"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/util"
	"regexp"
	"strings"
)

func CommaSeparatedValidation(fl validator.FieldLevel) bool {
	valuesInStr := fl.Field().Interface().(string)
	valuesInStr = strings.TrimSpace(valuesInStr)

	param := fl.Param()

	if util.IsEmptyString(valuesInStr) {
		return false
	}

	valuesWithoutComma := strings.ReplaceAll(valuesInStr, ",", "")
	if len(valuesWithoutComma) == 0 {
		return false
	}

	if !util.IsEmptyString(param) && param == appconstant.ParamValidationCommaSeparatedForNumber {
		pattern := regexp.MustCompile(`^\d+(,\d+)*$`)
		if !pattern.MatchString(valuesInStr) {
			return false
		}
	}

	return true
}
