package sources

import (
	"errors"
	"github.com/coldze/test/middleware"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/coldze/test/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
)

const (
	url = "https://test.url.com/v1/api"
)

func TestHttpDataSource_Get(t *testing.T) {

	target := url + "/" + key
	data := "some random data"

	t.Run("create request error is a failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockDo := mocks.NewMockHttpWrap(ctrl)
		createRequest := mocks.NewMockHttpRequestFactory(ctrl)
		createResponse := mocks.NewMockHttpResponseFactory(ctrl)

		c := httpDataSource{
			url:            url,
			do:             mockDo.Do,
			createResponse: createResponse.Create,
			createRequest:  createRequest.Create,
		}

		ctx, _ := makeLoggerContext(ctrl)
		expErr := errors.New("Test")

		createRequest.EXPECT().Create(ctx, nil, target, http.MethodGet).Return(nil, expErr).Times(1)
		r, err := c.Get(ctx, []byte(key))
		cmperr(t, err, expErr)
		if !cmp.Equal(r, nil) {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("http call error is a failure (resp not nil and closes)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		httpWrap := mocks.NewMockHttpWrap(ctrl)
		createRequest := mocks.NewMockHttpRequestFactory(ctrl)
		createResponse := mocks.NewMockHttpResponseFactory(ctrl)
		responseBody := mocks.NewMockReadCloser(ctrl)
		req := httptest.NewRequest(http.MethodGet, target, nil)
		respRec := httptest.NewRecorder()
		respRec.WriteHeader(http.StatusInternalServerError)
		_, err := respRec.WriteString(data)
		if err != nil {
			panic(err)
		}
		respRec.Flush()
		resp := respRec.Result()
		resp.Body = responseBody

		c := httpDataSource{
			url:            url,
			do:             httpWrap.Do,
			createResponse: createResponse.Create,
			createRequest:  createRequest.Create,
		}

		ctx, _ := makeLoggerContext(ctrl)
		expErr := errors.New("Test")

		createRequest.EXPECT().Create(ctx, nil, target, http.MethodGet).Return(req, nil).Times(1)
		httpWrap.EXPECT().Do(req).Return(resp, expErr).Times(1)
		responseBody.EXPECT().Close().Return(nil).Times(1)
		r, err := c.Get(ctx, []byte(key))
		cmperr(t, err, expErr)
		if !cmp.Equal(r, nil) {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("http call error is a failure (resp not nil and fails to close)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		httpWrap := mocks.NewMockHttpWrap(ctrl)
		createRequest := mocks.NewMockHttpRequestFactory(ctrl)
		createResponse := mocks.NewMockHttpResponseFactory(ctrl)
		responseBody := mocks.NewMockReadCloser(ctrl)
		req := httptest.NewRequest(http.MethodGet, target, nil)
		respRec := httptest.NewRecorder()
		respRec.WriteHeader(http.StatusInternalServerError)
		_, err := respRec.WriteString(data)
		if err != nil {
			panic(err)
		}
		respRec.Flush()
		resp := respRec.Result()
		resp.Body = responseBody

		c := httpDataSource{
			url:            url,
			do:             httpWrap.Do,
			createResponse: createResponse.Create,
			createRequest:  createRequest.Create,
		}

		ctx, logger := makeLoggerContext(ctrl)
		expErr := errors.New("Test")
		expErr2 := errors.New("Test2")

		createRequest.EXPECT().Create(ctx, nil, target, http.MethodGet).Return(req, nil).Times(1)
		httpWrap.EXPECT().Do(req).Return(resp, expErr).Times(1)
		responseBody.EXPECT().Close().Return(expErr2).Times(1)
		logger.EXPECT().Warningf(gomock.Any(), expErr2).Times(1)
		r, err := c.Get(ctx, []byte(key))
		cmperr(t, err, expErr)
		if !cmp.Equal(r, nil) {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("http call error is a failure (resp nil)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		httpWrap := mocks.NewMockHttpWrap(ctrl)
		createRequest := mocks.NewMockHttpRequestFactory(ctrl)
		createResponse := mocks.NewMockHttpResponseFactory(ctrl)
		req := httptest.NewRequest(http.MethodGet, target, nil)

		c := httpDataSource{
			url:            url,
			do:             httpWrap.Do,
			createResponse: createResponse.Create,
			createRequest:  createRequest.Create,
		}

		ctx, _ := makeLoggerContext(ctrl)
		expErr := errors.New("Test")

		createRequest.EXPECT().Create(ctx, nil, target, http.MethodGet).Return(req, nil).Times(1)
		httpWrap.EXPECT().Do(req).Return(nil, expErr).Times(1)
		r, err := c.Get(ctx, []byte(key))
		cmperr(t, err, expErr)
		if !cmp.Equal(r, nil) {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("http call error is a failure (resp body is nil)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		httpWrap := mocks.NewMockHttpWrap(ctrl)
		createRequest := mocks.NewMockHttpRequestFactory(ctrl)
		createResponse := mocks.NewMockHttpResponseFactory(ctrl)
		req := httptest.NewRequest(http.MethodGet, target, nil)
		respRec := httptest.NewRecorder()
		respRec.WriteHeader(http.StatusInternalServerError)
		_, err := respRec.WriteString(data)
		if err != nil {
			panic(err)
		}
		respRec.Flush()
		resp := respRec.Result()
		resp.Body = nil

		c := httpDataSource{
			url:            url,
			do:             httpWrap.Do,
			createResponse: createResponse.Create,
			createRequest:  createRequest.Create,
		}

		ctx, _ := makeLoggerContext(ctrl)
		expErr := errors.New("Test")

		createRequest.EXPECT().Create(ctx, nil, target, http.MethodGet).Return(req, nil).Times(1)
		httpWrap.EXPECT().Do(req).Return(resp, expErr).Times(1)
		r, err := c.Get(ctx, []byte(key))
		cmperr(t, err, expErr)
		if !cmp.Equal(r, nil) {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("http ok creates response", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		httpWrap := mocks.NewMockHttpWrap(ctrl)
		wrapResp := &DummyResponse{}
		createRequest := mocks.NewMockHttpRequestFactory(ctrl)
		createResponse := mocks.NewMockHttpResponseFactory(ctrl)
		req := httptest.NewRequest(http.MethodGet, target, nil)
		respRec := httptest.NewRecorder()
		respRec.WriteHeader(http.StatusOK)
		_, err := respRec.WriteString(data)
		if err != nil {
			panic(err)
		}
		respRec.Flush()
		resp := respRec.Result()
		resp.Body = nil

		c := httpDataSource{
			url:            url,
			do:             httpWrap.Do,
			createResponse: createResponse.Create,
			createRequest:  createRequest.Create,
		}

		ctx, _ := makeLoggerContext(ctrl)

		createRequest.EXPECT().Create(ctx, nil, target, http.MethodGet).Return(req, nil).Times(1)
		httpWrap.EXPECT().Do(req).Return(resp, nil).Times(1)
		createResponse.EXPECT().Create(resp).Return(wrapResp, nil).Times(1)
		r, err := c.Get(ctx, []byte(key))
		cmperr(t, err, nil)
		if !cmp.Equal(r, wrapResp) {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("http ok creates response (can be nil)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		httpWrap := mocks.NewMockHttpWrap(ctrl)
		var wrapResp middleware.Response = nil
		createRequest := mocks.NewMockHttpRequestFactory(ctrl)
		createResponse := mocks.NewMockHttpResponseFactory(ctrl)
		req := httptest.NewRequest(http.MethodGet, target, nil)
		respRec := httptest.NewRecorder()
		respRec.WriteHeader(http.StatusOK)
		_, err := respRec.WriteString(data)
		if err != nil {
			panic(err)
		}
		respRec.Flush()
		resp := respRec.Result()
		resp.Body = nil

		c := httpDataSource{
			url:            url,
			do:             httpWrap.Do,
			createResponse: createResponse.Create,
			createRequest:  createRequest.Create,
		}

		ctx, _ := makeLoggerContext(ctrl)

		createRequest.EXPECT().Create(ctx, nil, target, http.MethodGet).Return(req, nil).Times(1)
		httpWrap.EXPECT().Do(req).Return(resp, nil).Times(1)
		createResponse.EXPECT().Create(resp).Return(wrapResp, nil).Times(1)
		r, err := c.Get(ctx, []byte(key))
		cmperr(t, err, nil)
		if !cmp.Equal(r, wrapResp) {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("http ok creates response, if response is created with error - return resp and error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		httpWrap := mocks.NewMockHttpWrap(ctrl)
		wrapResp := &DummyResponse{}
		createRequest := mocks.NewMockHttpRequestFactory(ctrl)
		createResponse := mocks.NewMockHttpResponseFactory(ctrl)
		req := httptest.NewRequest(http.MethodGet, target, nil)
		respRec := httptest.NewRecorder()
		respRec.WriteHeader(http.StatusInternalServerError)
		_, err := respRec.WriteString(data)
		if err != nil {
			panic(err)
		}
		respRec.Flush()
		resp := respRec.Result()
		resp.Body = nil

		c := httpDataSource{
			url:            url,
			do:             httpWrap.Do,
			createResponse: createResponse.Create,
			createRequest:  createRequest.Create,
		}

		ctx, _ := makeLoggerContext(ctrl)
		expErr := errors.New("Test")

		createRequest.EXPECT().Create(ctx, nil, target, http.MethodGet).Return(req, nil).Times(1)
		httpWrap.EXPECT().Do(req).Return(resp, nil).Times(1)
		createResponse.EXPECT().Create(resp).Return(wrapResp, expErr).Times(1)
		r, err := c.Get(ctx, []byte(key))
		cmperr(t, err, expErr)
		if !cmp.Equal(r, wrapResp) {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("http ok creates response, if response is created without error, but status != 200 returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		httpWrap := mocks.NewMockHttpWrap(ctrl)
		wrapResp := &DummyResponse{}
		createRequest := mocks.NewMockHttpRequestFactory(ctrl)
		createResponse := mocks.NewMockHttpResponseFactory(ctrl)
		req := httptest.NewRequest(http.MethodGet, target, nil)
		respRec := httptest.NewRecorder()
		respRec.WriteHeader(http.StatusInternalServerError)
		_, err := respRec.WriteString(data)
		if err != nil {
			panic(err)
		}
		respRec.Flush()
		resp := respRec.Result()
		resp.Body = nil

		c := httpDataSource{
			url:            url,
			do:             httpWrap.Do,
			createResponse: createResponse.Create,
			createRequest:  createRequest.Create,
		}

		ctx, _ := makeLoggerContext(ctrl)

		createRequest.EXPECT().Create(ctx, nil, target, http.MethodGet).Return(req, nil).Times(1)
		httpWrap.EXPECT().Do(req).Return(resp, nil).Times(1)
		createResponse.EXPECT().Create(resp).Return(wrapResp, nil).Times(1)
		r, err := c.Get(ctx, []byte(key))
		if err == nil {
			t.Errorf("Failed. Code: %v. Error: %v", resp.StatusCode, err)
		}
		if !cmp.Equal(r, wrapResp) {
			t.Errorf("Expected correct response.")
		}
	})
}

func TestHttpDataSource_Create(t *testing.T) {
	target := url
	data := "some random data"

	t.Run("create request error is a failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockDo := mocks.NewMockHttpWrap(ctrl)
		createRequest := mocks.NewMockHttpRequestFactory(ctrl)
		createResponse := mocks.NewMockHttpResponseFactory(ctrl)

		c := httpDataSource{
			url:            url,
			do:             mockDo.Do,
			createResponse: createResponse.Create,
			createRequest:  createRequest.Create,
		}

		ctx, _ := makeLoggerContext(ctrl)
		expErr := errors.New("Test")

		createRequest.EXPECT().Create(ctx, []byte(key), target, http.MethodPost).Return(nil, expErr).Times(1)
		r, err := c.Create(ctx, []byte(key))
		cmperr(t, err, expErr)
		if !cmp.Equal(r, nil) {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("http call error is a failure (resp not nil and closes)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		httpWrap := mocks.NewMockHttpWrap(ctrl)
		createRequest := mocks.NewMockHttpRequestFactory(ctrl)
		createResponse := mocks.NewMockHttpResponseFactory(ctrl)
		responseBody := mocks.NewMockReadCloser(ctrl)
		req := httptest.NewRequest(http.MethodPost, target, nil)
		respRec := httptest.NewRecorder()
		respRec.WriteHeader(http.StatusInternalServerError)
		_, err := respRec.WriteString(data)
		if err != nil {
			panic(err)
		}
		respRec.Flush()
		resp := respRec.Result()
		resp.Body = responseBody

		c := httpDataSource{
			url:            url,
			do:             httpWrap.Do,
			createResponse: createResponse.Create,
			createRequest:  createRequest.Create,
		}

		ctx, _ := makeLoggerContext(ctrl)
		expErr := errors.New("Test")

		createRequest.EXPECT().Create(ctx, []byte(key), target, http.MethodPost).Return(req, nil).Times(1)
		httpWrap.EXPECT().Do(req).Return(resp, expErr).Times(1)
		responseBody.EXPECT().Close().Return(nil).Times(1)
		r, err := c.Create(ctx, []byte(key))
		cmperr(t, err, expErr)
		if !cmp.Equal(r, nil) {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("http call error is a failure (resp not nil and fails to close)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		httpWrap := mocks.NewMockHttpWrap(ctrl)
		createRequest := mocks.NewMockHttpRequestFactory(ctrl)
		createResponse := mocks.NewMockHttpResponseFactory(ctrl)
		responseBody := mocks.NewMockReadCloser(ctrl)
		req := httptest.NewRequest(http.MethodPost, target, nil)
		respRec := httptest.NewRecorder()
		respRec.WriteHeader(http.StatusInternalServerError)
		_, err := respRec.WriteString(data)
		if err != nil {
			panic(err)
		}
		respRec.Flush()
		resp := respRec.Result()
		resp.Body = responseBody

		c := httpDataSource{
			url:            url,
			do:             httpWrap.Do,
			createResponse: createResponse.Create,
			createRequest:  createRequest.Create,
		}

		ctx, logger := makeLoggerContext(ctrl)
		expErr := errors.New("Test")
		expErr2 := errors.New("Test2")

		createRequest.EXPECT().Create(ctx, []byte(key), target, http.MethodPost).Return(req, nil).Times(1)
		httpWrap.EXPECT().Do(req).Return(resp, expErr).Times(1)
		responseBody.EXPECT().Close().Return(expErr2).Times(1)
		logger.EXPECT().Warningf(gomock.Any(), expErr2).Times(1)
		r, err := c.Create(ctx, []byte(key))
		cmperr(t, err, expErr)
		if !cmp.Equal(r, nil) {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("http call error is a failure (resp nil)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		httpWrap := mocks.NewMockHttpWrap(ctrl)
		createRequest := mocks.NewMockHttpRequestFactory(ctrl)
		createResponse := mocks.NewMockHttpResponseFactory(ctrl)
		req := httptest.NewRequest(http.MethodPost, target, nil)

		c := httpDataSource{
			url:            url,
			do:             httpWrap.Do,
			createResponse: createResponse.Create,
			createRequest:  createRequest.Create,
		}

		ctx, _ := makeLoggerContext(ctrl)
		expErr := errors.New("Test")

		createRequest.EXPECT().Create(ctx, []byte(key), target, http.MethodPost).Return(req, nil).Times(1)
		httpWrap.EXPECT().Do(req).Return(nil, expErr).Times(1)
		r, err := c.Create(ctx, []byte(key))
		cmperr(t, err, expErr)
		if !cmp.Equal(r, nil) {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("http call error is a failure (resp body is nil)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		httpWrap := mocks.NewMockHttpWrap(ctrl)
		createRequest := mocks.NewMockHttpRequestFactory(ctrl)
		createResponse := mocks.NewMockHttpResponseFactory(ctrl)
		req := httptest.NewRequest(http.MethodPost, target, nil)
		respRec := httptest.NewRecorder()
		respRec.WriteHeader(http.StatusInternalServerError)
		_, err := respRec.WriteString(data)
		if err != nil {
			panic(err)
		}
		respRec.Flush()
		resp := respRec.Result()
		resp.Body = nil

		c := httpDataSource{
			url:            url,
			do:             httpWrap.Do,
			createResponse: createResponse.Create,
			createRequest:  createRequest.Create,
		}

		ctx, _ := makeLoggerContext(ctrl)
		expErr := errors.New("Test")

		createRequest.EXPECT().Create(ctx, []byte(key), target, http.MethodPost).Return(req, nil).Times(1)
		httpWrap.EXPECT().Do(req).Return(resp, expErr).Times(1)
		r, err := c.Create(ctx, []byte(key))
		cmperr(t, err, expErr)
		if !cmp.Equal(r, nil) {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("http ok creates response", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		httpWrap := mocks.NewMockHttpWrap(ctrl)
		wrapResp := &DummyResponse{}
		createRequest := mocks.NewMockHttpRequestFactory(ctrl)
		createResponse := mocks.NewMockHttpResponseFactory(ctrl)
		req := httptest.NewRequest(http.MethodPost, target, nil)
		respRec := httptest.NewRecorder()
		respRec.WriteHeader(http.StatusOK)
		_, err := respRec.WriteString(data)
		if err != nil {
			panic(err)
		}
		respRec.Flush()
		resp := respRec.Result()
		resp.Body = nil

		c := httpDataSource{
			url:            url,
			do:             httpWrap.Do,
			createResponse: createResponse.Create,
			createRequest:  createRequest.Create,
		}

		ctx, _ := makeLoggerContext(ctrl)

		createRequest.EXPECT().Create(ctx, []byte(key), target, http.MethodPost).Return(req, nil).Times(1)
		httpWrap.EXPECT().Do(req).Return(resp, nil).Times(1)
		createResponse.EXPECT().Create(resp).Return(wrapResp, nil).Times(1)
		r, err := c.Create(ctx, []byte(key))
		cmperr(t, err, nil)
		if !cmp.Equal(r, wrapResp) {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("http ok creates response (can be nil)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		httpWrap := mocks.NewMockHttpWrap(ctrl)
		var wrapResp middleware.Response = nil
		createRequest := mocks.NewMockHttpRequestFactory(ctrl)
		createResponse := mocks.NewMockHttpResponseFactory(ctrl)
		req := httptest.NewRequest(http.MethodPost, target, nil)
		respRec := httptest.NewRecorder()
		respRec.WriteHeader(http.StatusOK)
		_, err := respRec.WriteString(data)
		if err != nil {
			panic(err)
		}
		respRec.Flush()
		resp := respRec.Result()
		resp.Body = nil

		c := httpDataSource{
			url:            url,
			do:             httpWrap.Do,
			createResponse: createResponse.Create,
			createRequest:  createRequest.Create,
		}

		ctx, _ := makeLoggerContext(ctrl)

		createRequest.EXPECT().Create(ctx, []byte(key), target, http.MethodPost).Return(req, nil).Times(1)
		httpWrap.EXPECT().Do(req).Return(resp, nil).Times(1)
		createResponse.EXPECT().Create(resp).Return(wrapResp, nil).Times(1)
		r, err := c.Create(ctx, []byte(key))
		cmperr(t, err, nil)
		if !cmp.Equal(r, wrapResp) {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("http ok creates response, if response is created with error - return resp and error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		httpWrap := mocks.NewMockHttpWrap(ctrl)
		wrapResp := &DummyResponse{}
		createRequest := mocks.NewMockHttpRequestFactory(ctrl)
		createResponse := mocks.NewMockHttpResponseFactory(ctrl)
		req := httptest.NewRequest(http.MethodPost, target, nil)
		respRec := httptest.NewRecorder()
		respRec.WriteHeader(http.StatusInternalServerError)
		_, err := respRec.WriteString(data)
		if err != nil {
			panic(err)
		}
		respRec.Flush()
		resp := respRec.Result()
		resp.Body = nil

		c := httpDataSource{
			url:            url,
			do:             httpWrap.Do,
			createResponse: createResponse.Create,
			createRequest:  createRequest.Create,
		}

		ctx, _ := makeLoggerContext(ctrl)
		expErr := errors.New("Test")

		createRequest.EXPECT().Create(ctx, []byte(key), target, http.MethodPost).Return(req, nil).Times(1)
		httpWrap.EXPECT().Do(req).Return(resp, nil).Times(1)
		createResponse.EXPECT().Create(resp).Return(wrapResp, expErr).Times(1)
		r, err := c.Create(ctx, []byte(key))
		cmperr(t, err, expErr)
		if !cmp.Equal(r, wrapResp) {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("http ok creates response, if response is created without error, but status != 200 returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		httpWrap := mocks.NewMockHttpWrap(ctrl)
		wrapResp := &DummyResponse{}
		createRequest := mocks.NewMockHttpRequestFactory(ctrl)
		createResponse := mocks.NewMockHttpResponseFactory(ctrl)
		req := httptest.NewRequest(http.MethodPost, target, nil)
		respRec := httptest.NewRecorder()
		respRec.WriteHeader(http.StatusInternalServerError)
		_, err := respRec.WriteString(data)
		if err != nil {
			panic(err)
		}
		respRec.Flush()
		resp := respRec.Result()
		resp.Body = nil

		c := httpDataSource{
			url:            url,
			do:             httpWrap.Do,
			createResponse: createResponse.Create,
			createRequest:  createRequest.Create,
		}

		ctx, _ := makeLoggerContext(ctrl)

		createRequest.EXPECT().Create(ctx, []byte(key), target, http.MethodPost).Return(req, nil).Times(1)
		httpWrap.EXPECT().Do(req).Return(resp, nil).Times(1)
		createResponse.EXPECT().Create(resp).Return(wrapResp, nil).Times(1)
		r, err := c.Create(ctx, []byte(key))
		if err == nil {
			t.Errorf("Failed. Code: %v. Error: %v", resp.StatusCode, err)
		}
		if !cmp.Equal(r, wrapResp) {
			t.Errorf("Expected correct response.")
		}
	})
}

func TestHttpDataSource_Update(t *testing.T) {
	target := url
	data := "some random data"

	t.Run("create request error is a failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockDo := mocks.NewMockHttpWrap(ctrl)
		createRequest := mocks.NewMockHttpRequestFactory(ctrl)
		createResponse := mocks.NewMockHttpResponseFactory(ctrl)

		c := httpDataSource{
			url:            url,
			do:             mockDo.Do,
			createResponse: createResponse.Create,
			createRequest:  createRequest.Create,
		}

		ctx, _ := makeLoggerContext(ctrl)
		expErr := errors.New("Test")

		createRequest.EXPECT().Create(ctx, []byte(key), target, http.MethodPut).Return(nil, expErr).Times(1)
		r, err := c.Update(ctx, []byte(key))
		cmperr(t, err, expErr)
		if !cmp.Equal(r, nil) {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("http call error is a failure (resp not nil and closes)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		httpWrap := mocks.NewMockHttpWrap(ctrl)
		createRequest := mocks.NewMockHttpRequestFactory(ctrl)
		createResponse := mocks.NewMockHttpResponseFactory(ctrl)
		responseBody := mocks.NewMockReadCloser(ctrl)
		req := httptest.NewRequest(http.MethodPut, target, nil)
		respRec := httptest.NewRecorder()
		respRec.WriteHeader(http.StatusInternalServerError)
		_, err := respRec.WriteString(data)
		if err != nil {
			panic(err)
		}
		respRec.Flush()
		resp := respRec.Result()
		resp.Body = responseBody

		c := httpDataSource{
			url:            url,
			do:             httpWrap.Do,
			createResponse: createResponse.Create,
			createRequest:  createRequest.Create,
		}

		ctx, _ := makeLoggerContext(ctrl)
		expErr := errors.New("Test")

		createRequest.EXPECT().Create(ctx, []byte(key), target, http.MethodPut).Return(req, nil).Times(1)
		httpWrap.EXPECT().Do(req).Return(resp, expErr).Times(1)
		responseBody.EXPECT().Close().Return(nil).Times(1)
		r, err := c.Update(ctx, []byte(key))
		cmperr(t, err, expErr)
		if !cmp.Equal(r, nil) {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("http call error is a failure (resp not nil and fails to close)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		httpWrap := mocks.NewMockHttpWrap(ctrl)
		createRequest := mocks.NewMockHttpRequestFactory(ctrl)
		createResponse := mocks.NewMockHttpResponseFactory(ctrl)
		responseBody := mocks.NewMockReadCloser(ctrl)
		req := httptest.NewRequest(http.MethodPut, target, nil)
		respRec := httptest.NewRecorder()
		respRec.WriteHeader(http.StatusInternalServerError)
		_, err := respRec.WriteString(data)
		if err != nil {
			panic(err)
		}
		respRec.Flush()
		resp := respRec.Result()
		resp.Body = responseBody

		c := httpDataSource{
			url:            url,
			do:             httpWrap.Do,
			createResponse: createResponse.Create,
			createRequest:  createRequest.Create,
		}

		ctx, logger := makeLoggerContext(ctrl)
		expErr := errors.New("Test")
		expErr2 := errors.New("Test2")

		createRequest.EXPECT().Create(ctx, []byte(key), target, http.MethodPut).Return(req, nil).Times(1)
		httpWrap.EXPECT().Do(req).Return(resp, expErr).Times(1)
		responseBody.EXPECT().Close().Return(expErr2).Times(1)
		logger.EXPECT().Warningf(gomock.Any(), expErr2).Times(1)
		r, err := c.Update(ctx, []byte(key))
		cmperr(t, err, expErr)
		if !cmp.Equal(r, nil) {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("http call error is a failure (resp nil)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		httpWrap := mocks.NewMockHttpWrap(ctrl)
		createRequest := mocks.NewMockHttpRequestFactory(ctrl)
		createResponse := mocks.NewMockHttpResponseFactory(ctrl)
		req := httptest.NewRequest(http.MethodPut, target, nil)

		c := httpDataSource{
			url:            url,
			do:             httpWrap.Do,
			createResponse: createResponse.Create,
			createRequest:  createRequest.Create,
		}

		ctx, _ := makeLoggerContext(ctrl)
		expErr := errors.New("Test")

		createRequest.EXPECT().Create(ctx, []byte(key), target, http.MethodPut).Return(req, nil).Times(1)
		httpWrap.EXPECT().Do(req).Return(nil, expErr).Times(1)
		r, err := c.Update(ctx, []byte(key))
		cmperr(t, err, expErr)
		if !cmp.Equal(r, nil) {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("http call error is a failure (resp body is nil)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		httpWrap := mocks.NewMockHttpWrap(ctrl)
		createRequest := mocks.NewMockHttpRequestFactory(ctrl)
		createResponse := mocks.NewMockHttpResponseFactory(ctrl)
		req := httptest.NewRequest(http.MethodPut, target, nil)
		respRec := httptest.NewRecorder()
		respRec.WriteHeader(http.StatusInternalServerError)
		_, err := respRec.WriteString(data)
		if err != nil {
			panic(err)
		}
		respRec.Flush()
		resp := respRec.Result()
		resp.Body = nil

		c := httpDataSource{
			url:            url,
			do:             httpWrap.Do,
			createResponse: createResponse.Create,
			createRequest:  createRequest.Create,
		}

		ctx, _ := makeLoggerContext(ctrl)
		expErr := errors.New("Test")

		createRequest.EXPECT().Create(ctx, []byte(key), target, http.MethodPut).Return(req, nil).Times(1)
		httpWrap.EXPECT().Do(req).Return(resp, expErr).Times(1)
		r, err := c.Update(ctx, []byte(key))
		cmperr(t, err, expErr)
		if !cmp.Equal(r, nil) {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("http ok creates response", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		httpWrap := mocks.NewMockHttpWrap(ctrl)
		wrapResp := &DummyResponse{}
		createRequest := mocks.NewMockHttpRequestFactory(ctrl)
		createResponse := mocks.NewMockHttpResponseFactory(ctrl)
		req := httptest.NewRequest(http.MethodPut, target, nil)
		respRec := httptest.NewRecorder()
		respRec.WriteHeader(http.StatusOK)
		_, err := respRec.WriteString(data)
		if err != nil {
			panic(err)
		}
		respRec.Flush()
		resp := respRec.Result()
		resp.Body = nil

		c := httpDataSource{
			url:            url,
			do:             httpWrap.Do,
			createResponse: createResponse.Create,
			createRequest:  createRequest.Create,
		}

		ctx, _ := makeLoggerContext(ctrl)

		createRequest.EXPECT().Create(ctx, []byte(key), target, http.MethodPut).Return(req, nil).Times(1)
		httpWrap.EXPECT().Do(req).Return(resp, nil).Times(1)
		createResponse.EXPECT().Create(resp).Return(wrapResp, nil).Times(1)
		r, err := c.Update(ctx, []byte(key))
		cmperr(t, err, nil)
		if !cmp.Equal(r, wrapResp) {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("http ok creates response (can be nil)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		httpWrap := mocks.NewMockHttpWrap(ctrl)
		var wrapResp middleware.Response = nil
		createRequest := mocks.NewMockHttpRequestFactory(ctrl)
		createResponse := mocks.NewMockHttpResponseFactory(ctrl)
		req := httptest.NewRequest(http.MethodPut, target, nil)
		respRec := httptest.NewRecorder()
		respRec.WriteHeader(http.StatusOK)
		_, err := respRec.WriteString(data)
		if err != nil {
			panic(err)
		}
		respRec.Flush()
		resp := respRec.Result()
		resp.Body = nil

		c := httpDataSource{
			url:            url,
			do:             httpWrap.Do,
			createResponse: createResponse.Create,
			createRequest:  createRequest.Create,
		}

		ctx, _ := makeLoggerContext(ctrl)

		createRequest.EXPECT().Create(ctx, []byte(key), target, http.MethodPut).Return(req, nil).Times(1)
		httpWrap.EXPECT().Do(req).Return(resp, nil).Times(1)
		createResponse.EXPECT().Create(resp).Return(wrapResp, nil).Times(1)
		r, err := c.Update(ctx, []byte(key))
		cmperr(t, err, nil)
		if !cmp.Equal(r, wrapResp) {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("http ok creates response, if response is created with error - return resp and error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		httpWrap := mocks.NewMockHttpWrap(ctrl)
		wrapResp := &DummyResponse{}
		createRequest := mocks.NewMockHttpRequestFactory(ctrl)
		createResponse := mocks.NewMockHttpResponseFactory(ctrl)
		req := httptest.NewRequest(http.MethodPut, target, nil)
		respRec := httptest.NewRecorder()
		respRec.WriteHeader(http.StatusInternalServerError)
		_, err := respRec.WriteString(data)
		if err != nil {
			panic(err)
		}
		respRec.Flush()
		resp := respRec.Result()
		resp.Body = nil

		c := httpDataSource{
			url:            url,
			do:             httpWrap.Do,
			createResponse: createResponse.Create,
			createRequest:  createRequest.Create,
		}

		ctx, _ := makeLoggerContext(ctrl)
		expErr := errors.New("Test")

		createRequest.EXPECT().Create(ctx, []byte(key), target, http.MethodPut).Return(req, nil).Times(1)
		httpWrap.EXPECT().Do(req).Return(resp, nil).Times(1)
		createResponse.EXPECT().Create(resp).Return(wrapResp, expErr).Times(1)
		r, err := c.Update(ctx, []byte(key))
		cmperr(t, err, expErr)
		if !cmp.Equal(r, wrapResp) {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("http ok creates response, if response is created without error, but status != 200 returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		httpWrap := mocks.NewMockHttpWrap(ctrl)
		wrapResp := &DummyResponse{}
		createRequest := mocks.NewMockHttpRequestFactory(ctrl)
		createResponse := mocks.NewMockHttpResponseFactory(ctrl)
		req := httptest.NewRequest(http.MethodPut, target, nil)
		respRec := httptest.NewRecorder()
		respRec.WriteHeader(http.StatusInternalServerError)
		_, err := respRec.WriteString(data)
		if err != nil {
			panic(err)
		}
		respRec.Flush()
		resp := respRec.Result()
		resp.Body = nil

		c := httpDataSource{
			url:            url,
			do:             httpWrap.Do,
			createResponse: createResponse.Create,
			createRequest:  createRequest.Create,
		}

		ctx, _ := makeLoggerContext(ctrl)

		createRequest.EXPECT().Create(ctx, []byte(key), target, http.MethodPut).Return(req, nil).Times(1)
		httpWrap.EXPECT().Do(req).Return(resp, nil).Times(1)
		createResponse.EXPECT().Create(resp).Return(wrapResp, nil).Times(1)
		r, err := c.Update(ctx, []byte(key))
		if err == nil {
			t.Errorf("Failed. Code: %v. Error: %v", resp.StatusCode, err)
		}
		if !cmp.Equal(r, wrapResp) {
			t.Errorf("Expected correct response.")
		}
	})
}
