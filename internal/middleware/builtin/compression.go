package builtin

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"strings"

	"Qingyu_backend/internal/middleware/core"

	"github.com/gin-gonic/gin"
)

// CompressionMiddleware 压缩中间件
//
// 优先级: 12（业务层，最后执行）
// 用途: 使用gzip压缩响应体，减少传输数据量
type CompressionMiddleware struct {
	config *CompressionConfig
}

// CompressionConfig 压缩配置
type CompressionConfig struct {
	// Enabled 是否启用压缩
	// 默认: true
	Enabled bool `yaml:"enabled"`

	// Level 压缩级别
	// 范围: 0-9，0不压缩，9最高压缩比但最慢
	// 默认: 5
	// 建议值: 4-6 之间的平衡点
	Level int `yaml:"level"`

	// MinLength 最小压缩长度（字节）
	// 响应体小于此长度时不压缩
	// 默认: 1024（1KB）
	MinLength int `yaml:"min_length"`

	// Types 需要压缩的Content-Type列表
	// 空列表表示压缩所有类型
	// 默认: ["application/json", "text/html", "text/plain", "text/css", "application/javascript"]
	Types []string `yaml:"types"`

	// ExcludedTypes 不需要压缩的Content-Type列表
	// 优先级高于Types
	// 默认: ["image/*", "video/*", "application/octet-stream"]
	ExcludedTypes []string `yaml:"excluded_types"`
}

// DefaultCompressionConfig 返回默认压缩配置
func DefaultCompressionConfig() *CompressionConfig {
	return &CompressionConfig{
		Enabled:       true,
		Level:         5,
		MinLength:     1024,
		Types:         []string{
			"application/json",
			"text/html",
			"text/plain",
			"text/css",
			"text/javascript",
			"application/javascript",
			"application/xml",
			"text/xml",
		},
		ExcludedTypes: []string{
			"image/*",
			"video/*",
			"audio/*",
			"application/octet-stream",
		},
	}
}

// NewCompressionMiddleware 创建新的压缩中间件
func NewCompressionMiddleware() *CompressionMiddleware {
	return &CompressionMiddleware{
		config: DefaultCompressionConfig(),
	}
}

// Name 返回中间件名称
func (m *CompressionMiddleware) Name() string {
	return "compression"
}

// Priority 返回执行优先级
//
// 返回12，确保压缩在最后执行（在所有业务逻辑之后）
func (m *CompressionMiddleware) Priority() int {
	return 12
}

// Handler 返回Gin处理函数
func (m *CompressionMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 如果禁用压缩，直接跳过
		if !m.config.Enabled {
			c.Next()
			return
		}

		// 检查客户端是否支持gzip
		if !m.acceptsGzip(c) {
			c.Next()
			return
		}

		// 使用response writer包装器
		writer := &gzipWriter{
			ResponseWriter: c.Writer,
			config:        m.config,
		}

		// 设置新的writer
		c.Writer = writer

		// 执行后续处理
		c.Next()

		// 如果需要，压缩响应
		writer.compressIfNeeded()
	}
}

// acceptsGzip 检查客户端是否支持gzip
func (m *CompressionMiddleware) acceptsGzip(c *gin.Context) bool {
	acceptEncoding := c.Request.Header.Get("Accept-Encoding")
	return strings.Contains(acceptEncoding, "gzip")
}

