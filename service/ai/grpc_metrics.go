package ai

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

// 服务类型常量
const (
	ServiceExecuteAgent         = "ExecuteAgent"
	ServiceGenerateOutline      = "GenerateOutline"
	ServiceGenerateCharacters   = "GenerateCharacters"
	ServiceGeneratePlot         = "GeneratePlot"
	ServiceExecuteCreativeWorkflow = "ExecuteCreativeWorkflow"
	ServiceHealthCheck          = "HealthCheck"
)

// GRPCMetrics gRPC调用统计
type GRPCMetrics struct {
	mu    sync.RWMutex
	calls map[string]*ServiceStats
	perf  map[string]*PerformanceStats
	quota *QuotaMetrics // 配额监控
}

// ServiceStats 服务调用统计
type ServiceStats struct {
	Total   int64
	Success int64
	Failed  int64
}

// PerformanceStats 性能统计
type PerformanceStats struct {
	mu           sync.RWMutex
	TotalLatency int64     // 总延迟（毫秒）
	Count        int64     // 请求数
	Timeouts     int64     // 超时次数
	Retries      int64     // 重试次数
	MinLatency   int64     // 最小延迟
	MaxLatency   int64     // 最大延迟
	latencies    []int64   // 延迟记录（用于计算百分位）
	maxLatencies int       // 保留的延迟记录数量
}

// NewGRPCMetrics 创建新的gRPC指标收集器
func NewGRPCMetrics() *GRPCMetrics {
	return &GRPCMetrics{
		calls: make(map[string]*ServiceStats),
		perf:  make(map[string]*PerformanceStats),
	}
}

// RecordCall 记录调用
func (m *GRPCMetrics) RecordCall(serviceName string, success bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	stats, exists := m.calls[serviceName]
	if !exists {
		stats = &ServiceStats{}
		m.calls[serviceName] = stats
	}

	stats.Total++
	if success {
		stats.Success++
	} else {
		stats.Failed++
	}
}

// GetStats 获取所有统计信息
func (m *GRPCMetrics) GetStats() map[string]*ServiceStats {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// 创建副本避免外部修改
	result := make(map[string]*ServiceStats, len(m.calls))
	for k, v := range m.calls {
		result[k] = &ServiceStats{
			Total:   v.Total,
			Success: v.Success,
			Failed:  v.Failed,
		}
	}

	return result
}

// GetStatsByService 获取指定服务的统计信息
func (m *GRPCMetrics) GetStatsByService(serviceName string) (*ServiceStats, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	stats, exists := m.calls[serviceName]
	if !exists {
		return nil, fmt.Errorf("service %s not found", serviceName)
	}

	return &ServiceStats{
		Total:   stats.Total,
		Success: stats.Success,
		Failed:  stats.Failed,
	}, nil
}

// GetSuccessRate 获取成功率
func (s *ServiceStats) GetSuccessRate() float64 {
	if s.Total == 0 {
		return 0
	}
	return float64(s.Success) / float64(s.Total) * 100
}

// RecordLatency 记录延迟
func (m *GRPCMetrics) RecordLatency(serviceName string, duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	perf, exists := m.perf[serviceName]
	if !exists {
		perf = &PerformanceStats{
			MinLatency:   -1, // 初始化为-1表示未设置
			maxLatencies: 1000, // 默认保留1000条延迟记录
			latencies:    make([]int64, 0, 1000),
		}
		m.perf[serviceName] = perf
	}

	perf.mu.Lock()
	defer perf.mu.Unlock()

	latencyMs := duration.Milliseconds()
	perf.TotalLatency += latencyMs
	perf.Count++

	// 更新最小和最大延迟
	if perf.MinLatency == -1 || latencyMs < perf.MinLatency {
		perf.MinLatency = latencyMs
	}
	if latencyMs > perf.MaxLatency {
		perf.MaxLatency = latencyMs
	}

	// 记录延迟用于百分位计算
	if len(perf.latencies) < perf.maxLatencies {
		perf.latencies = append(perf.latencies, latencyMs)
	}
}

