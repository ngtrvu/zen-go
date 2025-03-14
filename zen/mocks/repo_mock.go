// Code generated by MockGen. DO NOT EDIT.
// Source: frameworks/common/go/zen/repo.go
//
// Generated by this command:
//
//	mockgen -source=frameworks/common/go/zen/repo.go -destination=frameworks/common/go/zen/mocks/repo_mock.go
//

// Package mock_zen is a generated GoMock package.
package mock_zen

import (
	context "context"
	reflect "reflect"

	uuid "github.com/google/uuid"
	gomock "go.uber.org/mock/gomock"
	gorm "gorm.io/gorm"
	clause "gorm.io/gorm/clause"
	gorm0 "github.com/ngtrvu/zen-go/gorm"
	zen "github.com/ngtrvu/zen-go/zen"
)

// MockRepoInterface is a mock of RepoInterface interface.
type MockRepoInterface struct {
	ctrl     *gomock.Controller
	recorder *MockRepoInterfaceMockRecorder
	isgomock struct{}
}

// MockRepoInterfaceMockRecorder is the mock recorder for MockRepoInterface.
type MockRepoInterfaceMockRecorder struct {
	mock *MockRepoInterface
}

// NewMockRepoInterface creates a new mock instance.
func NewMockRepoInterface(ctrl *gomock.Controller) *MockRepoInterface {
	mock := &MockRepoInterface{ctrl: ctrl}
	mock.recorder = &MockRepoInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepoInterface) EXPECT() *MockRepoInterfaceMockRecorder {
	return m.recorder
}

// BuildJoins mocks base method.
func (m *MockRepoInterface) BuildJoins(items any, query *gorm0.Query) []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BuildJoins", items, query)
	ret0, _ := ret[0].([]string)
	return ret0
}

// BuildJoins indicates an expected call of BuildJoins.
func (mr *MockRepoInterfaceMockRecorder) BuildJoins(items, query any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BuildJoins", reflect.TypeOf((*MockRepoInterface)(nil).BuildJoins), items, query)
}

// BulkUpdateOrCreate mocks base method.
func (m *MockRepoInterface) BulkUpdateOrCreate(ctx context.Context, items any, columns []clause.Column, updateColumns []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BulkUpdateOrCreate", ctx, items, columns, updateColumns)
	ret0, _ := ret[0].(error)
	return ret0
}

// BulkUpdateOrCreate indicates an expected call of BulkUpdateOrCreate.
func (mr *MockRepoInterfaceMockRecorder) BulkUpdateOrCreate(ctx, items, columns, updateColumns any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BulkUpdateOrCreate", reflect.TypeOf((*MockRepoInterface)(nil).BulkUpdateOrCreate), ctx, items, columns, updateColumns)
}

// CommitTransaction mocks base method.
func (m *MockRepoInterface) CommitTransaction(ctx context.Context) *gorm.DB {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CommitTransaction", ctx)
	ret0, _ := ret[0].(*gorm.DB)
	return ret0
}

// CommitTransaction indicates an expected call of CommitTransaction.
func (mr *MockRepoInterfaceMockRecorder) CommitTransaction(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CommitTransaction", reflect.TypeOf((*MockRepoInterface)(nil).CommitTransaction), ctx)
}

// Create mocks base method.
func (m *MockRepoInterface) Create(ctx context.Context, item any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, item)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockRepoInterfaceMockRecorder) Create(ctx, item any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockRepoInterface)(nil).Create), ctx, item)
}

// CreateMany mocks base method.
func (m *MockRepoInterface) CreateMany(ctx context.Context, items any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateMany", ctx, items)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateMany indicates an expected call of CreateMany.
func (mr *MockRepoInterfaceMockRecorder) CreateMany(ctx, items any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateMany", reflect.TypeOf((*MockRepoInterface)(nil).CreateMany), ctx, items)
}

// Delete mocks base method.
func (m *MockRepoInterface) Delete(ctx context.Context, item any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, item)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockRepoInterfaceMockRecorder) Delete(ctx, item any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockRepoInterface)(nil).Delete), ctx, item)
}

// DeleteMany mocks base method.
func (m *MockRepoInterface) DeleteMany(ctx context.Context, condition any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteMany", ctx, condition)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteMany indicates an expected call of DeleteMany.
func (mr *MockRepoInterfaceMockRecorder) DeleteMany(ctx, condition any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteMany", reflect.TypeOf((*MockRepoInterface)(nil).DeleteMany), ctx, condition)
}

// FirstOrCreate mocks base method.
func (m *MockRepoInterface) FirstOrCreate(ctx context.Context, item, conditions any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FirstOrCreate", ctx, item, conditions)
	ret0, _ := ret[0].(error)
	return ret0
}

// FirstOrCreate indicates an expected call of FirstOrCreate.
func (mr *MockRepoInterfaceMockRecorder) FirstOrCreate(ctx, item, conditions any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FirstOrCreate", reflect.TypeOf((*MockRepoInterface)(nil).FirstOrCreate), ctx, item, conditions)
}

// Get mocks base method.
func (m *MockRepoInterface) Get(ctx context.Context, query, item any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, query, item)
	ret0, _ := ret[0].(error)
	return ret0
}

// Get indicates an expected call of Get.
func (mr *MockRepoInterfaceMockRecorder) Get(ctx, query, item any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockRepoInterface)(nil).Get), ctx, query, item)
}

