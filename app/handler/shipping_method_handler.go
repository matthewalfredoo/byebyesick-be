package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/appvalidator"
	"halodeksik-be/app/dto"
	"halodeksik-be/app/dto/requestdto"
	"halodeksik-be/app/dto/responsedto"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/usecase"
	"net/http"
)

type ShippingMethodHandler struct {
	uc        usecase.ShippingMethodUseCase
	validator appvalidator.AppValidator
}

func NewShippingMethodHandler(uc usecase.ShippingMethodUseCase, validator appvalidator.AppValidator) *ShippingMethodHandler {
	return &ShippingMethodHandler{uc: uc, validator: validator}
}

func (h *ShippingMethodHandler) GetAll(ctx *gin.Context) {
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

	req := requestdto.CalculateShippingMethod{}
	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		return
	}

	err = h.validator.Validate(req)
	if err != nil {
		return
	}

	checkoutItems := make([]entity.CheckoutItem, 0)
	for _, cItem := range req.CheckoutItems {
		checkoutItems = append(checkoutItems, cItem.ToCheckoutItem())
	}

	paginatedItems, err := h.uc.GetAll(ctx, req.AddressId, checkoutItems)
	if err != nil {
		return
	}

	resps := make([]*responsedto.ShippingMethodResponse, 0)
	for _, shipMethod := range paginatedItems.Items.([]*entity.ShippingMethod) {
		resps = append(resps, shipMethod.ToResponse())
	}
	paginatedItems.Items = resps

	resp := dto.ResponseDto{Data: paginatedItems}
	ctx.JSON(http.StatusOK, resp)
}
