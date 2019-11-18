package handles

import (
	"github.com/coldze/test/mocks/mock_handles"
	"github.com/coldze/test/mocks/mock_sources"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestNewPutHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	loggerFactory := mock_handles.NewMockLoggerFactory(ctrl)
	dataSource := mock_sources.NewMockDataSource(ctrl)

	handler := NewPutHandler(loggerFactory.Create, dataSource)
	if handler == nil {
		t.Errorf("Put handler is nil")
	}
}
