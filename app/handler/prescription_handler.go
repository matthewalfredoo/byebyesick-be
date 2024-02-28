package handler

import (
	"github.com/gin-gonic/gin"
	"halodeksik-be/app/appvalidator"
	"halodeksik-be/app/dto"
	"halodeksik-be/app/dto/requestdto"
	"halodeksik-be/app/dto/uriparamdto"
	"halodeksik-be/app/usecase"
	"net/http"
)

type PrescriptionHandler struct {
	uc        usecase.PrescriptionUseCase
	validator appvalidator.AppValidator
}

func NewPrescriptionHandler(uc usecase.PrescriptionUseCase, validator appvalidator.AppValidator) *PrescriptionHandler {
	return &PrescriptionHandler{uc: uc, validator: validator}
}

func (h *PrescriptionHandler) Add(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()

	req := requestdto.AddPrescription{}
	err = ctx.ShouldBindJSON(&req)

	err = h.validator.Validate(req)
	if err != nil {
		return
	}

	added, err := h.uc.Add(ctx, req.ToPrescription())
	if err != nil {
		return
	}

	resp := dto.ResponseDto{Data: added.ToResponse()}
	ctx.JSON(http.StatusOK, resp)
}

func (h *PrescriptionHandler) GetBySessionId(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()

	uri := uriparamdto.PrescriptionBySessionId{}
	err = ctx.ShouldBindUri(&uri)

	err = h.validator.Validate(uri)
	if err != nil {
		return
	}

	prescription, err := h.uc.GetBySessionId(ctx, uri.SessionId)
	if err != nil {
		return
	}

	resp := dto.ResponseDto{Data: prescription.ToResponse()}
	ctx.JSON(http.StatusOK, resp)
}

func (h *PrescriptionHandler) EditBySessionId(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()

	uri := uriparamdto.PrescriptionBySessionId{}
	err = ctx.ShouldBindUri(&uri)

	err = h.validator.Validate(uri)
	if err != nil {
		return
	}

	req := requestdto.EditPrescription{}
	err = ctx.ShouldBindJSON(&req)

	err = h.validator.Validate(req)
	if err != nil {
		return
	}

	edited, err := h.uc.EditBySessionId(ctx, uri.SessionId, req.ToPrescription())
	if err != nil {
		return
	}

	resp := dto.ResponseDto{Data: edited.ToResponse()}
	ctx.JSON(http.StatusOK, resp)
}
