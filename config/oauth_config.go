package config

import (
	authModel "Qingyu_backend/models/auth"
)

// OAuthConfigManager OAuth配置管理器
type OAuthConfigManager struct {
	configs map[string]*authModel.OAuthConfig
}

// NewOAuthConfigManager 创建OAuth配置管理器
func NewOAuthConfigManager() *OAuthConfigManager {
	return &OAuthConfigManager{
		configs: make(map[string]*authModel.OAuthConfig),
	}
}

// LoadFromEnv 从环境变量加载OAuth配置
func (m *OAuthConfigManager) LoadFromEnv() {
	// Google OAuth
	if googleClientID := GetEnv("GOOGLE_CLIENT_ID", ""); googleClientID != "" {
		m.configs["google"] = &authModel.OAuthConfig{
			Provider:     authModel.OAuthProviderGoogle,
			ClientID:     googleClientID,
			ClientSecret: GetEnv("GOOGLE_CLIENT_SECRET", ""),
			AuthURL:      "https://accounts.google.com/o/oauth2/v2/auth",
			TokenURL:     "https://oauth2.googleapis.com/token",
			Scopes:       "openid email profile",
			Enabled:      true,
		}
	}

	// GitHub OAuth
	if githubClientID := GetEnv("GITHUB_CLIENT_ID", ""); githubClientID != "" {
		m.configs["github"] = &authModel.OAuthConfig{
			Provider:     authModel.OAuthProviderGitHub,
			ClientID:     githubClientID,
			ClientSecret: GetEnv("GITHUB_CLIENT_SECRET", ""),
			AuthURL:      "https://github.com/login/oauth/authorize",
			TokenURL:     "https://github.com/login/oauth/access_token",
			Scopes:       "read:user user:email",
			Enabled:      true,
		}
	}

	// QQ OAuth
	if qqClientID := GetEnv("QQ_CLIENT_ID", ""); qqClientID != "" {
		m.configs["qq"] = &authModel.OAuthConfig{
			Provider:     authModel.OAuthProviderQQ,
			ClientID:     qqClientID,
			ClientSecret: GetEnv("QQ_CLIENT_SECRET", ""),
			AuthURL:      "https://graph.qq.com/oauth2.0/authorize",
			TokenURL:     "https://graph.qq.com/oauth2.0/token",
			Scopes:       "get_user_info",
			Enabled:      true,
		}
	}

	// 微信 OAuth (预留)
	if wechatClientID := GetEnv("WECHAT_CLIENT_ID", ""); wechatClientID != "" {
		m.configs["wechat"] = &authModel.OAuthConfig{
			Provider:     authModel.OAuthProviderWeChat,
			ClientID:     wechatClientID,
			ClientSecret: GetEnv("WECHAT_CLIENT_SECRET", ""),
			AuthURL:      "https://open.weixin.qq.com/connect/qrconnect",
			TokenURL:     "https://api.weixin.qq.com/sns/oauth2/access_token",
			Scopes:       "snsapi_login",
			Enabled:      true,
		}
	}

	// 微博 OAuth (预留)
	if weiboClientID := GetEnv("WEIBO_CLIENT_ID", ""); weiboClientID != "" {
		m.configs["weibo"] = &authModel.OAuthConfig{
			Provider:     authModel.OAuthProviderWeibo,
			ClientID:     weiboClientID,
			ClientSecret: GetEnv("WEIBO_CLIENT_SECRET", ""),
			AuthURL:      "https://api.weibo.com/oauth2/authorize",
			TokenURL:     "https://api.weibo.com/oauth2/access_token",
			Scopes:       "email",
			Enabled:      true,
		}
	}
}

// LoadFromConfig 从配置文件加载OAuth配置
func (m *OAuthConfigManager) LoadFromConfig(cfg *Config) {
	// 从配置文件中读取OAuth配置
	if cfg.OAuth != nil {
		for provider, config := range cfg.OAuth {
			m.configs[provider] = config
		}
	}
}

// GetConfigs 获取所有OAuth配置
func (m *OAuthConfigManager) GetConfigs() map[string]*authModel.OAuthConfig {
	return m.configs
}

// GetConfig 获取指定提供商的OAuth配置
func (m *OAuthConfigManager) GetConfig(provider string) (*authModel.OAuthConfig, bool) {
	config, exists := m.configs[provider]
	return config, exists
}

// IsProviderEnabled 检查提供商是否启用
func (m *OAuthConfigManager) IsProviderEnabled(provider string) bool {
	config, exists := m.configs[provider]
	return exists && config.Enabled
}

// GetEnabledProviders 获取所有启用的提供商
func (m *OAuthConfigManager) GetEnabledProviders() []string {
	providers := make([]string, 0, len(m.configs))
	for provider, config := range m.configs {
		if config.Enabled {
			providers = append(providers, provider)
		}
	}
	return providers
}
