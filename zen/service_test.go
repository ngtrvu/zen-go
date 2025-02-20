package zen_test

import (
	"context"
	"testing"

	common_gorm "github.com/ngtrvu/zen-go/gorm"
	"github.com/ngtrvu/zen-go/zen"
	mocks "github.com/ngtrvu/zen-go/zen/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type TestItem struct {
	ID int
}

func TestService_GetAllSuccess(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepoInterface(ctrl)

	query := common_gorm.Query{
		Search: common_gorm.Search{
			SearchFields: []*common_gorm.SearchAttribute{{Field: "users.email", JoinedColumn: "user_id"}},
		},
		SortBy:    "created_at",
		SortOrder: common_gorm.QuerySortDESC,
	}

	items := []*TestItem{{ID: 1}, {ID: 2}}
	repo.EXPECT().
		GetCount(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(2, nil)

	repo.EXPECT().
		GetAll(gomock.Any(), gomock.Any(), gomock.Any()).
		SetArg(2, items).
		Return(nil)

	service := zen.NewBaseService(repo)

	var result interface{}
	service.GetAll(ctx, &query, &result)

	assert.Equal(t, 2, len(result.([]*TestItem)))
}
