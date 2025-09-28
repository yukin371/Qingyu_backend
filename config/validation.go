package config

import (
	"fmt"

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
		return fmt.Errorf("数据库配置不能为空")
	}

	// 使用database.go中的验证方法
	return cfg.Validate()
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
	// 数据库默认值 - 现在使用database.go中的高级配置
	if cfg.Database == nil {
		// 从viper加载数据库配置，如果失败则使用默认配置
		if dbConfig, err := LoadDatabaseConfigFromViper(); err == nil {
			cfg.Database = dbConfig
		} else {
			// 使用默认配置
			cfg.Database = getDefaultDatabaseConfig()
		}
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
