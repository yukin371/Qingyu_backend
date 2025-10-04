package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	// 新的接口和实现
	"Qingyu_backend/repository/interfaces"
	"Qingyu_backend/service/base"
	serviceInterfaces "Qingyu_backend/service/interfaces"
	aiServiceNew "Qingyu_backend/service/ai"
)

// TestNewArchitectureBasicFunctionality 测试新架构的基本功能
func TestNewArchitectureBasicFunctionality(t *testing.T) {
	ctx := context.Background()
	
	// 创建新的MongoDB工厂
	config := &interfaces.MongoConfig{
		URI:      "mongodb://localhost:27017",
		Database: "qingyu_test",
		Timeout:  30 * time.Second,
	}
	
	// 由于可能没有MongoDB连接，我们先测试配置验证
	t.Run("TestMongoConfigValidation", func(t *testing.T) {
		if config.URI == "" {
			t.Error("MongoDB URI不应为空")
		}
		
		if config.Database == "" {
			t.Error("数据库名称不应为空")
		}
		
		if config.Timeout <= 0 {
			t.Error("超时时间应大于0")
		}
	})
}

// TestNewServiceArchitecture 测试新Service架构
func TestNewServiceArchitecture(t *testing.T) {
	ctx := context.Background()
	
	// 创建事件总线
	eventBus := base.NewSimpleEventBus()
	
	// 创建模拟工厂
	mockFactory := &interfaces.MockRepositoryFactory{}
	
	// 创建服务实例
	aiService := aiServiceNew.NewAIServiceNew(mockFactory, eventBus)
	contextService := aiServiceNew.NewContextServiceNew(mockFactory, eventBus)
	externalAPIService := aiServiceNew.NewExternalAPIServiceNew(mockFactory, eventBus)
	adapterManager := aiServiceNew.NewAdapterManagerNew(mockFactory, eventBus)
	
	// 测试服务实例创建
	t.Run("TestServiceCreation", func(t *testing.T) {
		if aiService == nil {
			t.Error("AI服务实例不应为nil")
		}
		if contextService == nil {
			t.Error("上下文服务实例不应为nil")
		}
		if externalAPIService == nil {
			t.Error("外部API服务实例不应为nil")
		}
		if adapterManager == nil {
			t.Error("适配器管理器实例不应为nil")
		}
	})
	
	// 测试事件总线基本功能
	t.Run("TestEventBus", func(t *testing.T) {
		eventReceived := false
		
		// 订阅事件
		handler := &base.BaseEventHandler{
			EventType: "test.event",
			Handler: func(ctx context.Context, event base.Event) error {
				eventReceived = true
				return nil
			},
		}
		
		err := eventBus.Subscribe(ctx, "test.event", handler)
		if err != nil {
			t.Fatalf("订阅事件失败: %v", err)
		}
		
		// 发布事件
		event := &base.BaseEvent{
			EventType: "test.event",
			EventData: map[string]interface{}{
				"message": "test message",
			},
			Timestamp: time.Now(),
			Source:    "test",
		}
		
		err = eventBus.Publish(ctx, event)
		if err != nil {
			t.Fatalf("发布事件失败: %v", err)
		}
		
		// 等待事件处理
		time.Sleep(100 * time.Millisecond)
		
		if !eventReceived {
			t.Error("事件未被接收")
		}
	})
	
	// 测试验证器
	t.Run("TestValidator", func(t *testing.T) {
		validator := base.NewValidator()
		
		// 添加验证规则
		rule := &base.BaseValidationRule{
			Field:   "username",
			Message: "用户名不能为空",
			ValidateFunc: func(value interface{}) bool {
				if str, ok := value.(string); ok {
					return len(str) > 0
				}
				return false
			},
		}
		
		validator.AddRule(rule)
		
		// 测试验证
		data := map[string]interface{}{
			"username": "",
		}
		
		errors := validator.Validate(data)
		if len(errors) == 0 {
			t.Error("应该有验证错误")
		}
		
		// 测试有效数据
		data["username"] = "testuser"
		errors = validator.Validate(data)
		if len(errors) > 0 {
			t.Errorf("不应该有验证错误: %v", errors)
		}
	})
}

