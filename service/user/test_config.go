package user

import (
	"time"

	"Qingyu_backend/config"
)

// TestConfig 集成测试配置
type TestConfig struct {
	JWTSecret           string
	JWTExpiration       time.Duration
	JWTRefreshDuration  time.Duration
	MongoURI            string
	DatabaseName        string
	RedisAddr           string
	EmailServiceEnabled bool
}

// GetTestConfig 获取测试配置
func GetTestConfig() *TestConfig {
	return &TestConfig{
		JWTSecret:           "test-secret-key-for-integration-testing-only-do-not-use-in-production",
		JWTExpiration:       24 * time.Hour,
		JWTRefreshDuration:  7 * 24 * time.Hour,
		MongoURI:            "mongodb://localhost:27017",
		DatabaseName:        "qingyu_test",
		RedisAddr:           "localhost:6379",
		EmailServiceEnabled: false, // 测试中禁用真实邮件发送
	}
}

// GetTestJWTConfig 获取测试用JWT配置
func GetTestJWTConfig() *config.JWTConfigEnhanced {
	testCfg := GetTestConfig()
	return &config.JWTConfigEnhanced{
		SecretKey:       testCfg.JWTSecret,
		Issuer:          "qingyu-test",
		Expiration:      testCfg.JWTExpiration,
		RefreshDuration: testCfg.JWTRefreshDuration,
	}
}
