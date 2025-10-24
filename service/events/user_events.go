package events

import (
	"context"
	"fmt"
	"log"
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
	UserID   string    `json:"user_id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Action   string    `json:"action"`
	Time     time.Time `json:"time"`
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
	log.Printf("[WelcomeEmailHandler] 发送欢迎邮件给用户: %s (%s)", data.Username, data.Email)
	
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
	log.Printf("[UserActivityLog] 用户活动: %s - 用户: %s, 动作: %s, 时间: %s",
		event.GetEventType(),
		data.Username,
		data.Action,
		data.Time.Format("2006-01-02 15:04:05"))
	
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
		log.Printf("[UserStatistics] 新用户注册: %s, 更新总用户数", data.Username)
		// 实际项目中这里应该更新统计数据
		// statisticsRepo.IncrementTotalUsers()
		
	case EventTypeUserLoggedIn:
		log.Printf("[UserStatistics] 用户登录: %s, 更新活跃用户数", data.Username)
		// 实际项目中这里应该更新活跃用户统计
		// statisticsRepo.IncrementActiveUsers()
		
	case EventTypeUserDeleted:
		log.Printf("[UserStatistics] 用户删除: %s, 更新总用户数", data.Username)
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

