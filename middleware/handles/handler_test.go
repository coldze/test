package handles

import (
	"context"
	"errors"
	"github.com/coldze/test/mocks"
	"github.com/coldze/test/utils"
	"github.com/golang/mock/gomock"
	"net/http"
	"testing"
)

func cmperr(t *testing.T, err error, expErr error) {
	if err != expErr {
		t.Helper()
		t.Errorf("Unexpected error: (%T)%v", err, err)
		t.FailNow()
	}
}

func makeLoggerContext(controller *gomock.Controller) (context.Context, *mocks.MockLogger) {
	logger := mocks.NewMockLogger(controller)
	return utils.SetLogger(context.Background(), logger), logger
}


//TODO: I'm missing tests for panic handling. Will add them later this evening
func TestHttpHandler(t *testing.T) {
	data := "some random data"

	t.Run("error while getting data is a failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		getData := mocks.NewMockRequestDataExtractor(ctrl)
		handler := mocks.NewMockLogicHandler(ctrl)
		httpHandler := newWrappedHttpHandler(getData.Extract, handler.Handle)
		w := mocks.NewMockResponseWriter(ctrl)
		r := mocks.NewMockHttpRequest(ctrl)

		expErr := errors.New("Test")

		ctx, logger := makeLoggerContext(ctrl)
		r.EXPECT().GetContext().Return(ctx).Times(1)
		r.EXPECT().GetBody().Return(nil).Times(1)
		getData.EXPECT().Extract(r).Return(nil, expErr).Times(1)
		logger.EXPECT().Errorf(gomock.Any(), expErr).Times(1)
		w.EXPECT().WriteHeader(http.StatusInternalServerError)

		httpHandler(w, r)
	})

	t.Run("handling error is a failure if response is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		getData := mocks.NewMockRequestDataExtractor(ctrl)
		handler := mocks.NewMockLogicHandler(ctrl)
		httpHandler := newWrappedHttpHandler(getData.Extract, handler.Handle)
		w := mocks.NewMockResponseWriter(ctrl)
		r := mocks.NewMockHttpRequest(ctrl)

		expErr := errors.New("Test")

		ctx, logger := makeLoggerContext(ctrl)
		modCtx := utils.SetHeaders(ctx, http.Header{})
		r.EXPECT().GetContext().Return(ctx).Times(1)
		r.EXPECT().GetBody().Return(nil).Times(1)
		r.EXPECT().GetHeader().Return(http.Header{})
		getData.EXPECT().Extract(r).Return([]byte(data), nil).Times(1)
		handler.EXPECT().Handle(modCtx, []byte(data)).Return(nil, expErr).Times(1)
		logger.EXPECT().Errorf(gomock.Any(), expErr).Times(1)
		logger.EXPECT().Errorf(gomock.Any()).Times(1)
		w.EXPECT().WriteHeader(http.StatusInternalServerError)

		httpHandler(w, r)
	})

	t.Run("handling error is not a failure if response is not empty, but error is logged", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		resp := mocks.NewMockResponse(ctrl)
		getData := mocks.NewMockRequestDataExtractor(ctrl)
		handler := mocks.NewMockLogicHandler(ctrl)
		httpHandler := newWrappedHttpHandler(getData.Extract, handler.Handle)
		w := mocks.NewMockResponseWriter(ctrl)
		r := mocks.NewMockHttpRequest(ctrl)

		expErr := errors.New("Test")

		ctx, logger := makeLoggerContext(ctrl)
		modCtx := utils.SetHeaders(ctx, http.Header{})
		r.EXPECT().GetContext().Return(ctx).Times(1)
		r.EXPECT().GetBody().Return(nil).Times(1)
		r.EXPECT().GetHeader().Return(http.Header{})
		getData.EXPECT().Extract(r).Return([]byte(data), nil).Times(1)
		handler.EXPECT().Handle(modCtx, []byte(data)).Return(resp, expErr).Times(1)
		logger.EXPECT().Errorf(gomock.Any(), expErr).Times(1)
		resp.EXPECT().Write(w).Return(nil).Times(1)

		httpHandler(w, r)
	})

	t.Run("handling error is not a failure if response is not empty, but error is logged, resp write error is not a failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		resp := mocks.NewMockResponse(ctrl)
		getData := mocks.NewMockRequestDataExtractor(ctrl)
		handler := mocks.NewMockLogicHandler(ctrl)
		httpHandler := newWrappedHttpHandler(getData.Extract, handler.Handle)
		w := mocks.NewMockResponseWriter(ctrl)
		r := mocks.NewMockHttpRequest(ctrl)

		expErr := errors.New("Test")
		expErr2 := errors.New("Test 2")

		ctx, logger := makeLoggerContext(ctrl)
		modCtx := utils.SetHeaders(ctx, http.Header{})
		r.EXPECT().GetContext().Return(ctx).Times(1)
		r.EXPECT().GetBody().Return(nil).Times(1)
		r.EXPECT().GetHeader().Return(http.Header{})
		getData.EXPECT().Extract(r).Return([]byte(data), nil).Times(1)
		handler.EXPECT().Handle(modCtx, []byte(data)).Return(resp, expErr).Times(1)
		logger.EXPECT().Errorf(gomock.Any(), expErr).Times(1)
		resp.EXPECT().Write(w).Return(expErr2).Times(1)
		logger.EXPECT().Errorf(gomock.Any(), expErr2)

		httpHandler(w, r)
	})

	t.Run("resp write error is not a failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		resp := mocks.NewMockResponse(ctrl)
		getData := mocks.NewMockRequestDataExtractor(ctrl)
		handler := mocks.NewMockLogicHandler(ctrl)
		httpHandler := newWrappedHttpHandler(getData.Extract, handler.Handle)
		w := mocks.NewMockResponseWriter(ctrl)
		r := mocks.NewMockHttpRequest(ctrl)

		expErr := errors.New("Test")

		ctx, logger := makeLoggerContext(ctrl)
		modCtx := utils.SetHeaders(ctx, http.Header{})
		r.EXPECT().GetContext().Return(ctx).Times(1)
		r.EXPECT().GetBody().Return(nil).Times(1)
		r.EXPECT().GetHeader().Return(http.Header{})
		getData.EXPECT().Extract(r).Return([]byte(data), nil).Times(1)
		handler.EXPECT().Handle(modCtx, []byte(data)).Return(resp, nil).Times(1)
		resp.EXPECT().Write(w).Return(expErr).Times(1)
		logger.EXPECT().Errorf(gomock.Any(), expErr)

		httpHandler(w, r)
	})

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		resp := mocks.NewMockResponse(ctrl)
		getData := mocks.NewMockRequestDataExtractor(ctrl)
		handler := mocks.NewMockLogicHandler(ctrl)
		httpHandler := newWrappedHttpHandler(getData.Extract, handler.Handle)
		w := mocks.NewMockResponseWriter(ctrl)
		r := mocks.NewMockHttpRequest(ctrl)

		ctx, _ := makeLoggerContext(ctrl)
		modCtx := utils.SetHeaders(ctx, http.Header{})
		r.EXPECT().GetContext().Return(ctx).Times(1)
		r.EXPECT().GetBody().Return(nil).Times(1)
		r.EXPECT().GetHeader().Return(http.Header{})
		getData.EXPECT().Extract(r).Return([]byte(data), nil).Times(1)
		handler.EXPECT().Handle(modCtx, []byte(data)).Return(resp, nil).Times(1)
		resp.EXPECT().Write(w).Return(nil).Times(1)

		httpHandler(w, r)
	})
}
