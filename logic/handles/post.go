package handles

import (
	"net/http"

	"github.com/coldze/test/logic"
	"github.com/coldze/test/logic/sources"
)

func NewPostHandler(loggerFactory LoggerFactory, src sources.DataSource) http.HandlerFunc {
	lHandler := logicHandler(src.Create)
	getData := logic.GetRequestBodyData
	handler := newHttpHandler(getData, lHandler)
	return newCheckAndSetLoggerMiddleware(loggerFactory, handler)
}
