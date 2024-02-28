package requestdto

import (
	"mime/multipart"
)

type AddEditDoctorSpecializationImage struct {
	Image *multipart.FileHeader `form:"image" validate:"omitempty,filetype=png jpg jpeg,filesize=500"`
}
