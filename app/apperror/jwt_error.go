package apperror

import "fmt"

type AuthError struct {
	Err error
}

func (e *AuthError) Error() string {
	return fmt.Sprintf(e.Err.Error())
}
