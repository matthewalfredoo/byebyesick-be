package queryparamdto

import (
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/appdb"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/util"
	"strconv"
	"strings"
)

type GetAllDoctorsQuery struct {
	Search string `form:"search"`
	Limit  string `form:"limit"`
	Page   string `form:"page"`
}

func (q *GetAllDoctorsQuery) ToGetAllParams() (*GetAllParams, error) {
	param := NewGetAllParams()
	profile := new(entity.DoctorProfile)
	spec := new(entity.DoctorSpecialization)

	if q.Search != "" {
		words := strings.Split(q.Search, " ")
		wordToSearch := ""
		for _, word := range words {
			wordToSearch += "%" + word + "%"
		}
		param.WhereClauses = append(
			param.WhereClauses,
			appdb.NewWhere(profile.GetSqlColumnFromField("Name"), appdb.ILike, wordToSearch, appdb.OR),
			appdb.NewWhere(spec.GetSqlColumnFromField("Name"), appdb.ILike, wordToSearch),
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

	return param, nil
}
