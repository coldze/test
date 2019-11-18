package logic

import (
	"errors"
	"github.com/coldze/test/mocks"
	"github.com/coldze/test/mocks/mock_logic"
	"github.com/coldze/test/mocks/mock_std"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHttpResponse_Write(t *testing.T) {
	t.Run("empty header is not written", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		w := mock_std.NewMockResponseWriter(ctrl)
		r := httpResponse{
			data:    nil,
			headers: http.Header{},
			code:    http.StatusInternalServerError,
		}
		w.EXPECT().WriteHeader(r.code).Times(1)
		w.EXPECT().Write(r.data).Return(len(r.data), nil).Times(1)
		err := r.Write(w)
		if err != nil {
			t.Errorf("No errors expected. Error: %v", err)
		}
	})

	t.Run("error in write is a failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		w := mock_std.NewMockResponseWriter(ctrl)
		r := httpResponse{
			data:    nil,
			headers: http.Header{},
			code:    http.StatusInternalServerError,
		}
		expErr := errors.New("some test error")
		w.EXPECT().WriteHeader(r.code).Times(1)
		w.EXPECT().Write(r.data).Return(len(r.data), expErr).Times(1)
		err := r.Write(w)
		if err != expErr {
			t.Errorf("No errors expected. Error: %v", err)
		}
	})

	t.Run("empty headers set tries to be written", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		w := mock_std.NewMockResponseWriter(ctrl)
		r := httpResponse{
			data: nil,
			headers: http.Header{
				"Test": []string{},
			},
			code: http.StatusInternalServerError,
		}
		expErr := errors.New("some test error")
		w.EXPECT().WriteHeader(r.code).Times(1)
		w.EXPECT().Write(r.data).Return(len(r.data), expErr).Times(1)
		err := r.Write(w)
		if err != expErr {
			t.Errorf("No errors expected. Error: %v", err)
		}
	})

	t.Run("headers are written", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		w := mock_std.NewMockResponseWriter(ctrl)
		r := httpResponse{
			data: nil,
			headers: http.Header{
				"Test":  []string{"1", "2"},
				"Test2": []string{"4"},
			},
			code: http.StatusInternalServerError,
		}
		cnt := 0
		for _, v := range r.headers {
			cnt += len(v)
		}
		headers := http.Header{}
		expErr := errors.New("some test error")
		w.EXPECT().Header().Return(headers).Times(cnt)
		w.EXPECT().WriteHeader(r.code).Times(1)
		w.EXPECT().Write(r.data).Return(len(r.data), expErr).Times(1)
		err := r.Write(w)
		if err != expErr {
			t.Errorf("No errors expected. Error: %v", err)
		}
		if !cmp.Equal(headers, r.headers) {
			t.Errorf("Expected: %+v. Got: %+v", r.headers, headers)
		}
	})
}

func TestNewHttpResponse(t *testing.T) {
	t.Run("factory works", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		data := []byte{1, 2}
		headers := http.Header{
			"Test":  []string{"1", "2"},
			"Test2": []string{"4"},
		}
		code := http.StatusInternalServerError
		r, err := NewHttpResponse(data, headers, code)

		if r == nil {
			t.Errorf("Factory returned nil value")
		}
		if err != nil {
			t.Errorf("No errors expected. Got: %v", err)
		}
	})
}

func TestNewJsonOkResponse(t *testing.T) {
	t.Run("factory works", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		data := []byte{1, 2}
		r, err := NewJsonOkResponse(data)

		if r == nil {
			t.Errorf("Factory returned nil value")
		}
		if err != nil {
			t.Errorf("No errors expected. Got: %v", err)
		}
	})
}

func TestNewDefaultHttpResponseFactory(t *testing.T) {
	t.Run("factory works", func(t *testing.T) {
		r := NewDefaultHttpResponseFactory()

		if r == nil {
			t.Errorf("Factory returned nil value")
		}
	})
}

func TestNewHttpResponseFactory(t *testing.T) {
	t.Run("factory works", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		getData := mock_logic.NewMockResponseDataExtractor(ctrl)
		r := NewHttpResponseFactory(getData.Extract)

		if r == nil {
			t.Errorf("Factory returned nil value")
		}
	})

	t.Run("it is a failure if getData has error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		getData := mock_logic.NewMockResponseDataExtractor(ctrl)
		responseFactory := NewHttpResponseFactory(getData.Extract)
		rec := httptest.NewRecorder()
		rec.Header().Set("1", "2")
		_, err := rec.Write([]byte{1, 2, 3})
		if err != nil {
			panic(err)
		}

		if responseFactory == nil {
			t.Errorf("Factory returned nil value")
		}

		expErr := errors.New("some test error")

		response := rec.Result()
		getData.EXPECT().Extract(response).Return(nil, expErr).Times(1)
		res, err := responseFactory(response)
		mocks.CmpError(t, err, expErr)
		if res != nil {
			t.Errorf("Expected nil response. Got: %+v. Type: %T", res, res)
		}
	})

	t.Run("succeeds", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		getData := mock_logic.NewMockResponseDataExtractor(ctrl)
		responseFactory := NewHttpResponseFactory(getData.Extract)
		data := []byte{1, 2, 3}
		rec := httptest.NewRecorder()
		rec.Header().Set("1", "2")
		_, err := rec.Write(data)
		if err != nil {
			panic(err)
		}

		if responseFactory == nil {
			t.Errorf("Factory returned nil value")
		}

		response := rec.Result()
		getData.EXPECT().Extract(response).Return(data, nil).Times(1)
		res, err := responseFactory(response)
		mocks.CmpError(t, err, nil)
		if res == nil {
			t.Errorf("Expected nil response. Got: %+v. Type: %T", res, res)
		}
	})
}
