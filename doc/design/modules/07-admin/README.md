# 07 - 管理模块

> **模块编号**: 07
> **模块名称**: Admin & Management
> **负责功能**: 后台管理、内容审核、数据统计、系统配置
> **完成度**: 🟡 50%

## 📋 目录结构

```
管理模块/
├── api/v1/
│   └── admin/                    # 管理API
│       ├── user_api.go          # 用户管理
│       ├── content_api.go       # 内容审核
│       ├── statistics_api.go    # 数据统计
│       ├── system_api.go        # 系统配置
│       └── log_api.go           # 操作日志
├── service/admin/                # 管理服务层
│   ├── user_service.go         # 用户管理服务
│   ├── content_service.go      # 内容审核服务
│   ├── statistics_service.go   # 统计服务
│   └── system_service.go       # 系统服务
├── repository/interfaces/admin/  # 仓储接口
├── repository/mongodb/admin/     # MongoDB仓储实现
└── models/admin/                 # 数据模型
    ├── audit_log.go             # 审计日志
    ├── report.go                # 举报
    └── system_config.go         # 系统配置
```

## 🎯 核心功能

### 1. 用户管理

- **用户列表**: 查询、筛选用户
- **用户详情**: 查看用户详细信息
- **用户封禁**: 封禁/解封用户
- **权限调整**: 调整用户角色和权限
- **操作记录**: 用户操作日志

### 2. 内容审核

- **作品审核**: 审核待发布作品
- **评论审核**: 审核违规评论
- **举报处理**: 处理用户举报
- **敏感词管理**: 管理敏感词库
- **审核规则**: 配置审核规则

### 3. 数据统计

- **平台数据**: 用户数、作品数、阅读量
- **用户统计**: 日活、月活、新增用户
- **作品统计**: 发布量、完本率
- **收入统计**: 充值、消费、分成
- **作者排行**: 各项数据排行

### 4. 系统配置

- **参数配置**: 系统参数设置
- **功能开关**: 功能开关控制
- **公告管理**: 平台公告发布
- **反馈管理**: 用户反馈处理

### 5. 操作日志

- **审计日志**: 管理员操作记录
- **日志查询**: 查询操作日志
- **异常告警**: 异常操作告警

## 📊 数据模型

### AuditLog (审计日志)

```go
type AuditLog struct {
    ID              primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
    AdminID         primitive.ObjectID   `bson:"admin_id" json:"adminId"`
    Action          string               `bson:"action" json:"action"`
    TargetType      string               `bson:"target_type" json:"targetType"`
    TargetID        primitive.ObjectID   `bson:"target_id" json:"targetId"`
    Details         map[string]interface{} `bson:"details,omitempty" json:"details,omitempty"`
    IP              string               `bson:"ip" json:"ip"`
    UserAgent       string               `bson:"user_agent" json:"userAgent"`
    CreatedAt       time.Time            `bson:"created_at" json:"createdAt"`
}
```

### Report (举报)

```go
type Report struct {
    ID              primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
    ReporterID      primitive.ObjectID   `bson:"reporter_id" json:"reporterId"`
    TargetType      string               `bson:"target_type" json:"targetType"`
    TargetID        primitive.ObjectID   `bson:"target_id" json:"targetId"`
    Reason          string               `bson:"reason" json:"reason"`
    Description     string               `bson:"description" json:"description"`

    // 处理信息
    Status          ReportStatus         `bson:"status" json:"status"`
    HandlerID       *primitive.ObjectID  `bson:"handler_id,omitempty" json:"handlerId,omitempty"`
    HandleResult    string               `bson:"handle_result,omitempty" json:"handleResult,omitempty"`
    HandledAt       *time.Time           `bson:"handled_at,omitempty" json:"handledAt,omitempty"`

    // 时间戳
    CreatedAt       time.Time            `bson:"created_at" json:"createdAt"`
    UpdatedAt       time.Time            `bson:"updated_at" json:"updatedAt"`
}

type ReportStatus string
const (
    ReportStatusPending   ReportStatus = "pending"
    ReportStatusProcessing ReportStatus = "processing"
    ReportStatusResolved  ReportStatus = "resolved"
    ReportStatusRejected  ReportStatus = "rejected"
)
```

### SystemConfig (系统配置)

```go
type SystemConfig struct {
    ID              primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
    Key             string               `bson:"key" json:"key"`
    Value           interface{}          `bson:"value" json:"value"`
    Type            ConfigType           `bson:"type" json:"type"`
    Description     string               `bson:"description" json:"description"`
    UpdatedBy       primitive.ObjectID   `bson:"updated_by" json:"updatedBy"`
    UpdatedAt       time.Time            `bson:"updated_at" json:"updatedAt"`
}

type ConfigType string
const (
    ConfigTypeString   ConfigType = "string"
    ConfigTypeNumber   ConfigType = "number"
    ConfigTypeBoolean  ConfigType = "boolean"
    ConfigTypeJSON     ConfigType = "json"
)
```

## 🌐 API端点

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| GET | /api/v1/admin/users | 获取用户列表 | 是(管理员) |
| GET | /api/v1/admin/users/:id | 获取用户详情 | 是(管理员) |
| PUT | /api/v1/admin/users/:id/status | 更新用户状态 | 是(管理员) |
| GET | /api/v1/admin/contents/pending | 待审核内容 | 是(管理员) |
| PUT | /api/v1/admin/contents/:id/review | 审核内容 | 是(管理员) |
| GET | /api/v1/admin/reports | 获取举报列表 | 是(管理员) |
| PUT | /api/v1/admin/reports/:id/handle | 处理举报 | 是(管理员) |
| GET | /api/v1/admin/statistics/overview | 平台概览统计 | 是(管理员) |
| GET | /api/v1/admin/statistics/users | 用户统计 | 是(管理员) |
| GET | /api/v1/admin/statistics/works | 作品统计 | 是(管理员) |
| GET | /api/v1/admin/statistics/revenue | 收入统计 | 是(管理员) |
| GET | /api/v1/admin/configs | 获取系统配置 | 是(管理员) |
| PUT | /api/v1/admin/configs/:key | 更新系统配置 | 是(管理员) |
| GET | /api/v1/admin/audit-logs | 获取审计日志 | 是(管理员) |

## 🔐 权限控制

### 管理员角色

| 角色 | 描述 | 权限 |
|------|------|------|
| SuperAdmin | 超级管理员 | 所有权限 |
| Admin | 管理员 | 大部分管理权限 |
| Editor | 编辑 | 内容审核权限 |
| Operator | 运营 | 数据查看权限 |

### 权限粒度

- **模块级**: 访问特定模块
- **操作级**: 特定操作权限
- **数据级**: 数据范围限制

## 🔧 依赖关系

### 依赖的模块
- **01 - 认证授权**: 管理员身份验证
- **所有业务模块**: 获取数据进行管理

### 外部服务
- **日志服务**: 操作日志存储
- **监控服务**: 系统监控

## 📈 扩展点

1. **工作流引擎**
   - 审核流程配置
   - 多级审核
   - 自动流转

2. **数据大屏**
   - 实时数据展示
   - 可视化图表
   - 自定义报表

3. **批量操作**
   - 批量审核
   - 批量导入/导出
   - 批量处理

---

**文档维护**: 青羽后端架构团队
**最后更新**: 2025-01-06
**对应实现**: `../../Qingyu_backend/api/v1/admin/`
