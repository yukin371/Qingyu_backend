package ai

import (
	"context"
	"fmt"
	"time"

	"Qingyu_backend/pkg/errors"
	"Qingyu_backend/repository/interfaces"
	"Qingyu_backend/service/base"
	serviceInterfaces "Qingyu_backend/service/interfaces"
)

// AIServiceNew 新的AI服务实现
type AIServiceNew struct {
	repositoryFactory interfaces.RepositoryFactory
	contextService    serviceInterfaces.ContextService
	externalAPIService serviceInterfaces.ExternalAPIService
	adapterManager    serviceInterfaces.AdapterManager
	eventBus          base.EventBus
	validator         base.Validator
	serviceName       string
	version           string
	initialized       bool
}

// NewAIServiceNew 创建新的AI服务
func NewAIServiceNew(
	repositoryFactory interfaces.RepositoryFactory,
	contextService serviceInterfaces.ContextService,
	externalAPIService serviceInterfaces.ExternalAPIService,
	adapterManager serviceInterfaces.AdapterManager,
	eventBus base.EventBus,
) serviceInterfaces.AIService {
	service := &AIServiceNew{
		repositoryFactory:  repositoryFactory,
		contextService:     contextService,
		externalAPIService: externalAPIService,
		adapterManager:     adapterManager,
		eventBus:          eventBus,
		validator:         base.NewBaseValidator(),
		serviceName:        "AIService",
		version:            "2.0.0",
		initialized:        false,
	}
	
	// 添加验证规则
	service.setupValidationRules()
	
	return service
}

// setupValidationRules 设置验证规则
func (s *AIServiceNew) setupValidationRules() {
	// 为GenerateContent请求添加验证规则
	s.validator.AddRule(base.NewRequiredRule("model"))
	s.validator.AddRule(base.NewRequiredRule("prompt"))
	s.validator.AddRule(base.NewLengthRule("prompt", 1, 10000))
}

// Initialize 初始化服务
func (s *AIServiceNew) Initialize(ctx context.Context) error {
	if s.initialized {
		return nil
	}

	// 初始化上下文服务
	if err := s.contextService.Initialize(ctx); err != nil {
		return fmt.Errorf("初始化上下文服务失败: %w", err)
	}
	// 初始化外部API服务
	if err := s.externalAPIService.Initialize(ctx); err != nil {
		return fmt.Errorf("初始化外部API服务失败: %w", err)
	}
	// 初始化适配器管理器
	if err := s.adapterManager.Initialize(ctx); err != nil {
		return fmt.Errorf("初始化适配器管理器失败: %w", err)
	}

	// 设置验证规则
	s.setupValidationRules()

	s.initialized = true
	return nil
}

// Health 健康检查
func (s *AIServiceNew) Health(ctx context.Context) error {
	if !s.initialized {
		return fmt.Errorf("服务未初始化")
	}

	// 检查上下文服务健康状态
	if err := s.contextService.Health(ctx); err != nil {
		return fmt.Errorf("上下文服务健康检查失败: %w", err)
	}
	// 检查外部API服务健康状态
	if err := s.externalAPIService.Health(ctx); err != nil {
		return fmt.Errorf("外部API服务健康检查失败: %w", err)
	}
	// 检查适配器管理器健康状态
	if err := s.adapterManager.Health(ctx); err != nil {
		return fmt.Errorf("适配器管理器健康检查失败: %w", err)
	}
	// 检查Repository工厂健康状态
	if err := s.repositoryFactory.Health(ctx); err != nil {
		return fmt.Errorf("Repository工厂健康检查失败: %w", err)
	}

	return nil
}

// Close 关闭服务
func (s *AIServiceNew) Close(ctx context.Context) error {
	if !s.initialized {
		return nil
	}

	var lastErr error

	// 关闭上下文服务
	if err := s.contextService.Close(ctx); err != nil {
		lastErr = errors.AIFactory.InternalError("关闭上下文服务失败", err).
			WithOperation("Close")
	}
	// 关闭外部API服务
	if err := s.externalAPIService.Close(ctx); err != nil {
		lastErr = errors.AIFactory.InternalError("关闭外部API服务失败", err).
			WithOperation("Close")
	}
	// 关闭适配器管理器
	if err := s.adapterManager.Close(ctx); err != nil {
		lastErr = errors.AIFactory.InternalError("关闭适配器管理器失败", err).
			WithOperation("Close")
	}

	s.initialized = false
	return lastErr
}

