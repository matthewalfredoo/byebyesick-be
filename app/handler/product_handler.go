package handler

import (
	"context"
	"errors"
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

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	uc        usecase.ProductUseCase
	validator appvalidator.AppValidator
}

func NewProductHandler(uc usecase.ProductUseCase, validator appvalidator.AppValidator) *ProductHandler {
	return &ProductHandler{uc: uc, validator: validator}
}

func (h *ProductHandler) Add(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()

	req := requestdto.AddProduct{}
	err = ctx.ShouldBind(&req)
	if err != nil {
		return
	}

	err = h.validator.Validate(req)
	if err != nil {
		return
	}

	fileHeader, err := ctx.FormFile(appconstant.FormImage)
	if err != nil {
		return
	}

	reqCtx1 := ctx.Request.Context()
	reqCtx2 := context.WithValue(reqCtx1, appconstant.FormImage, fileHeader)
	ctx.Request = ctx.Request.WithContext(reqCtx2)

	added, err := h.uc.Add(ctx.Request.Context(), req.ToProduct())
	if err != nil {
		return
	}
	resp := dto.ResponseDto{Data: added.ToProductResponse()}
	ctx.JSON(http.StatusOK, resp)
}

func (h *ProductHandler) GetById(ctx *gin.Context) {
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

	product, err := h.uc.GetById(ctx.Request.Context(), uri.Id)
	if err != nil {
		return
	}
	resp := dto.ResponseDto{Data: product.ToProductResponse()}
	ctx.JSON(http.StatusOK, resp)
}

func (h *ProductHandler) GetByIdForUser(ctx *gin.Context) {
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

	getByIdProductQuery := queryparamdto.GetByIdProductQuery{}
	err = ctx.ShouldBindQuery(&getByIdProductQuery)
	if err != nil {
		return
	}

	err = h.validator.Validate(getByIdProductQuery)
	if err != nil {
		return
	}

	product, err := h.uc.GetByIdForUser(ctx.Request.Context(), uri.Id, getByIdProductQuery.ToGetAllParams())
	if err != nil {
		return
	}
	resp := dto.ResponseDto{Data: product.ToProductResponse()}
	ctx.JSON(http.StatusOK, resp)
}

func (h *ProductHandler) GetAll(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()

	getAllProductQuery := queryparamdto.GetAllProductsQuery{}
	err = ctx.ShouldBindQuery(&getAllProductQuery)
	if err != nil {
		return
	}

	err = h.validator.Validate(getAllProductQuery)
	if err != nil {
		return
	}

	param, lat, long, err := getAllProductQuery.ToGetAllParams()
	if err != nil {
		return
	}

	paginatedItems, err := h.uc.GetAllForUser(ctx.Request.Context(), lat, long, param)
	if err != nil {
		return
	}

	resps := make([]*responsedto.ProductResponse, 0)
	for _, product := range paginatedItems.Items.([]*entity.Product) {
		resps = append(resps, product.ToProductResponse())
	}
	paginatedItems.Items = resps

	resp := dto.ResponseDto{Data: paginatedItems}
	ctx.JSON(http.StatusOK, resp)
}

func (h *ProductHandler) GetAllForAdmin(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()

	getAllProductsForAdminQuery := queryparamdto.GetAllProductsForAdminQuery{}
	err = ctx.ShouldBindQuery(&getAllProductsForAdminQuery)
	if err != nil {
		return
	}

	err = h.validator.Validate(getAllProductsForAdminQuery)
	if err != nil {
		return
	}

	param, err := getAllProductsForAdminQuery.ToGetAllParams()
	if err != nil {
		return
	}

	paginatedItems, err := h.uc.GetAllForAdminByPharmacyId(ctx.Request.Context(), getAllProductsForAdminQuery.GetPharmacyId(), param)
	if err != nil {
		return
	}

	resps := make([]*responsedto.ProductResponse, 0)
	for _, product := range paginatedItems.Items.([]*entity.Product) {
		resps = append(resps, product.ToProductResponse())
	}
	paginatedItems.Items = resps

	resp := dto.ResponseDto{Data: paginatedItems}
	ctx.JSON(http.StatusOK, resp)
}

func (h *ProductHandler) Edit(ctx *gin.Context) {
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

	req := requestdto.EditProduct{}
	err = ctx.ShouldBind(&req)
	if err != nil {
		return
	}

	err = h.validator.Validate(req)
	if err != nil {
		return
	}

	fileHeader, err := ctx.FormFile(appconstant.FormImage)
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		return
	}
	if fileHeader != nil {
		reqImage := requestdto.EditProductImage{}
		err = ctx.ShouldBind(&reqImage)
		if err != nil {
			return
		}

		err = h.validator.Validate(reqImage)
		if err != nil {
			return
		}

		reqCtx1 := ctx.Request.Context()
		reqCtx2 := context.WithValue(reqCtx1, appconstant.FormImage, fileHeader)
		ctx.Request = ctx.Request.WithContext(reqCtx2)
	}

	updated, err := h.uc.Edit(ctx.Request.Context(), uri.Id, req.ToProduct())
	if err != nil {
		return
	}
	resp := dto.ResponseDto{Data: updated.ToProductResponse()}
	ctx.JSON(http.StatusOK, resp)
}

func (h *ProductHandler) Remove(ctx *gin.Context) {
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
