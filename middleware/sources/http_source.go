package sources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/coldze/test/middleware"
	"github.com/coldze/test/utils"
)

type httpDataSource struct {
	do             HttpDo
	url            string
	createRequest  RequestFactory
	createResponse middleware.HttpResponseFactory
}

func (h *httpDataSource) call(ctx context.Context, data []byte, url string, method string) (middleware.Response, error) {
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

func (h *httpDataSource) Get(ctx context.Context, key []byte) (middleware.Response, error) {
	return h.call(ctx, nil, h.url+"/"+string(key), http.MethodGet)
}

func (h *httpDataSource) Create(ctx context.Context, data []byte) (middleware.Response, error) {
	return h.call(ctx, data, h.url, http.MethodPost)
}

func (h *httpDataSource) Update(ctx context.Context, data []byte) (middleware.Response, error) {
	return h.call(ctx, data, h.url, http.MethodPut)
}

func NewHttpDataSource(client HttpDo, url string) DataSource {
	return &httpDataSource{
		url:            url,
		do:             client,
		createRequest:  defaultRequestFactory,
		createResponse: middleware.NewDefaultHttpResponseFactory(),
	}
}

func NewDefaultHttpDataSource(url string) DataSource {
	return NewHttpDataSource(NewHttpWrap(http.DefaultClient), url)
}
