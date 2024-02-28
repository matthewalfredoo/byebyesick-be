package entity

import "halodeksik-be/app/dto/responsedto"

type Profile struct {
	UserId       int64
	RoleId       int64
	Name         string
	ProfilePhoto string
}

func (e *Profile) ToResponse() *responsedto.ProfileResponse {
	if e == nil {
		return nil
	}
	return &responsedto.ProfileResponse{
		UserId:       e.UserId,
		Name:         e.Name,
		ProfilePhoto: e.ProfilePhoto,
	}
}
