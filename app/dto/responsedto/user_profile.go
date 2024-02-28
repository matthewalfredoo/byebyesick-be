package responsedto

type UserProfileResponse struct {
	Id           int64  `json:"id"`
	Email        string `json:"email"`
	UserRoleId   int64  `json:"user_role_id"`
	IsVerified   bool   `json:"is_verified"`
	Name         string `json:"name"`
	ProfilePhoto string `json:"profile_photo"`
	DateOfBirth  string `json:"date_of_birth"`
}
