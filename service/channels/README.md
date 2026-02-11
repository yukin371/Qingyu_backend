# Channels 通知渠道服务

## 概述

`service/channels` 模块负责处理各种通知渠道的发送服务，包括邮件、短信、站内信、推送通知等。

## 职责划分

在 Qingyu 后端架构中，有三个相关但职责不同的模块：

### 1. service/channels（本模块）
**职责**: 通知渠道服务 - 负责通过各种渠道发送通知

- **邮件服务** (`email_service.go`): SMTP邮件发送
- **短信服务**: 通过消息队列发送短信通知
- **站内通知服务** (`inbox_notification_service.go`): 应用内消息推送
- **消息队列服务** (`messaging_service.go`): Redis Stream消息队列
- **通知服务** (`notification_service.go`): 统一的通知发送接口

**核心功能**:
- 发送通知到各种渠道（邮件、短信、推送等）
- 管理通知模板
- 处理消息队列的发布和订阅
- 站内通知的CRUD操作

### 2. service/messaging
**职责**: 用户消息服务 - 管理用户之间的私信和对话

- 用户私信
- 对话管理
- 消息历史

### 3. service/notification
**职责**: 通知管理服务 - 管理用户的通知中心

- 通知聚合
- 通知偏好设置
- 通知分组和归档
- 通知统计

## 目录结构

```
service/channels/
├── interfaces.go                      # 核心接口定义
├── messaging_service.go               # 消息队列服务实现
├── notification_service.go            # 通知服务实现
├── notification_service_complete.go   # 完整的通知服务实现
├── email_service.go                   # 邮件服务实现
├── inbox_notification_service.go      # 站内通知服务实现
├── redis_queue_client.go              # Redis队列客户端
├── _migration/
│   └── shared_compat.go               # 向后兼容层
└── README.md                          # 本文档
```

## 核心接口

### MessagingService
消息队列服务接口，用于发布和订阅消息。

```go
type MessagingService interface {
    Publish(ctx context.Context, topic string, message interface{}) error
    PublishDelayed(ctx context.Context, topic string, message interface{}, delay time.Duration) error
    Subscribe(ctx context.Context, topic string, handler MessageHandler) error
    Unsubscribe(ctx context.Context, topic string) error
    CreateTopic(ctx context.Context, topic string) error
    DeleteTopic(ctx context.Context, topic string) error
    ListTopics(ctx context.Context) ([]string, error)
    Health(ctx context.Context) error
}
```

### EmailService
邮件服务接口，用于发送邮件通知。

```go
type EmailService interface {
    SendEmail(ctx context.Context, req *EmailRequest) error
    SendWithTemplate(ctx context.Context, to []string, template *MessageTemplate, variables map[string]string) error
    SendBatch(ctx context.Context, recipients []string, subject, body string) []EmailResult
    ValidateEmail(email string) bool
    Health(ctx context.Context) error
}
```

### NotificationService
通知服务接口，提供统一的通知发送方式。

```go
type NotificationService interface {
    CreateNotification(ctx context.Context, notification *NotificationDelivery) error
    GetNotification(ctx context.Context, notificationID string) (*NotificationDelivery, error)
    ListNotifications(ctx context.Context, userID string, page, pageSize int) ([]*NotificationDelivery, int64, error)
    MarkAsRead(ctx context.Context, notificationID string) error
    MarkAllAsRead(ctx context.Context, userID string) error
    DeleteNotification(ctx context.Context, notificationID string) error
    GetUnreadCount(ctx context.Context, userID string) (int64, error)
    SendEmailNotification(ctx context.Context, to, subject, content string) error
    SendTemplateEmail(ctx context.Context, to, templateName string, variables map[string]string) error
    CreateBatchNotifications(ctx context.Context, notifications []*NotificationDelivery) error
}
```

### InboxNotificationService
站内通知服务接口，管理应用内的通知消息。

```go
type InboxNotificationService interface {
    SendNotification(ctx context.Context, notification *InboxNotification) error
    SendBatchNotifications(ctx context.Context, notifications []*InboxNotification) error
    GetUserNotifications(ctx context.Context, userID string, filter *InboxNotificationFilter) ([]*InboxNotification, int64, error)
    GetNotificationByID(ctx context.Context, id string) (*InboxNotification, error)
    GetUnreadCount(ctx context.Context, userID string) (int64, error)
    MarkAsRead(ctx context.Context, id, userID string) error
    MarkAllAsRead(ctx context.Context, userID string) (int64, error)
    DeleteNotification(ctx context.Context, id, userID string) error
    PinNotification(ctx context.Context, id, userID string) error
    UnpinNotification(ctx context.Context, id, userID string) error
    CleanupExpired(ctx context.Context) (int64, error)
}
```

## 预定义主题

### 用户相关
- `user.registered`: 用户注册
- `user.logged_in`: 用户登录
- `user.updated`: 用户信息更新

