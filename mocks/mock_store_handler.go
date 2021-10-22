// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/handler/store_handler.go

// Package mocks is a generated GoMock package.
package mocks

import (
	http "net/http"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockStoreHandler is a mock of StoreHandler interface.
type MockStoreHandler struct {
	ctrl     *gomock.Controller
	recorder *MockStoreHandlerMockRecorder
}

// MockStoreHandlerMockRecorder is the mock recorder for MockStoreHandler.
type MockStoreHandlerMockRecorder struct {
	mock *MockStoreHandler
}

// NewMockStoreHandler creates a new mock instance.
func NewMockStoreHandler(ctrl *gomock.Controller) *MockStoreHandler {
	mock := &MockStoreHandler{ctrl: ctrl}
	mock.recorder = &MockStoreHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStoreHandler) EXPECT() *MockStoreHandlerMockRecorder {
	return m.recorder
}

// Flush mocks base method.
func (m *MockStoreHandler) Flush(w http.ResponseWriter, r *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Flush", w, r)
}

// Flush indicates an expected call of Flush.
func (mr *MockStoreHandlerMockRecorder) Flush(w, r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Flush", reflect.TypeOf((*MockStoreHandler)(nil).Flush), w, r)
}

// ServeHTTP mocks base method.
func (m *MockStoreHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ServeHTTP", w, r)
}

// ServeHTTP indicates an expected call of ServeHTTP.
func (mr *MockStoreHandlerMockRecorder) ServeHTTP(w, r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ServeHTTP", reflect.TypeOf((*MockStoreHandler)(nil).ServeHTTP), w, r)
}
