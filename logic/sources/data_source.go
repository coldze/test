package sources

import (
	"context"

	"github.com/coldze/test/logic"
)

type DataSource interface {
	Get(ctx context.Context, data []byte) (logic.Response, error)
	Create(ctx context.Context, data []byte) (logic.Response, error)
	Update(ctx context.Context, data []byte) (logic.Response, error)
}