### 钱包相关
- `wallet.recharged`: 钱包充值
- `wallet.consumed`: 钱包消费
- `wallet.withdraw_request`: 提现请求

### 内容相关
- `content.created`: 内容创建
- `content.updated`: 内容更新
- `content.deleted`: 内容删除
- `content.reviewed`: 内容审核

### 推荐相关
- `recommendation.behavior_recorded`: 行为记录
- `recommendation.refresh`: 推荐刷新

### 通知相关
- `notification.email`: 邮件通知
- `notification.sms`: 短信通知
- `notification.push`: 推送通知

## 使用示例

### 发送邮件通知

```go
import "Qingyu_backend/service/channels"

// 创建邮件服务
emailService := channels.NewEmailService(&channels.EmailConfig{
    SMTPHost:     "smtp.example.com",
    SMTPPort:     587,
    SMTPUsername: "user@example.com",
    SMTPPassword: "password",
    FromAddress:  "noreply@example.com",
    FromName:     "青羽",
})

// 发送邮件
err := emailService.SendEmail(ctx, &channels.EmailRequest{
    To:      []string{"recipient@example.com"},
    Subject: "欢迎注册青羽",
    Body:    "<h1>欢迎！</h1><p>感谢您的注册。</p>",
    IsHTML:  true,
})
```

### 使用模板发送邮件

```go
// 使用模板发送
err := emailService.SendWithTemplate(ctx, []string{"user@example.com"}, template, map[string]string{
    "username": "张三",
    "verifyUrl": "https://example.com/verify?code=123456",
})
```

### 发送站内通知

```go
// 创建站内通知服务
inboxService := channels.NewInboxNotificationService(repo)

// 发送通知
err := inboxService.SendNotification(ctx, &messagingModel.InboxNotification{
    ReceiverID: "user123",
    Type:       messagingModel.InboxNotificationTypeSystem,
    Title:      "系统通知",
    Content:    "您有新的消息",
    Priority:   messagingModel.InboxNotificationPriorityNormal,
})
```

### 发布消息到队列

```go
// 创建消息服务
msgService := channels.NewMessagingService(queueClient)

// 发布消息
err := msgService.Publish(ctx, "user.registered", map[string]interface{}{
    "user_id":  "user123",
    "username": "张三",
    "email":    "user@example.com",
})
```

### 订阅消息

```go
// 订阅用户注册事件
err := msgService.Subscribe(ctx, "user.registered", func(ctx context.Context, msg *channels.Message) error {
    // 处理用户注册事件
    userID := msg.Payload.(map[string]interface{})["user_id"].(string)
    // 发送欢迎邮件等...
    return nil
})
```

## 迁移指南

如果你正在从旧的 `service/shared/messaging` 迁移：

1. **更新import路径**:
   ```go
   // 旧路径
   import "Qingyu_backend/service/shared/messaging"

   // 新路径
   import "Qingyu_backend/service/channels"
   ```

2. **更新类型引用**:
   ```go
   // 旧代码
   service := messaging.NewNotificationService(...)

   // 新代码
   service := channels.NewNotificationService(...)
   ```

3. **使用兼容层**（过渡期）:
   ```go
   import "Qingyu_backend/service/channels/_migration"

   // 兼容层提供了与旧包相同的接口
   service := migration.NewNotificationService(...)
   ```

## 依赖关系

### 依赖的包
- `Qingyu_backend/models/messaging`: 消息相关数据模型
- `Qingyu_backend/repository/interfaces/shared`: 共享仓储接口
- `Qingyu_backend/repository/mongodb/messaging`: MongoDB消息仓储实现
- `github.com/redis/go-redis/v9`: Redis客户端

### 被依赖的包
- `service/notification`: 通知管理服务
- `service/user`: 用户服务
- `service/bookstore`: 书店服务
- API层 handlers

## 待实现功能 (Phase3)

- [ ] SMTP完整实现（当前为简化版）
- [ ] 短信服务集成
- [ ] 推送通知服务集成
- [ ] OAuth2邮件认证
- [ ] 邮件模板缓存
- [ ] 发送失败重试机制
- [ ] 发送队列管理
- [ ] 邮件发送统计
- [ ] 反垃圾邮件处理（SPF, DKIM, DMARC）
- [ ] 消息持久化和死信队列
- [ ] 消息优先级支持
- [ ] 消息批量发送优化

## 测试

```bash
# 运行单元测试
go test ./service/channels/...

# 运行集成测试
go test ./service/channels/... -tags=integration

# 查看测试覆盖率
go test ./service/channels/... -cover
```

## 维护者

- 架构重构团队
- 联系方式: 见项目主README

## 更新日志

### 2026-02-09
- 从 `service/shared/messaging` 迁移到独立模块
- 创建 `service/channels` 目录结构
- 添加向后兼容层 (`_migration/shared_compat.go`)
- 编写模块文档