// GetServiceName 获取服务名称
func (s *AIServiceNew) GetServiceName() string {
	return s.serviceName
}

// GetVersion 获取服务版本
func (s *AIServiceNew) GetVersion() string {
	return s.version
}

// GenerateContent 生成内容
func (s *AIServiceNew) GenerateContent(ctx context.Context, req *serviceInterfaces.GenerateContentRequest) (*serviceInterfaces.GenerateContentResponse, error) {
	startTime := time.Now()
	
	// 验证请求
	if err := s.validateGenerateContentRequest(req); err != nil {
		return nil, errors.AIFactory.ValidationError("请求验证失败", err).
			WithOperation("GenerateContent").
			WithUserID(req.UserID).
			WithMetadata(map[string]interface{}{
				"model": req.Model,
				"session_id": req.SessionID,
			})
	}
	
	// 发布请求开始事件
	event := &base.BaseEvent{
		EventType: "content.generation.started",
		EventData: map[string]interface{}{
			"model":     req.Model,
			"user_id":   req.UserID,
			"session_id": req.SessionID,
		},
		Timestamp: time.Now(),
		Source:    s.serviceName,
	}
	s.eventBus.PublishAsync(ctx, event)
	
	// 获取适配器
	adapterReq := &serviceInterfaces.GetAdapterRequest{
		AdapterID: req.Model,
	}
	adapterResp, err := s.adapterManager.GetAdapter(ctx, adapterReq)
	if err != nil {
		return nil, errors.AIFactory.NotFoundError("获取模型适配器失败", err).
			WithOperation("GenerateContent").
			WithUserID(req.UserID).
			WithMetadata(map[string]interface{}{
				"adapter_id": req.Model,
			})
	}
	
	// 调用外部API
	apiReq := &serviceInterfaces.CallAPIRequest{
		Provider: adapterResp.Adapter.Provider,
		Endpoint: "/v1/chat/completions",
		Method:   "POST",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: map[string]interface{}{
			"model":       req.Model,
			"messages":    []map[string]string{{"role": "user", "content": req.Prompt}},
			"max_tokens":  req.MaxTokens,
			"temperature": req.Temperature,
			"top_p":       req.TopP,
			"stop":        req.Stop,
		},
		Timeout: 30 * time.Second,
		UserID:  req.UserID,
	}
	
	apiResp, err := s.externalAPIService.CallAPI(ctx, apiReq)
	if err != nil {
		return nil, errors.AIFactory.ExternalError("调用外部API失败", err).
			WithOperation("GenerateContent").
			WithUserID(req.UserID).
			WithMetadata(map[string]interface{}{
				"provider": adapterResp.Adapter.Provider,
				"endpoint": "/v1/chat/completions",
				"model": req.Model,
			})
	}

	if !apiResp.Success {
		return nil, errors.AIFactory.ExternalError("外部API调用失败", fmt.Errorf("%s", apiResp.Error)).
			WithOperation("GenerateContent").
			WithUserID(req.UserID).
			WithMetadata(map[string]interface{}{
				"provider": adapterResp.Adapter.Provider,
				"error_message": apiResp.Error,
			})
	}
	
	// 解析API响应
	content, tokensUsed, finishReason := s.parseAPIResponse(apiResp.Body)
	
	// 构建响应
	response := &serviceInterfaces.GenerateContentResponse{
		Content:      content,
		Model:        req.Model,
		TokensUsed:   tokensUsed,
		FinishReason: finishReason,
		ResponseTime: time.Since(startTime),
		Metadata: map[string]string{
			"provider": adapterResp.Adapter.Provider,
			"version":  adapterResp.Adapter.Version,
		},
		RequestID: fmt.Sprintf("req_%d", time.Now().UnixNano()),
	}
	
	// 如果有会话ID，更新上下文
	if req.SessionID != "" {
		s.updateContextWithGeneration(ctx, req.SessionID, req.Prompt, content, req.UserID)
	}
	
	// 发布请求完成事件
	event = &base.BaseEvent{
		EventType: "content.generation.completed",
		EventData: map[string]interface{}{
			"model":        req.Model,
			"user_id":      req.UserID,
			"session_id":   req.SessionID,
			"tokens_used":  tokensUsed,
			"response_time": response.ResponseTime,
		},
		Timestamp: time.Now(),
		Source:    s.serviceName,
	}
	s.eventBus.PublishAsync(ctx, event)
	
	return response, nil
}

