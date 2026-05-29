package conf

type Config struct {
	DBConfig     DBConfig     `yaml:"db"`
	ServerConfig ServerConfig `yaml:"server"`
}

type DBConfig struct {
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	Database    string `yaml:"database"`
	DatabaseURL string `yaml:"database_url"`
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	SSLMode     string `yaml:"ssl_mode"`
}

type ServerConfig struct {
	LogLevel string `yaml:"log_level"`
	LogType  string `yaml:"log_type"`
	Addr     string `yaml:"addr"`
}
