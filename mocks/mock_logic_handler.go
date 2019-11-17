package mocks

import (
	"context"
	"github.com/coldze/test/middleware"
	"github.com/golang/mock/gomock"
	"reflect"
)

// MockLogicHandler is a mock of DataParser interface
type MockLogicHandler struct {
	ctrl     *gomock.Controller
	recorder *MockLogicHandlerMockRecorder
}

// MockLogicHandlerMockRecorder is the mock recorder for MockLogicHandler
type MockLogicHandlerMockRecorder struct {
	mock *MockLogicHandler
}

func NewMockLogicHandler(ctrl *gomock.Controller) *MockLogicHandler {
	mock := &MockLogicHandler{ctrl: ctrl}
	mock.recorder = &MockLogicHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockLogicHandler) EXPECT() *MockLogicHandlerMockRecorder {
	return m.recorder
}

// Create mocks base method
func (m *MockLogicHandler) Handle(arg0 context.Context, arg1 []byte) (middleware.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Handle", arg0, arg1)
	ret0, _ := ret[0].(middleware.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create
func (mr *MockLogicHandlerMockRecorder) Handle(arg0 interface{}, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Handle", reflect.TypeOf((*MockLogicHandler)(nil).Handle), arg0, arg1)
}
