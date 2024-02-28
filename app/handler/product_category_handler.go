package handler

import (
	"github.com/gin-gonic/gin"
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

type ProductCategoryHandler struct {
	uc        usecase.ProductCategoryUseCase
	validator appvalidator.AppValidator
}

func NewProductCategoryHandler(uc usecase.ProductCategoryUseCase, validator appvalidator.AppValidator) *ProductCategoryHandler {
	return &ProductCategoryHandler{uc: uc, validator: validator}
}

func (h *ProductCategoryHandler) Add(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()
	req := requestdto.AddEditProductCategory{}
	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		return
	}

	err = h.validator.Validate(req)
	if err != nil {
		return
	}

	added, err := h.uc.Add(ctx.Request.Context(), req.ToProductCategory())
	if err != nil {
		return
	}
	resp := dto.ResponseDto{Data: added.ToResponse()}
	ctx.JSON(http.StatusOK, resp)
}

func (h *ProductCategoryHandler) GetById(ctx *gin.Context) {
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

	category, err := h.uc.GetById(ctx.Request.Context(), uri.Id)
	if err != nil {
		return
	}
	resp := dto.ResponseDto{Data: category.ToResponse()}
	ctx.JSON(http.StatusOK, resp)
}

func (h *ProductCategoryHandler) GetAllWithoutParams(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()
	paginatedItems, err := h.uc.GetAllProductCategoriesWithoutParams(ctx.Request.Context())
	if err != nil {
		return
	}

	resps := make([]*responsedto.ProductCategoryResponse, 0)
	for _, drugClassification := range paginatedItems.Items.([]*entity.ProductCategory) {
		resps = append(resps, drugClassification.ToResponse())
	}
	paginatedItems.Items = resps

	resp := dto.ResponseDto{Data: paginatedItems}
	ctx.JSON(http.StatusOK, resp)
}

func (h *ProductCategoryHandler) GetAll(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()

	getAllProductCategoriesQuery := queryparamdto.GetAllProductCategoriesQuery{}
	_ = ctx.ShouldBindQuery(&getAllProductCategoriesQuery)

	param, err := getAllProductCategoriesQuery.ToGetAllParams()
	if err != nil {
		return
	}

	paginatedItems, err := h.uc.GetAll(ctx.Request.Context(), param)
	if err != nil {
		return
	}

	resps := make([]*responsedto.ProductCategoryResponse, 0)
	for _, category := range paginatedItems.Items.([]*entity.ProductCategory) {
		resps = append(resps, category.ToResponse())
	}
	paginatedItems.Items = resps

	resp := dto.ResponseDto{Data: paginatedItems}
	ctx.JSON(http.StatusOK, resp)
}

func (h *ProductCategoryHandler) Edit(ctx *gin.Context) {
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

	req := requestdto.AddEditProductCategory{}
	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		return
	}

	err = h.validator.Validate(req)
	if err != nil {
		return
	}

	edited, err := h.uc.Edit(ctx.Request.Context(), uri.Id, req.ToProductCategory())
	if err != nil {
		return
	}
	resp := dto.ResponseDto{Data: edited.ToResponse()}
	ctx.JSON(http.StatusOK, resp)
}

func (h *ProductCategoryHandler) Remove(ctx *gin.Context) {
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
