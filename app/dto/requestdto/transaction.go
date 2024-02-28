package requestdto

import "mime/multipart"

type RequestPaymentProof struct {
	PaymentProof *multipart.FileHeader `json:"payment_proof" form:"payment_proof" validate:"omitempty,filetype=png jpg jpeg,filesize=500"`
}
