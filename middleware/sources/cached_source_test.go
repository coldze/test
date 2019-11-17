package sources

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"

	"github.com/coldze/test/mocks"
	"github.com/coldze/test/utils"
)

const (
	key = "123"
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

type DummyResponse struct{}

func (d *DummyResponse) Write(w http.ResponseWriter) error {
	return nil
}

func TestCachedDataSource_Get(t *testing.T) {

	t.Run("cache error is logged, full chain executed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		dataSource := mocks.NewMockDataSource(ctrl)
		cache := mocks.NewMockCacheSource(ctrl)
		resp := &DummyResponse{}

		c := cachedDataSource{
			original: dataSource,
			cache:    cache,
		}

		ctx, logger := makeLoggerContext(ctrl)
		expErr := errors.New("Test")

		logger.EXPECT().Warningf(gomock.Any(), expErr).Times(1)
		cache.EXPECT().Get(key).Return(nil, expErr).Times(1)
		dataSource.EXPECT().Get(ctx, []byte(key)).Return(resp, nil).Times(1)
		cache.EXPECT().Insert(resp).Return(nil).Times(1)
		r, err := c.Get(ctx, []byte(key))
		cmperr(t, err, nil)
		if !cmp.Equal(r, resp) {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("cache miss triggers full chain", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		dataSource := mocks.NewMockDataSource(ctrl)
		cache := mocks.NewMockCacheSource(ctrl)
		resp := &DummyResponse{}

		c := cachedDataSource{
			original: dataSource,
			cache:    cache,
		}

		ctx, _ := makeLoggerContext(ctrl)

		cache.EXPECT().Get(key).Return(nil, nil).Times(1)
		dataSource.EXPECT().Get(ctx, []byte(key)).Return(resp, nil).Times(1)
		cache.EXPECT().Insert(resp).Return(nil).Times(1)
		r, err := c.Get(ctx, []byte(key))
		cmperr(t, err, nil)
		if !cmp.Equal(r, resp) {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("cache hit returns cached", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		dataSource := mocks.NewMockDataSource(ctrl)
		cache := mocks.NewMockCacheSource(ctrl)
		resp := &DummyResponse{}

		c := cachedDataSource{
			original: dataSource,
			cache:    cache,
		}
		ctx, _ := makeLoggerContext(ctrl)

		cache.EXPECT().Get(key).Return(resp, nil).Times(1)
		r, err := c.Get(ctx, []byte(key))
		cmperr(t, err, nil)
		if !cmp.Equal(r, resp) {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("main source error leads to error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		dataSource := mocks.NewMockDataSource(ctrl)
		cache := mocks.NewMockCacheSource(ctrl)
		resp := &DummyResponse{}

		c := cachedDataSource{
			original: dataSource,
			cache:    cache,
		}

		ctx, _ := makeLoggerContext(ctrl)
		expErr := errors.New("Test")

		cache.EXPECT().Get(key).Return(nil, nil).Times(1)
		dataSource.EXPECT().Get(ctx, []byte(key)).Return(resp, expErr).Times(1)
		r, err := c.Get(ctx, []byte(key))
		cmperr(t, err, expErr)

		if !cmp.Equal(r, resp) {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("inserting to cache failure is logged, call succeeds", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		dataSource := mocks.NewMockDataSource(ctrl)
		cache := mocks.NewMockCacheSource(ctrl)
		resp := &DummyResponse{}

		c := cachedDataSource{
			original: dataSource,
			cache:    cache,
		}

		ctx, logger := makeLoggerContext(ctrl)
		expErr := errors.New("Test")

		cache.EXPECT().Get(key).Return(nil, nil).Times(1)
		dataSource.EXPECT().Get(ctx, []byte(key)).Return(resp, nil).Times(1)
		cache.EXPECT().Insert(resp).Return(expErr).Times(1)
		logger.EXPECT().Warningf(gomock.Any(), expErr).Times(1)

		r, err := c.Get(ctx, []byte(key))
		cmperr(t, err, nil)

		if !cmp.Equal(r, resp) {
			t.Errorf("Expected correct response.")
		}
	})
}

func TestCachedDataSource_Create(t *testing.T) {
	t.Run("main source create error is a failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		dataSource := mocks.NewMockDataSource(ctrl)
		cache := mocks.NewMockCacheSource(ctrl)
		resp := &DummyResponse{}

		c := cachedDataSource{
			original: dataSource,
			cache:    cache,
		}

		ctx, _ := makeLoggerContext(ctrl)
		expErr := errors.New("Test")

		dataSource.EXPECT().Create(ctx, []byte(key)).Return(resp, expErr).Times(1)
		r, err := c.Create(ctx, []byte(key))
		cmperr(t, err, expErr)

		if !cmp.Equal(r, resp) {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("cache remove error is logged", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		dataSource := mocks.NewMockDataSource(ctrl)
		cache := mocks.NewMockCacheSource(ctrl)
		resp := &DummyResponse{}

		c := cachedDataSource{
			original: dataSource,
			cache:    cache,
		}

		ctx, logger := makeLoggerContext(ctrl)
		expErr := errors.New("Test")

		dataSource.EXPECT().Create(ctx, []byte(key)).Return(resp, nil).Times(1)
		cache.EXPECT().Remove(resp).Return(expErr).Times(1)
		logger.EXPECT().Warningf(gomock.Any(), expErr).Times(1)
		r, err := c.Create(ctx, []byte(key))
		cmperr(t, err, nil)

		if !cmp.Equal(r, resp) {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		dataSource := mocks.NewMockDataSource(ctrl)
		cache := mocks.NewMockCacheSource(ctrl)
		resp := &DummyResponse{}

		c := cachedDataSource{
			original: dataSource,
			cache:    cache,
		}

		ctx, _ := makeLoggerContext(ctrl)

		dataSource.EXPECT().Create(ctx, []byte(key)).Return(resp, nil).Times(1)
		cache.EXPECT().Remove(resp).Return(nil).Times(1)
		r, err := c.Create(ctx, []byte(key))
		cmperr(t, err, nil)

		if !cmp.Equal(r, resp) {
			t.Errorf("Expected correct response.")
		}
	})
}

func TestCachedDataSource_Update(t *testing.T) {
	t.Run("main source create error is a failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		dataSource := mocks.NewMockDataSource(ctrl)
		cache := mocks.NewMockCacheSource(ctrl)
		resp := &DummyResponse{}

		c := cachedDataSource{
			original: dataSource,
			cache:    cache,
		}

		ctx, _ := makeLoggerContext(ctrl)
		expErr := errors.New("Test")

		dataSource.EXPECT().Create(ctx, []byte(key)).Return(resp, expErr).Times(1)
		r, err := c.Update(ctx, []byte(key))
		cmperr(t, err, expErr)

		if !cmp.Equal(r, resp) {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("cache remove error is logged", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		dataSource := mocks.NewMockDataSource(ctrl)
		cache := mocks.NewMockCacheSource(ctrl)
		resp := &DummyResponse{}

		c := cachedDataSource{
			original: dataSource,
			cache:    cache,
		}

		ctx, logger := makeLoggerContext(ctrl)
		expErr := errors.New("Test")

		dataSource.EXPECT().Create(ctx, []byte(key)).Return(resp, nil).Times(1)
		cache.EXPECT().Remove(resp).Return(expErr).Times(1)
		logger.EXPECT().Warningf(gomock.Any(), expErr).Times(1)
		r, err := c.Update(ctx, []byte(key))
		cmperr(t, err, nil)

		if !cmp.Equal(r, resp) {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		dataSource := mocks.NewMockDataSource(ctrl)
		cache := mocks.NewMockCacheSource(ctrl)
		resp := &DummyResponse{}

		c := cachedDataSource{
			original: dataSource,
			cache:    cache,
		}

		ctx, _ := makeLoggerContext(ctrl)

		dataSource.EXPECT().Create(ctx, []byte(key)).Return(resp, nil).Times(1)
		cache.EXPECT().Remove(resp).Return(nil).Times(1)
		r, err := c.Update(ctx, []byte(key))
		cmperr(t, err, nil)

		if !cmp.Equal(r, resp) {
			t.Errorf("Expected correct response.")
		}
	})
}
