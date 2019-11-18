package sources

import (
	"context"
	"errors"
	"github.com/coldze/test/logic"
	"github.com/coldze/test/mocks"
	"github.com/coldze/test/mocks/mock_logs"
	"github.com/coldze/test/mocks/mock_sources"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"

	"github.com/coldze/test/utils"
)

func makeLoggerContext(controller *gomock.Controller) (context.Context, *mock_logs.MockLogger) {
	logger := mock_logs.NewMockLogger(controller)
	return utils.SetLogger(context.Background(), logger), logger
}

type DummyResponse struct{}

func (d *DummyResponse) Write(w http.ResponseWriter) error {
	return nil
}

type cacheSourceFixture struct {
	Key        string
	Error      error
	Ctx        context.Context
	Logger     *mock_logs.MockLogger
	Cache      *mock_sources.MockCacheSource
	DataSource *mock_sources.MockDataSource
	Response   logic.Response
}

func newCacheSourceFixture(ctrl *gomock.Controller) *cacheSourceFixture {
	ctx, logger := makeLoggerContext(ctrl)
	return &cacheSourceFixture{
		Key:        "some test key",
		Error:      errors.New("Some test error"),
		Ctx:        ctx,
		Logger:     logger,
		Cache:      mock_sources.NewMockCacheSource(ctrl),
		DataSource: mock_sources.NewMockDataSource(ctrl),
		Response:   &DummyResponse{},
	}
}

func newTestableCachedDataSource(f *cacheSourceFixture) *cachedDataSource {
	return &cachedDataSource{
		original: f.DataSource,
		cache:    f.Cache,
	}
}

