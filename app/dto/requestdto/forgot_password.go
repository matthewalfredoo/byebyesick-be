package requestdto

type ResetPasswordRequest struct {
	Password string `json:"password" validate:"required,min=8,max=72"`
}
