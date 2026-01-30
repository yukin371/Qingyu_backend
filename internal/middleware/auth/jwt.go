package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"Qingyu_backend/internal/middleware/core"
	appErrors "Qingyu_backend/pkg/errors"
)

const (
	// AuthPriority Auth中间件的默认优先级
	AuthPriority = 9
	// DefaultTokenHeader 默认Token所在Header
	DefaultTokenHeader = "Authorization"
	// DefaultTokenPrefix 默认Token前缀
	DefaultTokenPrefix = "Bearer"
)

// JWTAuthMiddleware JWT认证中间件
//
// 优先级: 9（认证层，在监控之后，权限之前）
// 用途: 验证JWT Token，提取用户信息
type JWTAuthMiddleware struct {
	config     *JWTConfig
	jwtManager JWTManager
	blacklist  Blacklist
	logger     *zap.Logger
}

// JWTConfig JWT配置
type JWTConfig struct {
	// Enabled 是否启用JWT认证
	// 默认: true
	Enabled bool `yaml:"enabled"`

	// Secret JWT密钥
	// 默认: ""（必须设置）
	// 示例: "your-secret-key"
	Secret string `yaml:"secret"`

	// AccessExpiration Access Token过期时间
	// 默认: 2小时
	// 示例: "2h", "30m"
	AccessExpiration time.Duration `yaml:"access_expiration"`

	// RefreshExpiration Refresh Token过期时间
	// 默认: 7天
	// 示例: "7d", "168h"
	RefreshExpiration time.Duration `yaml:"refresh_expiration"`

	// Issuer 签发者
	// 默认: "qingyu"
	Issuer string `yaml:"issuer"`

	// TokenHeader Token所在Header
	// 默认: "Authorization"
	TokenHeader string `yaml:"token_header"`

	// TokenPrefix Token前缀
	// 默认: "Bearer"
	TokenPrefix string `yaml:"token_prefix"`

	// SkipPaths 跳过认证的路径
	// 默认: ["/health", "/metrics"]
	SkipPaths []string `yaml:"skip_paths"`
}

// DefaultJWTConfig 返回默认JWT配置
func DefaultJWTConfig() *JWTConfig {
	return &JWTConfig{
		Enabled:          true,
		Secret:           "",
		AccessExpiration: DefaultAccessExpiration,
		RefreshExpiration: DefaultRefreshExpiration,
		Issuer:           DefaultIssuer,
		TokenHeader:      DefaultTokenHeader,
		TokenPrefix:      DefaultTokenPrefix,
		SkipPaths: []string{
			"/health",
			"/metrics",
		},
	}
}

// NewJWTAuthMiddleware 创建JWT认证中间件
func NewJWTAuthMiddleware(jwtManager JWTManager, blacklist Blacklist, logger *zap.Logger) *JWTAuthMiddleware {
	config := DefaultJWTConfig()

	if logger == nil {
		logger, _ = zap.NewDevelopment()
	}

	return &JWTAuthMiddleware{
		config:     config,
		jwtManager: jwtManager,
		blacklist:  blacklist,
		logger:     logger,
	}
}

// Name 返回中间件名称
func (m *JWTAuthMiddleware) Name() string {
	return "jwt_auth"
}

// Priority 返回执行优先级
func (m *JWTAuthMiddleware) Priority() int {
	return AuthPriority
}

// Handler 返回Gin处理函数
func (m *JWTAuthMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 如果认证被禁用，直接跳过
		if !m.config.Enabled {
			c.Next()
			return
		}

		// 检查是否跳过该路径
		if m.shouldSkipPath(c.Request.URL.Path) {
			c.Next()
			return
		}

		// 提取Token
		token, err := m.extractToken(c)
		if err != nil {
			m.respondWithError(c, err)
			c.Abort()
			return
		}

		// 验证Token
		claims, err := m.validateToken(token)
		if err != nil {
			m.respondWithError(c, err)
			c.Abort()
			return
		}

		// 检查Token是否在黑名单
		if m.blacklist != nil {
			ctx := c.Request.Context()
			isBlacklisted, err := m.blacklist.IsBlacklisted(ctx, token)
			if err != nil {
				m.logger.Error("Failed to check blacklist",
					zap.String("token", token),
					zap.Error(err),
				)
			}
			if isBlacklisted {
				m.respondWithError(c, appErrors.New(
					appErrors.TokenRevoked,
					"Token已被撤销",
				))
				c.Abort()
				return
			}
		}

		// 将用户信息注入到上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("roles", claims.Roles)
		c.Set("token_type", claims.TokenType)

		c.Next()
	}
}

// extractToken 从请求中提取Token
func (m *JWTAuthMiddleware) extractToken(c *gin.Context) (string, error) {
	authHeader := c.GetHeader(m.config.TokenHeader)

	// 检查Header是否为空
	if authHeader == "" {
		return "", errors.New("2010")
	}

	// 检查Token格式
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != m.config.TokenPrefix {
		return "", errors.New("2009")
	}

	token := parts[1]
	if token == "" {
		return "", errors.New("2009")
	}

	return token, nil
}

// validateToken 验证Token
func (m *JWTAuthMiddleware) validateToken(token string) (*Claims, error) {
	claims, err := m.jwtManager.ValidateToken(token)
	if err != nil {
		// 判断错误类型
		if strings.Contains(err.Error(), "expired") {
			return nil, errors.New("2007")
		}
		return nil, errors.New("2008")
	}

	return claims, nil
}

