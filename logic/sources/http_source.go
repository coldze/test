package sources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/coldze/test/logic"
	"github.com/coldze/test/utils"
)

type httpDataSource struct {
	do             HttpDo
	url            string
	createRequest  RequestFactory
	createResponse logic.HttpResponseFactory
}

func (h *httpDataSource) call(ctx context.Context, data []byte, url string, method string) (logic.Response, error) {
	req, err := h.createRequest(ctx, data, url, method)
	if err != nil {
		return nil, err
	}
	resp, err := h.do(req)
	if resp != nil && resp.Body != nil {
		defer func() {
			err := resp.Body.Close()
			if err != nil {
				logger := utils.GetLogger(ctx)
				logger.Warningf("Failed to close body. Error: %v", err)
			}
		}()
	}
	if err != nil {
		return nil, err
	}
	wrappedResp, err := h.createResponse(resp)
	if err == nil && resp.StatusCode != 200 {
		err = fmt.Errorf("response status code is not 200, code - %v, status - '%v'", resp.StatusCode, resp.Status)
	}
	return wrappedResp, err
}

func (h *httpDataSource) Get(ctx context.Context, key []byte) (logic.Response, error) {
	return h.call(ctx, nil, h.url+"/"+string(key), http.MethodGet)
}

func (h *httpDataSource) Create(ctx context.Context, data []byte) (logic.Response, error) {
	return h.call(ctx, data, h.url, http.MethodPost)
}

func (h *httpDataSource) Update(ctx context.Context, data []byte) (logic.Response, error) {
	return h.call(ctx, data, h.url, http.MethodPost)
}

func NewHttpDataSource(do HttpDo, url string) DataSource {
	return &httpDataSource{
		url:            url,
		do:             do,
		createRequest:  DefaultRequestFactory,
		createResponse: logic.NewDefaultHttpResponseFactory(),
	}
}

func NewDefaultHttpDataSource(url string) DataSource {
	return NewHttpDataSource(http.DefaultClient.Do, url)
}
