package v1

import (
	"github.com/gin-gonic/gin"
	"Qingyu_backend/pkg/response"
	"Qingyu_backend/middleware"
)

// VersionAPI 版本信息API
type VersionAPI struct {
	versionRegistry *middleware.VersionRegistry
}

// APIVersionInfo API版本信息
type APIVersionInfo struct {
	Version     string    `json:"version"`      // 版本号 (v1, v2)
	Status      string    `json:"status"`       // 状态 (stable, beta, deprecated)
	Path        string    `json:"path"`         // 路径前缀 (/api/v1)
	ReleaseDate string    `json:"release_date"` // 发布日期
	SunsetDate  *string   `json:"sunset_date,omitempty"`  // 废弃日期（可选）
	Replacement string    `json:"replacement,omitempty"`  // 替代版本（可选）
	Description string    `json:"description"`  // 描述
	DocsURL     string    `json:"docs_url"`      // 文档链接
}

// APIVersionsResponse API版本列表响应
type APIVersionsResponse struct {
	DefaultVersion string           `json:"default_version"` // 默认版本
	Versions       []APIVersionInfo `json:"versions"`        // 所有版本
}

// NewVersionAPI 创建版本信息API
func NewVersionAPI() *VersionAPI {
	registry := middleware.NewVersionRegistry()

	// 注册当前版本
	registry.RegisterVersion(&middleware.VersionConfig{
		Version:     "v1",
		Status:      "stable",
		Path:        "/api/v1",
		Description: "当前稳定版本",
	})

	// 设置默认版本
	registry.SetDefaultVersion("v1")

	return &VersionAPI{
		versionRegistry: registry,
	}
}

// GetVersions 获取所有API版本信息
// GET /api
// 返回所有可用的API版本及其状态
func (api *VersionAPI) GetVersions(c *gin.Context) {
	versions := api.buildVersionInfoList()

	respData := APIVersionsResponse{
		DefaultVersion: api.versionRegistry.GetDefaultVersion(),
		Versions:       versions,
	}

	response.Success(c, respData)
}

// GetVersionInfo 获取指定版本信息
// GET /api/v1
// 返回指定API版本的详细信息
func (api *VersionAPI) GetVersionInfo(c *gin.Context) {
	version := c.Param("version")
	if version == "" {
		version = middleware.DefaultAPIVersion
	}

	// 从注册表获取版本配置
	config, exists := api.versionRegistry.GetVersion(version)
	if !exists {
		response.NotFound(c, "版本不存在")
		return
	}

	versionInfo := api.buildVersionInfo(config)
	response.Success(c, versionInfo)
}

// buildVersionInfoList 构建版本信息列表
func (api *VersionAPI) buildVersionInfoList() []APIVersionInfo {
	versions := api.versionRegistry.GetAllVersions()
	result := make([]APIVersionInfo, 0, len(versions))

	for _, config := range versions {
		result = append(result, api.buildVersionInfo(config))
	}

	return result
}

// buildVersionInfo 构建单个版本信息
func (api *VersionAPI) buildVersionInfo(config *middleware.VersionConfig) APIVersionInfo {
	info := APIVersionInfo{
		Version:     config.Version,
		Status:      config.Status,
		Path:        config.Path,
		ReleaseDate: "2025-01-01", // TODO: 从配置读取
		Description: config.Description,
		DocsURL:     "/swagger/index.html",
	}

	// 如果是废弃状态，添加废除日期和替代版本
	if config.Status == "deprecated" {
		sunsetDate := "2026-12-31" // TODO: 从配置读取
		info.SunsetDate = &sunsetDate
		info.Replacement = "/api/v2" // TODO: 从配置读取
	}

	return info
}

// GetCurrentVersion 获取当前请求的API版本
// GET /api/version
// 返回当前请求使用的API版本
func (api *VersionAPI) GetCurrentVersion(c *gin.Context) {
	version := middleware.GetAPIVersion(c)

	response.Success(c, gin.H{
		"version":     version,
		"description": "当前请求使用的API版本",
	})
}

// RegisterCustomVersion 注册自定义版本（用于测试和扩展）
// 允许在运行时注册新版本
func (api *VersionAPI) RegisterCustomVersion(version, status, path, description string) {
	config := &middleware.VersionConfig{
		Version:     version,
		Status:      status,
		Path:        path,
		Description: description,
	}
	api.versionRegistry.RegisterVersion(config)
}

// GetVersionRegistry 获取版本注册表（供路由注册使用）
func (api *VersionAPI) GetVersionRegistry() *middleware.VersionRegistry {
	return api.versionRegistry
}

// VersionStatusHandler 版本状态处理器（用于健康检查）
// 返回各版本的状态信息
func (api *VersionAPI) VersionStatusHandler(c *gin.Context) {
	versions := api.versionRegistry.GetAllVersions()
	status := make(map[string]interface{})

	for _, config := range versions {
		versionStatus := map[string]interface{}{
			"status":      config.Status,
			"path":        config.Path,
			"description": config.Description,
		}
		status[config.Version] = versionStatus
	}

	response.Success(c, gin.H{
		"versions": status,
		"default":  api.versionRegistry.GetDefaultVersion(),
	})
}

// APIVersionMigrationGuide API版本迁移指南
// GET /api/migration-guide
// 返回版本迁移指南和最佳实践
func (api *VersionAPI) APIVersionMigrationGuide(c *gin.Context) {
	guide := gin.H{
		"versioning_strategy": gin.H{
			"url_based": "使用URL路径指定版本 (推荐): /api/v1/users, /api/v2/users",
			"header_based": "使用HTTP Header指定版本: X-API-Version: v1",
			"priority": "URL路径优先级高于Header",
		},
		"best_practices": []string{
			"优先使用URL路径方式指定版本",
			"客户端应缓存版本信息，减少查询",
			"废弃API会返回相应的响应头，请及时迁移",
			"新功能只在最新稳定版本中提供",
		},
		"response_headers": gin.H{
			"X-API-Deprecated":  "标记API是否已废弃",
			"X-API-Sunset-Date":  "API将移除的日期",
			"X-API-Replacement": "替代的API端点",
			"Warning":           "警告消息",
		},
		"example_requests": gin.H{
			"v1_api": "curl -H 'X-API-Version: v1' http://localhost:8080/api/users",
			"v2_api": "curl -H 'X-API-Version: v2' http://localhost:8080/api/users",
			"url_based": "curl http://localhost:8080/api/v1/users",
		},
	}

	response.Success(c, guide)
}

// InitVersionRoutes 初始化版本路由
// 在router中调用此函数注册版本相关路由
func InitVersionRoutes(router *gin.RouterGroup, api *VersionAPI) {
	// 注册到主路由组（不带版本前缀）
	// GET /api - 获取所有版本信息
	router.GET("", api.GetVersions)

	// GET /api/migration-guide - 获取迁移指南
	router.GET("/migration-guide", api.APIVersionMigrationGuide)

	// GET /api/version - 获取当前请求版本
	router.GET("/version", api.GetCurrentVersion)
}

// InitVersionRoutesWithVersionPrefix 初始化带版本前缀的路由
// 注册到特定版本组
func InitVersionRoutesWithVersionGroup(versionGroup *gin.RouterGroup, api *VersionAPI) {
	// GET /api/v1 - 获取v1版本信息
	versionGroup.GET("", api.GetVersionInfo)

	// GET /api/v1/status - 版本状态（健康检查）
	versionGroup.GET("/status", api.VersionStatusHandler)
}
