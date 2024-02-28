package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/appvalidator"
	"halodeksik-be/app/dto"
	"halodeksik-be/app/dto/queryparamdto"
	"halodeksik-be/app/dto/responsedto"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/usecase"
	"net/http"
)

type StockReportHandler struct {
	uc        usecase.ProductStockMutationUseCase
	validator appvalidator.AppValidator
}

func NewStockReportHandler(uc usecase.ProductStockMutationUseCase, validator appvalidator.AppValidator) *StockReportHandler {
	return &StockReportHandler{uc: uc, validator: validator}
}

func (h *StockReportHandler) FindAll(ctx *gin.Context) {
	var err error
	var notFound *apperror.NotFound
	defer func() {
		if err != nil {
			if errors.As(err, &notFound) {
				err = WrapError(err, http.StatusBadRequest)
			} else {
				err = WrapError(err)
			}
			_ = ctx.Error(err)
		}
	}()

	getAllStockMutationQuery := queryparamdto.GetAllStockMutationsQuery{}
	_ = ctx.ShouldBindQuery(&getAllStockMutationQuery)

	err = h.validator.Validate(getAllStockMutationQuery)
	if err != nil {
		return
	}

	pharmacyAdmin := ctx.Request.Context().Value(appconstant.ContextKeyUserId).(int64)
	param, pharmacyId, err := getAllStockMutationQuery.ToGetAllParams(pharmacyAdmin)
	if err != nil {
		return
	}

	paginatedItems, err := h.uc.GetAllByPharmacy(ctx.Request.Context(), pharmacyId, param)
	if err != nil {
		return
	}

	resps := make([]*responsedto.ProductStockMutationResponse, 0)
	for _, stockMutation := range paginatedItems.Items.([]*entity.ProductStockMutation) {
		resps = append(resps, stockMutation.ToResponse())
	}
	paginatedItems.Items = resps

	resp := dto.ResponseDto{Data: paginatedItems}
	ctx.JSON(http.StatusOK, resp)

}
