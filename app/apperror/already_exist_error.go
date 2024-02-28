package apperror

import (
	"fmt"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/entity"
)

type AlreadyExist struct {
	Resource        entity.Resourcer
	FieldInResource string
	Value           any
}

func NewAlreadyExist(resource entity.Resourcer, fieldInResource string, value any) *AlreadyExist {
	return &AlreadyExist{resource, fieldInResource, value}
}

func (e *AlreadyExist) Error() string {
	return fmt.Sprintf("'%s' with value of '%v' on field '%s' already exists",
		e.Resource.GetEntityName(), e.Value, e.Resource.GetFieldStructTag(e.FieldInResource, appconstant.JsonStructTag),
	)
}
