package mocks

import (
	"context"
	"github.com/golang/mock/gomock"
	"net/http"
	"reflect"
)

// MockDataBuilder is a mock of DataBuilder interface
type MockHttpRequestFactory struct {
	ctrl     *gomock.Controller
	recorder *MockHttpRequestFactoryMockRecorder
}

// MockDataBuilderMockRecorder is the mock recorder for MockDataBuilder
type MockHttpRequestFactoryMockRecorder struct {
	mock *MockHttpRequestFactory
}

func NewMockHttpRequestFactory(ctrl *gomock.Controller) *MockHttpRequestFactory {
	mock := &MockHttpRequestFactory{ctrl: ctrl}
	mock.recorder = &MockHttpRequestFactoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockHttpRequestFactory) EXPECT() *MockHttpRequestFactoryMockRecorder {
	return m.recorder
}

// Do mocks base method
func (m *MockHttpRequestFactory) Create(arg0 context.Context, arg1 []byte, arg2 string, arg3 string) (*http.Request, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(*http.Request)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Do indicates an expected call of Do
func (mr *MockHttpRequestFactoryMockRecorder) Create(arg0 interface{}, arg1 interface{}, arg2 interface{}, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockHttpRequestFactory)(nil).Create), arg0, arg1, arg2, arg3)
}
