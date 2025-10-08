package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	// 新的接口和实现
	"Qingyu_backend/repository/interfaces"
	"Qingyu_backend/repository/mongodb"
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
	
	newFactory, err := mongodb.NewMongoRepositoryFactoryNew(config)
	if err != nil {
		t.Skipf("跳过MongoDB测试，连接失败: %v", err)
		return
	}
	defer newFactory.Close(ctx)
	
	// 测试Repository工厂健康检查
	t.Run("TestRepositoryFactoryHealth", func(t *testing.T) {
		err := newFactory.Health(ctx)
		if err != nil {
			t.Errorf("Repository工厂健康检查失败: %v", err)
		}
	})
	
	// 测试创建UserRepository
	t.Run("TestCreateUserRepository", func(t *testing.T) {
		userRepo, err := newFactory.CreateUserRepository(ctx)
		if err != nil {
			t.Fatalf("创建UserRepository失败: %v", err)
		}
		
		if userRepo == nil {
			t.Error("创建的UserRepository为nil")
		}
	})
	
	// 测试创建ProjectRepository
	t.Run("TestCreateProjectRepository", func(t *testing.T) {
		projectRepo, err := newFactory.CreateProjectRepository(ctx)
		if err != nil {
			t.Fatalf("创建ProjectRepository失败: %v", err)
		}
		
		if projectRepo == nil {
			t.Error("创建的ProjectRepository为nil")
		}
	})
	
	// 测试创建RoleRepository
	t.Run("TestCreateRoleRepository", func(t *testing.T) {
		roleRepo, err := newFactory.CreateRoleRepository(ctx)
		if err != nil {
			t.Fatalf("创建RoleRepository失败: %v", err)
		}
		
		if roleRepo == nil {
			t.Error("创建的RoleRepository为nil")
		}
	})
}

// TestNewServiceArchitecture 测试新Service架构
func TestNewServiceArchitecture(t *testing.T) {
	ctx := context.Background()
	
	// 创建新的Repository工厂
	config := &interfaces.MongoConfig{
		URI:      "mongodb://localhost:27017",
		Database: "qingyu_test",
		Timeout:  30 * time.Second,
	}
	
	repositoryFactory, err := mongodb.NewMongoRepositoryFactoryNew(config)
	if err != nil {
		t.Skipf("跳过Service测试，Repository工厂创建失败: %v", err)
		return
	}
	defer repositoryFactory.Close(ctx)
	
	// 创建事件总线
	eventBus := base.NewSimpleEventBus()
	
	// 测试AI服务
	t.Run("TestAIService", func(t *testing.T) {
		aiService := aiServiceNew.NewAIService(repositoryFactory, eventBus)
		
		// 测试服务初始化
		err := aiService.Initialize(ctx)
		if err != nil {
			t.Fatalf("初始化AI服务失败: %v", err)
		}
		
		// 测试服务健康检查
		err = aiService.Health(ctx)
		if err != nil {
			t.Errorf("AI服务健康检查失败: %v", err)
		}
		
		// 测试获取服务信息
		serviceName := aiService.GetServiceName()
		if serviceName == "" {
			t.Error("服务名称为空")
		}
		
		version := aiService.GetVersion()
		if version == "" {
			t.Error("服务版本为空")
		}
		
		// 测试生成内容功能
		req := &serviceInterfaces.GenerateContentRequest{
			Prompt:      "写一篇关于人工智能的文章",
			MaxTokens:   1000,
			Temperature: 0.7,
			Model:       "gpt-4",
		}
		
		resp, err := aiService.GenerateContent(ctx, req)
		if err != nil {
			t.Errorf("生成内容失败: %v", err)
		} else {
			if resp.Content == "" {
				t.Error("生成的内容为空")
			}
		}
		
		// 关闭服务
		err = aiService.Close(ctx)
		if err != nil {
			t.Errorf("关闭AI服务失败: %v", err)
		}
	})
	
	// 测试Context服务
	t.Run("TestContextService", func(t *testing.T) {
		contextService := aiServiceNew.NewContextService(repositoryFactory, eventBus)
		
		// 测试服务初始化
		err := contextService.Initialize(ctx)
		if err != nil {
			t.Fatalf("初始化Context服务失败: %v", err)
		}
		
		// 测试创建上下文
		createReq := &serviceInterfaces.CreateContextRequest{
			Name:        "测试上下文",
			Description: "这是一个测试上下文",
			Type:        "conversation",
			UserID:      "test_user",
		}
		
		createResp, err := contextService.CreateContext(ctx, createReq)
		if err != nil {
			t.Errorf("创建上下文失败: %v", err)
		} else {
			if createResp.ContextID == "" {
				t.Error("创建的上下文ID为空")
			}
		}
		
		// 关闭服务
		err = contextService.Close(ctx)
		if err != nil {
			t.Errorf("关闭Context服务失败: %v", err)
		}
	})
	
	// 测试ExternalAPI服务
	t.Run("TestExternalAPIService", func(t *testing.T) {
		externalAPIService := aiServiceNew.NewExternalAPIService(repositoryFactory, eventBus)
		
		// 测试服务初始化
		err := externalAPIService.Initialize(ctx)
		if err != nil {
			t.Fatalf("初始化ExternalAPI服务失败: %v", err)
		}
		
		// 测试获取API状态
		statusReq := &serviceInterfaces.GetAPIStatusRequest{
			Provider: "openai",
		}
		
		statusResp, err := externalAPIService.GetAPIStatus(ctx, statusReq)
		if err != nil {
			t.Errorf("获取API状态失败: %v", err)
		} else {
			if statusResp.Provider != "openai" {
				t.Errorf("API提供商不匹配，期望: openai, 实际: %s", statusResp.Provider)
			}
		}
		
		// 关闭服务
		err = externalAPIService.Close(ctx)
		if err != nil {
			t.Errorf("关闭ExternalAPI服务失败: %v", err)
		}
	})
	
	// 测试AdapterManager
	t.Run("TestAdapterManager", func(t *testing.T) {
		adapterManager := aiServiceNew.NewAdapterManager(repositoryFactory, eventBus)
		
		// 测试服务初始化
		err := adapterManager.Initialize(ctx)
		if err != nil {
			t.Fatalf("初始化AdapterManager失败: %v", err)
		}
		
		// 测试列出适配器
		listReq := &serviceInterfaces.ListAdaptersRequest{}
		
		listResp, err := adapterManager.ListAdapters(ctx, listReq)
		if err != nil {
			t.Errorf("列出适配器失败: %v", err)
		} else {
			if listResp.Total <= 0 {
				t.Error("应该有默认的适配器")
			}
		}
		
		// 测试获取模型配置
		configReq := &serviceInterfaces.GetModelConfigRequest{
			ModelID: "gpt-4",
		}
		
		configResp, err := adapterManager.GetModelConfig(ctx, configReq)
		if err != nil {
			t.Errorf("获取模型配置失败: %v", err)
		} else {
			if configResp.Config.ModelID != "gpt-4" {
				t.Errorf("模型ID不匹配，期望: gpt-4, 实际: %s", configResp.Config.ModelID)
			}
		}
		
		// 关闭服务
		err = adapterManager.Close(ctx)
		if err != nil {
			t.Errorf("关闭AdapterManager失败: %v", err)
		}
	})
}