// GenerateContentStream 流式生成内容
func (s *AIServiceNew) GenerateContentStream(ctx context.Context, req *serviceInterfaces.GenerateContentRequest) (<-chan *serviceInterfaces.StreamResponse, error) {
	// 验证请求
	if err := s.validateGenerateContentRequest(req); err != nil {
		return nil, errors.AIFactory.ValidationError("请求验证失败", err).
			WithOperation("GenerateContentStream").
			WithUserID(req.UserID).
			WithMetadata(map[string]interface{}{
				"model": req.Model,
				"session_id": req.SessionID,
			})
	}

	// TODO: 实现流式生成逻辑
	responseChan := make(chan *serviceInterfaces.StreamResponse, 1)
	close(responseChan)
	return responseChan, nil
}

// AnalyzeContent 分析内容
func (s *AIServiceNew) AnalyzeContent(ctx context.Context, req *serviceInterfaces.AnalyzeContentRequest) (*serviceInterfaces.AnalyzeContentResponse, error) {
	// 验证请求
	if req.Content == "" {
		return nil, errors.AIFactory.ValidationError("内容不能为空", nil).
			WithOperation("AnalyzeContent").
			WithUserID(req.UserID)
	}
	if req.AnalysisType == "" {
		return nil, errors.AIFactory.ValidationError("分析类型不能为空", nil).
			WithOperation("AnalyzeContent").
			WithUserID(req.UserID)
	}

	// TODO: 实现内容分析逻辑
	return &serviceInterfaces.AnalyzeContentResponse{
		Analysis: "分析结果",
		Score:    0.8,
		Suggestions: []string{
			"建议1",
			"建议2",
		},
	}, nil
}

// ContinueWriting 续写内容
func (s *AIServiceNew) ContinueWriting(ctx context.Context, req *serviceInterfaces.ContinueWritingRequest) (*serviceInterfaces.ContinueWritingResponse, error) {
	// 验证请求
	if req.Content == "" {
		return nil, errors.AIFactory.ValidationError("内容不能为空", nil).
			WithOperation("ContinueWriting").
			WithUserID(req.UserID)
	}

	// TODO: 实现续写逻辑
	return &serviceInterfaces.ContinueWritingResponse{
		ContinuedContent: "续写内容",
		Confidence:       0.9,
	}, nil
}

// OptimizeText 优化文本
func (s *AIServiceNew) OptimizeText(ctx context.Context, req *serviceInterfaces.OptimizeTextRequest) (*serviceInterfaces.OptimizeTextResponse, error) {
	// 验证请求
	if req.Text == "" {
		return nil, errors.AIFactory.ValidationError("文本不能为空", nil).
			WithOperation("OptimizeText").
			WithUserID(req.UserID)
	}
	if req.OptimizationType == "" {
		return nil, errors.AIFactory.ValidationError("优化类型不能为空", nil).
			WithOperation("OptimizeText").
			WithUserID(req.UserID)
	}

	// TODO: 实现文本优化逻辑
	return &serviceInterfaces.OptimizeTextResponse{
		OptimizedText: "优化后的文本",
		Changes: []serviceInterfaces.TextChange{
			{
				Type:        "grammar",
				Original:    "原文",
				Optimized:   "优化后",
				Explanation: "语法优化",
			},
		},
	}, nil
}

// GenerateOutline 生成大纲
func (s *AIServiceNew) GenerateOutline(ctx context.Context, req *serviceInterfaces.GenerateOutlineRequest) (*serviceInterfaces.GenerateOutlineResponse, error) {
	// 验证请求
	if req.Topic == "" {
		return nil, errors.AIFactory.ValidationError("主题不能为空", nil).
			WithOperation("GenerateOutline").
			WithUserID(req.UserID)
	}

	// TODO: 实现大纲生成逻辑
	return &serviceInterfaces.GenerateOutlineResponse{
		Outline: []serviceInterfaces.OutlineItem{
			{
				Level:   1,
				Title:   "第一章",
				Content: "章节内容",
				Children: []serviceInterfaces.OutlineItem{
					{
						Level:   2,
						Title:   "第一节",
						Content: "节内容",
					},
				},
			},
		},
	}, nil
}