// GetLatencyStats 获取延迟统计
func (m *GRPCMetrics) GetLatencyStats(serviceName string) (*LatencyStats, error) {
	m.mu.RLock()
	perf, exists := m.perf[serviceName]
	m.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("service %s not found", serviceName)
	}

	perf.mu.RLock()
	defer perf.mu.RUnlock()

	stats := &LatencyStats{
		Average: perf.GetAverage(),
		Min:     perf.MinLatency,
		Max:     perf.MaxLatency,
		Count:   perf.Count,
	}

	// 计算百分位
	if len(perf.latencies) > 0 {
		stats.P50 = calculatePercentile(perf.latencies, 50)
		stats.P95 = calculatePercentile(perf.latencies, 95)
		stats.P99 = calculatePercentile(perf.latencies, 99)
	}

	return stats, nil
}

// GetAllLatencyStats 获取所有服务的延迟统计
func (m *GRPCMetrics) GetAllLatencyStats() map[string]*LatencyStats {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make(map[string]*LatencyStats)
	for serviceName := range m.perf {
		if stats, err := m.GetLatencyStats(serviceName); err == nil {
			result[serviceName] = stats
		}
	}

	return result
}

// RecordTimeout 记录超时
func (m *GRPCMetrics) RecordTimeout(serviceName string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	perf, exists := m.perf[serviceName]
	if !exists {
		perf = &PerformanceStats{
			MinLatency:   -1,
			maxLatencies: 1000,
			latencies:    make([]int64, 0, 1000),
		}
		m.perf[serviceName] = perf
	}

	perf.mu.Lock()
	perf.Timeouts++
	perf.mu.Unlock()
}

// RecordRetry 记录重试
func (m *GRPCMetrics) RecordRetry(serviceName string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	perf, exists := m.perf[serviceName]
	if !exists {
		perf = &PerformanceStats{
			MinLatency:   -1,
			maxLatencies: 1000,
			latencies:    make([]int64, 0, 1000),
		}
		m.perf[serviceName] = perf
	}

	perf.mu.Lock()
	perf.Retries++
	perf.mu.Unlock()
}

// GetTimeoutCount 获取超时次数
func (m *GRPCMetrics) GetTimeoutCount(serviceName string) int64 {
	m.mu.RLock()
	perf, exists := m.perf[serviceName]
	m.mu.RUnlock()

	if !exists {
		return 0
	}

	perf.mu.RLock()
	defer perf.mu.RUnlock()
	return perf.Timeouts
}

// GetRetryCount 获取重试次数
func (m *GRPCMetrics) GetRetryCount(serviceName string) int64 {
	m.mu.RLock()
	perf, exists := m.perf[serviceName]
	m.mu.RUnlock()

	if !exists {
		return 0
	}

	perf.mu.RLock()
	defer perf.mu.RUnlock()
	return perf.Retries
}

// GetAverage 获取平均延迟
func (p *PerformanceStats) GetAverage() float64 {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.Count == 0 {
		return 0
	}
	return float64(p.TotalLatency) / float64(p.Count)
}

// LatencyStats 延迟统计信息
type LatencyStats struct {
	Average float64 // 平均延迟（毫秒）
	Min     int64   // 最小延迟（毫秒）
	Max     int64   // 最大延迟（毫秒）
	P50     int64   // P50延迟（毫秒）
	P95     int64   // P95延迟（毫秒）
	P99     int64   // P99延迟（毫秒）
	Count   int64   // 请求数
}

// calculatePercentile 计算百分位
func calculatePercentile(latencies []int64, percentile int) int64 {
	if len(latencies) == 0 {
		return 0
	}

	// 复制并排序
	sorted := make([]int64, len(latencies))
	copy(sorted, latencies)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i] < sorted[j]
	})

	// 计算索引
	index := (len(sorted) * percentile) / 100
	if index >= len(sorted) {
		index = len(sorted) - 1
	}

	return sorted[index]
}

