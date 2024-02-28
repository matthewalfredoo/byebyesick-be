package responsedto

type LoginResponse struct {
	UserId     int64  `json:"user_id"`
	Email      string `json:"email"`
	UserRoleId int64  `json:"user_role_id"`
	Name       string `json:"name"`
	Image      string `json:"image"`
	Token      string `json:"token"`
}

type GenericProfileResponse struct {
	Image string `json:"image"`
	Name  string `json:"name"`
	Token string `json:"token"`
}
