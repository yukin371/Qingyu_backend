package validation

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"Qingyu_backend/internal/middleware/core"
	"Qingyu_backend/pkg/errors"

	"github.com/gin-gonic/gin"
)

const (
	// DefaultMaxBodySize 默认请求体大小限制（10MB）
	DefaultMaxBodySize = 10 * 1024 * 1024
	// ValidationPriority Validation中间件的默认优先级
	ValidationPriority = 11
)

// ValidationMiddleware 请求验证中间件
//
// 优先级: 11（业务层，在认证授权之后）
// 用途: 验证请求参数、Content-Type、请求体大小等
type ValidationMiddleware struct {
	config *ValidationConfig
}

// ValidationConfig Validation配置
type ValidationConfig struct {
	// Enabled 是否启用验证
	// 默认: true
	Enabled bool `yaml:"enabled"`

	// MaxBodySize 请求体最大大小（字节）
	// 默认: 10MB
	// 示例: 1048576 (1MB)
	MaxBodySize int64 `yaml:"max_body_size"`

	// AllowedContentTypes 允许的Content-Type列表
	// 默认: ["application/json", "multipart/form-data", "application/x-www-form-urlencoded"]
	// 支持通配符 "*" 表示允许所有Content-Type
	AllowedContentTypes []string `yaml:"allowed_content_types"`

	// RequiredQueryParams 必填的查询参数列表
	// 这些参数必须存在于请求的查询字符串中
	// 示例: ["id", "token"]
	RequiredQueryParams []string `yaml:"required_query_params,omitempty"`

	// RequiredFields 必填的JSON字段列表（对于JSON请求体）
	// 这些字段必须存在于请求体的JSON中
	// 示例: ["name", "email"]
	RequiredFields []string `yaml:"required_fields,omitempty"`
}

// DefaultValidationConfig 返回默认Validation配置
func DefaultValidationConfig() *ValidationConfig {
	return &ValidationConfig{
		Enabled: true,
		MaxBodySize: DefaultMaxBodySize,
		AllowedContentTypes: []string{
			"application/json",
			"multipart/form-data",
			"application/x-www-form-urlencoded",
		},
		RequiredQueryParams: []string{},
		RequiredFields:      []string{},
	}
}

// NewValidationMiddleware 创建新的Validation中间件
func NewValidationMiddleware() *ValidationMiddleware {
	return &ValidationMiddleware{
		config: DefaultValidationConfig(),
	}
}

// Name 返回中间件名称
func (m *ValidationMiddleware) Name() string {
	return "validation"
}

// Priority 返回执行优先级
func (m *ValidationMiddleware) Priority() int {
	return ValidationPriority
}

// Handler 返回Gin处理函数
func (m *ValidationMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 如果验证被禁用，直接跳过
		if !m.config.Enabled {
			c.Next()
			return
		}

		// 1. 验证Content-Type
		if err := m.validateContentType(c); err != nil {
			m.respondWithError(c, err)
			c.Abort()
			return
		}

		// 2. 验证请求体大小
		if err := m.validateBodySize(c); err != nil {
			m.respondWithError(c, err)
			c.Abort()
			return
		}

		// 3. 验证必填查询参数
		if err := m.validateQueryParams(c); err != nil {
			m.respondWithError(c, err)
			c.Abort()
			return
		}

		// 4. 对于有body的请求，读取并验证必填字段
		if c.Request.Body != nil && c.Request.ContentLength > 0 {
			bodyBytes, err := m.readAndValidateBody(c)
			if err != nil {
				m.respondWithError(c, err)
				c.Abort()
				return
			}

			// 替换request body以便后续handler可以读取
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		c.Next()
	}
}

// validateContentType 验证Content-Type
func (m *ValidationMiddleware) validateContentType(c *gin.Context) error {
	// 如果允许所有Content-Type，跳过验证
	if len(m.config.AllowedContentTypes) == 1 && m.config.AllowedContentTypes[0] == "*" {
		return nil
	}

	contentType := c.GetHeader("Content-Type")
	if contentType == "" {
		// GET等没有body的请求不需要Content-Type
		return nil
	}

	// 移除可能的参数（如charset）
	contentType = strings.Split(contentType, ";")[0]
	contentType = strings.TrimSpace(contentType)

	// 检查是否在允许列表中
	for _, allowed := range m.config.AllowedContentTypes {
		if contentType == allowed {
			return nil
		}
	}

	return errors.New(
		errors.InvalidParams,
		fmt.Sprintf("不支持的Content-Type: %s，允许的类型: %v", contentType, m.config.AllowedContentTypes),
	)
}

// validateBodySize 验证请求体大小
func (m *ValidationMiddleware) validateBodySize(c *gin.Context) error {
	if c.Request.ContentLength > m.config.MaxBodySize {
		return errors.New(
			errors.InvalidParams,
			fmt.Sprintf("请求体过大，最大允许 %d 字节", m.config.MaxBodySize),
		)
	}
	return nil
}

// validateQueryParams 验证必填查询参数
func (m *ValidationMiddleware) validateQueryParams(c *gin.Context) error {
	if len(m.config.RequiredQueryParams) == 0 {
		return nil
	}

	var missingParams []string
	for _, param := range m.config.RequiredQueryParams {
		if c.Query(param) == "" {
			missingParams = append(missingParams, param)
		}
	}

	if len(missingParams) > 0 {
		return errors.New(
			errors.InvalidParams,
			fmt.Sprintf("缺少必填查询参数: %s", strings.Join(missingParams, ", ")),
		)
	}

	return nil
}

