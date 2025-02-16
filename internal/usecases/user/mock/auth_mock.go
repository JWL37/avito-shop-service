// Code generated by MockGen. DO NOT EDIT.
// Source: auth.go

// Package mock is a generated GoMock package.
package mock

import (
	models "avito-shop-service/internal/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockUserAuthenticater is a mock of UserAuthenticater interface.
type MockUserAuthenticater struct {
	ctrl     *gomock.Controller
	recorder *MockUserAuthenticaterMockRecorder
}

// MockUserAuthenticaterMockRecorder is the mock recorder for MockUserAuthenticater.
type MockUserAuthenticaterMockRecorder struct {
	mock *MockUserAuthenticater
}

// NewMockUserAuthenticater creates a new mock instance.
func NewMockUserAuthenticater(ctrl *gomock.Controller) *MockUserAuthenticater {
	mock := &MockUserAuthenticater{ctrl: ctrl}
	mock.recorder = &MockUserAuthenticaterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserAuthenticater) EXPECT() *MockUserAuthenticaterMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockUserAuthenticater) Create(arg0, arg1 string) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockUserAuthenticaterMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockUserAuthenticater)(nil).Create), arg0, arg1)
}

// GetUserByUsername mocks base method.
func (m *MockUserAuthenticater) GetUserByUsername(arg0 string) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByUsername", arg0)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByUsername indicates an expected call of GetUserByUsername.
func (mr *MockUserAuthenticaterMockRecorder) GetUserByUsername(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByUsername", reflect.TypeOf((*MockUserAuthenticater)(nil).GetUserByUsername), arg0)
}
