package sources

import (
	"github.com/coldze/test/logic"
)

type CacheSource interface {
	Get(key string) (logic.Response, error)
	Insert(response logic.Response) error
	Remove(response logic.Response) error
}
