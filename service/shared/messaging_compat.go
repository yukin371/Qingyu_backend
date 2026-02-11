package shared

// 本文件提供向后兼容层，允许旧的import路径继续工作
// 当所有代码都迁移到新的import路径后，可以删除此文件

import (
	"Qingyu_backend/service/channels"
)

// ============ 类型别名（向后兼容）============

// MessagingService 消息队列服务接口
type MessagingService = channels.MessagingService

// MessageHandler 消息处理函数
type MessageHandler = channels.MessageHandler

// Message 消息
type Message = channels.Message

// NotificationService 通知服务接口
type NotificationService = channels.NotificationService

// EmailService 邮件服务接口
type EmailService = channels.EmailService

// InboxNotificationService 站内通知服务接口
type InboxNotificationService = channels.InboxNotificationService

// QueueClient 队列客户端接口
type QueueClient = channels.QueueClient

// StreamMessage Redis Stream消息
type StreamMessage = channels.StreamMessage

// EmailRequest 邮件请求
type EmailRequest = channels.EmailRequest

// EmailAttachment 邮件附件
type EmailAttachment = channels.EmailAttachment

// EmailResult 邮件发送结果
type EmailResult = channels.EmailResult

// EmailConfig 邮件配置
type EmailConfig = channels.EmailConfig

// InboxNotificationFilter 服务层过滤器
type InboxNotificationFilter = channels.InboxNotificationFilter

// ============ 实现类型别名 ============

// MessagingServiceImpl 消息服务实现
type MessagingServiceImpl = channels.MessagingServiceImpl

// NotificationServiceImpl 通知服务实现
type NotificationServiceImpl = channels.NotificationServiceImpl

// NotificationServiceComplete 完整的通知服务实现
type NotificationServiceComplete = channels.NotificationServiceComplete

// EmailServiceImpl 邮件服务实现
type EmailServiceImpl = channels.EmailServiceImpl

// InboxNotificationServiceImpl 站内通知服务实现
type InboxNotificationServiceImpl = channels.InboxNotificationServiceImpl

// RedisQueueClient Redis队列客户端实现
type RedisQueueClient = channels.RedisQueueClient

// TemplateStore 模板存储接口
type TemplateStore = channels.TemplateStore

// ============ 常量别名 ============

const (
	// 用户相关
	TopicUserRegistered = channels.TopicUserRegistered
	TopicUserLoggedIn   = channels.TopicUserLoggedIn
	TopicUserUpdated    = channels.TopicUserUpdated

	// 钱包相关
	TopicWalletRecharged = channels.TopicWalletRecharged
	TopicWalletConsumed  = channels.TopicWalletConsumed
	TopicWithdrawRequest = channels.TopicWithdrawRequest

	// 内容相关
	TopicContentCreated  = channels.TopicContentCreated
	TopicContentUpdated  = channels.TopicContentUpdated
	TopicContentDeleted  = channels.TopicContentDeleted
	TopicContentReviewed = channels.TopicContentReviewed

	// 推荐相关
	TopicBehaviorRecorded      = channels.TopicBehaviorRecorded
	TopicRecommendationRefresh = channels.TopicRecommendationRefresh

	// 通知相关
	TopicNotificationEmail = channels.TopicNotificationEmail
	TopicNotificationSMS   = channels.TopicNotificationSMS
	TopicNotificationPush  = channels.TopicNotificationPush
)

// ============ 构造函数别名（提供便捷创建）============

// NewMessagingService 创建消息服务
var NewMessagingService = channels.NewMessagingService

// NewNotificationService 创建通知服务
var NewNotificationService = channels.NewNotificationService

// NewNotificationServiceComplete 创建完整的通知服务
var NewNotificationServiceComplete = channels.NewNotificationServiceComplete

// NewEmailService 创建邮件服务
var NewEmailService = channels.NewEmailService

// NewInboxNotificationService 创建站内通知服务
var NewInboxNotificationService = channels.NewInboxNotificationService

// NewRedisQueueClient 创建Redis队列客户端
var NewRedisQueueClient = channels.NewRedisQueueClient

// ============ 迁移指南 ============

/*
旧的import路径:
  import "Qingyu_backend/service/shared/messaging"

新的import路径:
  import "Qingyu_backend/service/channels"

或者使用兼容层:
  import "Qingyu_backend/service/shared"

迁移步骤:
1. 更新import语句
2. 更新类型引用（例如: messaging.MessagingService -> shared.MessagingService）
3. 更新构造函数调用（例如: messaging.NewMessagingService -> shared.NewMessagingService）
4. 运行测试确保功能正常
5. 最终迁移到 channels 包

注意:
- 兼容层仅作为过渡期使用
- 建议尽快迁移到新的import路径 (service/channels)
- 迁移完成后可以删除此文件
*/
