package events

import (
	"context"
	"log"
	"time"

	"Qingyu_backend/service/base"
)

// ============ 社交相关事件 ============

// 社交事件类型常量
const (
	// 点赞事件
	EventLikeAdded   = "like.added"
	EventLikeRemoved = "like.removed"

	// 评论事件
	EventCommentCreated = "comment.created"
	EventCommentUpdated = "comment.updated"
	EventCommentDeleted = "comment.deleted"
	EventCommentReplied = "comment.replied"

	// 收藏事件
	EventCollectionAdded   = "collection.added"
	EventCollectionRemoved = "collection.removed"

	// 关注事件
	EventFollowAdded   = "follow.added"
	EventFollowRemoved = "follow.removed"
)

// SocialEventData 社交事件数据
type SocialEventData struct {
	UserID     string                 `json:"user_id"`
	TargetType string                 `json:"target_type"` // book/user/comment
	TargetID   string                 `json:"target_id"`
	Action     string                 `json:"action"`
	Time       time.Time              `json:"time"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

// ============ 点赞事件 ============

// NewLikeAddedEvent 创建点赞事件
func NewLikeAddedEvent(userID, targetType, targetID string) base.Event {
	return &base.BaseEvent{
		EventType: EventLikeAdded,
		EventData: SocialEventData{
			UserID:     userID,
			TargetType: targetType,
			TargetID:   targetID,
			Action:     "liked",
			Time:       time.Now(),
		},
		Timestamp: time.Now(),
		Source:    "SocialService",
	}
}

// NewLikeRemovedEvent 创建取消点赞事件
func NewLikeRemovedEvent(userID, targetType, targetID string) base.Event {
	return &base.BaseEvent{
		EventType: EventLikeRemoved,
		EventData: SocialEventData{
			UserID:     userID,
			TargetType: targetType,
			TargetID:   targetID,
			Action:     "unliked",
			Time:       time.Now(),
		},
		Timestamp: time.Now(),
		Source:    "SocialService",
	}
}

// ============ 评论事件 ============

// CommentEventData 评论事件数据（扩展的社交事件数据）
type CommentEventData struct {
	SocialEventData
	CommentID   string `json:"comment_id"`
	Content     string `json:"content"`
	ParentID    string `json:"parent_id,omitempty"`    // 父评论ID（用于回复）
	ReplyToID   string `json:"reply_to_id,omitempty"`  // 回复的评论ID
	ReplyToUserID string `json:"reply_to_user_id,omitempty"` // 回复的用户ID
}

// NewCommentCreatedEvent 创建评论事件
func NewCommentCreatedEvent(userID, targetType, targetID, commentID, content string) base.Event {
	return &base.BaseEvent{
		EventType: EventCommentCreated,
		EventData: CommentEventData{
			SocialEventData: SocialEventData{
				UserID:     userID,
				TargetType: targetType,
				TargetID:   targetID,
				Action:     "commented",
				Time:       time.Now(),
			},
			CommentID: commentID,
			Content:   content,
		},
		Timestamp: time.Now(),
		Source:    "SocialService",
	}
}

// NewCommentRepliedEvent 创建评论回复事件
func NewCommentRepliedEvent(userID, targetType, targetID, commentID, parentID, replyToID, replyToUserID, content string) base.Event {
	return &base.BaseEvent{
		EventType: EventCommentReplied,
		EventData: CommentEventData{
			SocialEventData: SocialEventData{
				UserID:     userID,
				TargetType: targetType,
				TargetID:   targetID,
				Action:     "replied",
				Time:       time.Now(),
			},
			CommentID:     commentID,
			ParentID:      parentID,
			ReplyToID:     replyToID,
			ReplyToUserID: replyToUserID,
			Content:       content,
		},
		Timestamp: time.Now(),
		Source:    "SocialService",
	}
}

// NewCommentDeletedEvent 创建评论删除事件
func NewCommentDeletedEvent(userID, targetType, targetID, commentID string) base.Event {
	return &base.BaseEvent{
		EventType: EventCommentDeleted,
		EventData: CommentEventData{
			SocialEventData: SocialEventData{
				UserID:     userID,
				TargetType: targetType,
				TargetID:   targetID,
				Action:     "comment_deleted",
				Time:       time.Now(),
			},
			CommentID: commentID,
		},
		Timestamp: time.Now(),
		Source:    "SocialService",
	}
}

// ============ 收藏事件 ============

// NewCollectionAddedEvent 创建收藏事件
func NewCollectionAddedEvent(userID, targetType, targetID string) base.Event {
	return &base.BaseEvent{
		EventType: EventCollectionAdded,
		EventData: SocialEventData{
			UserID:     userID,
			TargetType: targetType,
			TargetID:   targetID,
			Action:     "collected",
			Time:       time.Now(),
		},
		Timestamp: time.Now(),
		Source:    "SocialService",
	}
}

// NewCollectionRemovedEvent 创建取消收藏事件
func NewCollectionRemovedEvent(userID, targetType, targetID string) base.Event {
	return &base.BaseEvent{
		EventType: EventCollectionRemoved,
		EventData: SocialEventData{
			UserID:     userID,
			TargetType: targetType,
			TargetID:   targetID,
			Action:     "uncollected",
			Time:       time.Now(),
		},
		Timestamp: time.Now(),
		Source:    "SocialService",
	}
}

// ============ 关注事件 ============

// FollowEventData 关注事件数据
type FollowEventData struct {
	FollowerID  string    `json:"follower_id"`  // 关注者ID
	FolloweeID  string    `json:"followee_id"`  // 被关注者ID
	Action      string    `json:"action"`
	Time        time.Time `json:"time"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// NewFollowAddedEvent 创建关注事件
func NewFollowAddedEvent(followerID, followeeID string) base.Event {
	return &base.BaseEvent{
		EventType: EventFollowAdded,
		EventData: FollowEventData{
			FollowerID: followerID,
			FolloweeID: followeeID,
			Action:     "followed",
			Time:       time.Now(),
		},
		Timestamp: time.Now(),
		Source:    "SocialService",
	}
}

// NewFollowRemovedEvent 创建取消关注事件
func NewFollowRemovedEvent(followerID, followeeID string) base.Event {
	return &base.BaseEvent{
		EventType: EventFollowRemoved,
		EventData: FollowEventData{
			FollowerID: followerID,
			FolloweeID: followeeID,
			Action:     "unfollowed",
			Time:       time.Now(),
		},
		Timestamp: time.Now(),
		Source:    "SocialService",
	}
}

// ============ 事件处理器 ============

// SocialNotificationHandler 社交通知处理器
// 处理社交互动通知
type SocialNotificationHandler struct {
	name string
}

// NewSocialNotificationHandler 创建社交通知处理器
func NewSocialNotificationHandler() *SocialNotificationHandler {
	return &SocialNotificationHandler{
		name: "SocialNotificationHandler",
	}
}

// Handle 处理事件
func (h *SocialNotificationHandler) Handle(ctx context.Context, event base.Event) error {
	switch event.GetEventType() {
	case EventLikeAdded:
		data, _ := event.GetEventData().(SocialEventData)
		log.Printf("[SocialNotification] 用户 %s 点赞了 %s:%s", data.UserID, data.TargetType, data.TargetID)
		// 发送点赞通知给目标用户/作者

	case EventCommentCreated, EventCommentReplied:
		data, _ := event.GetEventData().(CommentEventData)
		log.Printf("[SocialNotification] 用户 %s 评论了 %s:%s", data.UserID, data.TargetType, data.TargetID)
		// 发送评论通知

	case EventFollowAdded:
		data, _ := event.GetEventData().(FollowEventData)
		log.Printf("[SocialNotification] 用户 %s 关注了用户 %s", data.FollowerID, data.FolloweeID)
		// 发送关注通知
	}

	return nil
}

// GetHandlerName 获取处理器名称
func (h *SocialNotificationHandler) GetHandlerName() string {
	return h.name
}

// GetSupportedEventTypes 获取支持的事件类型
func (h *SocialNotificationHandler) GetSupportedEventTypes() []string {
	return []string{
		EventLikeAdded,
		EventCommentCreated,
		EventCommentReplied,
		EventFollowAdded,
	}
}

// SocialStatisticsHandler 社交统计处理器
// 更新社交互动统计
type SocialStatisticsHandler struct {
	name string
}

// NewSocialStatisticsHandler 创建社交统计处理器
func NewSocialStatisticsHandler() *SocialStatisticsHandler {
	return &SocialStatisticsHandler{
		name: "SocialStatisticsHandler",
	}
}

// Handle 处理事件
func (h *SocialStatisticsHandler) Handle(ctx context.Context, event base.Event) error {
	switch event.GetEventType() {
	case EventLikeAdded, EventLikeRemoved:
		data, _ := event.GetEventData().(SocialEventData)
		log.Printf("[SocialStatistics] 更新 %s:%s 的点赞统计", data.TargetType, data.TargetID)
		// 更新点赞计数

	case EventCommentCreated, EventCommentDeleted:
		data, _ := event.GetEventData().(CommentEventData)
		log.Printf("[SocialStatistics] 更新 %s:%s 的评论统计", data.TargetType, data.TargetID)
		// 更新评论计数

	case EventFollowAdded, EventFollowRemoved:
		data, _ := event.GetEventData().(FollowEventData)
		log.Printf("[SocialStatistics] 更新用户 %s 的粉丝统计", data.FolloweeID)
		// 更新粉丝计数
	}

	return nil
}

// GetHandlerName 获取处理器名称
func (h *SocialStatisticsHandler) GetHandlerName() string {
	return h.name
}

// GetSupportedEventTypes 获取支持的事件类型
func (h *SocialStatisticsHandler) GetSupportedEventTypes() []string {
	return []string{
		EventLikeAdded,
		EventLikeRemoved,
		EventCommentCreated,
		EventCommentDeleted,
		EventFollowAdded,
		EventFollowRemoved,
	}
}

// SocialAchievementHandler 社交成就处理器
// 检查和颁发社交成就
type SocialAchievementHandler struct {
	name string
}

// NewSocialAchievementHandler 创建社交成就处理器
func NewSocialAchievementHandler() *SocialAchievementHandler {
	return &SocialAchievementHandler{
		name: "SocialAchievementHandler",
	}
}

// Handle 处理事件
func (h *SocialAchievementHandler) Handle(ctx context.Context, event base.Event) error {
	switch event.GetEventType() {
	case EventLikeAdded:
		data, _ := event.GetEventData().(SocialEventData)
		log.Printf("[SocialAchievement] 检查用户 %s 的点赞成就", data.UserID)
		// 检查是否达成"获得100个点赞"等成就

	case EventCommentCreated:
		data, _ := event.GetEventData().(CommentEventData)
		log.Printf("[SocialAchievement] 检查用户 %s 的评论成就", data.UserID)
		// 检查是否达成"发表10条评论"等成就

	case EventFollowAdded:
		data, _ := event.GetEventData().(FollowEventData)
		log.Printf("[SocialAchievement] 检查用户 %s 和 %s 的关注成就", data.FollowerID, data.FolloweeID)
		// 检查是否达成"关注10位作者"等成就
	}

	return nil
}

// GetHandlerName 获取处理器名称
func (h *SocialAchievementHandler) GetHandlerName() string {
	return h.name
}

// GetSupportedEventTypes 获取支持的事件类型
func (h *SocialAchievementHandler) GetSupportedEventTypes() []string {
	return []string{
		EventLikeAdded,
		EventCommentCreated,
		EventFollowAdded,
	}
}
