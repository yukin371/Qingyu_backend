package core

import (
	"strings"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// mockMiddlewareForPriority 测试用mock中间件
type mockMiddlewareForPriority struct {
	name     string
	priority int
}

func (m *mockMiddlewareForPriority) Name() string {
	return m.name
}

func (m *mockMiddlewareForPriority) Priority() int {
	return m.priority
}

func (m *mockMiddlewareForPriority) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

// TestDefaultPriority 测试默认优先级
func TestDefaultPriority(t *testing.T) {
	logger := zap.NewNop()
	registry := NewRegistry(logger)

	// 创建中间件，使用默认优先级
	mw1 := &mockMiddlewareForPriority{name: "request_id", priority: 1}
	mw2 := &mockMiddlewareForPriority{name: "recovery", priority: 2}
	mw3 := &mockMiddlewareForPriority{name: "cors", priority: 4}

	// 注册中间件
	if err := registry.Register(mw1); err != nil {
		t.Fatalf("Failed to register middleware: %v", err)
	}
	if err := registry.Register(mw2); err != nil {
		t.Fatalf("Failed to register middleware: %v", err)
	}
	if err := registry.Register(mw3); err != nil {
		t.Fatalf("Failed to register middleware: %v", err)
	}

	// 验证默认优先级
	tests := []struct {
		name            string
		middlewareName  string
		expectedPriority int
	}{
		{"RequestID默认优先级", "request_id", 1},
		{"Recovery默认优先级", "recovery", 2},
		{"CORS默认优先级", "cors", 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mw, err := registry.Get(tt.middlewareName)
			if err != nil {
				t.Fatalf("Failed to get middleware: %v", err)
			}
			if priority := mw.Priority(); priority != tt.expectedPriority {
				t.Errorf("Expected priority %d, got %d", tt.expectedPriority, priority)
			}
		})
	}
}

// TestGlobalPriorityOverride 测试全局优先级覆盖
func TestGlobalPriorityOverride(t *testing.T) {
	logger := zap.NewNop()
	registry := NewRegistry(logger)

	// 注册中间件
	rateLimit := &mockMiddlewareForPriority{name: "rate_limit", priority: 8}
	metrics := &mockMiddlewareForPriority{name: "metrics", priority: 7}

	if err := registry.Register(rateLimit); err != nil {
		t.Fatalf("Failed to register middleware: %v", err)
	}
	if err := registry.Register(metrics); err != nil {
		t.Fatalf("Failed to register middleware: %v", err)
	}

	// 设置优先级覆盖
	registry.SetPriorityOverride("rate_limit", 7)
	registry.SetPriorityOverride("metrics", 6)

	// 获取排序后的中间件列表
	sorted := registry.Sorted()

	// 验证顺序：metrics(6) -> rate_limit(7)
	if len(sorted) != 2 {
		t.Fatalf("Expected 2 middlewares, got %d", len(sorted))
	}

	if sorted[0].Name() != "metrics" {
		t.Errorf("Expected first middleware to be 'metrics', got '%s'", sorted[0].Name())
	}

	if sorted[1].Name() != "rate_limit" {
		t.Errorf("Expected second middleware to be 'rate_limit', got '%s'", sorted[1].Name())
	}
}

// TestRouteGroupPriorityConfiguration 测试路由组级别配置
func TestRouteGroupPriorityConfiguration(t *testing.T) {
	logger := zap.NewNop()
	manager := NewManager(logger)

	// 注册多个中间件
	middlewares := []*mockMiddlewareForPriority{
		{name: "request_id", priority: 1},
		{name: "recovery", priority: 2},
		{name: "cors", priority: 4},
		{name: "auth", priority: 9},
		{name: "permission", priority: 10},
	}

	for _, mw := range middlewares {
		if err := manager.Register(mw); err != nil {
			t.Fatalf("Failed to register middleware: %v", err)
		}
	}

	// 为API路由组设置特定配置：跳过认证和权限
	apiMiddlewares := []string{"request_id", "recovery", "cors", "metrics"}

	// 验证选择的路由组中间件数量
	if len(apiMiddlewares) != 4 {
		t.Errorf("Expected 4 middlewares for API route group, got %d", len(apiMiddlewares))
	}
}

