package handles

import (
	"context"
	"net/http"

	"github.com/coldze/test/logic"
	"github.com/coldze/test/utils"
)

type logicHandler func(ctx context.Context, data []byte) (logic.Response, error)

func newHttpHandler(getData logic.RequestDataExtractor, handler logicHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := utils.GetLogger(ctx)
		body := r.Body
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
			logger.Errorf("Panic occurred in handler. Unknown error: %+v. Type: %T.", r, r)
		}()
		data, err := getData(r)
		if err != nil {
			logger.Errorf("Failed to read body. Error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		ctx = utils.SetHeaders(ctx, r.Header)
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
