package gorm

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/ngtrvu/zen-go/log"
	"github.com/ngtrvu/zen-go/utils"
)

func GetTableName(src interface{}) string {
	if src == nil {
		return ""
	}

	var tableName string
	itemsValue := reflect.ValueOf(src)

	if itemsValue.Kind() == reflect.Ptr {
		itemsValue = itemsValue.Elem()
	}

	if itemsValue.Kind() == reflect.Interface {
		itemsValue = itemsValue.Elem()
	}

	if itemsValue.Kind() == reflect.Slice {
		elementType := itemsValue.Type().Elem()
		if elementType.Kind() == reflect.Ptr {
			elementType = elementType.Elem()
		}

		itemsValue = reflect.New(elementType).Elem()
	}

	instance := reflect.New(itemsValue.Type()).Interface()
	if tabler, ok := instance.(interface{ TableName() string }); ok {
		tableName = tabler.TableName()
	}

	return tableName
}

// Deprecated: use GetFKColumnName instead
func GetForeignKeyField(src interface{}, refTable string) string {
	var fkField string
	itemsValue := reflect.ValueOf(src)
	if itemsValue.Kind() == reflect.Ptr {
		itemsValue = itemsValue.Elem()
	}

	if itemsValue.Kind() == reflect.Slice {
		elementType := itemsValue.Type().Elem()
		if elementType.Kind() == reflect.Ptr {
			elementType = elementType.Elem()
		}
		instance := reflect.New(elementType).Interface()
		if tabler, ok := instance.(interface{ GetForeignKeyFieldName(refTable string) string }); ok {
			fkField = tabler.GetForeignKeyFieldName(refTable)
		}
	}

	return fkField
}

// GetFKCol returns the column name of the fk table
func GetFKCol(i interface{}, fkTableName string) (string, error) {
	vType := GetModelType(i)

	for i := 0; i < vType.NumField(); i++ {
		field := vType.Field(i)
		gormTag := field.Tag.Get("gorm")
		if !strings.Contains(gormTag, "foreignKey:") {
			continue
		}

		log.Info("gormTag: %s", gormTag)
		fkColumnName := strings.Split(gormTag, ":")[1]
		log.Info("fkColumnName: %s", fkColumnName)
		tableName := ""
		if tableName == fkTableName {
			jsonTag := field.Tag.Get("json")
			log.Info("gormTag: %s", jsonTag)
			return jsonTag, nil
		}
	}

	return "", fmt.Errorf("fk field not found for the table name %s", fkTableName)
}

// GetFKColumnName returns foreign key field name for ref table.
func GetFKColumnName(i interface{}, modelRefName string) (string, error) {
	fkFieldName := GetFKFieldName(i, modelRefName)
	if fkFieldName == "" {
		return "", fmt.Errorf("foreign key field name not found for model %s", modelRefName)
	}

	field, err := GetModelField(i, fkFieldName)
	if err == nil {
		gormTag := field.Tag.Get("gorm")
		if strings.Contains(gormTag, "column:") {
			return strings.Split(gormTag, ":")[1], nil
		}
	}
	return "", err
}

// GetModelType get model type from model variable or a list of models
func GetModelType(i interface{}) reflect.Type {
	v := reflect.ValueOf(i)

	// Check if i is a struct or a pointer to a struct
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	var vType reflect.Type
	if v.Kind() == reflect.Slice {
		vType = v.Type().Elem()
		if vType.Kind() == reflect.Ptr {
			vType = vType.Elem()
		}
	} else {
		vType = v.Type()
	}

	return vType
}

// GetFKFieldName returns foreign key model field name by giving field name. UserProgram
func GetFKFieldName(i interface{}, modelRefName string) string {
	vType := GetModelType(i)

	for i := 0; i < vType.NumField(); i++ {
		field := vType.Field(i)
		var fieldType reflect.Type

		// if the field is pointer, get its value's type
		if field.Type.Kind() == reflect.Ptr {
			fieldValueType := field.Type.Elem()
			fieldType = fieldValueType
		} else {
			fieldType = field.Type
		}

		if fieldType.String() == modelRefName || strings.HasSuffix(fieldType.String(), modelRefName) {
			gormTag := field.Tag.Get("gorm")
			if !strings.Contains(gormTag, "foreignKey:") {
				continue
			}

			return strings.Split(gormTag, ":")[1]
		}

	}

	return ""
}

// GetModelField returns model field by giving field name
func GetModelField(i interface{}, fieldName string) (reflect.StructField, error) {
	vType := GetModelType(i)

	for i := 0; i < vType.NumField(); i++ {
		field := vType.Field(i)
		if field.Name == fieldName {
			return field, nil
		}
	}

	return reflect.StructField{}, fmt.Errorf("field '%s' not found", fieldName)
}

// GetGormColumnNames returns column names of a gorm model
func GetGormColumnNames(modelStruct interface{}) []string {
	t := reflect.TypeOf(modelStruct)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() == reflect.Slice {
		t = t.Elem()
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
	}

	if t.Kind() != reflect.Struct {
		return []string{}
	}

	var columns []string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tags := strings.Split(field.Tag.Get("gorm"), ",")
		for _, tag := range tags {
			if strings.HasPrefix(tag, "column") {
				tagStr := strings.TrimPrefix(tag, "column:")
				tag := strings.Split(tagStr, ";")[0]
				columns = append(columns, tag)
			}
		}
	}

	baseT := utils.GetBaseStructType(modelStruct)
	if baseT == nil {
		return columns
	}

	for i := 0; i < baseT.NumField(); i++ {
		field := baseT.Field(i)
		gormTags := strings.Split(field.Tag.Get("gorm"), ",")
		for _, gormTag := range gormTags {
			if strings.HasPrefix(gormTag, "column") {
				columnTag := strings.TrimPrefix(gormTag, "column:")
				name := strings.Split(columnTag, ";")[0]
				columns = append(columns, name)
			}
		}
	}

	return columns
}
