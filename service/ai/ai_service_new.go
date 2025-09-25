package ai

import (
	"context"
	"fmt"
	"time"

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
	
	// 初始化依赖服务
	if err := s.contextService.Initialize(ctx); err != nil {
		return base.NewServiceError(s.serviceName, base.ErrorTypeInternal, "初始化上下文服务失败", err)
	}
	
	if err := s.externalAPIService.Initialize(ctx); err != nil {
		return base.NewServiceError(s.serviceName, base.ErrorTypeInternal, "初始化外部API服务失败", err)
	}
	
	if err := s.adapterManager.Initialize(ctx); err != nil {
		return base.NewServiceError(s.serviceName, base.ErrorTypeInternal, "初始化适配器管理器失败", err)
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
		// 事件发布失败不应该影响服务初始化
		fmt.Printf("发布服务初始化事件失败: %v\n", err)
	}
	
	return nil
}

// Health 健康检查
func (s *AIServiceNew) Health(ctx context.Context) error {
	if !s.initialized {
		return base.NewServiceError(s.serviceName, base.ErrorTypeInternal, "服务未初始化", nil)
	}
	
	// 检查依赖服务健康状态
	if err := s.contextService.Health(ctx); err != nil {
		return base.NewServiceError(s.serviceName, base.ErrorTypeInternal, "上下文服务健康检查失败", err)
	}
	
	if err := s.externalAPIService.Health(ctx); err != nil {
		return base.NewServiceError(s.serviceName, base.ErrorTypeInternal, "外部API服务健康检查失败", err)
	}
	
	if err := s.adapterManager.Health(ctx); err != nil {
		return base.NewServiceError(s.serviceName, base.ErrorTypeInternal, "适配器管理器健康检查失败", err)
	}
	
	// 检查Repository工厂健康状态
	if err := s.repositoryFactory.Health(ctx); err != nil {
		return base.NewServiceError(s.serviceName, base.ErrorTypeInternal, "Repository工厂健康检查失败", err)
	}
	
	return nil
}

// Close 关闭服务
func (s *AIServiceNew) Close(ctx context.Context) error {
	if !s.initialized {
		return nil
	}
	
	var lastErr error
	
	// 关闭依赖服务
	if err := s.contextService.Close(ctx); err != nil {
		lastErr = base.NewServiceError(s.serviceName, base.ErrorTypeInternal, "关闭上下文服务失败", err)
	}
	
	if err := s.externalAPIService.Close(ctx); err != nil {
		lastErr = base.NewServiceError(s.serviceName, base.ErrorTypeInternal, "关闭外部API服务失败", err)
	}
	
	if err := s.adapterManager.Close(ctx); err != nil {
		lastErr = base.NewServiceError(s.serviceName, base.ErrorTypeInternal, "关闭适配器管理器失败", err)
	}
	
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
		return nil, base.NewServiceError(s.serviceName, base.ErrorTypeValidation, "请求验证失败", err)
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
		return nil, base.NewServiceError(s.serviceName, base.ErrorTypeNotFound, "获取模型适配器失败", err)
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
		return nil, base.NewServiceError(s.serviceName, base.ErrorTypeExternal, "调用外部API失败", err)
	}
	
	if !apiResp.Success {
		return nil, base.NewServiceError(s.serviceName, base.ErrorTypeExternal, "外部API调用失败", fmt.Errorf("%s", apiResp.Error))
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
		return nil, base.NewServiceError(s.serviceName, base.ErrorTypeValidation, "请求验证失败", err)
	}
	
	// 创建响应通道
	responseChan := make(chan *serviceInterfaces.StreamResponse)
	
	go func() {
		defer close(responseChan)
		
		// 这里应该实现真正的流式调用
		// 为了简化，我们模拟流式响应
		content := "这是一个模拟的流式响应内容。"
		words := []string{"这是", "一个", "模拟的", "流式", "响应", "内容。"}
		
		for i, word := range words {
			select {
			case <-ctx.Done():
				responseChan <- &serviceInterfaces.StreamResponse{
					Error: "请求被取消",
					Done:  true,
				}
				return
			default:
				responseChan <- &serviceInterfaces.StreamResponse{
					Content:    content[:len(word)*(i+1)],
					Delta:      word,
					Done:       i == len(words)-1,
					TokensUsed: i + 1,
					Metadata: map[string]string{
						"model": req.Model,
					},
				}
				time.Sleep(100 * time.Millisecond) // 模拟延迟
			}
		}
	}()
	
	return responseChan, nil
}

