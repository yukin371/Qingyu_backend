# Qingyu Backend API架构

## 1. 概述

本文档描述 Qingyu Backend 系统中的所有 API 端点，按模块分类，包含权限要求和使用说明。

## 2. API 设计原则

- RESTful 风格设计
- 统一响应格式
- JWT 认证 + 角色权限验证
- 版本控制 (v1)
- 统一错误处理

## 3. 统一响应格式

```json
{
  "code": 200,
  "message": "操作成功",
  "data": {}
}
```

## 4. Admin 模块 API

### 4.1 统计分析 API (Analytics)

**路径前缀**: `/api/v1/admin/analytics`

| 端点 | 方法 | 描述 | 权限要求 |
|------|------|------|----------|
| `/user-growth` | GET | 获取用户增长趋势 | admin |
| `/content-statistics` | GET | 获取内容统计 | admin |
| `/revenue-report` | GET | 获取收入报告 | admin |
| `/active-users` | GET | 获取活跃用户报告 | admin |
| `/system-overview` | GET | 获取系统概览 | admin |
| `/export` | GET | 导出统计分析报告 | admin |
| `/dashboard` | GET | 获取统计分析仪表板 | admin |
| `/custom` | POST | 自定义统计分析查询 | admin |
| `/compare` | GET | 对比不同时期的统计数据 | admin |
| `/realtime` | GET | 获取实时统计数据 | admin |
| `/predict` | GET | 获取统计分析预测 | admin |

#### 用户增长趋势
```
GET /api/v1/admin/analytics/user-growth
Query参数:
  - start_date: 开始日期 (YYYY-MM-DD) [必填]
  - end_date: 结束日期 (YYYY-MM-DD) [必填]
  - interval: 间隔 (daily/weekly/monthly) [必填]
```

#### 内容统计
```
GET /api/v1/admin/analytics/content-statistics
Query参数:
  - start_date: 开始日期 (YYYY-MM-DD) [可选]
  - end_date: 结束日期 (YYYY-MM-DD) [可选]
```

#### 收入报告
```
GET /api/v1/admin/analytics/revenue-report
Query参数:
  - start_date: 开始日期 (YYYY-MM-DD) [必填]
  - end_date: 结束日期 (YYYY-MM-DD) [必填]
  - interval: 间隔 (daily/weekly/monthly) [必填]
```

#### 活跃用户报告
```
GET /api/v1/admin/analytics/active-users
Query参数:
  - start_date: 开始日期 (YYYY-MM-DD) [必填]
  - end_date: 结束日期 (YYYY-MM-DD) [必填]
  - type: 类型 (dau/wau/mau) [必填]
```

### 4.2 审计追踪 API (Audit)

**路径前缀**: `/api/v1/admin/audit`

| 端点 | 方法 | 描述 | 权限要求 |
|------|------|------|----------|
| `/trail` | GET | 获取审计追踪 | admin |
| `/trail/resource/{type}/{id}` | GET | 获取资源审计追踪 | admin |
| `/trail/export` | GET | 导出审计追踪 | admin |
| `/statistics` | GET | 获取审计统计 | admin |

#### 获取审计追踪
```
GET /api/v1/admin/audit/trail
Query参数:
  - admin_id: 管理员ID [可选]
  - operation: 操作类型 [可选]
  - resource_type: 资源类型 [可选]
  - resource_id: 资源ID [可选]
  - start_date: 开始日期 (YYYY-MM-DD) [可选]
  - end_date: 结束日期 (YYYY-MM-DD) [可选]
  - page: 页码 [默认: 1]
  - size: 每页数量 [默认: 20]
```

#### 获取资源审计追踪
```
GET /api/v1/admin/audit/trail/resource/{type}/{id}
路径参数:
  - type: 资源类型 [必填]
  - id: 资源ID [必填]
```

#### 导出审计追踪
```
GET /api/v1/admin/audit/trail/export
Query参数:
  - format: 导出格式 (csv/xlsx) [默认: csv]
  - start_date: 开始日期 (YYYY-MM-DD) [可选]
  - end_date: 结束日期 (YYYY-MM-DD) [可选]
  - admin_id: 管理员ID [可选]
  - operation: 操作类型 [可选]
```

### 4.3 权限模板 API (PermissionTemplate)

**路径前缀**: `/api/v1/admin/permission-templates`

