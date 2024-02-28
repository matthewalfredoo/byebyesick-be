package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"halodeksik-be/app/apperror"
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

type PharmacyProductHandler struct {
	uc        usecase.PharmacyProductUseCase
	validator appvalidator.AppValidator
}

func NewPharmacyProductHAndler(uc usecase.PharmacyProductUseCase, validator appvalidator.AppValidator) *PharmacyProductHandler {
	return &PharmacyProductHandler{uc: uc, validator: validator}
}

func (h *PharmacyProductHandler) Add(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			var notFoundError *apperror.NotFound
			if errors.As(err, &notFoundError) {
				err = WrapError(err, http.StatusBadRequest)
			} else {
				err = WrapError(err)
			}
			_ = ctx.Error(err)
		}
	}()

	req := requestdto.AddPharmacyProduct{}
	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		return
	}

	err = h.validator.Validate(req)
	if err != nil {
		return
	}

	added, err := h.uc.Add(ctx.Request.Context(), req.ToPharmacyProduct())
	if err != nil {
		return
	}
	resp := dto.ResponseDto{Data: added.ToPharmacyProductResponse()}
	ctx.JSON(http.StatusOK, resp)
}

func (h *PharmacyProductHandler) GetById(ctx *gin.Context) {
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
	resp := dto.ResponseDto{Data: pharmacy.ToPharmacyProductResponse()}
	ctx.JSON(http.StatusOK, resp)
}

func (h *PharmacyProductHandler) GetAllByProductId(ctx *gin.Context) {
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

	getAllPharmacyProductsByProductId := queryparamdto.GetAllPharmacyProductsByProductId{}
	err = ctx.ShouldBindQuery(&getAllPharmacyProductsByProductId)
	if err != nil {
		return
	}

	err = h.validator.Validate(getAllPharmacyProductsByProductId)
	if err != nil {
		return
	}

	param, err := getAllPharmacyProductsByProductId.ToGetAllParams()
	if err != nil {
		return
	}

	paginatedItems, err := h.uc.GetAllByProductId(ctx.Request.Context(), uri.Id, param)
	if err != nil {
		return
	}

	resps := make([]*responsedto.PharmacyProductResponse, 0)
	for _, pharmacyProduct := range paginatedItems.Items.([]*entity.PharmacyProduct) {
		resps = append(resps, pharmacyProduct.ToPharmacyProductResponse())
	}
	paginatedItems.Items = resps

	resp := dto.ResponseDto{Data: paginatedItems}
	ctx.JSON(http.StatusOK, resp)
}

func (h *PharmacyProductHandler) GetAllByPharmacy(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()

	getAllPharmacyProductQuery := queryparamdto.GetAllPharmacyProductsQuery{}
	_ = ctx.ShouldBindQuery(&getAllPharmacyProductQuery)

	param, err := getAllPharmacyProductQuery.ToGetAllParams()
	if err != nil {
		return
	}

	paginatedItems, err := h.uc.GetAllByPharmacy(ctx.Request.Context(), getAllPharmacyProductQuery.GetPharmacyId(), param)
	if err != nil {
		return
	}

	resps := make([]*responsedto.PharmacyProductResponse, 0)
	for _, pharmacyProduct := range paginatedItems.Items.([]*entity.PharmacyProduct) {
		resps = append(resps, pharmacyProduct.ToPharmacyProductResponse())
	}
	paginatedItems.Items = resps

	resp := dto.ResponseDto{Data: paginatedItems}
	ctx.JSON(http.StatusOK, resp)
}

func (h *PharmacyProductHandler) Edit(ctx *gin.Context) {
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

	req := requestdto.EditPharmacyProduct{}
	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		return
	}

	err = h.validator.Validate(req)
	if err != nil {
		return
	}

	updated, err := h.uc.Edit(ctx.Request.Context(), uri.Id, req.ToPharmacyProduct())
	if err != nil {
		return
	}
	resp := dto.ResponseDto{Data: updated.ToPharmacyProductResponse()}
	ctx.JSON(http.StatusOK, resp)
}
