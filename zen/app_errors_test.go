package zen_test

import (
	"testing"

	"github.com/ngtrvu/zen-go/zen"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAppErrorTranslation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	appError := zen.NewAppError("unauthorized", "unauthorized").AddTranslation("vi", "Không có quyền truy cập")

	assert.Equal(t, "unauthorized", appError.Code)
	assert.Equal(t, "unauthorized", appError.Err.Error())
	assert.Equal(t, "Không có quyền truy cập", appError.Message)

	appError2 := zen.NewAppError("unauthorized", "unauthorized").
		AddTranslation("en", "Unauthorized something else")
	assert.Equal(t, "unauthorized", appError2.Code)
	assert.Equal(t, "unauthorized", appError2.Err.Error())
	assert.Equal(t, "unauthorized", appError2.Message)

	appError3 := zen.NewAppError("unauthorized", "unauthorized").
		AddTranslation("en", "Unauthorized something else").
		AddTranslation("vi", "Không có quyền truy cập")

	assert.Equal(t, "unauthorized", appError3.Code)
	assert.Equal(t, "unauthorized", appError3.Err.Error())
	assert.Equal(t, "Không có quyền truy cập", appError3.Message)

}
