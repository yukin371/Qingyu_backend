# 青羽写作平台 - 通知系统实现文档

## 概述
本文档描述了为青羽写作平台实现的完整通知系统，包括站内通知、邮件通知、短信通知和推送通知功能。

---

## 1. 创建的文件列表

### 数据模型
- `D:\Github\青羽\Qingyu_backend\models\notification\notification.go`
  - Notification (通知模型)
  - NotificationPreference (通知偏好设置)
  - PushDevice (推送设备)
  - NotificationTemplate (通知模板)
  - NotificationStats (通知统计)

### 仓储接口
- `D:\Github\青羽\Qingyu_backend\repository\interfaces\notification\notification_repository.go`
  - NotificationRepository
  - NotificationPreferenceRepository
  - PushDeviceRepository
  - NotificationTemplateRepository

### 仓储实现 (MongoDB)
- `D:\Github\青羽\Qingyu_backend\repository\mongodb\notification\notification_repository_impl.go`
- `D:\Github\青羽\Qingyu_backend\repository\mongodb\notification\preference_repository_impl.go`
- `D:\Github\青羽\Qingyu_backend\repository\mongodb\notification\push_device_repository_impl.go`
- `D:\Github\青羽\Qingyu_backend\repository\mongodb\notification\template_repository_impl.go`

### 服务层
- `D:\Github\青羽\Qingyu_backend\service\notification\notification_service.go`
- `D:\Github\青羽\Qingyu_backend\service\notification\template_service.go`

### API层
- `D:\Github\青羽\Qingyu_backend\api\v1\notifications\notification_api.go`

### 路由层
- `D:\Github\青羽\Qingyu_backend\router\notifications\notification_router.go`

---

## 2. API端点完整列表

### 站内通知API

#### 2.1 获取通知列表
```
GET /api/v1/notifications
```
**查询参数:**
- `type` (可选): 通知类型 (system, social, content, reward, message, update, membership)
- `read` (可选): 是否已读 (true/false)
- `priority` (可选): 优先级 (low, normal, high, urgent)
- `keyword` (可选): 关键词搜索
- `limit` (可选): 每页数量，默认20，最大100
- `offset` (可选): 偏移量，默认0
- `sortBy` (可选): 排序字段 (created_at, priority, read_at)
- `sortDesc` (可选): 是否降序，默认true

