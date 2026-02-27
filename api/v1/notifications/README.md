# Notifications API 模块 - 通知

## 模块职责

**Notifications（通知）**模块负责系统中所有通知的管理和推送，是系统向单个用户传递事件驱动信息的核心渠道。

## 核心功能

### 1. 通知管理
- 获取通知列表（支持分页、筛选、排序）
- 获取通知详情
- 标记已读/批量标记已读
- 标记全部已读
- 删除通知/批量删除
- 清除已读通知

### 2. 通知偏好设置
- 获取通知偏好
- 更新通知偏好
- 重置为默认设置
- 支持按类型设置通知开关

### 3. 推送设备管理
- 注册推送设备
- 取消注册设备
- 获取设备列表
- 支持多设备管理

### 4. 实时推送
- WebSocket端点获取
- 实时推送新通知
- 邮件/短信通知
- 重新发送通知

### 5. 统计信息
- 获取未读数量
- 获取通知统计
- 按类型分组统计

## 文件结构

```
api/v1/notifications/
├── notification_api.go     # 通知API处理器
├── dto/                    # 数据传输对象
│   ├── batch.go           # 批量操作DTO
│   ├── notification.go    # 通知DTO
│   └── websocket.go       # WebSocket DTO
└── README.md               # 本文档
```

## API路由总览

### 通知管理

| 方法 | 路径 | 说明 | Handler |
|------|------|------|---------|
| GET | /api/v1/notifications | 获取通知列表 | NotificationAPI.GetNotifications |
| GET | /api/v1/notifications/:id | 获取通知详情 | NotificationAPI.GetNotification |
| POST | /api/v1/notifications/:id/read | 标记已读 | NotificationAPI.MarkAsRead |
| POST | /api/v1/notifications/batch-read | 批量标记已读 | NotificationAPI.MarkMultipleAsRead |
| POST | /api/v1/notifications/read-all | 标记全部已读 | NotificationAPI.MarkAllAsRead |
| DELETE | /api/v1/notifications/:id | 删除通知 | NotificationAPI.DeleteNotification |
| POST | /api/v1/notifications/batch-delete | 批量删除 | NotificationAPI.BatchDeleteNotifications |
| DELETE | /api/v1/notifications/delete-all | 删除全部 | NotificationAPI.DeleteAllNotifications |
| POST | /api/v1/notifications/clear-read | 清除已读 | NotificationAPI.ClearReadNotifications |

### 通知偏好

| 方法 | 路径 | 说明 | Handler |
|------|------|------|---------|
| GET | /api/v1/notifications/preferences | 获取偏好设置 | NotificationAPI.GetNotificationPreference |
| PUT | /api/v1/notifications/preferences | 更新偏好设置 | NotificationAPI.UpdateNotificationPreference |
| POST | /api/v1/notifications/preferences/reset | 重置偏好 | NotificationAPI.ResetNotificationPreference |
| GET | /api/v1/user-management/email-notifications | 获取邮件设置 | NotificationAPI.GetEmailNotificationSettings |
| PUT | /api/v1/user-management/email-notifications | 更新邮件设置 | NotificationAPI.UpdateEmailNotificationSettings |
| GET | /api/v1/user-management/sms-notifications | 获取短信设置 | NotificationAPI.GetSMSNotificationSettings |
| PUT | /api/v1/user-management/sms-notifications | 更新短信设置 | NotificationAPI.UpdateSMSNotificationSettings |

### 推送设备

| 方法 | 路径 | 说明 | Handler |
|------|------|------|---------|
| POST | /api/v1/notifications/push/register | 注册设备 | NotificationAPI.RegisterPushDevice |
| DELETE | /api/v1/notifications/push/unregister/:deviceId | 取消注册 | NotificationAPI.UnregisterPushDevice |
| GET | /api/v1/notifications/push/devices | 设备列表 | NotificationAPI.GetPushDevices |

### 统计与实时

| 方法 | 路径 | 说明 | Handler |
|------|------|------|---------|
| GET | /api/v1/notifications/unread-count | 未读数量 | NotificationAPI.GetUnreadCount |
| GET | /api/v1/notifications/stats | 通知统计 | NotificationAPI.GetNotificationStats |
| POST | /api/v1/notifications/:id/resend | 重新发送 | NotificationAPI.ResendNotification |
| GET | /api/v1/notifications/ws-endpoint | WebSocket端点 | NotificationAPI.GetWSEndpoint |

## 数据模型

### Notification（通知）

```go
type Notification struct {
    ID          string
    UserID      string              // 接收用户ID
    Type        NotificationType    // 通知类型
    Title       string
    Content     string
    Data        map[string]interface{} // 附加数据
    Priority    NotificationPriority   // 优先级
    Read        bool                // 是否已读
    ReadAt      *time.Time          // 已读时间
    ExpiresAt   *time.Time          // 过期时间
    CreatedAt   time.Time
}

type NotificationType string
const (
    NotificationTypeSystem    NotificationType = "system"
    NotificationTypeSocial    NotificationType = "social"
    NotificationTypeContent   NotificationType = "content"
    NotificationTypeReward    NotificationType = "reward"
    NotificationTypeMessage   NotificationType = "message"
    NotificationTypeUpdate    NotificationType = "update"
    NotificationTypeMembership NotificationType = "membership"
)

type NotificationPriority string
const (
    NotificationPriorityLow     NotificationPriority = "low"
    NotificationPriorityNormal  NotificationPriority = "normal"
    NotificationPriorityHigh    NotificationPriority = "high"
    NotificationPriorityUrgent  NotificationPriority = "urgent"
)
```