func TestCachedDataSource_Get(t *testing.T) {

	t.Run("cache error is logged, full chain executed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newCacheSourceFixture(ctrl)
		c := newTestableCachedDataSource(f)

		f.Logger.EXPECT().Warningf(gomock.Any(), f.Error).Times(1)
		f.Cache.EXPECT().Get(f.Key).Return(nil, f.Error).Times(1)
		f.DataSource.EXPECT().Get(f.Ctx, []byte(f.Key)).Return(f.Response, nil).Times(1)
		f.Cache.EXPECT().Insert(f.Response).Return(nil).Times(1)
		r, err := c.Get(f.Ctx, []byte(f.Key))
		mocks.CmpError(t, err, nil)
		if !cmp.Equal(r, f.Response) {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("cache miss triggers full chain", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newCacheSourceFixture(ctrl)
		c := newTestableCachedDataSource(f)

		f.Cache.EXPECT().Get(f.Key).Return(nil, nil).Times(1)
		f.DataSource.EXPECT().Get(f.Ctx, []byte(f.Key)).Return(f.Response, nil).Times(1)
		f.Cache.EXPECT().Insert(f.Response).Return(nil).Times(1)
		r, err := c.Get(f.Ctx, []byte(f.Key))
		mocks.CmpError(t, err, nil)
		if !cmp.Equal(r, f.Response) {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("cache hit returns cached", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newCacheSourceFixture(ctrl)
		c := newTestableCachedDataSource(f)

		f.Cache.EXPECT().Get(f.Key).Return(f.Response, nil).Times(1)
		r, err := c.Get(f.Ctx, []byte(f.Key))
		mocks.CmpError(t, err, nil)
		if !cmp.Equal(r, f.Response) {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("main source error leads to error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newCacheSourceFixture(ctrl)
		c := newTestableCachedDataSource(f)

		f.Cache.EXPECT().Get(f.Key).Return(nil, nil).Times(1)
		f.DataSource.EXPECT().Get(f.Ctx, []byte(f.Key)).Return(f.Response, f.Error).Times(1)
		r, err := c.Get(f.Ctx, []byte(f.Key))
		mocks.CmpError(t, err, f.Error)

		if !cmp.Equal(r, f.Response) {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("inserting to cache failure is logged, call succeeds", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newCacheSourceFixture(ctrl)
		c := newTestableCachedDataSource(f)

		f.Cache.EXPECT().Get(f.Key).Return(nil, nil).Times(1)
		f.DataSource.EXPECT().Get(f.Ctx, []byte(f.Key)).Return(f.Response, nil).Times(1)
		f.Cache.EXPECT().Insert(f.Response).Return(f.Error).Times(1)
		f.Logger.EXPECT().Warningf(gomock.Any(), f.Error).Times(1)

		r, err := c.Get(f.Ctx, []byte(f.Key))
		mocks.CmpError(t, err, nil)

		if !cmp.Equal(r, f.Response) {
			t.Errorf("Expected correct response.")
		}
	})
}

func TestCachedDataSource_Create(t *testing.T) {
	t.Run("main source create error is a failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newCacheSourceFixture(ctrl)
		c := newTestableCachedDataSource(f)

		f.DataSource.EXPECT().Create(f.Ctx, []byte(f.Key)).Return(f.Response, f.Error).Times(1)
		r, err := c.Create(f.Ctx, []byte(f.Key))
		mocks.CmpError(t, err, f.Error)

		if !cmp.Equal(r, f.Response) {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("cache remove error is logged", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newCacheSourceFixture(ctrl)
		c := newTestableCachedDataSource(f)

		f.DataSource.EXPECT().Create(f.Ctx, []byte(f.Key)).Return(f.Response, nil).Times(1)
		f.Cache.EXPECT().Remove(f.Response).Return(f.Error).Times(1)
		f.Logger.EXPECT().Warningf(gomock.Any(), f.Error).Times(1)
		r, err := c.Create(f.Ctx, []byte(f.Key))
		mocks.CmpError(t, err, nil)

		if !cmp.Equal(r, f.Response) {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newCacheSourceFixture(ctrl)
		c := newTestableCachedDataSource(f)

		f.DataSource.EXPECT().Create(f.Ctx, []byte(f.Key)).Return(f.Response, nil).Times(1)
		f.Cache.EXPECT().Remove(f.Response).Return(nil).Times(1)
		r, err := c.Create(f.Ctx, []byte(f.Key))
		mocks.CmpError(t, err, nil)

		if !cmp.Equal(r, f.Response) {
			t.Errorf("Expected correct response.")
		}
	})
}

func TestCachedDataSource_Update(t *testing.T) {
	t.Run("main source create error is a failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newCacheSourceFixture(ctrl)
		c := newTestableCachedDataSource(f)

		f.DataSource.EXPECT().Create(f.Ctx, []byte(f.Key)).Return(f.Response, f.Error).Times(1)
		r, err := c.Update(f.Ctx, []byte(f.Key))
		mocks.CmpError(t, err, f.Error)

		if !cmp.Equal(r, f.Response) {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("cache remove error is logged", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newCacheSourceFixture(ctrl)
		c := newTestableCachedDataSource(f)

		f.DataSource.EXPECT().Create(f.Ctx, []byte(f.Key)).Return(f.Response, nil).Times(1)
		f.Cache.EXPECT().Remove(f.Response).Return(f.Error).Times(1)
		f.Logger.EXPECT().Warningf(gomock.Any(), f.Error).Times(1)
		r, err := c.Update(f.Ctx, []byte(f.Key))
		mocks.CmpError(t, err, nil)

		if !cmp.Equal(r, f.Response) {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newCacheSourceFixture(ctrl)
		c := newTestableCachedDataSource(f)

		f.DataSource.EXPECT().Create(f.Ctx, []byte(f.Key)).Return(f.Response, nil).Times(1)
		f.Cache.EXPECT().Remove(f.Response).Return(nil).Times(1)
		r, err := c.Update(f.Ctx, []byte(f.Key))
		mocks.CmpError(t, err, nil)

		if !cmp.Equal(r, f.Response) {
			t.Errorf("Expected correct response.")
		}
	})
}

func TestNewCachedDataSource(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dataSource := mock_sources.NewMockDataSource(ctrl)
	cache := mock_sources.NewMockCacheSource(ctrl)

	res := NewCachedDataSource(dataSource, cache)
	if res == nil {
		t.Errorf("Factory returns nil")
	}
}
