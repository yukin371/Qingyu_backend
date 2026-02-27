# Admin API 模块 - 管理员接口

## 模块职责

**Admin（管理员）**模块负责系统管理员相关的所有功能，包括用户管理、权限管理、审核管理、统计分析、内容导出等。

## 核心功能

### 1. 用户管理
- 用户列表查询
- 用户详情查看
- 用户状态管理
- 用户数据导出

### 2. 权限管理
- 权限列表查询
- 权限模板管理
- 权限分配
- 权限检查

### 3. 审核管理
- 待审核内容查询
- 审核操作（批准/拒绝）
- 申诉处理
- 高风险审核

### 4. 统计分析
- 用户增长趋势
- 内容统计
- 收入报告
- 活跃用户报告
- 系统概览
- 自定义查询
- 数据预测

### 5. 内容导出
- 书籍数据导出（CSV/Excel）
- 章节数据导出
- 用户数据导出
- 导出历史查询

### 6. 系统管理
- 配置管理
- Banner管理
- 配额管理
- 事件重放

## 文件结构

```
api/v1/admin/
├── analytics_api.go              # 统计分析API
├── analytics_api_test.go         # 统计分析测试
├── announcement_api.go           # 公告管理API
├── announcement_api_test.go      # 公告管理测试
├── audit_admin_api.go            # 审核管理API（管理员）
├── audit_admin_api_test.go       # 审核管理测试
├── audit_api.go                  # 审计日志API
├── audit_api_test.go             # 审计日志测试
├── banner_api.go                 # Banner管理API
├── config_api.go                 # 配置管理API
├── content_export_api.go         # 内容导出API
├── content_export_api_test.go    # 内容导出测试
├── events_api.go                 # 事件管理API
├── events_api_test.go            # 事件管理测试
├── permission_api.go             # 权限管理API
├── permission_api_test.go        # 权限管理测试
├── permission_template_api.go    # 权限模板API
├── permission_template_api_test.go # 权限模板测试
├── quota_admin_api.go            # 配额管理API
├── system_admin_api.go           # 系统管理API
├── types.go                      # 共享类型定义
├── user_admin_api.go             # 用户管理API
├── user_admin_api_test.go        # 用户管理测试
├── user_export_api.go            # 用户导出API
├── user_export_api_test.go       # 用户导出测试
└── README.md                     # 本文档
```

## API路由总览

### 用户管理

| 方法 | 路径 | 说明 | Handler |
|------|------|------|---------|
| GET | /api/v1/admin/users | 用户列表 | UserAdminAPI.GetUsers |
| GET | /api/v1/admin/users/:id | 用户详情 | UserAdminAPI.GetUser |
| PUT | /api/v1/admin/users/:id/status | 更新用户状态 | UserAdminAPI.UpdateUserStatus |
| DELETE | /api/v1/admin/users/:id | 删除用户 | UserAdminAPI.DeleteUser |
| GET | /api/v1/admin/users/export | 导出用户数据 | UserExportAPI.ExportUsers |

### 权限管理

| 方法 | 路径 | 说明 | Handler |
|------|------|------|---------|
| GET | /api/v1/admin/permissions | 获取所有权限 | PermissionAPI.GetAllPermissions |
| GET | /api/v1/admin/permissions/:id | 获取权限详情 | PermissionAPI.GetPermission |
| POST | /api/v1/admin/permissions | 创建权限 | PermissionAPI.CreatePermission |
| PUT | /api/v1/admin/permissions/:id | 更新权限 | PermissionAPI.UpdatePermission |
| DELETE | /api/v1/admin/permissions/:id | 删除权限 | PermissionAPI.DeletePermission |
| GET | /api/v1/admin/permission-templates | 获取权限模板 | PermissionTemplateAPI.GetPermissionTemplates |
| POST | /api/v1/admin/permission-templates | 创建权限模板 | PermissionTemplateAPI.CreatePermissionTemplate |
| PUT | /api/v1/admin/permission-templates/:id | 更新权限模板 | PermissionTemplateAPI.UpdatePermissionTemplate |
| DELETE | /api/v1/admin/permission-templates/:id | 删除权限模板 | PermissionTemplateAPI.DeletePermissionTemplate |

### 审核管理

| 方法 | 路径 | 说明 | Handler |
|------|------|------|---------|
| GET | /api/v1/admin/audits/pending | 待审核内容 | AuditAdminAPI.GetPendingAudits |
| POST | /api/v1/admin/audits/:id/review | 审核操作 | AuditAdminAPI.ReviewAudit |
| POST | /api/v1/admin/audits/:id/appeal | 处理申诉 | AuditAdminAPI.ReviewAppeal |
| GET | /api/v1/admin/audits/high-risk | 高风险审核 | AuditAdminAPI.GetHighRiskAudits |
| GET | /api/v1/admin/audits/statistics | 审核统计 | AuditAdminAPI.GetAuditStatistics |
| POST | /api/v1/admin/audits/batch | 批量审核 | AuditAdminAPI.BatchReviewAudit |

### 统计分析

