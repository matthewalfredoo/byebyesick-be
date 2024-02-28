package handler

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/appvalidator"
	"halodeksik-be/app/dto"
	"halodeksik-be/app/dto/requestdto"
	"halodeksik-be/app/usecase"
	"net/http"
)

type ProfileHandler struct {
	uc        usecase.ProfileUseCase
	validator appvalidator.AppValidator
}

func NewProfileHandler(uc usecase.ProfileUseCase, v appvalidator.AppValidator) *ProfileHandler {
	return &ProfileHandler{uc: uc, validator: v}
}

func (h *ProfileHandler) GetProfile(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()

	roleId := ctx.Request.Context().Value(appconstant.ContextKeyRoleId)
	userId := ctx.Request.Context().Value(appconstant.ContextKeyUserId)

	resp := dto.ResponseDto{}
	if roleId.(int64) == appconstant.UserRoleIdDoctor {
		user, err := h.uc.GetDoctorProfileByUserId(ctx, userId.(int64))
		if err != nil {
			return
		}
		resp.Data = user.ToDoctorProfileResponse()
	} else if roleId.(int64) == appconstant.UserRoleIdUser {
		user, err := h.uc.GetUserProfileByUserId(ctx, userId.(int64))
		if err != nil {
			return
		}
		resp.Data = user.ToUserProfileResponse()
	} else {
		err = apperror.ErrUnauthorized
		return
	}

	ctx.JSON(http.StatusOK, resp)

}

func (h *ProfileHandler) EditDoctorProfile(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()

	roleId := ctx.Request.Context().Value(appconstant.ContextKeyRoleId)
	if roleId == nil || roleId.(int64) != appconstant.UserRoleIdDoctor {
		err = apperror.ErrUnauthorized
		return
	}

	if err = h.bindFile(ctx, appconstant.FormProfilePhoto); err != nil {
		return
	}

	if err = h.bindFile(ctx, appconstant.FormCertificate); err != nil {
		return
	}

	req := requestdto.RequestDoctorProfile{}
	err = ctx.ShouldBind(&req)
	if err != nil {
		return
	}

	err = h.validator.Validate(req)
	if err != nil {
		return
	}

	profile := req.ToDoctorProfile()

	doctorProfile, err := h.uc.UpdateDoctorProfile(ctx, profile)
	if err != nil {
		return
	}
	resp := dto.ResponseDto{Data: doctorProfile.ToDoctorProfileResponse()}
	ctx.JSON(http.StatusOK, resp)
}

func (h *ProfileHandler) EditDoctorIsOnline(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()

	req := requestdto.RequestDoctorIsOnline{}
	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		return
	}

	err = h.validator.Validate(req)
	if err != nil {
		return
	}

	doctorProfile, err := h.uc.UpdateDoctorIsOnline(ctx, *req.IsOnline)
	if err != nil {
		return
	}

	resp := dto.ResponseDto{Data: doctorProfile.ToDoctorProfileResponse()}
	ctx.JSON(http.StatusOK, resp)
}

func (h *ProfileHandler) EditUserProfile(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()

	roleId := ctx.Request.Context().Value(appconstant.ContextKeyRoleId)
	if roleId == nil || roleId.(int64) != appconstant.UserRoleIdUser {
		err = apperror.ErrUnauthorized
		return
	}

	if err = h.bindFile(ctx, appconstant.FormProfilePhoto); err != nil {
		return
	}

	req := requestdto.RequestUserProfile{}
	err = ctx.ShouldBind(&req)
	if err != nil {
		return
	}

	err = h.validator.Validate(req)
	if err != nil {
		return
	}

	profile, err := req.ToUserProfile()
	if err != nil {
		return
	}

	userProfile, err := h.uc.UpdateUserProfile(ctx, profile)
	if err != nil {
		return
	}
	resp := dto.ResponseDto{Data: userProfile.ToUserProfileResponse()}
	ctx.JSON(http.StatusOK, resp)
}

func (h *ProfileHandler) bindFile(ctx *gin.Context, fileFieldName string) error {
	file, err := ctx.FormFile(fileFieldName)
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		return err
	}

	if file != nil {
		req1 := requestdto.RequestProfilePhoto{}
		req2 := requestdto.RequestDoctorCertificate{}
		if fileFieldName == appconstant.FormProfilePhoto {
			err = ctx.ShouldBind(&req1)
			if err != nil {
				return err
			}
			err = h.validator.Validate(req1)
			if err != nil {
				return err
			}
		} else if fileFieldName == appconstant.FormCertificate {
			err = ctx.ShouldBind(&req2)
			if err != nil {
				return err
			}
			err = h.validator.Validate(req2)
			if err != nil {
				return err
			}
		} else {
			return apperror.ErrUnauthorized
		}

		reqCtx1 := ctx.Request.Context()
		reqCtx2 := context.WithValue(reqCtx1, fileFieldName, file)
		ctx.Request = ctx.Request.WithContext(reqCtx2)
	}

	return nil
}
