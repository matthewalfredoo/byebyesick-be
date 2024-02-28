package appdb

type WhereClause struct {
	Column           string
	Condition        DBCondition
	Value            interface{}
	Logic            DBLogic
	OpenParenthesis  bool
	CloseParenthesis bool
}

func NewWhere(column string, condition DBCondition, value interface{}, logic ...DBLogic) WhereClause {
	if len(logic) > 0 {
		return WhereClause{Column: column, Condition: condition, Value: value, Logic: logic[0]}
	}
	return WhereClause{Column: column, Condition: condition, Value: value}
}

func NewWhereParenthesis(column string, condition DBCondition, value interface{}, open bool, close bool, logic ...DBLogic) WhereClause {
	if len(logic) > 0 {
		return WhereClause{Column: column, Condition: condition, Value: value, OpenParenthesis: open, CloseParenthesis: close, Logic: logic[0]}
	}
	return WhereClause{Column: column, Condition: condition, Value: value, OpenParenthesis: open, CloseParenthesis: close}
}
