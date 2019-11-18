package handles

import "github.com/coldze/test/logs"

//This loggerFactory can generate unique logger for each request, thus they can be traced in logs.
type LoggerFactory func() logs.Logger

func NewDefaultLoggerFactory(defaultLogger logs.Logger) LoggerFactory {
	return func() logs.Logger {
		return defaultLogger
	}
}
