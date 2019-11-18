package handles

import (
	"github.com/coldze/test/mocks/mock_handles"
	"github.com/coldze/test/mocks/mock_logs"
	"github.com/coldze/test/mocks/mock_std"
	"github.com/coldze/test/utils"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

type middlewareFixture struct {
	Url       string
	W         *mock_std.MockResponseWriter
	NewLogger *mock_handles.MockLoggerFactory
	Logger    *mock_logs.MockLogger
	Next      *mock_std.MockHttpHandlerFunc
	Reader    *mock_std.MockReadCloser
}

func newMiddlewareFixture(ctrl *gomock.Controller) *middlewareFixture {
	return &middlewareFixture{
		Url:       "https://test.url.com/",
		W:         mock_std.NewMockResponseWriter(ctrl),
		NewLogger: mock_handles.NewMockLoggerFactory(ctrl),
		Logger:    mock_logs.NewMockLogger(ctrl),
		Next:      mock_std.NewMockHttpHandlerFunc(ctrl),
		Reader:    mock_std.NewMockReadCloser(ctrl),
	}
}

func TestHttpMiddleware(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newMiddlewareFixture(ctrl)

		next := func(w http.ResponseWriter, r *http.Request) {
			f.Next.Handle(w, r)
			logger := utils.GetLogger(r.Context())
			if logger != f.Logger {
				t.Errorf("Not expected logger value: %v", logger)
			}
		}

		handle := newCheckAndSetLoggerMiddleware(f.NewLogger.Create, next)
		r := httptest.NewRequest(http.MethodGet, "https://test.url.com/", f.Reader)

		f.NewLogger.EXPECT().Create().Return(f.Logger).Times(1)
		f.Next.EXPECT().Create(f.W, gomock.Any()).Times(1)

		handle(f.W, r)
	})

	t.Run("failed if request is nil", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newMiddlewareFixture(ctrl)

		next := func(w http.ResponseWriter, r *http.Request) {
			f.Next.Handle(w, r)
			logger := utils.GetLogger(r.Context())
			if logger != f.Logger {
				t.Errorf("Not expected logger value: %v", logger)
			}
		}

		handle := newCheckAndSetLoggerMiddleware(f.NewLogger.Create, next)

		f.W.EXPECT().WriteHeader(http.StatusInternalServerError).Times(1)

		handle(f.W, nil)
	})
}
