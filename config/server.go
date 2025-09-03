package config

// ServerConfig 服务器配置
type ServerConfig struct {
	Port string
	Mode string
}

// LoadServerConfig 加载服务器配置
func LoadServerConfig() *ServerConfig {
	return &ServerConfig{
		Port: getEnv("SERVER_PORT", "8080"),
		Mode: getEnv("SERVER_MODE", "debug"),
	}
}
