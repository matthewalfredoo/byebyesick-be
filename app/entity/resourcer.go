package entity

type Resourcer interface {
	GetEntityName() string
	GetFieldStructTag(fieldName string, structTag string) string
	GetSqlColumnFromField(fieldName string) string
}