// AnalyzeContent 分析内容
func (s *AIServiceNew) AnalyzeContent(ctx context.Context, req *serviceInterfaces.AnalyzeContentRequest) (*serviceInterfaces.AnalyzeContentResponse, error) {
	startTime := time.Now()
	
	// 验证请求
	if req.Content == "" {
		return nil, base.NewServiceError(s.serviceName, base.ErrorTypeValidation, "内容不能为空", nil)
	}
	
	if req.AnalysisType == "" {
		return nil, base.NewServiceError(s.serviceName, base.ErrorTypeValidation, "分析类型不能为空", nil)
	}
	
	// 模拟内容分析
	analysis := map[string]interface{}{
		"word_count":      len(req.Content),
		"character_count": len([]rune(req.Content)),
		"analysis_type":   req.AnalysisType,
	}
	
	response := &serviceInterfaces.AnalyzeContentResponse{
		Analysis:     analysis,
		Summary:      "这是一个内容分析的摘要",
		Keywords:     []string{"关键词1", "关键词2", "关键词3"},
		Sentiment:    "neutral",
		Topics:       []string{"主题1", "主题2"},
		Language:     "zh-CN",
		Confidence:   0.85,
		ResponseTime: time.Since(startTime),
	}
	
	return response, nil
}

// ContinueWriting 续写内容
func (s *AIServiceNew) ContinueWriting(ctx context.Context, req *serviceInterfaces.ContinueWritingRequest) (*serviceInterfaces.ContinueWritingResponse, error) {
	startTime := time.Now()
	
	// 验证请求
	if req.Content == "" {
		return nil, base.NewServiceError(s.serviceName, base.ErrorTypeValidation, "内容不能为空", nil)
	}
	
	// 模拟续写
	continuedContent := "这是续写的内容，基于原始内容进行扩展。"
	
	response := &serviceInterfaces.ContinueWritingResponse{
		ContinuedContent: continuedContent,
		OriginalLength:   len(req.Content),
		AddedLength:      len(continuedContent),
		ResponseTime:     time.Since(startTime),
	}
	
	return response, nil
}

// OptimizeText 优化文本
func (s *AIServiceNew) OptimizeText(ctx context.Context, req *serviceInterfaces.OptimizeTextRequest) (*serviceInterfaces.OptimizeTextResponse, error) {
	startTime := time.Now()
	
	// 验证请求
	if req.Text == "" {
		return nil, base.NewServiceError(s.serviceName, base.ErrorTypeValidation, "文本不能为空", nil)
	}
	
	if req.OptimizationType == "" {
		return nil, base.NewServiceError(s.serviceName, base.ErrorTypeValidation, "优化类型不能为空", nil)
	}
	
	// 模拟文本优化
	optimizedText := "这是优化后的文本内容。"
	
	response := &serviceInterfaces.OptimizeTextResponse{
		OptimizedText:   optimizedText,
		Changes:         []string{"修正了语法错误", "改进了表达方式"},
		Improvements:    []string{"提高了可读性", "增强了逻辑性"},
		OriginalLength:  len(req.Text),
		OptimizedLength: len(optimizedText),
		Suggestions:     []string{"建议1", "建议2"},
		Metadata: map[string]string{
			"optimization_type": req.OptimizationType,
		},
		ResponseTime: time.Since(startTime),
	}
	
	return response, nil
}