// readAndValidateBody 读取并验证请求体
func (m *ValidationMiddleware) readAndValidateBody(c *gin.Context) ([]byte, error) {
	// 使用LimitedReader限制读取大小
	limitedReader := io.LimitReader(c.Request.Body, m.config.MaxBodySize+1)
	bodyBytes, err := io.ReadAll(limitedReader)
	if err != nil {
		return nil, errors.New(errors.InvalidParams, "读取请求体失败")
	}

	// 检查是否超过大小限制
	if int64(len(bodyBytes)) > m.config.MaxBodySize {
		return nil, errors.New(
			errors.InvalidParams,
			fmt.Sprintf("请求体过大，最大允许 %d 字节", m.config.MaxBodySize),
		)
	}

	// 对于JSON请求体，验证JSON格式和必填字段
	contentType := c.GetHeader("Content-Type")
	if strings.Contains(contentType, "application/json") && len(bodyBytes) > 0 {
		// 先验证JSON格式
		if !json.Valid(bodyBytes) {
			return nil, errors.New(errors.InvalidParams, "无效的JSON格式")
		}

		// 如果有配置必填字段，验证它们
		if len(m.config.RequiredFields) > 0 {
			if err := m.validateRequiredFields(bodyBytes); err != nil {
				return nil, err
			}
		}
	}

	return bodyBytes, nil
}

// validateRequiredFields 验证JSON必填字段
func (m *ValidationMiddleware) validateRequiredFields(bodyBytes []byte) error {
	var jsonData map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &jsonData); err != nil {
		return errors.New(errors.InvalidParams, "无效的JSON格式")
	}

	var missingFields []string
	for _, field := range m.config.RequiredFields {
		if _, exists := jsonData[field]; !exists {
			missingFields = append(missingFields, field)
		}
	}

	if len(missingFields) > 0 {
		return errors.New(
			errors.InvalidParams,
			fmt.Sprintf("缺少必填字段: %s", strings.Join(missingFields, ", ")),
		)
	}

	return nil
}

// respondWithError 返回错误响应
func (m *ValidationMiddleware) respondWithError(c *gin.Context, err error) {
	// 尝试转换为我们的错误类型
	if appErr, ok := err.(*errors.UnifiedError); ok {
		c.JSON(int(appErr.HTTPStatus), gin.H{
			"code":    appErr.Code,
			"message": appErr.Message,
		})
		return
	}

	// 其他错误，返回通用错误
	c.JSON(http.StatusBadRequest, gin.H{
		"code":    errors.InvalidParams,
		"message": err.Error(),
	})
}

// LoadConfig 从配置加载参数
func (m *ValidationMiddleware) LoadConfig(config map[string]interface{}) error {
	if m.config == nil {
		m.config = &ValidationConfig{}
	}

	// 加载Enabled
	if enabled, ok := config["enabled"].(bool); ok {
		m.config.Enabled = enabled
	}

	// 加载MaxBodySize
	if maxSize, ok := config["max_body_size"].(int64); ok {
		m.config.MaxBodySize = maxSize
	}
	if maxSize, ok := config["max_body_size"].(int); ok {
		m.config.MaxBodySize = int64(maxSize)
	}
	if maxSize, ok := config["max_body_size"].(float64); ok {
		m.config.MaxBodySize = int64(maxSize)
	}

	// 加载AllowedContentTypes
	if contentTypes, ok := config["allowed_content_types"].([]interface{}); ok {
		m.config.AllowedContentTypes = make([]string, len(contentTypes))
		for i, ct := range contentTypes {
			if str, ok := ct.(string); ok {
				m.config.AllowedContentTypes[i] = str
			}
		}
	}

	// 加载RequiredQueryParams
	if queryParams, ok := config["required_query_params"].([]interface{}); ok {
		m.config.RequiredQueryParams = make([]string, len(queryParams))
		for i, param := range queryParams {
			if str, ok := param.(string); ok {
				m.config.RequiredQueryParams[i] = str
			}
		}
	}

	// 加载RequiredFields
	if fields, ok := config["required_fields"].([]interface{}); ok {
		m.config.RequiredFields = make([]string, len(fields))
		for i, field := range fields {
			if str, ok := field.(string); ok {
				m.config.RequiredFields[i] = str
			}
		}
	}

	return nil
}

// ValidateConfig 验证配置有效性
func (m *ValidationMiddleware) ValidateConfig() error {
	if m.config == nil {
		m.config = DefaultValidationConfig()
	}

	// 验证MaxBodySize
	if m.config.MaxBodySize < 0 {
		return fmt.Errorf("max_body_size不能为负数")
	}

	// 验证AllowedContentTypes
	if len(m.config.AllowedContentTypes) == 0 {
		return fmt.Errorf("allowed_content_types不能为空")
	}

	return nil
}

// GetConfig 获取配置
func (m *ValidationMiddleware) GetConfig() *ValidationConfig {
	return m.config
}

// 确保ValidationMiddleware实现了ConfigurableMiddleware接口
var _ core.ConfigurableMiddleware = (*ValidationMiddleware)(nil)
