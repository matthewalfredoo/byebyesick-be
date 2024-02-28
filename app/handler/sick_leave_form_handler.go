package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/appvalidator"
	"halodeksik-be/app/dto"
	"halodeksik-be/app/dto/requestdto"
	"halodeksik-be/app/dto/uriparamdto"
	"halodeksik-be/app/usecase"
	"net/http"
)

type SickLeaveFormHandler struct {
	uc        usecase.SickLeaveFormUseCase
	validator appvalidator.AppValidator
}

func NewSickLeaveFormHandler(uc usecase.SickLeaveFormUseCase, validator appvalidator.AppValidator) *SickLeaveFormHandler {
	return &SickLeaveFormHandler{uc: uc, validator: validator}
}

func (h *SickLeaveFormHandler) Add(ctx *gin.Context) {
	var err error
	defer func() {
		var errNotFound *apperror.NotFound
		if err != nil {
			if errors.As(err, &errNotFound) {
				err = WrapError(err, http.StatusBadRequest)
				_ = ctx.Error(err)
				return
			}
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()

	req := requestdto.AddSickLeaveForm{}
	err = ctx.ShouldBind(&req)
	if err != nil {
		return
	}

	err = h.validator.Validate(req)
	if err != nil {
		return
	}

	toAdd, err := req.ToSickLeaveForm()
	if err != nil {
		return
	}

	added, err := h.uc.Add(ctx.Request.Context(), toAdd)
	if err != nil {
		return
	}

	resp := dto.ResponseDto{Data: added.ToResponse()}
	ctx.JSON(http.StatusOK, resp)
}

func (h *SickLeaveFormHandler) GetBySessionId(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()

	uri := uriparamdto.SickLeaveFormBySessionId{}
	err = ctx.ShouldBindUri(&uri)
	if err != nil {
		return
	}

	err = h.validator.Validate(uri)
	if err != nil {
		return
	}

	form, err := h.uc.GetBySessionId(ctx, uri.SessionId)
	if err != nil {
		return
	}

	resp := dto.ResponseDto{Data: form.ToResponse()}
	ctx.JSON(http.StatusOK, resp)
}

func (h *SickLeaveFormHandler) EditBySessionId(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()

	uri := uriparamdto.SickLeaveFormBySessionId{}
	err = ctx.ShouldBindUri(&uri)
	if err != nil {
		return
	}

	err = h.validator.Validate(uri)
	if err != nil {
		return
	}

	req := requestdto.EditSickLeaveForm{}
	err = ctx.ShouldBind(&req)
	if err != nil {
		return
	}

	err = h.validator.Validate(req)
	if err != nil {
		return
	}

	toEdit, err := req.ToSickLeaveForm()
	if err != nil {
		return
	}

	edited, err := h.uc.EditBySessionId(ctx, uri.SessionId, toEdit)
	if err != nil {
		return
	}

	resp := dto.ResponseDto{Data: edited.ToResponse()}
	ctx.JSON(http.StatusOK, resp)
}
