package monitoring

import (
	"fmt"
	"strconv"
	"time"

	"Qingyu_backend/internal/middleware/core"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	// DefaultMetricsNamespace 默认的指标命名空间
	DefaultMetricsNamespace = "qingyu"
	// DefaultMetricsPath 默认的指标暴露路径
	DefaultMetricsPath = "/metrics"
	// MetricsPriority Metrics中间件的默认优先级
	MetricsPriority = 7
)

// MetricsMiddleware Prometheus指标采集中间件
//
// 优先级: 7（监控层，在Logger之后）
// 用途: 采集请求指标，包括计数器、延迟直方图、活跃连接数
type MetricsMiddleware struct {
	config           *MetricsConfig
	registry         prometheus.Registerer
	requestCounter   *prometheus.CounterVec
	requestDuration  *prometheus.HistogramVec
	activeConnections prometheus.Gauge
}

// MetricsConfig Metrics配置
type MetricsConfig struct {
	// Namespace 指标命名空间
	// 默认: "qingyu"
	// 示例: "qingyu", "myapp"
	Namespace string `yaml:"namespace"`

	// MetricsPath 指标暴露路径
	// 默认: "/metrics"
	// 示例: "/metrics", "/admin/metrics"
	MetricsPath string `yaml:"metrics_path"`

	// Enabled 是否启用metrics
	// 默认: true
	Enabled bool `yaml:"enabled"`

	// Buckets 延迟直方图的桶配置（可选）
	// 默认使用prometheus.DefBuckets
	Buckets []float64 `yaml:"buckets,omitempty"`
}

// DefaultMetricsConfig 返回默认Metrics配置
func DefaultMetricsConfig() *MetricsConfig {
	return &MetricsConfig{
		Namespace:   DefaultMetricsNamespace,
		MetricsPath: DefaultMetricsPath,
		Enabled:     true,
		Buckets:     nil, // 使用默认桶
	}
}

// NewMetricsMiddleware 创建新的Metrics中间件
func NewMetricsMiddleware() *MetricsMiddleware {
	return NewMetricsMiddlewareWithRegistry(prometheus.NewRegistry())
}

// NewMetricsMiddlewareWithRegistry 使用自定义registry创建Metrics中间件
func NewMetricsMiddlewareWithRegistry(registry prometheus.Registerer) *MetricsMiddleware {
	config := DefaultMetricsConfig()
	middleware := &MetricsMiddleware{
		config:   config,
		registry: registry,
	}

	// 初始化Prometheus指标
	middleware.initMetrics()

	return middleware
}

// Name 返回中间件名称
func (m *MetricsMiddleware) Name() string {
	return "metrics"
}

// Priority 返回执行优先级
func (m *MetricsMiddleware) Priority() int {
	return MetricsPriority
}

// Handler 返回Gin处理函数
func (m *MetricsMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 如果metrics被禁用，直接跳过
		if !m.config.Enabled {
			c.Next()
			return
		}

		// 增加活跃连接数
		m.activeConnections.Inc()
		defer m.activeConnections.Dec()

		// 记录开始时间
		start := time.Now()

		// 处理请求
		c.Next()

		// 计算请求持续时间
		duration := time.Since(start).Seconds()

		// 更新指标
		m.updateMetrics(c, duration)
	}
}

// initMetrics 初始化Prometheus指标
func (m *MetricsMiddleware) initMetrics() {
	// 确定使用哪个桶配置
	buckets := m.config.Buckets
	if buckets == nil {
		buckets = prometheus.DefBuckets
	}

	// 创建请求计数器
	m.requestCounter = promauto.With(m.registry).NewCounterVec(
		prometheus.CounterOpts{
			Namespace: m.config.Namespace,
			Name:      "http_requests_total",
			Help:      "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status_code"},
	)

	// 创建请求延迟直方图
	m.requestDuration = promauto.With(m.registry).NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: m.config.Namespace,
			Name:      "http_request_duration_seconds",
			Help:      "HTTP request latency in seconds",
			Buckets:   buckets,
		},
		[]string{"method", "endpoint", "status_code"},
	)

	// 创建活跃连接数仪表
	m.activeConnections = promauto.With(m.registry).NewGauge(
		prometheus.GaugeOpts{
			Namespace: m.config.Namespace,
			Name:      "http_active_connections",
			Help:      "Current number of active HTTP connections",
		},
	)
}

// updateMetrics 更新指标
func (m *MetricsMiddleware) updateMetrics(c *gin.Context, duration float64) {
	// 获取请求信息
	method := c.Request.Method
	endpoint := c.FullPath()
	statusCode := strconv.Itoa(c.Writer.Status())

	// 如果FullPath为空（例如404），使用路径
	if endpoint == "" {
		endpoint = c.Request.URL.Path
	}

	// 更新计数器
	m.requestCounter.WithLabelValues(method, endpoint, statusCode).Inc()

	// 更新延迟直方图
	m.requestDuration.WithLabelValues(method, endpoint, statusCode).Observe(duration)
}

// LoadConfig 从配置加载参数
func (m *MetricsMiddleware) LoadConfig(config map[string]interface{}) error {
	if m.config == nil {
		m.config = &MetricsConfig{}
	}

	// 加载Namespace
	if namespace, ok := config["namespace"].(string); ok {
		m.config.Namespace = namespace
	}

	// 加载MetricsPath
	if metricsPath, ok := config["metrics_path"].(string); ok {
		m.config.MetricsPath = metricsPath
	}

	// 加载Enabled
	if enabled, ok := config["enabled"].(bool); ok {
		m.config.Enabled = enabled
	}

	// 加载Buckets（可选）
	if buckets, ok := config["buckets"].([]interface{}); ok {
		m.config.Buckets = make([]float64, len(buckets))
		for i, bucket := range buckets {
			if f, ok := bucket.(float64); ok {
				m.config.Buckets[i] = f
			}
		}
	}

	return nil
}

// ValidateConfig 验证配置有效性
func (m *MetricsMiddleware) ValidateConfig() error {
	if m.config == nil {
		m.config = DefaultMetricsConfig()
	}

	// 验证Namespace
	if m.config.Namespace == "" {
		return fmt.Errorf("namespace不能为空")
	}

	// 验证MetricsPath
	if m.config.MetricsPath == "" {
		return fmt.Errorf("metrics_path不能为空")
	}

	// 验证Buckets（如果提供了）
	if m.config.Buckets != nil {
		if len(m.config.Buckets) == 0 {
			return fmt.Errorf("buckets不能为空数组")
		}
		// 验证桶是递增的
		for i := 1; i < len(m.config.Buckets); i++ {
			if m.config.Buckets[i] <= m.config.Buckets[i-1] {
				return fmt.Errorf("buckets必须是递增的")
			}
		}
	}

	return nil
}

// GetRegistry 获取Prometheus Registry
func (m *MetricsMiddleware) GetRegistry() prometheus.Registerer {
	return m.registry
}

// GetConfig 获取配置
func (m *MetricsMiddleware) GetConfig() *MetricsConfig {
	return m.config
}

// 确保MetricsMiddleware实现了ConfigurableMiddleware接口
var _ core.ConfigurableMiddleware = (*MetricsMiddleware)(nil)
