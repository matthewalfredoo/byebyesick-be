package handler

import (
	"github.com/gin-gonic/gin"
	"halodeksik-be/app/appvalidator"
	"halodeksik-be/app/dto"
	"halodeksik-be/app/dto/requestdto"
	"halodeksik-be/app/usecase"
	"net/http"
)

type ForgotTokenHandler struct {
	ucToken   usecase.ForgotTokenUseCase
	validator appvalidator.AppValidator
}

func NewForgotTokenHandler(uc usecase.ForgotTokenUseCase, v appvalidator.AppValidator) *ForgotTokenHandler {
	return &ForgotTokenHandler{ucToken: uc, validator: v}
}

func (h *ForgotTokenHandler) SendForgotToken(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()

	req := requestdto.RequestToken{}
	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		return
	}

	err = h.validator.Validate(req)
	if err != nil {
		return
	}

	_, err = h.ucToken.SendForgotToken(ctx.Request.Context(), req.Email)
	if err != nil {
		return
	}
	resp := dto.ResponseDto{Data: "Forgot token link has been sent to email."}
	ctx.JSON(http.StatusOK, resp)

}

func (h *ForgotTokenHandler) VerifyForgotToken(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()

	req := requestdto.RequestTokenUrl{}
	err = ctx.ShouldBindQuery(&req)
	if err != nil {
		return
	}

	err = h.validator.Validate(req)
	if err != nil {
		return
	}

	_, err = h.ucToken.VerifyForgetToken(ctx.Request.Context(), req.Token)
	if err != nil {
		return
	}
	resp := dto.ResponseDto{Data: "Token is valid."}
	ctx.JSON(http.StatusOK, resp)

}
