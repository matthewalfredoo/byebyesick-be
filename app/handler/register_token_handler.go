package handler

import (
	"github.com/gin-gonic/gin"
	"halodeksik-be/app/appvalidator"
	"halodeksik-be/app/dto"
	"halodeksik-be/app/dto/requestdto"
	"halodeksik-be/app/usecase"
	"net/http"
)

type RegisterTokenHandler struct {
	ucToken   usecase.RegisterTokenUseCase
	validator appvalidator.AppValidator
}

func NewRegisterTokenHandler(uc usecase.RegisterTokenUseCase, v appvalidator.AppValidator) *RegisterTokenHandler {
	return &RegisterTokenHandler{ucToken: uc, validator: v}
}

func (h *RegisterTokenHandler) SendRegisterToken(ctx *gin.Context) {
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

	_, err = h.ucToken.SendRegisterToken(ctx.Request.Context(), req.Email)
	if err != nil {
		return
	}
	resp := dto.ResponseDto{Data: "Verification link has been sent to email."}
	ctx.JSON(http.StatusOK, resp)

}

func (h *RegisterTokenHandler) VerifyRegisterToken(ctx *gin.Context) {
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

	token, err := h.ucToken.VerifyRegisterToken(ctx.Request.Context(), req.Token)
	if err != nil {
		return
	}
	resp := dto.ResponseDto{Data: token.Email}
	ctx.JSON(http.StatusOK, resp)

}
