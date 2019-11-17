package utils

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
)

type serviceImpl struct {
	srv     *http.Server
	lock    sync.RWMutex
	stopped bool
}

type Service interface {
	Stop() error
}

func (s *serviceImpl) Stop() error {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.stopped {
		return errors.New("already stopped")
	}
	s.stopped = true
	err := s.srv.Shutdown(context.Background())
	if err == nil {
		return nil
	}
	return fmt.Errorf("failed to stop http-service: %v", err)
}

func (s *serviceImpl) isStopped() bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.stopped
}

func NewService(bindAddress string, handler http.Handler) (Service, error) {
	s := &serviceImpl{
		srv: &http.Server{
			Addr:    bindAddress,
			Handler: handler,
		},
	}

	go func() {
		err := s.srv.ListenAndServe()
		if err == nil {
			return
		}
		if s.isStopped() {
			return
		}
		panic(fmt.Errorf("http-service failed to listen: %v", err))
	}()

	return s, nil
}
