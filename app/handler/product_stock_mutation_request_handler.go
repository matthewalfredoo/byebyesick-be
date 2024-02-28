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

type ProductStockMutationRequestHandler struct {
	uc        usecase.ProductStockMutationRequestUseCase
	validator appvalidator.AppValidator
}

func NewProductStockMutationRequestHandler(uc usecase.ProductStockMutationRequestUseCase, validator appvalidator.AppValidator) *ProductStockMutationRequestHandler {
	return &ProductStockMutationRequestHandler{uc: uc, validator: validator}
}

func (h *ProductStockMutationRequestHandler) Add(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()

	req := requestdto.AddProductStockMutationRequest{}
	if err = ctx.ShouldBindJSON(&req); err != nil {
		return
	}

	if err = h.validator.Validate(req); err != nil {
		return
	}

	added, err := h.uc.Add(ctx, req.ToProductStockMutationRequest())
	if err != nil {
		return
	}
	resp := dto.ResponseDto{Data: added.ToResponse()}
	ctx.JSON(http.StatusOK, resp)
}

func (h *ProductStockMutationRequestHandler) GetAllIncoming(ctx *gin.Context) {
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

	getAllMutationRequestQuery := queryparamdto.GetAllMutationRequestQuery{}
	err = ctx.ShouldBindQuery(&getAllMutationRequestQuery)
	if err != nil {
		return
	}

	err = h.validator.Validate(getAllMutationRequestQuery)
	if err != nil {
		return
	}

	incomingMutationReqQuery := queryparamdto.GetAllIncomingMutationRequestQuery{PharmacyOriginId: getAllMutationRequestQuery.PharmacyOriginId}
	if err = h.validator.Validate(incomingMutationReqQuery); err != nil {
		return
	}

	param, pharmacyOriginId, err := getAllMutationRequestQuery.ToGetAllParams(true)
	if err != nil {
		return
	}

	paginatedItems, err := h.uc.GetAllIncoming(ctx, pharmacyOriginId, param)
	if err != nil {
		return
	}

	resps := make([]*responsedto.ProductStockMutationRequestResponse, 0)
	for _, mutationRequest := range paginatedItems.Items.([]*entity.ProductStockMutationRequest) {
		mutationRequest.PharmacyProductOrigin = nil
		resps = append(resps, mutationRequest.ToResponse())
	}
	paginatedItems.Items = resps

	resp := dto.ResponseDto{Data: paginatedItems}
	ctx.JSON(http.StatusOK, resp)
}

func (h *ProductStockMutationRequestHandler) GetAllOutgoing(ctx *gin.Context) {
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

	getAllMutationRequestQuery := queryparamdto.GetAllMutationRequestQuery{}
	err = ctx.ShouldBindQuery(&getAllMutationRequestQuery)
	if err != nil {
		return
	}

	err = h.validator.Validate(getAllMutationRequestQuery)
	if err != nil {
		return
	}

	outgoingMutationReqQuery := queryparamdto.GetAllOutgoingMutationRequestQuery{PharmacyDestId: getAllMutationRequestQuery.PharmacyDestId}
	if err = h.validator.Validate(outgoingMutationReqQuery); err != nil {
		return
	}

	param, pharmacyDestId, err := getAllMutationRequestQuery.ToGetAllParams(false)
	if err != nil {
		return
	}

	paginatedItems, err := h.uc.GetAllOutgoing(ctx, pharmacyDestId, param)
	if err != nil {
		return
	}

	resps := make([]*responsedto.ProductStockMutationRequestResponse, 0)
	for _, mutationRequest := range paginatedItems.Items.([]*entity.ProductStockMutationRequest) {
		mutationRequest.PharmacyProductDest = nil
		resps = append(resps, mutationRequest.ToResponse())
	}
	paginatedItems.Items = resps

	resp := dto.ResponseDto{Data: paginatedItems}
	ctx.JSON(http.StatusOK, resp)
}

func (h *ProductStockMutationRequestHandler) EditStatus(ctx *gin.Context) {
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

	req := requestdto.EditProductStockMutationRequest{}
	if err = ctx.ShouldBindJSON(&req); err != nil {
		return
	}

	if err = h.validator.Validate(req); err != nil {
		return
	}

	updated, err := h.uc.EditStatus(ctx, uri.Id, req.ToProductStockMutationRequest())
	if err != nil {
		return
	}
	resp := dto.ResponseDto{Data: updated.ToResponse()}
	ctx.JSON(http.StatusOK, resp)
}
