package handles

import (
	"net/http"

	"github.com/coldze/test/logs"
	"github.com/coldze/test/utils"
)

type loggerFactory func() logs.Logger

func newCheckAndSetLoggerMiddleware(newLogger loggerFactory, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r == nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		if r.Body == nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		ctx := utils.SetLogger(r.Context(), newLogger())
		next(w, r.WithContext(ctx))
	}
}
