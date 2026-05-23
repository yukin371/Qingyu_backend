package events

import (
	"context"
	"fmt"
	"time"

	"Qingyu_backend/service/base"
)

// ============ 用户相关事件 ============

// UserEvent 事件类型常量
const (
	EventTypeUserRegistered = "user.registered"
	EventTypeUserLoggedIn   = "user.logged_in"
	EventTypeUserLoggedOut  = "user.logged_out"
	EventTypeUserUpdated    = "user.updated"
	EventTypeUserDeleted    = "user.deleted"
)

// UserEventData 用户事件数据
type UserEventData struct {
	UserID   string                 `json:"user_id"`
	Username string                 `json:"username"`
	Email    string                 `json:"email"`
	Action   string                 `json:"action"`
	Time     time.Time              `json:"time"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// NewUserRegisteredEvent 创建用户注册事件
func NewUserRegisteredEvent(userID, username, email string) base.Event {
	return &base.BaseEvent{
		EventType: EventTypeUserRegistered,
		EventData: UserEventData{
			UserID:   userID,
			Username: username,
			Email:    email,
			Action:   "registered",
			Time:     time.Now(),
		},
		Timestamp: time.Now(),
		Source:    "UserService",
	}
}

// NewUserLoggedInEvent 创建用户登录事件
func NewUserLoggedInEvent(userID, username string) base.Event {
	return &base.BaseEvent{
		EventType: EventTypeUserLoggedIn,
		EventData: UserEventData{
			UserID:   userID,
			Username: username,
			Action:   "logged_in",
			Time:     time.Now(),
		},
		Timestamp: time.Now(),
		Source:    "UserService",
	}
}

// ============ 事件处理器 ============

// WelcomeEmailHandler 欢迎邮件处理器
// 当用户注册时发送欢迎邮件
type WelcomeEmailHandler struct {
	name string
}

// NewWelcomeEmailHandler 创建欢迎邮件处理器
func NewWelcomeEmailHandler() *WelcomeEmailHandler {
	return &WelcomeEmailHandler{
		name: "WelcomeEmailHandler",
	}
}

// Handle 处理事件
func (h *WelcomeEmailHandler) Handle(ctx context.Context, event base.Event) error {
	// 解析事件数据
	data, ok := event.GetEventData().(UserEventData)
	if !ok {
		return fmt.Errorf("事件数据类型错误")
	}

	// 发送欢迎邮件（这里只是模拟）
	logEventRuntime("info", "发送欢迎邮件", map[string]interface{}{
		"handler":  h.name,
		"user_id":  data.UserID,
		"username": data.Username,
		"email":    data.Email,
	})

	// 实际项目中这里应该调用邮件服务
	// emailService.SendWelcomeEmail(data.Email, data.Username)

	return nil
}

// GetHandlerName 获取处理器名称
func (h *WelcomeEmailHandler) GetHandlerName() string {
	return h.name
}

// GetSupportedEventTypes 获取支持的事件类型
func (h *WelcomeEmailHandler) GetSupportedEventTypes() []string {
	return []string{EventTypeUserRegistered}
}

// UserActivityLogHandler 用户活动日志处理器
// 记录所有用户活动
type UserActivityLogHandler struct {
	name string
}

// NewUserActivityLogHandler 创建用户活动日志处理器
func NewUserActivityLogHandler() *UserActivityLogHandler {
	return &UserActivityLogHandler{
		name: "UserActivityLogHandler",
	}
}

// Handle 处理事件
func (h *UserActivityLogHandler) Handle(ctx context.Context, event base.Event) error {
	// 解析事件数据
	data, ok := event.GetEventData().(UserEventData)
	if !ok {
		return fmt.Errorf("事件数据类型错误")
	}

	// 记录日志
	logEventRuntime("info", "记录用户活动", map[string]interface{}{
		"handler":    h.name,
		"event_type": event.GetEventType(),
		"user_id":    data.UserID,
		"username":   data.Username,
		"action":     data.Action,
		"event_time": data.Time,
	})

	// 实际项目中这里应该将日志写入数据库或日志系统
	// activityLogRepo.Create(...)

	return nil
}

// GetHandlerName 获取处理器名称
func (h *UserActivityLogHandler) GetHandlerName() string {
	return h.name
}

// GetSupportedEventTypes 获取支持的事件类型
func (h *UserActivityLogHandler) GetSupportedEventTypes() []string {
	return []string{
		EventTypeUserRegistered,
		EventTypeUserLoggedIn,
		EventTypeUserLoggedOut,
		EventTypeUserUpdated,
		EventTypeUserDeleted,
	}
}

// UserStatisticsHandler 用户统计处理器
// 更新用户统计信息
type UserStatisticsHandler struct {
	name string
}

// NewUserStatisticsHandler 创建用户统计处理器
func NewUserStatisticsHandler() *UserStatisticsHandler {
	return &UserStatisticsHandler{
		name: "UserStatisticsHandler",
	}
}

// Handle 处理事件
func (h *UserStatisticsHandler) Handle(ctx context.Context, event base.Event) error {
	// 解析事件数据
	data, ok := event.GetEventData().(UserEventData)
	if !ok {
		return fmt.Errorf("事件数据类型错误")
	}

	// 更新统计信息
	switch event.GetEventType() {
	case EventTypeUserRegistered:
		logEventRuntime("info", "更新用户统计", map[string]interface{}{
			"handler":    h.name,
			"event_type": event.GetEventType(),
			"user_id":    data.UserID,
			"username":   data.Username,
			"action":     "increment_total_users",
		})
		// 实际项目中这里应该更新统计数据
		// statisticsRepo.IncrementTotalUsers()

	case EventTypeUserLoggedIn:
		logEventRuntime("info", "更新用户统计", map[string]interface{}{
			"handler":    h.name,
			"event_type": event.GetEventType(),
			"user_id":    data.UserID,
			"username":   data.Username,
			"action":     "increment_active_users",
		})
		// 实际项目中这里应该更新活跃用户统计
		// statisticsRepo.IncrementActiveUsers()

	case EventTypeUserDeleted:
		logEventRuntime("info", "更新用户统计", map[string]interface{}{
			"handler":    h.name,
			"event_type": event.GetEventType(),
			"user_id":    data.UserID,
			"username":   data.Username,
			"action":     "decrement_total_users",
		})
		// 实际项目中这里应该更新统计数据
		// statisticsRepo.DecrementTotalUsers()
	}

	return nil
}

// GetHandlerName 获取处理器名称
func (h *UserStatisticsHandler) GetHandlerName() string {
	return h.name
}

// GetSupportedEventTypes 获取支持的事件类型
func (h *UserStatisticsHandler) GetSupportedEventTypes() []string {
	return []string{
		EventTypeUserRegistered,
		EventTypeUserLoggedIn,
		EventTypeUserDeleted,
	}
}
