package zen_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	common_gorm "github.com/ngtrvu/zen-go/gorm"
	"github.com/ngtrvu/zen-go/zen"
	mocks "github.com/ngtrvu/zen-go/zen/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type ModelTest struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func TestControllerSet_GetAllSuccess(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	httpHandler, _ := zen.NewHttpHandler(&zen.ZenConfig{})
	baseServiceMock := mocks.NewMockBaseServiceInterface(ctrl)

	config := &zen.ControllerConfig{
		SearchFields: []*common_gorm.SearchField{{Field: "id"}},
		DefaultSort:  []*common_gorm.SortField{{SortBy: "created_at", SortOrder: common_gorm.QuerySortDESC}},
	}
	controller := zen.NewControllerSet(httpHandler, baseServiceMock, ModelTest{}, config)
	items := []*ModelTest{
		{ID: uuid.New()},
		{ID: uuid.New()},
	}
	baseServiceMock.EXPECT().GetAll(gomock.Any(), gomock.Any(), gomock.Any()).SetArg(2, items)

	client := zen.NewTestClient(ctx)
	req := client.MakeRequest("GET", "/admin/v1/users", nil)
	controller.GetAll(&client.Writer, req)
	require.Equal(t, 200, client.Writer.Code)
}

func TestControllerSet_GetInvalid(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	httpHandler, _ := zen.NewHttpHandler(&zen.ZenConfig{})
	baseServiceMock := mocks.NewMockBaseServiceInterface(ctrl)

	config := &zen.ControllerConfig{}
	controller := zen.NewControllerSet(httpHandler, baseServiceMock, ModelTest{}, config)

	client := zen.NewTestClient(ctx)
	req := client.MakeRequest("GET", "/admin/v1/users/{id}", nil)

	controller.Get(&client.Writer, req)
	require.Equal(t, 400, client.Writer.Code)
}

func TestControllerSet_GetSuccess(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	httpHandler, _ := zen.NewHttpHandler(&zen.ZenConfig{})
	baseServiceMock := mocks.NewMockBaseServiceInterface(ctrl)

	config := &zen.ControllerConfig{}
	controller := zen.NewControllerSet(httpHandler, baseServiceMock, ModelTest{}, config)

	// Create http test client
	client := zen.NewTestClient(ctx)

	// test get all fund accounts
	id := uuid.New()
	client.RouterContext.URLParams.Add("id", id.String())

	req := client.MakeRequest("GET", "/admin/v1/users/{id}", nil)

	item := ModelTest{ID: id, Name: "test"}
	baseServiceMock.EXPECT().Get(gomock.Any(), id, gomock.Any()).SetArg(2, item).Times(1)

	controller.Get(&client.Writer, req)
	require.Equal(t, 200, client.Writer.Code)

	res, err := client.ResponseJSON()
	assert.NoError(t, err)

	data := res.Data.(map[string]interface{})
	assert.Equal(t, item.ID.String(), data["id"])
	assert.Equal(t, res.Error, "")
	assert.Equal(t, res.ErrorCode, "")
	assert.Equal(t, res.Pagination, nil)
}
