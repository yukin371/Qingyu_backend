# 05 - 通信消息模块

> **模块编号**: 05
> **模块名称**: Communication & Messaging
> **负责功能**: 实时消息推送和通知功能
> **完成度**: 🟡 45%

## 📋 目录结构

```
通信消息模块/
├── api/v1/
│   └── communication/            # 通信API
│       ├── websocket_api.go     # WebSocket连接
│       ├── notification_api.go  # 通知管理
│       └── message_api.go       # 消息管理
├── service/communication/        # 通信服务层
│   ├── websocket_service.go    # WebSocket服务
│   ├── notification_service.go # 通知服务
│   └── email_service.go        # 邮件服务
├── repository/interfaces/communication/ # 仓储接口
├── repository/mongodb/communication/    # MongoDB仓储实现
└── models/communication/                # 数据模型
    ├── notification.go          # 通知
    ├── message.go               # 消息
    └── template.go              # 模板
```

## 🎯 核心功能

### 1. 实时消息

- **WebSocket连接**: 建立实时双向通信
- **消息推送**: 服务端主动推送消息
- **在线状态**: 用户在线状态管理
- **消息确认**: 消息已读确认

### 2. 站内通知

- **系统通知**: 平台公告、系统消息
- **互动通知**: 评论、点赞、关注提醒
- **业务通知**: 订单、充值、订阅提醒
- **通知设置**: 通知开关和偏好设置

### 3. 邮件通知

- **验证邮件**: 注册验证、邮箱验证
- **提醒邮件**: 重要事件提醒
- **营销邮件**: 活动推广（需订阅）
- **邮件模板**: 可配置邮件模板

### 4. 短信通知

- **验证码**: 登录、注册验证码
- **重要提醒**: 账户安全、重要操作
- **短信模板**: 可配置短信模板

## 📊 数据模型

### Notification (通知)

```go
type Notification struct {
    ID              primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
    UserID          primitive.ObjectID   `bson:"user_id" json:"userId"`
    Type            NotificationType     `bson:"type" json:"type"`
    Title           string               `bson:"title" json:"title"`
    Content         string               `bson:"content" json:"content"`
    Data            map[string]interface{} `bson:"data,omitempty" json:"data,omitempty"`

    // 状态
    IsRead          bool                 `bson:"is_read" json:"isRead"`
    ReadAt          *time.Time           `bson:"read_at,omitempty" json:"readAt,omitempty"`

    // 时间戳
    CreatedAt       time.Time            `bson:"created_at" json:"createdAt"`
    ExpiresAt       *time.Time           `bson:"expires_at,omitempty" json:"expiresAt,omitempty"`
}

type NotificationType string
const (
    NotificationTypeSystem     NotificationType = "system"
    NotificationTypeComment    NotificationType = "comment"
    NotificationTypeLike       NotificationType = "like"
    NotificationTypeFollow     NotificationType = "follow"
    NotificationTypeOrder      NotificationType = "order"
    NotificationTypePayment    NotificationType = "payment"
)
```

### Message (消息)

```go
type Message struct {
    ID              primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
    ConversationID  primitive.ObjectID   `bson:"conversation_id" json:"conversationId"`
    SenderID        primitive.ObjectID   `bson:"sender_id" json:"senderId"`
    ReceiverID      primitive.ObjectID   `bson:"receiver_id" json:"receiverId"`
    Content         string               `bson:"content" json:"content"`
    MessageType     MessageType          `bson:"message_type" json:"messageType"`

    // 状态
    IsRead          bool                 `bson:"is_read" json:"isRead"`
    ReadAt          *time.Time           `bson:"read_at,omitempty" json:"readAt,omitempty"`

    // 时间戳
    CreatedAt       time.Time            `bson:"created_at" json:"createdAt"`
}

type MessageType string
const (
    MessageTypeText      MessageType = "text"
    MessageTypeImage     MessageType = "image"
    MessageTypeSystem    MessageType = "system"
)
```

## 🌐 API端点

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| WS | /api/v1/communication/ws | WebSocket连接 | 是 |
| GET | /api/v1/communication/notifications | 获取通知列表 | 是 |
| PUT | /api/v1/communication/notifications/:id/read | 标记通知已读 | 是 |
| PUT | /api/v1/communication/notifications/read-all | 全部标记已读 | 是 |
| DELETE | /api/v1/communication/notifications/:id | 删除通知 | 是 |
| GET | /api/v1/communication/messages | 获取消息列表 | 是 |
| POST | /api/v1/communication/messages | 发送消息 | 是 |
| PUT | /api/v1/communication/messages/:id/read | 标记消息已读 | 是 |

## 🔧 依赖关系

### 依赖的模块
- **01 - 认证授权**: 用户身份验证
- **04 - 社交互动**: 生成互动通知

### 外部服务
- **邮件服务**: SMTP服务或邮件API
- **短信服务**: 短信网关API
- **Redis**: WebSocket连接管理、消息队列

## 📈 扩展点

1. **消息模板**
   - 可配置消息模板
   - 多语言支持
   - 个性化消息

2. **推送通知**
   - 移动端推送
   - 桌面通知
   - 浏览器推送

3. **消息统计**
   - 消息送达率
   - 消息打开率
   - 用户活跃度

---

**文档维护**: 青羽后端架构团队
**最后更新**: 2025-01-06
**对应实现**: `../../Qingyu_backend/api/v1/communication/`
**相关设计**: [通信设计文档](../../communication/)
