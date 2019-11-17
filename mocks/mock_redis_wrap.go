// Code generated by MockGen. DO NOT EDIT.
// Source: middleware/sources/redis_wrap.go

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
	time "time"
)

// MockRedisWrap is a mock of RedisWrap interface
type MockRedisWrap struct {
	ctrl     *gomock.Controller
	recorder *MockRedisWrapMockRecorder
}

// MockRedisWrapMockRecorder is the mock recorder for MockRedisWrap
type MockRedisWrapMockRecorder struct {
	mock *MockRedisWrap
}

// NewMockRedisWrap creates a new mock instance
func NewMockRedisWrap(ctrl *gomock.Controller) *MockRedisWrap {
	mock := &MockRedisWrap{ctrl: ctrl}
	mock.recorder = &MockRedisWrapMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRedisWrap) EXPECT() *MockRedisWrapMockRecorder {
	return m.recorder
}

// Set mocks base method
func (m *MockRedisWrap) Set(key string, data interface{}, ttl time.Duration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", key, data, ttl)
	ret0, _ := ret[0].(error)
	return ret0
}

// Set indicates an expected call of Set
func (mr *MockRedisWrapMockRecorder) Set(key, data, ttl interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockRedisWrap)(nil).Set), key, data, ttl)
}

// Del mocks base method
func (m *MockRedisWrap) Del(key string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Del", key)
	ret0, _ := ret[0].(error)
	return ret0
}

// Del indicates an expected call of Del
func (mr *MockRedisWrapMockRecorder) Del(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Del", reflect.TypeOf((*MockRedisWrap)(nil).Del), key)
}

// Get mocks base method
func (m *MockRedisWrap) Get(key string) (interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", key)
	ret0, _ := ret[0].(interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockRedisWrapMockRecorder) Get(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockRedisWrap)(nil).Get), key)
}

// Close mocks base method
func (m *MockRedisWrap) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockRedisWrapMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockRedisWrap)(nil).Close))
}
