package requestdto

import "halodeksik-be/app/entity"

type EditAdmin struct {
	Email    string `json:"email" validate:"omitempty,email"`
	Password string `json:"password" validate:"omitempty,min=8,max=72"`
}

func (r EditAdmin) ToUser() entity.User {
	return entity.User{
		Email:    r.Email,
		Password: r.Password,
	}
}
