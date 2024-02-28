package repository

import (
	"fmt"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/appdb"
	"halodeksik-be/app/dto/queryparamdto"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/util"
	"strings"
)

func buildQuery(initQuery string, resourcer entity.Resourcer, param *queryparamdto.GetAllParams, setLimit bool, setPaginated bool, initIndex ...int) (string, []interface{}) {
	var query strings.Builder
	var values []interface{}

	query.WriteString(initQuery)

	if len(param.WhereClauses) > 0 {
		query.WriteString(appdb.AND + " ")
	}

	indexPreparedStatement := 0
	if len(initIndex) > 0 {
		indexPreparedStatement = initIndex[0]
	}

	for index, whereClause := range param.WhereClauses {
		if whereClause.OpenParenthesis {
			query.WriteString(" ( ")
		}
		if whereClause.Condition == appdb.In {
			query.WriteString(fmt.Sprintf("%s %s (", whereClause.Column, whereClause.Condition))
			val := strings.Split(whereClause.Value.(string), ",")
			for idx, v := range val {
				indexPreparedStatement++
				query.WriteString(fmt.Sprintf("$%d", indexPreparedStatement))
				if idx != len(val)-1 {
					query.WriteString(",")
				}
				values = append(values, v)
			}
			query.WriteString(string(") " + whereClause.Logic))

			if index != len(param.WhereClauses)-1 {
				if util.IsEmptyString(string(whereClause.Logic)) {
					query.WriteString(appdb.AND + " ")
				}
				if !util.IsEmptyString(string(whereClause.Logic)) {
					query.WriteString(string(whereClause.Logic) + " ")
				}
			}
			continue
		}

		if whereClause.Value != nil {
			indexPreparedStatement++
			query.WriteString(fmt.Sprintf("%s %s $%d %s ", whereClause.Column, whereClause.Condition, indexPreparedStatement, whereClause.Logic))
		}
		if whereClause.Value == nil {
			query.WriteString(fmt.Sprintf("%s %s %s %s ", whereClause.Column, whereClause.Condition, appdb.Null, whereClause.Logic))
		}

		if whereClause.CloseParenthesis {
			query.WriteString(" ) ")
		}

		if index != len(param.WhereClauses)-1 && util.IsEmptyString(string(whereClause.Logic)) {
			query.WriteString(appdb.AND + " ")
		}

		if whereClause.Value != nil {
			values = append(values, whereClause.Value)
		}
	}

	if len(param.GroupClauses) > 0 {
		query.WriteString(" GROUP BY ")
		for index, groupClause := range param.GroupClauses {
			query.WriteString(fmt.Sprintf("%s", groupClause.Column))
			if index != len(param.GroupClauses)-1 {
				query.WriteString(", ")
			}
		}
	}

	if setLimit {
		for index, sortClause := range param.SortClauses {
			if index == 0 {
				query.WriteString(" ORDER BY ")
			}
			query.WriteString(fmt.Sprintf("%s %s", sortClause.Column, sortClause.Order))
			if index != len(param.SortClauses)-1 {
				query.WriteString(", ")
			}
			if index == len(param.SortClauses)-1 && setPaginated {
				query.WriteString(fmt.Sprintf(", %s ASC ", resourcer.GetSqlColumnFromField("Id")))
			}
		}
		if len(param.SortClauses) == 0 && setPaginated {
			query.WriteString(fmt.Sprintf(" ORDER BY %s ASC ", resourcer.GetSqlColumnFromField("Id")))
		}
	}

	if param.PageSize != nil {
		if *param.PageSize > appconstant.MaxGetAllPageSize {
			*param.PageSize = appconstant.MaxGetAllPageSize
		}
	}

	if setLimit && param.PageSize != nil {
		size := *param.PageSize
		if size > appconstant.MaxGetAllPageSize {
			size = appconstant.MaxGetAllPageSize
		}
		query.WriteString(fmt.Sprintf(" LIMIT $%d ", indexPreparedStatement+1))

		indexPreparedStatement += 1
		values = append(values, size)
	}

	if setLimit && param.PageId != nil && param.PageSize != nil {
		size := *param.PageSize
		offset := (*param.PageId - 1) * size
		query.WriteString(fmt.Sprintf(" OFFSET $%d ", indexPreparedStatement+1))

		indexPreparedStatement += 1
		values = append(values, offset)
	}

	return query.String(), values
}
