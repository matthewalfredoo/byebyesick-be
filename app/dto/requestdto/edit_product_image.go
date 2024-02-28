package requestdto

import "mime/multipart"

type EditProductImage struct {
	Image *multipart.FileHeader `json:"image" form:"image" validate:"omitempty,filetype=png jpg jpeg,filesize=500"`
}