// FormatReport 格式化统计报告
func (m *GRPCMetrics) FormatReport() string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var report string
	report += "=== gRPC调用统计报告 ===\n\n"

	// 服务调用统计
	report += "## 服务调用统计\n"
	for name, stats := range m.calls {
		report += fmt.Sprintf("%s:\n", name)
		report += fmt.Sprintf("  总调用: %d\n", stats.Total)
		report += fmt.Sprintf("  成功: %d\n", stats.Success)
		report += fmt.Sprintf("  失败: %d\n", stats.Failed)
		report += fmt.Sprintf("  成功率: %.2f%%\n", stats.GetSuccessRate())
		report += "\n"
	}

	// 性能统计
	report += "## 性能统计\n"
	for name, perf := range m.perf {
		perf.mu.RLock()
		report += fmt.Sprintf("%s:\n", name)
		report += fmt.Sprintf("  平均延迟: %.2fms\n", perf.GetAverage())
		report += fmt.Sprintf("  最小延迟: %dms\n", perf.MinLatency)
		report += fmt.Sprintf("  最大延迟: %dms\n", perf.MaxLatency)
		report += fmt.Sprintf("  请求数: %d\n", perf.Count)
		report += fmt.Sprintf("  超时次数: %d\n", perf.Timeouts)
		report += fmt.Sprintf("  重试次数: %d\n", perf.Retries)
		perf.mu.RUnlock()
		report += "\n"
	}

	return report
}

// Reset 重置统计信息
func (m *GRPCMetrics) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.calls = make(map[string]*ServiceStats)
	m.perf = make(map[string]*PerformanceStats)
}

// ResetService 重置指定服务的统计信息
func (m *GRPCMetrics) ResetService(serviceName string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.calls, serviceName)
	delete(m.perf, serviceName)
}

// ============ 配额监控 ============

// QuotaConsumptionRecord 配额消费记录
type QuotaConsumptionRecord struct {
	Timestamp time.Time
	UserID    string
	Service   string
	Model     string
	Tokens    int64
}

// QuotaMetrics 配额使用统计
type QuotaMetrics struct {
	mu                   sync.RWMutex
	TotalConsumed        int64
	ByService            map[string]int64
	ByModel              map[string]int64
	ShortageCount        int64
	ConsumptionHistory   []QuotaConsumptionRecord
	maxHistoryRecords    int // 最大历史记录数
}

// QuotaReport 配额使用报告
type QuotaReport struct {
	TotalConsumed     int64                         `json:"total_consumed"`
	ByService         map[string]int64              `json:"by_service"`
	ByModel           map[string]int64              `json:"by_model"`
	ShortageCount     int64                         `json:"shortage_count"`
	RecentConsumption []QuotaConsumptionRecord      `json:"recent_consumption"`
	Statistics        *QuotaStatistics              `json:"statistics"`
}

// QuotaStatistics 配额统计信息
type QuotaStatistics struct {
	AveragePerService float64 `json:"average_per_service"`
	AveragePerModel   float64 `json:"average_per_model"`
	TopService        string  `json:"top_service"`
	TopModel          string  `json:"top_model"`
}

// NewQuotaMetrics 创建配额监控
func NewQuotaMetrics() *QuotaMetrics {
	return &QuotaMetrics{
		ByService:          make(map[string]int64),
		ByModel:            make(map[string]int64),
		ConsumptionHistory: make([]QuotaConsumptionRecord, 0, 100),
		maxHistoryRecords:  1000, // 默认保留最近1000条记录
	}
}

