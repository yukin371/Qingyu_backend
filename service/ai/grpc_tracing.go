package ai

import (
	"fmt"
	"sync"
	"time"
)

// Tracer 请求追踪器
type Tracer struct {
	mu         sync.RWMutex
	traces     map[string]*RequestTrace
	activeTraces map[string]*ActiveTrace
	maxSize    int
	maxAge     time.Duration
}

// RequestTrace 请求追踪信息
type RequestTrace struct {
	RequestID   string
	ServiceName string
	StartTime   time.Time
	EndTime     time.Time
	Duration    time.Duration
	Status      string
	Error       string
	Metadata    map[string]string
}

// ActiveTrace 正在进行的追踪
type ActiveTrace struct {
	RequestID   string
	ServiceName string
	StartTime   time.Time
	Metadata    map[string]string
}

// TraceStatus 追踪状态常量
const (
	TraceStatusSuccess   = "success"
	TraceStatusFailed    = "failed"
	TraceStatusTimeout   = "timeout"
	TraceStatusCancelled = "cancelled"
	TraceStatusRunning   = "running"
)

// NewTracer 创建新的追踪器
func NewTracer(maxSize int) *Tracer {
	if maxSize <= 0 {
		maxSize = 1000 // 默认保留1000条追踪记录
	}

	return &Tracer{
		traces:       make(map[string]*RequestTrace),
		activeTraces: make(map[string]*ActiveTrace),
		maxSize:      maxSize,
		maxAge:       24 * time.Hour, // 默认保留24小时
	}
}

// StartTrace 开始追踪
func (t *Tracer) StartTrace(serviceName, requestID string) {
	t.StartTraceWithMetadata(serviceName, requestID, nil)
}

