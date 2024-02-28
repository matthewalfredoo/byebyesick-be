package handler

import (
	"github.com/gin-gonic/gin"
	"halodeksik-be/app/dto"
	"halodeksik-be/app/dto/responsedto"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/usecase"
	"net/http"
)

type DrugClassificationHandler struct {
	uc usecase.DrugClassificationUseCase
}

func NewDrugClassificationHandler(uc usecase.DrugClassificationUseCase) *DrugClassificationHandler {
	return &DrugClassificationHandler{uc: uc}
}

func (h *DrugClassificationHandler) GetAllWithoutParams(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()
	paginatedItems, err := h.uc.GetAllDrugsWithoutParams(ctx.Request.Context())
	if err != nil {
		return
	}

	resps := make([]*responsedto.DrugClassificationResponse, 0)
	for _, drugClassification := range paginatedItems.Items.([]*entity.DrugClassification) {
		resps = append(resps, drugClassification.ToResponse())
	}
	paginatedItems.Items = resps

	resp := dto.ResponseDto{Data: paginatedItems}
	ctx.JSON(http.StatusOK, resp)
}
