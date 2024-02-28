package uriparamdto

type PrescriptionBySessionId struct {
	SessionId int64 `uri:"sessionId" validate:"required,number"`
}
