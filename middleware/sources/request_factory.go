package sources

import (
	"bytes"
	"context"
	"net/http"

	"github.com/coldze/test/consts"
	"github.com/coldze/test/utils"
)

type RequestFactory func(ctx context.Context, data []byte, url string, method string) (*http.Request, error)

func defaultRequestFactory(ctx context.Context, data []byte, url string, method string) (*http.Request, error) {
	r := bytes.NewReader(data)
	req, err := http.NewRequest(method, url, r)
	if err != nil {
		return nil, err
	}
	headers := utils.GetHeaders(ctx)
	if headers != nil {
		req.Header = headers
	}
	req.Header.Set(consts.HEADER_CONTENT_TYPE, consts.MIME_APPLICATION_JSON)
	return req, nil
}
