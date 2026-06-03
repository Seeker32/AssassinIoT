package application

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

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
}

func NewServer(dep dependence.Dep) *server {
	return &server{
		dep:    dep,
		logger: dep.Logger(),
	}
}

func (s *server) Start() error {
	addr := s.resolveAddr()

	r := gin.New()
	r.Use(gin.Recovery())

	httpServer := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	errCh := make(chan error, 1)
	go func() {
		s.logger.Info("starting server", "addr", addr)
		errCh <- s.runHTTPServer(httpServer)
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	select {
	case err := <-errCh:
		if err == nil || err == http.ErrServerClosed {
			return nil
		}
		return err
	case sig := <-sigCh:
		s.logger.Info("shutting down server", "signal", sig.String())
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		if err := httpServer.Shutdown(ctx); err != nil {
			return err
		}
		if err := <-errCh; err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
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
