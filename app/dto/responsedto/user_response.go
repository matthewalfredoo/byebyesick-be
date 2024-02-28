package responsedto

type UserResponse struct {
	Id         int64  `json:"id"`
	Email      string `json:"email"`
	UserRoleId int64  `json:"user_role_id"`
	IsVerified bool   `json:"is_verified"`
}
