package gorm

import (
	"fmt"
	"strings"
)

type QuerySortOrder string

const (
	QuerySortASC  = "asc"
	QuerySortDESC = "desc"
)

const (
	FieldTypeInt      = "int"
	FieldTypeString   = "string"
	FieldTypeUUID     = "uuid"
	FieldTypeDatetime = "datetime"
	FieldTypeBoolean  = "boolean"
	FieldTypeNumberic = "numberic"
)

const (
	OperatorEqual        = "="
	OperatorNotEqual     = "!="
	OperatorIn           = "IN"
	OperatorNotIn        = "NOT IN"
	OperatorGreaterEqual = ">="
	OperatorGreater      = ">"
	OperatorLessEqual    = "<="
	OperatorLess         = "<"
	OperatorLike         = "LIKE"
	OperatorIsNull       = "IS NULL"
	OperatorIsNotNull    = "IS NOT NULL"
)

const (
	LogicalOperatorAND = "AND"
	LogicalOperatorOR  = "OR"
)

type Query struct {
	Limit      int
	Offset     int
	SortBy     string // order by
	SortOrder  QuerySortOrder
	Filter     Filter
	Search     Search
	SortFields []*SortField
}

type SearchField struct {
	Field        string
	JoinedColumn string
	Type         string
	Operator     string
}

type SortField struct {
	SortBy    string
	SortOrder QuerySortOrder
}

type FilterConfig struct {
	SearchFields []string
}

type Search struct {
	SearchFields []*SearchAttribute
}

type Filter struct {
	Filters []*FilterAttribute
}

type FilterAttribute struct {
	Field     string
	Type      string
	Operator  string
	Value     interface{}
	LogicalOp string // AND or OR
}

type SearchAttribute struct {
	Field        string
	JoinedColumn string
	Type         string
	Operator     string
	Value        interface{}
}

func (s *Search) QueryStatement() (string, []interface{}) {
	searchQueries := []string{}
	params := []interface{}{}
	for _, fa := range s.SearchFields {
		if fa.Operator == OperatorLike {
			params = append(params, fmt.Sprintf("%%%v%%", fa.Value))
		} else {
			params = append(params, fa.Value)
		}
		searchQueries = append(searchQueries, fa.QueryStatement())
	}

	searchQueryStr := strings.Join(searchQueries, " OR ")

	return searchQueryStr, params
}

func (s *SearchAttribute) QueryStatement() string {
	if s.Type == FieldTypeString && s.Operator == OperatorLike {
		return fmt.Sprintf("LOWER(UNACCENT(%s)) LIKE LOWER(UNACCENT(?))", s.Field)
	} else if s.Type == FieldTypeString {
		return fmt.Sprintf("LOWER(%s) %s LOWER(?)", s.Field, s.Operator)
	}

	return fmt.Sprintf("%s %s ?", s.Field, s.Operator)
}

func (f *Filter) QueryStatement() (string, []interface{}) {
	if len(f.Filters) == 0 {
		return "", []interface{}{}
	}

	filterQueries := []string{}
	params := []interface{}{}

	for i, fa := range f.Filters {
		// Build the condition
		condition := ""
		if fa.Operator == OperatorIn || fa.Operator == OperatorNotIn {
			condition = fmt.Sprintf("%s %s (?)", fa.Field, fa.Operator)
			params = append(params, fa.Value)
		} else if fa.Operator == OperatorIsNull || fa.Operator == OperatorIsNotNull {
			condition = fmt.Sprintf("%s %s", fa.Field, fa.Operator)
		} else {
			condition = fmt.Sprintf("%s %s ?", fa.Field, fa.Operator)
			params = append(params, fa.Value)
		}

		// Add logical operator if not the first condition
		if i > 0 {
			if fa.LogicalOp != "" {
				filterQueries = append(filterQueries, fa.LogicalOp)
			} else {
				filterQueries = append(filterQueries, LogicalOperatorAND)
			}
		}

		filterQueries = append(filterQueries, condition)
	}

	return strings.Join(filterQueries, " "), params
}

func (f *Filter) AddFilter(attr *FilterAttribute) {
	f.Filters = append(f.Filters, attr)
}

func (f *Search) AddSearchField(attr *SearchAttribute) {
	f.SearchFields = append(f.SearchFields, attr)
}

func (q *Query) SortStatement() string {
	if q.SortFields != nil && len(q.SortFields) > 0 {
		sortFields := []string{}
		for _, sortField := range q.SortFields {
			sortFields = append(sortFields, fmt.Sprintf("%s %s", sortField.SortBy, sortField.SortOrder))
		}
		return strings.Join(sortFields, ", ")
	} else if q.SortBy != "" && q.SortOrder != "" {
		// TODO: to be removed after refactoring, we use SortFields instead
		return fmt.Sprintf("%s %s", q.SortBy, q.SortOrder)
	}

	// TODO: remove default sort by created_at, should be handled by controller
	return "created_at DESC, id DESC"
}
