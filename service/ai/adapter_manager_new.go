package ai

import (
	"context"
	"fmt"
	"sync"
	"time"

	"Qingyu_backend/repository/interfaces"
	"Qingyu_backend/service/base"
	serviceInterfaces "Qingyu_backend/service/interfaces"
)

// AdapterManagerNew 新的适配器管理器实现
type AdapterManagerNew struct {
	repositoryFactory interfaces.RepositoryFactory
	eventBus          base.EventBus
	validator         base.Validator
	serviceName       string
	version           string
	initialized       bool

	// 适配器管理
	adapters     map[string]serviceInterfaces.ModelAdapter
	adapterMutex sync.RWMutex

	// 模型配置
	modelConfigs map[string]*serviceInterfaces.ModelConfig
	configMutex  sync.RWMutex
}

// ModelAdapterImpl 模型适配器实现
type ModelAdapterImpl struct {
	modelID     string
	provider    string
	config      *serviceInterfaces.ModelConfig
	initialized bool
}

// NewAdapterManagerNew 创建新的适配器管理器
func NewAdapterManagerNew(
	repositoryFactory interfaces.RepositoryFactory,
	eventBus base.EventBus,
) serviceInterfaces.AdapterManager {
	service := &AdapterManagerNew{
		repositoryFactory: repositoryFactory,
		eventBus:          eventBus,
		validator:         base.NewBaseValidator(),
		adapters:          make(map[string]serviceInterfaces.ModelAdapter),
		modelConfigs:      make(map[string]*serviceInterfaces.ModelConfig),
		serviceName:       "AdapterManager",
		version:           "1.0.0",
		initialized:       false,
	}

	// 设置默认模型配置
	service.setupDefaultConfigs()

	return service
}

