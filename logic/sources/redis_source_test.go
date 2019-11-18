package sources

import (
	"errors"
	"github.com/coldze/test/logic"
	"github.com/coldze/test/mocks"
	"github.com/coldze/test/mocks/mock_sources"
	"github.com/go-redis/redis"
	"github.com/golang/mock/gomock"
	"testing"
	"time"
)

type redisCacheFixture struct {
	Contact            logic.Contact
	Data               string
	Key                string
	Error              error
	Ttl                time.Duration
	Response           *mocks.MockResponse
	DataBuilder        *mock_sources.MockDataBuilder
	RedisWrap          *mock_sources.MockRedisWrap
	CreateResponse     *mock_sources.MockResponseFactory
	DataParser         *mock_sources.MockDataParser
	DataBuilderFactory *mock_sources.MockDataBuilderFactory
}

func newRedisCacheFixture(ctrl *gomock.Controller) *redisCacheFixture {
	return &redisCacheFixture{
		Contact:            logic.Contact{ID: "some random key"},
		Data:               "some random data",
		Key:                "some test key",
		Error:              errors.New("some test error"),
		Ttl:                1 * time.Second,
		Response:           mocks.NewMockResponse(ctrl),
		DataBuilder:        mock_sources.NewMockDataBuilder(ctrl),
		RedisWrap:          mock_sources.NewMockRedisWrap(ctrl),
		CreateResponse:     mock_sources.NewMockResponseFactory(ctrl),
		DataParser:         mock_sources.NewMockDataParser(ctrl),
		DataBuilderFactory: mock_sources.NewMockDataBuilderFactory(ctrl),
	}
}

func newRedisCacheSource(f *redisCacheFixture) *redisCacheSource {
	return &redisCacheSource{
		cache:          f.RedisWrap,
		ttl:            f.Ttl,
		parse:          f.DataParser.Parse,
		createBuilder:  f.DataBuilderFactory.Create,
		createResponse: f.CreateResponse.Create,
	}
}

func TestRedisCacheSource_Get(t *testing.T) {

	t.Run("get error is a failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newRedisCacheFixture(ctrl)
		c := newRedisCacheSource(f)

		f.RedisWrap.EXPECT().Get(f.Key).Return(nil, f.Error).Times(1)
		r, err := c.Get(f.Key)
		mocks.CmpError(t, err, f.Error)
		if r != nil {
			t.Errorf("Response should be nil.")
		}
	})

	t.Run("cache miss is not an error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newRedisCacheFixture(ctrl)
		c := newRedisCacheSource(f)

		f.RedisWrap.EXPECT().Get(f.Key).Return(nil, redis.Nil).Times(1)
		r, err := c.Get(f.Key)
		mocks.CmpError(t, err, nil)
		if r != nil {
			t.Errorf("Response should be nil.")
		}
	})

	t.Run("not string response is not expected", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newRedisCacheFixture(ctrl)
		c := newRedisCacheSource(f)

		f.RedisWrap.EXPECT().Get(f.Key).Return([]byte(f.Data), nil).Times(1)
		r, err := c.Get(f.Key)
		if err == nil {
			t.Errorf("Error is nil.")
		}
		if r != nil {
			t.Errorf("Response should be nil.")
		}
	})

	t.Run("string response is passed to create response", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newRedisCacheFixture(ctrl)
		c := newRedisCacheSource(f)

		f.RedisWrap.EXPECT().Get(f.Key).Return(f.Data, nil).Times(1)
		f.CreateResponse.EXPECT().Create([]byte(f.Data)).Return(f.Response, nil).Times(1)
		r, err := c.Get(f.Key)
		mocks.CmpError(t, err, nil)
		if r != f.Response {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("string response is passed to create response (result can be nil)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newRedisCacheFixture(ctrl)
		c := newRedisCacheSource(f)

		f.RedisWrap.EXPECT().Get(f.Key).Return(f.Data, nil).Times(1)
		f.CreateResponse.EXPECT().Create([]byte(f.Data)).Return(nil, nil).Times(1)
		r, err := c.Get(f.Key)
		mocks.CmpError(t, err, nil)
		if r != nil {
			t.Errorf("Response should be nil.")
		}
	})

	t.Run("create response error is a failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newRedisCacheFixture(ctrl)
		c := newRedisCacheSource(f)

		f.RedisWrap.EXPECT().Get(f.Key).Return(f.Data, nil).Times(1)
		f.CreateResponse.EXPECT().Create([]byte(f.Data)).Return(nil, f.Error).Times(1)
		r, err := c.Get(f.Key)
		mocks.CmpError(t, err, f.Error)
		if r != nil {
			t.Errorf("Response should be nil.")
		}
	})
}

