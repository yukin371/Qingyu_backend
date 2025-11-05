package base

import (
	"context"
	"fmt"
	"time"

	"Qingyu_backend/repository/interfaces"
	baseInterface "Qingyu_backend/service/interfaces/base"
)

// 类型别名：方便使用
type (
	BaseService  = baseInterface.BaseService
	Event        = baseInterface.Event
	EventHandler = baseInterface.EventHandler
	EventBus     = baseInterface.EventBus
)

// ServiceContainer 服务容器
// 用于依赖注入和服务管理
type ServiceContainer struct {
	repositoryFactory interfaces.RepositoryFactory
	services          map[string]BaseService
	initialized       bool
}

// NewServiceContainer 创建服务容器
func NewServiceContainer(repositoryFactory interfaces.RepositoryFactory) *ServiceContainer {
	return &ServiceContainer{
		repositoryFactory: repositoryFactory,
		services:          make(map[string]BaseService),
		initialized:       false,
	}
}

// RegisterService 注册服务
func (c *ServiceContainer) RegisterService(name string, service BaseService) error {
	if c.services[name] != nil {
		return fmt.Errorf("服务 %s 已存在", name)
	}

	c.services[name] = service
	return nil
}

// GetService 获取服务
func (c *ServiceContainer) GetService(name string) (BaseService, error) {
	service, exists := c.services[name]
	if !exists {
		return nil, fmt.Errorf("服务 %s 不存在", name)
	}

	return service, nil
}

// Initialize 初始化所有服务
func (c *ServiceContainer) Initialize(ctx context.Context) error {
	if c.initialized {
		return nil
	}

	for name, service := range c.services {
		if err := service.Initialize(ctx); err != nil {
			return fmt.Errorf("初始化服务 %s 失败: %w", name, err)
		}
	}

	c.initialized = true
	return nil
}

// Health 检查所有服务健康状态
func (c *ServiceContainer) Health(ctx context.Context) error {
	// 检查Repository工厂健康状态
	if err := c.repositoryFactory.Health(ctx); err != nil {
		return fmt.Errorf("Repository工厂健康检查失败: %w", err)
	}

	// 检查所有服务健康状态
	for name, service := range c.services {
		if err := service.Health(ctx); err != nil {
			return fmt.Errorf("服务 %s 健康检查失败: %w", name, err)
		}
	}

	return nil
}

// Close 关闭所有服务
func (c *ServiceContainer) Close(ctx context.Context) error {
	var lastErr error

	// 关闭所有服务
	for name, service := range c.services {
		if err := service.Close(ctx); err != nil {
			lastErr = fmt.Errorf("关闭服务 %s 失败: %w", name, err)
		}
	}

	// 关闭Repository工厂
	if err := c.repositoryFactory.Close(); err != nil {
		lastErr = fmt.Errorf("关闭Repository工厂失败: %w", err)
	}

	c.initialized = false
	return lastErr
}

// GetRepositoryFactory 获取Repository工厂
func (c *ServiceContainer) GetRepositoryFactory() interfaces.RepositoryFactory {
	return c.repositoryFactory
}

// ValidationRule 验证规则接口
type ValidationRule interface {
	Validate(ctx context.Context, value interface{}) error
	GetFieldName() string
	GetErrorMessage() string
}

// Validator 验证器接口
type Validator interface {
	AddRule(rule ValidationRule) Validator
	Validate(ctx context.Context, data interface{}) error
	ValidateField(ctx context.Context, fieldName string, value interface{}) error
}

// BaseValidator 基础验证器实现
type BaseValidator struct {
	rules map[string][]ValidationRule
}

// NewBaseValidator 创建基础验证器
func NewBaseValidator() *BaseValidator {
	return &BaseValidator{
		rules: make(map[string][]ValidationRule),
	}
}

// AddRule 添加验证规则
func (v *BaseValidator) AddRule(rule ValidationRule) Validator {
	fieldName := rule.GetFieldName()
	v.rules[fieldName] = append(v.rules[fieldName], rule)
	return v
}

// Validate 验证数据
func (v *BaseValidator) Validate(ctx context.Context, data interface{}) error {
	// 这里需要使用反射来验证结构体字段
	// 为了简化，暂时返回nil
	return nil
}

// ValidateField 验证单个字段
func (v *BaseValidator) ValidateField(ctx context.Context, fieldName string, value interface{}) error {
	rules, exists := v.rules[fieldName]
	if !exists {
		return nil
	}

	for _, rule := range rules {
		if err := rule.Validate(ctx, value); err != nil {
			return err
		}
	}

	return nil
}