// GetContextInfo 获取上下文信息
func (s *AIServiceNew) GetContextInfo(ctx context.Context, req *serviceInterfaces.GetContextInfoRequest) (*serviceInterfaces.GetContextInfoResponse, error) {
	// 获取上下文
	contextReq := &serviceInterfaces.GetContextRequest{
		ContextID: req.ContextID,
		UserID:    req.UserID,
	}
	contextResp, err := s.contextService.GetContext(ctx, contextReq)
	if err != nil {
		return nil, errors.AIFactory.NotFoundError("获取上下文失败", err).
			WithOperation("GetContextInfo").
			WithUserID(req.UserID).
			WithMetadata(map[string]interface{}{
				"context_id": req.ContextID,
			})
	}

	// 获取消息列表
	messagesReq := &serviceInterfaces.GetMessagesRequest{
		ContextID: req.ContextID,
		UserID:    req.UserID,
		Limit:     req.MessageLimit,
	}
	messagesResp, err := s.contextService.GetMessages(ctx, messagesReq)
	if err != nil {
		return nil, errors.AIFactory.InternalError("获取消息失败", err).
			WithOperation("GetContextInfo").
			WithUserID(req.UserID).
			WithMetadata(map[string]interface{}{
				"context_id": req.ContextID,
			})
	}

	return &serviceInterfaces.GetContextInfoResponse{
		Context:  contextResp.Context,
		Messages: messagesResp.Messages,
		Total:    messagesResp.Total,
	}, nil
}

// UpdateContextWithFeedback 使用反馈更新上下文
func (s *AIServiceNew) UpdateContextWithFeedback(ctx context.Context, req *serviceInterfaces.UpdateContextWithFeedbackRequest) (*serviceInterfaces.UpdateContextWithFeedbackResponse, error) {
	// 验证请求
	if req.ContextID == "" {
		return nil, errors.AIFactory.ValidationError("上下文ID不能为空", nil).
			WithOperation("UpdateContextWithFeedback").
			WithUserID(req.UserID)
	}
	if req.Feedback == "" {
		return nil, errors.AIFactory.ValidationError("反馈内容不能为空", nil).
			WithOperation("UpdateContextWithFeedback").
			WithUserID(req.UserID)
	}

	// 添加反馈消息
	addMessageReq := &serviceInterfaces.AddMessageRequest{
		ContextID: req.ContextID,
		UserID:    req.UserID,
		Role:      "feedback",
		Content:   req.Feedback,
		Metadata: map[string]string{
			"feedback_type": req.FeedbackType,
			"rating":        fmt.Sprintf("%d", req.Rating),
		},
	}
	_, err := s.contextService.AddMessage(ctx, addMessageReq)
	if err != nil {
		return nil, errors.AIFactory.InternalError("添加反馈消息失败", err).
			WithOperation("UpdateContextWithFeedback").
			WithUserID(req.UserID).
			WithMetadata(map[string]interface{}{
				"context_id": req.ContextID,
				"feedback_type": req.FeedbackType,
			})
	}

	return &serviceInterfaces.UpdateContextWithFeedbackResponse{
		Success: true,
		Message: "反馈已记录",
	}, nil
}

// GetSupportedModels 获取支持的模型列表
func (s *AIServiceNew) GetSupportedModels(ctx context.Context) (*serviceInterfaces.GetSupportedModelsResponse, error) {
	// 获取适配器列表
	adaptersResp, err := s.adapterManager.ListAdapters(ctx, &serviceInterfaces.ListAdaptersRequest{})
	if err != nil {
		return nil, errors.AIFactory.InternalError("获取适配器列表失败", err).
			WithOperation("GetSupportedModels")
	}

	var models []serviceInterfaces.ModelInfo
	for _, adapter := range adaptersResp.Adapters {
		for _, model := range adapter.SupportedModels {
			models = append(models, serviceInterfaces.ModelInfo{
				ID:          model.ID,
				Name:        model.Name,
				Provider:    adapter.Provider,
				Type:        model.Type,
				MaxTokens:   model.MaxTokens,
				Description: model.Description,
				Pricing:     model.Pricing,
			})
		}
	}

	return &serviceInterfaces.GetSupportedModelsResponse{
		Models: models,
		Total:  len(models),
	}, nil
}

