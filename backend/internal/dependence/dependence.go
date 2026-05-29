package dependence

import (
	"log/slog"

	"github.com/Seeker32/AssassinIoT/backend/internal/conf"
	"github.com/Seeker32/AssassinIoT/backend/internal/logging"
)

type Dep interface {
	ConfigProvider() conf.ConfigProvider
	Logger() *slog.Logger
}

type dependence struct {
	configProvider conf.ConfigProvider
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
