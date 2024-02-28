package requestdto

import (
	"halodeksik-be/app/appdb"
	"halodeksik-be/app/entity"
)

type WsConsultationMessage struct {
	IsTyping    bool   `json:"is_typing"`
	MessageType int64  `json:"message_type"`
	Message     string `json:"message"`
	Attachment  string `json:"attachment"`
	SenderId    int64  `json:"-"`
	SessionId   int64  `json:"-"`
}

func (r *WsConsultationMessage) ToConsultationMessage() *entity.ConsultationMessage {
	return &entity.ConsultationMessage{
		MessageType: appdb.NewSqlNullInt64(r.MessageType),
		Message:     appdb.NewSqlNullString(r.Message),
		Attachment:  appdb.NewSqlNullString(r.Attachment),
		SenderId:    appdb.NewSqlNullInt64(r.SenderId),
		SessionId:   appdb.NewSqlNullInt64(r.SessionId),
	}
}