func TestRedisCacheSource_Remove(t *testing.T) {

	t.Run("nil builder is a failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newRedisCacheFixture(ctrl)
		c := newRedisCacheSource(f)

		f.DataBuilderFactory.EXPECT().Create().Return(nil).Times(1)

		err := c.Remove(f.Response)
		if err == nil {
			t.Errorf("Error is nil")
		}
	})

	t.Run("error while writing response is a failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newRedisCacheFixture(ctrl)
		c := newRedisCacheSource(f)

		f.DataBuilderFactory.EXPECT().Create().Return(f.DataBuilder).Times(1)
		f.Response.EXPECT().Write(f.DataBuilder).Return(f.Error).Times(1)

		err := c.Remove(f.Response)
		mocks.CmpError(t, err, f.Error)
	})

	t.Run("builder error is a failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newRedisCacheFixture(ctrl)
		c := newRedisCacheSource(f)

		f.DataBuilderFactory.EXPECT().Create().Return(f.DataBuilder).Times(1)
		f.Response.EXPECT().Write(f.DataBuilder).Return(nil).Times(1)
		f.DataBuilder.EXPECT().Build().Return([]byte(f.Data), f.Error).Times(1)

		err := c.Remove(f.Response)
		mocks.CmpError(t, err, f.Error)
	})

	t.Run("parser error is a failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newRedisCacheFixture(ctrl)
		c := newRedisCacheSource(f)

		f.DataBuilderFactory.EXPECT().Create().Return(f.DataBuilder).Times(1)
		f.Response.EXPECT().Write(f.DataBuilder).Return(nil).Times(1)
		f.DataBuilder.EXPECT().Build().Return([]byte(f.Data), nil).Times(1)
		f.DataParser.EXPECT().Create([]byte(f.Data)).Return(f.Contact, f.Error).Times(1)

		err := c.Remove(f.Response)
		mocks.CmpError(t, err, f.Error)
	})

	t.Run("remove error is a failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newRedisCacheFixture(ctrl)
		c := newRedisCacheSource(f)

		f.DataBuilderFactory.EXPECT().Create().Return(f.DataBuilder).Times(1)
		f.Response.EXPECT().Write(f.DataBuilder).Return(nil).Times(1)
		f.DataBuilder.EXPECT().Build().Return([]byte(f.Data), nil).Times(1)
		f.DataParser.EXPECT().Create([]byte(f.Data)).Return(f.Contact, nil).Times(1)
		f.RedisWrap.EXPECT().Del(f.Contact.ID).Return(f.Error)

		err := c.Remove(f.Response)
		mocks.CmpError(t, err, f.Error)
	})

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newRedisCacheFixture(ctrl)
		c := newRedisCacheSource(f)

		f.DataBuilderFactory.EXPECT().Create().Return(f.DataBuilder).Times(1)
		f.Response.EXPECT().Write(f.DataBuilder).Return(nil).Times(1)
		f.DataBuilder.EXPECT().Build().Return([]byte(f.Data), nil).Times(1)
		f.DataParser.EXPECT().Create([]byte(f.Data)).Return(f.Contact, nil).Times(1)
		f.RedisWrap.EXPECT().Del(f.Contact.ID).Return(nil)

		err := c.Remove(f.Response)
		mocks.CmpError(t, err, nil)
	})
}