| 方法 | 路径 | 说明 | Handler |
|------|------|------|---------|
| GET | /api/v1/admin/analytics/user-growth | 用户增长趋势 | AnalyticsAPI.GetUserGrowthTrend |
| GET | /api/v1/admin/analytics/content | 内容统计 | AnalyticsAPI.GetContentStatistics |
| GET | /api/v1/admin/analytics/revenue | 收入报告 | AnalyticsAPI.GetRevenueReport |
| GET | /api/v1/admin/analytics/active-users | 活跃用户报告 | AnalyticsAPI.GetActiveUsersReport |
| GET | /api/v1/admin/analytics/overview | 系统概览 | AnalyticsAPI.GetSystemOverview |
| GET | /api/v1/admin/analytics/export | 导出分析报告 | AnalyticsAPI.ExportAnalyticsReport |
| GET | /api/v1/admin/analytics/dashboard | 数据看板 | AnalyticsAPI.GetAnalyticsDashboard |
| POST | /api/v1/admin/analytics/query | 自定义查询 | AnalyticsAPI.GetCustomAnalyticsQuery |
| GET | /api/v1/admin/analytics/compare | 周期对比 | AnalyticsAPI.CompareAnalyticsPeriods |
| GET | /api/v1/admin/analytics/realtime | 实时统计 | AnalyticsAPI.GetRealTimeStats |
| GET | /api/v1/admin/analytics/predictions | 数据预测 | AnalyticsAPI.GetAnalyticsPredictions |

### 内容导出

| 方法 | 路径 | 说明 | Handler |
|------|------|------|---------|
| POST | /api/v1/admin/export/books/csv | 导出书籍为CSV | ContentExportAPI.ExportBooksToCSV |
| POST | /api/v1/admin/export/books/excel | 导出书籍为Excel | ContentExportAPI.ExportBooksToExcel |
| POST | /api/v1/admin/export/chapters/csv | 导出章节为CSV | ContentExportAPI.ExportChaptersToCSV |
| GET | /api/v1/admin/export/books/template | 获取书籍导出模板 | ContentExportAPI.GetBookExportTemplate |
| GET | /api/v1/admin/export/chapters/template | 获取章节导出模板 | ContentExportAPI.GetChapterExportTemplate |
| GET | /api/v1/admin/export/history | 导出历史 | - |

### 审计日志

| 方法 | 路径 | 说明 | Handler |
|------|------|------|---------|
| GET | /api/v1/admin/audit/trail | 审计日志查询 | AuditAPI.GetAuditTrail |
| GET | /api/v1/admin/audit/resource | 资源审计日志 | AuditAPI.GetAuditTrailByResource |
| GET | /api/v1/admin/audit/action | 操作审计日志 | AuditAPI.GetAuditTrailByAction |
| GET | /api/v1/admin/audit/statistics | 审计统计 | AuditAPI.GetAuditStatistics |
| GET | /api/v1/admin/audit/export | 导出审计日志 | AuditAPI.ExportAuditTrail |

### 系统管理

| 方法 | 路径 | 说明 | Handler |
|------|------|------|---------|
| GET | /api/v1/admin/system/config | 获取系统配置 | ConfigAPI.GetConfig |
| PUT | /api/v1/admin/system/config | 更新系统配置 | ConfigAPI.UpdateConfig |
| GET | /api/v1/admin/system/info | 系统信息 | SystemAdminAPI.GetSystemInfo |
| GET | /api/v1/admin/system/health | 健康检查 | SystemAdminAPI.HealthCheck |

### 事件管理

| 方法 | 路径 | 说明 | Handler |
|------|------|------|---------|
| POST | /api/v1/admin/events/replay | 重放事件 | EventsAPI.Replay |

## 技术特点

### 1. 权限控制
- 基于角色的访问控制（RBAC）
- 权限模板系统
- 细粒度权限检查
- 动态权限验证

### 2. 审计追踪
- 完整的操作日志
- 敏感操作记录
- 审计日志导出

### 3. 数据安全
- 敏感数据脱敏
- 导出权限控制
- 操作审批流程

### 4. 性能优化
- 分页查询
- 缓存策略
- 异步导出

## 测试覆盖

当前测试覆盖率: **59.1%**

主要测试文件：
- `analytics_api_test.go` - 统计分析测试
- `announcement_api_test.go` - 公告管理测试
- `audit_admin_api_test.go` - 审核管理测试
- `audit_api_test.go` - 审计日志测试
- `content_export_api_test.go` - 内容导出测试
- `events_api_test.go` - 事件管理测试
- `permission_api_test.go` - 权限管理测试
- `permission_template_api_test.go` - 权限模板测试
- `user_admin_api_test.go` - 用户管理测试
- `user_export_api_test.go` - 用户导出测试

## 依赖关系

### 依赖的服务
- `service/admin` - 管理员服务层
- `service/auth` - 认证授权服务
- `service/audit` - 审计服务

### 被依赖的模块
- 前端管理系统
- 移动端管理界面

## 相关文档

- [Admin Service 设计](../../../service/admin/README.md)
- [权限系统设计](../../../docs/design/modules/01-auth/README.md)
- [审计系统设计](../../../docs/design/modules/07-admin/README.md)

---

**版本**: v1.0
**更新日期**: 2026-02-27
**维护者**: Admin模块开发组
