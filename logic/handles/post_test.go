package handles

import (
	"github.com/coldze/test/mocks/mock_handles"
	"github.com/coldze/test/mocks/mock_sources"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestNewPostHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	loggerFactory := mock_handles.NewMockLoggerFactory(ctrl)
	dataSource := mock_sources.NewMockDataSource(ctrl)

	handler := NewPostHandler(loggerFactory.Create, dataSource)
	if handler == nil {
		t.Errorf("Post handler is nil")
	}
}
