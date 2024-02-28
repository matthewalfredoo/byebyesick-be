package appvalidator

import (
	"github.com/go-playground/validator/v10"
	"halodeksik-be/app/appconstant"
	"mime/multipart"
	"strconv"
)

func FileSizeValidation(fl validator.FieldLevel) bool {
	fileHeader := fl.Field().Interface().(multipart.FileHeader)
	maxFileSizeInKbStr := fl.Param()
	maxFileSizeInKb, err := strconv.ParseInt(maxFileSizeInKbStr, 10, 64)
	if err != nil {
		panic("filesize must be an integer in kb unit")
	}
	maxFileSizeInKb *= appconstant.BytesToKilobyte

	if fileHeader.Size > maxFileSizeInKb {
		return false
	}
	return true
}
