package metrics

import (
	"sync"
	"time"
)

// ServiceMetrics 服务指标
type ServiceMetrics struct {
	mu sync.RWMutex

	// 基础指标
	CallCount    int64         `json:"call_count"`    // 总调用次数
	SuccessCount int64         `json:"success_count"` // 成功次数
	FailureCount int64         `json:"failure_count"` // 失败次数
	TotalTime    time.Duration `json:"total_time"`    // 总响应时间
	AvgTime      time.Duration `json:"avg_time"`      // 平均响应时间
	MaxTime      time.Duration `json:"max_time"`      // 最大响应时间
	MinTime      time.Duration `json:"min_time"`      // 最小响应时间

	// 健康状态
	LastHealthCheck time.Time `json:"last_health_check"` // 最后健康检查时间
	IsHealthy       bool      `json:"is_healthy"`        // 是否健康

	// 服务信息
	ServiceName string    `json:"service_name"` // 服务名称
	Version     string    `json:"version"`      // 版本号
	StartTime   time.Time `json:"start_time"`   // 启动时间
}

// NewServiceMetrics 创建服务指标
func NewServiceMetrics(serviceName, version string) *ServiceMetrics {
	return &ServiceMetrics{
		ServiceName: serviceName,
		Version:     version,
		StartTime:   time.Now(),
		IsHealthy:   false,
		MinTime:     time.Duration(1<<63 - 1), // 设置为最大值
	}
}

// RecordCall 记录一次调用
func (m *ServiceMetrics) RecordCall(duration time.Duration, success bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.CallCount++
	if success {
		m.SuccessCount++
	} else {
		m.FailureCount++
	}

	m.TotalTime += duration
	m.AvgTime = m.TotalTime / time.Duration(m.CallCount)

	if duration > m.MaxTime {
		m.MaxTime = duration
	}
	if duration < m.MinTime {
		m.MinTime = duration
	}
}

// RecordHealthCheck 记录健康检查
func (m *ServiceMetrics) RecordHealthCheck(healthy bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.LastHealthCheck = time.Now()
	m.IsHealthy = healthy
}

// GetSnapshot 获取指标快照
func (m *ServiceMetrics) GetSnapshot() ServiceMetrics {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return ServiceMetrics{
		CallCount:       m.CallCount,
		SuccessCount:    m.SuccessCount,
		FailureCount:    m.FailureCount,
		TotalTime:       m.TotalTime,
		AvgTime:         m.AvgTime,
		MaxTime:         m.MaxTime,
		MinTime:         m.MinTime,
		LastHealthCheck: m.LastHealthCheck,
		IsHealthy:       m.IsHealthy,
		ServiceName:     m.ServiceName,
		Version:         m.Version,
		StartTime:       m.StartTime,
	}
}

// Reset 重置指标
func (m *ServiceMetrics) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.CallCount = 0
	m.SuccessCount = 0
	m.FailureCount = 0
	m.TotalTime = 0
	m.AvgTime = 0
	m.MaxTime = 0
	m.MinTime = time.Duration(1<<63 - 1)
}

// GetSuccessRate 获取成功率
func (m *ServiceMetrics) GetSuccessRate() float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.CallCount == 0 {
		return 0
	}
	return float64(m.SuccessCount) / float64(m.CallCount) * 100
}

// GetFailureRate 获取失败率
func (m *ServiceMetrics) GetFailureRate() float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.CallCount == 0 {
		return 0
	}
	return float64(m.FailureCount) / float64(m.CallCount) * 100
}

// GetUptime 获取运行时间
func (m *ServiceMetrics) GetUptime() time.Duration {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return time.Since(m.StartTime)
}

// MetricsResponse 指标响应结构
type MetricsResponse struct {
	ServiceName     string  `json:"service_name"`
	Version         string  `json:"version"`
	Uptime          string  `json:"uptime"`
	IsHealthy       bool    `json:"is_healthy"`
	CallCount       int64   `json:"call_count"`
	SuccessCount    int64   `json:"success_count"`
	FailureCount    int64   `json:"failure_count"`
	SuccessRate     float64 `json:"success_rate"`
	FailureRate     float64 `json:"failure_rate"`
	AvgTime         string  `json:"avg_time"`
	MaxTime         string  `json:"max_time"`
	MinTime         string  `json:"min_time"`
	LastHealthCheck string  `json:"last_health_check"`
}

// ToResponse 转换为响应格式
func (m *ServiceMetrics) ToResponse() *MetricsResponse {
	snapshot := m.GetSnapshot()

	// 处理 MinTime 显示（如果没有调用，则为 N/A）
	minTimeStr := "N/A"
	if snapshot.CallCount > 0 && snapshot.MinTime < time.Duration(1<<63-1) {
		minTimeStr = snapshot.MinTime.String()
	}

	return &MetricsResponse{
		ServiceName:     snapshot.ServiceName,
		Version:         snapshot.Version,
		Uptime:          m.GetUptime().String(),
		IsHealthy:       snapshot.IsHealthy,
		CallCount:       snapshot.CallCount,
		SuccessCount:    snapshot.SuccessCount,
		FailureCount:    snapshot.FailureCount,
		SuccessRate:     m.GetSuccessRate(),
		FailureRate:     m.GetFailureRate(),
		AvgTime:         snapshot.AvgTime.String(),
		MaxTime:         snapshot.MaxTime.String(),
		MinTime:         minTimeStr,
		LastHealthCheck: snapshot.LastHealthCheck.Format(time.RFC3339),
	}
}
