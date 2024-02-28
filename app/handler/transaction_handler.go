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

type TransactionHandler struct {
	uc        usecase.TransactionUseCase
	validator appvalidator.AppValidator
}

func NewTransactionHandler(uc usecase.TransactionUseCase, validator appvalidator.AppValidator) *TransactionHandler {
	return &TransactionHandler{uc: uc, validator: validator}
}

func (h *TransactionHandler) GetTransactionById(ctx *gin.Context) {
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

	transaction, err := h.uc.GetTransactionById(ctx.Request.Context(), uri.Id)
	if err != nil {
		return
	}

	resp := dto.ResponseDto{Data: transaction.ToTransactionJoinResponse()}
	ctx.JSON(http.StatusOK, resp)
}

func (h *TransactionHandler) GetAllUserTransactions(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()

	getAllProductQuery := queryparamdto.GetAllTransactionsQuery{}
	err = ctx.ShouldBindQuery(&getAllProductQuery)
	if err != nil {
		return
	}

	err = h.validator.Validate(getAllProductQuery)
	if err != nil {
		return
	}

	param := getAllProductQuery.ToGetAllParams()

	paginatedItems, err := h.uc.GetAllTransactions(ctx.Request.Context(), param)
	if err != nil {
		return
	}

	resps := make([]*responsedto.TransactionJoinResponse, 0)
	for _, transaction := range paginatedItems.Items.([]*entity.Transaction) {
		transResp := transaction.ToTransactionJoinResponse()
		resps = append(resps, &transResp)
	}
	paginatedItems.Items = resps

	resp := dto.ResponseDto{Data: paginatedItems}
	ctx.JSON(http.StatusOK, resp)

}

func (h *TransactionHandler) UploadPaymentProof(ctx *gin.Context) {
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

	file, err := ctx.FormFile(appconstant.FormPaymentProof)
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		return
	}

	if file != nil {
		req := requestdto.RequestPaymentProof{}
		err = ctx.ShouldBind(&req)
		if err != nil {
			return
		}

		err = h.validator.Validate(req)
		if err != nil {
			return
		}

		reqCtx1 := ctx.Request.Context()
		reqCtx2 := context.WithValue(reqCtx1, appconstant.FormPaymentProof, file)
		ctx.Request = ctx.Request.WithContext(reqCtx2)
	}

	transaction, err := h.uc.UploadTransactionPayment(ctx, uri.Id)
	if err != nil {
		return
	}
	resp := dto.ResponseDto{Data: transaction.ToTransactionResponse()}
	ctx.JSON(http.StatusOK, resp)
}

func (h *TransactionHandler) AcceptTransaction(ctx *gin.Context) {
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

	transaction, err := h.uc.UpdateTransactionStatus(ctx, uri.Id, true)
	if err != nil {
		return
	}

	resp := dto.ResponseDto{Data: transaction.ToTransactionResponse()}
	ctx.JSON(http.StatusOK, resp)
}

func (h *TransactionHandler) RejectTransaction(ctx *gin.Context) {
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

	transaction, err := h.uc.UpdateTransactionStatus(ctx, uri.Id, false)
	if err != nil {
		return
	}

	resp := dto.ResponseDto{Data: transaction.ToTransactionResponse()}
	ctx.JSON(http.StatusOK, resp)
}

func (h *TransactionHandler) CancelTransaction(ctx *gin.Context) {
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

	transaction, err := h.uc.CancelTransaction(ctx, uri.Id)
	if err != nil {
		return
	}

	resp := dto.ResponseDto{Data: transaction.ToTransactionResponse()}
	ctx.JSON(http.StatusOK, resp)
}

func (h *TransactionHandler) GetPayment(ctx *gin.Context) {
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
	payment, err := h.uc.FindTotalPaymentByTransactionId(ctx, uri.Id)
	if err != nil {
		return
	}
	resp := dto.ResponseDto{Data: payment}
	ctx.JSON(http.StatusOK, resp)
}

func (h *TransactionHandler) AddTransaction(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()

	req := requestdto.AddTransaction{}
	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		return
	}

	err = h.validator.Validate(req)
	if err != nil {
		return
	}

	transaction, err := h.uc.AddTransaction(ctx.Request.Context(), req)
	if err != nil {
		return
	}
	resp := dto.ResponseDto{Data: transaction.ToTransactionResponse()}
	ctx.JSON(http.StatusOK, resp)
}
