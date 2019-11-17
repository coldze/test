package handles

import (
	"context"
	"net/http"

	"github.com/coldze/test/middleware"
	"github.com/coldze/test/utils"
)

type logicHandler func(ctx context.Context, data []byte) (middleware.Response, error)
type wrappedHttpHandler func(w http.ResponseWriter, r middleware.HttpRequest)

func newWrappedHttpHandler(getData middleware.RequestDataExtractor, handler logicHandler) wrappedHttpHandler {
	return func(w http.ResponseWriter, r middleware.HttpRequest) {
		ctx := r.GetContext()
		logger := utils.GetLogger(ctx)
		body := r.GetBody()
		if body != nil {
			defer func() {
				err := body.Close()
				if err != nil {
					logger.Errorf("Failed to close body. Error: %v", err)
				}
			}()
		}
		defer func() {
			r := recover()
			if r == nil {
				return
			}
			err, ok := r.(error)
			if ok {
				logger.Errorf("Panic occurred in handler. Error: %v", err)
				return
			}
			logger.Errorf("Panic occurred in handler. Unknown error: %+v. Type: %T.", err, err)
		}()
		data, err := getData(r)
		if err != nil {
			logger.Errorf("Failed to read body. Error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		ctx = utils.SetHeaders(ctx, r.GetHeader())
		res, err := handler(ctx, data)
		if err != nil {
			logger.Errorf("Failed to process. Error: %v", err)
		}
		if res == nil {
			logger.Errorf("Response is empty.")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = res.Write(w)
		if err != nil {
			logger.Errorf("Failed to write response. Error: %v", err)
		}
	}
}

//Description: this is a handler, where http handling code is gathered. It strips that off and passes data down to business logic.
func newHttpHandler(getData middleware.RequestDataExtractor, handler logicHandler) http.HandlerFunc {
	wrap := newWrappedHttpHandler(getData, handler)
	return func(w http.ResponseWriter, r *http.Request) {
		wrap(w, middleware.NewHttpRequestWrap(r))
	}
}
