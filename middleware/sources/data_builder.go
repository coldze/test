package sources

import (
	"net/http"

	"github.com/coldze/test/middleware"
)

type httpDataBuilder struct {
	data []byte
}

func (h *httpDataBuilder) Header() http.Header {
	return http.Header{}
}

func (h *httpDataBuilder) Write(data []byte) (int, error) {
	h.data = data
	return len(h.data), nil
}

func (h *httpDataBuilder) WriteHeader(code int) {
}

func (h *httpDataBuilder) Build() ([]byte, error) {
	return h.data, nil
}

func NewHttpDataBuilder() middleware.DataBuilder {
	return &httpDataBuilder{}
}
