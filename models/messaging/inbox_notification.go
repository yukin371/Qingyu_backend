package messaging

import (
	"time"

	"Qingyu_backend/models/messaging/base"
)

// InboxNotificationType 站内通知类型
type InboxNotificationType string

const (
	InboxNotificationTypeComment      InboxNotificationType = "comment"       // 评论通知
	InboxNotificationTypeLike         InboxNotificationType = "like"          // 点赞通知
	InboxNotificationTypeFollow       InboxNotificationType = "follow"        // 关注通知
	InboxNotificationTypeMention      InboxNotificationType = "mention"       // @提醒
	InboxNotificationTypeSystem       InboxNotificationType = "system"        // 系统通知
	InboxNotificationTypeAnnouncement InboxNotificationType = "announcement"  // 公告通知
	InboxNotificationTypeUpdate       InboxNotificationType = "update"        // 内容更新通知
	InboxNotificationTypeInvite       InboxNotificationType = "invite"        // 邀请通知
)

// InboxNotificationPriority 通知优先级
type InboxNotificationPriority string

const (
	InboxNotificationPriorityLow    InboxNotificationPriority = "low"    // 低优先级
	InboxNotificationPriorityNormal InboxNotificationPriority = "normal" // 普通优先级
	InboxNotificationPriorityHigh   InboxNotificationPriority = "high"   // 高优先级
	InboxNotificationPriorityUrgent InboxNotificationPriority = "urgent" // 紧急
)

// InboxNotification 站内通知模型
type InboxNotification struct {
	base.IdentifiedEntity      `bson:",inline"`
	base.Timestamps            `bson:",inline"`
	base.CommunicationBase     `bson:",inline"`
	base.Expirable             `bson:",inline"`
	base.TargetEntity          `bson:",inline"`
	base.Pinned                `bson:",inline"`

	Type        InboxNotificationType    `bson:"type" json:"type" validate:"required"` // 通知类型
	Priority    InboxNotificationPriority `bson:"priority" json:"priority" validate:"required,oneof=low normal high urgent"` // 优先级
	Title       string              `bson:"title" json:"title" validate:"required,min=1,max=200"`  // 通知标题
	Content     string              `bson:"content" json:"content" validate:"required,min=1,max=1000"` // 通知内容
	ActionURL   string              `bson:"action_url,omitempty" json:"actionUrl,omitempty"`  // 操作链接
	ActionText  string              `bson:"action_text,omitempty" json:"actionText,omitempty"` // 操作按钮文字

	// 发送者快照
	ActorSnapshot *NotificationActorSnapshot `bson:"actor_snapshot,omitempty" json:"actorSnapshot,omitempty"`
}

// NotificationActorSnapshot 通知发起者快照
type NotificationActorSnapshot struct {
	ID       string `bson:"id" json:"id"`
	Username string `bson:"username" json:"username"`
	Avatar   string `bson:"avatar,omitempty" json:"avatar,omitempty"`
}

// InboxNotificationFilter 通知查询过滤器
type InboxNotificationFilter struct {
	UserID     *string                    `json:"userId,omitempty"`
	Type       *InboxNotificationType     `json:"type,omitempty"`
	Priority   *InboxNotificationPriority `json:"priority,omitempty"`
	IsRead     *bool                      `json:"isRead,omitempty"`
	IsPinned   *bool                      `json:"isPinned,omitempty"`
	TargetType *string                    `json:"targetType,omitempty"`
	TargetID   *string                    `json:"targetId,omitempty"`
	StartTime  *time.Time                 `json:"startTime,omitempty"`
	EndTime    *time.Time                 `json:"endTime,omitempty"`
	SortBy     string                     `json:"sortBy,omitempty"`
	SortOrder  string                     `json:"sortOrder,omitempty"`
	Limit      int                        `json:"limit,omitempty"`
	Offset     int                        `json:"offset,omitempty"`
}

// GetConditions 获取查询条件
func (f *InboxNotificationFilter) GetConditions() map[string]interface{} {
	conditions := make(map[string]interface{})

	// 通知只能由接收者查询，必须指定用户ID
	if f.UserID != nil {
		conditions["receiver_id"] = *f.UserID
	}

	if f.Type != nil {
		conditions["type"] = *f.Type
	}
	if f.Priority != nil {
		conditions["priority"] = *f.Priority
	}
	if f.IsRead != nil {
		conditions["is_read"] = *f.IsRead
	}
	if f.IsPinned != nil {
		conditions["is_pinned"] = *f.IsPinned
	}
	if f.TargetType != nil {
		conditions["target_type"] = *f.TargetType
	}
	if f.TargetID != nil {
		conditions["target_id"] = *f.TargetID
	}

	// 时间范围
	if f.StartTime != nil || f.EndTime != nil {
		timeCondition := make(map[string]interface{})
		if f.StartTime != nil {
			timeCondition["$gte"] = *f.StartTime
		}
		if f.EndTime != nil {
			timeCondition["$lte"] = *f.EndTime
		}
		conditions["created_at"] = timeCondition
	}

	// 排除过期通知
	conditions["$or"] = []map[string]interface{}{
		{"expires_at": nil},
		{"expires_at": map[string]interface{}{"$gt": time.Now()}},
	}

	return conditions
}