**响应示例:**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "notifications": [...],
    "total": 100,
    "unreadCount": 10
  }
}
```

#### 2.2 获取通知详情
```
GET /api/v1/notifications/:id
```

#### 2.3 标记单个通知为已读
```
PUT /api/v1/notifications/:id/read
```

#### 2.4 批量标记通知为已读
```
PUT /api/v1/notifications/mark-read
```
**请求体:**
```json
{
  "ids": ["id1", "id2", "id3"]
}
```

#### 2.5 标记所有通知为已读
```
PUT /api/v1/notifications/read-all
```

#### 2.6 删除单个通知
```
DELETE /api/v1/notifications/:id
```

#### 2.7 批量删除通知
```
DELETE /api/v1/notifications/batch-delete
```
**请求体:**
```json
{
  "ids": ["id1", "id2", "id3"]
}
```

#### 2.8 删除所有通知
```
DELETE /api/v1/notifications/delete-all
```

#### 2.9 获取未读通知数量
```
GET /api/v1/notifications/unread-count
```

#### 2.10 获取通知统计
```
GET /api/v1/notifications/stats
```
**响应示例:**
```json
{
  "totalCount": 100,
  "unreadCount": 10,
  "typeCounts": {
    "system": 20,
    "social": 30,
    "content": 25,
    "reward": 5,
    "message": 10,
    "update": 8,
    "membership": 2
  },
  "priorityCounts": {
    "low": 50,
    "normal": 30,
    "high": 15,
    "urgent": 5
  }
}
```

---

### 通知偏好设置API

#### 2.11 获取通知偏好设置
```
GET /api/v1/notifications/preferences
```

#### 2.12 更新通知偏好设置
```
PUT /api/v1/notifications/preferences
```
**请求体:**
```json
{
  "enableSystem": true,
  "enableSocial": true,
  "enableContent": true,
  "enableReward": true,
  "enableMessage": true,
  "enableUpdate": true,
  "enableMembership": true,
  "emailNotification": {
    "enabled": false,
    "types": ["system", "reward"],
    "frequency": "immediate"
  },
  "smsNotification": {
    "enabled": false,
    "types": ["system"]
  },
  "pushNotification": true,
  "quietHoursStart": "22:00",
  "quietHoursEnd": "08:00"
}
```

#### 2.13 重置通知偏好设置
```
POST /api/v1/notifications/preferences/reset
```

---

### 邮件通知设置API

#### 2.14 获取邮件通知设置
```
GET /api/v1/user-management/email-notifications
```

#### 2.15 更新邮件通知设置
```
PUT /api/v1/user-management/email-notifications
```
**请求体:**
```json
{
  "enabled": true,
  "types": ["system", "content", "reward", "membership"],
  "frequency": "immediate"
}
```
**频率选项:**
- `immediate`: 立即发送
- `hourly`: 每小时汇总
- `daily`: 每天汇总

---

### 短信通知设置API

#### 2.16 获取短信通知设置
```
GET /api/v1/user-management/sms-notifications
```

#### 2.17 更新短信通知设置
```
PUT /api/v1/user-management/sms-notifications
```
**请求体:**
```json
{
  "enabled": true,
  "types": ["system", "reward", "membership"]
}
```

---

### 推送设备管理API

#### 2.18 注册推送设备
```
POST /api/v1/notifications/push/register
```
**请求体:**
```json
{
  "deviceType": "ios",
  "deviceToken": "device_token_here",
  "deviceId": "unique_device_id"
}
```
**设备类型:**
- `ios`: iOS设备
- `android`: Android设备
- `web`: Web浏览器

#### 2.19 取消注册推送设备
```
DELETE /api/v1/notifications/push/unregister/:deviceId
```

#### 2.20 获取用户的推送设备列表
```
GET /api/v1/notifications/push/devices
```

---

## 3. 数据模型定义

### 3.1 Notification (通知)
```go
type Notification struct {
    ID        string                 // 通知ID
    UserID    string                 // 用户ID
    Type      NotificationType       // 通知类型
    Priority  NotificationPriority   // 优先级
    Title     string                 // 标题
    Content   string                 // 内容
    Data      map[string]interface{} // 扩展数据
    Read      bool                   // 是否已读
    ReadAt    *time.Time             // 已读时间
    CreatedAt time.Time              // 创建时间
    ExpiresAt *time.Time             // 过期时间
}
```

### 3.2 NotificationPreference (通知偏好设置)
```go
type NotificationPreference struct {
    ID                      string                      // ID
    UserID                  string                      // 用户ID
    EnableSystem            bool                        // 启用系统通知
    EnableSocial            bool                        // 启用社交通知
    EnableContent           bool                        // 启用内容通知
    EnableReward            bool                        // 启用打赏通知
    EnableMessage           bool                        // 启用私信通知
    EnableUpdate            bool                        // 启用更新通知
    EnableMembership        bool                        // 启用会员通知
    EmailNotification       EmailNotificationSettings   // 邮件通知设置
    SMSNotification         SMSNotificationSettings     // 短信通知设置
    PushNotification        bool                        // 启用推送通知
    QuietHoursStart         *string                     // 免打扰开始时间 (HH:MM)
    QuietHoursEnd           *string                     // 免打扰结束时间 (HH:MM)
    CreatedAt               time.Time                   // 创建时间
    UpdatedAt               time.Time                   // 更新时间
}
```

### 3.3 PushDevice (推送设备)
```go
type PushDevice struct {
    ID          string    // 设备ID
    UserID      string    // 用户ID
    DeviceType  string    // 设备类型 (ios, android, web)
    DeviceToken string    // 设备令牌
    DeviceID    string    // 设备唯一标识
    IsActive    bool      // 是否激活
    LastUsedAt  time.Time // 最后使用时间
    CreatedAt   time.Time // 创建时间
}
```

### 3.4 NotificationTemplate (通知模板)
```go
type NotificationTemplate struct {
    ID          string                 // 模板ID
    Type        NotificationType       // 通知类型
    Action      string                 // 操作类型
    Title       string                 // 标题模板
    Content     string                 // 内容模板
    Variables   []string               // 模板变量列表
    Data        map[string]interface{} // 扩展数据
    Language    string                 // 语言 (zh-CN, en-US)
    IsActive    bool                   // 是否激活
    CreatedAt   time.Time              // 创建时间
    UpdatedAt   time.Time              // 更新时间
}
```

---

## 4. 通知类型列表

| 类型代码 | 类型名称 | 描述 | 使用场景 |
|---------|---------|------|---------|
| `system` | 系统通知 | 平台公告、活动通知 | 系统维护、平台公告、活动推广 |
| `social` | 社交通知 | 关注、点赞、评论 | 新关注、收到点赞、收到评论 |
| `content` | 内容通知 | 作品审核、上架、下架 | 作品审核结果、上架/下架通知 |
| `reward` | 打赏通知 | 收到打赏 | 用户收到打赏时 |
| `message` | 私信通知 | 收到私信 | 收到新私信时 |
| `update` | 更新通知 | 关注作品更新 | 关注的作品更新章节 |
| `membership` | 会员通知 | 会员到期、续费提醒 | 会员即将到期、已到期、续费成功 |

### 通知优先级

| 优先级 | 代码 | 描述 |
|--------|------|------|
| 低 | `low` | 一般性通知，不紧急 |
| 普通 | `normal` | 普通通知，默认优先级 |
| 高 | `high` | 重要通知，需要用户关注 |
| 紧急 | `urgent` | 紧急通知，需要立即处理 |

---

## 5. 通知模板示例

### 5.1 系统通知模板

#### 平台公告
```
类型: system
操作: announcement
标题: 平台公告
内容: {{title}}

