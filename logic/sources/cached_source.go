package sources

import (
	"context"

	"github.com/coldze/test/logic"
	"github.com/coldze/test/utils"
)

type cachedDataSource struct {
	original DataSource
	cache    CacheSource
}

func (c *cachedDataSource) Get(ctx context.Context, key []byte) (logic.Response, error) {
	logger := utils.GetLogger(ctx)
	res, err := c.cache.Get(string(key))
	if err != nil {
		logger.Warningf("Error occurred while getting data from cache. Error: %v", err)
	} else if res != nil {
		return res, nil
	}
	res, err = c.original.Get(ctx, key)
	if err != nil {
		return res, err
	}
	err = c.cache.Insert(res)
	if err != nil {
		logger.Warningf("Error occurred while inserting data to cache. Error: %v", err)
	}
	return res, nil
}

func (c *cachedDataSource) Create(ctx context.Context, data []byte) (logic.Response, error) {
	res, err := c.original.Create(ctx, data)
	if err != nil {
		return res, err
	}
	err = c.cache.Remove(res)
	if err != nil {
		logger := utils.GetLogger(ctx)
		logger.Warningf("Failed to remove value from cache. Error: %v", err)
	}
	return res, nil
}

func (c *cachedDataSource) Update(ctx context.Context, data []byte) (logic.Response, error) {
	return c.Create(ctx, data)
}

func NewCachedDataSource(original DataSource, cache CacheSource) DataSource {
	return &cachedDataSource{
		original: original,
		cache:    cache,
	}
}