// RecordQuotaConsumed 记录配额消费
func (m *GRPCMetrics) RecordQuotaConsumed(userID, service, model string, tokens int64) {
	if m.quota == nil {
		m.quota = NewQuotaMetrics()
	}

	m.quota.mu.Lock()
	defer m.quota.mu.Unlock()

	// 更新总计
	m.quota.TotalConsumed += tokens

	// 按服务统计
	if _, exists := m.quota.ByService[service]; !exists {
		m.quota.ByService[service] = 0
	}
	m.quota.ByService[service] += tokens

	// 按模型统计
	if model != "" {
		if _, exists := m.quota.ByModel[model]; !exists {
			m.quota.ByModel[model] = 0
		}
		m.quota.ByModel[model] += tokens
	}

	// 添加历史记录
	record := QuotaConsumptionRecord{
		Timestamp: time.Now(),
		UserID:    userID,
		Service:   service,
		Model:     model,
		Tokens:    tokens,
	}

	m.quota.ConsumptionHistory = append(m.quota.ConsumptionHistory, record)

	// 限制历史记录数量
	if len(m.quota.ConsumptionHistory) > m.quota.maxHistoryRecords {
		// 保留最近的记录
		m.quota.ConsumptionHistory = m.quota.ConsumptionHistory[len(m.quota.ConsumptionHistory)-m.quota.maxHistoryRecords:]
	}
}

// RecordQuotaShortage 记录配额不足
func (m *GRPCMetrics) RecordQuotaShortage(userID string) {
	if m.quota == nil {
		m.quota = NewQuotaMetrics()
	}

	m.quota.mu.Lock()
	defer m.quota.mu.Unlock()

	m.quota.ShortageCount++
}

// GetQuotaReport 获取配额使用报告
func (m *GRPCMetrics) GetQuotaReport() *QuotaReport {
	if m.quota == nil {
		return &QuotaReport{
			TotalConsumed:     0,
			ByService:         make(map[string]int64),
			ByModel:           make(map[string]int64),
			ShortageCount:     0,
			RecentConsumption: []QuotaConsumptionRecord{},
			Statistics:        &QuotaStatistics{},
		}
	}

	m.quota.mu.RLock()
	defer m.quota.mu.RUnlock()

	// 复制数据避免外部修改
	byService := make(map[string]int64, len(m.quota.ByService))
	for k, v := range m.quota.ByService {
		byService[k] = v
	}

	byModel := make(map[string]int64, len(m.quota.ByModel))
	for k, v := range m.quota.ByModel {
		byModel[k] = v
	}

	// 获取最近20条记录
	recentCount := 20
	if len(m.quota.ConsumptionHistory) < recentCount {
		recentCount = len(m.quota.ConsumptionHistory)
	}
	recentConsumption := make([]QuotaConsumptionRecord, recentCount)
	copy(recentConsumption, m.quota.ConsumptionHistory[len(m.quota.ConsumptionHistory)-recentCount:])

	// 计算统计信息
	stats := m.calculateQuotaStatistics()

	return &QuotaReport{
		TotalConsumed:     m.quota.TotalConsumed,
		ByService:         byService,
		ByModel:           byModel,
		ShortageCount:     m.quota.ShortageCount,
		RecentConsumption: recentConsumption,
		Statistics:        stats,
	}
}

// calculateQuotaStatistics 计算配额统计信息
func (m *GRPCMetrics) calculateQuotaStatistics() *QuotaStatistics {
	if m.quota == nil {
		return &QuotaStatistics{}
	}

	stats := &QuotaStatistics{}

	// 计算平均每服务消耗
	serviceCount := len(m.quota.ByService)
	if serviceCount > 0 {
		stats.AveragePerService = float64(m.quota.TotalConsumed) / float64(serviceCount)
	}

	// 计算平均每模型消耗
	modelCount := len(m.quota.ByModel)
	if modelCount > 0 {
		stats.AveragePerModel = float64(m.quota.TotalConsumed) / float64(modelCount)
	}

	// 找出消耗最大的服务
	maxServiceTokens := int64(0)
	for service, tokens := range m.quota.ByService {
		if tokens > maxServiceTokens {
			maxServiceTokens = tokens
			stats.TopService = service
		}
	}

	// 找出消耗最大的模型
	maxModelTokens := int64(0)
	for model, tokens := range m.quota.ByModel {
		if tokens > maxModelTokens {
			maxModelTokens = tokens
			stats.TopModel = model
		}
	}

	return stats
}

