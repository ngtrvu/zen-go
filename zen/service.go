package zen

import (
	"context"

	"github.com/google/uuid"
	common_gorm "github.com/ngtrvu/zen-go/gorm"
)

type BaseServiceInterface interface {
	Get(ctx context.Context, id uuid.UUID, item interface{}) error
	GetAll(ctx context.Context, query *common_gorm.Query, items interface{}) (int, error)
}

type BaseService struct {
	Repo RepoInterface
}

func NewBaseService(repo RepoInterface) *BaseService {
	return &BaseService{
		Repo: repo,
	}
}

func (service *BaseService) Get(ctx context.Context, id uuid.UUID, item interface{}) error {
	err := service.Repo.GetByUUID(ctx, id, item)
	if err != nil {
		return err
	}
	return nil
}

func (service *BaseService) GetAll(ctx context.Context, query *common_gorm.Query, items interface{}) (int, error) {
	if query == nil {
		query = &common_gorm.Query{
			Offset: 0,
			Limit:  10,
		}
	}

	count, err := service.Repo.GetCount(ctx, query, items)
	if err != nil {
		return 0, err
	}

	err = service.Repo.GetAll(ctx, query, items)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (service *BaseService) Create(ctx context.Context, item interface{}) error {
	// Save the new instance using the repository layer
	if err := service.Repo.Create(ctx, item); err != nil {
		return err
	}

	return nil
}

func (service *BaseService) Update(ctx context.Context, item interface{}, params map[string]interface{}) error {
	// Update the instance using the repository layer
	if err := service.Repo.UpdatePartial(ctx, item, params); err != nil {
		return err
	}

	return nil
}
