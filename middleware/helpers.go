package middleware

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
)

type HttpRequest interface {
	GetContext() context.Context
	GetHeader() http.Header
	GetBody() io.ReadCloser
	GetRawRequest() *http.Request
}

type RequestDataExtractor func(r HttpRequest) ([]byte, error)
type ResponseDataExtractor func(r *http.Response) ([]byte, error)

func getBodyData(r io.ReadCloser) ([]byte, error) {
	if r == nil {
		return []byte{}, nil
	}
	return ioutil.ReadAll(r)
}

func GetRequestBodyData(r HttpRequest) ([]byte, error) {
	return getBodyData(r.GetBody())
}

func GetResponseBodyData(r *http.Response) ([]byte, error) {
	return getBodyData(r.Body)
}

type httpRequestWrap struct {
	r *http.Request
}

func (h *httpRequestWrap) GetContext() context.Context {
	return h.r.Context()
}

func (h *httpRequestWrap) GetHeader() http.Header {
	return h.r.Header
}

func (h *httpRequestWrap) GetBody() io.ReadCloser {
	return h.r.Body
}

func (h *httpRequestWrap) GetRawRequest() *http.Request {
	return h.r
}

func NewHttpRequestWrap(r *http.Request) HttpRequest {
	return &httpRequestWrap{
		r: r,
	}
}