| 端点 | 方法 | 描述 | 权限要求 |
|------|------|------|----------|
| `` | GET | 获取权限模板列表 | admin |
| `/{id}` | GET | 获取权限模板详情 | admin |
| `` | POST | 创建权限模板 | admin |
| `/{id}` | PUT | 更新权限模板 | admin |
| `/{id}` | DELETE | 删除权限模板 | admin |
| `/system` | GET | 获取系统模板 | admin |
| `/{id}/apply` | POST | 应用权限模板到角色 | admin |

#### 获取权限模板列表
```
GET /api/v1/admin/permission-templates
Query参数:
  - category: 分类 [可选]
  - page: 页码 [默认: 1]
  - size: 每页数量 [默认: 20]
```

#### 创建权限模板
```
POST /api/v1/admin/permission-templates
Body:
{
  "name": "模板名称",
  "code": "template_code",
  "description": "模板描述",
  "permissions": ["permission1", "permission2"],
  "category": "custom"
}
```

### 4.4 内容导出 API (ContentExport)

**路径前缀**: `/api/v1/admin/content`

| 端点 | 方法 | 描述 | 权限要求 |
|------|------|------|----------|
| `/books/export` | GET | 导出书籍数据 | admin |
| `/books/export/template` | GET | 获取书籍导出模板 | admin |
| `/chapters/export` | GET | 导出章节数据 | admin |
| `/chapters/export/template` | GET | 获取章节导出模板 | admin |
| `/users/export` | GET | 导出用户数据 | admin |

#### 导出书籍数据
```
GET /api/v1/admin/content/books/export
Query参数:
  - format: 导出格式 (csv/excel) [默认: csv]
  - status: 状态筛选 [可选]
  - author: 作者筛选 [可选]
  - start_date: 开始日期 (YYYY-MM-DD) [可选]
  - end_date: 结束日期 (YYYY-MM-DD) [可选]
```

#### 导出章节数据
```
GET /api/v1/admin/content/chapters/export
Query参数:
  - format: 导出格式 (csv/excel) [默认: csv]
  - book_id: 书籍ID [可选]
  - start_date: 开始日期 (YYYY-MM-DD) [可选]
  - end_date: 结束日期 (YYYY-MM-DD) [可选]
```

## 5. Writer 模块 API

### 5.1 导出 API (Export)

**路径前缀**: `/api/v1/writer`

| 端点 | 方法 | 描述 | 权限要求 |
|------|------|------|----------|
| `/export/books` | GET | 导出作者的书籍 | author |
| `/export/chapters` | GET | 导出书籍的章节 | author |
| `/export/progress/{id}` | GET | 获取导出进度 | author |

### 5.2 审核申诉 API (Audit)

**路径前缀**: `/api/v1/writer/audit`

| 端点 | 方法 | 描述 | 权限要求 |
|------|------|------|----------|
| `/appeals` | GET | 获取审核申诉列表 | author |
| `/appeals` | POST | 创建审核申诉 | author |
| `/appeals/{id}` | GET | 获取申诉详情 | author |

## 6. 认证中间件

### 6.1 JWT 认证
所有 API 端点（除登录注册外）都需要 JWT 认证。

```go
// 使用示例
router.Use(auth.JWTAuth())
```

### 6.2 角色权限验证
管理员 API 需要管理员角色权限。

```go
// 使用示例
adminGroup.Use(auth.RequireRole("admin"))
```

### 6.3 动态权限检查
支持动态权限检查的端点会验证具体权限代码。

```go
// 使用示例
router.GET("/admin/permission-templates",
    auth.RequirePermission("permission.template.view"),
    handler.ListTemplates)
```

## 7. API 权限矩阵

| 功能模块 | 读者 | 作者 | 管理员 |
|----------|------|------|--------|
| 用户管理 | - | - | ✓ |
| 统计分析 | - | - | ✓ |
| 审计追踪 | - | - | ✓ |
| 权限管理 | - | - | ✓ |
| 内容导出 | - | - | ✓ |
| 书籍管理 | - | ✓ | ✓ |
| 章节管理 | - | ✓ | ✓ |
| 阅读历史 | ✓ | - | ✓ |
| 书架管理 | ✓ | - | ✓ |
| 社交互动 | ✓ | ✓ | ✓ |

## 8. API 版本控制

当前版本: `v1`

路径格式: `/api/v1/{module}/{endpoint}`

## 9. 错误码定义

| 错误码 | 描述 |
|--------|------|
| 200 | 操作成功 |
| 400 | 请求参数错误 |
| 401 | 未认证 |
| 403 | 权限不足 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

---

**文档版本**: v1.0
**最后更新**: 2026-02-27
**维护者**: yukin371
