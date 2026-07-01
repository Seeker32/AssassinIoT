// 配置管理模块
package conf

import (
	"os"

	"gopkg.in/yaml.v3"
)

// ConfigProvider defines the interface for providing configuration data.
type ConfigProvider interface {
	DatabaseConfig() DBConfig
	ServerConfig() ServerConfig
}

type initConfigProvider struct {
	databaseConfig DBConfig
	serverConfig   ServerConfig
}

func (i *initConfigProvider) DatabaseConfig() DBConfig {
	return i.databaseConfig
}

func (i *initConfigProvider) ServerConfig() ServerConfig {
	return i.serverConfig
}

func NewConfigProvider(configPath string) (ConfigProvider, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &initConfigProvider{
		databaseConfig: cfg.DBConfig,
		serverConfig:   cfg.ServerConfig,
	}, nil
}
