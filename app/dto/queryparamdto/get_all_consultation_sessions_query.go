package queryparamdto

import (
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/appdb"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/util"
	"strconv"
	"strings"
)

type GetAllConsultationSessions struct {
	Status string `form:"status" validate:"omitempty,number,oneof=1 2"`
	SortBy string `form:"sort_by"`
	Sort   string `form:"sort"`
	Limit  string `form:"limit"`
	Page   string `form:"page"`
}

func (q GetAllConsultationSessions) ToGetAllParams() *GetAllParams {
	const sortByDate = "date"

	param := NewGetAllParams()
	session := new(entity.ConsultationSession)

	switch q.SortBy {
	case sortByDate:
		q.SortBy = session.GetSqlColumnFromField("UpdatedAt")
	default:
		q.SortBy = session.GetSqlColumnFromField("UpdatedAt")
	}
	sortClause := appdb.NewSort(q.SortBy)
	switch q.Sort {
	case strings.ToLower(string(appdb.OrderAsc)):
		sortClause.Order = appdb.OrderAsc
	default:
		sortClause.Order = appdb.OrderDesc
	}
	if !util.IsEmptyString(q.SortBy) {
		param.SortClauses = append(param.SortClauses, sortClause)
	}

	if !util.IsEmptyString(q.Status) {
		column := session.GetSqlColumnFromField("ConsultationSessionStatusId")
		status, _ := util.ParseInt64(q.Status)
		param.WhereClauses = append(
			param.WhereClauses,
			appdb.NewWhere(column, appdb.EqualTo, status),
		)
	}
	if util.IsEmptyString(q.Status) {
		column := session.GetSqlColumnFromField("ConsultationSessionStatusId")
		param.WhereClauses = append(
			param.WhereClauses,
			appdb.NewWhere(column, appdb.EqualTo, appconstant.ConsultationSessionStatusOngoing),
		)
	}

	pageSize := appconstant.DefaultGetAllPageSize
	if !util.IsEmptyString(q.Limit) {
		noPageSize, err := strconv.Atoi(q.Limit)
		if err == nil && noPageSize > 0 {
			pageSize = noPageSize
		}
	}
	param.PageSize = &pageSize

	pageId := 1
	if !util.IsEmptyString(q.Page) {
		noPageId, err := strconv.Atoi(q.Page)
		if err == nil && noPageId > 0 {
			pageId = noPageId
		}
	}
	param.PageId = &pageId

	return param
}
