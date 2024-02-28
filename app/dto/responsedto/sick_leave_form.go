package responsedto

import (
	"time"
)

type SickLeaveFormResponse struct {
	Id           int64                                     `json:"id"`
	SessionId    int64                                     `json:"session_id"`
	StartingDate time.Time                                 `json:"starting_date"`
	EndingDate   time.Time                                 `json:"ending_date"`
	Description  string                                    `json:"description"`
	CreatedAt    string                                    `json:"created_at,omitempty"`
	UpdatedAt    string                                    `json:"updated_at,omitempty"`
	Prescription *PrescriptionResponse                     `json:"prescription,omitempty"`
	User         *PrescriptionSickLeaveUserProfileResponse `json:"user,omitempty"`
	Doctor       *PrescriptionSickLeaveUserProfileResponse `json:"doctor,omitempty"`
}
