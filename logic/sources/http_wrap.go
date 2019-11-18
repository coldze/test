package sources

import "net/http"

type HttpDo func(r *http.Request) (*http.Response, error)