// TestEventBusIntegration 测试事件总线集成
func TestEventBusIntegration(t *testing.T) {
	ctx := context.Background()
	
	// 创建事件总线
	eventBus := base.NewSimpleEventBus()
	
	// 测试事件发布和订阅
	t.Run("TestEventPublishSubscribe", func(t *testing.T) {
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
}

// BenchmarkNewArchitecturePerformance 测试新架构的性能
func BenchmarkNewArchitecturePerformance(b *testing.B) {
	ctx := context.Background()
	
	// 创建新的Repository工厂
	config := &interfaces.MongoConfig{
		URI:      "mongodb://localhost:27017",
		Database: "qingyu_benchmark",
		Timeout:  30 * time.Second,
	}
	
	repositoryFactory, err := mongodb.NewMongoRepositoryFactoryNew(config)
	if err != nil {
		b.Skipf("跳过性能测试，Repository工厂创建失败: %v", err)
		return
	}
	defer repositoryFactory.Close(ctx)
	
	// 创建事件总线
	eventBus := base.NewSimpleEventBus()
	
	// 基准测试AI服务初始化
	b.Run("AIServiceInitialization", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			aiService := aiServiceNew.NewAIService(repositoryFactory, eventBus)
			err := aiService.Initialize(ctx)
			if err != nil {
				b.Fatalf("初始化AI服务失败: %v", err)
			}
			aiService.Close(ctx)
		}
	})
	
	// 基准测试Context服务操作
	b.Run("ContextServiceOperations", func(b *testing.B) {
		contextService := aiServiceNew.NewContextService(repositoryFactory, eventBus)
		err := contextService.Initialize(ctx)
		if err != nil {
			b.Fatalf("初始化Context服务失败: %v", err)
		}
		defer contextService.Close(ctx)
		
		b.ResetTimer()
		
		for i := 0; i < b.N; i++ {
			req := &serviceInterfaces.CreateContextRequest{
				Name:        fmt.Sprintf("benchmark_context_%d", i),
				Description: "基准测试上下文",
				Type:        "conversation",
				UserID:      "benchmark_user",
			}
			
			_, err := contextService.CreateContext(ctx, req)
			if err != nil {
				b.Fatalf("创建上下文失败: %v", err)
			}
		}
	})
}