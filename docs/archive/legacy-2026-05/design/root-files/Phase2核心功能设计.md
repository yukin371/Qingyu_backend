# Phase 2 核心功能设计文档

**文档版本**: v1.0  
**创建日期**: 2025-10-27  
**设计者**: AI Assistant  
**状态**: 🟢 实施中

---

## 文档概述

本文档描述Phase 2核心功能的设计方案，包括文件存储、搜索、消息通知和数据统计四大模块。

### 设计原则

1. **快速实施**：优先实现核心功能，高级功能用TODO标记
2. **可扩展性**：框架设计支持后续扩展
3. **接口优先**：通过接口定义清晰边界
4. **渐进式增强**：P0 → P1 → P2逐步完善

---

## 一、文件存储系统设计

### 1.1 系统架构

```
┌─────────────────┐
│   Storage API   │ ← API层
├─────────────────┤
│ StorageService  │ ← Service层
├─────────────────┤
│ StorageBackend  │ ← Backend抽象层
├─────┬─────┬─────┤
│Local│MinIO│ OSS │ ← 具体实现
└─────┴─────┴─────┘
```

### 1.2 核心接口

```go
// StorageRepository 文件元数据Repository
type StorageRepository interface {
    // 文件元数据管理
    CreateFile(ctx context.Context, file *FileInfo) error
    GetFile(ctx context.Context, fileID string) (*FileInfo, error)
    UpdateFile(ctx context.Context, fileID string, updates map[string]interface{}) error
    DeleteFile(ctx context.Context, fileID string) error
    
    // Health检查
    Health(ctx context.Context) error
}

// StorageBackend 存储后端接口
type StorageBackend interface {
    Save(ctx context.Context, path string, reader io.Reader) error
    Load(ctx context.Context, path string) (io.ReadCloser, error)
    Delete(ctx context.Context, path string) error
    GetURL(ctx context.Context, path string, expiresIn time.Duration) (string, error)
}
```

### 1.3 P0功能（必须实现）

- ✅ 本地文件存储
- ✅ MinIO基础集成
- ✅ 小文件上传（<5MB）
- ✅ 文件下载
- ✅ 文件删除
- ✅ 文件信息查询

### 1.4 P1功能（框架+TODO）

- 🔵 大文件分片上传
- 🔵 断点续传
- 🔵 图片处理（缩略图、压缩）
- 🔵 CDN加速

### 1.5 数据模型

```go
type FileInfo struct {
    ID           string    `json:"id" bson:"_id"`
    UserID       string    `json:"user_id" bson:"user_id"`
    Filename     string    `json:"filename" bson:"filename"`
    Size         int64     `json:"size" bson:"size"`
    MimeType     string    `json:"mime_type" bson:"mime_type"`
    StoragePath  string    `json:"storage_path" bson:"storage_path"`
    MD5          string    `json:"md5" bson:"md5"`
    Status       string    `json:"status" bson:"status"`
    CreatedAt    time.Time `json:"created_at" bson:"created_at"`
    UpdatedAt    time.Time `json:"updated_at" bson:"updated_at"`
}
```

---

## 二、搜索功能设计

### 2.1 系统架构

```
┌─────────────────┐
│   Search API    │ ← API层
├─────────────────┤
│  SearchService  │ ← Service层
├─────────────────┤
│ SearchBackend   │ ← Backend抽象层
├─────┬───────────┤
│MongoDB│  ES     │ ← 具体实现
└─────┴───────────┘
```

### 2.2 核心接口

```go
// SearchService 搜索服务
type SearchService interface {
    // 书籍搜索
    SearchBooks(ctx context.Context, req *SearchRequest) (*SearchResult, error)
    
    // 文档搜索
    SearchDocuments(ctx context.Context, req *SearchRequest) (*SearchResult, error)
    
    // 搜索建议
    GetSuggestions(ctx context.Context, keyword string) ([]string, error)
    
    // 索引管理
    CreateIndex(ctx context.Context, collection string) error
    UpdateIndex(ctx context.Context, collection string, docID string) error
}
```

