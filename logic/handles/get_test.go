package handles

import (
	"github.com/coldze/test/mocks/mock_handles"
	"github.com/coldze/test/mocks/mock_logic"
	"github.com/coldze/test/mocks/mock_sources"
	"github.com/golang/mock/gomock"
	"testing"
)

func TestNewGetHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	loggerFactory := mock_handles.NewMockLoggerFactory(ctrl)
	dataSource := mock_sources.NewMockDataSource(ctrl)
	getData := mock_logic.NewMockRequestDataExtractor(ctrl)

	handler := NewGetHandler(loggerFactory.Create, dataSource, getData.Extract)
	if handler == nil {
		t.Errorf("Get handler is nil")
	}
}
