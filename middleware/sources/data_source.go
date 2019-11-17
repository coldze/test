package sources

import (
	"context"

	"github.com/coldze/test/middleware"
)

type DataSource interface {
	Get(ctx context.Context, data []byte) (middleware.Response, error)
	Create(ctx context.Context, data []byte) (middleware.Response, error)
	Update(ctx context.Context, data []byte) (middleware.Response, error)
}
