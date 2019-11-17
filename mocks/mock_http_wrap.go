package mocks

import (
	"github.com/golang/mock/gomock"
	"net/http"
	"reflect"
)

// MockDataBuilder is a mock of DataBuilder interface
type MockHttpWrap struct {
	ctrl     *gomock.Controller
	recorder *MockHttpWrapMockRecorder
}

// MockDataBuilderMockRecorder is the mock recorder for MockDataBuilder
type MockHttpWrapMockRecorder struct {
	mock *MockHttpWrap
}

func NewMockHttpWrap(ctrl *gomock.Controller) *MockHttpWrap {
	mock := &MockHttpWrap{ctrl: ctrl}
	mock.recorder = &MockHttpWrapMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockHttpWrap) EXPECT() *MockHttpWrapMockRecorder {
	return m.recorder
}

// Do mocks base method
func (m *MockHttpWrap) Do(arg0 *http.Request) (*http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Do", arg0)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Do indicates an expected call of Do
func (mr *MockHttpWrapMockRecorder) Do(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Do", reflect.TypeOf((*MockHttpWrap)(nil).Do), arg0)
}