// StartTraceWithMetadata 开始追踪（带元数据）
func (t *Tracer) StartTraceWithMetadata(serviceName, requestID string, metadata map[string]string) {
	if requestID == "" {
		requestID = generateRequestID()
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	// 创建活动追踪
	activeTrace := &ActiveTrace{
		RequestID:   requestID,
		ServiceName: serviceName,
		StartTime:   time.Now(),
		Metadata:    metadata,
	}

	t.activeTraces[requestID] = activeTrace

	// 清理过期的追踪记录
	t.cleanupOldTraces()
}

// EndTrace 结束追踪
func (t *Tracer) EndTrace(requestID string, status string, err error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	// 获取活动追踪
	activeTrace, exists := t.activeTraces[requestID]
	if !exists {
		return // 没有对应的追踪记录
	}

	// 计算持续时间
	endTime := time.Now()
	duration := endTime.Sub(activeTrace.StartTime)

	// 创建完成追踪
	trace := &RequestTrace{
		RequestID:   requestID,
		ServiceName: activeTrace.ServiceName,
		StartTime:   activeTrace.StartTime,
		EndTime:     endTime,
		Duration:    duration,
		Status:      status,
		Metadata:    activeTrace.Metadata,
	}

	// 记录错误信息
	if err != nil {
		trace.Error = err.Error()
	}

	// 删除活动追踪
	delete(t.activeTraces, requestID)

	// 检查是否超过最大容量
	if len(t.traces) >= t.maxSize {
		// 删除最旧的追踪记录
		t.removeOldestTrace()
	}

	// 保存追踪记录
	t.traces[requestID] = trace
}

// GetTrace 获取追踪信息
func (t *Tracer) GetTrace(requestID string) (*RequestTrace, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	// 先查找活动追踪
	if activeTrace, exists := t.activeTraces[requestID]; exists {
		return &RequestTrace{
			RequestID:   activeTrace.RequestID,
			ServiceName: activeTrace.ServiceName,
			StartTime:   activeTrace.StartTime,
			Status:      TraceStatusRunning,
			Metadata:    activeTrace.Metadata,
		}, nil
	}

	// 查找已完成追踪
	if trace, exists := t.traces[requestID]; exists {
		// 返回副本避免外部修改
		return &RequestTrace{
			RequestID:   trace.RequestID,
			ServiceName: trace.ServiceName,
			StartTime:   trace.StartTime,
			EndTime:     trace.EndTime,
			Duration:    trace.Duration,
			Status:      trace.Status,
			Error:       trace.Error,
			Metadata:    trace.Metadata,
		}, nil
	}

	return nil, fmt.Errorf("trace %s not found", requestID)
}

// GetRecentTraces 获取最近的追踪记录
func (t *Tracer) GetRecentTraces(limit int) []*RequestTrace {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if limit <= 0 {
		limit = 10 // 默认返回10条
	}

	// 收集所有追踪记录
	allTraces := make([]*RequestTrace, 0, len(t.traces))
	for _, trace := range t.traces {
		allTraces = append(allTraces, trace)
	}

	// 按开始时间排序（最新的在前）
	sortTracesByStartTime(allTraces)

	// 返回最近的N条
	if len(allTraces) > limit {
		allTraces = allTraces[:limit]
	}

	return allTraces
}

// GetTracesByService 获取指定服务的追踪记录
func (t *Tracer) GetTracesByService(serviceName string, limit int) []*RequestTrace {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if limit <= 0 {
		limit = 10
	}

	var result []*RequestTrace
	for _, trace := range t.traces {
		if trace.ServiceName == serviceName {
			result = append(result, trace)
		}
	}

	// 按开始时间排序
	sortTracesByStartTime(result)

	if len(result) > limit {
		result = result[:limit]
	}

	return result
}

// GetTracesByStatus 获取指定状态的追踪记录
func (t *Tracer) GetTracesByStatus(status string, limit int) []*RequestTrace {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if limit <= 0 {
		limit = 10
	}

	var result []*RequestTrace
	for _, trace := range t.traces {
		if trace.Status == status {
			result = append(result, trace)
		}
	}

	// 按开始时间排序
	sortTracesByStartTime(result)

	if len(result) > limit {
		result = result[:limit]
	}

	return result
}

// GetActiveTraces 获取所有活动的追踪
func (t *Tracer) GetActiveTraces() []*RequestTrace {
	t.mu.RLock()
	defer t.mu.RUnlock()

	result := make([]*RequestTrace, 0, len(t.activeTraces))
	for _, active := range t.activeTraces {
		result = append(result, &RequestTrace{
			RequestID:   active.RequestID,
			ServiceName: active.ServiceName,
			StartTime:   active.StartTime,
			Status:      TraceStatusRunning,
			Metadata:    active.Metadata,
		})
	}

	return result
}

// GetStats 获取统计信息
func (t *Tracer) GetStats() *TraceStats {
	t.mu.RLock()
	defer t.mu.RUnlock()

	stats := &TraceStats{
		TotalTraces:   len(t.traces),
		ActiveTraces:  len(t.activeTraces),
		MaxSize:       t.maxSize,
		ServiceCounts: make(map[string]int),
		StatusCounts:  make(map[string]int),
	}

	for _, trace := range t.traces {
		stats.ServiceCounts[trace.ServiceName]++
		stats.StatusCounts[trace.Status]++
	}

	return stats
}

// cleanupOldTraces 清理过期的追踪记录
func (t *Tracer) cleanupOldTraces() {
	now := time.Now()
	for requestID, trace := range t.traces {
		if now.Sub(trace.EndTime) > t.maxAge {
			delete(t.traces, requestID)
		}
	}
}

// removeOldestTrace 删除最旧的追踪记录
func (t *Tracer) removeOldestTrace() {
	var oldestID string
	var oldestTime time.Time

	for id, trace := range t.traces {
		if oldestID == "" || trace.EndTime.Before(oldestTime) {
			oldestID = id
			oldestTime = trace.EndTime
		}
	}

	if oldestID != "" {
		delete(t.traces, oldestID)
	}
}

// sortTracesByStartTime 按开始时间排序追踪记录（最新的在前）
func sortTracesByStartTime(traces []*RequestTrace) {
	for i := 0; i < len(traces)-1; i++ {
		for j := i + 1; j < len(traces); j++ {
			if traces[i].StartTime.Before(traces[j].StartTime) {
				traces[i], traces[j] = traces[j], traces[i]
			}
		}
	}
}

// TraceStats 追踪统计信息
type TraceStats struct {
	TotalTraces  int            // 总追踪数
	ActiveTraces int            // 活动追踪数
	MaxSize      int            // 最大容量
	ServiceCounts map[string]int // 服务调用次数
	StatusCounts map[string]int // 状态分布
}

// generateRequestID 生成请求ID
func generateRequestID() string {
	return fmt.Sprintf("req-%d", time.Now().UnixNano())
}

// FormatTrace 格式化追踪信息为字符串
func (t *RequestTrace) FormatTrace() string {
	var result string
	result += fmt.Sprintf("RequestID: %s\n", t.RequestID)
	result += fmt.Sprintf("Service: %s\n", t.ServiceName)
	result += fmt.Sprintf("StartTime: %s\n", t.StartTime.Format(time.RFC3339))
	result += fmt.Sprintf("EndTime: %s\n", t.EndTime.Format(time.RFC3339))
	result += fmt.Sprintf("Duration: %v\n", t.Duration)
	result += fmt.Sprintf("Status: %s\n", t.Status)
	if t.Error != "" {
		result += fmt.Sprintf("Error: %s\n", t.Error)
	}
	if len(t.Metadata) > 0 {
		result += "Metadata:\n"
		for k, v := range t.Metadata {
			result += fmt.Sprintf("  %s: %s\n", k, v)
		}
	}
	return result
}

// Clear 清空所有追踪记录
func (t *Tracer) Clear() {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.traces = make(map[string]*RequestTrace)
	t.activeTraces = make(map[string]*ActiveTrace)
}

// SetMaxAge 设置最大保留时间
func (t *Tracer) SetMaxAge(maxAge time.Duration) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.maxAge = maxAge
}
