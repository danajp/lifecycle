// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/buildpack/lifecycle (interfaces: Cache)

// Package testmock is a generated GoMock package.
package testmock

import (
	io "io"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"

	lifecycle "github.com/buildpack/lifecycle"
)

// MockCache is a mock of Cache interface
type MockCache struct {
	ctrl     *gomock.Controller
	recorder *MockCacheMockRecorder
}

// MockCacheMockRecorder is the mock recorder for MockCache
type MockCacheMockRecorder struct {
	mock *MockCache
}

// NewMockCache creates a new mock instance
func NewMockCache(ctrl *gomock.Controller) *MockCache {
	mock := &MockCache{ctrl: ctrl}
	mock.recorder = &MockCacheMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCache) EXPECT() *MockCacheMockRecorder {
	return m.recorder
}

// AddLayerFile mocks base method
func (m *MockCache) AddLayerFile(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddLayerFile", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddLayerFile indicates an expected call of AddLayerFile
func (mr *MockCacheMockRecorder) AddLayerFile(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddLayerFile", reflect.TypeOf((*MockCache)(nil).AddLayerFile), arg0, arg1)
}

// Commit mocks base method
func (m *MockCache) Commit() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Commit")
	ret0, _ := ret[0].(error)
	return ret0
}

// Commit indicates an expected call of Commit
func (mr *MockCacheMockRecorder) Commit() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Commit", reflect.TypeOf((*MockCache)(nil).Commit))
}

// Name mocks base method
func (m *MockCache) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name
func (mr *MockCacheMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockCache)(nil).Name))
}

// RetrieveLayer mocks base method
func (m *MockCache) RetrieveLayer(arg0 string) (io.ReadCloser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RetrieveLayer", arg0)
	ret0, _ := ret[0].(io.ReadCloser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RetrieveLayer indicates an expected call of RetrieveLayer
func (mr *MockCacheMockRecorder) RetrieveLayer(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RetrieveLayer", reflect.TypeOf((*MockCache)(nil).RetrieveLayer), arg0)
}

// RetrieveMetadata mocks base method
func (m *MockCache) RetrieveMetadata() (lifecycle.CacheMetadata, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RetrieveMetadata")
	ret0, _ := ret[0].(lifecycle.CacheMetadata)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RetrieveMetadata indicates an expected call of RetrieveMetadata
func (mr *MockCacheMockRecorder) RetrieveMetadata() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RetrieveMetadata", reflect.TypeOf((*MockCache)(nil).RetrieveMetadata))
}

// ReuseLayer mocks base method
func (m *MockCache) ReuseLayer(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReuseLayer", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReuseLayer indicates an expected call of ReuseLayer
func (mr *MockCacheMockRecorder) ReuseLayer(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReuseLayer", reflect.TypeOf((*MockCache)(nil).ReuseLayer), arg0)
}

// SetMetadata mocks base method
func (m *MockCache) SetMetadata(arg0 lifecycle.CacheMetadata) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetMetadata", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetMetadata indicates an expected call of SetMetadata
func (mr *MockCacheMockRecorder) SetMetadata(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetMetadata", reflect.TypeOf((*MockCache)(nil).SetMetadata), arg0)
}
