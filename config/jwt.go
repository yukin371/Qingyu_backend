package config

// JWTConfig JWT配置
type JWTConfig struct {
	Secret          string
	ExpirationHours int
}

// LoadJWTConfig 加载JWT配置
func LoadJWTConfig() *JWTConfig {
	return &JWTConfig{
		Secret:          getEnv("JWT_SECRET", "qingyu_secret_key"),
		ExpirationHours: getEnvAsInt("JWT_EXPIRATION_HOURS", 24),
	}
}