{{content}}
变量: ["title", "content"]
```

#### 系统维护
```
类型: system
操作: maintenance
标题: 系统维护通知
内容: 尊敬的用户，系统将于{{startTime}}至{{endTime}}进行维护，期间部分功能可能无法使用，敬请谅解。
变量: ["startTime", "endTime"]
```

### 5.2 社交通知模板

#### 新关注
```
类型: social
操作: follow
标题: 您有新的关注者
内容: {{followerName}}关注了您，点击查看TA的主页。
变量: ["followerName", "followerId"]
```

#### 收到点赞
```
类型: social
操作: like
标题: 作品收到点赞
内容: {{likerName}}点赞了您的作品《{{bookTitle}}》。
变量: ["likerName", "bookTitle", "bookId"]
```

#### 收到评论
```
类型: social
操作: comment
标题: 作品收到新评论
内容: {{commenterName}}评论了您的作品《{{bookTitle}}》：{{commentContent}}
变量: ["commenterName", "bookTitle", "bookId", "commentContent"]
```

### 5.3 内容通知模板

#### 审核通过
```
类型: content
操作: review_approved
标题: 作品审核通过
内容: 恭喜！您的作品《{{bookTitle}}》已通过审核，现已上架。
变量: ["bookTitle", "bookId"]
```

#### 审核未通过
```
类型: content
操作: review_rejected
标题: 作品审核未通过
内容: 很遗憾，您的作品《{{bookTitle}}》未通过审核。原因：{{reason}}
变量: ["bookTitle", "bookId", "reason"]
```

#### 作品下架
```
类型: content
操作: book_offline
标题: 作品下架通知
内容: 您的作品《{{bookTitle}}》已被下架。原因：{{reason}}
变量: ["bookTitle", "bookId", "reason"]
```

### 5.4 打赏通知模板

#### 收到打赏
```
类型: reward
操作: received
标题: 收到打赏
内容: {{senderName}}打赏了您的作品《{{bookTitle}}》，金额：{{amount}}书币。
变量: ["senderName", "bookTitle", "bookId", "amount"]
```

### 5.5 私信通知模板

#### 收到私信
```
类型: message
操作: received
标题: 收到新私信
内容: {{senderName}}给您发送了一条私信，点击查看。
变量: ["senderName", "senderId"]
```

### 5.6 更新通知模板

#### 作品更新
```
类型: update
操作: chapter_update
标题: 关注作品更新
内容: 您关注的《{{bookTitle}}》更新了第{{chapterNumber}}章：{{chapterTitle}}
变量: ["bookTitle", "bookId", "chapterNumber", "chapterTitle", "chapterId"]
```

### 5.7 会员通知模板

#### 会员即将到期
```
类型: membership
操作: expiring_soon
标题: 会员即将到期
内容: 您的会员将在{{days}}天后到期，请及时续费以享受会员权益。
变量: ["days", "expireDate"]
```

#### 会员已到期
```
类型: membership
操作: expired
标题: 会员已到期
内容: 您的会员已到期，续费后可继续享受会员权益。
变量: ["expireDate"]
```

#### 会员续费成功
```
类型: membership
操作: renewed
标题: 会员续费成功
内容: 您的会员已成功续费，有效期至{{expireDate}}。
变量: ["expireDate"]
```

---

## 6. 使用示例

### 6.1 发送简单通知
```go
// 获取通知服务
notificationService, _ := serviceContainer.GetNotificationService()

