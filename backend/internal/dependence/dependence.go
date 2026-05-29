package dependence

import (
	"fmt"
	"log/slog"
	"net"
	"net/url"
	"strconv"
	"strings"

	"github.com/Seeker32/AssassinIoT/backend/ent"
	"github.com/Seeker32/AssassinIoT/backend/internal/conf"
	"github.com/Seeker32/AssassinIoT/backend/internal/logging"
	_ "github.com/lib/pq"

	"entgo.io/ent/dialect"
)

// Dep defines the interface for dependency injection in the application.
type Dep interface {
	ConfigProvider() conf.ConfigProvider
	DBClient() *ent.Client
	Logger() *slog.Logger
}

type dependence struct {
	configProvider conf.ConfigProvider
	dbClient       *ent.Client
	logger         *slog.Logger

	configPath string
}

func NewDependence(opts ...Option) Dep {
	d := &dependence{}
	for _, opt := range opts {
		opt.apply(d)
	}

	return d
}

func (d *dependence) ConfigProvider() conf.ConfigProvider {
	if d.configProvider != nil {
		return d.configProvider
	}

	if d.configPath == "" {
		d.configPath = "config.yaml"
	}

	provider, err := conf.NewConfigProvider(d.configPath)
	if err != nil {
		panic(err)
	}
	d.configProvider = provider

	return d.configProvider
}

func (d *dependence) Logger() *slog.Logger {
	if d.logger != nil {
		return d.logger
	}
	cfg := d.ConfigProvider().ServerConfig()
	var level slog.Level
	if err := level.UnmarshalText([]byte(cfg.LogLevel)); err != nil {
		level = slog.LevelInfo
	}

	d.logger = logging.NewLogger(level, logging.HandlerType(cfg.LogType), nil)

	return d.logger
}

func (d *dependence) DBClient() *ent.Client {
	if d.dbClient != nil {
		return d.dbClient
	}

	driverName, dataSourceName, err := resolveDatabaseConnection(d.ConfigProvider().DatabaseConfig())
	if err != nil {
		panic(err)
	}

	client, err := ent.Open(driverName, dataSourceName)
	if err != nil {
		panic(err)
	}
	d.dbClient = client

	return d.dbClient
}

func resolveDatabaseConnection(cfg conf.DBConfig) (string, string, error) {
	if cfg.DatabaseURL != "" {
		return resolveDatabaseURL(cfg.DatabaseURL)
	}

	dataSourceName, err := buildPostgresDSN(cfg)
	if err != nil {
		return "", "", err
	}

	return dialect.Postgres, dataSourceName, nil
}

func resolveDatabaseURL(databaseURL string) (string, string, error) {
	parsed, err := url.Parse(databaseURL)
	if err != nil {
		return "", "", err
	}

	switch strings.ToLower(parsed.Scheme) {
	case "postgres", "postgresql":
		return dialect.Postgres, databaseURL, nil
	case "sqlite":
		return dialect.SQLite, strings.TrimPrefix(databaseURL, "sqlite://"), nil
	}

	if strings.HasPrefix(databaseURL, "file:") || strings.HasPrefix(databaseURL, ":memory:") {
		return dialect.SQLite, databaseURL, nil
	}

	return "", "", fmt.Errorf("unsupported database URL scheme: %q", parsed.Scheme)
}

func buildPostgresDSN(cfg conf.DBConfig) (string, error) {
	if cfg.Host == "" {
		return "", fmt.Errorf("missing database host")
	}
	if cfg.Database == "" {
		return "", fmt.Errorf("missing database name")
	}

	port := cfg.Port
	if port == 0 {
		port = 5432
	}

	dsn := &url.URL{
		Scheme: "postgres",
		Host:   net.JoinHostPort(cfg.Host, strconv.Itoa(port)),
		Path:   cfg.Database,
	}
	switch {
	case cfg.Username != "" && cfg.Password != "":
		dsn.User = url.UserPassword(cfg.Username, cfg.Password)
	case cfg.Username != "":
		dsn.User = url.User(cfg.Username)
	}
	if cfg.SSLMode != "" {
		query := dsn.Query()
		query.Set("sslmode", cfg.SSLMode)
		dsn.RawQuery = query.Encode()
	}

	return dsn.String(), nil
}
