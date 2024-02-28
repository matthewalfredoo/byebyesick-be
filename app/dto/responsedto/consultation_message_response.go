package responsedto

import (
	"time"
)

type ConsultationMessageResponse struct {
	Id         int64     `json:"id"`
	SessionId  int64     `json:"session_id,omitempty"`
	SenderId   int64     `json:"sender_id"`
	Message    string    `json:"message"`
	Attachment string    `json:"attachment"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
