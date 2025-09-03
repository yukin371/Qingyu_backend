package config

// Config 存储应用配置
type Config struct {
	Database *DatabaseConfig
	Server   *ServerConfig
	JWT      *JWTConfig
}

// LoadConfig 从环境变量加载配置
func LoadConfig() *Config {
	return &Config{
		Database: LoadDatabaseConfig(),
		Server:   LoadServerConfig(),
		JWT:      LoadJWTConfig(),
	}
}
