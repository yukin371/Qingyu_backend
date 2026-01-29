package middleware

import (
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

// 版本上下文key
const (
	APIVersionKey      = "api_version"
	APIVersionHeader   = "X-API-Version"
	DefaultAPIVersion  = "v1"
)

// 版本路径正则表达式
var versionPathRegex = regexp.MustCompile(`^/api/([^/]+)`)

// VersionRoutingMiddleware 版本路由中间件
// 从URL路径或Header中提取API版本信息，并存入context
func VersionRoutingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 优先从URL路径提取版本
		version := extractVersionFromPath(c.Request.URL.Path)

		// 2. 如果URL中没有版本，尝试从Header获取
		if version == "" {
			version = c.GetHeader(APIVersionHeader)
		}

		// 3. 如果都没有，使用默认版本
		if version == "" {
			version = DefaultAPIVersion
		}

		// 4. 验证版本格式（简单验证：v+数字）
		if !isValidVersion(version) {
			version = DefaultAPIVersion
		}

		// 5. 将版本信息存入context
		c.Set(APIVersionKey, version)

		c.Next()
	}
}

// extractVersionFromPath 从URL路径提取版本号
// 例如: /api/v1/users -> v1
func extractVersionFromPath(path string) string {
	matches := versionPathRegex.FindStringSubmatch(path)
	if len(matches) >= 2 {
		return matches[1]
	}
	return ""
}

// isValidVersion 验证版本格式
// 支持格式: v1, v2, v1.1 等
func isValidVersion(version string) bool {
	if version == "" {
		return false
	}
	// 简单验证：以v或V开头，后面跟数字
	matched, _ := regexp.MatchString(`^[vV][0-9]+(\.[0-9]+)?$`, version)
	return matched
}

// GetAPIVersion 从context获取API版本
func GetAPIVersion(c *gin.Context) string {
	if version, exists := c.Get(APIVersionKey); exists {
		if v, ok := version.(string); ok {
			return v
		}
	}
	return DefaultAPIVersion
}

// SetAPIVersion 设置API版本到context
func SetAPIVersion(c *gin.Context, version string) {
	c.Set(APIVersionKey, version)
}

// VersionConfig 版本配置
type VersionConfig struct {
	Version     string `json:"version" yaml:"version"`         // 版本号 (v1, v2)
	Status      string `json:"status" yaml:"status"`           // 状态 (stable, beta, deprecated, sunset)
	Path        string `json:"path" yaml:"path"`               // 路径前缀
	Description string `json:"description" yaml:"description"` // 描述
}

// VersionRegistry 版本注册表
type VersionRegistry struct {
	defaultVersion string
	versions       map[string]*VersionConfig
}

// NewVersionRegistry 创建版本注册表
func NewVersionRegistry() *VersionRegistry {
	return &VersionRegistry{
		defaultVersion: DefaultAPIVersion,
		versions:       make(map[string]*VersionConfig),
	}
}

// RegisterVersion 注册版本
func (r *VersionRegistry) RegisterVersion(config *VersionConfig) {
	r.versions[config.Version] = config
}

// GetVersion 获取版本配置
func (r *VersionRegistry) GetVersion(version string) (*VersionConfig, bool) {
	config, exists := r.versions[version]
	return config, exists
}

// SetDefaultVersion 设置默认版本
func (r *VersionRegistry) SetDefaultVersion(version string) {
	r.defaultVersion = version
}

// GetDefaultVersion 获取默认版本
func (r *VersionRegistry) GetDefaultVersion() string {
	return r.defaultVersion
}

// GetAllVersions 获取所有版本
func (r *VersionRegistry) GetAllVersions() []*VersionConfig {
	versions := make([]*VersionConfig, 0, len(r.versions))
	for _, config := range r.versions {
		versions = append(versions, config)
	}
	return versions
}

// GetVersionFromPath 从路径获取版本配置
func (r *VersionRegistry) GetVersionFromPath(path string) (*VersionConfig, bool) {
	version := extractVersionFromPath(path)
	if version == "" {
		return r.GetVersion(r.defaultVersion)
	}
	return r.GetVersion(version)
}

// IsVersionAvailable 检查版本是否可用
func (r *VersionRegistry) IsVersionAvailable(version string) bool {
	_, exists := r.versions[version]
	return exists
}

// ParseVersionFromHeader 从Header解析版本
// 支持格式: v1, V1, latest
func ParseVersionFromHeader(header string) string {
	if header == "" {
		return ""
	}

	// 转为小写处理
	header = strings.ToLower(strings.TrimSpace(header))

	// 特殊值: latest
	if header == "latest" {
		// TODO: 实现获取最新稳定版本的逻辑
		return DefaultAPIVersion
	}

	// 标准格式: v1, v2 等
	if strings.HasPrefix(header, "v") {
		return header
	}

	// 如果只有数字，添加v前缀
	// 例如: 1 -> v1
	return "v" + header
}
