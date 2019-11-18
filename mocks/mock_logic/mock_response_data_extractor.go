package mock_logic

import (
	"github.com/golang/mock/gomock"
	"net/http"
	"reflect"
)

// MockResponseDataExtractor is a mock of DataParser interface
type MockResponseDataExtractor struct {
	ctrl     *gomock.Controller
	recorder *MockResponseDataExtractorMockRecorder
}

// MockResponseDataExtractorMockRecorder is the mock recorder for MockResponseDataExtractor
type MockResponseDataExtractorMockRecorder struct {
	mock *MockResponseDataExtractor
}

func NewMockResponseDataExtractor(ctrl *gomock.Controller) *MockResponseDataExtractor {
	mock := &MockResponseDataExtractor{ctrl: ctrl}
	mock.recorder = &MockResponseDataExtractorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockResponseDataExtractor) EXPECT() *MockResponseDataExtractorMockRecorder {
	return m.recorder
}

// Create mocks base method
func (m *MockResponseDataExtractor) Extract(arg0 *http.Response) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Extract", arg0)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create
func (mr *MockResponseDataExtractorMockRecorder) Extract(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Extract", reflect.TypeOf((*MockResponseDataExtractor)(nil).Extract), arg0)
}
