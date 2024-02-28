package uriparamdto

type ResourceById struct {
	Id int64 `uri:"id" validate:"required,numeric"`
}
