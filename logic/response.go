package logic

import (
	"net/http"

	"github.com/coldze/test/consts"
)

type Response interface {
	Write(w http.ResponseWriter) error
}

type HttpResponseFactory func(response *http.Response) (Response, error)

type httpResponse struct {
	data    []byte
	headers http.Header
	code    int
}

func (r *httpResponse) Write(w http.ResponseWriter) error {
	for k, v := range r.headers {
		for _, h := range v {
			w.Header().Add(k, h)
		}
	}
	w.WriteHeader(r.code)
	_, err := w.Write(r.data)
	return err
}

func NewHttpResponse(data []byte, headers http.Header, code int) (Response, error) {
	return &httpResponse{
		data:    data,
		headers: headers,
		code:    code,
	}, nil
}

func NewJsonOkResponse(data []byte) (Response, error) {
	headers := http.Header{}
	headers.Set(consts.HEADER_CONTENT_TYPE, consts.MIME_APPLICATION_JSON)
	return NewHttpResponse(data, headers, http.StatusOK)
}

func NewHttpResponseFactory(getData ResponseDataExtractor) HttpResponseFactory {
	return func(response *http.Response) (Response, error) {
		data, err := getData(response)
		if err != nil {
			return nil, err
		}
		return NewHttpResponse(data, response.Header, response.StatusCode)
	}
}

func NewDefaultHttpResponseFactory() HttpResponseFactory {
	return NewHttpResponseFactory(GetResponseBodyData)
}
