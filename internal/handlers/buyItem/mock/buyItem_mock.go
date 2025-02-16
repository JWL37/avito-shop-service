// Code generated by MockGen. DO NOT EDIT.
// Source: buyItem.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockItemBuyer is a mock of ItemBuyer interface.
type MockItemBuyer struct {
	ctrl     *gomock.Controller
	recorder *MockItemBuyerMockRecorder
}

// MockItemBuyerMockRecorder is the mock recorder for MockItemBuyer.
type MockItemBuyerMockRecorder struct {
	mock *MockItemBuyer
}

// NewMockItemBuyer creates a new mock instance.
func NewMockItemBuyer(ctrl *gomock.Controller) *MockItemBuyer {
	mock := &MockItemBuyer{ctrl: ctrl}
	mock.recorder = &MockItemBuyerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockItemBuyer) EXPECT() *MockItemBuyerMockRecorder {
	return m.recorder
}

// BuyItem mocks base method.
func (m *MockItemBuyer) BuyItem(arg0 context.Context, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BuyItem", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// BuyItem indicates an expected call of BuyItem.
func (mr *MockItemBuyerMockRecorder) BuyItem(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BuyItem", reflect.TypeOf((*MockItemBuyer)(nil).BuyItem), arg0, arg1, arg2)
}
