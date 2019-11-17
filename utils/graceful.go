package utils

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/coldze/test/logs"
)

const (
	shutdown_channel_size = 10
)

type MainFunc func(logger logs.Logger, stopping <-chan struct{}) int

type gracefulShutdown struct {
	WaitForShutdown  <-chan struct{}
	ShutdownComplete chan<- int
	ReturnCode       <-chan int
}

func runGracefully(timeout time.Duration, logger logs.Logger) *gracefulShutdown {
	shutdown := make(chan struct{}, shutdown_channel_size)
	shutdownComplete := make(chan int, shutdown_channel_size)
	returnCode := make(chan int, shutdown_channel_size)

	gracefulStop := make(chan os.Signal, shutdown_channel_size)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	go func() {
		defer close(shutdown)
		defer close(gracefulStop)
		exitCode := 0
		select {
		case exitCode = <-shutdownComplete:
			{
				logger.Infof("Business logic completed. Exit code: %+v", exitCode)
				returnCode <- exitCode
				return
			}

		case sig := <-gracefulStop:
			logger.Infof("Caught sig: %+v", sig)
			shutdown <- struct{}{}
		}

		select {
		case <-time.After(timeout):
			logger.Warningf("Shutdown timeout occured. Terminating.")
			returnCode <- 1
		case exitCode := <-shutdownComplete:
			logger.Infof("Shutdown complete.")
			returnCode <- exitCode
		}
	}()

	return &gracefulShutdown{
		WaitForShutdown:  shutdown,
		ShutdownComplete: shutdownComplete,
		ReturnCode:       returnCode,
	}
}

func safeRunAppLogic(appLogic MainFunc, stopChan <-chan struct{}, logger logs.Logger) (res int) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}
		err, ok := r.(error)
		if ok {
			logger.Errorf("mainFunc failed: %v", err)
			res = 1
			return
		}
		logger.Errorf("mainFunc failed. Unknown error: %+v. Type: %T", r, r)
		res = 1
	}()
	return appLogic(logger, stopChan)
}

func Run(timeout time.Duration, appLogic MainFunc, logger logs.Logger) {

	graceful := runGracefully(timeout, logger)
	go func() {
		graceful.ShutdownComplete <- safeRunAppLogic(appLogic, graceful.WaitForShutdown, logger)
	}()
	exitCode := <-graceful.ReturnCode
	logger.Infof("Exiting application. Code: %+v", exitCode)
	os.Exit(exitCode)
}
