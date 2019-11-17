package mocks

import (
	"github.com/coldze/test/middleware"
	"github.com/golang/mock/gomock"
	"reflect"
)

// MockRequestDataExtractor is a mock of DataParser interface
type MockRequestDataExtractor struct {
	ctrl     *gomock.Controller
	recorder *MockRequestDataExtractorMockRecorder
}

// MockRequestDataExtractorMockRecorder is the mock recorder for MockRequestDataExtractor
type MockRequestDataExtractorMockRecorder struct {
	mock *MockRequestDataExtractor
}

func NewMockRequestDataExtractor(ctrl *gomock.Controller) *MockRequestDataExtractor {
	mock := &MockRequestDataExtractor{ctrl: ctrl}
	mock.recorder = &MockRequestDataExtractorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRequestDataExtractor) EXPECT() *MockRequestDataExtractorMockRecorder {
	return m.recorder
}

// Create mocks base method
func (m *MockRequestDataExtractor) Extract(arg0 middleware.HttpRequest) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Extract", arg0)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create
func (mr *MockRequestDataExtractorMockRecorder) Extract(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Extract", reflect.TypeOf((*MockRequestDataExtractor)(nil).Extract), arg0)
}
