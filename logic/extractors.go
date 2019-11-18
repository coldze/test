package logic

import (
	"io"
	"io/ioutil"
	"net/http"
)

type RequestDataExtractor func(r *http.Request) ([]byte, error)
type ResponseDataExtractor func(r *http.Response) ([]byte, error)

func getBodyData(r io.ReadCloser) ([]byte, error) {
	if r == nil {
		return []byte{}, nil
	}
	return ioutil.ReadAll(r)
}

func GetRequestBodyData(r *http.Request) ([]byte, error) {
	if r == nil {
		return []byte{}, nil
	}
	return getBodyData(r.Body)
}

func GetResponseBodyData(r *http.Response) ([]byte, error) {
	if r == nil {
		return []byte{}, nil
	}
	return getBodyData(r.Body)
}
