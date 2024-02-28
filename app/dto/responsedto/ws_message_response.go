package responsedto

import "time"

type WsConsultationMessage struct {
	IsTyping    bool      `json:"is_typing"`
	MessageType int64     `json:"message_type"`
	Message     string    `json:"message"`
	Attachment  string    `json:"attachment"`
	CreatedAt   time.Time `json:"created_at"`
	SenderId    int64     `json:"sender_id"`
	SessionId   int64     `json:"session_id"`
}
