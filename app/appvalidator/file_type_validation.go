package appvalidator

import (
	"github.com/go-playground/validator/v10"
	"halodeksik-be/app/util"
	"mime/multipart"
	"net/http"
	"strings"
)

func FileTypeValidation(fl validator.FieldLevel) bool {
	fileHeader := fl.Field().Interface().(multipart.FileHeader)
	fileType := fl.Param()
	file, err := fileHeader.Open()
	if err != nil {
		return false
	}

	extension, err := GetFileContentType(&file)
	if err != nil {
		return false
	}

	if util.IsEmptyString(extension) {
		return false
	}
	if !strings.Contains(fileType, extension) {
		return false
	}
	return true
}

func GetFileContentType(file *multipart.File) (string, error) {
	buffer := make([]byte, 512)

	_, err := (*file).Read(buffer)
	if err != nil {
		return "", err
	}

	contentType := http.DetectContentType(buffer)
	return strings.Split(contentType, "/")[1], nil
}
