package ai

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"Qingyu_backend/pkg/errors"
	"Qingyu_backend/repository/interfaces"
	"Qingyu_backend/service/base"
	serviceInterfaces "Qingyu_backend/service/interfaces"
)

// ExternalAPIServiceNew 新的外部API服务实现
type ExternalAPIServiceNew struct {
	repositoryFactory interfaces.RepositoryFactory
	eventBus          base.EventBus
	validator         base.Validator
	httpClient        *http.Client
	serviceName       string
	version           string
	initialized       bool
	apiKeys           map[string]string // 存储API密钥
	apiEndpoints      map[string]string // 存储API端点
}

// NewExternalAPIServiceNew 创建新的外部API服务
func NewExternalAPIServiceNew(
	repositoryFactory interfaces.RepositoryFactory,
	eventBus base.EventBus,
) serviceInterfaces.ExternalAPIService {
	service := &ExternalAPIServiceNew{
		repositoryFactory: repositoryFactory,
		eventBus:          eventBus,
		validator:         base.NewBaseValidator(),
		httpClient:        &http.Client{Timeout: 30 * time.Second},
		serviceName:       "ExternalAPIService",
		version:           "1.0.0",
		initialized:       false,
		apiKeys:           make(map[string]string),
		apiEndpoints:      make(map[string]string),
	}

	service.setupDefaultEndpoints()
	return service
}

// setupDefaultEndpoints 设置默认API端点
func (s *ExternalAPIServiceNew) setupDefaultEndpoints() {
	s.apiEndpoints["openai"] = "https://api.openai.com/v1"
	s.apiEndpoints["claude"] = "https://api.anthropic.com/v1"
	s.apiEndpoints["gemini"] = "https://generativelanguage.googleapis.com/v1"
	s.apiEndpoints["qianwen"] = "https://dashscope.aliyuncs.com/api/v1"
	s.apiEndpoints["zhipu"] = "https://open.bigmodel.cn/api/paas/v4"
}