// GetSort 获取排序
func (f *InboxNotificationFilter) GetSort() map[string]int {
	sort := make(map[string]int)

	sortValue := -1 // 默认降序
	if f.SortOrder == "asc" {
		sortValue = 1
	}

	switch f.SortBy {
	case "created_at":
		sort["created_at"] = sortValue
	case "priority":
		sort["priority"] = sortValue
	case "is_read":
		sort["is_read"] = sortValue
	default:
		// 默认：置顶优先，然后按优先级高到低，最后按创建时间
		sort["is_pinned"] = -1
		sort["priority"] = -1
		sort["created_at"] = -1
	}

	return sort
}

// GetLimit 获取限制
func (f *InboxNotificationFilter) GetLimit() int {
	return f.Limit
}

// GetOffset 获取偏移
func (f *InboxNotificationFilter) GetOffset() int {
	return f.Offset
}

// GetFields 获取字段
func (f *InboxNotificationFilter) GetFields() []string {
	return []string{}
}

// Validate 验证
func (f *InboxNotificationFilter) Validate() error {
	if f.Limit < 0 {
		return &ValidationError{Message: "limit不能为负数"}
	}
	if f.Offset < 0 {
		return &ValidationError{Message: "offset不能为负数"}
	}
	if f.UserID == nil {
		return &ValidationError{Message: "必须指定用户ID"}
	}
	return nil
}

// IsEffective 判断通知是否有效（未过期且未删除）
func (n *InboxNotification) IsEffective() bool {
	if n.IsDeleted() {
		return false
	}
	return !n.IsExpired()
}

// ShouldShow 判断通知是否应该显示给用户
func (n *InboxNotification) ShouldShow() bool {
	return n.IsEffective() && n.ReceiverID != ""
}

// GetPriorityLevel 获取优先级等级（用于排序）
func (n *InboxNotification) GetPriorityLevel() int {
	switch n.Priority {
	case InboxNotificationPriorityUrgent:
		return 4
	case InboxNotificationPriorityHigh:
		return 3
	case InboxNotificationPriorityNormal:
		return 2
	case InboxNotificationPriorityLow:
		return 1
	default:
		return 0
	}
}

// SetAutoExpiration 设置自动过期时间
// 根据优先级设置不同的过期时间
func (n *InboxNotification) SetAutoExpiration() {
	var duration time.Duration

	switch n.Priority {
	case InboxNotificationPriorityUrgent:
		duration = 30 * 24 * time.Hour // 30天
	case InboxNotificationPriorityHigh:
		duration = 14 * 24 * time.Hour // 14天
	case InboxNotificationPriorityNormal:
		duration = 7 * 24 * time.Hour  // 7天
	case InboxNotificationPriorityLow:
		duration = 3 * 24 * time.Hour  // 3天
	default:
		duration = 7 * 24 * time.Hour
	}

	n.SetExpiration(duration)
}

// MarkAsReadAndTouch 标记为已读并更新时间戳
func (n *InboxNotification) MarkAsReadAndTouch() {
	n.MarkAsRead()
	n.Touch()
}

// CreateCommentNotification 创建评论通知的辅助函数
func CreateInboxCommentNotification(receiverID, senderID, username, avatar, targetID, targetType string, content string) *InboxNotification {
	return &InboxNotification{
		CommunicationBase: base.CommunicationBase{
			SenderID:   senderID,
			ReceiverID: receiverID,
		},
		TargetEntity: base.TargetEntity{
			TargetType: targetType,
			TargetID:   targetID,
		},
		Type:     InboxNotificationTypeComment,
		Priority: InboxNotificationPriorityNormal,
		Title:    "新评论通知",
		Content:  content,
		ActorSnapshot: &NotificationActorSnapshot{
			ID:       senderID,
			Username: username,
			Avatar:   avatar,
		},
	}
}

// CreateLikeNotification 创建点赞通知的辅助函数
func CreateInboxLikeNotification(receiverID, senderID, username, avatar, targetID, targetType string) *InboxNotification {
	return &InboxNotification{
		CommunicationBase: base.CommunicationBase{
			SenderID:   senderID,
			ReceiverID: receiverID,
		},
		TargetEntity: base.TargetEntity{
			TargetType: targetType,
			TargetID:   targetID,
		},
		Type:     InboxNotificationTypeLike,
		Priority: InboxNotificationPriorityLow,
		Title:    "收到新的点赞",
		Content:  "赞了你的内容",
		ActorSnapshot: &NotificationActorSnapshot{
			ID:       senderID,
			Username: username,
			Avatar:   avatar,
		},
	}
}

// CreateSystemNotification 创建系统通知的辅助函数
func CreateInboxSystemNotification(receiverID, title, content string, priority InboxNotificationPriority) *InboxNotification {
	return &InboxNotification{
		CommunicationBase: base.CommunicationBase{
			SenderID:   "system", // 系统通知使用固定ID
			ReceiverID: receiverID,
		},
		Type:     InboxNotificationTypeSystem,
		Priority: priority,
		Title:    title,
		Content:  content,
	}
}