// GenerateOutline 生成大纲
func (s *AIServiceNew) GenerateOutline(ctx context.Context, req *serviceInterfaces.GenerateOutlineRequest) (*serviceInterfaces.GenerateOutlineResponse, error) {
	startTime := time.Now()
	
	// 验证请求
	if req.Topic == "" {
		return nil, base.NewServiceError(s.serviceName, base.ErrorTypeValidation, "主题不能为空", nil)
	}
	
	// 模拟生成大纲
	outline := []serviceInterfaces.OutlineItem{
		{
			Level:       1,
			Title:       "引言",
			Description: "介绍主题背景",
		},
		{
			Level:       1,
			Title:       "主体内容",
			Description: "详细阐述主题",
			Children: []serviceInterfaces.OutlineItem{
				{
					Level:       2,
					Title:       "子主题1",
					Description: "第一个子主题",
				},
				{
					Level:       2,
					Title:       "子主题2",
					Description: "第二个子主题",
				},
			},
		},
		{
			Level:       1,
			Title:       "结论",
			Description: "总结全文",
		},
	}
	
	response := &serviceInterfaces.GenerateOutlineResponse{
		Outline:         outline,
		Title:           fmt.Sprintf("关于%s的大纲", req.Topic),
		Summary:         "这是一个关于指定主题的详细大纲",
		EstimatedLength: 1500,
		ResponseTime:    time.Since(startTime),
	}
	
	return response, nil
}

// GetContextInfo 获取上下文信息
func (s *AIServiceNew) GetContextInfo(ctx context.Context, req *serviceInterfaces.GetContextInfoRequest) (*serviceInterfaces.GetContextInfoResponse, error) {
	// 委托给上下文服务
	contextReq := &serviceInterfaces.GetContextRequest{
		ContextID: req.ContextID,
		UserID:    req.UserID,
	}
	
	contextResp, err := s.contextService.GetContext(ctx, contextReq)
	if err != nil {
		return nil, base.NewServiceError(s.serviceName, base.ErrorTypeNotFound, "获取上下文失败", err)
	}
	
	// 获取消息数量
	messagesReq := &serviceInterfaces.GetMessagesRequest{
		ContextID: req.ContextID,
		UserID:    req.UserID,
	}
	
	messagesResp, err := s.contextService.GetMessages(ctx, messagesReq)
	if err != nil {
		return nil, base.NewServiceError(s.serviceName, base.ErrorTypeInternal, "获取消息失败", err)
	}
	
	response := &serviceInterfaces.GetContextInfoResponse{
		ContextID:    contextResp.Context.ID,
		MessageCount: messagesResp.Total,
		TokensUsed:   0, // 这里需要计算实际使用的token数量
		CreatedAt:    contextResp.Context.CreatedAt,
		UpdatedAt:    contextResp.Context.UpdatedAt,
		Metadata:     contextResp.Context.Metadata,
		Status:       contextResp.Context.Status,
	}
	
	return response, nil
}

// UpdateContextWithFeedback 根据反馈更新上下文
func (s *AIServiceNew) UpdateContextWithFeedback(ctx context.Context, req *serviceInterfaces.UpdateContextWithFeedbackRequest) (*serviceInterfaces.UpdateContextWithFeedbackResponse, error) {
	startTime := time.Now()
	
	// 验证请求
	if req.ContextID == "" {
		return nil, base.NewServiceError(s.serviceName, base.ErrorTypeValidation, "上下文ID不能为空", nil)
	}
	
	if req.Feedback == "" {
		return nil, base.NewServiceError(s.serviceName, base.ErrorTypeValidation, "反馈内容不能为空", nil)
	}
	
	// 添加反馈消息到上下文
	addMessageReq := &serviceInterfaces.AddMessageRequest{
		ContextID: req.ContextID,
		Role:      "feedback",
		Content:   req.Feedback,
		Metadata: map[string]string{
			"rating": fmt.Sprintf("%d", req.Rating),
			"type":   "user_feedback",
		},
		UserID: req.UserID,
	}
	
	_, err := s.contextService.AddMessage(ctx, addMessageReq)
	if err != nil {
		return nil, base.NewServiceError(s.serviceName, base.ErrorTypeInternal, "添加反馈消息失败", err)
	}
	
	response := &serviceInterfaces.UpdateContextWithFeedbackResponse{
		ContextID:    req.ContextID,
		Updated:      true,
		Changes:      []string{"添加了用户反馈"},
		ResponseTime: time.Since(startTime),
	}
	
	return response, nil
}

