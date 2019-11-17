package sources

import (
	"github.com/coldze/test/middleware"
)

type CacheSource interface {
	Get(key string) (middleware.Response, error)
	Insert(response middleware.Response) error
	Remove(response middleware.Response) error
}
