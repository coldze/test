package mock_handles

import (
	"github.com/coldze/test/logs"
	"github.com/golang/mock/gomock"
	"reflect"
)

// MockHttpHandler is a mock of http.HandlerFunc interface
type MockLoggerFactory struct {
	ctrl     *gomock.Controller
	recorder *MockLoggerFactoryMockRecorder
}

// MockHttpHandlerMockRecorder is the mock recorder for MockHttpHandler
type MockLoggerFactoryMockRecorder struct {
	mock *MockLoggerFactory
}

func NewMockLoggerFactory(ctrl *gomock.Controller) *MockLoggerFactory {
	mock := &MockLoggerFactory{ctrl: ctrl}
	mock.recorder = &MockLoggerFactoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockLoggerFactory) EXPECT() *MockLoggerFactoryMockRecorder {
	return m.recorder
}

// Do mocks base method
func (m *MockLoggerFactory) Create() logs.Logger {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create")
	ret0, _ := ret[0].(logs.Logger)
	return ret0
}

// Do indicates an expected call of Do
func (mr *MockLoggerFactoryMockRecorder) Create() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockLoggerFactory)(nil).Create))
}
