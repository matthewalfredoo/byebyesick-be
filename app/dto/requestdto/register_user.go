package requestdto

import (
	"halodeksik-be/app/entity"
	"mime/multipart"
)

type RequestToken struct {
	Email string `json:"email" validate:"required,email"`
}

type RequestTokenUrl struct {
	Token string `form:"token" validate:"required"`
}

type RequestRegisterUser struct {
	Email      string `json:"email" form:"email" validate:"required"`
	Name       string `json:"name" form:"name" validate:"required"`
	Password   string `json:"password" form:"password" validate:"required,min=8,max=72"`
	UserRoleId int64  `json:"user_role_id" form:"user_role_id" validate:"required"`
}

type RequestRegisterDoctorCertificate struct {
	Certificate *multipart.FileHeader `json:"certificate" form:"certificate" validate:"omitempty,filetype=pdf png jpg jpeg,filesize=500"`
}

func (u *RequestRegisterUser) ToUser() entity.User {
	return entity.User{
		Email:      u.Email,
		Password:   u.Password,
		UserRoleId: u.UserRoleId,
	}
}
