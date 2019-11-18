package sources

import (
	"context"
	"errors"
	"github.com/coldze/test/logic"
	"github.com/coldze/test/mocks"
	"github.com/coldze/test/mocks/mock_logs"
	"github.com/coldze/test/mocks/mock_sources"
	"github.com/coldze/test/mocks/mock_std"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
)

type httpDataSourceFixture struct {
	Url            string
	Target         string
	Key            string
	Data           string
	Error          error
	Ctx            context.Context
	Logger         *mock_logs.MockLogger
	Do             *mock_std.MockHttpWrap
	CreateRequest  *mock_sources.MockHttpRequestFactory
	CreateResponse *mock_sources.MockHttpResponseFactory
}

func newHttpDataSourceFixture(ctrl *gomock.Controller) *httpDataSourceFixture {
	url := "https://test.url.com/v1/api"
	key := "some_random_key"
	target := url + "/" + key
	ctx, logger := makeLoggerContext(ctrl)
	return &httpDataSourceFixture{
		Url:            url,
		Target:         target,
		Key:            key,
		Data:           "some random data",
		Error:          errors.New("Some test error"),
		Ctx:            ctx,
		Logger:         logger,
		Do:             mock_std.NewMockHttpWrap(ctrl),
		CreateRequest:  mock_sources.NewMockHttpRequestFactory(ctrl),
		CreateResponse: mock_sources.NewMockHttpResponseFactory(ctrl),
	}
}

func newTestableHttpDataSource(f *httpDataSourceFixture) *httpDataSource {
	return &httpDataSource{
		url:            f.Url,
		do:             f.Do.Do,
		createResponse: f.CreateResponse.Create,
		createRequest:  f.CreateRequest.Create,
	}
}

