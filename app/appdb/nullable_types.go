package appdb

import (
	"database/sql"
	"time"
)

func NewSqlNullString(val ...string) sql.NullString {
	if len(val) > 0 {
		return sql.NullString{String: val[0], Valid: true}
	}
	return sql.NullString{}
}

func NewSqlNullInt64(val ...int64) sql.NullInt64 {
	if len(val) > 0 {
		return sql.NullInt64{Int64: val[0], Valid: true}
	}
	return sql.NullInt64{}
}

func NewSqlNullTime(val ...time.Time) sql.NullTime {
	if len(val) > 0 {
		return sql.NullTime{Time: val[0], Valid: true}
	}
	return sql.NullTime{}
}
