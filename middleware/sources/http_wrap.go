package sources

import "net/http"

type HttpDo func(r *http.Request) (*http.Response, error)

func NewHttpWrap(client *http.Client) HttpDo {
	return func(r *http.Request) (*http.Response, error) {
		return client.Do(r)
	}
}