// RequiredRule 必填验证规则
type RequiredRule struct {
	fieldName    string
	errorMessage string
}

// NewRequiredRule 创建必填验证规则
func NewRequiredRule(fieldName string) *RequiredRule {
	return &RequiredRule{
		fieldName:    fieldName,
		errorMessage: fmt.Sprintf("%s 不能为空", fieldName),
	}
}

// Validate 验证值
func (r *RequiredRule) Validate(ctx context.Context, value interface{}) error {
	if value == nil {
		return fmt.Errorf("%s", r.errorMessage)
	}

	// 检查字符串是否为空
	if str, ok := value.(string); ok && str == "" {
		return fmt.Errorf("%s", r.errorMessage)
	}

	return nil
}

// GetFieldName 获取字段名
func (r *RequiredRule) GetFieldName() string {
	return r.fieldName
}

// GetErrorMessage 获取错误消息
func (r *RequiredRule) GetErrorMessage() string {
	return r.errorMessage
}

// LengthRule 长度验证规则
type LengthRule struct {
	fieldName    string
	minLength    int
	maxLength    int
	errorMessage string
}

// NewLengthRule 创建长度验证规则
func NewLengthRule(fieldName string, minLength, maxLength int) *LengthRule {
	return &LengthRule{
		fieldName:    fieldName,
		minLength:    minLength,
		maxLength:    maxLength,
		errorMessage: fmt.Sprintf("%s 长度必须在 %d 到 %d 之间", fieldName, minLength, maxLength),
	}
}

// Validate 验证值
func (r *LengthRule) Validate(ctx context.Context, value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("%s 必须是字符串类型", r.fieldName)
	}

	length := len(str)
	if length < r.minLength || length > r.maxLength {
		return fmt.Errorf("%s", r.errorMessage)
	}

	return nil
}

// GetFieldName 获取字段名
func (r *LengthRule) GetFieldName() string {
	return r.fieldName
}

// GetErrorMessage 获取错误消息
func (r *LengthRule) GetErrorMessage() string {
	return r.errorMessage
}

// BaseEvent 基础事件实现
type BaseEvent struct {
	EventType string      `json:"event_type"`
	EventData interface{} `json:"event_data"`
	Timestamp time.Time   `json:"timestamp"`
	Source    string      `json:"source"`
}

// GetEventType 获取事件类型
func (e *BaseEvent) GetEventType() string {
	return e.EventType
}

// GetEventData 获取事件数据
func (e *BaseEvent) GetEventData() interface{} {
	return e.EventData
}

// GetTimestamp 获取时间戳
func (e *BaseEvent) GetTimestamp() time.Time {
	return e.Timestamp
}

// GetSource 获取事件源
func (e *BaseEvent) GetSource() string {
	return e.Source
}

// SimpleEventBus 简单事件总线实现
type SimpleEventBus struct {
	handlers map[string][]EventHandler
}

// NewSimpleEventBus 创建简单事件总线
func NewSimpleEventBus() *SimpleEventBus {
	return &SimpleEventBus{
		handlers: make(map[string][]EventHandler),
	}
}

// Subscribe 订阅事件
func (bus *SimpleEventBus) Subscribe(eventType string, handler EventHandler) error {
	bus.handlers[eventType] = append(bus.handlers[eventType], handler)
	return nil
}

// Unsubscribe 取消订阅事件
func (bus *SimpleEventBus) Unsubscribe(eventType string, handlerName string) error {
	handlers := bus.handlers[eventType]
	for i, handler := range handlers {
		if handler.GetHandlerName() == handlerName {
			bus.handlers[eventType] = append(handlers[:i], handlers[i+1:]...)
			break
		}
	}
	return nil
}

// Publish 发布事件（同步）
func (bus *SimpleEventBus) Publish(ctx context.Context, event Event) error {
	handlers := bus.handlers[event.GetEventType()]
	for _, handler := range handlers {
		if err := handler.Handle(ctx, event); err != nil {
			return fmt.Errorf("事件处理器 %s 处理事件失败: %w", handler.GetHandlerName(), err)
		}
	}
	return nil
}

// PublishAsync 发布事件（异步）
func (bus *SimpleEventBus) PublishAsync(ctx context.Context, event Event) error {
	handlers := bus.handlers[event.GetEventType()]
	for _, handler := range handlers {
		go func(h EventHandler) {
			if err := h.Handle(ctx, event); err != nil {
				// 这里可以记录日志或发送到错误处理系统
				fmt.Printf("异步事件处理器 %s 处理事件失败: %v\n", h.GetHandlerName(), err)
			}
		}(handler)
	}
	return nil
}
