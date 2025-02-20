package gorm_test

import (
	"reflect"
	"testing"

	"github.com/ngtrvu/zen-go/gorm"
	"github.com/ngtrvu/zen-go/utils"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type TestBase struct {
	CreatedAt int `gorm:"column:created_at" json:"created_at"`
	UpdatedAt int `gorm:"column:updated_at" json:"updated_at"`
}

type FKModel struct {
	ID int `gorm:"column:id;primary_key" json:"id"`
}

type TestModel struct {
	TestBase
	ID           int      `gorm:"column:id;primary_key"                        json:"id"`
	Code         string   `gorm:"column:code;unique;not null"                  json:"code"`
	Email        string   `gorm:"uniqueIndex:uidx_email,column:email,not null" json:"email"`
	Email2       string   `gorm:"column:email2;not null"                       json:"email2"`
	FKID         int      `gorm:"column:fk_id"                                 json:"fk_id"`
	FKModelField *FKModel `gorm:"foreignKey:FKID"                              json:"fk_model_field"`
}

func (TestModel) TableName() string {
	return "test_models"
}

func TestUtils_GetFKField(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	items := []*TestModel{}

	modelType := gorm.GetModelType(items)
	assert.Equal(t, "gorm_test.TestModel", modelType.String())

	fkField := gorm.GetFKFieldName(items, "FKModel")
	assert.Equal(t, "FKID", fkField)

	field, err := gorm.GetModelField(items, "FKID")
	assert.Equal(t, "FKID", field.Name)
	assert.Nil(t, err)

	fkColumnName, err := gorm.GetFKColumnName(items, "FKModel")
	assert.Equal(t, "fk_id", fkColumnName)
	assert.Nil(t, err)
}

func TestUtils_GetGormColumnNames(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	columnNames := gorm.GetGormColumnNames(TestModel{})
	assert.Equal(t, []string{"id", "code", "email", "email2", "fk_id", "created_at", "updated_at"}, columnNames)

	columnNames = gorm.GetGormColumnNames(FKModel{})
	assert.Equal(t, []string{"id"}, columnNames)
}

func TestUtils_GetTableName_GetSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	items := TestModel{}
	tableName := gorm.GetTableName(items)
	assert.Equal(t, "test_models", tableName)

	tableName = gorm.GetTableName(&items)
	assert.Equal(t, "test_models", tableName)

	items2 := []*TestModel{}
	tableName = gorm.GetTableName(items2)
	assert.Equal(t, "test_models", tableName)

	items3 := []*TestModel{{ID: 12}}
	tableName = gorm.GetTableName(items3)
	assert.Equal(t, "test_models", tableName)

	items4 := []TestModel{{ID: 12}}
	tableName = gorm.GetTableName(items4)
	assert.Equal(t, "test_models", tableName)
	tableName = gorm.GetTableName(&items4)
	assert.Equal(t, "test_models", tableName)
	tableName = gorm.GetTableName(nil)
	assert.Equal(t, "", tableName)

	items5 := reflect.ValueOf(items4).Interface()
	tableName = gorm.GetTableName(items5)
	assert.Equal(t, "test_models", tableName)

	items6 := utils.CreateArrayFromObject(items4)
	tableName = gorm.GetTableName(items6)
	assert.Equal(t, "test_models", tableName)
}
