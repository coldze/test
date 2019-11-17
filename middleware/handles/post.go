package handles

import (
	"net/http"

	"github.com/coldze/test/logs"
	"github.com/coldze/test/middleware"
	"github.com/coldze/test/middleware/sources"
)

func NewPostHandler(logger logs.Logger, src sources.DataSource) http.HandlerFunc {
	logic := logicHandler(src.Create)
	getData := middleware.GetRequestBodyData
	handler := newHttpHandler(getData, logic)
	getLogger := logs.NewPrefixedLogger(logger, "[POST]")
	//This factory can be used to generate logger with ID request as a prefix, but we will just use common logger here to reduce complexity
	dummyLoggerFactory := func() logs.Logger {
		return getLogger
	}
	return newCheckAndSetLoggerMiddleware(dummyLoggerFactory, handler)
}
