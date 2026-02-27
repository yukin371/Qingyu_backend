# Social API 模块 - 社交

## 模块职责

**Social（社交）**模块负责用户之间的社交互动功能，包括私信、关注、评论、点赞等。

## 核心功能

### 1. Messaging（私信）
- 基于会话的消息系统
- 一对一私聊
- 实时消息推送（WebSocket）
- 消息已读状态
- 附件支持
- 消息回复

### 2. Follow（关注）
- 关注用户
- 取消关注
- 获取关注列表
- 获取粉丝列表

### 3. Comment（评论）
- 发布评论
- 获取评论列表
- 删除评论
- 评论点赞

### 4. Like（点赞）
- 点赞内容
- 取消点赞
- 获取点赞列表

## 文件结构

```
api/v1/social/
├── message_api.go          # 消息API（MessageAPIV2，新版本）
├── follow_api.go           # 关注API
├── comment_api.go          # 评论API
├── like_api.go             # 点赞API
├── dto/                    # 数据传输对象
│   ├── message.go         # 消息DTO
│   ├── conversation.go    # 会话DTO
│   ├── follow.go          # 关注DTO
│   ├── comment.go         # 评论DTO
│   └── like.go            # 点赞DTO
└── README.md               # 本文档
```

## Messaging模块

### API路由

| 方法 | 路径 | 说明 | Handler |
|------|------|------|---------|
| GET | /api/v1/social/messages/conversations/:conversationId/messages | 获取消息列表 | MessageAPIV2.GetMessages |
| POST | /api/v1/social/messages/conversations/:conversationId/messages | 发送消息 | MessageAPIV2.SendMessage |
| POST | /api/v1/social/messages/conversations | 创建会话 | MessageAPIV2.CreateConversation |
| POST | /api/v1/social/messages/conversations/:conversationId/read | 标记已读 | MessageAPIV2.MarkConversationRead |

### 数据模型

#### Conversation（会话）

```go
type Conversation struct {
    ID             string
    ParticipantIDs []string    // 参与者ID列表
    Type           string      // direct（一对一）/ group（群组，预留）
    LastMessageAt  time.Time   // 最后消息时间
    CreatedAt      time.Time
    UpdatedAt      time.Time
}

// HasParticipant 检查用户是否参与会话
func (c *Conversation) HasParticipant(userID string) bool
```

#### DirectMessage（消息）

```go
type DirectMessage struct {
    ID             string
    ConversationID string
    SenderID       string      // 发送者ID
    ReceiverID     string      // 接收者ID
    Content        string      // 消息内容
    Type           MessageType // 消息类型
    Status         MessageStatus
    Extra          map[string]interface{} // 附件等额外信息
    ParentID       *string      // 回复的消息ID
    IsRead         bool         // 是否已读
    ReadAt         *time.Time   // 已读时间
    CreatedAt      time.Time
}

type MessageType string
const (
    MessageTypeText     MessageType = "text"
    MessageTypeImage    MessageType = "image"
    MessageTypeFile     MessageType = "file"
    MessageTypeAudio    MessageType = "audio"
    MessageTypeVideo    MessageType = "video"
    MessageTypeSystem   MessageType = "system"
)

type MessageStatus string
const (
    MessageStatusNormal   MessageStatus = "normal"
    MessageStatusDeleted  MessageStatus = "deleted"
    MessageStatusRecalled MessageStatus = "recalled"
)
```

### 技术特点

#### 1. 基于会话的消息系统
```go
// 创建会话
POST /api/v1/social/messages/conversations
{
    "participantIDs": ["user1", "user2"]
}

// 发送消息
POST /api/v1/social/messages/conversations/:conversationId/messages
{
    "content": "你好",
    "type": "text",
    "attachments": [...],
    "replyTo": "msg123"  // 可选，回复消息
}
```

#### 2. 实时推送（WebSocket）
```go
// 通过MessagingWSHub实时推送消息
if api.wsHub != nil {
    api.wsHub.SendMessage(conversationID, wsMessage, senderID)
}
```