### NotificationPreference（通知偏好）

```go
type NotificationPreference struct {
    UserID                    string
    SystemEnabled             bool
    SocialEnabled             bool
    ContentEnabled            bool
    RewardEnabled             bool
    MessageEnabled            bool
    UpdateEnabled             bool
    MembershipEnabled         bool
    EmailEnabled              bool
    SMSEnabled                bool
    PushEnabled               bool
    QuietHoursEnabled         bool
    QuietHoursStart           string
    QuietHoursEnd             string
}
```

### PushDevice（推送设备）

```go
type PushDevice struct {
    ID         string
    UserID     string
    Platform   string    // ios/android/web
    Token      string
    Active     bool
    LastUsed   time.Time
    CreatedAt  time.Time
}
```

## 技术特点

### 1. 统一参数验证
```go
// 使用shared验证器
if err := shared.GetValidator().Struct(req); err != nil {
    shared.HandleValidationError(c, err)
    return
}
```

### 2. 批量操作支持
```go
// 批量标记已读响应
type BatchOperationResponse struct {
    Success   bool     `json:"success"`
    Total     int      `json:"total"`
    Succeeded int      `json:"succeeded"`
    Failed    int      `json:"failed"`
    Errors    []string `json:"errors,omitempty"`
}
```

### 3. 分页与筛选
```go
// 支持的查询参数
?type=system           // 按类型筛选
&read=false            // 按已读状态筛选
&priority=high         // 按优先级筛选
&keyword=更新          // 关键词搜索
&limit=20              // 每页数量
&offset=0              // 偏移量
&sortBy=created_at     // 排序字段
&sortDesc=true         // 降序
```

### 4. WebSocket实时推送
```go
// WebSocket端点
ws://host/ws/notifications

// 通过Authorization Header传递token
// 支持实时推送新通知
```

### 5. 多渠道推送
- 应用内推送（WebSocket）
- 邮件推送
- 短信推送
- 移动端推送（APNs/FCM）

## 使用场景

### 场景1：获取并处理通知
```
1. 获取通知列表 → GET /notifications?type=social&read=false
2. 查看未读数量 → GET /notifications/unread-count
3. 标记已读 → POST /notifications/:id/read
4. 删除已读 → POST /notifications/clear-read
```

### 场景2：配置通知偏好
```
1. 获取当前偏好 → GET /notifications/preferences
2. 关闭不需要的通知 → PUT /notifications/preferences
3. 设置免打扰时段 → quietHoursEnabled=true
```

### 场景3：实时通知推送
```
1. 获取WebSocket端点 → GET /notifications/ws-endpoint
2. 建立WebSocket连接 → ws://host/ws/notifications
3. 接收实时推送
4. 处理通知操作
```

### 场景4：推送设备管理
```
1. 注册新设备 → POST /notifications/push/register
2. 查看设备列表 → GET /notifications/push/devices
3. 取消旧设备 → DELETE /notifications/push/unregister/:deviceId
```

## 与其他模块的关系

| 模块 | 关系 | 说明 |
|------|------|------|
| Announcements | 独立 | 通知是私有的，公告是公开的 |
| Messages | 独立 | 通知是系统触发，消息是用户发送 |
| Service Layer | 依赖 | 使用notification.Service和notification.ChannelService |
| WebSocket | 集成 | 通过NotificationWSHub实时推送 |

## 通信系统定位

在三个通信系统中，**Notifications** 的定位是：
- **方向**: System → User（系统向单个用户）
- **可见性**: 私有（仅接收者可见）
- **模式**: 事件驱动（系统事件触发）
- **存储**: 每个用户独立的通知集合
- **推送**: 主动推送（WebSocket/邮件/短信）

## 通知类型说明

| 类型 | 说明 | 示例 |
|------|------|------|
| system | 系统通知 | 账号安全、系统维护 |
| social | 社交通知 | 点赞、评论、关注 |
| content | 内容通知 | 作品更新、章节发布 |
| reward | 奖励通知 | 打赏收益、活动奖励 |
| message | 消息通知 | 新私信提醒 |
| update | 更新通知 | 版本更新、功能上线 |
| membership | 会员通知 | 会员续费、权益变动 |

## 重构改进

### Phase 3 完成的优化
1. 统一使用 `shared.GetValidator()` 进行参数验证
2. 统一使用 `shared.HandleValidationError()` 处理验证错误
3. 批量操作返回详细的结果统计
4. 完善的WebSocket实时推送支持

### 测试覆盖
- 单元测试：24个测试全部通过
- 功能覆盖：通知CRUD、批量操作、偏好设置、设备管理、WebSocket

## 相关文档

- [通信模块架构设计](../../../architecture/api_architecture.md#通信模块架构)
- [Announcements API](../announcements/README.md)
- [Messages API](../social/README.md#messaging模块)
- [Notification Service](../../../service/notification/README.md)

---

**版本**: v1.1
**更新日期**: 2026-02-27
**维护者**: Backend Communication Team
**测试覆盖率**: 26.1%（需补充测试）
