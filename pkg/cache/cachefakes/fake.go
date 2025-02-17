// Code generated by MockGen. DO NOT EDIT.
// SourceURL: github.com/ainsleydev/webkit/pkg/cache (interfaces: Storage)
//
// Generated by this command:
//
//	mockgen -package=cachefakes -destination=cachefakes/fake.go . Storage
//
// Package cachefakes is a generated GoMock package.
package cachefakes

import (
	context "context"
	reflect "reflect"

	cache "github.com/ainsleydev/webkit/pkg/cache"
	gomock "go.uber.org/mock/gomock"
)

// MockStore is a mock of Storage interface.
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore.
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance.
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockStore) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockStoreMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockStore)(nil).Close))
}

// Delete mocks base method.
func (m *MockStore) Delete(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockStoreMockRecorder) Delete(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockStore)(nil).Delete), arg0, arg1)
}

// Flush mocks base method.
func (m *MockStore) Flush(arg0 context.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Flush", arg0)
}

// Flush indicates an expected call of Flush.
func (mr *MockStoreMockRecorder) Flush(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Flush", reflect.TypeOf((*MockStore)(nil).Flush), arg0)
}

// Get mocks base method.
func (m *MockStore) Get(arg0 context.Context, arg1 string, arg2 any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Get indicates an expected call of Get.
func (mr *MockStoreMockRecorder) Get(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockStore)(nil).Get), arg0, arg1, arg2)
}

// Invalidate mocks base method.
func (m *MockStore) Invalidate(arg0 context.Context, arg1 []string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Invalidate", arg0, arg1)
}

// Invalidate indicates an expected call of Invalidate.
func (mr *MockStoreMockRecorder) Invalidate(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Invalidate", reflect.TypeOf((*MockStore)(nil).Invalidate), arg0, arg1)
}

// Ping mocks base method.
func (m *MockStore) Ping(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockStoreMockRecorder) Ping(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockStore)(nil).Ping), arg0)
}

// Set mocks base method.
func (m *MockStore) Set(arg0 context.Context, arg1 string, arg2 any, arg3 cache.Options) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Set", arg0, arg1, arg2, arg3)
}

// Set indicates an expected call of Set.
func (mr *MockStoreMockRecorder) Set(arg0, arg1, arg2, arg3 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockStore)(nil).Set), arg0, arg1, arg2, arg3)
}
