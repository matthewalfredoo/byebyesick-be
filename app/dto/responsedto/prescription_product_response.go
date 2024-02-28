package responsedto

import (
	"time"
)

type PrescriptionProductResponse struct {
	Id             int64            `json:"id"`
	PrescriptionId int64            `json:"prescription_id,omitempty"`
	ProductId      int64            `json:"product_id"`
	Note           string           `json:"note"`
	CreatedAt      time.Time        `json:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at"`
	Product        *ProductResponse `json:"product,omitempty"`
}
