package middleware

import "net/http"

type DataBuilder interface {
	http.ResponseWriter
	Build() ([]byte, error)
}
