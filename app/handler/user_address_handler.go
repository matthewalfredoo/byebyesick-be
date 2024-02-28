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

type UserAddressHandler struct {
	uc        usecase.AddressUseCase
	validator appvalidator.AppValidator
}

func NewAddressHandler(uc usecase.AddressUseCase, validator appvalidator.AppValidator) *UserAddressHandler {
	return &UserAddressHandler{uc: uc, validator: validator}
}

func (h *UserAddressHandler) GetAll(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()

	getAllProductQuery := queryparamdto.GetAllAddressesQuery{}
	err = ctx.ShouldBindQuery(&getAllProductQuery)
	if err != nil {
		return
	}

	err = h.validator.Validate(getAllProductQuery)
	if err != nil {
		return
	}

	param, err := getAllProductQuery.ToGetAllParams()
	if err != nil {
		return
	}

	paginatedItems, err := h.uc.GetAll(ctx.Request.Context(), param)
	if err != nil {
		return
	}

	resps := make([]*responsedto.AddressResponse, 0)
	for _, address := range paginatedItems.Items.([]*entity.Address) {
		resps = append(resps, address.ToAddressResponse())
	}
	paginatedItems.Items = resps

	resp := dto.ResponseDto{Data: paginatedItems}
	ctx.JSON(http.StatusOK, resp)
}

func (h *UserAddressHandler) GetMain(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()

	address, err := h.uc.GetMain(ctx.Request.Context())
	if err != nil {
		return
	}

	resp := dto.ResponseDto{Data: address.ToAddressResponse()}
	ctx.JSON(http.StatusOK, resp)
}

func (h *UserAddressHandler) SetMain(ctx *gin.Context) {
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

	added, err := h.uc.SetMain(ctx.Request.Context(), uri.Id)
	if err != nil {
		return
	}

	resp := dto.ResponseDto{Data: added.ToAddressResponse()}
	ctx.JSON(http.StatusOK, resp)

}

func (h *UserAddressHandler) Add(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()

	req := requestdto.RequestAddress{}
	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		return
	}

	err = h.validator.Validate(req)
	if err != nil {
		return
	}

	added, err := h.uc.Add(ctx.Request.Context(), req.ToAddress())
	if err != nil {
		return
	}

	resp := dto.ResponseDto{Data: added.ToAddressResponse()}
	ctx.JSON(http.StatusOK, resp)

}

func (h *UserAddressHandler) Update(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()

	req := requestdto.RequestAddress{}
	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		return
	}

	err = h.validator.Validate(req)
	if err != nil {
		return
	}

	uri := uriparamdto.ResourceById{}
	err = ctx.ShouldBindUri(&uri)
	if err != nil {
		return
	}

	err = h.validator.Validate(uri)
	if err != nil {
		return
	}

	address := req.ToAddress()
	address.Id = uri.Id

	added, err := h.uc.Edit(ctx.Request.Context(), address)
	if err != nil {
		return
	}

	resp := dto.ResponseDto{Data: added.ToAddressResponse()}
	ctx.JSON(http.StatusOK, resp)

}

func (h *UserAddressHandler) GetById(ctx *gin.Context) {
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

	address, err := h.uc.GetById(ctx.Request.Context(), uri.Id)
	if err != nil {
		return
	}

	resp := dto.ResponseDto{Data: address.ToAddressResponse()}
	ctx.JSON(http.StatusOK, resp)
}

func (h *UserAddressHandler) Remove(ctx *gin.Context) {
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
