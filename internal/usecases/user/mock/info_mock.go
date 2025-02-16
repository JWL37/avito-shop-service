// Code generated by MockGen. DO NOT EDIT.
// Source: info.go

// Package mock is a generated GoMock package.
package mock

import (
	models "avito-shop-service/internal/models"
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockUserInfoGetter is a mock of UserInfoGetter interface.
type MockUserInfoGetter struct {
	ctrl     *gomock.Controller
	recorder *MockUserInfoGetterMockRecorder
}

// MockUserInfoGetterMockRecorder is the mock recorder for MockUserInfoGetter.
type MockUserInfoGetterMockRecorder struct {
	mock *MockUserInfoGetter
}

// NewMockUserInfoGetter creates a new mock instance.
func NewMockUserInfoGetter(ctrl *gomock.Controller) *MockUserInfoGetter {
	mock := &MockUserInfoGetter{ctrl: ctrl}
	mock.recorder = &MockUserInfoGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserInfoGetter) EXPECT() *MockUserInfoGetterMockRecorder {
	return m.recorder
}

// GetUserBalance mocks base method.
func (m *MockUserInfoGetter) GetUserBalance(arg0 context.Context, arg1 string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserBalance", arg0, arg1)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserBalance indicates an expected call of GetUserBalance.
func (mr *MockUserInfoGetterMockRecorder) GetUserBalance(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserBalance", reflect.TypeOf((*MockUserInfoGetter)(nil).GetUserBalance), arg0, arg1)
}

// GetUserInventory mocks base method.
func (m *MockUserInfoGetter) GetUserInventory(arg0 context.Context, arg1 string) ([]models.InventoryItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserInventory", arg0, arg1)
	ret0, _ := ret[0].([]models.InventoryItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserInventory indicates an expected call of GetUserInventory.
func (mr *MockUserInfoGetterMockRecorder) GetUserInventory(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserInventory", reflect.TypeOf((*MockUserInfoGetter)(nil).GetUserInventory), arg0, arg1)
}

// GetUserTransactions mocks base method.
func (m *MockUserInfoGetter) GetUserTransactions(arg0 context.Context, arg1 string) (models.CoinHistory, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserTransactions", arg0, arg1)
	ret0, _ := ret[0].(models.CoinHistory)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserTransactions indicates an expected call of GetUserTransactions.
func (mr *MockUserInfoGetterMockRecorder) GetUserTransactions(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserTransactions", reflect.TypeOf((*MockUserInfoGetter)(nil).GetUserTransactions), arg0, arg1)
}