// GetModelInfo 获取模型信息
func (s *AIServiceNew) GetModelInfo(ctx context.Context, req *serviceInterfaces.GetModelInfoRequest) (*serviceInterfaces.GetModelInfoResponse, error) {
	// 验证请求
	if req.ModelID == "" {
		return nil, errors.AIFactory.ValidationError("模型ID不能为空", nil).
			WithOperation("GetModelInfo")
	}

	// 获取模型配置
	modelReq := &serviceInterfaces.GetModelConfigRequest{
		ModelID: req.ModelID,
	}
	modelResp, err := s.adapterManager.GetModelConfig(ctx, modelReq)
	if err != nil {
		return nil, errors.AIFactory.NotFoundError("模型不存在", err).
			WithOperation("GetModelInfo").
			WithMetadata(map[string]interface{}{
				"model_id": req.ModelID,
			})
	}

	return &serviceInterfaces.GetModelInfoResponse{
		Model: serviceInterfaces.ModelInfo{
			ID:          modelResp.Model.ID,
			Name:        modelResp.Model.Name,
			Provider:    modelResp.Model.Provider,
			Type:        modelResp.Model.Type,
			MaxTokens:   modelResp.Model.MaxTokens,
			Description: modelResp.Model.Description,
			Pricing:     modelResp.Model.Pricing,
		},
	}, nil
}

// ValidateAPIKey 验证API密钥
func (s *AIServiceNew) ValidateAPIKey(ctx context.Context, req *serviceInterfaces.ValidateAPIKeyRequest) (*serviceInterfaces.ValidateAPIKeyResponse, error) {
	// 验证请求
	if req.Provider == "" {
		return nil, errors.AIFactory.ValidationError("提供商不能为空", nil).
			WithOperation("ValidateAPIKey")
	}
	if req.APIKey == "" {
		return nil, errors.AIFactory.ValidationError("API密钥不能为空", nil).
			WithOperation("ValidateAPIKey")
	}

	// TODO: 实现API密钥验证逻辑
	// 这里应该调用相应提供商的API来验证密钥有效性
	
	return &serviceInterfaces.ValidateAPIKeyResponse{
		Valid:   true,
		Message: "API密钥有效",
		Details: map[string]interface{}{
			"provider": req.Provider,
			"quota":    "unlimited",
		},
	}, nil
}

// 辅助方法

// validateGenerateContentRequest 验证生成内容请求
func (s *AIServiceNew) validateGenerateContentRequest(req *serviceInterfaces.GenerateContentRequest) error {
	if req.Model == "" {
		return fmt.Errorf("模型不能为空")
	}
	
	if req.Prompt == "" {
		return fmt.Errorf("提示不能为空")
	}
	
	if len(req.Prompt) > 10000 {
		return fmt.Errorf("提示长度不能超过10000个字符")
	}
	
	if req.MaxTokens < 0 {
		return fmt.Errorf("最大token数不能为负数")
	}
	
	if req.Temperature < 0 || req.Temperature > 2 {
		return fmt.Errorf("温度值必须在0到2之间")
	}
	
	return nil
}

// parseAPIResponse 解析API响应
func (s *AIServiceNew) parseAPIResponse(body interface{}) (content string, tokensUsed int, finishReason string) {
	// 这里应该根据实际的API响应格式进行解析
	// 为了简化，我们返回模拟数据
	return "这是AI生成的内容", 50, "stop"
}

// updateContextWithGeneration 更新上下文
func (s *AIServiceNew) updateContextWithGeneration(ctx context.Context, sessionID, prompt, response, userID string) {
	// 添加用户消息
	userMessageReq := &serviceInterfaces.AddMessageRequest{
		ContextID: sessionID,
		Role:      "user",
		Content:   prompt,
		UserID:    userID,
	}
	s.contextService.AddMessage(ctx, userMessageReq)
	
	// 添加AI响应
	aiMessageReq := &serviceInterfaces.AddMessageRequest{
		ContextID: sessionID,
		Role:      "assistant",
		Content:   response,
		UserID:    userID,
	}
	s.contextService.AddMessage(ctx, aiMessageReq)
}