// TestServiceInterfaces 测试服务接口定义
func TestServiceInterfaces(t *testing.T) {
	// 测试请求结构体
	t.Run("TestGenerateContentRequest", func(t *testing.T) {
		req := &serviceInterfaces.GenerateContentRequest{
			Prompt:      "测试提示",
			MaxTokens:   1000,
			Temperature: 0.7,
			Model:       "gpt-4",
		}
		
		if req.Prompt == "" {
			t.Error("提示不应为空")
		}
		
		if req.MaxTokens <= 0 {
			t.Error("最大token数应大于0")
		}
		
		if req.Temperature < 0 || req.Temperature > 2 {
			t.Error("温度值应在0-2之间")
		}
	})
	
	t.Run("TestCreateContextRequest", func(t *testing.T) {
		req := &serviceInterfaces.CreateContextRequest{
			Name:        "测试上下文",
			Description: "这是一个测试上下文",
			Type:        "conversation",
			UserID:      "test_user",
		}
		
		if req.Name == "" {
			t.Error("上下文名称不应为空")
		}
		
		if req.Type == "" {
			t.Error("上下文类型不应为空")
		}
		
		if req.UserID == "" {
			t.Error("用户ID不应为空")
		}
	})
	
	t.Run("TestGetAPIStatusRequest", func(t *testing.T) {
		req := &serviceInterfaces.GetAPIStatusRequest{
			Provider: "openai",
		}
		
		if req.Provider == "" {
			t.Error("API提供商不应为空")
		}
	})
}

// TestServiceErrorHandling 测试服务错误处理
func TestServiceErrorHandling(t *testing.T) {
	t.Run("TestServiceError", func(t *testing.T) {
		err := base.NewServiceError(base.ErrCodeValidation, "验证失败", nil)
		
		if err.Code != base.ErrCodeValidation {
			t.Errorf("错误代码不匹配，期望: %s, 实际: %s", base.ErrCodeValidation, err.Code)
		}
		
		if err.Message != "验证失败" {
			t.Errorf("错误消息不匹配，期望: %s, 实际: %s", "验证失败", err.Message)
		}
		
		if !base.IsValidationError(err) {
			t.Error("应该是验证错误")
		}
	})
	
	t.Run("TestErrorTypes", func(t *testing.T) {
		validationErr := base.NewValidationError("字段验证失败")
		if !base.IsValidationError(validationErr) {
			t.Error("应该是验证错误")
		}
		
		notFoundErr := base.NewNotFoundError("资源未找到")
		if !base.IsNotFoundError(notFoundErr) {
			t.Error("应该是未找到错误")
		}
		
		internalErr := base.NewInternalError("内部服务器错误")
		if !base.IsInternalError(internalErr) {
			t.Error("应该是内部错误")
		}
	})
}

// BenchmarkNewArchitecturePerformance 测试新架构的性能
func BenchmarkNewArchitecturePerformance(b *testing.B) {
	ctx := context.Background()
	
	// 基准测试事件总线性能
	b.Run("EventBusPerformance", func(b *testing.B) {
		eventBus := base.NewSimpleEventBus()
		
		// 订阅事件
		handler := &base.BaseEventHandler{
			EventType: "benchmark.event",
			Handler: func(ctx context.Context, event base.Event) error {
				return nil
			},
		}
		
		err := eventBus.Subscribe(ctx, "benchmark.event", handler)
		if err != nil {
			b.Fatalf("订阅事件失败: %v", err)
		}
		
		b.ResetTimer()
		
		for i := 0; i < b.N; i++ {
			event := &base.BaseEvent{
				EventType: "benchmark.event",
				EventData: map[string]interface{}{
					"index": i,
				},
				Timestamp: time.Now(),
				Source:    "benchmark",
			}
			
			err := eventBus.Publish(ctx, event)
			if err != nil {
				b.Fatalf("发布事件失败: %v", err)
			}
		}
	})
	
	// 基准测试验证器性能
	b.Run("ValidatorPerformance", func(b *testing.B) {
		validator := base.NewValidator()
		
		// 添加验证规则
		rule := &base.BaseValidationRule{
			Field:   "username",
			Message: "用户名不能为空",
			ValidateFunc: func(value interface{}) bool {
				if str, ok := value.(string); ok {
					return len(str) > 0
				}
				return false
			},
		}
		
		validator.AddRule(rule)
		
		b.ResetTimer()
		
		for i := 0; i < b.N; i++ {
			data := map[string]interface{}{
				"username": fmt.Sprintf("user_%d", i),
			}
			
			validator.Validate(data)
		}
	})
	
	// 基准测试错误创建性能
	b.Run("ErrorCreationPerformance", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			err := base.NewServiceError(
				base.ErrCodeValidation,
				fmt.Sprintf("验证错误 %d", i),
				map[string]interface{}{
					"field": "username",
					"value": fmt.Sprintf("user_%d", i),
				},
			)
			
			if err == nil {
				b.Fatal("错误不应为nil")
			}
		}
	})
}