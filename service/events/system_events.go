package events

import (
	"context"
	"log"
	"time"

	"Qingyu_backend/service/base"
)

// ============ 系统相关事件 ============

// 系统事件类型常量
const (
	// 配置事件
	EventConfigChanged = "config.changed"

	// 权限事件
	EventPermissionChanged = "permission.changed"
	EventRoleChanged       = "role.changed"

	// 审核事件
	EventReviewSubmitted = "review.submitted"
	EventReviewApproved  = "review.approved"
	EventReviewRejected  = "review.rejected"

	// 内容事件
	EventContentPublished   = "content.published"
	EventContentUnpublished = "content.unpublished"

	// 安全事件
	EventSecurityAlert     = "security.alert"
	EventSuspiciousActivity = "security.suspicious_activity"
	EventAccountLocked     = "security.account_locked"

	// 系统事件
	EventSystemMaintenance = "system.maintenance"
	EventSystemError       = "system.error"
)

// SystemEventData 系统事件数据
type SystemEventData struct {
	OperatorID string                 `json:"operator_id,omitempty"`
	TargetType string                 `json:"target_type"`
	TargetID   string                 `json:"target_id,omitempty"`
	Action     string                 `json:"action"`
	Time       time.Time              `json:"time"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

// ============ 配置事件 ============

// ConfigEventData 配置事件数据
type ConfigEventData struct {
	SystemEventData
	ConfigKey   string                 `json:"config_key"`
	OldValue    interface{}            `json:"old_value,omitempty"`
	NewValue    interface{}            `json:"new_value,omitempty"`
	Reason      string                 `json:"reason,omitempty"`
}

// NewConfigChangedEvent 创建配置变更事件
func NewConfigChangedEvent(operatorID, configKey string, oldValue, newValue interface{}, reason string) base.Event {
	return &base.BaseEvent{
		EventType: EventConfigChanged,
		EventData: ConfigEventData{
			SystemEventData: SystemEventData{
				OperatorID: operatorID,
				TargetType: "config",
				TargetID:   configKey,
				Action:     "changed",
				Time:       time.Now(),
			},
			ConfigKey: configKey,
			OldValue:  oldValue,
			NewValue:  newValue,
			Reason:    reason,
		},
		Timestamp: time.Now(),
		Source:    "System",
	}
}

// ============ 权限事件 ============

// PermissionEventData 权限事件数据
type PermissionEventData struct {
	SystemEventData
	UserID       string                 `json:"user_id"`
	Permission   string                 `json:"permission"`
	OldRole      string                 `json:"old_role,omitempty"`
	NewRole      string                 `json:"new_role,omitempty"`
	Granted      bool                   `json:"granted,omitempty"`
	Reason       string                 `json:"reason,omitempty"`
}

// NewPermissionChangedEvent 创建权限变更事件
func NewPermissionChangedEvent(operatorID, userID, permission string, granted bool, reason string) base.Event {
	return &base.BaseEvent{
		EventType: EventPermissionChanged,
		EventData: PermissionEventData{
			SystemEventData: SystemEventData{
				OperatorID: operatorID,
				TargetType:  "permission",
				Action:      "changed",
				Time:        time.Now(),
			},
			UserID:     userID,
			Permission: permission,
			Granted:    granted,
			Reason:     reason,
		},
		Timestamp: time.Now(),
		Source:    "System",
	}
}

// NewRoleChangedEvent 创建角色变更事件
func NewRoleChangedEvent(operatorID, userID, oldRole, newRole string, reason string) base.Event {
	return &base.BaseEvent{
		EventType: EventRoleChanged,
		EventData: PermissionEventData{
			SystemEventData: SystemEventData{
				OperatorID: operatorID,
				TargetType:  "role",
				TargetID:    userID,
				Action:      "changed",
				Time:        time.Now(),
			},
			UserID:   userID,
			OldRole:  oldRole,
			NewRole:  newRole,
			Reason:   reason,
		},
		Timestamp: time.Now(),
		Source:    "System",
	}
}

// ============ 审核事件 ============

// ReviewEventData 审核事件数据
type ReviewEventData struct {
	SystemEventData
	ReviewID     string                 `json:"review_id"`
	ContentType  string                 `json:"content_type"`   // book/chapter/comment
	ContentID    string                 `json:"content_id"`
	SubmitterID  string                 `json:"submitter_id"`
	ReviewerID   string                 `json:"reviewer_id,omitempty"`
	Status       string                 `json:"status"`         // pending/approved/rejected
	Reason       string                 `json:"reason,omitempty"`
	ReviewedAt   time.Time              `json:"reviewed_at,omitempty"`
}

// NewReviewSubmittedEvent 创建审核提交事件
func NewReviewSubmittedEvent(reviewID, contentType, contentID, submitterID string) base.Event {
	return &base.BaseEvent{
		EventType: EventReviewSubmitted,
		EventData: ReviewEventData{
			SystemEventData: SystemEventData{
				TargetType: "review",
				TargetID:   reviewID,
				Action:     "submitted",
				Time:       time.Now(),
			},
			ReviewID:    reviewID,
			ContentType: contentType,
			ContentID:   contentID,
			SubmitterID: submitterID,
			Status:      "pending",
		},
		Timestamp: time.Now(),
		Source:    "ModerationService",
	}
}

// NewReviewApprovedEvent 创建审核批准事件
func NewReviewApprovedEvent(reviewID, contentType, contentID, reviewerID string) base.Event {
	return &base.BaseEvent{
		EventType: EventReviewApproved,
		EventData: ReviewEventData{
			SystemEventData: SystemEventData{
				TargetType: "review",
				TargetID:   reviewID,
				Action:     "approved",
				Time:       time.Now(),
			},
			ReviewID:    reviewID,
			ContentType: contentType,
			ContentID:   contentID,
			ReviewerID:  reviewerID,
			Status:      "approved",
			ReviewedAt:  time.Now(),
		},
		Timestamp: time.Now(),
		Source:    "ModerationService",
	}
}

// NewReviewRejectedEvent 创建审核拒绝事件
func NewReviewRejectedEvent(reviewID, contentType, contentID, reviewerID, reason string) base.Event {
	return &base.BaseEvent{
		EventType: EventReviewRejected,
		EventData: ReviewEventData{
			SystemEventData: SystemEventData{
				TargetType: "review",
				TargetID:   reviewID,
				Action:     "rejected",
				Time:       time.Now(),
			},
			ReviewID:    reviewID,
			ContentType: contentType,
			ContentID:   contentID,
			ReviewerID:  reviewerID,
			Status:      "rejected",
			Reason:      reason,
			ReviewedAt:  time.Now(),
		},
		Timestamp: time.Now(),
		Source:    "ModerationService",
	}
}

// ============ 内容事件 ============

// ContentEventData 内容事件数据
type ContentEventData struct {
	SystemEventData
	ContentType  string                 `json:"content_type"`  // book/chapter/comment
	ContentID    string                 `json:"content_id"`
	AuthorID     string                 `json:"author_id,omitempty"`
	Title        string                 `json:"title,omitempty"`
	PublishTime  time.Time              `json:"publish_time,omitempty"`
	Reason       string                 `json:"reason,omitempty"`
}

// NewContentPublishedEvent 创建内容发布事件
func NewContentPublishedEvent(contentType, contentID, authorID, title string) base.Event {
	return &base.BaseEvent{
		EventType: EventContentPublished,
		EventData: ContentEventData{
			SystemEventData: SystemEventData{
				TargetType: contentType,
				TargetID:   contentID,
				Action:     "published",
				Time:       time.Now(),
			},
			ContentType: contentType,
			ContentID:   contentID,
			AuthorID:    authorID,
			Title:       title,
			PublishTime: time.Now(),
		},
		Timestamp: time.Now(),
		Source:    "ModerationService",
	}
}

// NewContentUnpublishedEvent 创建内容下架事件
func NewContentUnpublishedEvent(contentType, contentID, reason string) base.Event {
	return &base.BaseEvent{
		EventType: EventContentUnpublished,
		EventData: ContentEventData{
			SystemEventData: SystemEventData{
				TargetType: contentType,
				TargetID:   contentID,
				Action:     "unpublished",
				Time:       time.Now(),
			},
			ContentType: contentType,
			ContentID:   contentID,
			Reason:      reason,
		},
		Timestamp: time.Now(),
		Source:    "ModerationService",
	}
}

// ============ 安全事件 ============

// SecurityEventData 安全事件数据
type SecurityEventData struct {
	SystemEventData
	IPAddress    string                 `json:"ip_address,omitempty"`
	UserAgent    string                 `json:"user_agent,omitempty"`
	EventType    string                 `json:"event_type"`
	Severity     string                 `json:"severity"`      // low/medium/high/critical
	Description  string                 `json:"description"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// NewSecurityAlertEvent 创建安全告警事件
func NewSecurityAlertEvent(eventType, severity, description string, metadata map[string]interface{}) base.Event {
	return &base.BaseEvent{
		EventType: EventSecurityAlert,
		EventData: SecurityEventData{
			SystemEventData: SystemEventData{
				TargetType: "security",
				Action:     "alert",
				Time:       time.Now(),
				Metadata:   metadata,
			},
			EventType:   eventType,
			Severity:    severity,
			Description: description,
		},
		Timestamp: time.Now(),
		Source:    "SecurityService",
	}
}

// NewSuspiciousActivityEvent 创建可疑活动事件
func NewSuspiciousActivityEvent(userID, ipAddress, userAgent, activityType string) base.Event {
	return &base.BaseEvent{
		EventType: EventSuspiciousActivity,
		EventData: SecurityEventData{
			SystemEventData: SystemEventData{
				TargetType: "security",
				TargetID:   userID,
				Action:     "suspicious_activity",
				Time:       time.Now(),
			},
			EventType:   activityType,
			IPAddress:   ipAddress,
			UserAgent:   userAgent,
			Severity:    "medium",
			Description: "检测到可疑活动",
		},
		Timestamp: time.Now(),
		Source:    "SecurityService",
	}
}

// NewAccountLockedEvent 创建账户锁定事件
func NewAccountLockedEvent(userID, reason string) base.Event {
	return &base.BaseEvent{
		EventType: EventAccountLocked,
		EventData: SecurityEventData{
			SystemEventData: SystemEventData{
				TargetType: "account",
				TargetID:   userID,
				Action:     "locked",
				Time:       time.Now(),
			},
			EventType:   "account_locked",
			Severity:    "high",
			Description: reason,
		},
		Timestamp: time.Now(),
		Source:    "SecurityService",
	}
}

// ============ 系统维护事件 ============

// MaintenanceEventData 系统维护事件数据
type MaintenanceEventData struct {
	SystemEventData
	MaintenanceType string                 `json:"maintenance_type"` // scheduled/emergency
	StartTime      time.Time              `json:"start_time"`
	EndTime        time.Time              `json:"end_time"`
	Description    string                 `json:"description"`
	AffectedServices []string             `json:"affected_services,omitempty"`
}

// NewSystemMaintenanceEvent 创建系统维护事件
func NewSystemMaintenanceEvent(maintenanceType string, startTime, endTime time.Time, description string, affectedServices []string) base.Event {
	return &base.BaseEvent{
		EventType: EventSystemMaintenance,
		EventData: MaintenanceEventData{
			SystemEventData: SystemEventData{
				TargetType: "system",
				Action:     "maintenance",
				Time:       time.Now(),
			},
			MaintenanceType:  maintenanceType,
			StartTime:       startTime,
			EndTime:         endTime,
			Description:     description,
			AffectedServices: affectedServices,
		},
		Timestamp: time.Now(),
		Source:    "System",
	}
}

// ============ 事件处理器 ============

// AuditLogHandler 审计日志处理器
// 记录所有系统操作审计日志
type AuditLogHandler struct {
	name string
}

// NewAuditLogHandler 创建审计日志处理器
func NewAuditLogHandler() *AuditLogHandler {
	return &AuditLogHandler{
		name: "AuditLogHandler",
	}
}

// Handle 处理事件
func (h *AuditLogHandler) Handle(ctx context.Context, event base.Event) error {
	switch event.GetEventType() {
	case EventConfigChanged:
		data, _ := event.GetEventData().(ConfigEventData)
		log.Printf("[AuditLog] 配置变更: %s, 操作者: %s, 旧值: %v, 新值: %v, 原因: %s",
			data.ConfigKey, data.OperatorID, data.OldValue, data.NewValue, data.Reason)

	case EventPermissionChanged, EventRoleChanged:
		data, _ := event.GetEventData().(PermissionEventData)
		log.Printf("[AuditLog] 权限/角色变更: 用户 %s, 操作者: %s, 原因: %s",
			data.UserID, data.OperatorID, data.Reason)

	case EventReviewApproved, EventReviewRejected:
		data, _ := event.GetEventData().(ReviewEventData)
		log.Printf("[AuditLog] 审核: %s %s, 审核者: %s, 状态: %s, 原因: %s",
			data.ContentType, data.ContentID, data.ReviewerID, data.Status, data.Reason)
	}

	return nil
}

// GetHandlerName 获取处理器名称
func (h *AuditLogHandler) GetHandlerName() string {
	return h.name
}

// GetSupportedEventTypes 获取支持的事件类型
func (h *AuditLogHandler) GetSupportedEventTypes() []string {
	return []string{
		EventConfigChanged,
		EventPermissionChanged,
		EventRoleChanged,
		EventReviewApproved,
		EventReviewRejected,
		EventContentPublished,
		EventContentUnpublished,
	}
}

// CacheInvalidationHandler 缓存失效处理器
// 处理缓存失效
type CacheInvalidationHandler struct {
	name string
}

// NewCacheInvalidationHandler 创建缓存失效处理器
func NewCacheInvalidationHandler() *CacheInvalidationHandler {
	return &CacheInvalidationHandler{
		name: "CacheInvalidationHandler",
	}
}

// Handle 处理事件
func (h *CacheInvalidationHandler) Handle(ctx context.Context, event base.Event) error {
	switch event.GetEventType() {
	case EventConfigChanged:
		data, _ := event.GetEventData().(ConfigEventData)
		log.Printf("[CacheInvalidation] 清除配置缓存: %s", data.ConfigKey)
		// 清除相关缓存

	case EventPermissionChanged, EventRoleChanged:
		data, _ := event.GetEventData().(PermissionEventData)
		log.Printf("[CacheInvalidation] 清除用户权限缓存: %s", data.UserID)
		// 清除用户权限缓存

	case EventContentPublished, EventContentUnpublished:
		data, _ := event.GetEventData().(ContentEventData)
		log.Printf("[CacheInvalidation] 清除内容缓存: %s:%s", data.ContentType, data.ContentID)
		// 清除内容缓存
	}

	return nil
}

// GetHandlerName 获取处理器名称
func (h *CacheInvalidationHandler) GetHandlerName() string {
	return h.name
}

// GetSupportedEventTypes 获取支持的事件类型
func (h *CacheInvalidationHandler) GetSupportedEventTypes() []string {
	return []string{
		EventConfigChanged,
		EventPermissionChanged,
		EventRoleChanged,
		EventContentPublished,
		EventContentUnpublished,
	}
}

// SecurityAlertHandler 安全告警处理器
// 处理安全告警
type SecurityAlertHandler struct {
	name string
}

// NewSecurityAlertHandler 创建安全告警处理器
func NewSecurityAlertHandler() *SecurityAlertHandler {
	return &SecurityAlertHandler{
		name: "SecurityAlertHandler",
	}
}

// Handle 处理事件
func (h *SecurityAlertHandler) Handle(ctx context.Context, event base.Event) error {
	switch event.GetEventType() {
	case EventSecurityAlert:
		data, _ := event.GetEventData().(SecurityEventData)
		log.Printf("[SecurityAlert] 安全告警: 类型=%s, 严重程度=%s, 描述=%s",
			data.EventType, data.Severity, data.Description)
		// 根据严重程度采取不同措施
		// - high/critical: 立即通知管理员
		// - medium: 记录并监控
		// - low: 仅记录

	case EventSuspiciousActivity:
		data, _ := event.GetEventData().(SecurityEventData)
		log.Printf("[SecurityAlert] 可疑活动: 用户=%s, IP=%s, 类型=%s",
			data.TargetID, data.IPAddress, data.EventType)
		// 标记账户，要求额外验证
		// 发送安全通知给用户

	case EventAccountLocked:
		data, _ := event.GetEventData().(SecurityEventData)
		log.Printf("[SecurityAlert] 账户锁定: 用户=%s, 原因=%s", data.TargetID, data.Description)
		// 锁定账户
		// 发送锁定通知
		// 记录安全事件
	}

	return nil
}

// GetHandlerName 获取处理器名称
func (h *SecurityAlertHandler) GetHandlerName() string {
	return h.name
}

// GetSupportedEventTypes 获取支持的事件类型
func (h *SecurityAlertHandler) GetSupportedEventTypes() []string {
	return []string{
		EventSecurityAlert,
		EventSuspiciousActivity,
		EventAccountLocked,
	}
}

// ContentModerationNotificationHandler 内容审核通知处理器
// 发送内容审核通知
type ContentModerationNotificationHandler struct {
	name string
}

// NewContentModerationNotificationHandler 创建内容审核通知处理器
func NewContentModerationNotificationHandler() *ContentModerationNotificationHandler {
	return &ContentModerationNotificationHandler{
		name: "ContentModerationNotificationHandler",
	}
}

// Handle 处理事件
func (h *ContentModerationNotificationHandler) Handle(ctx context.Context, event base.Event) error {
	switch event.GetEventType() {
	case EventReviewApproved:
		data, _ := event.GetEventData().(ReviewEventData)
		log.Printf("[ModerationNotification] 审核通过通知: %s %s，发送给用户 %s",
			data.ContentType, data.ContentID, data.SubmitterID)
		// 发送审核通过通知

	case EventReviewRejected:
		data, _ := event.GetEventData().(ReviewEventData)
		log.Printf("[ModerationNotification] 审核拒绝通知: %s %s，原因: %s，发送给用户 %s",
			data.ContentType, data.ContentID, data.Reason, data.SubmitterID)
		// 发送审核拒绝通知，包含拒绝原因

	case EventContentUnpublished:
		data, _ := event.GetEventData().(ContentEventData)
		log.Printf("[ModerationNotification] 内容下架通知: %s %s，原因: %s",
			data.ContentType, data.ContentID, data.Reason)
		// 发送内容下架通知
	}

	return nil
}

// GetHandlerName 获取处理器名称
func (h *ContentModerationNotificationHandler) GetHandlerName() string {
	return h.name
}

// GetSupportedEventTypes 获取支持的事件类型
func (h *ContentModerationNotificationHandler) GetSupportedEventTypes() []string {
	return []string{
		EventReviewApproved,
		EventReviewRejected,
		EventContentUnpublished,
	}
}
