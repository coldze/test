package handles

import (
	"net/http"

	"github.com/coldze/test/logic"
	"github.com/coldze/test/logic/sources"
)

func NewPutHandler(loggerFactory LoggerFactory, src sources.DataSource) http.HandlerFunc {
	lHandler := logicHandler(src.Update)
	getData := logic.GetRequestBodyData
	handler := newHttpHandler(getData, lHandler)
	return newCheckAndSetLoggerMiddleware(loggerFactory, handler)
}