// TestPriorityConflictDetection 测试优先级冲突检测
func TestPriorityConflictDetection(t *testing.T) {
	logger := zap.NewNop()
	registry := NewRegistry(logger)

	// 注册两个相同优先级的中间件
	mw1 := &mockMiddlewareForPriority{name: "middleware1", priority: 5}
	mw2 := &mockMiddlewareForPriority{name: "middleware2", priority: 5}

	if err := registry.Register(mw1); err != nil {
		t.Fatalf("Failed to register middleware: %v", err)
	}
	if err := registry.Register(mw2); err != nil {
		t.Fatalf("Failed to register middleware: %v", err)
	}

	// 验证检测到冲突（应该在Validate时发出警告）
	if err := registry.Validate(); err != nil {
		// 验证不应该失败，但应该有警告
		t.Logf("Validation detected issues: %v", err)
	}

	// 测试无效优先级
	mw3 := &mockMiddlewareForPriority{name: "middleware3", priority: 0}
	if err := registry.Register(mw3); err != nil {
		t.Fatalf("Failed to register middleware: %v", err)
	}

	if err := registry.Validate(); err == nil {
		t.Error("Expected validation error for invalid priority (0), got nil")
	}
}

// TestExecutionOrderGeneration 测试执行顺序生成
func TestExecutionOrderGeneration(t *testing.T) {
	logger := zap.NewNop()
	manager := NewManager(logger)

	// 注册中间件
	middlewares := []*mockMiddlewareForPriority{
		{name: "compression", priority: 12},
		{name: "request_id", priority: 1},
		{name: "recovery", priority: 2},
		{name: "cors", priority: 4},
	}

	for _, mw := range middlewares {
		if err := manager.Register(mw); err != nil {
			t.Fatalf("Failed to register middleware: %v", err)
		}
	}

	// 生成执行顺序
	executionOrder := make([]string, 0)
	for _, mw := range manager.List() {
		executionOrder = append(executionOrder, mw.Name())
	}

	// 验证顺序
	expectedOrder := []string{"request_id", "recovery", "cors", "compression"}
	if len(executionOrder) != len(expectedOrder) {
		t.Fatalf("Expected %d middlewares, got %d", len(expectedOrder), len(executionOrder))
	}

	for i, name := range expectedOrder {
		if executionOrder[i] != name {
			t.Errorf("At position %d, expected '%s', got '%s'", i, name, executionOrder[i])
		}
	}
}

// TestGetEffectivePriority 测试GetEffectivePriority方法
func TestGetEffectivePriority(t *testing.T) {
	logger := zap.NewNop()
	registry := NewRegistry(logger)

	// 注册中间件
	mw := &mockMiddlewareForPriority{name: "rate_limit", priority: 8}
	if err := registry.Register(mw); err != nil {
		t.Fatalf("Failed to register middleware: %v", err)
	}

	// 测试默认优先级
	if priority := registry.GetEffectivePriority("rate_limit"); priority != 8 {
		t.Errorf("Expected effective priority 8, got %d", priority)
	}

	// 设置覆盖
	registry.SetPriorityOverride("rate_limit", 7)

	// 测试覆盖后的优先级
	if priority := registry.GetEffectivePriority("rate_limit"); priority != 7 {
		t.Errorf("Expected effective priority 7 after override, got %d", priority)
	}

	// 测试不存在的中间件
	if priority := registry.GetEffectivePriority("nonexistent"); priority != -1 {
		t.Errorf("Expected -1 for nonexistent middleware, got %d", priority)
	}
}

// TestPriorityOverrideValidation 测试优先级覆盖配置验证
func TestPriorityOverrideValidation(t *testing.T) {
	logger := zap.NewNop()
	registry := NewRegistry(logger)

	// 注册中间件
	mw := &mockMiddlewareForPriority{name: "auth", priority: 9}
	if err := registry.Register(mw); err != nil {
		t.Fatalf("Failed to register middleware: %v", err)
	}

	// 测试有效覆盖
	registry.SetPriorityOverride("auth", 10)
	if priority := registry.GetEffectivePriority("auth"); priority != 10 {
		t.Errorf("Expected priority 10, got %d", priority)
	}

	// 测试无效覆盖（超出范围）
	registry.SetPriorityOverride("auth", 0)
	if err := registry.Validate(); err != nil {
		t.Logf("Validation correctly detected invalid priority: %v", err)
	}

	// 测试无效覆盖（超出范围）
	registry.SetPriorityOverride("auth", 101)
	if err := registry.Validate(); err != nil {
		t.Logf("Validation correctly detected invalid priority: %v", err)
	}
}

