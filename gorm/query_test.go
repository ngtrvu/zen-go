package gorm_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/ngtrvu/zen-go/gorm"
)

type ModelTest struct {
	ID          uuid.UUID   `json:"id"`
	Name        string      `json:"name"`
	FKModelID   uuid.UUID   `json:"fk_model_test_id"`
	FKModelTest FKModelTest `json:"fk_model_test"    gorm:"foreignKey:fk_model_test_id"`
}

type FKModelTest struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Code string    `json:"code"`
}

func TestQueryStatement(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testcases := []struct {
		Desc                 string
		Query                *gorm.Query
		ExpectFilterQueryStr string
		ExpectFilterParams   []interface{}
		ExpectSearchQueryStr string
		ExpcetSearchParams   []interface{}
		ExpectQueryStr       string
		ExpectParams         []interface{}
	}{
		{
			Desc: "Case 1: Filter 1 criteria type int",
			Query: &gorm.Query{
				Filter: gorm.Filter{
					Filters: []*gorm.FilterAttribute{
						{
							Field:    "status",
							Type:     "int",
							Operator: "=",
							Value:    1,
						},
					},
				},
			},
			ExpectFilterQueryStr: "status = ?",
			ExpectFilterParams:   []interface{}{1},
			ExpectSearchQueryStr: "",
			ExpcetSearchParams:   []interface{}{},
			ExpectQueryStr:       "status = ?",
			ExpectParams:         []interface{}{1},
		},
		{
			Desc: "Case 2: Filter 2 criteria type int + uuid",
			Query: &gorm.Query{
				Filter: gorm.Filter{
					Filters: []*gorm.FilterAttribute{
						{
							Field:    "status",
							Type:     "int",
							Operator: "=",
							Value:    1,
						},
						{
							Field:    "user_id",
							Type:     "uuid",
							Operator: "=",
							Value:    "148afffc-e35b-4aff-b8cc-b59e7b30a24c",
						},
					},
				},
			},
			ExpectFilterQueryStr: "status = ? AND user_id = ?",
			ExpectFilterParams:   []interface{}{1, "148afffc-e35b-4aff-b8cc-b59e7b30a24c"},
			ExpectSearchQueryStr: "",
			ExpcetSearchParams:   []interface{}{},
			ExpectQueryStr:       "status = ? AND user_id = ?",
			ExpectParams:         []interface{}{1, "148afffc-e35b-4aff-b8cc-b59e7b30a24c"},
		},
		{
			Desc: "Case 3: Add search filter",
			Query: &gorm.Query{
				Filter: gorm.Filter{
					Filters: []*gorm.FilterAttribute{
						{
							Field:    "status",
							Type:     "int",
							Operator: "=",
							Value:    1,
						},
						{
							Field:    "user_id",
							Type:     "uuid",
							Operator: "=",
							Value:    "148afffc-e35b-4aff-b8cc-b59e7b30a24c",
						},
					},
				},
				Search: gorm.Search{
					SearchFields: []*gorm.SearchAttribute{
						{
							Field:    "email",
							Type:     "string",
							Operator: "LIKE",
							Value:    "ngtrvu@gma",
						},
					},
				},
			},
			ExpectFilterQueryStr: "status = ? AND user_id = ?",
			ExpectFilterParams:   []interface{}{1, "148afffc-e35b-4aff-b8cc-b59e7b30a24c"},
			ExpectSearchQueryStr: "LOWER(UNACCENT(email)) LIKE LOWER(UNACCENT(?))",
			ExpcetSearchParams:   []interface{}{"%ngtrvu@gma%"},
			ExpectQueryStr:       "(status = ? AND user_id = ?) AND (LOWER(UNACCENT(email)) LIKE LOWER(UNACCENT(?)))",
			ExpectParams:         []interface{}{1, "148afffc-e35b-4aff-b8cc-b59e7b30a24c", "%ngtrvu@gma%"},
		},
		{
			Desc: "Case 4: Add search filter",
			Query: &gorm.Query{
				Filter: gorm.Filter{
					Filters: []*gorm.FilterAttribute{
						{
							Field:    "status",
							Type:     "int",
							Operator: "=",
							Value:    1,
						},
						{
							Field:    "user_id",
							Type:     "uuid",
							Operator: "=",
							Value:    "148afffc-e35b-4aff-b8cc-b59e7b30a24c",
						},
					},
				},
				Search: gorm.Search{
					SearchFields: []*gorm.SearchAttribute{
						{
							Field:    "email",
							Type:     "string",
							Operator: "LIKE",
							Value:    "ngtrvu@gma",
						},
						{
							Field:    "name",
							Type:     "string",
							Operator: "LIKE",
							Value:    "vu",
						},
						{
							Field:    "full_name",
							Type:     "string",
							Operator: "LIKE",
							Value:    "vu",
						},
					},
				},
			},
			ExpectFilterQueryStr: "status = ? AND user_id = ?",
			ExpectFilterParams:   []interface{}{1, "148afffc-e35b-4aff-b8cc-b59e7b30a24c"},
			ExpectSearchQueryStr: "LOWER(UNACCENT(email)) LIKE LOWER(UNACCENT(?)) OR LOWER(UNACCENT(name)) LIKE LOWER(UNACCENT(?)) OR LOWER(UNACCENT(full_name)) LIKE LOWER(UNACCENT(?))",
			ExpcetSearchParams:   []interface{}{"%ngtrvu@gma%", "%vu%", "%vu%"},
			ExpectQueryStr:       "(status = ? AND user_id = ?) AND (LOWER(UNACCENT(email)) LIKE LOWER(UNACCENT(?)) OR LOWER(UNACCENT(name)) LIKE LOWER(UNACCENT(?)) OR LOWER(UNACCENT(full_name)) LIKE LOWER(UNACCENT(?)))",
			ExpectParams: []interface{}{
				1,
				"148afffc-e35b-4aff-b8cc-b59e7b30a24c",
				"%ngtrvu@gma%",
				"%vu%",
				"%vu%",
			},
		},
		{
			Desc: "Case 5: search without filter",
			Query: &gorm.Query{
				Filter: gorm.Filter{
					Filters: []*gorm.FilterAttribute{},
				},
				Search: gorm.Search{
					SearchFields: []*gorm.SearchAttribute{
						{
							Field:    "email",
							Type:     "string",
							Operator: "LIKE",
							Value:    "ngtrvu@gma",
						},
						{
							Field:    "name",
							Type:     "string",
							Operator: "LIKE",
							Value:    "vu",
						},
						{
							Field:    "full_name",
							Type:     "string",
							Operator: "LIKE",
							Value:    "vu",
						},
					},
				},
			},
			ExpectFilterQueryStr: "",
			ExpectFilterParams:   []interface{}{},
			ExpectSearchQueryStr: "LOWER(UNACCENT(email)) LIKE LOWER(UNACCENT(?)) OR LOWER(UNACCENT(name)) LIKE LOWER(UNACCENT(?)) OR LOWER(UNACCENT(full_name)) LIKE LOWER(UNACCENT(?))",
			ExpcetSearchParams:   []interface{}{"%ngtrvu@gma%", "%vu%", "%vu%"},
			ExpectQueryStr:       "LOWER(UNACCENT(email)) LIKE LOWER(UNACCENT(?)) OR LOWER(UNACCENT(name)) LIKE LOWER(UNACCENT(?)) OR LOWER(UNACCENT(full_name)) LIKE LOWER(UNACCENT(?))",
			ExpectParams:         []interface{}{"%ngtrvu@gma%", "%vu%", "%vu%"},
		},
		{
			Desc: "Case 6: Filter is null",
			Query: &gorm.Query{
				Filter: gorm.Filter{
					Filters: []*gorm.FilterAttribute{
						{
							Field:    "data",
							Type:     "byte",
							Operator: gorm.OperatorIsNull,
							Value:    nil,
						},
					},
				},
			},
			ExpectFilterQueryStr: "data IS NULL",
			ExpectFilterParams:   []interface{}{},
			ExpectSearchQueryStr: "",
			ExpcetSearchParams:   []interface{}{},
			ExpectQueryStr:       "data IS NULL",
			ExpectParams:         []interface{}{},
		},
		{
			Desc: "Case 7: Filter is datetime",
			Query: &gorm.Query{
				Filter: gorm.Filter{
					Filters: []*gorm.FilterAttribute{
						{
							Field:    "data",
							Type:     gorm.FieldTypeDatetime,
							Operator: gorm.OperatorGreaterEqual,
							Value:    time.Date(2024, 1, 24, 14, 56, 28, 0, time.UTC).Add(-1 * time.Hour),
						},
					},
				},
			},
			ExpectFilterQueryStr: "data >= ?",
			ExpectFilterParams:   []interface{}{time.Date(2024, 1, 24, 13, 56, 28, 0, time.UTC)},
			ExpectSearchQueryStr: "",
			ExpcetSearchParams:   []interface{}{},
			ExpectQueryStr:       "data >= ?",
			ExpectParams:         []interface{}{time.Date(2024, 1, 24, 13, 56, 28, 0, time.UTC)},
		},
		{
			Desc: "Case 8: OR condition between filters",
			Query: &gorm.Query{
				Filter: gorm.Filter{
					Filters: []*gorm.FilterAttribute{
						{
							Field:    "status",
							Type:     "int",
							Operator: "=",
							Value:    1,
						},
						{
							Field:     "user_id",
							Type:      "uuid",
							Operator:  "=",
							Value:     "148afffc-e35b-4aff-b8cc-b59e7b30a24c",
							LogicalOp: gorm.LogicalOperatorOR,
						},
					},
				},
			},
			ExpectFilterQueryStr: "status = ? OR user_id = ?",
			ExpectFilterParams:   []interface{}{1, "148afffc-e35b-4aff-b8cc-b59e7b30a24c"},
			ExpectSearchQueryStr: "",
			ExpcetSearchParams:   []interface{}{},
			ExpectQueryStr:       "status = ? OR user_id = ?",
			ExpectParams:         []interface{}{1, "148afffc-e35b-4aff-b8cc-b59e7b30a24c"},
		},
		{
			Desc: "Case 9: Multiple conditions with AND/OR",
			Query: &gorm.Query{
				Filter: gorm.Filter{
					Filters: []*gorm.FilterAttribute{
						{
							Field:    "status",
							Type:     "int",
							Operator: "=",
							Value:    1,
						},
						{
							Field:     "user_id",
							Type:      "uuid",
							Operator:  "=",
							Value:     "148afffc-e35b-4aff-b8cc-b59e7b30a24c",
							LogicalOp: gorm.LogicalOperatorOR,
						},
						{
							Field:     "name",
							Type:      "string",
							Operator:  "=",
							Value:     "test",
							LogicalOp: gorm.LogicalOperatorAND,
						},
						{
							Field:     "age",
							Type:      "int",
							Operator:  ">",
							Value:     18,
							LogicalOp: gorm.LogicalOperatorOR,
						},
					},
				},
			},
			ExpectFilterQueryStr: "status = ? OR user_id = ? AND name = ? OR age > ?",
			ExpectFilterParams:   []interface{}{1, "148afffc-e35b-4aff-b8cc-b59e7b30a24c", "test", 18},
			ExpectSearchQueryStr: "",
			ExpcetSearchParams:   []interface{}{},
			ExpectQueryStr:       "status = ? OR user_id = ? AND name = ? OR age > ?",
			ExpectParams:         []interface{}{1, "148afffc-e35b-4aff-b8cc-b59e7b30a24c", "test", 18},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.Desc, func(t *testing.T) {
			// verify filter query
			filter := gorm.Filter{Filters: tc.Query.Filter.Filters}
			actualQueryStr, actualParams := filter.QueryStatement()

			assert.Equal(t, tc.ExpectFilterQueryStr, actualQueryStr)
			assert.Equal(t, tc.ExpectFilterParams, actualParams)

			// verify search query

			search := gorm.Search{SearchFields: tc.Query.Search.SearchFields}
			actualQueryStr, actualParams = search.QueryStatement()

			assert.Equal(t, tc.ExpectSearchQueryStr, actualQueryStr)
			assert.Equal(t, tc.ExpcetSearchParams, actualParams)
		})
	}
}