// setupDefaultConfigs 设置默认模型配置
func (s *AdapterManagerNew) setupDefaultConfigs() {
	s.configMutex.Lock()
	defer s.configMutex.Unlock()

	// OpenAI GPT-4
	s.modelConfigs["gpt-4"] = &serviceInterfaces.ModelConfig{
		ID:          "gpt-4",
		Name:        "GPT-4",
		Provider:    "openai",
		Type:        "text-generation",
		MaxTokens:   8192,
		Temperature: 0.7,
		TopP:        1.0,
		TopK:        0,
		InputPrice:  0.03,
		OutputPrice: 0.06,
		Features:    []string{"text-generation", "conversation", "code-generation"},
		Status:      "active",
		Description: "OpenAI GPT-4 模型",
		Metadata:    map[string]string{"version": "gpt-4-0613"},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Claude 3 Sonnet
	s.modelConfigs["claude-3-sonnet"] = &serviceInterfaces.ModelConfig{
		ID:          "claude-3-sonnet",
		Name:        "Claude 3 Sonnet",
		Provider:    "anthropic",
		Type:        "text-generation",
		MaxTokens:   200000,
		Temperature: 0.7,
		TopP:        1.0,
		TopK:        0,
		InputPrice:  0.003,
		OutputPrice: 0.015,
		Features:    []string{"text-generation", "conversation", "analysis"},
		Status:      "active",
		Description: "Anthropic Claude 3 Sonnet 模型",
		Metadata:    map[string]string{"version": "claude-3-sonnet-20240229"},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 阿里云通义千问
	s.modelConfigs["qwen-turbo"] = &serviceInterfaces.ModelConfig{
		ID:          "qwen-turbo",
		Name:        "通义千问 Turbo",
		Provider:    "alibaba",
		Type:        "text-generation",
		MaxTokens:   8192,
		Temperature: 0.7,
		TopP:        0.8,
		TopK:        50,
		InputPrice:  0.002,
		OutputPrice: 0.006,
		Features:    []string{"text-generation", "conversation", "chinese"},
		Status:      "active",
		Description: "阿里云通义千问 Turbo 模型",
		Metadata:    map[string]string{"version": "qwen-turbo-latest"},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// Initialize 初始化服务
func (s *AdapterManagerNew) Initialize(ctx context.Context) error {
	if s.initialized {
		return nil
	}

	// 检查Repository工厂健康状态
	if err := s.repositoryFactory.Health(ctx); err != nil {
		return base.NewServiceError(s.serviceName, base.ErrorTypeInternal, "Repository工厂不可用", err)
	}

	// 初始化默认适配器
	if err := s.initializeDefaultAdapters(ctx); err != nil {
		return base.NewServiceError(s.serviceName, base.ErrorTypeInternal, "初始化默认适配器失败", err)
	}

	s.initialized = true

	// 发布服务初始化事件
	event := &base.BaseEvent{
		EventType: "service.initialized",
		EventData: map[string]interface{}{
			"service": s.serviceName,
			"version": s.version,
		},
		Timestamp: time.Now(),
		Source:    s.serviceName,
	}

	if err := s.eventBus.PublishAsync(ctx, event); err != nil {
		fmt.Printf("发布服务初始化事件失败: %v\n", err)
	}

	return nil
}

// Health 健康检查
func (s *AdapterManagerNew) Health(ctx context.Context) error {
	if !s.initialized {
		return base.NewServiceError(s.serviceName, base.ErrorTypeInternal, "服务未初始化", nil)
	}

	// 检查Repository工厂健康状态
	if err := s.repositoryFactory.Health(ctx); err != nil {
		return base.NewServiceError(s.serviceName, base.ErrorTypeInternal, "Repository工厂健康检查失败", err)
	}

	// 检查适配器状态
	s.adapterMutex.RLock()
	adapterCount := len(s.adapters)
	s.adapterMutex.RUnlock()

	if adapterCount == 0 {
		return base.NewServiceError(s.serviceName, base.ErrorTypeInternal, "没有可用的模型适配器", nil)
	}

	return nil
}

// Close 关闭服务
func (s *AdapterManagerNew) Close(ctx context.Context) error {
	if !s.initialized {
		return nil
	}

	// 关闭所有适配器
	s.adapterMutex.Lock()
	for modelID, adapter := range s.adapters {
		if err := adapter.Close(ctx); err != nil {
			fmt.Printf("关闭适配器 %s 失败: %v\n", modelID, err)
		}
	}
	s.adapters = make(map[string]serviceInterfaces.ModelAdapter)
	s.adapterMutex.Unlock()

	s.initialized = false

	// 发布服务关闭事件
	event := &base.BaseEvent{
		EventType: "service.closed",
		EventData: map[string]interface{}{
			"service": s.serviceName,
		},
		Timestamp: time.Now(),
		Source:    s.serviceName,
	}

	if err := s.eventBus.PublishAsync(ctx, event); err != nil {
		fmt.Printf("发布服务关闭事件失败: %v\n", err)
	}

	return nil
}

// GetServiceName 获取服务名称
func (s *AdapterManagerNew) GetServiceName() string {
	return s.serviceName
}

// GetVersion 获取服务版本
func (s *AdapterManagerNew) GetVersion() string {
	return s.version
}

// GetAdapter 获取适配器
func (s *AdapterManagerNew) GetAdapter(ctx context.Context, req *serviceInterfaces.GetAdapterRequest) (*serviceInterfaces.GetAdapterResponse, error) {
	if req.AdapterID == "" {
		return nil, base.NewServiceError(s.serviceName, base.ErrorTypeValidation, "适配器ID不能为空", nil)
	}

	s.adapterMutex.RLock()
	defer s.adapterMutex.RUnlock()

	adapter, exists := s.adapters[req.AdapterID]
	if !exists {
		return nil, base.NewServiceError(s.serviceName, base.ErrorTypeNotFound, fmt.Sprintf("适配器 %s 不存在", req.AdapterID), nil)
	}

	// 获取模型配置
	s.configMutex.RLock()
	config, configExists := s.modelConfigs[req.AdapterID]
	s.configMutex.RUnlock()

	// 构建适配器信息
	adapterInfo := serviceInterfaces.AdapterInfo{
		ID:          req.AdapterID,
		Name:        adapter.GetModelID(),
		Type:        "model_adapter",
		Provider:    adapter.GetProvider(),
		Version:     "1.0.0",
		Status:      "active",
		Config:      map[string]string{},
		Features:    []string{"text-generation"},
		Description: fmt.Sprintf("适配器 %s", req.AdapterID),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 如果有配置信息，使用配置信息更新适配器信息
	if configExists {
		adapterInfo.Name = config.Name
		adapterInfo.Type = config.Type
		adapterInfo.Provider = config.Provider
		adapterInfo.Status = config.Status
		adapterInfo.Features = config.Features
		adapterInfo.Description = config.Description
		adapterInfo.CreatedAt = config.CreatedAt
		adapterInfo.UpdatedAt = config.UpdatedAt
		adapterInfo.Config = map[string]string{
			"provider": config.Provider,
			"type":     config.Type,
		}
	}

	return &serviceInterfaces.GetAdapterResponse{
		Adapter: adapterInfo,
	}, nil
}

// RegisterAdapter 注册适配器
func (s *AdapterManagerNew) RegisterAdapter(ctx context.Context, req *serviceInterfaces.RegisterAdapterRequest) (*serviceInterfaces.RegisterAdapterResponse, error) {
	if req.Name == "" || req.Type == "" || req.Provider == "" {
		return nil, base.NewServiceError(s.serviceName, base.ErrorTypeValidation, "名称、类型和提供商不能为空", nil)
	}

	// 生成适配器ID
	adapterID := fmt.Sprintf("%s_%s_%d", req.Provider, req.Type, time.Now().Unix())

	// 创建模型配置
	config := &serviceInterfaces.ModelConfig{
		ID:          adapterID,
		Name:        req.Name,
		Provider:    req.Provider,
		Type:        req.Type,
		Status:      "active",
		Description: req.Description,
		Features:    req.Features,
		Metadata:    req.Config,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 创建适配器实例
	adapter := &ModelAdapterImpl{
		modelID:  adapterID,
		provider: req.Provider,
		config:   config,
	}

	s.adapterMutex.Lock()
	s.adapters[adapterID] = adapter
	s.modelConfigs[adapterID] = config
	s.adapterMutex.Unlock()

	// 发布适配器注册事件
	event := &base.BaseEvent{
		EventType: "adapter.registered",
		EventData: map[string]interface{}{
			"adapter_id": adapterID,
			"provider":   req.Provider,
			"type":       req.Type,
		},
		Timestamp: time.Now(),
		Source:    s.serviceName,
	}

	if err := s.eventBus.Publish(ctx, event); err != nil {
		// 记录日志但不返回错误
	}

	return &serviceInterfaces.RegisterAdapterResponse{
		AdapterID:    adapterID,
		Registered:   true,
		RegisteredAt: time.Now(),
	}, nil
}

// UnregisterAdapter 注销适配器
func (s *AdapterManagerNew) UnregisterAdapter(ctx context.Context, req *serviceInterfaces.UnregisterAdapterRequest) (*serviceInterfaces.UnregisterAdapterResponse, error) {
	if req.AdapterID == "" {
		return nil, base.NewServiceError(s.serviceName, base.ErrorTypeValidation, "适配器ID不能为空", nil)
	}

	s.adapterMutex.Lock()
	adapter, exists := s.adapters[req.AdapterID]
	if exists {
		delete(s.adapters, req.AdapterID)
	}
	s.adapterMutex.Unlock()

	if !exists {
		return nil, base.NewServiceError(s.serviceName, base.ErrorTypeNotFound, fmt.Sprintf("未找到适配器: %s", req.AdapterID), nil)
	}

	// 关闭适配器
	if err := adapter.Close(ctx); err != nil {
		fmt.Printf("关闭适配器失败: %v\n", err)
	}

	// 删除模型配置
	s.adapterMutex.Lock()
	delete(s.modelConfigs, req.AdapterID)
	s.adapterMutex.Unlock()

	// 发布适配器注销事件
	event := &base.BaseEvent{
		EventType: "adapter.unregistered",
		EventData: map[string]interface{}{
			"adapter_id": req.AdapterID,
		},
		Timestamp: time.Now(),
		Source:    s.serviceName,
	}

	if err := s.eventBus.Publish(ctx, event); err != nil {
		// 记录日志但不返回错误
	}

	return &serviceInterfaces.UnregisterAdapterResponse{
		Unregistered:   true,
		UnregisteredAt: time.Now(),
	}, nil
}

// ListAdapters 列出所有适配器
func (s *AdapterManagerNew) ListAdapters(ctx context.Context, req *serviceInterfaces.ListAdaptersRequest) (*serviceInterfaces.ListAdaptersResponse, error) {
	s.adapterMutex.RLock()
	adapters := make([]serviceInterfaces.AdapterInfo, 0, len(s.adapters))

	for modelID, adapter := range s.adapters {
		s.configMutex.RLock()
		config, exists := s.modelConfigs[modelID]
		s.configMutex.RUnlock()

		adapterInfo := serviceInterfaces.AdapterInfo{
			ID:          modelID,
			Name:        "Default Adapter",
			Type:        "text",
			Provider:    "default",
			Version:     "1.0.0",
			Status:      "active",
			Config:      map[string]string{},
			Features:    []string{},
			Description: "Default adapter",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		if exists {
			adapterInfo.Name = config.Name
			adapterInfo.Type = config.Type
			adapterInfo.Provider = config.Provider
			adapterInfo.Status = config.Status
			adapterInfo.Features = config.Features
			adapterInfo.Description = config.Description
			adapterInfo.CreatedAt = config.CreatedAt
			adapterInfo.UpdatedAt = config.UpdatedAt
			adapterInfo.Config = map[string]string{
				"provider": config.Provider,
				"type":     config.Type,
			}
		}

		// 检查适配器健康状态
		if err := adapter.Health(ctx); err != nil {
			adapterInfo.Status = "inactive"
		}

		adapters = append(adapters, adapterInfo)
	}
	s.adapterMutex.RUnlock()

	response := &serviceInterfaces.ListAdaptersResponse{
		Adapters: adapters,
		Total:    len(adapters),
	}

	return response, nil
}

// UpdateAdapter 更新适配器
func (s *AdapterManagerNew) UpdateAdapter(ctx context.Context, req *serviceInterfaces.UpdateAdapterRequest) (*serviceInterfaces.UpdateAdapterResponse, error) {
	if req.AdapterID == "" {
		return nil, base.NewServiceError(s.serviceName, base.ErrorTypeValidation, "适配器ID不能为空", nil)
	}

	s.adapterMutex.Lock()
	defer s.adapterMutex.Unlock()

	adapter, exists := s.adapters[req.AdapterID]
	if !exists {
		return nil, base.NewServiceError(s.serviceName, base.ErrorTypeNotFound, "适配器不存在", nil)
	}

	// 更新适配器信息（这里是模拟实现）
	_ = adapter // 使用适配器进行更新操作

	// 发布适配器更新事件
	event := &base.BaseEvent{
		EventType: "adapter.updated",
		EventData: map[string]interface{}{
			"adapter_id": req.AdapterID,
			"updates":    req,
		},
		Timestamp: time.Now(),
		Source:    s.serviceName,
	}

	if err := s.eventBus.Publish(ctx, event); err != nil {
		// 记录日志但不返回错误
	}

	return &serviceInterfaces.UpdateAdapterResponse{
		Updated:   true,
		UpdatedAt: time.Now(),
	}, nil
}
func (s *AdapterManagerNew) GetModelConfig(ctx context.Context, req *serviceInterfaces.GetModelConfigRequest) (*serviceInterfaces.GetModelConfigResponse, error) {
	// 验证请求
	if req.ModelID == "" {
		return nil, base.NewServiceError(s.serviceName, base.ErrorTypeValidation, "模型ID不能为空", nil)
	}

	s.configMutex.RLock()
	config, exists := s.modelConfigs[req.ModelID]
	s.configMutex.RUnlock()

	if !exists {
		return nil, base.NewServiceError(s.serviceName, base.ErrorTypeNotFound, fmt.Sprintf("未找到模型配置: %s", req.ModelID), nil)
	}

	response := &serviceInterfaces.GetModelConfigResponse{
		Config: config,
	}

	return response, nil
}

// UpdateModelConfig 更新模型配置
func (s *AdapterManagerNew) UpdateModelConfig(ctx context.Context, req *serviceInterfaces.UpdateModelConfigRequest) (*serviceInterfaces.UpdateModelConfigResponse, error) {
	if req.ModelID == "" {
		return nil, base.NewServiceError(s.serviceName, base.ErrorTypeValidation, "模型ID不能为空", nil)
	}

	s.configMutex.Lock()
	defer s.configMutex.Unlock()

	config, exists := s.modelConfigs[req.ModelID]
	if !exists {
		return nil, base.NewServiceError(s.serviceName, base.ErrorTypeNotFound, "模型配置不存在", nil)
	}

	// 更新配置字段
	if req.Name != "" {
		config.Name = req.Name
	}
	if req.MaxTokens > 0 {
		config.MaxTokens = req.MaxTokens
	}
	if req.Temperature >= 0 {
		config.Temperature = req.Temperature
	}
	if req.TopP >= 0 {
		config.TopP = req.TopP
	}
	if req.TopK >= 0 {
		config.TopK = req.TopK
	}
	if req.Features != nil {
		config.Features = req.Features
	}
	if req.Status != "" {
		config.Status = req.Status
	}
	if req.Description != "" {
		config.Description = req.Description
	}
	if req.Metadata != nil {
		config.Metadata = req.Metadata
	}

	config.UpdatedAt = time.Now()
	s.modelConfigs[req.ModelID] = config

	return &serviceInterfaces.UpdateModelConfigResponse{
		Updated:   true,
		UpdatedAt: time.Now(),
	}, nil
}

// 辅助方法

// initializeDefaultAdapters 初始化默认适配器
func (s *AdapterManagerNew) initializeDefaultAdapters(ctx context.Context) error {
	s.configMutex.RLock()
	configs := make(map[string]*serviceInterfaces.ModelConfig)
	for k, v := range s.modelConfigs {
		configs[k] = v
	}
	s.configMutex.RUnlock()

	for modelID, config := range configs {
		adapter := &ModelAdapterImpl{
			modelID:     modelID,
			provider:    config.Provider,
			config:      config,
			initialized: true,
		}

		s.adapterMutex.Lock()
		s.adapters[modelID] = adapter
		s.adapterMutex.Unlock()
	}

	return nil
}

// validateRegisterAdapterRequest 验证注册适配器请求
func (s *AdapterManagerNew) validateRegisterAdapterRequest(req *serviceInterfaces.RegisterAdapterRequest) error {
	if req.Name == "" {
		return fmt.Errorf("适配器名称不能为空")
	}

	if req.Type == "" {
		return fmt.Errorf("适配器类型不能为空")
	}

	if req.Provider == "" {
		return fmt.Errorf("提供商不能为空")
	}

	if req.Version == "" {
		return fmt.Errorf("版本不能为空")
	}

	return nil
}

// validateUpdateModelConfigRequest 验证更新模型配置请求
func (s *AdapterManagerNew) validateUpdateModelConfigRequest(req *serviceInterfaces.UpdateModelConfigRequest) error {
	if req.ModelID == "" {
		return fmt.Errorf("模型ID不能为空")
	}

	return nil
}

// ModelAdapterImpl 方法实现

// Initialize 初始化适配器
func (a *ModelAdapterImpl) Initialize(ctx context.Context) error {
	if a.initialized {
		return nil
	}

	a.initialized = true
	return nil
}

// Health 健康检查
func (a *ModelAdapterImpl) Health(ctx context.Context) error {
	if !a.initialized {
		return fmt.Errorf("适配器未初始化")
	}

	return nil
}

// Close 关闭适配器
func (a *ModelAdapterImpl) Close(ctx context.Context) error {
	a.initialized = false
	return nil
}

// GetModelID 获取模型ID
func (a *ModelAdapterImpl) GetModelID() string {
	return a.modelID
}

// GetProvider 获取提供商
func (a *ModelAdapterImpl) GetProvider() string {
	return a.provider
}

// GetConfig 获取配置
func (a *ModelAdapterImpl) GetConfig() *serviceInterfaces.ModelConfig {
	return a.config
}

// GenerateContent 生成内容
func (a *ModelAdapterImpl) GenerateContent(ctx context.Context, req *serviceInterfaces.GenerateContentRequest) (*serviceInterfaces.GenerateContentResponse, error) {
	startTime := time.Now()

	// 这里应该调用实际的AI模型API
	// 为了简化，我们返回模拟响应
	return &serviceInterfaces.GenerateContentResponse{
		Content:      "这是一个模拟的生成内容响应",
		Model:        a.modelID,
		TokensUsed:   100,
		FinishReason: "stop",
		ResponseTime: time.Since(startTime),
		RequestID:    "mock-request-id",
	}, nil
}

// GenerateContentStream 流式生成内容
func (a *ModelAdapterImpl) GenerateContentStream(ctx context.Context, req *serviceInterfaces.GenerateContentRequest) (<-chan *serviceInterfaces.StreamResponse, error) {
	// 创建响应通道
	responseChan := make(chan *serviceInterfaces.StreamResponse, 10)

	go func() {
		defer close(responseChan)

		// 模拟流式响应
		content := "这是一个模拟的流式响应内容"
		words := []string{"这是", "一个", "模拟的", "流式", "响应", "内容"}

		for i, word := range words {
			select {
			case <-ctx.Done():
				return
			case responseChan <- &serviceInterfaces.StreamResponse{
				Content:    content[:len(word)*(i+1)],
				Delta:      word,
				Done:       i == len(words)-1,
				TokensUsed: i + 1,
			}:
				time.Sleep(100 * time.Millisecond) // 模拟延迟
			}
		}
	}()

	return responseChan, nil
}
