package application

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Seeker32/AssassinIoT/backend/internal/dependence"
)

const (
	defaultServerAddr = ":8080"
	shutdownTimeout   = 5 * time.Second
)

type server struct {
	dep            dependence.Dep
	logger         *slog.Logger
	listenAndServe func(*http.Server) error
	shutdown       func(*http.Server) error
	signals        func() (<-chan os.Signal, func())
}

func NewServer(dep dependence.Dep) *server {
	return &server{
		dep:    dep,
		logger: dep.Logger(),
	}
}

func (s *server) Start() error {
	addr := s.resolveAddr()
	httpServer := &http.Server{
		Addr:    addr,
		Handler: http.NewServeMux(),
	}

	errCh := make(chan error, 1)
	go func() {
		s.logger.Info("starting server", "addr", addr)
		errCh <- s.runHTTPServer(httpServer)
	}()

	sigCh, stopSignals := s.signalSource()
	defer stopSignals()

	select {
	case err := <-errCh:
		return s.closeDependencies(normalizeServerError(err))
	case sig := <-sigCh:
		s.logger.Info("shutting down server", "signal", sig.String())
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		err := s.shutdownHTTPServer(ctx, httpServer)
		if err == nil {
			if serverErr := <-errCh; serverErr != nil && serverErr != http.ErrServerClosed {
				err = serverErr
			}
		}
		return s.closeDependencies(err)
	}
}

func (s *server) resolveAddr() string {
	if cfgAddr := s.dep.ConfigProvider().ServerConfig().Addr; cfgAddr != "" {
		return cfgAddr
	}

	return defaultServerAddr
}

func (s *server) runHTTPServer(httpServer *http.Server) error {
	if s.listenAndServe != nil {
		return s.listenAndServe(httpServer)
	}
	return httpServer.ListenAndServe()
}

func (s *server) signalSource() (<-chan os.Signal, func()) {
	if s.signals != nil {
		return s.signals()
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	return sigCh, func() {
		signal.Stop(sigCh)
	}
}

func (s *server) shutdownHTTPServer(ctx context.Context, httpServer *http.Server) error {
	if s.shutdown != nil {
		return s.shutdown(httpServer)
	}

	return httpServer.Shutdown(ctx)
}

func (s *server) closeDependencies(err error) error {
	return errors.Join(err, s.dep.Close())
}

func normalizeServerError(err error) error {
	if err == nil || err == http.ErrServerClosed {
		return nil
	}

	return err
}
