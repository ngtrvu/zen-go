package utils_test

import (
	"testing"

	"github.com/ngtrvu/zen-go/utils"
	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Name string
}

func TestCreateVarFromObject(t *testing.T) {
	obj := TestStruct{Name: "test"}
	instance := utils.CreateInstanceFromObject(obj)
	assert.Equal(t, &TestStruct{}, instance)

	var obj2 TestStruct
	instance2 := utils.CreateInstanceFromObject(obj2)
	assert.Equal(t, &TestStruct{}, instance2)

	var obj3 *TestStruct
	instance3 := utils.CreateInstanceFromObject(obj3)
	assert.IsType(t, &TestStruct{}, instance3)

	var obj4 []*TestStruct
	instance4 := utils.CreateInstanceFromObject(obj4)
	assert.IsType(t, &TestStruct{}, instance4)

	instance5 := utils.CreateInstanceFromObject(&TestStruct{})
	assert.IsType(t, &TestStruct{}, instance5)
}

func TestCreateArrayFromObject(t *testing.T) {
	obj := TestStruct{Name: "test"}
	instance := utils.CreateArrayFromObject(obj)
	assert.IsType(t, []TestStruct{}, instance)

	var obj2 TestStruct
	instance2 := utils.CreateArrayFromObject(obj2)
	assert.IsType(t, []TestStruct{}, instance2)

	var obj3 *TestStruct
	instance3 := utils.CreateArrayFromObject(obj3)
	assert.IsType(t, []*TestStruct{}, instance3)

	var obj4 []*TestStruct
	instance4 := utils.CreateArrayFromObject(obj4)
	assert.IsType(t, []*TestStruct{}, instance4)

	instance5 := utils.CreateArrayFromObject(&TestStruct{})
	assert.IsType(t, []*TestStruct{}, instance5)
}
