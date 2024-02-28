package handler

import (
	"github.com/gin-gonic/gin"
	"halodeksik-be/app/appvalidator"
	"halodeksik-be/app/dto"
	"halodeksik-be/app/dto/queryparamdto"
	"halodeksik-be/app/dto/responsedto"
	"halodeksik-be/app/dto/uriparamdto"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/usecase"
	"net/http"
)

type OrderHandler struct {
	uc        usecase.OrderUseCase
	validator appvalidator.AppValidator
}

func NewOrderHandler(uc usecase.OrderUseCase, validator appvalidator.AppValidator) *OrderHandler {
	return &OrderHandler{uc: uc, validator: validator}
}

func (h *OrderHandler) GetAllPharmacyAdminOrders(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()

	getAllOrderQuery := queryparamdto.GetAllOrderAdminRequestQuery{}
	err = ctx.ShouldBindQuery(&getAllOrderQuery)
	if err != nil {
		return
	}

	err = h.validator.Validate(getAllOrderQuery)
	if err != nil {
		return
	}

	param := getAllOrderQuery.ToGetAllParams()
	if err != nil {
		return
	}

	paginatedItems, err := h.uc.GetAllOrdersByPharmacyAdminId(ctx.Request.Context(), param)
	if err != nil {
		return
	}

	resps := make([]*responsedto.OrderListResponse, 0)
	for _, order := range paginatedItems.Items.([]*entity.Order) {
		transResp := order.ToOrderListResponse()
		resps = append(resps, &transResp)
	}
	paginatedItems.Items = resps

	resp := dto.ResponseDto{Data: paginatedItems}
	ctx.JSON(http.StatusOK, resp)
}

func (h *OrderHandler) GetAllUserOrders(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()

	getAllOrderQuery := queryparamdto.GetAllOrderUserRequestQuery{}
	err = ctx.ShouldBindQuery(&getAllOrderQuery)
	if err != nil {
		return
	}

	err = h.validator.Validate(getAllOrderQuery)
	if err != nil {
		return
	}

	param := getAllOrderQuery.ToGetAllParams()
	if err != nil {
		return
	}

	paginatedItems, err := h.uc.GetAllOrdersByUserId(ctx.Request.Context(), param)
	if err != nil {
		return
	}

	resps := make([]*responsedto.OrderListResponse, 0)
	for _, order := range paginatedItems.Items.([]*entity.Order) {
		transResp := order.ToOrderListResponse()
		resps = append(resps, &transResp)
	}
	paginatedItems.Items = resps

	resp := dto.ResponseDto{Data: paginatedItems}
	ctx.JSON(http.StatusOK, resp)
}

func (h *OrderHandler) GetAllAdminOrders(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()

	getAllOrderQuery := queryparamdto.GetAllOrderAdminRequestQuery{}
	err = ctx.ShouldBindQuery(&getAllOrderQuery)
	if err != nil {
		return
	}

	err = h.validator.Validate(getAllOrderQuery)
	if err != nil {
		return
	}

	param := getAllOrderQuery.ToGetAllParams()
	if err != nil {
		return
	}

	paginatedItems, err := h.uc.GetAllOrders(ctx.Request.Context(), param)
	if err != nil {
		return
	}

	resps := make([]*responsedto.OrderListResponse, 0)
	for _, order := range paginatedItems.Items.([]*entity.Order) {
		transResp := order.ToOrderListResponse()
		resps = append(resps, &transResp)
	}
	paginatedItems.Items = resps

	resp := dto.ResponseDto{Data: paginatedItems}
	ctx.JSON(http.StatusOK, resp)
}

func (h *OrderHandler) GetById(ctx *gin.Context) {
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

	order, err := h.uc.GetOrderById(ctx.Request.Context(), uri.Id)
	if err != nil {
		return
	}

	resp := dto.ResponseDto{Data: order.ToOrderDetailResponse()}
	ctx.JSON(http.StatusOK, resp)
}

func (h *OrderHandler) ConfirmOrder(ctx *gin.Context) {
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

	order, err := h.uc.ConfirmOrder(ctx.Request.Context(), uri.Id)
	if err != nil {
		return
	}

	resp := dto.ResponseDto{Data: order.ToOrderStatusLogResponse()}
	ctx.JSON(http.StatusOK, resp)
}

func (h *OrderHandler) RejectOrder(ctx *gin.Context) {
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

	order, err := h.uc.RejectOrder(ctx.Request.Context(), uri.Id)
	if err != nil {
		return
	}

	resp := dto.ResponseDto{Data: order.ToOrderStatusLogResponse()}
	ctx.JSON(http.StatusOK, resp)
}

func (h *OrderHandler) ShipOrder(ctx *gin.Context) {
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

	order, err := h.uc.ShipOrder(ctx.Request.Context(), uri.Id)
	if err != nil {
		return
	}

	resp := dto.ResponseDto{Data: order.ToOrderStatusLogResponse()}
	ctx.JSON(http.StatusOK, resp)
}

func (h *OrderHandler) ReceiveOrder(ctx *gin.Context) {
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

	order, err := h.uc.ReceiveOrder(ctx.Request.Context(), uri.Id)
	if err != nil {
		return
	}

	resp := dto.ResponseDto{Data: order.ToOrderStatusLogResponse()}
	ctx.JSON(http.StatusOK, resp)
}

func (h *OrderHandler) CancelOrder(ctx *gin.Context) {
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

	order, err := h.uc.CancelOrder(ctx.Request.Context(), uri.Id)
	if err != nil {
		return
	}

	resp := dto.ResponseDto{Data: order.ToOrderStatusLogResponse()}
	ctx.JSON(http.StatusOK, resp)
}

func (h *OrderHandler) GetOrderLogs(ctx *gin.Context) {
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

	logs, err := h.uc.GetAllOrderLogsByOrderId(ctx.Request.Context(), uri.Id)
	if err != nil {
		return
	}
	var resps []*responsedto.OrderHistoryResponse
	for _, log := range logs {
		res := log.ToOrderStatusHistoryResponse()
		resps = append(resps, &res)
	}

	resp := dto.ResponseDto{Data: resps}
	ctx.JSON(http.StatusOK, resp)

}
