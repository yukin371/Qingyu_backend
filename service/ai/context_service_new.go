package ai

import (
	"context"
	"fmt"
	"time"

	"Qingyu_backend/repository/interfaces"
	"Qingyu_backend/service/base"
	serviceInterfaces "Qingyu_backend/service/interfaces"
)

// ContextServiceNew 新的上下文服务实现
type ContextServiceNew struct {
	repositoryFactory interfaces.RepositoryFactory
	eventBus          base.EventBus
	validator         base.Validator
	serviceName       string
	version           string
	initialized       bool
}

// NewContextServiceNew 创建新的上下文服务
func NewContextServiceNew(
	repositoryFactory interfaces.RepositoryFactory,
	eventBus base.EventBus,
) serviceInterfaces.ContextService {
	service := &ContextServiceNew{
		repositoryFactory: repositoryFactory,
		eventBus:          eventBus,
		validator:         base.NewBaseValidator(),
		serviceName:       "ContextService",
		version:           "1.0.0",
		initialized:       false,
	}
	
	// 添加验证规则
	service.setupValidationRules()
	
	return service
}

// setupValidationRules 设置验证规则
func (s *ContextServiceNew) setupValidationRules() {
	s.validator.AddRule(base.NewRequiredRule("name"))
	s.validator.AddRule(base.NewRequiredRule("user_id"))
	s.validator.AddRule(base.NewLengthRule("name", 1, 100))
	s.validator.AddRule(base.NewLengthRule("description", 0, 500))
}

