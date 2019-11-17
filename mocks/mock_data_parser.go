package mocks

import (
	"github.com/coldze/test/middleware"
	"github.com/golang/mock/gomock"
	"reflect"
)

// MockDataParser is a mock of DataParser interface
type MockDataParser struct {
	ctrl     *gomock.Controller
	recorder *MockDataParserMockRecorder
}

// MockDataParserMockRecorder is the mock recorder for MockDataParser
type MockDataParserMockRecorder struct {
	mock *MockDataParser
}

func NewMockDataParser(ctrl *gomock.Controller) *MockDataParser {
	mock := &MockDataParser{ctrl: ctrl}
	mock.recorder = &MockDataParserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDataParser) EXPECT() *MockDataParserMockRecorder {
	return m.recorder
}

// Create mocks base method
func (m *MockDataParser) Parse(arg0 []byte) (middleware.Contact, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Parse", arg0)
	ret0, _ := ret[0].(middleware.Contact)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create
func (mr *MockDataParserMockRecorder) Create(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Parse", reflect.TypeOf((*MockDataParser)(nil).Parse), arg0)
}
