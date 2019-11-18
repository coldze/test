package handles

import (
	"net/http"

	"github.com/coldze/test/utils"
)

func newCheckAndSetLoggerMiddleware(newLogger LoggerFactory, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r == nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		ctx := utils.SetLogger(r.Context(), newLogger())
		next(w, r.WithContext(ctx))
	}
}
