package apperror

import (
	"fmt"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/entity"
)

type Forbidden struct {
	Resource               entity.Resourcer
	FieldCheckedInResource string
	ExpectedValue          any
	GotValue               any
}

func NewForbidden(resource entity.Resourcer, fieldInResource string, expectedValue, gotValue any) *Forbidden {
	return &Forbidden{
		Resource:               resource,
		FieldCheckedInResource: fieldInResource,
		ExpectedValue:          expectedValue,
		GotValue:               gotValue,
	}
}

func (e *Forbidden) Error() string {
	return fmt.Sprintf("'%s' with field '%s' of '%v' is expected to do the action but got '%v' instead",
		e.Resource.GetEntityName(),
		e.Resource.GetFieldStructTag(e.FieldCheckedInResource, appconstant.JsonStructTag),
		e.ExpectedValue,
		e.GotValue,
	)
}
