package handler

import (
	"context"
	"errors"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/appvalidator"
	"halodeksik-be/app/dto"
	"halodeksik-be/app/dto/requestdto"
	"halodeksik-be/app/dto/responsedto"
	"halodeksik-be/app/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	ucAuth    usecase.AuthUsecase
	validator appvalidator.AppValidator
}

func NewAuthHandler(uc usecase.AuthUsecase, v appvalidator.AppValidator) *AuthHandler {
	return &AuthHandler{ucAuth: uc, validator: v}
}

func (h *AuthHandler) ResetPassword(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()

	reqUri := requestdto.RequestTokenUrl{}
	err = ctx.ShouldBindQuery(&reqUri)
	if err != nil {
		return
	}

	err = h.validator.Validate(reqUri)
	if err != nil {
		return
	}

	var req requestdto.ResetPasswordRequest
	if err = ctx.ShouldBindJSON(&req); err != nil {
		return
	}

	err = h.validator.Validate(req)
	if err != nil {
		return
	}

	_, err = h.ucAuth.ChangePassword(ctx, req.Password, reqUri.Token)
	if err != nil {
		return
	}

	resp := dto.ResponseDto{Data: "Password has been changed."}
	ctx.JSON(http.StatusOK, resp)

}

func (h *AuthHandler) Register(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()

	req := requestdto.RequestRegisterUser{}
	err = ctx.ShouldBind(&req)
	if err != nil {
		return
	}

	err = h.validator.Validate(req)
	if err != nil {
		return
	}

	reqUri := requestdto.RequestTokenUrl{}
	err = ctx.ShouldBindQuery(&reqUri)
	if err != nil {
		return
	}

	err = h.validator.Validate(reqUri)
	if err != nil {
		return
	}

	fileHeader, err := ctx.FormFile(appconstant.FormCertificate)
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		return
	}

	if fileHeader != nil && req.UserRoleId == appconstant.UserRoleIdDoctor {
		reqFile := requestdto.RequestRegisterDoctorCertificate{}
		err = ctx.ShouldBind(&reqFile)
		if err != nil {
			return
		}

		err = h.validator.Validate(reqFile)
		if err != nil {
			return
		}

		reqCtx1 := ctx.Request.Context()
		reqCtx2 := context.WithValue(reqCtx1, appconstant.FormCertificate, fileHeader)
		ctx.Request = ctx.Request.WithContext(reqCtx2)
	}

	user, err := h.ucAuth.Register(ctx.Request.Context(), req.ToUser(), reqUri.Token, req.Name)
	if err != nil {
		return
	}
	resp := dto.ResponseDto{Data: user.ToUserResponse()}
	ctx.JSON(http.StatusOK, resp)

}

func (h *AuthHandler) Login(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()

	var req requestdto.LoginRequest
	if err = ctx.ShouldBindJSON(&req); err != nil {
		return
	}

	err = h.validator.Validate(req)
	if err != nil {
		return
	}

	user, profile, err := h.ucAuth.Login(ctx.Request.Context(), req)
	if errors.Is(err, apperror.ErrRecordNotFound) {
		err = apperror.ErrWrongCredentials
		return
	}

	if err != nil {
		return
	}

	resp := dto.ResponseDto{Data: responsedto.LoginResponse{
		UserId:     user.Id,
		Email:      user.Email,
		UserRoleId: user.UserRoleId,
		Name:       profile.Name,
		Image:      profile.Image,
		Token:      profile.Token,
	}}
	ctx.JSON(http.StatusOK, resp)
}
