package handler

import (
	"github.com/gin-gonic/gin"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/appvalidator"
	"halodeksik-be/app/dto"
	"halodeksik-be/app/dto/queryparamdto"
	"halodeksik-be/app/dto/requestdto"
	"halodeksik-be/app/dto/responsedto"
	"halodeksik-be/app/dto/uriparamdto"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/usecase"
	"net/http"
)

type PharmacyHandler struct {
	uc        usecase.PharmacyUseCase
	validator appvalidator.AppValidator
}

func NewPharmacyHandler(uc usecase.PharmacyUseCase, validator appvalidator.AppValidator) *PharmacyHandler {
	return &PharmacyHandler{uc: uc, validator: validator}
}

func (h *PharmacyHandler) Add(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()

	req := requestdto.AddEditPharmacy{}
	req.PharmacyAdminId = ctx.Request.Context().Value(appconstant.ContextKeyUserId).(int64)
	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		return
	}

	err = h.validator.Validate(req)
	if err != nil {
		return
	}

	added, err := h.uc.Add(ctx.Request.Context(), req.ToPharmacy())
	if err != nil {
		return
	}
	resp := dto.ResponseDto{Data: added.ToPharmacyResponse()}
	ctx.JSON(http.StatusOK, resp)
}

func (h *PharmacyHandler) GetById(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()

	uri := uriparamdto.ResourceById{}
	err = ctx.ShouldBindUri(&uri)
	if err != nil {
		return
	}

	err = h.validator.Validate(uri)
	if err != nil {
		return
	}

	pharmacy, err := h.uc.GetById(ctx.Request.Context(), uri.Id)
	if err != nil {
		return
	}
	resp := dto.ResponseDto{Data: pharmacy.ToPharmacyResponse()}
	ctx.JSON(http.StatusOK, resp)
}

func (h *PharmacyHandler) GetAll(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()

	getAllPharmacyQuery := queryparamdto.GetAllPharmaciesQuery{}
	_ = ctx.ShouldBindQuery(&getAllPharmacyQuery)

	param, err := getAllPharmacyQuery.ToGetAllParams(ctx.Request.Context().Value(appconstant.ContextKeyUserId).(int64))
	if err != nil {
		return
	}

	paginatedItems, err := h.uc.GetAll(ctx.Request.Context(), param)
	if err != nil {
		return
	}

	resps := make([]*responsedto.PharmacyResponse, 0)
	for _, pharmacy := range paginatedItems.Items.([]*entity.Pharmacy) {
		resps = append(resps, pharmacy.ToPharmacyResponse())
	}
	paginatedItems.Items = resps

	resp := dto.ResponseDto{Data: paginatedItems}
	ctx.JSON(http.StatusOK, resp)
}

func (h *PharmacyHandler) Edit(ctx *gin.Context) {
	var err error

	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()

	uri := uriparamdto.ResourceById{}
	err = ctx.ShouldBindUri(&uri)
	if err != nil {
		return
	}

	err = h.validator.Validate(uri)
	if err != nil {
		return
	}

	req := requestdto.AddEditPharmacy{}
	req.PharmacyAdminId = ctx.Request.Context().Value(appconstant.ContextKeyUserId).(int64)
	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		return
	}

	err = h.validator.Validate(req)
	if err != nil {
		return
	}

	updated, err := h.uc.Edit(ctx.Request.Context(), uri.Id, req.ToPharmacy())
	if err != nil {
		return
	}
	resp := dto.ResponseDto{Data: updated.ToPharmacyResponse()}
	ctx.JSON(http.StatusOK, resp)
}

func (h *PharmacyHandler) Remove(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()

	uri := uriparamdto.ResourceById{}
	err = ctx.ShouldBindUri(&uri)
	if err != nil {
		return
	}

	err = h.validator.Validate(uri)
	if err != nil {
		return
	}

	err = h.uc.Remove(ctx.Request.Context(), uri.Id)
	if err != nil {
		return
	}
	ctx.JSON(http.StatusNoContent, dto.ResponseDto{})
}
