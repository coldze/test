package sources

import (
	"errors"
	"github.com/coldze/test/middleware"
	"github.com/coldze/test/mocks"
	"github.com/go-redis/redis"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"testing"
	"time"
)

const (
	ttl = 10 * time.Second
)

func TestRedisCacheSource_Get(t *testing.T) {
	data := "some random data"

	t.Run("get error is a failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		redisWrap := mocks.NewMockRedisWrap(ctrl)
		createResponse := mocks.NewMockResponseFactory(ctrl)
		dataParser := mocks.NewMockDataParser(ctrl)
		dataBuilderFactory := mocks.NewMockDataBuilderFactory(ctrl)

		c := redisCacheSource{
			cache:          redisWrap,
			ttl:            ttl,
			parse:          dataParser.Parse,
			createBuilder:  dataBuilderFactory.Create,
			createResponse: createResponse.Create,
		}

		expErr := errors.New("Test")

		redisWrap.EXPECT().Get(key).Return(nil, expErr).Times(1)
		r, err := c.Get(key)
		cmperr(t, err, expErr)
		if !cmp.Equal(r, nil) {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("cache miss is not an error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		redisWrap := mocks.NewMockRedisWrap(ctrl)
		createResponse := mocks.NewMockResponseFactory(ctrl)
		dataParser := mocks.NewMockDataParser(ctrl)
		dataBuilderFactory := mocks.NewMockDataBuilderFactory(ctrl)

		c := redisCacheSource{
			cache:          redisWrap,
			ttl:            ttl,
			parse:          dataParser.Parse,
			createBuilder:  dataBuilderFactory.Create,
			createResponse: createResponse.Create,
		}

		redisWrap.EXPECT().Get(key).Return(nil, redis.Nil).Times(1)
		r, err := c.Get(key)
		cmperr(t, err, nil)
		if !cmp.Equal(r, nil) {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("not string response is not expected", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		redisWrap := mocks.NewMockRedisWrap(ctrl)
		createResponse := mocks.NewMockResponseFactory(ctrl)
		dataParser := mocks.NewMockDataParser(ctrl)
		dataBuilderFactory := mocks.NewMockDataBuilderFactory(ctrl)

		c := redisCacheSource{
			cache:          redisWrap,
			ttl:            ttl,
			parse:          dataParser.Parse,
			createBuilder:  dataBuilderFactory.Create,
			createResponse: createResponse.Create,
		}

		redisWrap.EXPECT().Get(key).Return([]byte(data), nil).Times(1)
		r, err := c.Get(key)
		if err == nil {
			t.Errorf("Error is nil.")
		}
		if !cmp.Equal(r, nil) {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("string response is passed to create response", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		resp := &DummyResponse{}
		redisWrap := mocks.NewMockRedisWrap(ctrl)
		createResponse := mocks.NewMockResponseFactory(ctrl)
		dataParser := mocks.NewMockDataParser(ctrl)
		dataBuilderFactory := mocks.NewMockDataBuilderFactory(ctrl)

		c := redisCacheSource{
			cache:          redisWrap,
			ttl:            ttl,
			parse:          dataParser.Parse,
			createBuilder:  dataBuilderFactory.Create,
			createResponse: createResponse.Create,
		}

		redisWrap.EXPECT().Get(key).Return(data, nil).Times(1)
		createResponse.EXPECT().Create([]byte(data)).Return(resp, nil).Times(1)
		r, err := c.Get(key)
		cmperr(t, err, nil)
		if !cmp.Equal(r, resp) {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("string response is passed to create response (result can be nil)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		redisWrap := mocks.NewMockRedisWrap(ctrl)
		createResponse := mocks.NewMockResponseFactory(ctrl)
		dataParser := mocks.NewMockDataParser(ctrl)
		dataBuilderFactory := mocks.NewMockDataBuilderFactory(ctrl)

		c := redisCacheSource{
			cache:          redisWrap,
			ttl:            ttl,
			parse:          dataParser.Parse,
			createBuilder:  dataBuilderFactory.Create,
			createResponse: createResponse.Create,
		}

		redisWrap.EXPECT().Get(key).Return(data, nil).Times(1)
		createResponse.EXPECT().Create([]byte(data)).Return(nil, nil).Times(1)
		r, err := c.Get(key)
		cmperr(t, err, nil)
		if !cmp.Equal(r, nil) {
			t.Errorf("Expected correct response.")
		}
	})

	t.Run("create response error is a failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		redisWrap := mocks.NewMockRedisWrap(ctrl)
		createResponse := mocks.NewMockResponseFactory(ctrl)
		dataParser := mocks.NewMockDataParser(ctrl)
		dataBuilderFactory := mocks.NewMockDataBuilderFactory(ctrl)

		c := redisCacheSource{
			cache:          redisWrap,
			ttl:            ttl,
			parse:          dataParser.Parse,
			createBuilder:  dataBuilderFactory.Create,
			createResponse: createResponse.Create,
		}

		expErr := errors.New("Test")

		redisWrap.EXPECT().Get(key).Return(data, nil).Times(1)
		createResponse.EXPECT().Create([]byte(data)).Return(nil, expErr).Times(1)
		r, err := c.Get(key)
		cmperr(t, err, expErr)
		if !cmp.Equal(r, nil) {
			t.Errorf("Expected correct response.")
		}
	})
}

func TestRedisCacheSource_Remove(t *testing.T) {
	data := "some random data"
	contact := middleware.Contact{ID: key}

	t.Run("nil builder is a failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		resp := mocks.NewMockResponse(ctrl)
		redisWrap := mocks.NewMockRedisWrap(ctrl)
		createResponse := mocks.NewMockResponseFactory(ctrl)
		dataParser := mocks.NewMockDataParser(ctrl)
		dataBuilderFactory := mocks.NewMockDataBuilderFactory(ctrl)

		c := redisCacheSource{
			cache:          redisWrap,
			ttl:            ttl,
			parse:          dataParser.Parse,
			createBuilder:  dataBuilderFactory.Create,
			createResponse: createResponse.Create,
		}

		dataBuilderFactory.EXPECT().Create().Return(nil).Times(1)

		err := c.Remove(resp)
		if err == nil {
			t.Errorf("Error is nil")
		}
	})

	t.Run("error while writing response is a failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		resp := mocks.NewMockResponse(ctrl)
		builder := mocks.NewMockDataBuilder(ctrl)
		redisWrap := mocks.NewMockRedisWrap(ctrl)
		createResponse := mocks.NewMockResponseFactory(ctrl)
		dataParser := mocks.NewMockDataParser(ctrl)
		dataBuilderFactory := mocks.NewMockDataBuilderFactory(ctrl)

		c := redisCacheSource{
			cache:          redisWrap,
			ttl:            ttl,
			parse:          dataParser.Parse,
			createBuilder:  dataBuilderFactory.Create,
			createResponse: createResponse.Create,
		}

		expErr := errors.New("Test")

		dataBuilderFactory.EXPECT().Create().Return(builder).Times(1)
		resp.EXPECT().Write(builder).Return(expErr).Times(1)

		err := c.Remove(resp)
		cmperr(t, err, expErr)
	})

	t.Run("builder error is a failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		resp := mocks.NewMockResponse(ctrl)
		builder := mocks.NewMockDataBuilder(ctrl)
		redisWrap := mocks.NewMockRedisWrap(ctrl)
		createResponse := mocks.NewMockResponseFactory(ctrl)
		dataParser := mocks.NewMockDataParser(ctrl)
		dataBuilderFactory := mocks.NewMockDataBuilderFactory(ctrl)

		c := redisCacheSource{
			cache:          redisWrap,
			ttl:            ttl,
			parse:          dataParser.Parse,
			createBuilder:  dataBuilderFactory.Create,
			createResponse: createResponse.Create,
		}

		expErr := errors.New("Test")

		dataBuilderFactory.EXPECT().Create().Return(builder).Times(1)
		resp.EXPECT().Write(builder).Return(nil).Times(1)
		builder.EXPECT().Build().Return([]byte(data), expErr).Times(1)

		err := c.Remove(resp)
		cmperr(t, err, expErr)
	})

	t.Run("parser error is a failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		resp := mocks.NewMockResponse(ctrl)
		builder := mocks.NewMockDataBuilder(ctrl)
		redisWrap := mocks.NewMockRedisWrap(ctrl)
		createResponse := mocks.NewMockResponseFactory(ctrl)
		dataParser := mocks.NewMockDataParser(ctrl)
		dataBuilderFactory := mocks.NewMockDataBuilderFactory(ctrl)

		c := redisCacheSource{
			cache:          redisWrap,
			ttl:            ttl,
			parse:          dataParser.Parse,
			createBuilder:  dataBuilderFactory.Create,
			createResponse: createResponse.Create,
		}

		expErr := errors.New("Test")

		dataBuilderFactory.EXPECT().Create().Return(builder).Times(1)
		resp.EXPECT().Write(builder).Return(nil).Times(1)
		builder.EXPECT().Build().Return([]byte(data), nil).Times(1)
		dataParser.EXPECT().Create([]byte(data)).Return(contact, expErr).Times(1)

		err := c.Remove(resp)
		cmperr(t, err, expErr)
	})

	t.Run("remove error is a failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		resp := mocks.NewMockResponse(ctrl)
		builder := mocks.NewMockDataBuilder(ctrl)
		redisWrap := mocks.NewMockRedisWrap(ctrl)
		createResponse := mocks.NewMockResponseFactory(ctrl)
		dataParser := mocks.NewMockDataParser(ctrl)
		dataBuilderFactory := mocks.NewMockDataBuilderFactory(ctrl)

		c := redisCacheSource{
			cache:          redisWrap,
			ttl:            ttl,
			parse:          dataParser.Parse,
			createBuilder:  dataBuilderFactory.Create,
			createResponse: createResponse.Create,
		}

		expErr := errors.New("Test")

		dataBuilderFactory.EXPECT().Create().Return(builder).Times(1)
		resp.EXPECT().Write(builder).Return(nil).Times(1)
		builder.EXPECT().Build().Return([]byte(data), nil).Times(1)
		dataParser.EXPECT().Create([]byte(data)).Return(contact, nil).Times(1)
		redisWrap.EXPECT().Del(contact.ID).Return(expErr)

		err := c.Remove(resp)
		cmperr(t, err, expErr)
	})

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		resp := mocks.NewMockResponse(ctrl)
		builder := mocks.NewMockDataBuilder(ctrl)
		redisWrap := mocks.NewMockRedisWrap(ctrl)
		createResponse := mocks.NewMockResponseFactory(ctrl)
		dataParser := mocks.NewMockDataParser(ctrl)
		dataBuilderFactory := mocks.NewMockDataBuilderFactory(ctrl)

		c := redisCacheSource{
			cache:          redisWrap,
			ttl:            ttl,
			parse:          dataParser.Parse,
			createBuilder:  dataBuilderFactory.Create,
			createResponse: createResponse.Create,
		}

		dataBuilderFactory.EXPECT().Create().Return(builder).Times(1)
		resp.EXPECT().Write(builder).Return(nil).Times(1)
		builder.EXPECT().Build().Return([]byte(data), nil).Times(1)
		dataParser.EXPECT().Create([]byte(data)).Return(contact, nil).Times(1)
		redisWrap.EXPECT().Del(contact.ID).Return(nil)

		err := c.Remove(resp)
		cmperr(t, err, nil)
	})
}

