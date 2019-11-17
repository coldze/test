package utils

import (
	"context"
	"net/http"

	"github.com/coldze/test/logs"
)

type loggerKey struct{}

type headerKey struct{}

var (
	//it is recommended to use structs as keys for values in context - not to overlap with other packages by accident.
	loggerCtxKey loggerKey
	headerCtxKey headerKey

	//global variables are bad, but this one is not that bad - it's not exported outside and is used as a default logger, in case nothing was set in context - to remove checking == nil every single time.
	defaultLogger logs.Logger
)

func SetLogger(ctx context.Context, logger logs.Logger) context.Context {
	return context.WithValue(ctx, loggerCtxKey, logger)
}

func GetLogger(ctx context.Context) logs.Logger {
	res := ctx.Value(loggerCtxKey)
	if res == nil {
		return defaultLogger
	}
	logger, ok := res.(logs.Logger)
	if !ok {
		return defaultLogger
	}
	return logger
}

func SetHeaders(ctx context.Context, headers http.Header) context.Context {
	return context.WithValue(ctx, headerCtxKey, headers)
}

func GetHeaders(ctx context.Context) http.Header {
	res := ctx.Value(headerCtxKey)
	if res == nil {
		return nil
	}
	headers, ok := res.(http.Header)
	if !ok {
		return nil
	}
	return headers
}

func init() {
	defaultLogger = logs.NewStdLogger()
}
