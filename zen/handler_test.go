package zen_test

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	common_gorm "github.com/ngtrvu/zen-go/gorm"
	"github.com/ngtrvu/zen-go/zen"
	"github.com/stretchr/testify/assert"
)

func TestHandlerSuccessResponse(t *testing.T) {
	httpHandler, err := zen.NewHttpHandler(&zen.ZenConfig{})
	assert.Nil(t, err)

	var resp zen.Response

	// testing 200 ok
	request := httptest.NewRequest("GET", "/test", nil)
	writer := httptest.NewRecorder()
	httpHandler.Success(writer, request, nil)

	assert.Equal(t, 200, writer.Code)
	json.Unmarshal(writer.Body.Bytes(), &resp)
	assert.Equal(t, "", resp.Error)
	assert.Equal(t, true, resp.Success)

	// testing 201 created
	request = httptest.NewRequest("POST", "/test", nil)
	writer = httptest.NewRecorder()

	httpHandler.SuccessCreated(writer, request, nil)
	assert.Equal(t, 201, writer.Code)
	json.Unmarshal(writer.Body.Bytes(), &resp)
	assert.Equal(t, "", resp.Error)
	assert.Equal(t, true, resp.Success)
}

func TestHandlerFailedResponse(t *testing.T) {
	httpHandler, err := zen.NewHttpHandler(&zen.ZenConfig{})
	assert.Nil(t, err)

	var resp zen.Response

	// testing 400 bad request
	writer := httptest.NewRecorder()
	httpHandler.BadRequest(writer, nil)

	assert.Equal(t, 400, writer.Code)
	json.Unmarshal(writer.Body.Bytes(), &resp)
	assert.Equal(t, zen.ErrBadRequest.Message, resp.Error)
	assert.Equal(t, false, resp.Success)
}

func TestGetGormSearch(t *testing.T) {
	ctx := context.Background()
	httpHandler, err := zen.NewHttpHandler(&zen.ZenConfig{})
	assert.Nil(t, err)

	searchFields := []*common_gorm.SearchField{
		{Field: "users.email", JoinedColumn: "user_id"},
		{Field: "users.phone", JoinedColumn: "user_id"},
		{Field: "name"},
	}

	client := zen.NewTestClient(ctx)

	// testing with empty query string
	req := client.MakeRequest("GET", "/test", nil)
	searchQuery := httpHandler.GetGormSearch(req, searchFields)
	assert.Equal(t, 0, len(searchQuery.SearchFields))

	// testing with search query string
	req = client.MakeRequest("GET", "/test?search=dev@test.local", nil)
	searchQuery = httpHandler.GetGormSearch(req, searchFields)
	assert.Equal(t, 3, len(searchQuery.SearchFields))
	assert.Equal(t, "users.email", searchQuery.SearchFields[0].Field)
	assert.Equal(t, "dev@test.local", searchQuery.SearchFields[0].Value)
	assert.Equal(t, "user_id", searchQuery.SearchFields[0].JoinedColumn)
}

func TestGormQuery(t *testing.T) {
	ctx := context.Background()
	httpHandler, err := zen.NewHttpHandler(&zen.ZenConfig{})
	assert.Nil(t, err)

	searchFields := []*common_gorm.SearchField{
		{Field: "users.email", JoinedColumn: "user_id"},
		{Field: "users.phone", JoinedColumn: "user_id"},
		{Field: "name"},
	}

	client := zen.NewTestClient(ctx)

	// testing with empty query string
	req := client.MakeRequest("GET", "/test?page=3", nil)
	query := httpHandler.GetFilteringQueryset(req, searchFields, nil)
	assert.Equal(t, 0, len(query.Search.SearchFields))

	req = client.MakeRequest("GET", "/test?page=3&search=test", nil)
	query = httpHandler.GetFilteringQueryset(req, searchFields, nil)

	// verify paging
	assert.Equal(t, 24, query.Limit)
	assert.Equal(t, 48, query.Offset)
	assert.Equal(t, "", query.SortBy)

	// verify search query
	assert.Equal(t, common_gorm.QuerySortOrder(""), query.SortOrder)
	assert.Equal(t, 3, len(query.Search.SearchFields))
	assert.Equal(t, "users.email", query.Search.SearchFields[0].Field)
	assert.Equal(t, "user_id", query.Search.SearchFields[0].JoinedColumn)
	assert.Equal(t, "test", query.Search.SearchFields[0].Value)
	assert.Equal(t, "users.phone", query.Search.SearchFields[1].Field)
	assert.Equal(t, "user_id", query.Search.SearchFields[1].JoinedColumn)
	assert.Equal(t, "test", query.Search.SearchFields[1].Value)
	assert.Equal(t, "name", query.Search.SearchFields[2].Field)
}
