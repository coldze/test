package mocks

import (
	"github.com/coldze/test/middleware"
	"github.com/golang/mock/gomock"
	"reflect"
)

// MockDataBuilder is a mock of DataBuilder interface
type MockResponseFactory struct {
	ctrl     *gomock.Controller
	recorder *MockResponseFactoryMockRecorder
}

// MockDataBuilderMockRecorder is the mock recorder for MockDataBuilder
type MockResponseFactoryMockRecorder struct {
	mock *MockResponseFactory
}

func NewMockResponseFactory(ctrl *gomock.Controller) *MockResponseFactory {
	mock := &MockResponseFactory{ctrl: ctrl}
	mock.recorder = &MockResponseFactoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockResponseFactory) EXPECT() *MockResponseFactoryMockRecorder {
	return m.recorder
}

// Do mocks base method
func (m *MockResponseFactory) Create(arg0 []byte) (middleware.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0)
	ret0, _ := ret[0].(middleware.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Do indicates an expected call of Do
func (mr *MockResponseFactoryMockRecorder) Create(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockResponseFactory)(nil).Create), arg0)
}