func TestRedisCacheSource_Insert(t *testing.T) {
	data := "some random data"
	contact := middleware.Contact{ID: key}

	t.Run("nil builder is a failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		resp := mocks.NewMockResponse(ctrl)
		redisWrap := mocks.NewMockRedisWrap(ctrl)
		createResponse := mocks.NewMockResponseFactory(ctrl)
		dataParser := mocks.NewMockDataParser(ctrl)
		dataBuilderFactory := mocks.NewMockDataBuilderFactory(ctrl)

		c := redisCacheSource{
			cache:          redisWrap,
			ttl:            ttl,
			parse:          dataParser.Parse,
			createBuilder:  dataBuilderFactory.Create,
			createResponse: createResponse.Create,
		}

		dataBuilderFactory.EXPECT().Create().Return(nil).Times(1)

		err := c.Insert(resp)
		if err == nil {
			t.Errorf("Error is nil")
		}
	})

	t.Run("error while writing response is a failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		resp := mocks.NewMockResponse(ctrl)
		builder := mocks.NewMockDataBuilder(ctrl)
		redisWrap := mocks.NewMockRedisWrap(ctrl)
		createResponse := mocks.NewMockResponseFactory(ctrl)
		dataParser := mocks.NewMockDataParser(ctrl)
		dataBuilderFactory := mocks.NewMockDataBuilderFactory(ctrl)

		c := redisCacheSource{
			cache:          redisWrap,
			ttl:            ttl,
			parse:          dataParser.Parse,
			createBuilder:  dataBuilderFactory.Create,
			createResponse: createResponse.Create,
		}

		expErr := errors.New("Test")

		dataBuilderFactory.EXPECT().Create().Return(builder).Times(1)
		resp.EXPECT().Write(builder).Return(expErr).Times(1)

		err := c.Insert(resp)
		cmperr(t, err, expErr)
	})

	t.Run("builder error is a failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		resp := mocks.NewMockResponse(ctrl)
		builder := mocks.NewMockDataBuilder(ctrl)
		redisWrap := mocks.NewMockRedisWrap(ctrl)
		createResponse := mocks.NewMockResponseFactory(ctrl)
		dataParser := mocks.NewMockDataParser(ctrl)
		dataBuilderFactory := mocks.NewMockDataBuilderFactory(ctrl)

		c := redisCacheSource{
			cache:          redisWrap,
			ttl:            ttl,
			parse:          dataParser.Parse,
			createBuilder:  dataBuilderFactory.Create,
			createResponse: createResponse.Create,
		}

		expErr := errors.New("Test")

		dataBuilderFactory.EXPECT().Create().Return(builder).Times(1)
		resp.EXPECT().Write(builder).Return(nil).Times(1)
		builder.EXPECT().Build().Return([]byte(data), expErr).Times(1)

		err := c.Insert(resp)
		cmperr(t, err, expErr)
	})

	t.Run("parser error is a failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		resp := mocks.NewMockResponse(ctrl)
		builder := mocks.NewMockDataBuilder(ctrl)
		redisWrap := mocks.NewMockRedisWrap(ctrl)
		createResponse := mocks.NewMockResponseFactory(ctrl)
		dataParser := mocks.NewMockDataParser(ctrl)
		dataBuilderFactory := mocks.NewMockDataBuilderFactory(ctrl)

		c := redisCacheSource{
			cache:          redisWrap,
			ttl:            ttl,
			parse:          dataParser.Parse,
			createBuilder:  dataBuilderFactory.Create,
			createResponse: createResponse.Create,
		}

		expErr := errors.New("Test")

		dataBuilderFactory.EXPECT().Create().Return(builder).Times(1)
		resp.EXPECT().Write(builder).Return(nil).Times(1)
		builder.EXPECT().Build().Return([]byte(data), nil).Times(1)
		dataParser.EXPECT().Create([]byte(data)).Return(contact, expErr).Times(1)

		err := c.Insert(resp)
		cmperr(t, err, expErr)
	})

	t.Run("insert error is a failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		resp := mocks.NewMockResponse(ctrl)
		builder := mocks.NewMockDataBuilder(ctrl)
		redisWrap := mocks.NewMockRedisWrap(ctrl)
		createResponse := mocks.NewMockResponseFactory(ctrl)
		dataParser := mocks.NewMockDataParser(ctrl)
		dataBuilderFactory := mocks.NewMockDataBuilderFactory(ctrl)

		c := redisCacheSource{
			cache:          redisWrap,
			ttl:            ttl,
			parse:          dataParser.Parse,
			createBuilder:  dataBuilderFactory.Create,
			createResponse: createResponse.Create,
		}

		expErr := errors.New("Test")

		dataBuilderFactory.EXPECT().Create().Return(builder).Times(1)
		resp.EXPECT().Write(builder).Return(nil).Times(1)
		builder.EXPECT().Build().Return([]byte(data), nil).Times(1)
		dataParser.EXPECT().Create([]byte(data)).Return(contact, nil).Times(1)
		redisWrap.EXPECT().Set(contact.ID, []byte(data), ttl).Return(expErr)

		err := c.Insert(resp)
		cmperr(t, err, expErr)
	})

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		resp := mocks.NewMockResponse(ctrl)
		builder := mocks.NewMockDataBuilder(ctrl)
		redisWrap := mocks.NewMockRedisWrap(ctrl)
		createResponse := mocks.NewMockResponseFactory(ctrl)
		dataParser := mocks.NewMockDataParser(ctrl)
		dataBuilderFactory := mocks.NewMockDataBuilderFactory(ctrl)

		c := redisCacheSource{
			cache:          redisWrap,
			ttl:            ttl,
			parse:          dataParser.Parse,
			createBuilder:  dataBuilderFactory.Create,
			createResponse: createResponse.Create,
		}

		dataBuilderFactory.EXPECT().Create().Return(builder).Times(1)
		resp.EXPECT().Write(builder).Return(nil).Times(1)
		builder.EXPECT().Build().Return([]byte(data), nil).Times(1)
		dataParser.EXPECT().Create([]byte(data)).Return(contact, nil).Times(1)
		redisWrap.EXPECT().Set(contact.ID, []byte(data), ttl).Return(nil)

		err := c.Insert(resp)
		cmperr(t, err, nil)
	})
}
