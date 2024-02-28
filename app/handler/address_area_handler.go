package handler

import (
	"github.com/gin-gonic/gin"
	"halodeksik-be/app/appvalidator"
	"halodeksik-be/app/dto"
	"halodeksik-be/app/dto/requestdto"
	"halodeksik-be/app/dto/responsedto"
	"halodeksik-be/app/usecase"
	"net/http"
)

type AddressAreaHandler struct {
	uc        usecase.AddressAreaUseCase
	validator appvalidator.AppValidator
}

func NewAddressAreaHandler(uc usecase.AddressAreaUseCase) *AddressAreaHandler {
	return &AddressAreaHandler{uc: uc}
}

func (h *AddressAreaHandler) GetAllProvince(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()

	provinces, err := h.uc.GetAllProvinces(ctx.Request.Context())
	if err != nil {
		return
	}

	resps := make([]*responsedto.ProvinceResponse, 0)
	for _, province := range provinces {
		resps = append(resps, province.ToResponse())
	}
	resp := dto.ResponseDto{Data: resps}
	ctx.JSON(http.StatusOK, resp)
}

func (h *AddressAreaHandler) GetAllCities(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()

	cities, err := h.uc.GetAllCities(ctx.Request.Context())
	if err != nil {
		return
	}

	resps := make([]*responsedto.CityResponse, 0)
	for _, city := range cities {
		resps = append(resps, city.ToResponse())
	}
	resp := dto.ResponseDto{Data: resps}
	ctx.JSON(http.StatusOK, resp)
}

func (h *AddressAreaHandler) ValidateLatLong(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()

	req := requestdto.RequestValidateLatLong{}
	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		return
	}

	err = h.validator.Validate(req)
	if err != nil {
		return
	}

	err = h.uc.ValidateCityWithLatLong(ctx.Request.Context(), req.CityId, req.ProvinceId, req.Latitude, req.Longitude)
	if err != nil {
		return
	}

	resp := dto.ResponseDto{Data: "Latitude and Longitude is valid."}
	ctx.JSON(http.StatusOK, resp)
}