// GetAll mocks base method.
func (m *MockRepoInterface) GetAll(ctx context.Context, query *gorm0.Query, items any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", ctx, query, items)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetAll indicates an expected call of GetAll.
func (mr *MockRepoInterfaceMockRecorder) GetAll(ctx, query, items any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockRepoInterface)(nil).GetAll), ctx, query, items)
}

// GetByUUID mocks base method.
func (m *MockRepoInterface) GetByUUID(ctx context.Context, id uuid.UUID, item any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByUUID", ctx, id, item)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetByUUID indicates an expected call of GetByUUID.
func (mr *MockRepoInterfaceMockRecorder) GetByUUID(ctx, id, item any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByUUID", reflect.TypeOf((*MockRepoInterface)(nil).GetByUUID), ctx, id, item)
}

// GetCount mocks base method.
func (m *MockRepoInterface) GetCount(ctx context.Context, query *gorm0.Query, items any) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCount", ctx, query, items)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCount indicates an expected call of GetCount.
func (mr *MockRepoInterfaceMockRecorder) GetCount(ctx, query, items any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCount", reflect.TypeOf((*MockRepoInterface)(nil).GetCount), ctx, query, items)
}

// GetDB mocks base method.
func (m *MockRepoInterface) GetDB(ctx context.Context) *gorm.DB {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDB", ctx)
	ret0, _ := ret[0].(*gorm.DB)
	return ret0
}

// GetDB indicates an expected call of GetDB.
func (mr *MockRepoInterfaceMockRecorder) GetDB(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDB", reflect.TypeOf((*MockRepoInterface)(nil).GetDB), ctx)
}

// GroupByField mocks base method.
func (m *MockRepoInterface) GroupByField(ctx context.Context, db *gorm.DB, field, fieldCount string, result any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GroupByField", ctx, db, field, fieldCount, result)
	ret0, _ := ret[0].(error)
	return ret0
}

// GroupByField indicates an expected call of GroupByField.
func (mr *MockRepoInterfaceMockRecorder) GroupByField(ctx, db, field, fieldCount, result any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GroupByField", reflect.TypeOf((*MockRepoInterface)(nil).GroupByField), ctx, db, field, fieldCount, result)
}

// RollbackTransaction mocks base method.
func (m *MockRepoInterface) RollbackTransaction(ctx context.Context) *gorm.DB {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RollbackTransaction", ctx)
	ret0, _ := ret[0].(*gorm.DB)
	return ret0
}

// RollbackTransaction indicates an expected call of RollbackTransaction.
func (mr *MockRepoInterfaceMockRecorder) RollbackTransaction(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RollbackTransaction", reflect.TypeOf((*MockRepoInterface)(nil).RollbackTransaction), ctx)
}

// StartTransaction mocks base method.
func (m *MockRepoInterface) StartTransaction(ctx context.Context) *gorm.DB {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StartTransaction", ctx)
	ret0, _ := ret[0].(*gorm.DB)
	return ret0
}

// StartTransaction indicates an expected call of StartTransaction.
func (mr *MockRepoInterfaceMockRecorder) StartTransaction(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartTransaction", reflect.TypeOf((*MockRepoInterface)(nil).StartTransaction), ctx)
}

// Update mocks base method.
func (m *MockRepoInterface) Update(ctx context.Context, item any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, item)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockRepoInterfaceMockRecorder) Update(ctx, item any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockRepoInterface)(nil).Update), ctx, item)
}

// UpdateLocking mocks base method.
func (m *MockRepoInterface) UpdateLocking(ctx context.Context, item any, params map[string]any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateLocking", ctx, item, params)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateLocking indicates an expected call of UpdateLocking.
func (mr *MockRepoInterfaceMockRecorder) UpdateLocking(ctx, item, params any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateLocking", reflect.TypeOf((*MockRepoInterface)(nil).UpdateLocking), ctx, item, params)
}

// UpdateOrCreate mocks base method.
func (m *MockRepoInterface) UpdateOrCreate(ctx context.Context, query, getItem, updateItem any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateOrCreate", ctx, query, getItem, updateItem)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateOrCreate indicates an expected call of UpdateOrCreate.
func (mr *MockRepoInterfaceMockRecorder) UpdateOrCreate(ctx, query, getItem, updateItem any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateOrCreate", reflect.TypeOf((*MockRepoInterface)(nil).UpdateOrCreate), ctx, query, getItem, updateItem)
}

// UpdatePartial mocks base method.
func (m *MockRepoInterface) UpdatePartial(ctx context.Context, item any, params map[string]any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePartial", ctx, item, params)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdatePartial indicates an expected call of UpdatePartial.
func (mr *MockRepoInterfaceMockRecorder) UpdatePartial(ctx, item, params any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePartial", reflect.TypeOf((*MockRepoInterface)(nil).UpdatePartial), ctx, item, params)
}

// WithScopes mocks base method.
func (m *MockRepoInterface) WithScopes(scopes ...zen.Scope) zen.RepoInterface {
	m.ctrl.T.Helper()
	varargs := []any{}
	for _, a := range scopes {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "WithScopes", varargs...)
	ret0, _ := ret[0].(zen.RepoInterface)
	return ret0
}

// WithScopes indicates an expected call of WithScopes.
func (mr *MockRepoInterfaceMockRecorder) WithScopes(scopes ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithScopes", reflect.TypeOf((*MockRepoInterface)(nil).WithScopes), scopes...)
}
