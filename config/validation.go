package config

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidateConfig 验证配置
func ValidateConfig(cfg *Config) error {
	if err := validate.Struct(cfg); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}
	return validateConfigDetails(cfg)
}

// validateConfigDetails 验证配置细节
func validateConfigDetails(cfg *Config) error {
	if err := validateDatabaseConfig(cfg.Database); err != nil {
		return fmt.Errorf("database config validation failed: %w", err)
	}

	if err := validateServerConfig(cfg.Server); err != nil {
		return fmt.Errorf("server config validation failed: %w", err)
	}

	if err := validateJWTConfig(cfg.JWT); err != nil {
		return fmt.Errorf("JWT config validation failed: %w", err)
	}

	if err := validateAIConfig(cfg.AI); err != nil {
		return fmt.Errorf("AI config validation failed: %w", err)
	}

	return nil
}

// validateDatabaseConfig 验证数据库配置
func validateDatabaseConfig(cfg *DatabaseConfig) error {
	if cfg == nil {
		return fmt.Errorf("database config is required")
	}

	if cfg.MongoURI == "" {
		return fmt.Errorf("database URI is required")
	}

	if cfg.DBName == "" {
		return fmt.Errorf("database name is required")
	}

	if cfg.ConnectTimeout < time.Second {
		return fmt.Errorf("connect timeout must be at least 1 second")
	}

	if cfg.MaxPoolSize < cfg.MinPoolSize {
		return fmt.Errorf("max pool size must be greater than or equal to min pool size")
	}

	if cfg.MaxConnIdleTime < time.Second {
		return fmt.Errorf("max connection idle time must be at least 1 second")
	}

	return nil
}

// validateServerConfig 验证服务器配置
func validateServerConfig(cfg *ServerConfig) error {
	if cfg == nil {
		return fmt.Errorf("server config is required")
	}

	if cfg.Port == "" {
		return fmt.Errorf("server port is required")
	}

	if cfg.Mode != "debug" && cfg.Mode != "release" {
		return fmt.Errorf("server mode must be either 'debug' or 'release'")
	}

	return nil
}

// validateJWTConfig 验证JWT配置
func validateJWTConfig(cfg *JWTConfig) error {
	if cfg == nil {
		return fmt.Errorf("JWT config is required")
	}

	if cfg.Secret == "" {
		return fmt.Errorf("JWT secret is required")
	}

	if cfg.ExpirationHours <= 0 {
		return fmt.Errorf("JWT expiration hours must be positive")
	}

	return nil
}

// validateAIConfig 验证AI配置
func validateAIConfig(cfg *AIConfig) error {
	if cfg == nil {
		return fmt.Errorf("AI config is required")
	}

	if cfg.APIKey == "" {
		return fmt.Errorf("AI API key is required")
	}

	if cfg.BaseURL == "" {
		return fmt.Errorf("AI base URL is required")
	}

	if cfg.MaxTokens <= 0 {
		return fmt.Errorf("AI max tokens must be positive")
	}

	if cfg.Temperature < 0 || cfg.Temperature > 10 {
		return fmt.Errorf("AI temperature must be between 0 and 10")
	}

	return nil
}

// setConfigDefaults 设置配置默认值
func setConfigDefaults(cfg *Config) {
	// 数据库默认值
	if cfg.Database == nil {
		cfg.Database = &DatabaseConfig{}
	}
	if cfg.Database.ConnectTimeout == 0 {
		cfg.Database.ConnectTimeout = 10 * time.Second
	}
	if cfg.Database.MaxPoolSize == 0 {
		cfg.Database.MaxPoolSize = 100
	}
	if cfg.Database.MinPoolSize == 0 {
		cfg.Database.MinPoolSize = 10
	}
	if cfg.Database.MaxConnIdleTime == 0 {
		cfg.Database.MaxConnIdleTime = 60 * time.Second
	}
	if !cfg.Database.RetryWrites {
		cfg.Database.RetryWrites = true
	}
	if !cfg.Database.RetryReads {
		cfg.Database.RetryReads = true
	}

	// 服务器默认值
	if cfg.Server == nil {
		cfg.Server = &ServerConfig{}
	}
	if cfg.Server.Port == "" {
		cfg.Server.Port = "8080"
	}
	if cfg.Server.Mode == "" {
		cfg.Server.Mode = "debug"
	}

	// JWT默认值
	if cfg.JWT == nil {
		cfg.JWT = &JWTConfig{}
	}
	if cfg.JWT.ExpirationHours == 0 {
		cfg.JWT.ExpirationHours = 24
	}

	// AI默认值
	if cfg.AI == nil {
		cfg.AI = &AIConfig{}
	}
	if cfg.AI.BaseURL == "" {
		cfg.AI.BaseURL = "https://api.openai.com/v1"
	}
	if cfg.AI.MaxTokens == 0 {
		cfg.AI.MaxTokens = 2000
	}
	if cfg.AI.Temperature == 0 {
		cfg.AI.Temperature = 7
	}
}