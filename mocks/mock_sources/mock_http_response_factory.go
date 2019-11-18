package mock_sources

import (
	"github.com/coldze/test/logic"
	"github.com/golang/mock/gomock"
	"net/http"
	"reflect"
)

// MockDataBuilder is a mock of DataBuilder interface
type MockHttpResponseFactory struct {
	ctrl     *gomock.Controller
	recorder *MockHttpResponseFactoryMockRecorder
}

// MockDataBuilderMockRecorder is the mock recorder for MockDataBuilder
type MockHttpResponseFactoryMockRecorder struct {
	mock *MockHttpResponseFactory
}

func NewMockHttpResponseFactory(ctrl *gomock.Controller) *MockHttpResponseFactory {
	mock := &MockHttpResponseFactory{ctrl: ctrl}
	mock.recorder = &MockHttpResponseFactoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockHttpResponseFactory) EXPECT() *MockHttpResponseFactoryMockRecorder {
	return m.recorder
}

// Do mocks base method
func (m *MockHttpResponseFactory) Create(arg0 *http.Response) (logic.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0)
	ret0, _ := ret[0].(logic.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Do indicates an expected call of Do
func (mr *MockHttpResponseFactoryMockRecorder) Create(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockHttpResponseFactory)(nil).Create), arg0)
}
