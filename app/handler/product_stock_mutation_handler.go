package handler

import (
	"github.com/gin-gonic/gin"
	"halodeksik-be/app/appvalidator"
	"halodeksik-be/app/dto"
	"halodeksik-be/app/dto/requestdto"
	"halodeksik-be/app/usecase"
	"net/http"
)

type ProductStockMutationHandler struct {
	uc        usecase.ProductStockMutationUseCase
	validator appvalidator.AppValidator
}

func NewProductStockMutationHandler(uc usecase.ProductStockMutationUseCase, validator appvalidator.AppValidator) *ProductStockMutationHandler {
	return &ProductStockMutationHandler{uc: uc, validator: validator}
}

func (h *ProductStockMutationHandler) Add(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()

	req := requestdto.AddProductStockMutation{}
	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		return
	}

	err = h.validator.Validate(req)
	if err != nil {
		return
	}

	added, err := h.uc.Add(ctx.Request.Context(), req.ToProductStockMutation())
	if err != nil {
		return
	}
	resp := dto.ResponseDto{Data: added.ToResponse()}
	ctx.JSON(http.StatusOK, resp)
}
