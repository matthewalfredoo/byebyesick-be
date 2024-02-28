package middleware

import (
	"errors"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/dto"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ErrorHandler(ctx *gin.Context) {
	ctx.Next()
	var (
		errValidation validator.ValidationErrors
		errWrapper    *apperror.Wrapper
	)

	err := ctx.Errors.Last()
	if err == nil {
		return
	}
	resp := dto.ResponseDto{}

	switch {
	case errors.As(err, &errWrapper):
		resp.Errors = strings.Split(errWrapper.ErrorStored.Error(), "\n")
		if errors.As(errWrapper.ErrorStored, &errValidation) {
			resp.Errors = strings.Split(errWrapper.Message, "\n")
		}
		ctx.AbortWithStatusJSON(errWrapper.Code, resp)
	default:
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, resp)
	}
}