// respondWithError 返回错误响应
func (m *JWTAuthMiddleware) respondWithError(c *gin.Context, err error) {
	// 解析错误码和消息
	var code string
	var message string
	var httpStatus int

	errStr := err.Error()
	switch errStr {
	case "2007":
		code = "2007"
		message = appErrors.GetDefaultMessage(appErrors.TokenExpired)
		httpStatus = appErrors.GetHTTPStatus(appErrors.TokenExpired)
	case "2008":
		code = "2008"
		message = appErrors.GetDefaultMessage(appErrors.TokenInvalid)
		httpStatus = appErrors.GetHTTPStatus(appErrors.TokenInvalid)
	case "2009":
		code = "2009"
		message = appErrors.GetDefaultMessage(appErrors.TokenFormatError)
		httpStatus = appErrors.GetHTTPStatus(appErrors.TokenFormatError)
	case "2010":
		code = "2010"
		message = appErrors.GetDefaultMessage(appErrors.TokenMissing)
		httpStatus = appErrors.GetHTTPStatus(appErrors.TokenMissing)
	case "2016":
		code = "2016"
		message = appErrors.GetDefaultMessage(appErrors.TokenRevoked)
		httpStatus = appErrors.GetHTTPStatus(appErrors.TokenRevoked)
	default:
		code = "2008"
		message = "认证失败"
		httpStatus = http.StatusUnauthorized
	}

	c.JSON(httpStatus, gin.H{
		"code":    code,
		"message": message,
	})
}

// shouldSkipPath 检查是否应该跳过认证
func (m *JWTAuthMiddleware) shouldSkipPath(path string) bool {
	for _, skipPath := range m.config.SkipPaths {
		// 支持前缀匹配
		if len(path) >= len(skipPath) && path[:len(skipPath)] == skipPath {
			return true
		}
	}
	return false
}

// LoadConfig 从配置加载参数
func (m *JWTAuthMiddleware) LoadConfig(config map[string]interface{}) error {
	if m.config == nil {
		m.config = &JWTConfig{}
	}

	// 加载Enabled
	if enabled, ok := config["enabled"].(bool); ok {
		m.config.Enabled = enabled
	}

	// 加载Secret
	if secret, ok := config["secret"].(string); ok {
		m.config.Secret = secret
		// 重新创建JWT管理器
		if m.jwtManager != nil && m.config.Secret != "" {
			jwtManager, err := NewJWTManager(
				m.config.Secret,
				m.config.AccessExpiration,
				m.config.RefreshExpiration,
			)
			if err == nil {
				m.jwtManager = jwtManager
			}
		}
	}

	// 加载AccessExpiration
	if accessExp, ok := config["access_expiration"].(int); ok {
		m.config.AccessExpiration = time.Duration(accessExp) * time.Second
	}
	if accessExp, ok := config["access_expiration"].(float64); ok {
		m.config.AccessExpiration = time.Duration(accessExp) * time.Second
	}

	// 加载RefreshExpiration
	if refreshExp, ok := config["refresh_expiration"].(int); ok {
		m.config.RefreshExpiration = time.Duration(refreshExp) * time.Second
	}
	if refreshExp, ok := config["refresh_expiration"].(float64); ok {
		m.config.RefreshExpiration = time.Duration(refreshExp) * time.Second
	}

	// 加载Issuer
	if issuer, ok := config["issuer"].(string); ok {
		m.config.Issuer = issuer
	}

	// 加载TokenHeader
	if tokenHeader, ok := config["token_header"].(string); ok {
		m.config.TokenHeader = tokenHeader
	}

	// 加载TokenPrefix
	if tokenPrefix, ok := config["token_prefix"].(string); ok {
		m.config.TokenPrefix = tokenPrefix
	}

	// 加载SkipPaths
	if skipPaths, ok := config["skip_paths"].([]interface{}); ok {
		m.config.SkipPaths = make([]string, len(skipPaths))
		for i, v := range skipPaths {
			if str, ok := v.(string); ok {
				m.config.SkipPaths[i] = str
			}
		}
	}

	return nil
}

// ValidateConfig 验证配置有效性
func (m *JWTAuthMiddleware) ValidateConfig() error {
	if m.config == nil {
		m.config = DefaultJWTConfig()
	}

	// 验证Secret
	if m.config.Secret == "" {
		return errors.New("secret不能为空")
	}

	// 验证AccessExpiration
	if m.config.AccessExpiration <= 0 {
		return errors.New("access_expiration必须大于0")
	}

	// 验证RefreshExpiration
	if m.config.RefreshExpiration <= 0 {
		return errors.New("refresh_expiration必须大于0")
	}

	return nil
}

// GetConfig 获取配置
func (m *JWTAuthMiddleware) GetConfig() *JWTConfig {
	return m.config
}

// GetJWTManager 获取JWT管理器
func (m *JWTAuthMiddleware) GetJWTManager() JWTManager {
	return m.jwtManager
}

// GetBlacklist 获取黑名单
func (m *JWTAuthMiddleware) GetBlacklist() Blacklist {
	return m.blacklist
}

// 确保JWTAuthMiddleware实现了ConfigurableMiddleware接口
var _ core.ConfigurableMiddleware = (*JWTAuthMiddleware)(nil)