### 2.3 MongoDB全文索引实现（P0）

**索引创建**：
```javascript
// 书籍索引
db.books.createIndex(
  {
    "title": "text",
    "author": "text",
    "description": "text",
    "tags": "text"
  },
  {
    weights: {
      title: 10,
      author: 5,
      tags: 3,
      description: 1
    },
    default_language: "none" // 支持中文
  }
)
```

**查询示例**：
```go
filter := bson.M{
    "$text": bson.M{"$search": keyword},
}
opts := options.Find().
    SetProjection(bson.M{
        "score": bson.M{"$meta": "textScore"},
    }).
    SetSort(bson.M{
        "score": bson.M{"$meta": "textScore"},
    })
```

### 2.4 P1功能（TODO）

- 🔵 Elasticsearch集成
- 🔵 智能搜索建议（拼音、同义词）
- 🔵 搜索历史
- 🔵 搜索结果高亮

---

## 三、消息通知系统设计

### 3.1 系统架构

```
┌──────────────────────┐
│    Message API       │
├──────────────────────┤
│  MessagingService    │
├──────────┬───────────┤
│站内消息  │ 邮件通知   │
├──────────┼───────────┤
│  Redis   │   SMTP    │
└──────────┴───────────┘
```

### 3.2 核心接口

```go
// MessageRepository 消息Repository
type MessageRepository interface {
    // 消息队列
    CreateMessage(ctx context.Context, message *Message) error
    GetMessage(ctx context.Context, messageID string) (*Message, error)
    
    // 通知记录
    CreateNotification(ctx context.Context, notification *Notification) error
    ListNotifications(ctx context.Context, filter *NotificationFilter) ([]*Notification, int64, error)
    MarkAsRead(ctx context.Context, notificationID string) error
    GetUnreadCount(ctx context.Context, userID string) (int64, error)
    
    // 消息模板
    CreateTemplate(ctx context.Context, template *MessageTemplate) error
    GetTemplateByName(ctx context.Context, name string) (*MessageTemplate, error)
}
```

### 3.3 消息流程

```
发送消息
   ↓
创建通知记录
   ↓
推送到队列(Redis)
   ↓
消费者处理
   ├→ 站内消息
   ├→ 邮件发送
   └→ 短信发送(TODO)
```

### 3.4 P0功能

- ✅ 站内消息CRUD
- ✅ 邮件通知（SMTP）
- ✅ 消息模板管理
- ✅ Redis队列

### 3.5 P1功能（TODO）

- 🔵 短信通知
- 🔵 WebSocket实时推送
- 🔵 消息重试机制
- 🔵 RabbitMQ集成

### 3.6 邮件模板示例

```html
<!-- templates/email/welcome.html -->
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>欢迎加入青羽</title>
</head>
<body>
    <h1>欢迎, {{.Username}}!</h1>
    <p>感谢您注册青羽写作平台。</p>
    <!-- TODO(Phase3): 添加更丰富的样式 -->
</body>
</html>
```

---

## 四、数据统计系统设计

### 4.1 系统架构

```
┌─────────────────┐
│   Stats API     │
├─────────────────┤
│  StatsService   │
├─────────────────┤
│StatsRepository  │
├─────────────────┤
│    MongoDB      │
└─────────────────┘
```

### 4.2 核心接口