func TestRedisCacheSource_Insert(t *testing.T) {
	t.Run("nil builder is a failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newRedisCacheFixture(ctrl)
		c := newRedisCacheSource(f)

		f.DataBuilderFactory.EXPECT().Create().Return(nil).Times(1)

		err := c.Insert(f.Response)
		if err == nil {
			t.Errorf("Error is nil")
		}
	})

	t.Run("error while writing response is a failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newRedisCacheFixture(ctrl)
		c := newRedisCacheSource(f)

		f.DataBuilderFactory.EXPECT().Create().Return(f.DataBuilder).Times(1)
		f.Response.EXPECT().Write(f.DataBuilder).Return(f.Error).Times(1)

		err := c.Insert(f.Response)
		mocks.CmpError(t, err, f.Error)
	})

	t.Run("builder error is a failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newRedisCacheFixture(ctrl)
		c := newRedisCacheSource(f)

		f.DataBuilderFactory.EXPECT().Create().Return(f.DataBuilder).Times(1)
		f.Response.EXPECT().Write(f.DataBuilder).Return(nil).Times(1)
		f.DataBuilder.EXPECT().Build().Return([]byte(f.Data), f.Error).Times(1)

		err := c.Insert(f.Response)
		mocks.CmpError(t, err, f.Error)
	})

	t.Run("parser error is a failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newRedisCacheFixture(ctrl)
		c := newRedisCacheSource(f)

		f.DataBuilderFactory.EXPECT().Create().Return(f.DataBuilder).Times(1)
		f.Response.EXPECT().Write(f.DataBuilder).Return(nil).Times(1)
		f.DataBuilder.EXPECT().Build().Return([]byte(f.Data), nil).Times(1)
		f.DataParser.EXPECT().Create([]byte(f.Data)).Return(f.Contact, f.Error).Times(1)

		err := c.Insert(f.Response)
		mocks.CmpError(t, err, f.Error)
	})

	t.Run("insert error is a failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newRedisCacheFixture(ctrl)
		c := newRedisCacheSource(f)

		f.DataBuilderFactory.EXPECT().Create().Return(f.DataBuilder).Times(1)
		f.Response.EXPECT().Write(f.DataBuilder).Return(nil).Times(1)
		f.DataBuilder.EXPECT().Build().Return([]byte(f.Data), nil).Times(1)
		f.DataParser.EXPECT().Create([]byte(f.Data)).Return(f.Contact, nil).Times(1)
		f.RedisWrap.EXPECT().Set(f.Contact.ID, []byte(f.Data), f.Ttl).Return(f.Error)

		err := c.Insert(f.Response)
		mocks.CmpError(t, err, f.Error)
	})

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newRedisCacheFixture(ctrl)
		c := newRedisCacheSource(f)

		f.DataBuilderFactory.EXPECT().Create().Return(f.DataBuilder).Times(1)
		f.Response.EXPECT().Write(f.DataBuilder).Return(nil).Times(1)
		f.DataBuilder.EXPECT().Build().Return([]byte(f.Data), nil).Times(1)
		f.DataParser.EXPECT().Create([]byte(f.Data)).Return(f.Contact, nil).Times(1)
		f.RedisWrap.EXPECT().Set(f.Contact.ID, []byte(f.Data), f.Ttl).Return(nil)

		err := c.Insert(f.Response)
		mocks.CmpError(t, err, nil)
	})
}

func TestNewRedisCacheSource(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	redisWrap := mock_sources.NewMockRedisWrap(ctrl)

	res := NewRedisCacheSource(redisWrap, 1*time.Second)
	if res == nil {
		t.Errorf("Factory returns nil")
	}
}
