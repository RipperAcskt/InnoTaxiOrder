// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/RipperAcskt/innotaxiorder/internal/service (interfaces: DriverService)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	proto "github.com/RipperAcskt/innotaxi/pkg/proto"
	gomock "github.com/golang/mock/gomock"
)

// MockDriverService is a mock of DriverService interface.
type MockDriverService struct {
	ctrl     *gomock.Controller
	recorder *MockDriverServiceMockRecorder
}

// MockDriverServiceMockRecorder is the mock recorder for MockDriverService.
type MockDriverServiceMockRecorder struct {
	mock *MockDriverService
}

// NewMockDriverService creates a new mock instance.
func NewMockDriverService(ctrl *gomock.Controller) *MockDriverService {
	mock := &MockDriverService{ctrl: ctrl}
	mock.recorder = &MockDriverServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDriverService) EXPECT() *MockDriverServiceMockRecorder {
	return m.recorder
}

// SetRaiting mocks base method.
func (m *MockDriverService) SetRaiting(arg0 context.Context, arg1 proto.Raiting, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetRaiting", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetRaiting indicates an expected call of SetRaiting.
func (mr *MockDriverServiceMockRecorder) SetRaiting(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetRaiting", reflect.TypeOf((*MockDriverService)(nil).SetRaiting), arg0, arg1, arg2)
}

// SyncDriver mocks base method.
func (m *MockDriverService) SyncDriver(arg0 context.Context, arg1 []*proto.Driver) ([]*proto.Driver, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SyncDriver", arg0, arg1)
	ret0, _ := ret[0].([]*proto.Driver)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SyncDriver indicates an expected call of SyncDriver.
func (mr *MockDriverServiceMockRecorder) SyncDriver(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SyncDriver", reflect.TypeOf((*MockDriverService)(nil).SyncDriver), arg0, arg1)
}
