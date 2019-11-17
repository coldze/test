package sources

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis"

	"github.com/coldze/test/middleware"
)

func parseContact(data []byte) (middleware.Contact, error) {
	contact := middleware.Contact{}
	err := json.Unmarshal(data, &contact)
	return contact, err
}

type ResponseFactory func([]byte) (middleware.Response, error)
type DataBuilderFactory func() middleware.DataBuilder
type DataParser func([]byte) (middleware.Contact, error)

type redisCacheSource struct {
	createResponse ResponseFactory
	createBuilder  DataBuilderFactory
	parse          DataParser
	cache          RedisWrap
	ttl            time.Duration
}

func (r *redisCacheSource) Get(key string) (middleware.Response, error) {
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

func (r *redisCacheSource) decode(response middleware.Response) ([]byte, *middleware.Contact, error) {
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

func (r *redisCacheSource) Remove(response middleware.Response) error {
	_, contact, err := r.decode(response)
	if err != nil {
		return err
	}
	return r.cache.Del(contact.ID)
}

func (r *redisCacheSource) Insert(response middleware.Response) error {
	data, contact, err := r.decode(response)
	if err != nil {
		return err
	}
	return r.cache.Set(contact.ID, data, r.ttl)
}

func NewRedisCacheSource(cache RedisWrap, ttl time.Duration) CacheSource {
	return &redisCacheSource{
		cache:          cache,
		createResponse: middleware.NewJsonOkResponse,
		createBuilder:  NewHttpDataBuilder,
		parse:          parseContact,
		ttl:            ttl,
	}
}
