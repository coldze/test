// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/coldze/test/middleware/sources (interfaces: DataBuilder)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	http "net/http"
	reflect "reflect"
)

// MockDataBuilder is a mock of DataBuilder interface
type MockDataBuilder struct {
	ctrl     *gomock.Controller
	recorder *MockDataBuilderMockRecorder
}

// MockDataBuilderMockRecorder is the mock recorder for MockDataBuilder
type MockDataBuilderMockRecorder struct {
	mock *MockDataBuilder
}

// NewMockDataBuilder creates a new mock instance
func NewMockDataBuilder(ctrl *gomock.Controller) *MockDataBuilder {
	mock := &MockDataBuilder{ctrl: ctrl}
	mock.recorder = &MockDataBuilderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDataBuilder) EXPECT() *MockDataBuilderMockRecorder {
	return m.recorder
}

// Build mocks base method
func (m *MockDataBuilder) Build() ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Build")
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Build indicates an expected call of Build
func (mr *MockDataBuilderMockRecorder) Build() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Build", reflect.TypeOf((*MockDataBuilder)(nil).Build))
}

// Header mocks base method
func (m *MockDataBuilder) Header() http.Header {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Header")
	ret0, _ := ret[0].(http.Header)
	return ret0
}

// Header indicates an expected call of Header
func (mr *MockDataBuilderMockRecorder) Header() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Header", reflect.TypeOf((*MockDataBuilder)(nil).Header))
}

// Write mocks base method
func (m *MockDataBuilder) Write(arg0 []byte) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Write", arg0)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Write indicates an expected call of Write
func (mr *MockDataBuilderMockRecorder) Write(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Write", reflect.TypeOf((*MockDataBuilder)(nil).Write), arg0)
}

// WriteHeader mocks base method
func (m *MockDataBuilder) WriteHeader(arg0 int) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "WriteHeader", arg0)
}

// WriteHeader indicates an expected call of WriteHeader
func (mr *MockDataBuilderMockRecorder) WriteHeader(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteHeader", reflect.TypeOf((*MockDataBuilder)(nil).WriteHeader), arg0)
}