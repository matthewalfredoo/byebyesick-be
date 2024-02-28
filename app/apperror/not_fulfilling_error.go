package apperror

import (
	"fmt"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/entity"
)

type NotFulfilling struct {
	Resource          entity.Resourcer
	FieldInResource   string
	ExpectedCondition string
	ExpectedValue     any
	GotValue          any
}

func NewNotFulfilling(resource entity.Resourcer, fieldInResource string, expCond string, expVal, gotVal any) *NotFulfilling {
	return &NotFulfilling{resource, fieldInResource, expCond, expVal, gotVal}
}

func (e *NotFulfilling) Error() string {
	return fmt.Sprintf("resource '%s' field '%s' expects value that is '%s %v', but got '%v'",
		e.Resource.GetEntityName(), e.Resource.GetFieldStructTag(e.FieldInResource, appconstant.JsonStructTag),
		e.ExpectedCondition, e.ExpectedValue, e.GotValue,
	)
}