```go
// StatsRepository 统计Repository
type StatsRepository interface {
    // 用户统计
    CountUsers(ctx context.Context, filter *UserStatsFilter) (int64, error)
    GetActiveUsers(ctx context.Context, startDate, endDate time.Time) (int64, error)
    
    // 内容统计
    CountBooks(ctx context.Context) (int64, error)
    CountDocuments(ctx context.Context) (int64, error)
    
    // AI统计
    GetAIUsageStats(ctx context.Context, filter *AIStatsFilter) (*AIUsageStats, error)
}

// StatsService 统计服务
type StatsService interface {
    GetOverview(ctx context.Context) (*OverviewStats, error)
    GetUserStats(ctx context.Context) (*UserStats, error)
    GetContentStats(ctx context.Context) (*ContentStats, error)
    GetAIStats(ctx context.Context) (*AIStats, error)
}
```

### 4.3 数据模型

```go
type OverviewStats struct {
    TotalUsers      int64 `json:"total_users"`
    TotalBooks      int64 `json:"total_books"`
    TotalDocuments  int64 `json:"total_documents"`
    TotalAICalls    int64 `json:"total_ai_calls"`
    DailyActiveUsers int64 `json:"daily_active_users"`
}

type UserStats struct {
    TotalUsers       int64 `json:"total_users"`
    NewUsersToday    int64 `json:"new_users_today"`
    ActiveUsersToday int64 `json:"active_users_today"`
    // TODO(Phase3): 添加详细用户画像
}
```

### 4.4 P0功能

- ✅ 基础统计查询（总数、今日新增）
- ✅ 用户统计
- ✅ 内容统计
- ✅ AI使用统计

### 4.5 P1功能（TODO）

- 🔵 趋势分析（周对比、月对比）
- 🔵 报表生成
- 🔵 数据导出（Excel、CSV）
- 🔵 实时数据看板

---

## 五、数据库设计

### 5.1 新增集合

**files集合**：
```json
{
  "_id": "ObjectId",
  "user_id": "string",
  "filename": "string",
  "size": "int64",
  "mime_type": "string",
  "storage_path": "string",
  "md5": "string",
  "status": "string",
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

**messages集合**：
```json
{
  "_id": "ObjectId",
  "topic": "string",
  "payload": "object",
  "status": "string",
  "retry": "int",
  "created_at": "timestamp"
}
```

**notifications集合**：
```json
{
  "_id": "ObjectId",
  "user_id": "string",
  "type": "string",
  "title": "string",
  "content": "string",
  "is_read": "boolean",
  "created_at": "timestamp",
  "read_at": "timestamp"
}
```

### 5.2 索引设计

```javascript
// files集合索引
db.files.createIndex({"user_id": 1, "created_at": -1})
db.files.createIndex({"md5": 1}) // 去重

// notifications集合索引
db.notifications.createIndex({"user_id": 1, "is_read": 1, "created_at": -1})

// 书籍全文索引
db.books.createIndex({
  "title": "text",
  "author": "text",
  "description": "text",
  "tags": "text"
})
```

---

## 六、配置管理

### 6.1 config.yaml新增配置

```yaml
# 文件存储配置
storage:
  backend: "minio"  # local, minio
  minio:
    endpoint: "localhost:9000"
    access_key: "${MINIO_ACCESS_KEY}"
    secret_key: "${MINIO_SECRET_KEY}"
    bucket: "qingyu-files"
    use_ssl: false
  # TODO(Phase3): 支持阿里云OSS、腾讯云COS
  # oss:
  #   endpoint: ""
  #   access_key: ""
  
# 邮件配置
email:
  smtp:
    host: "smtp.gmail.com"
    port: 587
    username: "${SMTP_USERNAME}"
    password: "${SMTP_PASSWORD}"
    from: "noreply@qingyu.com"
  # TODO(Phase3): 短信配置
  # sms:
  #   provider: "aliyun"
  
# 搜索配置
search:
  backend: "mongodb"  # mongodb, elasticsearch
  # TODO(Phase3): Elasticsearch配置
  # elasticsearch:
  #   hosts: ["localhost:9200"]
