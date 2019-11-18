package logic

import (
	"bytes"
	"errors"
	"github.com/coldze/test/mocks/mock_std"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type wrapIoReadCloser struct {
	data io.Reader
	sub  *mock_std.MockReadCloser
}

func (m *wrapIoReadCloser) Close() error {
	return m.sub.Close()
}

func (m *wrapIoReadCloser) Read(arg0 []byte) (int, error) {
	res, err := m.sub.Read(arg0)
	_, _ = m.data.Read(arg0)
	return res, err
}

func TestGetBodyData(t *testing.T) {
	t.Run("returns empty data if reader is nil", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		data := []byte{}

		res, err := getBodyData(nil)
		if err != nil {
			t.Errorf("Empty reader is fine, but got error: %v", err)
		}
		if !cmp.Equal(res, data) {
			t.Errorf("Expected empty data. Received: %v", res)
		}
	})

	t.Run("returns read data", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		data := []byte{1, 2}

		b := mock_std.NewMockReadCloser(ctrl)

		w := &wrapIoReadCloser{
			data: bytes.NewReader(data),
			sub:  b,
		}

		b.EXPECT().Read(gomock.Any()).Return(len(data), io.EOF).Times(1)
		res, err := getBodyData(w)
		if err != nil {
			t.Errorf("Empty reader is fine, but got error: %v", err)
		}
		if !cmp.Equal(res, data) {
			t.Errorf("Expected %v. Received: %v", data, res)
		}
	})

	t.Run("returns error if failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		data := []byte{1, 2}

		b := mock_std.NewMockReadCloser(ctrl)

		w := &wrapIoReadCloser{
			data: bytes.NewReader(data),
			sub:  b,
		}

		expErr := errors.New("failed")

		b.EXPECT().Read(gomock.Any()).Return(len(data), expErr).Times(1)
		res, err := getBodyData(w)
		if err == nil {
			t.Errorf("Empty reader is fine, but got error: %v", err)
		}
		if !cmp.Equal(res, data) {
			t.Errorf("Expected %v. Received: %v", data, res)
		}
	})
}

func TestGetRequestBodyData(t *testing.T) {
	t.Run("returns empty data if request is nil", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		data := []byte{}

		res, err := GetRequestBodyData(nil)
		if err != nil {
			t.Errorf("Empty reader is fine, but got error: %v", err)
		}
		if !cmp.Equal(res, data) {
			t.Errorf("Expected empty data. Received: %v", res)
		}
	})

	t.Run("returns empty data if body is nil", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		data := []byte{}
		r := httptest.NewRequest(http.MethodGet, "https://test.com.au", nil)

		res, err := GetRequestBodyData(r)
		if err != nil {
			t.Errorf("Empty reader is fine, but got error: %v", err)
		}
		if !cmp.Equal(res, data) {
			t.Errorf("Expected empty data. Received: %v", res)
		}
	})

	t.Run("returns read data", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		data := []byte{1, 2}

		b := mock_std.NewMockReadCloser(ctrl)

		w := &wrapIoReadCloser{
			data: bytes.NewReader(data),
			sub:  b,
		}

		r := httptest.NewRequest(http.MethodGet, "https://test.com.au", w)

		b.EXPECT().Read(gomock.Any()).Return(len(data), io.EOF).Times(1)
		res, err := GetRequestBodyData(r)
		if err != nil {
			t.Errorf("Empty reader is fine, but got error: %v", err)
		}
		if !cmp.Equal(res, data) {
			t.Errorf("Expected %v. Received: %v", data, res)
		}
	})

	t.Run("returns error if failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		data := []byte{1, 2}

		b := mock_std.NewMockReadCloser(ctrl)

		w := &wrapIoReadCloser{
			data: bytes.NewReader(data),
			sub:  b,
		}

		expErr := errors.New("failed")

		r := httptest.NewRequest(http.MethodGet, "https://test.com.au", w)

		b.EXPECT().Read(gomock.Any()).Return(len(data), expErr).Times(1)
		res, err := GetRequestBodyData(r)
		if err == nil {
			t.Errorf("Empty reader is fine, but got error: %v", err)
		}
		if !cmp.Equal(res, data) {
			t.Errorf("Expected %v. Received: %v", data, res)
		}
	})
}

func TestGetResponseBodyData(t *testing.T) {
	t.Run("returns empty data if request is nil", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		data := []byte{}

		res, err := GetResponseBodyData(nil)
		if err != nil {
			t.Errorf("Empty reader is fine, but got error: %v", err)
		}
		if !cmp.Equal(res, data) {
			t.Errorf("Expected empty data. Received: %v", res)
		}
	})

	t.Run("returns empty data if body is nil", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		data := []byte{}
		rec := httptest.NewRecorder()
		resp := rec.Result()
		err := resp.Body.Close()
		if err != nil {
			panic(err)
		}
		resp.Body = nil

		res, err := GetResponseBodyData(resp)
		if err != nil {
			t.Errorf("Empty reader is fine, but got error: %v", err)
		}
		if !cmp.Equal(res, data) {
			t.Errorf("Expected empty data. Received: %v", res)
		}
	})

	t.Run("returns read data", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		data := []byte{1, 2}

		b := mock_std.NewMockReadCloser(ctrl)

		rec := httptest.NewRecorder()
		_, err := rec.Write(data)
		if err != nil {
			panic(err)
		}
		resp := rec.Result()
		resp.Body = &wrapIoReadCloser{
			data: resp.Body,
			sub:  b,
		}

		b.EXPECT().Read(gomock.Any()).Return(len(data), io.EOF).Times(1)
		res, err := GetResponseBodyData(rec.Result())
		if err != nil {
			t.Errorf("Empty reader is fine, but got error: %v", err)
		}
		if !cmp.Equal(res, data) {
			t.Errorf("Expected %v. Received: %v", data, res)
		}
	})

	t.Run("returns error if failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		data := []byte{1, 2}

		b := mock_std.NewMockReadCloser(ctrl)

		expErr := errors.New("failed")

		rec := httptest.NewRecorder()
		_, err := rec.Write(data)
		if err != nil {
			panic(err)
		}

		resp := rec.Result()
		resp.Body = &wrapIoReadCloser{
			data: resp.Body,
			sub:  b,
		}

		b.EXPECT().Read(gomock.Any()).Return(len(data), expErr).Times(1)
		res, err := GetResponseBodyData(rec.Result())
		if err == nil {
			t.Errorf("Empty reader is fine, but got error: %v", err)
		}
		if !cmp.Equal(res, data) {
			t.Errorf("Expected %v. Received: %v", data, res)
		}
	})
}
