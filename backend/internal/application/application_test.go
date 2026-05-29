package application

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"testing"

	"github.com/Seeker32/AssassinIoT/backend/internal/conf"
	"github.com/Seeker32/AssassinIoT/backend/internal/dependence"
)

var errListenStub = errors.New("listen stub")

type testConfigProvider struct {
	serverConfig conf.ServerConfig
}

func (t testConfigProvider) DatabaseConfig() conf.DBConfig {
	return conf.DBConfig{}
}

func (t testConfigProvider) ServerConfig() conf.ServerConfig {
	return t.serverConfig
}

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestServerStartUsesConfiguredAddr(t *testing.T) {
	t.Helper()

	dep := dependence.NewDependence(
		dependence.WithConfigProvider(testConfigProvider{
			serverConfig: conf.ServerConfig{Addr: "127.0.0.1:8088"},
		}),
		dependence.WithLogger(testLogger()),
	)

	srv := NewServer(dep)

	var gotAddr string
	srv.listenAndServe = func(httpServer *http.Server) error {
		gotAddr = httpServer.Addr
		return errListenStub
	}

	err := srv.Start()
	if !errors.Is(err, errListenStub) {
		t.Fatalf("Start() error = %v, want %v", err, errListenStub)
	}
	if gotAddr != "127.0.0.1:8088" {
		t.Fatalf("Start() addr = %q, want %q", gotAddr, "127.0.0.1:8088")
	}
}

func TestServerStartFallsBackToDefaultAddr(t *testing.T) {
	t.Helper()

	dep := dependence.NewDependence(
		dependence.WithConfigProvider(testConfigProvider{}),
		dependence.WithLogger(testLogger()),
	)

	srv := NewServer(dep)

	var gotAddr string
	srv.listenAndServe = func(httpServer *http.Server) error {
		gotAddr = httpServer.Addr
		return errListenStub
	}

	err := srv.Start()
	if !errors.Is(err, errListenStub) {
		t.Fatalf("Start() error = %v, want %v", err, errListenStub)
	}
	if gotAddr != defaultServerAddr {
		t.Fatalf("Start() addr = %q, want %q", gotAddr, defaultServerAddr)
	}
}
