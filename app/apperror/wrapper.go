package apperror

type Wrapper struct {
	Code        int
	ErrorStored error
	Message     string
}

func NewWrapper(errorStored error, message ...string) *Wrapper {
	if len(message) > 0 {
		return &Wrapper{ErrorStored: errorStored, Message: message[0]}
	}
	return &Wrapper{ErrorStored: errorStored}
}

func (e *Wrapper) IsMessageEmpty() bool {
	return e.Message == ""
}

func (e *Wrapper) Error() string {
	if !e.IsMessageEmpty() {
		return e.Message
	}
	return e.ErrorStored.Error()
}
