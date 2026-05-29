package dependence

import (
	"log/slog"

	"github.com/Seeker32/AssassinIoT/backend/ent"
	"github.com/Seeker32/AssassinIoT/backend/internal/conf"
)

type Option interface {
	apply(*dependence)
}

type optionFunc func(*dependence)

func (f optionFunc) apply(d *dependence) {
	f(d)
}

func WithConfigPath(path string) Option {
	return optionFunc(func(d *dependence) {
		d.configPath = path
	})
}

func WithConfigProvider(c conf.ConfigProvider) Option {
	return optionFunc(func(d *dependence) {
		d.configProvider = c
	})
}

func WithLogger(logger *slog.Logger) Option {
	return optionFunc(func(d *dependence) {
		d.logger = logger
	})
}

func WithDBClient(client *ent.Client) Option {
	return optionFunc(func(d *dependence) {
		d.dbClient = client
	})
}