// shouldCompress 检查是否应该压缩响应
func (m *CompressionMiddleware) shouldCompress(contentType string, contentLength int) bool {
	// 检查Content-Length
	if contentLength < m.config.MinLength {
		return false
	}

	// 检查Content-Type
	if len(m.config.Types) > 0 {
		// 如果有指定类型列表，检查是否在列表中
		found := false
		for _, t := range m.config.Types {
			if strings.Contains(contentType, t) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// 检查排除列表
	for _, excluded := range m.config.ExcludedTypes {
		if excluded == "*" {
			return false
		}
		if strings.HasSuffix(excluded, "*") {
			// 通配符匹配
			prefix := strings.TrimSuffix(excluded, "*")
			if strings.HasPrefix(contentType, prefix) {
				return false
			}
		} else if contentType == excluded {
			return false
		}
	}

	return true
}

// LoadConfig 从配置加载参数
//
// 实现ConfigurableMiddleware接口
func (m *CompressionMiddleware) LoadConfig(config map[string]interface{}) error {
	if m.config == nil {
		m.config = &CompressionConfig{}
	}

	// 加载Enabled
	if enabled, ok := config["enabled"].(bool); ok {
		m.config.Enabled = enabled
	}

	// 加载Level
	if level, ok := config["level"].(int); ok {
		m.config.Level = level
	}

	// 加载MinLength
	if minLength, ok := config["min_length"].(int); ok {
		m.config.MinLength = minLength
	}

	// 加载Types
	if types, ok := config["types"].([]interface{}); ok {
		m.config.Types = make([]string, len(types))
		for i, v := range types {
			if str, ok := v.(string); ok {
				m.config.Types[i] = str
			}
		}
	}

	// 加载ExcludedTypes
	if excludedTypes, ok := config["excluded_types"].([]interface{}); ok {
		m.config.ExcludedTypes = make([]string, len(excludedTypes))
		for i, v := range excludedTypes {
			if str, ok := v.(string); ok {
				m.config.ExcludedTypes[i] = str
			}
		}
	}

	return nil
}

// ValidateConfig 验证配置有效性
//
// 实现ConfigurableMiddleware接口
func (m *CompressionMiddleware) ValidateConfig() error {
	if m.config == nil {
		m.config = DefaultCompressionConfig()
	}

	// 验证Level
	if m.config.Level < 0 || m.config.Level > 9 {
		return fmt.Errorf("压缩级别必须在0-9之间，当前值: %d", m.config.Level)
	}

	// 验证MinLength
	if m.config.MinLength < 0 {
		return fmt.Errorf("min_length不能为负数")
	}

	return nil
}

// gzipWriter gzip响应写入器
type gzipWriter struct {
	gin.ResponseWriter
	config  *CompressionConfig
	buffer  *bytes.Buffer
	written bool
}

// Write 写入响应数据
func (w *gzipWriter) Write(data []byte) (int, error) {
	// 第一次写入时，初始化buffer
	if w.buffer == nil {
		w.buffer = bytes.NewBuffer(nil)
	}
	return w.buffer.Write(data)
}

// WriteString 写入字符串响应数据
func (w *gzipWriter) WriteString(s string) (int, error) {
	// 第一次写入时，初始化buffer
	if w.buffer == nil {
		w.buffer = bytes.NewBuffer(nil)
	}
	return w.buffer.WriteString(s)
}

// compressIfNeeded 如果需要，压缩响应
func (w *gzipWriter) compressIfNeeded() {
	// 如果没有数据或已经写入，直接返回
	if w.buffer == nil || w.written {
		return
	}

	// 标记已写入
	w.written = true

	// 获取响应数据
	data := w.buffer.Bytes()
	contentType := w.Header().Get("Content-Type")
	contentLength := len(data)

	// 检查是否需要压缩
	middleware := &CompressionMiddleware{config: w.config}
	if !middleware.shouldCompress(contentType, contentLength) {
		// 不压缩，直接写入原始数据
		w.ResponseWriter.Write(data)
		return
	}

	// 压缩数据
	var compressed bytes.Buffer
	gz, err := gzip.NewWriterLevel(&compressed, w.config.Level)
	if err != nil {
		// 压缩失败，写入原始数据
		w.ResponseWriter.Write(data)
		return
	}

	if _, err := gz.Write(data); err != nil {
		gz.Close()
		w.ResponseWriter.Write(data)
		return
	}

	if err := gz.Close(); err != nil {
		w.ResponseWriter.Write(data)
		return
	}

	// 设置压缩相关的响应头
	w.Header().Set("Content-Encoding", "gzip")
	w.Header().Del("Content-Length") // 删除原始长度

	// 写入压缩后的数据
	w.ResponseWriter.Write(compressed.Bytes())
}

// 确保CompressionMiddleware实现了ConfigurableMiddleware接口
var _ core.ConfigurableMiddleware = (*CompressionMiddleware)(nil)
