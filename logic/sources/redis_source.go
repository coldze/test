package sources

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis"

	"github.com/coldze/test/logic"
)

type ResponseFactory func([]byte) (logic.Response, error)
type DataBuilderFactory func() logic.DataBuilder
type DataParser func([]byte) (logic.Contact, error)

type redisCacheSource struct {
	createResponse ResponseFactory
	createBuilder  DataBuilderFactory
	parse          DataParser
	cache          RedisWrap
	ttl            time.Duration
}

func (r *redisCacheSource) Get(key string) (logic.Response, error) {
	rawData, err := r.cache.Get(key)
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	data, ok := rawData.(string)
	if !ok {
		return nil, fmt.Errorf("cached data is not of type string, it's type is: %T", rawData)
	}
	return r.createResponse([]byte(data))
}

func (r *redisCacheSource) decode(response logic.Response) ([]byte, *logic.Contact, error) {
	b := r.createBuilder()
	if b == nil {
		return nil, nil, errors.New("internal error - builder is nil")
	}
	err := response.Write(b)
	if err != nil {
		return nil, nil, err
	}
	data, err := b.Build()
	if err != nil {
		return nil, nil, err
	}
	contact, err := r.parse(data)
	if err != nil {
		return nil, nil, err
	}
	return data, &contact, nil
}

func (r *redisCacheSource) Remove(response logic.Response) error {
	_, contact, err := r.decode(response)
	if err != nil {
		return err
	}
	return r.cache.Del(contact.ID)
}

func (r *redisCacheSource) Insert(response logic.Response) error {
	data, contact, err := r.decode(response)
	if err != nil {
		return err
	}
	return r.cache.Set(contact.ID, data, r.ttl)
}

func NewRedisCacheSource(cache RedisWrap, ttl time.Duration) CacheSource {
	return &redisCacheSource{
		cache:          cache,
		createResponse: logic.NewJsonOkResponse,
		createBuilder:  NewHttpDataBuilder,
		parse:          logic.ParseContact,
		ttl:            ttl,
	}
}
