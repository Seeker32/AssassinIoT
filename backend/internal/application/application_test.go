package application

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"os"
	"reflect"
	"syscall"
	"testing"

	"github.com/Seeker32/AssassinIoT/backend/ent"
	"github.com/Seeker32/AssassinIoT/backend/internal/conf"
	"github.com/Seeker32/AssassinIoT/backend/internal/dependence"
)

var errListenStub = errors.New("listen stub")
var errCloseStub = errors.New("close stub")
var errShutdownStub = errors.New("shutdown stub")

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

type testDep struct {
	configProvider conf.ConfigProvider
	logger         *slog.Logger
	closeErr       error
	events         *[]string
}

func (d *testDep) ConfigProvider() conf.ConfigProvider {
	return d.configProvider
}

func (d *testDep) DBClient() *ent.Client {
	return nil
}

func (d *testDep) Logger() *slog.Logger {
	return d.logger
}

func (d *testDep) Close() error {
	if d.events != nil {
		*d.events = append(*d.events, "close")
	}
	return d.closeErr
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

func TestServerStartClosesDependenciesWhenListenReturns(t *testing.T) {
	t.Helper()

	events := []string{}
	srv := NewServer(&testDep{
		configProvider: testConfigProvider{},
		logger:         testLogger(),
		events:         &events,
	})

	srv.listenAndServe = func(httpServer *http.Server) error {
		return errListenStub
	}

	err := srv.Start()
	if !errors.Is(err, errListenStub) {
		t.Fatalf("Start() error = %v, want error containing %v", err, errListenStub)
	}
	if !reflect.DeepEqual(events, []string{"close"}) {
		t.Fatalf("events = %v, want [close]", events)
	}
}

func TestServerStartReturnsJoinedListenAndCloseErrors(t *testing.T) {
	t.Helper()

	srv := NewServer(&testDep{
		configProvider: testConfigProvider{},
		logger:         testLogger(),
		closeErr:       errCloseStub,
	})

	srv.listenAndServe = func(httpServer *http.Server) error {
		return errListenStub
	}

	err := srv.Start()
	if !errors.Is(err, errListenStub) {
		t.Fatalf("Start() error = %v, want error containing %v", err, errListenStub)
	}
	if !errors.Is(err, errCloseStub) {
		t.Fatalf("Start() error = %v, want error containing %v", err, errCloseStub)
	}
}

func TestServerStartClosesDependenciesAfterSignalShutdown(t *testing.T) {
	t.Helper()

	events := []string{}
	sigCh := make(chan os.Signal, 1)
	shutdownDone := make(chan struct{})
	sigCh <- syscall.SIGTERM

	srv := NewServer(&testDep{
		configProvider: testConfigProvider{},
		logger:         testLogger(),
		events:         &events,
	})

	srv.listenAndServe = func(httpServer *http.Server) error {
		<-shutdownDone
		return http.ErrServerClosed
	}
	srv.shutdown = func(httpServer *http.Server) error {
		events = append(events, "shutdown")
		close(shutdownDone)
		return nil
	}
	srv.signals = func() (<-chan os.Signal, func()) {
		return sigCh, func() {}
	}

	if err := srv.Start(); err != nil {
		t.Fatalf("Start() error = %v, want nil", err)
	}
	if !reflect.DeepEqual(events, []string{"shutdown", "close"}) {
		t.Fatalf("events = %v, want [shutdown close]", events)
	}
}

func TestServerStartReturnsJoinedShutdownAndCloseErrors(t *testing.T) {
	t.Helper()

	sigCh := make(chan os.Signal, 1)
	sigCh <- syscall.SIGTERM

	srv := NewServer(&testDep{
		configProvider: testConfigProvider{},
		logger:         testLogger(),
		closeErr:       errCloseStub,
	})

	srv.listenAndServe = func(httpServer *http.Server) error {
		return http.ErrServerClosed
	}
	srv.shutdown = func(httpServer *http.Server) error {
		return errShutdownStub
	}
	srv.signals = func() (<-chan os.Signal, func()) {
		return sigCh, func() {}
	}

	err := srv.Start()
	if !errors.Is(err, errShutdownStub) {
		t.Fatalf("Start() error = %v, want error containing %v", err, errShutdownStub)
	}
	if !errors.Is(err, errCloseStub) {
		t.Fatalf("Start() error = %v, want error containing %v", err, errCloseStub)
	}
}