// 发送系统通知
err := notificationService.SendNotification(
    ctx,
    userID,
    notification.NotificationTypeSystem,
    "系统维护通知",
    "系统将于今晚22:00进行维护",
    map[string]interface{}{
        "maintenanceTime": "2024-01-01 22:00:00",
        "duration": "2小时",
    },
)
```

### 6.2 使用模板发送通知
```go
// 使用模板发送通知
err := notificationService.SendNotificationWithTemplate(
    ctx,
    userID,
    notification.NotificationTypeSocial,
    "follow",
    map[string]interface{}{
        "followerName": "张三",
        "followerId":   "user123",
    },
)
```

### 6.3 批量发送通知
```go
// 批量发送通知（如系统公告）
userIDs := []string{"user1", "user2", "user3"}
err := notificationService.BatchSendNotification(
    ctx,
    userIDs,
    notification.NotificationTypeSystem,
    "平台活动",
    "春节期间活动火热进行中！",
    map[string]interface{}{
        "activityId": "act001",
        "startDate":  "2024-02-10",
    },
)
```

---

## 7. 数据库集合

系统使用MongoDB存储通知相关数据，包含以下集合：

| 集合名称 | 描述 |
|---------|------|
| `notifications` | 通知记录 |
| `notification_preferences` | 通知偏好设置 |
| `push_devices` | 推送设备信息 |
| `notification_templates` | 通知模板 |

---

## 8. 扩展功能说明

### 8.1 实时推送 (WebSocket)
服务层已预留接口，可集成WebSocket实现实时推送通知到前端。

### 8.2 邮件通知
已实现邮件通知设置接口，实际发送邮件需要集成邮件服务提供商（如SMTP、SendGrid、阿里云邮件推送等）。

### 8.3 短信通知
已实现短信通知设置接口，实际发送短信需要集成短信服务提供商（如阿里云短信、腾讯云短信等）。

### 8.4 推送通知
已实现推送设备管理接口，实际推送需要集成APNs（iOS）、FCM（Android）等服务。

---

## 9. 定期清理任务

系统提供了清理方法，可以定期执行：

### 9.1 清理过期通知
```go
count, err := notificationService.CleanupExpiredNotifications(ctx)
fmt.Printf("清理了 %d 条过期通知", count)
```

### 9.2 清理旧通知
```go
// 清理90天前的通知
count, err := notificationService.CleanupOldNotifications(ctx, 90)
fmt.Printf("清理了 %d 条旧通知", count)
```

### 9.3 清理不活跃的推送设备
```go
// 清理180天未使用的设备
beforeDate := time.Now().AddDate(0, 0, -180)
count, err := pushDeviceRepo.DeleteInactiveDevices(ctx, beforeDate)
fmt.Printf("清理了 %d 个不活跃设备", count)
```

建议在后台任务中定期执行这些清理操作（如每天凌晨执行）。

---

## 10. 安全性说明

1. **权限验证**: 所有通知API都需要用户认证，用户只能操作自己的通知
2. **数据隔离**: 所有查询都自动添加用户ID过滤，确保数据隔离
3. **输入验证**: 所有请求参数都经过验证，防止注入攻击
4. **错误处理**: 统一的错误处理机制，不泄露敏感信息

---

## 11. 性能优化建议

1. **索引优化**: 建议在MongoDB中为以下字段创建索引
   - `notifications.user_id`
   - `notifications.created_at`
   - `notifications.read`
   - `notifications.type`
   - `push_devices.user_id`
   - `push_devices.device_id`

2. **分页加载**: 通知列表使用分页加载，避免一次性加载大量数据

3. **缓存策略**: 对于未读数量等统计信息，可以使用Redis缓存

4. **批量操作**: 批量标记已读、批量删除等操作使用批量数据库操作，提高性能

---

## 12. 测试建议

建议编写以下测试用例：

1. **单元测试**
   - 通知创建、读取、更新、删除
   - 通知偏好设置的更新
   - 模板渲染功能

2. **集成测试**
   - API端点测试
   - 数据库操作测试
   - 服务层逻辑测试

3. **性能测试**
   - 大量通知的创建和查询
   - 批量操作性能测试

---

## 13. 未来扩展

可以考虑的功能扩展：

1. **通知分组**: 相似类型的通知可以合并显示
2. **通知提醒**: 在特定时间提醒用户查看未读通知
3. **智能推荐**: 根据用户行为智能推送相关通知
4. **通知历史**: 保留用户的通知历史记录
5. **多语言支持**: 模板系统已支持多语言，可以扩展更多语言
6. **通知统计**: 通知点击率、阅读率等统计分析

---

## 14. 联系与支持

如有问题或建议，请联系开发团队。

---

**文档版本**: 1.0
**最后更新**: 2026-01-03
**作者**: Claude Code
**状态**: ✅ 已完成实现
