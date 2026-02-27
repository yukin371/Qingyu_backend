# Announcements API 模块 - 公告

## 模块职责

**Announcements（公告）**模块负责平台公告的发布和管理，是平台向用户传递重要信息的官方渠道。

## 核心功能

### 1. 公告发布
- 管理员创建公告
- 设置目标受众（全部/读者/作者/管理员）
- 设置优先级和有效期
- 富文本内容支持

### 2. 公告展示
- 获取有效公告列表
- 按角色筛选
- 公告详情查看
- 查看次数统计

### 3. 公告管理
- 更新公告内容
- 批量更新状态
- 删除公告
- 批量删除

## 文件结构

```
api/v1/announcements/
├── announcement_api.go    # 公开API处理器
└── README.md              # 本文档

api/v1/admin/
├── announcement_api.go    # 管理员API处理器
└── README.md              # 管理员文档
```

## API路由总览

### 公开接口（无需认证）

| 方法 | 路径 | 说明 | Handler |
|------|------|------|---------|
| GET | /api/v1/announcements/effective | 获取有效公告 | AnnouncementPublicAPI.GetEffectiveAnnouncements |
| GET | /api/v1/announcements/:id | 获取公告详情 | AnnouncementPublicAPI.GetAnnouncementByID |
| POST | /api/v1/announcements/:id/view | 增加查看次数 | AnnouncementPublicAPI.IncrementViewCount |

### 管理员接口（需要认证和管理员权限）

| 方法 | 路径 | 说明 | Handler |
|------|------|------|---------|
| GET | /api/v1/admin/announcements | 获取公告列表 | AdminAPI.GetAnnouncements |
| GET | /api/v1/admin/announcements/:id | 获取公告详情 | AdminAPI.GetAnnouncementByID |
| POST | /api/v1/admin/announcements | 创建公告 | AdminAPI.CreateAnnouncement |
| PUT | /api/v1/admin/announcements/:id | 更新公告 | AdminAPI.UpdateAnnouncement |
| DELETE | /api/v1/admin/announcements/:id | 删除公告 | AdminAPI.DeleteAnnouncement |
| POST | /api/v1/admin/announcements/batch/status | 批量更新状态 | AdminAPI.BatchUpdateStatus |
| DELETE | /api/v1/admin/announcements/batch | 批量删除 | AdminAPI.BatchDeleteAnnouncements |

## 数据模型

### Announcement（公告）

```go
type Announcement struct {
    ID            string
    Title         string
    Content       string
    TargetRole    string    // all/reader/writer/admin
    Priority      int       // 优先级
    Status        string    // draft/active/archived
    EffectiveFrom time.Time
    EffectiveTo   time.Time
    ViewCount     int64
    CreatedBy     string
    CreatedAt     time.Time
    UpdatedAt     time.Time
}
```

## 技术特点

### 1. 辅助函数应用
- 使用 `shared.GetRequiredParam` 获取必需参数
- 使用 `shared.GetIntParam` 获取整数参数（带范围验证）
- 统一使用 `pkg/response` 包处理响应

### 2. 参数验证
```go
// 获取必需参数
id, ok := shared.GetRequiredParam(c, "id", "公告ID")
if !ok {
    return
}

// 获取整数参数（带范围验证）
limit := shared.GetIntParam(c, "limit", true, 10, 1, 50)
```

### 3. 响应格式统一
```go
// 成功响应
response.Success(c, data)

// 错误响应
response.InternalError(c, err)
response.NotFound(c, "公告不存在")
```

### 4. 角色筛选
支持按角色筛选公告：
- `all`: 所有用户
- `reader`: 仅读者
- `writer`: 仅作者
- `admin`: 仅管理员

## 使用场景

### 场景1：用户获取有效公告
```
1. 访问公告列表 → GET /announcements/effective?targetRole=reader
2. 浏览公告标题和摘要
3. 点击查看详情 → GET /announcements/:id
4. 记录查看次数 → POST /announcements/:id/view
```

### 场景2：管理员发布公告
```
1. 创建公告 → POST /admin/announcements
2. 设置标题、内容、目标角色
3. 设置有效期和优先级
4. 发布后立即可见（状态为active）
```

### 场景3：批量管理公告
```
1. 批量更新状态 → POST /admin/announcements/batch/status
2. 批量删除公告 → DELETE /admin/announcements/batch
3. 归档过期公告
```

## 与其他模块的关系

| 模块 | 关系 | 说明 |
|------|------|------|
| Notifications | 独立 | 公告是公开的，通知是私有的 |
| Messages | 独立 | 公告是一对多，消息是点对点 |
| Admin | 依赖 | 管理员API在admin模块中 |

## 通信系统定位

在三个通信系统中，**Announcements** 的定位是：
- **方向**: Platform → Users（平台向用户）
- **可见性**: 公开（所有用户可见）
- **模式**: 一对多（一个公告，多个接收者）
- **存储**: 集中式存储在Announcement集合
- **推送**: 被动获取（用户主动获取）+ WebSocket广播

## 重构改进

### Phase 3 完成的优化
1. 应用 `GetRequiredParam` 辅助函数 - 提高代码复用
2. 应用 `GetIntParam` 辅助函数 - 统一参数验证
3. 统一使用 `pkg/response` 包 - 规范响应格式
4. 移除重复的参数验证代码 - 减少代码冗余

### 废弃文件
- ~~`api/v1/shared/notification_api.go`~~ (已删除，252行)

## 相关文档

- [通信模块架构设计](../../../architecture/api_architecture.md#通信模块架构)
- [Notifications API](../notifications/README.md)
- [Messages API](../social/README.md#messaging模块)
- [Admin API](../admin/README.md)

---

**版本**: v1.1
**更新日期**: 2026-02-27
**维护者**: Backend Communication Team
**测试覆盖率**: 良好（所有核心功能已测试）
