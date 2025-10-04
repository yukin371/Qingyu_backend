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

	// 检查Repository工厂是否可用
	if err := s.repositoryFactory.Health(ctx); err != nil {
		return errors.ContextFactory.InternalError("Repository工厂不可用", err).
			WithOperation("Initialize")
	}

	s.initialized = true
	return nil
}

// Health 健康检查
func (s *ContextServiceNew) Health(ctx context.Context) error {
	if !s.initialized {
		return errors.ContextFactory.InternalError("服务未初始化", nil).
			WithOperation("Health")
	}

	// 检查Repository工厂健康状态
	if err := s.repositoryFactory.Health(ctx); err != nil {
		return errors.ContextFactory.InternalError("Repository工厂健康检查失败", err).
			WithOperation("Health")
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
		return nil, errors.ContextFactory.ValidationError("请求验证失败", err).
			WithOperation("CreateContext").
			WithUserID(req.UserID).
			WithMetadata(map[string]interface{}{
				"name": req.Name,
				"type": req.Type,
			})
	}

	// TODO: 实现创建上下文逻辑
	return &serviceInterfaces.CreateContextResponse{
		Context: &serviceInterfaces.Context{
			ID:          fmt.Sprintf("ctx_%d", time.Now().UnixNano()),
			Name:        req.Name,
			Type:        req.Type,
			UserID:      req.UserID,
			Description: req.Description,
			Metadata:    req.Metadata,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}, nil
}

// GetContext 获取上下文
func (s *ContextServiceNew) GetContext(ctx context.Context, req *serviceInterfaces.GetContextRequest) (*serviceInterfaces.GetContextResponse, error) {
	// 验证请求
	if req.ContextID == "" {
		return nil, errors.ContextFactory.ValidationError("上下文ID不能为空", nil).
			WithOperation("GetContext").
			WithUserID(req.UserID)
	}

	// TODO: 实现获取上下文逻辑
	return &serviceInterfaces.GetContextResponse{
		Context: &serviceInterfaces.Context{
			ID:          req.ContextID,
			Name:        "示例上下文",
			Type:        "conversation",
			UserID:      req.UserID,
			Description: "这是一个示例上下文",
			CreatedAt:   time.Now().Add(-time.Hour),
			UpdatedAt:   time.Now(),
		},
	}, nil
}

// UpdateContext 更新上下文
func (s *ContextServiceNew) UpdateContext(ctx context.Context, req *serviceInterfaces.UpdateContextRequest) (*serviceInterfaces.UpdateContextResponse, error) {
	// 验证请求
	if req.ContextID == "" {
		return nil, errors.ContextFactory.ValidationError("上下文ID不能为空", nil).
			WithOperation("UpdateContext").
			WithUserID(req.UserID)
	}

	// TODO: 实现更新上下文逻辑
	return &serviceInterfaces.UpdateContextResponse{
		Context: &serviceInterfaces.Context{
			ID:          req.ContextID,
			Name:        req.Name,
			UserID:      req.UserID,
			Description: req.Description,
			Metadata:    req.Metadata,
			UpdatedAt:   time.Now(),
		},
	}, nil
}

// DeleteContext 删除上下文
func (s *ContextServiceNew) DeleteContext(ctx context.Context, req *serviceInterfaces.DeleteContextRequest) (*serviceInterfaces.DeleteContextResponse, error) {
	// 验证请求
	if req.ContextID == "" {
		return nil, errors.ContextFactory.ValidationError("上下文ID不能为空", nil).
			WithOperation("DeleteContext").
			WithUserID(req.UserID)
	}

	// TODO: 实现删除上下文逻辑
	return &serviceInterfaces.DeleteContextResponse{
		Success: true,
		Message: "上下文已删除",
	}, nil
}

// ListContexts 列出上下文
func (s *ContextServiceNew) ListContexts(ctx context.Context, req *serviceInterfaces.ListContextsRequest) (*serviceInterfaces.ListContextsResponse, error) {
	// 验证请求
	if req.UserID == "" {
		return nil, errors.ContextFactory.ValidationError("用户ID不能为空", nil).
			WithOperation("ListContexts")
	}

	// TODO: 实现列出上下文逻辑
	contexts := []*serviceInterfaces.Context{
		{
			ID:          "ctx_1",
			Name:        "上下文1",
			Type:        "conversation",
			UserID:      req.UserID,
			Description: "第一个上下文",
			CreatedAt:   time.Now().Add(-2 * time.Hour),
			UpdatedAt:   time.Now().Add(-time.Hour),
		},
		{
			ID:          "ctx_2",
			Name:        "上下文2",
			Type:        "document",
			UserID:      req.UserID,
			Description: "第二个上下文",
			CreatedAt:   time.Now().Add(-time.Hour),
			UpdatedAt:   time.Now(),
		},
	}

	return &serviceInterfaces.ListContextsResponse{
		Contexts: contexts,
		Total:    len(contexts),
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// AddMessage 添加消息
func (s *ContextServiceNew) AddMessage(ctx context.Context, req *serviceInterfaces.AddMessageRequest) (*serviceInterfaces.AddMessageResponse, error) {
	// 验证请求
	if err := s.validateAddMessageRequest(req); err != nil {
		return nil, errors.ContextFactory.ValidationError("请求验证失败", err).
			WithOperation("AddMessage").
			WithUserID(req.UserID).
			WithMetadata(map[string]interface{}{
				"context_id": req.ContextID,
				"role": req.Role,
			})
	}

	// TODO: 实现添加消息逻辑
	return &serviceInterfaces.AddMessageResponse{
		Message: &serviceInterfaces.Message{
			ID:        fmt.Sprintf("msg_%d", time.Now().UnixNano()),
			ContextID: req.ContextID,
			Role:      req.Role,
			Content:   req.Content,
			UserID:    req.UserID,
			Metadata:  req.Metadata,
			CreatedAt: time.Now(),
		},
	}, nil
}

// GetMessages 获取消息
func (s *ContextServiceNew) GetMessages(ctx context.Context, req *serviceInterfaces.GetMessagesRequest) (*serviceInterfaces.GetMessagesResponse, error) {
	// 验证请求
	if req.ContextID == "" {
		return nil, errors.ContextFactory.ValidationError("上下文ID不能为空", nil).
			WithOperation("GetMessages").
			WithUserID(req.UserID)
	}

	// TODO: 实现获取消息逻辑
	messages := []*serviceInterfaces.Message{
		{
			ID:        "msg_1",
			ContextID: req.ContextID,
			Role:      "user",
			Content:   "用户消息1",
			UserID:    req.UserID,
			CreatedAt: time.Now().Add(-time.Hour),
		},
		{
			ID:        "msg_2",
			ContextID: req.ContextID,
			Role:      "assistant",
			Content:   "助手回复1",
			UserID:    req.UserID,
			CreatedAt: time.Now().Add(-30 * time.Minute),
		},
	}

	return &serviceInterfaces.GetMessagesResponse{
		Messages: messages,
		Total:    len(messages),
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// ClearMessages 清空消息
func (s *ContextServiceNew) ClearMessages(ctx context.Context, req *serviceInterfaces.ClearMessagesRequest) (*serviceInterfaces.ClearMessagesResponse, error) {
	// 验证请求
	if req.ContextID == "" {
		return nil, errors.ContextFactory.ValidationError("上下文ID不能为空", nil).
			WithOperation("ClearMessages").
			WithUserID(req.UserID)
	}

	// TODO: 实现清空消息逻辑
	return &serviceInterfaces.ClearMessagesResponse{
		Success: true,
		Message: "消息已清空",
		Cleared: 10, // 示例数量
	}, nil
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