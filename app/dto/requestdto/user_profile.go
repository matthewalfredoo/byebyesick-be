package requestdto

import (
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/util"
)

type RequestUserProfile struct {
	Name        string `json:"name" form:"name" validate:"required"`
	DateOfBirth string `json:"date_of_birth" form:"date_of_birth" validate:"required,datetime=2006-01-02"`
}

func (p *RequestUserProfile) ToUserProfile() (entity.UserProfile, error) {
	dateTime, _ := util.ParseDateTime(p.DateOfBirth, appconstant.TimeFormatQueryParam)
	return entity.UserProfile{
		Name:        p.Name,
		DateOfBirth: dateTime,
	}, nil
}
