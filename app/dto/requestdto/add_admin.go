package requestdto

import "halodeksik-be/app/entity"

type AddAdmin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

func (r AddAdmin) ToUser() entity.User {
	return entity.User{
		Email:    r.Email,
		Password: r.Password,
	}
}