// GetSupportedModels 获取支持的AI模型列表
func (s *AIServiceNew) GetSupportedModels(ctx context.Context) (*serviceInterfaces.GetSupportedModelsResponse, error) {
	// 从适配器管理器获取所有适配器
	listReq := &serviceInterfaces.ListAdaptersRequest{
		Type:   "ai_model",
		Status: "active",
	}
	
	listResp, err := s.adapterManager.ListAdapters(ctx, listReq)
	if err != nil {
		return nil, base.NewServiceError(s.serviceName, base.ErrorTypeInternal, "获取适配器列表失败", err)
	}
	
	// 转换为模型信息
	models := make([]serviceInterfaces.ModelInfo, len(listResp.Adapters))
	for i, adapter := range listResp.Adapters {
		models[i] = serviceInterfaces.ModelInfo{
			ID:          adapter.ID,
			Name:        adapter.Name,
			Provider:    adapter.Provider,
			Type:        adapter.Type,
			MaxTokens:   4096, // 默认值，应该从配置中获取
			InputPrice:  0.001,
			OutputPrice: 0.002,
			Features:    adapter.Features,
			Status:      adapter.Status,
			Description: adapter.Description,
			Metadata:    adapter.Config,
		}
	}
	
	response := &serviceInterfaces.GetSupportedModelsResponse{
		Models: models,
	}
	
	return response, nil
}

// GetModelInfo 获取模型信息
func (s *AIServiceNew) GetModelInfo(ctx context.Context, req *serviceInterfaces.GetModelInfoRequest) (*serviceInterfaces.GetModelInfoResponse, error) {
	// 验证请求
	if req.ModelID == "" {
		return nil, base.NewServiceError(s.serviceName, base.ErrorTypeValidation, "模型ID不能为空", nil)
	}
	
	// 从适配器管理器获取适配器信息
	adapterReq := &serviceInterfaces.GetAdapterRequest{
		AdapterID: req.ModelID,
	}
	
	adapterResp, err := s.adapterManager.GetAdapter(ctx, adapterReq)
	if err != nil {
		return nil, base.NewServiceError(s.serviceName, base.ErrorTypeNotFound, "模型不存在", err)
	}
	
	model := serviceInterfaces.ModelInfo{
		ID:          adapterResp.Adapter.ID,
		Name:        adapterResp.Adapter.Name,
		Provider:    adapterResp.Adapter.Provider,
		Type:        adapterResp.Adapter.Type,
		MaxTokens:   4096,
		InputPrice:  0.001,
		OutputPrice: 0.002,
		Features:    adapterResp.Adapter.Features,
		Status:      adapterResp.Adapter.Status,
		Description: adapterResp.Adapter.Description,
		Metadata:    adapterResp.Adapter.Config,
	}
	
	response := &serviceInterfaces.GetModelInfoResponse{
		Model: model,
	}
	
	return response, nil
}

// ValidateAPIKey 验证API密钥
func (s *AIServiceNew) ValidateAPIKey(ctx context.Context, req *serviceInterfaces.ValidateAPIKeyRequest) (*serviceInterfaces.ValidateAPIKeyResponse, error) {
	startTime := time.Now()
	
	// 验证请求
	if req.Provider == "" {
		return nil, base.NewServiceError(s.serviceName, base.ErrorTypeValidation, "提供商不能为空", nil)
	}
	
	if req.APIKey == "" {
		return nil, base.NewServiceError(s.serviceName, base.ErrorTypeValidation, "API密钥不能为空", nil)
	}
	
	// 调用外部API验证密钥
	apiReq := &serviceInterfaces.CallAPIRequest{
		Provider: req.Provider,
		Endpoint: "/v1/models", // 通用的验证端点
		Method:   "GET",
		Headers: map[string]string{
			"Authorization": "Bearer " + req.APIKey,
		},
		Timeout: 10 * time.Second,
	}
	
	apiResp, err := s.externalAPIService.CallAPI(ctx, apiReq)
	if err != nil {
		return &serviceInterfaces.ValidateAPIKeyResponse{
			Valid:        false,
			Provider:     req.Provider,
			Message:      "API密钥验证失败",
			ResponseTime: time.Since(startTime),
		}, nil
	}
	
	response := &serviceInterfaces.ValidateAPIKeyResponse{
		Valid:        apiResp.Success,
		Provider:     req.Provider,
		Message:      "API密钥验证成功",
		ResponseTime: time.Since(startTime),
	}
	
	if apiResp.Success {
		// 如果验证成功，尝试获取配额信息
		response.Quota = &serviceInterfaces.APIQuota{
			Used:      100,
			Limit:     1000,
			Remaining: 900,
			ResetAt:   time.Now().Add(24 * time.Hour),
		}
	}
	
	return response, nil
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