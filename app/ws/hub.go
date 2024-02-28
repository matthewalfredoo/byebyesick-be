package ws

import (
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/dto/responsedto"
)

type ConsultationSession struct {
	Id        int64             `json:"id"`
	DoctorId  int64             `json:"doctor_id"`
	PatientId int64             `json:"patient_id"`
	Clients   map[int64]*Client `json:"clients"`
}

type Hub struct {
	ConsultationSessions map[int64]*ConsultationSession
	Register             chan *Client
	Unregister           chan *Client
	Broadcast            chan *responsedto.WsConsultationMessage
}

func NewHub() *Hub {
	return &Hub{
		ConsultationSessions: make(map[int64]*ConsultationSession),
		Register:             make(chan *Client),
		Unregister:           make(chan *Client),
		Broadcast:            make(chan *responsedto.WsConsultationMessage, appconstant.BroadcastChannelBufferSize),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			if _, isRoomExist := h.ConsultationSessions[client.SessionId]; isRoomExist {
				r := h.ConsultationSessions[client.SessionId]

				if _, isClientExist := r.Clients[client.SenderId]; !isClientExist {
					r.Clients[client.SenderId] = client
				}
			}
		case client := <-h.Unregister:
			if _, isRoomExist := h.ConsultationSessions[client.SessionId]; isRoomExist {
				if _, isClientExist := h.ConsultationSessions[client.SessionId].Clients[client.SenderId]; isClientExist {
					delete(h.ConsultationSessions[client.SessionId].Clients, client.SenderId)
					close(client.Message)
				}
			}
		case message := <-h.Broadcast:
			if _, isRoomExist := h.ConsultationSessions[message.SessionId]; isRoomExist {

				for _, client := range h.ConsultationSessions[message.SessionId].Clients {
					client.Message <- message
				}
			}
		}
	}
}