// Initialize 初始化服务
func (s *ContextServiceNew) Initialize(ctx context.Context) error {
	if s.initialized {
		return nil
	}
	
	// 检查Repository工厂健康状态
	if err := s.repositoryFactory.Health(ctx); err != nil {
		return base.NewServiceError(s.serviceName, base.ErrorTypeInternal, "Repository工厂不可用", err)
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
func (s *ContextServiceNew) Health(ctx context.Context) error {
	if !s.initialized {
		return base.NewServiceError(s.serviceName, base.ErrorTypeInternal, "服务未初始化", nil)
	}
	
	// 检查Repository工厂健康状态
	if err := s.repositoryFactory.Health(ctx); err != nil {
		return base.NewServiceError(s.serviceName, base.ErrorTypeInternal, "Repository工厂健康检查失败", err)
	}
	
	return nil
}

// Close 关闭服务
func (s *ContextServiceNew) Close(ctx context.Context) error {
	if !s.initialized {
		return nil
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
	
	return nil
}

// GetServiceName 获取服务名称
func (s *ContextServiceNew) GetServiceName() string {
	return s.serviceName
}

// GetVersion 获取服务版本
func (s *ContextServiceNew) GetVersion() string {
	return s.version
}

// CreateContext 创建上下文
func (s *ContextServiceNew) CreateContext(ctx context.Context, req *serviceInterfaces.CreateContextRequest) (*serviceInterfaces.CreateContextResponse, error) {
	// 验证请求
	if err := s.validateCreateContextRequest(req); err != nil {
		return nil, base.NewServiceError(s.serviceName, base.ErrorTypeValidation, "请求验证失败", err)
	}
	
	// 生成上下文ID
	contextID := fmt.Sprintf("ctx_%d_%s", time.Now().UnixNano(), req.UserID)
	
	// 这里应该将上下文信息保存到数据库
	// 为了简化，我们模拟保存操作
	
	// 发布上下文创建事件
	event := &base.BaseEvent{
		EventType: "context.created",
		EventData: map[string]interface{}{
			"context_id": contextID,
			"user_id":    req.UserID,
			"name":       req.Name,
			"type":       req.Type,
		},
		Timestamp: time.Now(),
		Source:    s.serviceName,
	}
	s.eventBus.PublishAsync(ctx, event)
	
	response := &serviceInterfaces.CreateContextResponse{
		ContextID: contextID,
		CreatedAt: time.Now(),
	}
	
	return response, nil
}

// GetContext 获取上下文
func (s *ContextServiceNew) GetContext(ctx context.Context, req *serviceInterfaces.GetContextRequest) (*serviceInterfaces.GetContextResponse, error) {
	// 验证请求
	if req.ContextID == "" {
		return nil, base.NewServiceError(s.serviceName, base.ErrorTypeValidation, "上下文ID不能为空", nil)
	}
	
	// 这里应该从数据库获取上下文信息
	// 为了简化，我们返回模拟数据
	contextInfo := serviceInterfaces.ContextInfo{
		ID:          req.ContextID,
		Name:        "示例上下文",
		Description: "这是一个示例上下文",
		Type:        "conversation",
		Status:      "active",
		UserID:      req.UserID,
		Metadata: map[string]string{
			"created_by": "system",
		},
		CreatedAt: time.Now().Add(-24 * time.Hour),
		UpdatedAt: time.Now(),
	}
	
	response := &serviceInterfaces.GetContextResponse{
		Context: contextInfo,
	}
	
	return response, nil
}

// UpdateContext 更新上下文
func (s *ContextServiceNew) UpdateContext(ctx context.Context, req *serviceInterfaces.UpdateContextRequest) (*serviceInterfaces.UpdateContextResponse, error) {
	// 验证请求
	if req.ContextID == "" {
		return nil, base.NewServiceError(s.serviceName, base.ErrorTypeValidation, "上下文ID不能为空", nil)
	}
	
	// 这里应该更新数据库中的上下文信息
	// 为了简化，我们模拟更新操作
	
	// 发布上下文更新事件
	event := &base.BaseEvent{
		EventType: "context.updated",
		EventData: map[string]interface{}{
			"context_id": req.ContextID,
			"user_id":    req.UserID,
		},
		Timestamp: time.Now(),
		Source:    s.serviceName,
	}
	s.eventBus.PublishAsync(ctx, event)
	
	response := &serviceInterfaces.UpdateContextResponse{
		Updated:   true,
		UpdatedAt: time.Now(),
	}
	
	return response, nil
}

// DeleteContext 删除上下文
func (s *ContextServiceNew) DeleteContext(ctx context.Context, req *serviceInterfaces.DeleteContextRequest) (*serviceInterfaces.DeleteContextResponse, error) {
	// 验证请求
	if req.ContextID == "" {
		return nil, base.NewServiceError(s.serviceName, base.ErrorTypeValidation, "上下文ID不能为空", nil)
	}
	
	// 这里应该从数据库删除上下文信息
	// 为了简化，我们模拟删除操作
	
	// 发布上下文删除事件
	event := &base.BaseEvent{
		EventType: "context.deleted",
		EventData: map[string]interface{}{
			"context_id": req.ContextID,
			"user_id":    req.UserID,
		},
		Timestamp: time.Now(),
		Source:    s.serviceName,
	}
	s.eventBus.PublishAsync(ctx, event)
	
	response := &serviceInterfaces.DeleteContextResponse{
		Deleted:   true,
		DeletedAt: time.Now(),
	}
	
	return response, nil
}

// ListContexts 列出上下文
func (s *ContextServiceNew) ListContexts(ctx context.Context, req *serviceInterfaces.ListContextsRequest) (*serviceInterfaces.ListContextsResponse, error) {
	// 验证请求
	if req.UserID == "" {
		return nil, base.NewServiceError(s.serviceName, base.ErrorTypeValidation, "用户ID不能为空", nil)
	}
	
	// 设置默认分页参数
	page := req.Page
	if page <= 0 {
		page = 1
	}
	
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	
	// 这里应该从数据库查询上下文列表
	// 为了简化，我们返回模拟数据
	contexts := []serviceInterfaces.ContextInfo{
		{
			ID:          "ctx_1",
			Name:        "对话上下文1",
			Description: "第一个对话上下文",
			Type:        "conversation",
			Status:      "active",
			UserID:      req.UserID,
			Metadata:    map[string]string{"tag": "important"},
			CreatedAt:   time.Now().Add(-48 * time.Hour),
			UpdatedAt:   time.Now().Add(-1 * time.Hour),
		},
		{
			ID:          "ctx_2",
			Name:        "对话上下文2",
			Description: "第二个对话上下文",
			Type:        "conversation",
			Status:      "active",
			UserID:      req.UserID,
			Metadata:    map[string]string{"tag": "normal"},
			CreatedAt:   time.Now().Add(-24 * time.Hour),
			UpdatedAt:   time.Now(),
		},
	}
	
	total := len(contexts)
	totalPages := (total + pageSize - 1) / pageSize
	
	response := &serviceInterfaces.ListContextsResponse{
		Contexts:   contexts,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
	
	return response, nil
}

// AddMessage 添加消息到上下文
func (s *ContextServiceNew) AddMessage(ctx context.Context, req *serviceInterfaces.AddMessageRequest) (*serviceInterfaces.AddMessageResponse, error) {
	// 验证请求
	if err := s.validateAddMessageRequest(req); err != nil {
		return nil, base.NewServiceError(s.serviceName, base.ErrorTypeValidation, "请求验证失败", err)
	}
	
	// 生成消息ID
	messageID := fmt.Sprintf("msg_%d", time.Now().UnixNano())
	
	// 这里应该将消息保存到数据库
	// 为了简化，我们模拟保存操作
	
	// 发布消息添加事件
	event := &base.BaseEvent{
		EventType: "message.added",
		EventData: map[string]interface{}{
			"context_id": req.ContextID,
			"message_id": messageID,
			"role":       req.Role,
			"user_id":    req.UserID,
		},
		Timestamp: time.Now(),
		Source:    s.serviceName,
	}
	s.eventBus.PublishAsync(ctx, event)
	
	response := &serviceInterfaces.AddMessageResponse{
		MessageID: messageID,
		AddedAt:   time.Now(),
	}
	
	return response, nil
}

// GetMessages 获取上下文消息
func (s *ContextServiceNew) GetMessages(ctx context.Context, req *serviceInterfaces.GetMessagesRequest) (*serviceInterfaces.GetMessagesResponse, error) {
	// 验证请求
	if req.ContextID == "" {
		return nil, base.NewServiceError(s.serviceName, base.ErrorTypeValidation, "上下文ID不能为空", nil)
	}
	
	// 设置默认参数
	limit := req.Limit
	if limit <= 0 {
		limit = 50
	}
	
	offset := req.Offset
	if offset < 0 {
		offset = 0
	}
	
	// 这里应该从数据库查询消息列表
	// 为了简化，我们返回模拟数据
	messages := []serviceInterfaces.MessageInfo{
		{
			ID:        "msg_1",
			ContextID: req.ContextID,
			Role:      "user",
			Content:   "你好，请帮我写一篇文章",
			Metadata:  map[string]string{"type": "text"},
			CreatedAt: time.Now().Add(-1 * time.Hour),
		},
		{
			ID:        "msg_2",
			ContextID: req.ContextID,
			Role:      "assistant",
			Content:   "好的，我来帮您写文章。请告诉我文章的主题和要求。",
			Metadata:  map[string]string{"type": "text"},
			CreatedAt: time.Now().Add(-30 * time.Minute),
		},
	}
	
	response := &serviceInterfaces.GetMessagesResponse{
		Messages: messages,
		Total:    len(messages),
	}
	
	return response, nil
}

// ClearMessages 清空上下文消息
func (s *ContextServiceNew) ClearMessages(ctx context.Context, req *serviceInterfaces.ClearMessagesRequest) (*serviceInterfaces.ClearMessagesResponse, error) {
	// 验证请求
	if req.ContextID == "" {
		return nil, base.NewServiceError(s.serviceName, base.ErrorTypeValidation, "上下文ID不能为空", nil)
	}
	
	// 这里应该从数据库清空消息
	// 为了简化，我们模拟清空操作
	
	// 发布消息清空事件
	event := &base.BaseEvent{
		EventType: "messages.cleared",
		EventData: map[string]interface{}{
			"context_id": req.ContextID,
			"user_id":    req.UserID,
		},
		Timestamp: time.Now(),
		Source:    s.serviceName,
	}
	s.eventBus.PublishAsync(ctx, event)
	
	response := &serviceInterfaces.ClearMessagesResponse{
		Cleared:   true,
		ClearedAt: time.Now(),
	}
	
	return response, nil
}

// 辅助方法

// validateCreateContextRequest 验证创建上下文请求
func (s *ContextServiceNew) validateCreateContextRequest(req *serviceInterfaces.CreateContextRequest) error {
	if req.Name == "" {
		return fmt.Errorf("上下文名称不能为空")
	}
	
	if len(req.Name) > 100 {
		return fmt.Errorf("上下文名称长度不能超过100个字符")
	}
	
	if req.UserID == "" {
		return fmt.Errorf("用户ID不能为空")
	}
	
	if len(req.Description) > 500 {
		return fmt.Errorf("描述长度不能超过500个字符")
	}
	
	return nil
}

// validateAddMessageRequest 验证添加消息请求
func (s *ContextServiceNew) validateAddMessageRequest(req *serviceInterfaces.AddMessageRequest) error {
	if req.ContextID == "" {
		return fmt.Errorf("上下文ID不能为空")
	}
	
	if req.Role == "" {
		return fmt.Errorf("角色不能为空")
	}
	
	if req.Content == "" {
		return fmt.Errorf("消息内容不能为空")
	}
	
	if len(req.Content) > 10000 {
		return fmt.Errorf("消息内容长度不能超过10000个字符")
	}
	
	// 验证角色类型
	validRoles := map[string]bool{
		"user":      true,
		"assistant": true,
		"system":    true,
		"feedback":  true,
	}
	
	if !validRoles[req.Role] {
		return fmt.Errorf("无效的角色类型: %s", req.Role)
	}
	
	return nil
}