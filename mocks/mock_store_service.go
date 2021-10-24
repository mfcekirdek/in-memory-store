// Code generated by MockGen. DO NOT EDIT.
// Source: ./pkg/service/store_service.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockStoreService is a mock of StoreService interface.
type MockStoreService struct {
	ctrl     *gomock.Controller
	recorder *MockStoreServiceMockRecorder
}

// MockStoreServiceMockRecorder is the mock recorder for MockStoreService.
type MockStoreServiceMockRecorder struct {
	mock *MockStoreService
}

// NewMockStoreService creates a new mock instance.
func NewMockStoreService(ctrl *gomock.Controller) *MockStoreService {
	mock := &MockStoreService{ctrl: ctrl}
	mock.recorder = &MockStoreServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStoreService) EXPECT() *MockStoreServiceMockRecorder {
	return m.recorder
}

// Flush mocks base method.
func (m *MockStoreService) Flush() map[string]string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Flush")
	ret0, _ := ret[0].(map[string]string)
	return ret0
}

// Flush indicates an expected call of Flush.
func (mr *MockStoreServiceMockRecorder) Flush() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Flush", reflect.TypeOf((*MockStoreService)(nil).Flush))
}

// Get mocks base method.
func (m *MockStoreService) Get(key string) (map[string]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", key)
	ret0, _ := ret[0].(map[string]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockStoreServiceMockRecorder) Get(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockStoreService)(nil).Get), key)
}

// Set mocks base method.
func (m *MockStoreService) Set(key, value string) (map[string]string, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", key, value)
	ret0, _ := ret[0].(map[string]string)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// Set indicates an expected call of Set.
func (mr *MockStoreServiceMockRecorder) Set(key, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockStoreService)(nil).Set), key, value)
}
