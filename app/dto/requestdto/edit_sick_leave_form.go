package requestdto

import (
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/util"
)

type EditSickLeaveForm struct {
	StartingDate string `json:"starting_date" validate:"required,datetime=2006-01-02"`
	EndingDate   string `json:"ending_date" validate:"required,datetime=2006-01-02"`
	Description  string `json:"description" validate:"required"`
}

func (r EditSickLeaveForm) ToSickLeaveForm() (entity.SickLeaveForm, error) {
	startingDate, _ := util.ParseDateTime(r.StartingDate, appconstant.TimeFormatQueryParam)
	endingDate, _ := util.ParseDateTime(r.EndingDate, appconstant.TimeFormatQueryParam)

	if startingDate.After(endingDate) {
		return entity.SickLeaveForm{}, apperror.ErrSickLeaveStartingDateShouldBeBeforeEndingDate
	}

	return entity.SickLeaveForm{
		StartingDate: startingDate,
		EndingDate:   endingDate,
		Description:  r.Description,
	}, nil
}
