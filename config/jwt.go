package config

import "time"

// GetJWTConfigEnhanced 获取增强的JWT配置（从现有JWTConfig转换）
func GetJWTConfigEnhanced() *JWTConfigEnhanced {
	// 从全局配置获取基础JWT配置
	var baseJWT *JWTConfig
	if GlobalConfig != nil && GlobalConfig.JWT != nil {
		baseJWT = GlobalConfig.JWT
	} else {
		// 默认配置
		baseJWT = &JWTConfig{
			Secret:          "qingyu-secret-key-change-in-production",
			ExpirationHours: 24,
		}
	}

	// 转换为增强配置
	return &JWTConfigEnhanced{
		SecretKey:       baseJWT.Secret,
		Issuer:          "qingyu-backend",
		Expiration:      time.Duration(baseJWT.ExpirationHours) * time.Hour,
		RefreshDuration: 7 * 24 * time.Hour, // 默认7天
	}
}

// JWTConfigEnhanced 增强的JWT配置（用于Auth模块）
type JWTConfigEnhanced struct {
	SecretKey       string        // JWT密钥
	Issuer          string        // 签发者
	Expiration      time.Duration // Token过期时间
	RefreshDuration time.Duration // 刷新Token过期时间
}
