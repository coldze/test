package handles

import (
	"github.com/coldze/test/mocks/mock_logs"
	"github.com/golang/mock/gomock"
	"testing"
)

func TestDefaultLoggerFactory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mock_logs.NewMockLogger(ctrl)
	lf := NewDefaultLoggerFactory(logger)
	res := lf()
	if res != logger {
		t.Errorf("Logger was modified")
	}
}