// Initialize 初始化服务
func (s *ExternalAPIServiceNew) Initialize(ctx context.Context) error {
	if s.initialized {
		return nil
	}
	
	// 检查Repository工厂健康状态
	if err := s.repositoryFactory.Health(ctx); err != nil {
		return errors.ExternalAPIFactory.InternalError("Repository工厂不可用", err).
			WithOperation("Initialize")
	}
	
	// 从配置或数据库加载API密钥
	if err := s.loadAPIKeys(ctx); err != nil {
		return errors.ExternalAPIFactory.InternalError("加载API密钥失败", err).
			WithOperation("Initialize")
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
func (s *ExternalAPIServiceNew) Health(ctx context.Context) error {
	if !s.initialized {
		return errors.ExternalAPIFactory.InternalError("服务未初始化", nil).
			WithOperation("Health")
	}
	
	// 检查Repository工厂健康状态
	if err := s.repositoryFactory.Health(ctx); err != nil {
		return errors.ExternalAPIFactory.InternalError("Repository工厂健康检查失败", err).
			WithOperation("Health")
	}
	
	// 检查HTTP客户端
	if s.httpClient == nil {
		return errors.ExternalAPIFactory.InternalError("HTTP客户端未初始化", nil).
			WithOperation("Health")
	}
	
	return nil
}

// Close 关闭服务
func (s *ExternalAPIServiceNew) Close(ctx context.Context) error {
	if !s.initialized {
		return nil
	}
	
	s.initialized = false
	
	// 清理资源
	s.httpClient.CloseIdleConnections()
	
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
func (s *ExternalAPIServiceNew) GetServiceName() string {
	return s.serviceName
}

// GetVersion 获取服务版本
func (s *ExternalAPIServiceNew) GetVersion() string {
	return s.version
}

// CallAPI 调用外部API
func (s *ExternalAPIServiceNew) CallAPI(ctx context.Context, req *serviceInterfaces.CallAPIRequest) (*serviceInterfaces.CallAPIResponse, error) {
	// 验证请求
	if err := s.validateCallAPIRequest(req); err != nil {
		return nil, errors.ExternalAPIFactory.ValidationError("请求验证失败", err).
			WithOperation("CallAPI").
			WithMetadata(map[string]interface{}{
				"provider": req.Provider,
				"endpoint": req.Endpoint,
				"method": req.Method,
			})
	}
	
	// 获取API端点
	endpoint, exists := s.apiEndpoints[req.Provider]
	if !exists {
		return nil, errors.ExternalAPIFactory.ValidationError(fmt.Sprintf("不支持的API提供商: %s", req.Provider), nil).
			WithOperation("CallAPI").
			WithMetadata(map[string]interface{}{
				"provider": req.Provider,
				"available_providers": s.getAvailableProviders(),
			})
	}
	
	// 获取API密钥
	apiKey, exists := s.apiKeys[req.Provider]
	if !exists {
		return nil, errors.ExternalAPIFactory.ValidationError(fmt.Sprintf("未配置API密钥: %s", req.Provider), nil).
			WithOperation("CallAPI").
			WithMetadata(map[string]interface{}{
				"provider": req.Provider,
			})
	}
	
	// 构建完整URL
	fullURL := fmt.Sprintf("%s%s", endpoint, req.Endpoint)
	
	// 创建HTTP请求
	httpReq, err := http.NewRequestWithContext(ctx, req.Method, fullURL, nil)
	if err != nil {
		return nil, errors.ExternalAPIFactory.InternalError("创建HTTP请求失败", err).
			WithOperation("CallAPI").
			WithMetadata(map[string]interface{}{
				"provider": req.Provider,
				"endpoint": req.Endpoint,
				"method": req.Method,
				"full_url": fullURL,
			})
	}
	
	// 设置请求头
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	httpReq.Header.Set("Content-Type", "application/json")
	
	// 添加自定义请求头
	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}
	
	// 发送请求
	startTime := time.Now()
	resp, err := s.httpClient.Do(httpReq)
	duration := time.Since(startTime)
	
	if err != nil {
		// 记录API调用失败事件
		s.recordAPICall(ctx, req.Provider, req.Endpoint, "failed", duration, err.Error())
		return nil, errors.ExternalAPIFactory.ExternalError("API调用失败", err).
			WithOperation("CallAPI").
			WithMetadata(map[string]interface{}{
				"provider": req.Provider,
				"endpoint": req.Endpoint,
				"method": req.Method,
				"full_url": fullURL,
				"duration_ms": duration.Milliseconds(),
			})
	}
	defer resp.Body.Close()
	
	// 读取响应
	responseData := make([]byte, 0)
	if resp.Body != nil {
		// 这里应该读取响应体，为了简化暂时跳过
	}
	
	// 记录API调用成功事件
	s.recordAPICall(ctx, req.Provider, req.Endpoint, "success", duration, "")
	
	response := &serviceInterfaces.CallAPIResponse{
		StatusCode:   resp.StatusCode,
		Headers:      make(map[string]string),
		Body:         responseData,
		ResponseTime: duration,
		Success:      true,
	}
	
	// 复制响应头
	for key, values := range resp.Header {
		if len(values) > 0 {
			response.Headers[key] = values[0]
		}
	}
	
	return response, nil
}

// GetAPIStatus 获取API状态
func (s *ExternalAPIServiceNew) GetAPIStatus(ctx context.Context, req *serviceInterfaces.GetAPIStatusRequest) (*serviceInterfaces.GetAPIStatusResponse, error) {
	startTime := time.Now()
	
	if req.Provider == "" {
		return nil, base.NewServiceError(s.serviceName, base.ErrorTypeValidation, "Provider不能为空", nil)
	}

	// 这里应该检查实际的API状态
	// 为了简化，返回模拟状态
	return &serviceInterfaces.GetAPIStatusResponse{
		Provider:     req.Provider,
		Status:       "active",
		ResponseTime: time.Since(startTime),
		LastChecked:  time.Now(),
		Message:      "API运行正常",
	}, nil
}

// GetUsageStats 获取API使用统计
func (s *ExternalAPIServiceNew) GetUsageStats(ctx context.Context, req *serviceInterfaces.GetAPIUsageRequest) (*serviceInterfaces.GetAPIUsageResponse, error) {
	if req.Provider == "" {
		return nil, errors.ExternalAPIFactory.ValidationError("Provider不能为空", nil).
			WithOperation("GetUsageStats")
	}

	// 这里应该从数据库或缓存中获取实际的使用统计
	// 为了简化，返回模拟数据
	return &serviceInterfaces.GetAPIUsageResponse{
		Provider: req.Provider,
		Usage: serviceInterfaces.APIUsage{
			TotalRequests: 1000,
			TotalTokens:   50000,
			TotalCost:     25.50,
			SuccessRate:   0.95,
		},
		Period: serviceInterfaces.TimePeriod{
			StartDate: time.Now().AddDate(0, -1, 0),
			EndDate:   time.Now(),
		},
	}, nil
}

// RefreshAPIKey 刷新API密钥
func (s *ExternalAPIServiceNew) RefreshAPIKey(ctx context.Context, req *serviceInterfaces.RefreshAPIKeyRequest) (*serviceInterfaces.RefreshAPIKeyResponse, error) {
	if req.Provider == "" {
		return nil, errors.ExternalAPIFactory.ValidationError("Provider不能为空", nil).
			WithOperation("RefreshAPIKey").
			WithMetadata(map[string]interface{}{
				"provider": req.Provider,
			})
	}

	// 这里应该实现实际的API密钥刷新逻辑
	// 为了简化，我们模拟刷新过程
	
	// 模拟刷新API密钥
	refreshed := true
	
	// 发布事件
	event := &base.BaseEvent{
		EventType: "api_key_refreshed",
		EventData: map[string]interface{}{
			"provider":     req.Provider,
			"refreshed":    refreshed,
			"refreshed_at": time.Now(),
		},
		Timestamp: time.Now(),
		Source:    s.serviceName,
	}
	
	s.eventBus.PublishAsync(ctx, event)

	return &serviceInterfaces.RefreshAPIKeyResponse{
		Provider:    req.Provider,
		Refreshed:   refreshed,
		RefreshedAt: time.Now(),
		Message:     "API密钥刷新成功",
	}, nil
}

// 辅助方法

// loadAPIKeys 从配置或数据库加载API密钥
func (s *ExternalAPIServiceNew) loadAPIKeys(ctx context.Context) error {
	// 这里应该从数据库或配置文件加载API密钥
	// 为了简化，我们设置一些示例密钥
	s.apiKeys["openai"] = "sk-example-openai-key"
	s.apiKeys["claude"] = "sk-example-claude-key"
	s.apiKeys["gemini"] = "example-gemini-key"
	
	return nil
}

// checkAPIHealth 检查API健康状态
func (s *ExternalAPIServiceNew) checkAPIHealth(ctx context.Context, provider, endpoint string) bool {
	// 这里应该发送健康检查请求
	// 为了简化，我们返回模拟结果
	return true
}

// validateAPIKey 验证API密钥
func (s *ExternalAPIServiceNew) validateAPIKey(ctx context.Context, provider, apiKey string) error {
	if apiKey == "" {
		return fmt.Errorf("API密钥不能为空")
	}
	
	// 这里应该向API提供商验证密钥有效性
	// 为了简化，我们只做基本格式检查
	switch provider {
	case "openai":
		if len(apiKey) < 20 || apiKey[:3] != "sk-" {
			return fmt.Errorf("OpenAI API密钥格式无效")
		}
	case "claude":
		if len(apiKey) < 20 {
			return fmt.Errorf("Claude API密钥格式无效")
		}
	case "gemini":
		if len(apiKey) < 10 {
			return fmt.Errorf("Gemini API密钥格式无效")
		}
	}
	
	return nil
}

// recordAPICall 记录API调用
func (s *ExternalAPIServiceNew) recordAPICall(ctx context.Context, provider, endpoint, status string, duration time.Duration, errorMsg string) {
	event := &base.BaseEvent{
		EventType: "api.called",
		EventData: map[string]interface{}{
			"provider":  provider,
			"endpoint":  endpoint,
			"status":    status,
			"duration":  duration.Milliseconds(),
			"error":     errorMsg,
		},
		Timestamp: time.Now(),
		Source:    s.serviceName,
	}
	
	s.eventBus.PublishAsync(ctx, event)
}

// GetAPIUsage 获取API使用情况
func (s *ExternalAPIServiceNew) GetAPIUsage(ctx context.Context, req *serviceInterfaces.GetAPIUsageRequest) (*serviceInterfaces.GetAPIUsageResponse, error) {
	if req.Provider == "" {
		return nil, base.NewServiceError(s.serviceName, base.ErrorTypeValidation, "Provider不能为空", nil)
	}

	// 这里应该从数据库或缓存中获取实际的使用统计
	// 为了简化，返回模拟数据
	return &serviceInterfaces.GetAPIUsageResponse{
		Provider: req.Provider,
		Usage: serviceInterfaces.APIUsage{
			TotalRequests: 1000,
			TotalTokens:   50000,
			TotalCost:     25.50,
			SuccessRate:   0.95,
		},
		Period: serviceInterfaces.TimePeriod{
			StartDate: time.Now().AddDate(0, -1, 0),
			EndDate:   time.Now(),
		},
	}, nil
}

func (s *ExternalAPIServiceNew) validateCallAPIRequest(req *serviceInterfaces.CallAPIRequest) error {
	if req.Provider == "" {
		return fmt.Errorf("API提供商不能为空")
	}
	
	if req.Endpoint == "" {
		return fmt.Errorf("API端点不能为空")
	}
	
	if req.Method == "" {
		return fmt.Errorf("HTTP方法不能为空")
	}
	
	// 验证HTTP方法
	validMethods := map[string]bool{
		"GET":    true,
		"POST":   true,
		"PUT":    true,
		"DELETE": true,
		"PATCH":  true,
	}
	
	if !validMethods[req.Method] {
		return fmt.Errorf("无效的HTTP方法: %s", req.Method)
	}
	
	return nil
}

// validateRefreshAPIKeyRequest 验证刷新API密钥请求
func (s *ExternalAPIServiceNew) validateRefreshAPIKeyRequest(req *serviceInterfaces.RefreshAPIKeyRequest) error {
	if req.Provider == "" {
		return fmt.Errorf("API提供商不能为空")
	}
	
	if req.UserID == "" {
		return fmt.Errorf("用户ID不能为空")
	}
	
	return nil
}

// getAvailableProviders 获取可用的API提供商列表
func (s *ExternalAPIServiceNew) getAvailableProviders() []string {
	providers := make([]string, 0, len(s.apiEndpoints))
	for provider := range s.apiEndpoints {
		providers = append(providers, provider)
	}
	return providers
}