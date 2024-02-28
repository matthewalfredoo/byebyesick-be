package apperror

import (
	"fmt"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/entity"
)

type NotFound struct {
	Resource        entity.Resourcer
	FieldInResource string
	Value           any
}

func NewNotFound(resource entity.Resourcer, fieldInResource string, value any) *NotFound {
	return &NotFound{resource, fieldInResource, value}
}

func (e *NotFound) Error() string {
	return fmt.Sprintf("'%s' with '%s' of '%v' is not found",
		e.Resource.GetEntityName(), e.Resource.GetFieldStructTag(e.FieldInResource, appconstant.JsonStructTag), e.Value,
	)
}