```

---

## 七、API设计

### 7.1 文件存储API

```
POST   /api/v1/storage/upload           上传文件
GET    /api/v1/storage/download/:id     下载文件
DELETE /api/v1/storage/:id              删除文件
GET    /api/v1/storage/:id/info         文件信息
```

### 7.2 搜索API

```
GET    /api/v1/search/books?q=keyword   搜索书籍
GET    /api/v1/search/documents?q=keyword  搜索文档
GET    /api/v1/search/suggest?q=key     搜索建议
```

### 7.3 消息通知API

```
GET    /api/v1/messages                 消息列表
GET    /api/v1/messages/:id             消息详情
PUT    /api/v1/messages/:id/read        标记已读
DELETE /api/v1/messages/:id             删除消息
POST   /api/v1/admin/messages/broadcast 广播消息
```

### 7.4 统计API

```
GET    /api/v1/stats/overview           总体统计
GET    /api/v1/stats/users              用户统计
GET    /api/v1/stats/content            内容统计
GET    /api/v1/stats/ai                 AI统计
```

---

## 八、技术选型

### 8.1 文件存储

**选择**: MinIO（开源对象存储）

**理由**:
- 兼容S3 API
- 易于部署
- 支持分布式
- 后续可迁移到云OSS

**替代方案**:
- 阿里云OSS（商业云存储）
- 腾讯云COS（商业云存储）

### 8.2 搜索

**选择**: MongoDB全文索引

**理由**:
- 无需额外服务
- 满足MVP需求
- 实施快速

**后续升级**: Elasticsearch（Phase 3）

### 8.3 消息队列

**选择**: Redis List

**理由**:
- 已有Redis
- 简单可靠
- 满足基础需求

**后续升级**: RabbitMQ（Phase 3）

---

## 九、性能要求

| 功能 | 指标 | 目标值 |
|------|------|--------|
| 文件上传 | 成功率 | ≥99% |
| 文件下载 | 速度 | ≥5MB/s |
| 搜索响应 | 延迟 | <500ms |
| 消息送达 | 成功率 | ≥95% |
| 消息延迟 | 时间 | <1s |
| 统计查询 | 延迟 | <200ms |

---

## 十、安全设计

### 10.1 文件安全

- 文件类型验证
- 文件大小限制
- 病毒扫描（TODO）
- 访问权限控制（TODO）

### 10.2 搜索安全

- SQL注入防护（MongoDB自动）
- 关键词过滤
- 搜索频率限制

### 10.3 消息安全

- 用户身份验证
- 消息内容过滤
- 发送频率限制

---

## 十一、TODO清单

### 高优先级（Phase 3）

- [ ] 大文件分片上传
- [ ] Elasticsearch集成
- [ ] WebSocket实时推送
- [ ] 短信通知集成

### 中优先级（Phase 4）

- [ ] 图片处理服务
- [ ] CDN加速
- [ ] 搜索历史和热门搜索
- [ ] 趋势分析和报表

### 低优先级（后续）

- [ ] 文件版本管理
- [ ] 高级图片处理（水印）
- [ ] 语义搜索
- [ ] 实时数据看板

---

## 十二、测试策略

### 12.1 单元测试

- Repository层测试覆盖率 >70%
- Service层测试覆盖率 >70%
- Mock外部依赖（MinIO、SMTP）

### 12.2 集成测试

- 文件上传下载完整流程
- 搜索功能端到端
- 邮件发送流程

### 12.3 性能测试

- 文件上传并发测试
- 搜索压力测试
- 消息队列吞吐量测试

---

## 附录

### A. 参考文档

- MinIO文档: https://docs.min.io/
- MongoDB全文索引: https://docs.mongodb.com/manual/core/index-text/
- SMTP邮件发送: Go net/smtp包

### B. 版本历史

| 版本 | 日期 | 变更内容 | 作者 |
|------|------|---------|------|
| v1.0 | 2025-10-27 | 初始版本 | AI Assistant |

---

**文档维护者**: AI Assistant  
**最后更新**: 2025-10-27  
**审核状态**: 待审核

