package requestdto

import "halodeksik-be/app/entity"

type AddEditPrescriptionProduct struct {
	ProductId int64  `json:"product_id" validate:"required"`
	Note      string `json:"note" validate:"required"`
}

func (r AddEditPrescriptionProduct) ToPrescriptionProduct() *entity.PrescriptionProduct {
	return &entity.PrescriptionProduct{
		ProductId: r.ProductId,
		Note:      r.Note,
	}
}
