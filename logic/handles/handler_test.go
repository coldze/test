package handles

import (
	"context"
	"errors"
	"github.com/coldze/test/mocks/mock_handles"
	"github.com/coldze/test/mocks/mock_logic"
	"github.com/coldze/test/mocks/mock_logs"
	"github.com/coldze/test/mocks/mock_std"

	"github.com/coldze/test/mocks"
	"github.com/coldze/test/utils"
	"github.com/golang/mock/gomock"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func makeLoggerContext(controller *gomock.Controller) (context.Context, *mock_logs.MockLogger) {
	logger := mock_logs.NewMockLogger(controller)
	return utils.SetLogger(context.Background(), logger), logger
}

type handlerFixture struct {
	Url           string
	Ctx           context.Context
	CtxWithHeader context.Context
	Data          string
	Error         error
	Response      *mocks.MockResponse
	Logger        *mock_logs.MockLogger
	GetData       *mock_logic.MockRequestDataExtractor
	LogicHandler  *mock_handles.MockLogicHandler
	W             *mock_std.MockResponseWriter
}

func newHandlerFixture(ctrl *gomock.Controller) *handlerFixture {
	logger := mock_logs.NewMockLogger(ctrl)
	ctx := utils.SetLogger(context.Background(), logger)
	ctxWithHeader := utils.SetHeaders(ctx, http.Header{})

	return &handlerFixture{
		Url:           "https://test.com.au/",
		Ctx:           ctx,
		CtxWithHeader: ctxWithHeader,
		Data:          "some random data",
		Error:         errors.New("Some test error"),
		Logger:        logger,
		Response:      mocks.NewMockResponse(ctrl),
		LogicHandler:  mock_handles.NewMockLogicHandler(ctrl),
		GetData:       mock_logic.NewMockRequestDataExtractor(ctrl),
		W:             mock_std.NewMockResponseWriter(ctrl),
	}
}

func newRequest(ctx context.Context, method string, url string, body io.Reader, headers http.Header) *http.Request {
	r := httptest.NewRequest(method, url, body)
	r.Header = headers
	return r.WithContext(ctx)
}

func newTestableHttpHandler(f *handlerFixture) http.HandlerFunc {
	return newHttpHandler(f.GetData.Extract, f.LogicHandler.Handle)
}

func TestHttpHandler(t *testing.T) {
	t.Run("error while getting data is a failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newHandlerFixture(ctrl)
		httpHandler := newTestableHttpHandler(f)
		r := httptest.NewRequest(http.MethodGet, f.Url, nil)
		r = r.WithContext(f.Ctx)

		f.GetData.EXPECT().Extract(r).Return(nil, f.Error).Times(1)
		f.Logger.EXPECT().Errorf(gomock.Any(), f.Error).Times(1)
		f.W.EXPECT().WriteHeader(http.StatusInternalServerError)

		httpHandler(f.W, r)
	})

	t.Run("handling error is a failure if response is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newHandlerFixture(ctrl)
		httpHandler := newTestableHttpHandler(f)

		r := newRequest(f.Ctx, http.MethodGet, f.Url, nil, http.Header{})

		f.GetData.EXPECT().Extract(r).Return([]byte(f.Data), nil).Times(1)
		f.LogicHandler.EXPECT().Handle(f.CtxWithHeader, []byte(f.Data)).Return(nil, f.Error).Times(1)
		f.Logger.EXPECT().Errorf(gomock.Any(), f.Error).Times(1)
		f.Logger.EXPECT().Errorf(gomock.Any()).Times(1)
		f.W.EXPECT().WriteHeader(http.StatusInternalServerError)

		httpHandler(f.W, r)
	})

	t.Run("handling error is not a failure if response is not empty, but error is logged", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newHandlerFixture(ctrl)
		httpHandler := newTestableHttpHandler(f)

		r := newRequest(f.Ctx, http.MethodGet, f.Url, nil, http.Header{})

		f.GetData.EXPECT().Extract(r).Return([]byte(f.Data), nil).Times(1)
		f.LogicHandler.EXPECT().Handle(f.CtxWithHeader, []byte(f.Data)).Return(f.Response, f.Error).Times(1)
		f.Logger.EXPECT().Errorf(gomock.Any(), f.Error).Times(1)
		f.Response.EXPECT().Write(f.W).Return(nil).Times(1)

		httpHandler(f.W, r)
	})

	t.Run("handling error is not a failure if response is not empty, but error is logged, resp write error is not a failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newHandlerFixture(ctrl)
		httpHandler := newTestableHttpHandler(f)

		expErr2 := errors.New("Test 2")

		r := newRequest(f.Ctx, http.MethodGet, f.Url, nil, http.Header{})

		f.GetData.EXPECT().Extract(r).Return([]byte(f.Data), nil).Times(1)
		f.LogicHandler.EXPECT().Handle(f.CtxWithHeader, []byte(f.Data)).Return(f.Response, f.Error).Times(1)
		f.Logger.EXPECT().Errorf(gomock.Any(), f.Error).Times(1)
		f.Response.EXPECT().Write(f.W).Return(expErr2).Times(1)
		f.Logger.EXPECT().Errorf(gomock.Any(), expErr2)

		httpHandler(f.W, r)
	})

	t.Run("resp write error is not a failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newHandlerFixture(ctrl)
		httpHandler := newTestableHttpHandler(f)

		r := newRequest(f.Ctx, http.MethodGet, f.Url, nil, http.Header{})

		f.GetData.EXPECT().Extract(r).Return([]byte(f.Data), nil).Times(1)
		f.LogicHandler.EXPECT().Handle(f.CtxWithHeader, []byte(f.Data)).Return(f.Response, nil).Times(1)
		f.Response.EXPECT().Write(f.W).Return(f.Error).Times(1)
		f.Logger.EXPECT().Errorf(gomock.Any(), f.Error)

		httpHandler(f.W, r)
	})

	t.Run("logs error on throw of error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newHandlerFixture(ctrl)
		getData := func(arg *http.Request) ([]byte, error) {
			_, _ = f.GetData.Extract(arg)
			panic(f.Error)
		}
		httpHandler := newHttpHandler(getData, f.LogicHandler.Handle)

		r := newRequest(f.Ctx, http.MethodGet, f.Url, nil, http.Header{})

		f.GetData.EXPECT().Extract(r).Return([]byte(f.Data), nil).Times(1)
		f.Logger.EXPECT().Errorf(gomock.Any(), f.Error)

		httpHandler(f.W, r)
	})

	t.Run("logs error on throw of unknown", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		unknownErr := "unknown error type"
		f := newHandlerFixture(ctrl)
		getData := func(arg *http.Request) ([]byte, error) {
			_, _ = f.GetData.Extract(arg)
			panic(unknownErr)
		}
		httpHandler := newHttpHandler(getData, f.LogicHandler.Handle)

		r := newRequest(f.Ctx, http.MethodGet, f.Url, nil, http.Header{})

		f.GetData.EXPECT().Extract(r).Return([]byte(f.Data), nil).Times(1)
		f.Logger.EXPECT().Errorf(gomock.Any(), unknownErr, unknownErr)

		httpHandler(f.W, r)
	})

	t.Run("if body not nil, it gets closed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		unknownErr := "unknown error type"
		f := newHandlerFixture(ctrl)
		getData := func(arg *http.Request) ([]byte, error) {
			_, _ = f.GetData.Extract(arg)
			panic(unknownErr)
		}
		httpHandler := newHttpHandler(getData, f.LogicHandler.Handle)
		body := mock_std.NewMockReadCloser(ctrl)

		r := newRequest(f.Ctx, http.MethodGet, f.Url, body, http.Header{})

		body.EXPECT().Close().Return(nil).Times(1)
		f.GetData.EXPECT().Extract(r).Return([]byte(f.Data), nil).Times(1)
		f.Logger.EXPECT().Errorf(gomock.Any(), unknownErr, unknownErr)

		httpHandler(f.W, r)
	})

	t.Run("if body returns error when closed, it gets logged", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		unknownErr := "unknown error type"
		f := newHandlerFixture(ctrl)
		getData := func(arg *http.Request) ([]byte, error) {
			_, _ = f.GetData.Extract(arg)
			panic(unknownErr)
		}
		httpHandler := newHttpHandler(getData, f.LogicHandler.Handle)
		body := mock_std.NewMockReadCloser(ctrl)

		r := newRequest(f.Ctx, http.MethodGet, f.Url, body, http.Header{})

		body.EXPECT().Close().Return(f.Error).Times(1)
		f.Logger.EXPECT().Errorf(gomock.Any(), f.Error).Times(1)
		f.GetData.EXPECT().Extract(r).Return([]byte(f.Data), nil).Times(1)
		f.Logger.EXPECT().Errorf(gomock.Any(), unknownErr, unknownErr)

		httpHandler(f.W, r)
	})

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newHandlerFixture(ctrl)
		httpHandler := newTestableHttpHandler(f)

		r := newRequest(f.Ctx, http.MethodGet, f.Url, nil, http.Header{})

		f.GetData.EXPECT().Extract(r).Return([]byte(f.Data), nil).Times(1)
		f.LogicHandler.EXPECT().Handle(f.CtxWithHeader, []byte(f.Data)).Return(f.Response, nil).Times(1)
		f.Response.EXPECT().Write(f.W).Return(nil).Times(1)

		httpHandler(f.W, r)
	})
}

func TestNewHttpHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	f := newHandlerFixture(ctrl)
	httpHandler := newHttpHandler(f.GetData.Extract, f.LogicHandler.Handle)
	if httpHandler == nil {
		t.Errorf("Factory returns nil")
	}
}
