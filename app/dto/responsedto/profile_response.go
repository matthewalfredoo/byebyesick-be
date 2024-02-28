package responsedto

type ProfileResponse struct {
	UserId       int64  `json:"user_id"`
	Name         string `json:"name"`
	ProfilePhoto string `json:"profile_photo"`
}
