package uriparamdto

type SickLeaveFormBySessionId struct {
	SessionId int64 `uri:"sessionId" validate:"required,number"`
}