func TestSearchQueryStatement(t *testing.T) {
	query := gorm.SearchAttribute{
		Field:    "email",
		Type:     gorm.FieldTypeString,
		Operator: gorm.OperatorLike,
		Value:    "ngtrvu@gmail.com",
	}
	assert.Equal(t, "LOWER(UNACCENT(email)) LIKE LOWER(UNACCENT(?))", query.QueryStatement())

	query = gorm.SearchAttribute{
		Field:    "users.email",
		Type:     gorm.FieldTypeString,
		Operator: gorm.OperatorLike,
		Value:    "ngtrvu@gmail.com",
	}
	assert.Equal(t, "LOWER(UNACCENT(users.email)) LIKE LOWER(UNACCENT(?))", query.QueryStatement())

	query1 := gorm.Search{
		SearchFields: []*gorm.SearchAttribute{
			{
				Field:    "email",
				Type:     gorm.FieldTypeString,
				Operator: gorm.OperatorEqual,
				Value:    "ngtrvu@gmail.com",
			},
			{
				Field:    "users.email",
				Type:     gorm.FieldTypeString,
				Operator: gorm.OperatorLike,
				Value:    "ngtrvu@gmail.com",
			},
		},
	}
	queryStr, params := query1.QueryStatement()
	assert.Equal(t, "LOWER(email) = LOWER(?) OR LOWER(UNACCENT(users.email)) LIKE LOWER(UNACCENT(?))", queryStr)
	assert.Equal(t, []interface{}{"ngtrvu@gmail.com", "%ngtrvu@gmail.com%"}, params)
}