// FormatQuotaReport 格式化配额报告
func (m *GRPCMetrics) FormatQuotaReport() string {
	report := m.GetQuotaReport()

	var result string
	result += "=== 配额使用统计报告 ===\n\n"

	// 总体统计
	result += fmt.Sprintf("总消耗: %d tokens\n", report.TotalConsumed)
	result += fmt.Sprintf("配额不足次数: %d\n\n", report.ShortageCount)

	// 按服务统计
	result += "## 按服务统计\n"
	for service, tokens := range report.ByService {
		result += fmt.Sprintf("  %s: %d tokens (%.2f%%)\n",
			service, tokens, float64(tokens)*100/float64(report.TotalConsumed))
	}
	result += "\n"

	// 按模型统计
	if len(report.ByModel) > 0 {
		result += "## 按模型统计\n"
		for model, tokens := range report.ByModel {
			result += fmt.Sprintf("  %s: %d tokens (%.2f%%)\n",
				model, tokens, float64(tokens)*100/float64(report.TotalConsumed))
		}
		result += "\n"
	}

	// 统计信息
	if report.Statistics != nil {
		result += "## 统计摘要\n"
		result += fmt.Sprintf("  平均每服务: %.2f tokens\n", report.Statistics.AveragePerService)
		result += fmt.Sprintf("  平均每模型: %.2f tokens\n", report.Statistics.AveragePerModel)
		if report.Statistics.TopService != "" {
			result += fmt.Sprintf("  最高服务: %s\n", report.Statistics.TopService)
		}
		if report.Statistics.TopModel != "" {
			result += fmt.Sprintf("  最高模型: %s\n", report.Statistics.TopModel)
		}
		result += "\n"
	}

	// 最近消费记录
	if len(report.RecentConsumption) > 0 {
		result += "## 最近消费记录\n"
		for _, record := range report.RecentConsumption {
			result += fmt.Sprintf("  [%s] %s - %s/%s: %d tokens\n",
				record.Timestamp.Format("15:04:05"),
				record.UserID,
				record.Service,
				record.Model,
				record.Tokens)
		}
	}

	return result
}

// ResetQuotaMetrics 重置配额统计
func (m *GRPCMetrics) ResetQuotaMetrics() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.quota != nil {
		m.quota.mu.Lock()
		defer m.quota.mu.Unlock()

		m.quota.TotalConsumed = 0
		m.quota.ByService = make(map[string]int64)
		m.quota.ByModel = make(map[string]int64)
		m.quota.ShortageCount = 0
		m.quota.ConsumptionHistory = make([]QuotaConsumptionRecord, 0, 100)
	}
}

// GetTotalQuotaConsumed 获取总配额消耗
func (m *GRPCMetrics) GetTotalQuotaConsumed() int64 {
	if m.quota == nil {
		return 0
	}

	m.quota.mu.RLock()
	defer m.quota.mu.RUnlock()
	return m.quota.TotalConsumed
}

// GetQuotaConsumptionByService 获取指定服务的配额消耗
func (m *GRPCMetrics) GetQuotaConsumptionByService(serviceName string) int64 {
	if m.quota == nil {
		return 0
	}

	m.quota.mu.RLock()
	defer m.quota.mu.RUnlock()

	return m.quota.ByService[serviceName]
}

// GetQuotaShortageCount 获取配额不足次数
func (m *GRPCMetrics) GetQuotaShortageCount() int64 {
	if m.quota == nil {
		return 0
	}

	m.quota.mu.RLock()
	defer m.quota.mu.RUnlock()
	return m.quota.ShortageCount
}
