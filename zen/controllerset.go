package zen

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	common_gorm "github.com/ngtrvu/zen-go/gorm"
	"github.com/ngtrvu/zen-go/utils"
	"gorm.io/gorm"
)

type ControllerConfig struct {
	SearchFields []*common_gorm.SearchField
	DefaultSort  []*common_gorm.SortField
}

type ControllerSet struct {
	HttpHandler
	Service          BaseServiceInterface
	Model            interface{}
	ControllerConfig *ControllerConfig
}

func NewControllerSet(
	httpHandler *HttpHandler,
	service BaseServiceInterface,
	model interface{},
	controllerConfig *ControllerConfig,
) *ControllerSet {
	if model == nil {
		panic("model is required")
	}

	if controllerConfig == nil {
		controllerConfig = &ControllerConfig{}
	}

	return &ControllerSet{
		HttpHandler:      *httpHandler,
		Service:          service,
		Model:            model,
		ControllerConfig: controllerConfig,
	}
}

func (ctrl ControllerSet) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	query := ctrl.GetFilteringQueryset(r, ctrl.ControllerConfig.SearchFields, ctrl.ControllerConfig.DefaultSort)
	items := utils.CreateArrayFromObject(ctrl.Model)

	count, err := ctrl.Service.GetAll(ctx, query, &items)
	if err != nil {
		ctrl.ServerError(w, err)
		return
	}

	ctrl.SuccessWithPagination(w, r, items, count)
}

func (ctrl ControllerSet) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		ctrl.BadRequest(w, ErrBadRequest)
		return
	}

	item := utils.CreateInstanceFromObject(ctrl.Model)
	err = ctrl.Service.Get(ctx, id, &item)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		ctrl.NotFound(w, err)
		return
	}

	if err != nil {
		ctrl.ServerError(w, err)
		return
	}

	ctrl.Success(w, r, item)
}
