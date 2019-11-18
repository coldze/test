package mock_sources

import (
	"github.com/coldze/test/logic"
	"github.com/golang/mock/gomock"
	"reflect"
)

// MockDataBuilderFactory is a mock of DataBuilderFactory interface
type MockDataBuilderFactory struct {
	ctrl     *gomock.Controller
	recorder *MockDataBuilderFactoryMockRecorder
}

// MockDataBuilderFactoryMockRecorder is the mock recorder for MockDataBuilderFactory
type MockDataBuilderFactoryMockRecorder struct {
	mock *MockDataBuilderFactory
}

func NewMockDataBuilderFactory(ctrl *gomock.Controller) *MockDataBuilderFactory {
	mock := &MockDataBuilderFactory{ctrl: ctrl}
	mock.recorder = &MockDataBuilderFactoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDataBuilderFactory) EXPECT() *MockDataBuilderFactoryMockRecorder {
	return m.recorder
}

// Create mocks base method
func (m *MockDataBuilderFactory) Create() logic.DataBuilder {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create")
	ret0, _ := ret[0].(logic.DataBuilder)
	return ret0
}

// Create indicates an expected call of Create
func (mr *MockDataBuilderFactoryMockRecorder) Create() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockDataBuilderFactory)(nil).Create))
}
