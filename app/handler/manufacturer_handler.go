package handler

import (
	"context"
	"errors"
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

type ManufacturerHandler struct {
	uc        usecase.ManufacturerUseCase
	validator appvalidator.AppValidator
}

func NewManufacturerHandler(uc usecase.ManufacturerUseCase, validator appvalidator.AppValidator) *ManufacturerHandler {
	return &ManufacturerHandler{uc: uc, validator: validator}
}

func (h *ManufacturerHandler) Add(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()

	req := requestdto.AddEditManufacturer{}
	err = ctx.Bind(&req)
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
		reqImage := requestdto.AddEditManufacturerImage{}
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

	added, err := h.uc.Add(ctx.Request.Context(), req.ToManufacturer())
	if err != nil {
		return
	}
	resp := dto.ResponseDto{Data: added.ToResponse()}
	ctx.JSON(http.StatusOK, resp)
}

func (h *ManufacturerHandler) GetById(ctx *gin.Context) {
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

	manufacturer, err := h.uc.GetById(ctx.Request.Context(), uri.Id)
	if err != nil {
		return
	}
	resp := dto.ResponseDto{Data: manufacturer.ToResponse()}
	ctx.JSON(http.StatusOK, resp)
}

func (h *ManufacturerHandler) GetAll(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()

	getAllManufacturersQuery := queryparamdto.GetAllManufacturersQuery{}
	err = ctx.ShouldBindQuery(&getAllManufacturersQuery)
	if err != nil {
		return
	}

	param := getAllManufacturersQuery.ToGetAllParams()
	paginatedItems, err := h.uc.GetAll(ctx.Request.Context(), param)
	if err != nil {
		return
	}

	resps := make([]*responsedto.ManufacturerResponse, 0)
	for _, manufacturer := range paginatedItems.Items.([]*entity.Manufacturer) {
		resps = append(resps, manufacturer.ToResponse())
	}
	paginatedItems.Items = resps

	resp := dto.ResponseDto{Data: paginatedItems}
	ctx.JSON(http.StatusOK, resp)
}

func (h *ManufacturerHandler) GetAllWithoutParams(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()
	paginatedItems, err := h.uc.GetAllManufacturersWithoutParams(ctx.Request.Context())
	if err != nil {
		return
	}

	resps := make([]*responsedto.ManufacturerResponse, 0)
	for _, manufacturer := range paginatedItems.Items.([]*entity.Manufacturer) {
		resps = append(resps, manufacturer.ToResponse())
	}
	paginatedItems.Items = resps

	resp := dto.ResponseDto{Data: paginatedItems}
	ctx.JSON(http.StatusOK, resp)
}

func (h *ManufacturerHandler) Edit(ctx *gin.Context) {
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

	req := requestdto.AddEditManufacturer{}
	err = ctx.Bind(&req)
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
		reqImage := requestdto.AddEditManufacturerImage{}
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

	edited, err := h.uc.Edit(ctx.Request.Context(), uri.Id, req.ToManufacturer())
	if err != nil {
		return
	}
	resp := dto.ResponseDto{Data: edited.ToResponse()}
	ctx.JSON(http.StatusOK, resp)
}

func (h *ManufacturerHandler) Remove(ctx *gin.Context) {
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

	err = h.uc.Remove(ctx, uri.Id)
	if err != nil {
		return
	}
	ctx.JSON(http.StatusNoContent, dto.ResponseDto{})
}
