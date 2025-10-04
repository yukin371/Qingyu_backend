package messaging

import (
	"context"
	"time"
)

// MessagingService 消息队列服务接口（对外暴露）
type MessagingService interface {
	// 发布消息
	Publish(ctx context.Context, topic string, message interface{}) error
	PublishDelayed(ctx context.Context, topic string, message interface{}, delay time.Duration) error

	// 订阅消息
	Subscribe(ctx context.Context, topic string, handler MessageHandler) error
	Unsubscribe(ctx context.Context, topic string) error

	// 管理
	CreateTopic(ctx context.Context, topic string) error
	DeleteTopic(ctx context.Context, topic string) error
	ListTopics(ctx context.Context) ([]string, error)

	// 健康检查
	Health(ctx context.Context) error
}

// MessageHandler 消息处理函数
type MessageHandler func(ctx context.Context, message *Message) error

// ============ 数据结构 ============

// Message 消息
type Message struct {
	ID        string                 `json:"id"`
	Topic     string                 `json:"topic"`
	Payload   interface{}            `json:"payload"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt time.Time              `json:"created_at"`
}

// ============ 常用Topic定义 ============

const (
	// 用户相关
	TopicUserRegistered = "user.registered"
	TopicUserLoggedIn   = "user.logged_in"
	TopicUserUpdated    = "user.updated"

	// 钱包相关
	TopicWalletRecharged = "wallet.recharged"
	TopicWalletConsumed  = "wallet.consumed"
	TopicWithdrawRequest = "wallet.withdraw_request"

	// 内容相关
	TopicContentCreated  = "content.created"
	TopicContentUpdated  = "content.updated"
	TopicContentDeleted  = "content.deleted"
	TopicContentReviewed = "content.reviewed"

	// 推荐相关
	TopicBehaviorRecorded      = "recommendation.behavior_recorded"
	TopicRecommendationRefresh = "recommendation.refresh"

	// 通知相关
	TopicNotificationEmail = "notification.email"
	TopicNotificationSMS   = "notification.sms"
	TopicNotificationPush  = "notification.push"
)