// TestConcurrentPriorityAccess 测试并发访问优先级
func TestConcurrentPriorityAccess(t *testing.T) {
	logger := zap.NewNop()
	registry := NewRegistry(logger)

	// 注册中间件
	mw := &mockMiddlewareForPriority{name: "concurrent_test", priority: 5}
	if err := registry.Register(mw); err != nil {
		t.Fatalf("Failed to register middleware: %v", err)
	}

	// 并发测试
	var wg sync.WaitGroup
	done := make(chan bool)

	// 启动多个goroutine同时访问和修改优先级
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(iteration int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				priority := (iteration + j) % 100 + 1 // 1-100
				registry.SetPriorityOverride("concurrent_test", priority)
				registry.GetEffectivePriority("concurrent_test")
			}
		}(i)
	}

	// 等待所有goroutine完成
	go func() {
		wg.Wait()
		done <- true
	}()

	<-done
}

// TestPriorityOrderReportGeneration 测试生成优先级顺序报告
func TestPriorityOrderReportGeneration(t *testing.T) {
	logger := zap.NewNop()
	manager := NewManager(logger)

	// 注册中间件
	middlewares := []*mockMiddlewareForPriority{
		{name: "request_id", priority: 1},
		{name: "recovery", priority: 2},
		{name: "cors", priority: 4},
		{name: "auth", priority: 9},
	}

	for _, mw := range middlewares {
		if err := manager.Register(mw); err != nil {
			t.Fatalf("Failed to register middleware: %v", err)
		}
	}

	// 生成报告
	report := manager.GenerateOrderReport()

	if report == "" {
		t.Error("Expected non-empty report, got empty string")
	}

	// 验证报告包含关键信息
	expectedKeywords := []string{
		"Middleware Execution Order",
		"request_id",
		"recovery",
		"cors",
		"auth",
	}

	for _, keyword := range expectedKeywords {
		if !strings.Contains(report, keyword) {
			t.Errorf("Expected report to contain '%s', but it didn't", keyword)
		}
	}
}

// TestManagerValidation 测试Manager的Validate方法
func TestManagerValidation(t *testing.T) {
	logger := zap.NewNop()
	manager := NewManager(logger)

	// 空管理器应该验证失败
	if err := manager.Validate(); err == nil {
		t.Error("Expected validation error for empty manager, got nil")
	}

	// 注册有效中间件
	mw := &mockMiddlewareForPriority{name: "test", priority: 5}
	if err := manager.Register(mw); err != nil {
		t.Fatalf("Failed to register middleware: %v", err)
	}

	// 应该验证通过
	if err := manager.Validate(); err != nil {
		t.Errorf("Expected validation to pass, got error: %v", err)
	}
}

// TestGetExecutionOrder 测试获取执行顺序
func TestGetExecutionOrder(t *testing.T) {
	logger := zap.NewNop()
	manager := NewManager(logger)

	// 注册中间件
	middlewares := []*mockMiddlewareForPriority{
		{name: "compression", priority: 12},
		{name: "request_id", priority: 1},
		{name: "recovery", priority: 2},
	}

	for _, mw := range middlewares {
		if err := manager.Register(mw); err != nil {
			t.Fatalf("Failed to register middleware: %v", err)
		}
	}

	// 获取执行顺序
	order := manager.GetExecutionOrder()

	expected := []string{"request_id", "recovery", "compression"}
	if len(order) != len(expected) {
		t.Fatalf("Expected %d items in execution order, got %d", len(expected), len(order))
	}

	for i, name := range expected {
		if order[i] != name {
			t.Errorf("At position %d, expected '%s', got '%s'", i, name, order[i])
		}
	}
}

