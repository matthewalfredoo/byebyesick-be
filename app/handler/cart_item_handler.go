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
	"halodeksik-be/app/entity"
	"halodeksik-be/app/usecase"
	"net/http"
)

type CartItemHandler struct {
	uc        usecase.CartItemUseCase
	validator appvalidator.AppValidator
}

func NewCartItemHandler(uc usecase.CartItemUseCase, validator appvalidator.AppValidator) *CartItemHandler {
	return &CartItemHandler{uc: uc, validator: validator}
}

func (h *CartItemHandler) Add(ctx *gin.Context) {
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

	req := requestdto.AddEditCartItem{}
	err = ctx.ShouldBind(&req)
	if err != nil {
		return
	}

	err = h.validator.Validate(req)
	if err != nil {
		return
	}

	_, err = h.uc.Add(ctx.Request.Context(), req.ToCartItem())
	if err != nil {
		return
	}

	paginatedItems, err := h.uc.GetAllByUserId(ctx.Request.Context())
	if err != nil {
		return
	}

	resps := make([]*responsedto.CartItemResponse, 0)
	for _, cartItem := range paginatedItems.Items.([]*entity.CartItem) {
		resps = append(resps, cartItem.ToResponse())
	}
	paginatedItems.Items = resps

	resp := dto.ResponseDto{Data: paginatedItems}
	ctx.JSON(http.StatusOK, resp)
}

func (h *CartItemHandler) GetAllByUserId(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()

	paginatedItems, err := h.uc.GetAllByUserId(ctx.Request.Context())
	if err != nil {
		return
	}

	resps := make([]*responsedto.CartItemResponse, 0)
	for _, cartItem := range paginatedItems.Items.([]*entity.CartItem) {
		resps = append(resps, cartItem.ToResponse())
	}
	paginatedItems.Items = resps

	resp := dto.ResponseDto{Data: paginatedItems}
	ctx.JSON(http.StatusOK, resp)
}

func (h *CartItemHandler) Remove(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()

	req := queryparamdto.DeleteCartItem{}
	err = ctx.ShouldBindQuery(&req)
	if err != nil {
		return
	}

	productIds, err := req.ToSliceOfInt64()
	if err != nil {
		return
	}

	err = h.validator.Validate(req)
	if err != nil {
		return
	}

	err = h.uc.Remove(ctx.Request.Context(), productIds)
	if err != nil {
		return
	}
	ctx.JSON(http.StatusNoContent, dto.ResponseDto{})
}

func (h *CartItemHandler) Checkout(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()

	getCartItemCheckoutQuery := queryparamdto.GetCartItemCheckoutQuery{}
	err = ctx.ShouldBindQuery(&getCartItemCheckoutQuery)
	if err != nil {
		return
	}

	err = h.validator.Validate(getCartItemCheckoutQuery)
	if err != nil {
		return
	}

	ids, err := getCartItemCheckoutQuery.GetCartItemIds()
	if err != nil {
		return
	}

	param := getCartItemCheckoutQuery.ToGetAllParams()
	paginatedItems, err := h.uc.Checkout(ctx.Request.Context(), param, ids...)
	if err != nil {
		return
	}

	resps := make([]*responsedto.CartItemResponse, 0)
	for _, cartItem := range paginatedItems.Items.([]*entity.CartItem) {
		resps = append(resps, cartItem.ToResponse())
	}
	paginatedItems.Items = resps

	ctx.JSON(http.StatusOK, paginatedItems)
}
