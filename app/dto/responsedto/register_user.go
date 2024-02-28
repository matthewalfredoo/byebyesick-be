package responsedto

type ResponseUser struct {
	Id         int64  `json:"id"`
	Email      string `json:"email"`
	UserRoleID int64  `json:"user_role_id"`
}