func TestHttpDataSource_Get(t *testing.T) {

	t.Run("create request error is a failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newHttpDataSourceFixture(ctrl)
		c := newTestableHttpDataSource(f)

		f.CreateRequest.EXPECT().Create(f.Ctx, nil, f.Target, http.MethodGet).Return(nil, f.Error).Times(1)
		r, err := c.Get(f.Ctx, []byte(f.Key))
		mocks.CmpError(t, err, f.Error)
		if !cmp.Equal(r, nil) {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("http call error is a failure (resp not nil and closes)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newHttpDataSourceFixture(ctrl)
		c := newTestableHttpDataSource(f)

		responseBody := mock_std.NewMockReadCloser(ctrl)
		req := httptest.NewRequest(http.MethodGet, f.Target, nil)
		respRec := httptest.NewRecorder()
		respRec.WriteHeader(http.StatusInternalServerError)
		_, err := respRec.WriteString(f.Data)
		if err != nil {
			panic(err)
		}
		respRec.Flush()
		resp := respRec.Result()
		resp.Body = responseBody

		f.CreateRequest.EXPECT().Create(f.Ctx, nil, f.Target, http.MethodGet).Return(req, nil).Times(1)
		f.Do.EXPECT().Do(req).Return(resp, f.Error).Times(1)
		responseBody.EXPECT().Close().Return(nil).Times(1)
		r, err := c.Get(f.Ctx, []byte(f.Key))
		mocks.CmpError(t, err, f.Error)
		if !cmp.Equal(r, nil) {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("http call error is a failure (resp not nil and fails to close)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newHttpDataSourceFixture(ctrl)
		c := newTestableHttpDataSource(f)

		responseBody := mock_std.NewMockReadCloser(ctrl)
		req := httptest.NewRequest(http.MethodGet, f.Target, nil)
		respRec := httptest.NewRecorder()
		respRec.WriteHeader(http.StatusInternalServerError)
		_, err := respRec.WriteString(f.Data)
		if err != nil {
			panic(err)
		}
		respRec.Flush()
		resp := respRec.Result()
		resp.Body = responseBody

		expErr := errors.New("Test")

		f.CreateRequest.EXPECT().Create(f.Ctx, nil, f.Target, http.MethodGet).Return(req, nil).Times(1)
		f.Do.EXPECT().Do(req).Return(resp, f.Error).Times(1)
		responseBody.EXPECT().Close().Return(expErr).Times(1)
		f.Logger.EXPECT().Warningf(gomock.Any(), expErr).Times(1)
		r, err := c.Get(f.Ctx, []byte(f.Key))
		mocks.CmpError(t, err, f.Error)
		if !cmp.Equal(r, nil) {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("http call error is a failure (resp nil)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newHttpDataSourceFixture(ctrl)
		c := newTestableHttpDataSource(f)

		req := httptest.NewRequest(http.MethodGet, f.Target, nil)

		f.CreateRequest.EXPECT().Create(f.Ctx, nil, f.Target, http.MethodGet).Return(req, nil).Times(1)
		f.Do.EXPECT().Do(req).Return(nil, f.Error).Times(1)
		r, err := c.Get(f.Ctx, []byte(f.Key))
		mocks.CmpError(t, err, f.Error)
		if !cmp.Equal(r, nil) {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("http call error is a failure (resp body is nil)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newHttpDataSourceFixture(ctrl)
		c := newTestableHttpDataSource(f)

		req := httptest.NewRequest(http.MethodGet, f.Target, nil)
		respRec := httptest.NewRecorder()
		respRec.WriteHeader(http.StatusInternalServerError)
		_, err := respRec.WriteString(f.Data)
		if err != nil {
			panic(err)
		}
		respRec.Flush()
		resp := respRec.Result()
		resp.Body = nil

		f.CreateRequest.EXPECT().Create(f.Ctx, nil, f.Target, http.MethodGet).Return(req, nil).Times(1)
		f.Do.EXPECT().Do(req).Return(resp, f.Error).Times(1)
		r, err := c.Get(f.Ctx, []byte(f.Key))
		mocks.CmpError(t, err, f.Error)
		if !cmp.Equal(r, nil) {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("http ok creates response", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newHttpDataSourceFixture(ctrl)
		c := newTestableHttpDataSource(f)

		wrapResp := &DummyResponse{}
		req := httptest.NewRequest(http.MethodGet, f.Target, nil)
		respRec := httptest.NewRecorder()
		respRec.WriteHeader(http.StatusOK)
		_, err := respRec.WriteString(f.Data)
		if err != nil {
			panic(err)
		}
		respRec.Flush()
		resp := respRec.Result()
		resp.Body = nil

		f.CreateRequest.EXPECT().Create(f.Ctx, nil, f.Target, http.MethodGet).Return(req, nil).Times(1)
		f.Do.EXPECT().Do(req).Return(resp, nil).Times(1)
		f.CreateResponse.EXPECT().Create(resp).Return(wrapResp, nil).Times(1)
		r, err := c.Get(f.Ctx, []byte(f.Key))
		mocks.CmpError(t, err, nil)
		if !cmp.Equal(r, wrapResp) {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("http ok creates response (can be nil)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newHttpDataSourceFixture(ctrl)
		c := newTestableHttpDataSource(f)

		var wrapResp logic.Response = nil
		req := httptest.NewRequest(http.MethodGet, f.Target, nil)
		respRec := httptest.NewRecorder()
		respRec.WriteHeader(http.StatusOK)
		_, err := respRec.WriteString(f.Data)
		if err != nil {
			panic(err)
		}
		respRec.Flush()
		resp := respRec.Result()
		resp.Body = nil

		f.CreateRequest.EXPECT().Create(f.Ctx, nil, f.Target, http.MethodGet).Return(req, nil).Times(1)
		f.Do.EXPECT().Do(req).Return(resp, nil).Times(1)
		f.CreateResponse.EXPECT().Create(resp).Return(wrapResp, nil).Times(1)
		r, err := c.Get(f.Ctx, []byte(f.Key))
		mocks.CmpError(t, err, nil)
		if !cmp.Equal(r, wrapResp) {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("http ok creates response, if response is created with error - return resp and error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newHttpDataSourceFixture(ctrl)
		c := newTestableHttpDataSource(f)

		wrapResp := &DummyResponse{}
		req := httptest.NewRequest(http.MethodGet, f.Target, nil)
		respRec := httptest.NewRecorder()
		respRec.WriteHeader(http.StatusInternalServerError)
		_, err := respRec.WriteString(f.Data)
		if err != nil {
			panic(err)
		}
		respRec.Flush()
		resp := respRec.Result()
		resp.Body = nil

		f.CreateRequest.EXPECT().Create(f.Ctx, nil, f.Target, http.MethodGet).Return(req, nil).Times(1)
		f.Do.EXPECT().Do(req).Return(resp, nil).Times(1)
		f.CreateResponse.EXPECT().Create(resp).Return(wrapResp, f.Error).Times(1)
		r, err := c.Get(f.Ctx, []byte(f.Key))
		mocks.CmpError(t, err, f.Error)
		if !cmp.Equal(r, wrapResp) {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("http ok creates response, if response is created without error, but status != 200 returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newHttpDataSourceFixture(ctrl)
		c := newTestableHttpDataSource(f)

		wrapResp := &DummyResponse{}
		req := httptest.NewRequest(http.MethodGet, f.Target, nil)
		respRec := httptest.NewRecorder()
		respRec.WriteHeader(http.StatusInternalServerError)
		_, err := respRec.WriteString(f.Data)
		if err != nil {
			panic(err)
		}
		respRec.Flush()
		resp := respRec.Result()
		resp.Body = nil

		f.CreateRequest.EXPECT().Create(f.Ctx, nil, f.Target, http.MethodGet).Return(req, nil).Times(1)
		f.Do.EXPECT().Do(req).Return(resp, nil).Times(1)
		f.CreateResponse.EXPECT().Create(resp).Return(wrapResp, nil).Times(1)
		r, err := c.Get(f.Ctx, []byte(f.Key))
		if err == nil {
			t.Errorf("Failed. Code: %v. Error: %v", resp.StatusCode, err)
		}
		if !cmp.Equal(r, wrapResp) {
			t.Errorf("Expected correct response.")
		}
	})
}

func TestHttpDataSource_Create(t *testing.T) {
	t.Run("create request error is a failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newHttpDataSourceFixture(ctrl)
		c := newTestableHttpDataSource(f)

		f.CreateRequest.EXPECT().Create(f.Ctx, []byte(f.Key), f.Url, http.MethodPost).Return(nil, f.Error).Times(1)
		r, err := c.Create(f.Ctx, []byte(f.Key))
		mocks.CmpError(t, err, f.Error)
		if !cmp.Equal(r, nil) {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("http call error is a failure (resp not nil and closes)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newHttpDataSourceFixture(ctrl)
		c := newTestableHttpDataSource(f)

		responseBody := mock_std.NewMockReadCloser(ctrl)
		req := httptest.NewRequest(http.MethodPost, f.Target, nil)
		respRec := httptest.NewRecorder()
		respRec.WriteHeader(http.StatusInternalServerError)
		_, err := respRec.WriteString(f.Data)
		if err != nil {
			panic(err)
		}
		respRec.Flush()
		resp := respRec.Result()
		resp.Body = responseBody

		f.CreateRequest.EXPECT().Create(f.Ctx, []byte(f.Key), f.Url, http.MethodPost).Return(req, nil).Times(1)
		f.Do.EXPECT().Do(req).Return(resp, f.Error).Times(1)
		responseBody.EXPECT().Close().Return(nil).Times(1)
		r, err := c.Create(f.Ctx, []byte(f.Key))
		mocks.CmpError(t, err, f.Error)
		if !cmp.Equal(r, nil) {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("http call error is a failure (resp not nil and fails to close)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newHttpDataSourceFixture(ctrl)
		c := newTestableHttpDataSource(f)

		responseBody := mock_std.NewMockReadCloser(ctrl)
		req := httptest.NewRequest(http.MethodPost, f.Target, nil)
		respRec := httptest.NewRecorder()
		respRec.WriteHeader(http.StatusInternalServerError)
		_, err := respRec.WriteString(f.Data)
		if err != nil {
			panic(err)
		}
		respRec.Flush()
		resp := respRec.Result()
		resp.Body = responseBody

		expErr := errors.New("Test")

		f.CreateRequest.EXPECT().Create(f.Ctx, []byte(f.Key), f.Url, http.MethodPost).Return(req, nil).Times(1)
		f.Do.EXPECT().Do(req).Return(resp, f.Error).Times(1)
		responseBody.EXPECT().Close().Return(expErr).Times(1)
		f.Logger.EXPECT().Warningf(gomock.Any(), expErr).Times(1)
		r, err := c.Create(f.Ctx, []byte(f.Key))
		mocks.CmpError(t, err, f.Error)
		if !cmp.Equal(r, nil) {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("http call error is a failure (resp nil)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newHttpDataSourceFixture(ctrl)
		c := newTestableHttpDataSource(f)

		req := httptest.NewRequest(http.MethodPost, f.Target, nil)

		f.CreateRequest.EXPECT().Create(f.Ctx, []byte(f.Key), f.Url, http.MethodPost).Return(req, nil).Times(1)
		f.Do.EXPECT().Do(req).Return(nil, f.Error).Times(1)
		r, err := c.Create(f.Ctx, []byte(f.Key))
		mocks.CmpError(t, err, f.Error)
		if !cmp.Equal(r, nil) {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("http call error is a failure (resp body is nil)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newHttpDataSourceFixture(ctrl)
		c := newTestableHttpDataSource(f)

		req := httptest.NewRequest(http.MethodPost, f.Target, nil)
		respRec := httptest.NewRecorder()
		respRec.WriteHeader(http.StatusInternalServerError)
		_, err := respRec.WriteString(f.Data)
		if err != nil {
			panic(err)
		}
		respRec.Flush()
		resp := respRec.Result()
		resp.Body = nil

		f.CreateRequest.EXPECT().Create(f.Ctx, []byte(f.Key), f.Url, http.MethodPost).Return(req, nil).Times(1)
		f.Do.EXPECT().Do(req).Return(resp, f.Error).Times(1)
		r, err := c.Create(f.Ctx, []byte(f.Key))
		mocks.CmpError(t, err, f.Error)
		if !cmp.Equal(r, nil) {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("http ok creates response", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newHttpDataSourceFixture(ctrl)
		c := newTestableHttpDataSource(f)

		wrapResp := &DummyResponse{}
		req := httptest.NewRequest(http.MethodPost, f.Target, nil)
		respRec := httptest.NewRecorder()
		respRec.WriteHeader(http.StatusOK)
		_, err := respRec.WriteString(f.Data)
		if err != nil {
			panic(err)
		}
		respRec.Flush()
		resp := respRec.Result()
		resp.Body = nil

		f.CreateRequest.EXPECT().Create(f.Ctx, []byte(f.Key), f.Url, http.MethodPost).Return(req, nil).Times(1)
		f.Do.EXPECT().Do(req).Return(resp, nil).Times(1)
		f.CreateResponse.EXPECT().Create(resp).Return(wrapResp, nil).Times(1)
		r, err := c.Create(f.Ctx, []byte(f.Key))
		mocks.CmpError(t, err, nil)
		if !cmp.Equal(r, wrapResp) {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("http ok creates response (can be nil)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newHttpDataSourceFixture(ctrl)
		c := newTestableHttpDataSource(f)

		var wrapResp logic.Response = nil
		req := httptest.NewRequest(http.MethodPost, f.Target, nil)
		respRec := httptest.NewRecorder()
		respRec.WriteHeader(http.StatusOK)
		_, err := respRec.WriteString(f.Data)
		if err != nil {
			panic(err)
		}
		respRec.Flush()
		resp := respRec.Result()
		resp.Body = nil

		f.CreateRequest.EXPECT().Create(f.Ctx, []byte(f.Key), f.Url, http.MethodPost).Return(req, nil).Times(1)
		f.Do.EXPECT().Do(req).Return(resp, nil).Times(1)
		f.CreateResponse.EXPECT().Create(resp).Return(wrapResp, nil).Times(1)
		r, err := c.Create(f.Ctx, []byte(f.Key))
		mocks.CmpError(t, err, nil)
		if !cmp.Equal(r, wrapResp) {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("http ok creates response, if response is created with error - return resp and error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newHttpDataSourceFixture(ctrl)
		c := newTestableHttpDataSource(f)

		wrapResp := &DummyResponse{}
		req := httptest.NewRequest(http.MethodPost, f.Target, nil)
		respRec := httptest.NewRecorder()
		respRec.WriteHeader(http.StatusInternalServerError)
		_, err := respRec.WriteString(f.Data)
		if err != nil {
			panic(err)
		}
		respRec.Flush()
		resp := respRec.Result()
		resp.Body = nil

		f.CreateRequest.EXPECT().Create(f.Ctx, []byte(f.Key), f.Url, http.MethodPost).Return(req, nil).Times(1)
		f.Do.EXPECT().Do(req).Return(resp, nil).Times(1)
		f.CreateResponse.EXPECT().Create(resp).Return(wrapResp, f.Error).Times(1)
		r, err := c.Create(f.Ctx, []byte(f.Key))
		mocks.CmpError(t, err, f.Error)
		if !cmp.Equal(r, wrapResp) {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("http ok creates response, if response is created without error, but status != 200 returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newHttpDataSourceFixture(ctrl)
		c := newTestableHttpDataSource(f)

		wrapResp := &DummyResponse{}
		req := httptest.NewRequest(http.MethodPost, f.Target, nil)
		respRec := httptest.NewRecorder()
		respRec.WriteHeader(http.StatusInternalServerError)
		_, err := respRec.WriteString(f.Data)
		if err != nil {
			panic(err)
		}
		respRec.Flush()
		resp := respRec.Result()
		resp.Body = nil

		f.CreateRequest.EXPECT().Create(f.Ctx, []byte(f.Key), f.Url, http.MethodPost).Return(req, nil).Times(1)
		f.Do.EXPECT().Do(req).Return(resp, nil).Times(1)
		f.CreateResponse.EXPECT().Create(resp).Return(wrapResp, nil).Times(1)
		r, err := c.Create(f.Ctx, []byte(f.Key))
		if err == nil {
			t.Errorf("Failed. Code: %v. Error: %v", resp.StatusCode, err)
		}
		if !cmp.Equal(r, wrapResp) {
			t.Errorf("Expected correct response.")
		}
	})
}

func TestHttpDataSource_Update(t *testing.T) {
	t.Run("create request error is a failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newHttpDataSourceFixture(ctrl)
		c := newTestableHttpDataSource(f)

		f.CreateRequest.EXPECT().Create(f.Ctx, []byte(f.Key), f.Url, http.MethodPost).Return(nil, f.Error).Times(1)
		r, err := c.Update(f.Ctx, []byte(f.Key))
		mocks.CmpError(t, err, f.Error)
		if !cmp.Equal(r, nil) {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("http call error is a failure (resp not nil and closes)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newHttpDataSourceFixture(ctrl)
		c := newTestableHttpDataSource(f)

		responseBody := mock_std.NewMockReadCloser(ctrl)
		req := httptest.NewRequest(http.MethodPost, f.Target, nil)
		respRec := httptest.NewRecorder()
		respRec.WriteHeader(http.StatusInternalServerError)
		_, err := respRec.WriteString(f.Data)
		if err != nil {
			panic(err)
		}
		respRec.Flush()
		resp := respRec.Result()
		resp.Body = responseBody

		f.CreateRequest.EXPECT().Create(f.Ctx, []byte(f.Key), f.Url, http.MethodPost).Return(req, nil).Times(1)
		f.Do.EXPECT().Do(req).Return(resp, f.Error).Times(1)
		responseBody.EXPECT().Close().Return(nil).Times(1)
		r, err := c.Update(f.Ctx, []byte(f.Key))
		mocks.CmpError(t, err, f.Error)
		if !cmp.Equal(r, nil) {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("http call error is a failure (resp not nil and fails to close)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newHttpDataSourceFixture(ctrl)
		c := newTestableHttpDataSource(f)

		responseBody := mock_std.NewMockReadCloser(ctrl)
		req := httptest.NewRequest(http.MethodPost, f.Target, nil)
		respRec := httptest.NewRecorder()
		respRec.WriteHeader(http.StatusInternalServerError)
		_, err := respRec.WriteString(f.Data)
		if err != nil {
			panic(err)
		}
		respRec.Flush()
		resp := respRec.Result()
		resp.Body = responseBody

		expErr := errors.New("Test")

		f.CreateRequest.EXPECT().Create(f.Ctx, []byte(f.Key), f.Url, http.MethodPost).Return(req, nil).Times(1)
		f.Do.EXPECT().Do(req).Return(resp, expErr).Times(1)
		responseBody.EXPECT().Close().Return(expErr).Times(1)
		f.Logger.EXPECT().Warningf(gomock.Any(), expErr).Times(1)
		r, err := c.Update(f.Ctx, []byte(f.Key))
		mocks.CmpError(t, err, expErr)
		if !cmp.Equal(r, nil) {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("http call error is a failure (resp nil)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newHttpDataSourceFixture(ctrl)
		c := newTestableHttpDataSource(f)

		req := httptest.NewRequest(http.MethodPost, f.Target, nil)

		f.CreateRequest.EXPECT().Create(f.Ctx, []byte(f.Key), f.Url, http.MethodPost).Return(req, nil).Times(1)
		f.Do.EXPECT().Do(req).Return(nil, f.Error).Times(1)
		r, err := c.Update(f.Ctx, []byte(f.Key))
		mocks.CmpError(t, err, f.Error)
		if !cmp.Equal(r, nil) {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("http call error is a failure (resp body is nil)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newHttpDataSourceFixture(ctrl)
		c := newTestableHttpDataSource(f)

		req := httptest.NewRequest(http.MethodPost, f.Target, nil)
		respRec := httptest.NewRecorder()
		respRec.WriteHeader(http.StatusInternalServerError)
		_, err := respRec.WriteString(f.Data)
		if err != nil {
			panic(err)
		}
		respRec.Flush()
		resp := respRec.Result()
		resp.Body = nil

		f.CreateRequest.EXPECT().Create(f.Ctx, []byte(f.Key), f.Url, http.MethodPost).Return(req, nil).Times(1)
		f.Do.EXPECT().Do(req).Return(resp, f.Error).Times(1)
		r, err := c.Update(f.Ctx, []byte(f.Key))
		mocks.CmpError(t, err, f.Error)
		if !cmp.Equal(r, nil) {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("http ok creates response", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newHttpDataSourceFixture(ctrl)
		c := newTestableHttpDataSource(f)

		wrapResp := &DummyResponse{}
		req := httptest.NewRequest(http.MethodPost, f.Target, nil)
		respRec := httptest.NewRecorder()
		respRec.WriteHeader(http.StatusOK)
		_, err := respRec.WriteString(f.Data)
		if err != nil {
			panic(err)
		}
		respRec.Flush()
		resp := respRec.Result()
		resp.Body = nil

		f.CreateRequest.EXPECT().Create(f.Ctx, []byte(f.Key), f.Url, http.MethodPost).Return(req, nil).Times(1)
		f.Do.EXPECT().Do(req).Return(resp, nil).Times(1)
		f.CreateResponse.EXPECT().Create(resp).Return(wrapResp, nil).Times(1)
		r, err := c.Update(f.Ctx, []byte(f.Key))
		mocks.CmpError(t, err, nil)
		if !cmp.Equal(r, wrapResp) {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("http ok creates response (can be nil)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newHttpDataSourceFixture(ctrl)
		c := newTestableHttpDataSource(f)

		var wrapResp logic.Response = nil
		req := httptest.NewRequest(http.MethodPost, f.Target, nil)
		respRec := httptest.NewRecorder()
		respRec.WriteHeader(http.StatusOK)
		_, err := respRec.WriteString(f.Data)
		if err != nil {
			panic(err)
		}
		respRec.Flush()
		resp := respRec.Result()
		resp.Body = nil

		f.CreateRequest.EXPECT().Create(f.Ctx, []byte(f.Key), f.Url, http.MethodPost).Return(req, nil).Times(1)
		f.Do.EXPECT().Do(req).Return(resp, nil).Times(1)
		f.CreateResponse.EXPECT().Create(resp).Return(wrapResp, nil).Times(1)
		r, err := c.Update(f.Ctx, []byte(f.Key))
		mocks.CmpError(t, err, nil)
		if !cmp.Equal(r, wrapResp) {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("http ok creates response, if response is created with error - return resp and error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newHttpDataSourceFixture(ctrl)
		c := newTestableHttpDataSource(f)

		wrapResp := &DummyResponse{}
		req := httptest.NewRequest(http.MethodPost, f.Target, nil)
		respRec := httptest.NewRecorder()
		respRec.WriteHeader(http.StatusInternalServerError)
		_, err := respRec.WriteString(f.Data)
		if err != nil {
			panic(err)
		}
		respRec.Flush()
		resp := respRec.Result()
		resp.Body = nil

		f.CreateRequest.EXPECT().Create(f.Ctx, []byte(f.Key), f.Url, http.MethodPost).Return(req, nil).Times(1)
		f.Do.EXPECT().Do(req).Return(resp, nil).Times(1)
		f.CreateResponse.EXPECT().Create(resp).Return(wrapResp, f.Error).Times(1)
		r, err := c.Update(f.Ctx, []byte(f.Key))
		mocks.CmpError(t, err, f.Error)
		if !cmp.Equal(r, wrapResp) {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("http ok creates response, if response is created without error, but status != 200 returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newHttpDataSourceFixture(ctrl)
		c := newTestableHttpDataSource(f)

		wrapResp := &DummyResponse{}
		req := httptest.NewRequest(http.MethodPost, f.Target, nil)
		respRec := httptest.NewRecorder()
		respRec.WriteHeader(http.StatusInternalServerError)
		_, err := respRec.WriteString(f.Data)
		if err != nil {
			panic(err)
		}
		respRec.Flush()
		resp := respRec.Result()
		resp.Body = nil

		f.CreateRequest.EXPECT().Create(f.Ctx, []byte(f.Key), f.Url, http.MethodPost).Return(req, nil).Times(1)
		f.Do.EXPECT().Do(req).Return(resp, nil).Times(1)
		f.CreateResponse.EXPECT().Create(resp).Return(wrapResp, nil).Times(1)
		r, err := c.Update(f.Ctx, []byte(f.Key))
		if err == nil {
			t.Errorf("Failed. Code: %v. Error: %v", resp.StatusCode, err)
		}
		if !cmp.Equal(r, wrapResp) {
			t.Errorf("Expected correct response.")
		}
	})
}

func TestNewHttpDataSource(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	httpWrap := mock_std.NewMockHttpWrap(ctrl)
	dataSource := NewHttpDataSource(httpWrap.Do, "test_url")
	if dataSource == nil {
		t.Errorf("Factory returned nil")
	}
}

func TestNewDefaultHttpDataSource(t *testing.T) {
	dataSource := NewDefaultHttpDataSource("test_url")
	if dataSource == nil {
		t.Errorf("Factory returned nil")
	}
}
