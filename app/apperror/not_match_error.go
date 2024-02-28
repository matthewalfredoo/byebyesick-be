package apperror

import (
	"fmt"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/entity"
)

type NotMatch struct {
	Resource              entity.Resourcer
	UniqueFieldInResource string
	UniqueFieldValue      any
	FieldInResource       string
	Value                 any
}

func NewNotMatch(resource entity.Resourcer, uniqueField string, uniqueFieldValue any, fieldInResource string, value any) *NotMatch {
	return &NotMatch{
		Resource:              resource,
		UniqueFieldInResource: uniqueField,
		UniqueFieldValue:      uniqueFieldValue,
		FieldInResource:       fieldInResource,
		Value:                 value,
	}
}

func (e *NotMatch) Error() string {
	return fmt.Sprintf("resource '%s' with field '%s' of '%v' does not match value '%v' on field '%s'",
		e.Resource.GetEntityName(),
		e.Resource.GetFieldStructTag(e.UniqueFieldInResource, appconstant.JsonStructTag),
		e.UniqueFieldValue,
		e.Value,
		e.Resource.GetFieldStructTag(e.FieldInResource, appconstant.JsonStructTag),
	)
}
