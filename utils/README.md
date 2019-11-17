## Utils

[Go to main](../README.md)

This package contains utility functions.

### Graceful shutdown (`graceful.go`).
As we're running an http-service, it would be nice to have a graceful shutdown - our main function can be notified, when application is going to be interrupted.

Function `Run` accepts arguments
- `timeout` - how long should it wait for main function to complete. When timeout occures, panic is thrown, application will be terminated.
- `appLogic` - that's your main function, which is provided with channel `stopping`, that signals when application has to stop.
- `logger` - logger

### Service (`service.go`)
Creates an implementation of http-service, that can be stopped (method `Stop`).

### Context (`context.go`)
Helper functions to set values to context and retrieve values from context.
* can set/get a logger to/from context
* can set/get http.Header to/from context