#### 3. 分页获取消息
```go
// 支持向上/向下翻页
GET /api/v1/social/messages/conversations/:conversationId/messages?page=1&page_size=20
GET /api/v1/social/messages/conversations/:conversationId/messages?before=msg123
GET /api/v1/social/messages/conversations/:conversationId/messages?after=msg456
```

#### 4. 消息附件
```go
type MessageAttachmentDTO struct {
    Type       string `json:"type"`       // image/file/audio/video
    URL        string `json:"url"`        // 资源URL
    Name       string `json:"name"`       // 文件名
    Size       int64  `json:"size"`       // 文件大小（字节）
    Duration   int    `json:"duration"`   // 时长（秒，音频/视频）
    Width      int    `json:"width"`      // 宽度（图片/视频）
    Height     int    `json:"height"`     // 高度（图片/视频）
    Thumbnail  string `json:"thumbnail"`  // 缩略图URL
}
```

### 使用场景

#### 场景1：发送私信
```
1. 创建会话 → POST /conversations
2. 发送消息 → POST /conversations/:id/messages
3. WebSocket实时推送给接收者
4. 接收者标记已读 → POST /conversations/:id/read
```

#### 场景2：获取历史消息
```
1. 获取消息列表 → GET /conversations/:id/messages?page=1
2. 向上翻页加载更早消息 → ?before=msg123
3. 向下翻页加载更新消息 → ?after=msg456
```

#### 场景3：回复消息
```
1. 发送消息时指定replyTo字段
2. 显示消息回复关系
3. 支持多级回复
```

### 版本说明

- **MessageAPIV2**: 当前版本（`api/v1/social/message_api.go`）
  - 基于会话的消息系统
  - 支持WebSocket实时推送
  - 支持附件和回复

- **MessageAPI (旧版)**: 已废弃（`api/v1/messages/message_api.go`）
  - 预留@提醒功能待迁移
  - 计划在Phase 4整合

## 与其他模块的关系

| 模块 | 关系 | 说明 |
|------|------|------|
| Announcements | 独立 | 消息是点对点，公告是一对多 |
| Notifications | 独立 | 消息是用户发送，通知是系统触发 |
| WebSocket | 集成 | 通过MessagingWSHub实时推送 |
| Service Layer | 依赖 | 使用messaging.MessageService和ConversationService |

## 通信系统定位

在三个通信系统中，**Messages** 的定位是：
- **方向**: User ↔ User（用户之间）
- **可见性**: 私有（仅参与者可见）
- **模式**: 点对点（一对一私聊）
- **存储**: 按会话组织的消息集合
- **推送**: 实时推送（WebSocket）+ 被动获取

## Follow模块（待完善）

### 计划功能
- 关注/取消关注用户
- 获取关注列表
- 获取粉丝列表
- 检查关注状态

### 计划路由
```
POST   /api/v1/social/follow/:userId
DELETE /api/v1/social/follow/:userId
GET    /api/v1/social/following
GET    /api/v1/social/followers
GET    /api/v1/social/follow/:userId/status
```

## Comment模块（待完善）

### 计划功能
- 发布评论
- 获取评论列表
- 删除评论
- 评论点赞

## Like模块（待完善）

### 计划功能
- 点赞内容
- 取消点赞
- 获取点赞列表
- 检查点赞状态

## 重构改进

### Phase 3 完成的优化
1. MessageAPIV2代码规范化
2. 完善的DTO定义
3. WebSocket实时推送集成
4. 支持消息附件和回复
5. 会话权限验证

### 测试覆盖
- 单元测试：30个测试全部通过
- 功能覆盖：消息CRUD、会话管理、权限检查、参数验证

## 相关文档

- [通信模块架构设计](../../../architecture/api_architecture.md#通信模块架构)
- [Announcements API](../announcements/README.md)
- [Notifications API](../notifications/README.md)
- [Message Service](../../../service/messaging/README.md)

---

**版本**: v1.0
**更新日期**: 2026-02-27
**维护者**: Backend Communication Team
