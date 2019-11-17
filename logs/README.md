## Logs

[Go to main](../README.md)

This package contains logger interface and two implementations of logger:
* std_logger.go - uses `log` package for output.
* prefixed_logger - wraps provided logger with prefix. For example, we can create a logger with prefix that contains ID of http request and this will give us an ability to track logs, related to one particular http-request.

StdLogger could be split up in to two separate implementations to fit into SOLID principles:
1. LeveledLogger, that will wrap another logger and ad prefix `[ERROR]` or `[WARNING]`, depending on method being called
2. StdLogger that is just a proxy to "log" package.

I decided not to do that here just to keep it a bit simpler + we will get rid of additional call.
