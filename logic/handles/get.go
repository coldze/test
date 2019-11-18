package handles

import (
	"net/http"

	"github.com/coldze/test/logic"
	"github.com/coldze/test/logic/sources"
)

func NewGetHandler(loggerFactory LoggerFactory, src sources.DataSource, getData logic.RequestDataExtractor) http.HandlerFunc {
	lHandler := logicHandler(src.Get)
	handler := newHttpHandler(getData, lHandler)
	return newCheckAndSetLoggerMiddleware(loggerFactory, handler)
}
