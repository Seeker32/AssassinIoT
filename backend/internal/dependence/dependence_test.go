package dependence

import (
	"context"
	"testing"

	"github.com/Seeker32/AssassinIoT/backend/ent"
	"github.com/Seeker32/AssassinIoT/backend/internal/conf"
	"entgo.io/ent/dialect"
)

type testConfigProvider struct {
	dbConfig conf.DBConfig
}

func (t testConfigProvider) DatabaseConfig() conf.DBConfig {
	return t.dbConfig
}

func (t testConfigProvider) ServerConfig() conf.ServerConfig {
	return conf.ServerConfig{}
}

type testDriver struct {
	closeErr error
	closed   bool
}

func (d *testDriver) Exec(context.Context, string, any, any) error {
	return nil
}

func (d *testDriver) Query(context.Context, string, any, any) error {
	return nil
}

func (d *testDriver) Tx(context.Context) (dialect.Tx, error) {
	return dialect.NopTx(d), nil
}

func (d *testDriver) Close() error {
	d.closed = true
	return d.closeErr
}

func (d *testDriver) Dialect() string {
	return dialect.SQLite
}

func TestDBClientReturnsInjectedClient(t *testing.T) {
	t.Helper()

	want := &ent.Client{}
	dep := NewDependence(WithDBClient(want))

	got := dep.DBClient()
	if got != want {
		t.Fatalf("DBClient() = %p, want %p", got, want)
	}
}

func TestDBClientUsesDatabaseURLAndCachesResult(t *testing.T) {
	t.Helper()

	dep := NewDependence(WithConfigProvider(testConfigProvider{
		dbConfig: conf.DBConfig{
			Host:        "invalid-host",
			Port:        1,
			Database:    "ignored",
			Username:    "ignored",
			Password:    "ignored",
			SSLMode:     "disable",
			DatabaseURL: "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable",
		},
	}))

	first := dep.DBClient()
	t.Cleanup(func() {
		_ = first.Close()
	})

	second := dep.DBClient()
	if first == nil {
		t.Fatal("DBClient() returned nil client")
	}
	if second != first {
		t.Fatalf("DBClient() did not cache client: first=%p second=%p", first, second)
	}
}

func TestDBClientBuildsClientFromDatabaseConfigFields(t *testing.T) {
	t.Helper()

	dep := NewDependence(WithConfigProvider(testConfigProvider{
		dbConfig: conf.DBConfig{
			Host:     "localhost",
			Port:     5432,
			Database: "postgres",
			Username: "postgres",
			Password: "postgres",
			SSLMode:  "disable",
		},
	}))

	client := dep.DBClient()
	t.Cleanup(func() {
		_ = client.Close()
	})

	if client == nil {
		t.Fatal("DBClient() returned nil client")
	}
}

func TestDBClientPanicsOnInvalidDatabaseURL(t *testing.T) {
	t.Helper()

	dep := NewDependence(WithConfigProvider(testConfigProvider{
		dbConfig: conf.DBConfig{
			DatabaseURL: "postgres://%zz",
		},
	}))

	defer func() {
		if recover() == nil {
			t.Fatal("DBClient() did not panic")
		}
	}()

	_ = dep.DBClient()
}

func TestCloseReturnsNilWhenDBClientWasNeverInitialized(t *testing.T) {
	t.Helper()

	dep := NewDependence()

	if err := dep.Close(); err != nil {
		t.Fatalf("Close() error = %v, want nil", err)
	}
}

func TestCloseClosesCachedClientAndClearsIt(t *testing.T) {
	t.Helper()

	drv := &testDriver{}
	client := ent.NewClient(ent.Driver(drv))
	dep := NewDependence(WithDBClient(client))

	if err := dep.Close(); err != nil {
		t.Fatalf("Close() error = %v, want nil", err)
	}
	if !drv.closed {
		t.Fatal("Close() did not close the driver")
	}

	if got := dep.(*dependence).dbClient; got != nil {
		t.Fatalf("dbClient after Close() = %v, want nil", got)
	}
}

func TestCloseIsSafeToCallMoreThanOnce(t *testing.T) {
	t.Helper()

	client := ent.NewClient(ent.Driver(&testDriver{}))
	dep := NewDependence(WithDBClient(client))

	if err := dep.Close(); err != nil {
		t.Fatalf("first Close() error = %v, want nil", err)
	}
	if err := dep.Close(); err != nil {
		t.Fatalf("second Close() error = %v, want nil", err)
	}
}
