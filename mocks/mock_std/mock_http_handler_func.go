package mock_std

import (
	"github.com/golang/mock/gomock"
	"net/http"
	"reflect"
)

// MockHttpHandler is a mock of http.HandlerFunc interface
type MockHttpHandlerFunc struct {
	ctrl     *gomock.Controller
	recorder *MockHttpHandlerFuncMockRecorder
}

// MockHttpHandlerMockRecorder is the mock recorder for MockHttpHandler
type MockHttpHandlerFuncMockRecorder struct {
	mock *MockHttpHandlerFunc
}

func NewMockHttpHandlerFunc(ctrl *gomock.Controller) *MockHttpHandlerFunc {
	mock := &MockHttpHandlerFunc{ctrl: ctrl}
	mock.recorder = &MockHttpHandlerFuncMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockHttpHandlerFunc) EXPECT() *MockHttpHandlerFuncMockRecorder {
	return m.recorder
}

// Do mocks base method
func (m *MockHttpHandlerFunc) Handle(arg0 http.ResponseWriter, arg1 *http.Request) {
	m.ctrl.T.Helper()
	_ = m.ctrl.Call(m, "Handle", arg0, arg1)
}

// Do indicates an expected call of Do
func (mr *MockHttpHandlerFuncMockRecorder) Create(arg0 interface{}, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Handle", reflect.TypeOf((*MockHttpHandlerFunc)(nil).Handle), arg0, arg1)
}
