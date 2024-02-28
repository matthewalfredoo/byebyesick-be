package requestdto

import (
	"mime/multipart"
)

type AddEditManufacturerImage struct {
	Image *multipart.FileHeader `form:"image" validate:"omitempty,filetype=png jpg jpeg,filesize=500"`
}
