package entity

import (
	"database/sql"
	"halodeksik-be/app/dto/responsedto"
	"time"
)

type ConsultationSessionStatus struct {
	Id        int64        `json:"id"`
	Name      string       `json:"name"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt sql.NullTime `json:"deleted_at"`
}

func (e *ConsultationSessionStatus) ToResponse() *responsedto.ConsultationSessionStatusResponse {
	if e == nil {
		return nil
	}
	return &responsedto.ConsultationSessionStatusResponse{
		Id:   e.Id,
		Name: e.Name,
	}
}
