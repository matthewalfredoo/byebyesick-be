package handler

import (
	"github.com/gin-gonic/gin"
	"halodeksik-be/app/appvalidator"
	"halodeksik-be/app/dto"
	"halodeksik-be/app/dto/queryparamdto"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/usecase"
	"net/http"
)

type ReportHandler struct {
	uc        usecase.ReportUseCase
	validator appvalidator.AppValidator
}

func NewReportHandler(uc usecase.ReportUseCase, v appvalidator.AppValidator) *ReportHandler {
	return &ReportHandler{uc: uc, validator: v}
}

func (h ReportHandler) GetAllSellPharmacy(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()

	getAllUserQuery := queryparamdto.GetAllPharmacySellForAdminQuery{}
	_ = ctx.ShouldBindQuery(&getAllUserQuery)

	param, err := getAllUserQuery.ToGetAllParams()
	if err != nil {
		return
	}

	paginatedItems, err := h.uc.GetSellsAllPharmacy(ctx.Request.Context(), getAllUserQuery.Year, param)
	if err != nil {
		return
	}

	resps := make([]*entity.SellReport, 0)
	for _, report := range paginatedItems.Items.([]*entity.SellReport) {
		resps = append(resps, report)
	}

	paginatedItems.Items = resps

	resp := dto.ResponseDto{Data: paginatedItems}
	ctx.JSON(http.StatusOK, resp)
}

func (h ReportHandler) GetAllSellPharmacyMonthly(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()

	getAllUserQuery := queryparamdto.GetAllPharmacySellsMonthlyForAdminQuery{}
	_ = ctx.ShouldBindQuery(&getAllUserQuery)

	param, err := getAllUserQuery.ToGetAllParams()
	if err != nil {
		return
	}

	paginatedItems, err := h.uc.GetSellsAllPharmacyMonthly(ctx.Request.Context(), getAllUserQuery.Year, param)
	if err != nil {
		return
	}

	resps := make([]*entity.SellReportMonthly, 0)
	for _, report := range paginatedItems.Items.([]*entity.SellReportMonthly) {
		resps = append(resps, report)
	}

	paginatedItems.Items = resps

	resp := dto.ResponseDto{Data: paginatedItems}
	ctx.JSON(http.StatusOK, resp)
}

func (h ReportHandler) GetAllSellsPharmacyAdmin(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()

	getAllUserQuery := queryparamdto.GetAllPharmacySellForAdminQuery{}
	_ = ctx.ShouldBindQuery(&getAllUserQuery)

	param, err := getAllUserQuery.ToGetAllParams()
	if err != nil {
		return
	}

	paginatedItems, err := h.uc.GetSellsAllAdminPharmacy(ctx.Request.Context(), getAllUserQuery.Year, param)
	if err != nil {
		return
	}

	resps := make([]*entity.SellReport, 0)
	for _, report := range paginatedItems.Items.([]*entity.SellReport) {
		resps = append(resps, report)
	}

	paginatedItems.Items = resps

	resp := dto.ResponseDto{Data: paginatedItems}
	ctx.JSON(http.StatusOK, resp)
}

func (h ReportHandler) GetAllSellPharmacyAdminMonthly(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()

	getAllUserQuery := queryparamdto.GetAllPharmacySellsMonthlyForAdminQuery{}
	_ = ctx.ShouldBindQuery(&getAllUserQuery)

	param, err := getAllUserQuery.ToGetAllParams()
	if err != nil {
		return
	}

	paginatedItems, err := h.uc.GetSellsAllAdminPharmacyMonthly(ctx.Request.Context(), getAllUserQuery.Year, param)
	if err != nil {
		return
	}

	resps := make([]*entity.SellReportMonthly, 0)
	for _, report := range paginatedItems.Items.([]*entity.SellReportMonthly) {
		resps = append(resps, report)
	}

	paginatedItems.Items = resps

	resp := dto.ResponseDto{Data: paginatedItems}
	ctx.JSON(http.StatusOK, resp)
}
