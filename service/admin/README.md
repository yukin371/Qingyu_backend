# Admin Service 模块 - 管理员服务

## 模块职责

**Admin Service（管理员服务）**模块是管理员功能的核心服务层，负责处理所有管理员相关的业务逻辑，包括统计分析、审计日志、内容导出、敏感操作管理等。

## 核心功能

### 1. 统计分析服务 (analytics_service)
- 用户增长趋势分析
- 内容统计分析
- 收入报告生成
- 活跃用户统计
- 系统概览数据
- 自定义查询
- 周期对比分析
- 实时统计
- 数据预测

### 2. 审计日志服务 (audit_log_service)
- 审计日志记录
- 审计日志查询
- 审计日志统计
- 审计日志导出

### 3. 导出服务 (export_service)
- 书籍数据导出
- 章节数据导出
- 用户数据导出
- 导出历史记录
- 导出任务管理

### 4. 敏感操作服务 (sensitive_operation_service)
- 敏感操作记录
- 敏感操作验证
- 敏感操作审计

### 5. 用户管理服务 (user_admin_service)
- 用户查询
- 用户状态管理
- 用户信息更新

### 6. 管理员服务 (admin_service)
- 管理员基础功能
- 管理员权限验证

## 文件结构

```
service/admin/
├── admin_service.go                    # 管理员基础服务
├── analytics_service.go                # 统计分析服务
├── analytics_service_test.go           # 统计分析测试
├── audit_log_service.go                # 审计日志服务
├── audit_log_service_test.go           # 审计日志测试
├── export_service.go                   # 导出服务
├── export_service_test.go              # 导出测试
├── interfaces.go                       # 服务接口定义
├── sensitive_operation_service.go      # 敏感操作服务
├── sensitive_operation_service_test.go # 敏感操作测试
└── user_admin_service.go               # 用户管理服务
```

## 服务接口

### AnalyticsService

统计分析服务接口，提供各类数据分析功能。

```go
type AnalyticsService interface {
    // 用户增长趋势
    GetUserGrowthTrend(ctx context.Context, startDate, endDate string, interval string) (*UserGrowthData, error)

    // 内容统计
    GetContentStatistics(ctx context.Context, startDate, endDate string) (*ContentStatistics, error)

    // 收入报告
    GetRevenueReport(ctx context.Context, startDate, endDate string) (*RevenueReport, error)

    // 活跃用户报告
    GetActiveUsersReport(ctx context.Context, reportType string) (*ActiveUsersReport, error)

    // 系统概览
    GetSystemOverview(ctx context.Context) (*SystemOverview, error)

    // 导出分析报告
    ExportAnalyticsReport(ctx context.Context, reportType, format string) ([]byte, error)

    // 数据看板
    GetAnalyticsDashboard(ctx context.Context) (*AnalyticsDashboard, error)

    // 自定义查询
    GetCustomAnalyticsQuery(ctx context.Context, queryJSON string) (*QueryResult, error)

    // 周期对比
    CompareAnalyticsPeriods(ctx context.Context, period1, period2 AnalyticsPeriod) (*ComparisonResult, error)

    // 实时统计
    GetRealTimeStats(ctx context.Context) (*RealTimeStats, error)

    // 数据预测
    GetAnalyticsPredictions(ctx context.Context, predictionType string, days int) (*PredictionResult, error)
}
```

### AuditLogService

审计日志服务接口，负责审计日志的记录和查询。

```go
type AuditLogService interface {
    // 记录审计日志
    CreateAuditLog(ctx context.Context, log *AuditLog) error

    // 查询审计日志
    GetAuditTrail(ctx context.Context, filter AuditLogFilter) (*AuditLogPage, error)

    // 获取资源审计日志
    GetAuditTrailByResource(ctx context.Context, resourceType, resourceID string) ([]*AuditLog, error)

    // 获取操作审计日志
    GetAuditTrailByAction(ctx context.Context, action string) ([]*AuditLog, error)

    // 获取审计统计
    GetAuditStatistics(ctx context.Context) (*AuditStatistics, error)

    // 导出审计日志
    ExportAuditTrail(ctx context.Context, filter AuditLogFilter) ([]byte, error)
}
```

### ExportService

导出服务接口，处理各类数据导出功能。

```go
type ExportService interface {
    // 导出书籍为CSV
    ExportBooksToCSV(ctx context.Context, filter *BookFilter) ([]byte, error)

    // 导出书籍为Excel
    ExportBooksToExcel(ctx context.Context, filter *BookFilter) ([]byte, error)

    // 导出章节为CSV
    ExportChaptersToCSV(ctx context.Context, bookID string) ([]byte, error)

    // 导出用户数据
    ExportUsers(ctx context.Context, filter *UserFilter) ([]byte, error)

    // 获取导出历史
    GetExportHistory(ctx context.Context, userID string) ([]*ExportHistory, error)
}
```

### SensitiveOperationService

敏感操作服务接口，处理敏感操作的记录和验证。

```go
type SensitiveOperationService interface {
    // 记录敏感操作
    RecordOperation(ctx context.Context, op *SensitiveOperation) error

    // 验证敏感操作
    ValidateOperation(ctx context.Context, operationID string) (bool, error)

    // 获取用户敏感操作历史
    GetUserOperations(ctx context.Context, userID string) ([]*SensitiveOperation, error)
}
```

## 数据模型

### AnalyticsData

```go
type UserGrowthData struct {
    Periods   []string `json:"periods"`
    NewUsers  []int64  `json:"newUsers"`
    ActiveUsers []int64 `json:"activeUsers"`
    TotalUsers int64   `json:"totalUsers"`
}

type ContentStatistics struct {
    TotalBooks     int64 `json:"totalBooks"`
    PublishedBooks int64 `json:"publishedBooks"`
    DraftBooks     int64 `json:"draftBooks"`
    TotalChapters  int64 `json:"totalChapters"`
    TotalWords     int64 `json:"totalWords"`
}

type RevenueReport struct {
    Period      string  `json:"period"`
    TotalRevenue float64 `json:"totalRevenue"`
    Subscriptions int64  `json:"subscriptions"`
    Purchases     int64  `json:"purchases"`
}
```

### AuditLog

```go
type AuditLog struct {
    ID          string    `json:"id" bson:"_id"`
    UserID      string    `json:"userId" bson:"user_id"`
    Action      string    `json:"action" bson:"action"`
    ResourceType string   `json:"resourceType" bson:"resource_type"`
    ResourceID  string    `json:"resourceId" bson:"resource_id"`
    Details     string    `json:"details" bson:"details"`
    IPAddress   string    `json:"ipAddress" bson:"ip_address"`
    UserAgent   string    `json:"userAgent" bson:"user_agent"`
    CreatedAt   time.Time `json:"createdAt" bson:"created_at"`
}
```

### ExportHistory

```go
type ExportHistory struct {
    ID          string    `json:"id" bson:"_id"`
    UserID      string    `json:"userId" bson:"user_id"`
    ExportType  string    `json:"exportType" bson:"export_type"`
    Format      string    `json:"format" bson:"format"`
    Status      string    `json:"status" bson:"status"`
    FileURL     string    `json:"fileUrl" bson:"file_url"`
    RecordCount int       `json:"recordCount" bson:"record_count"`
    CreatedAt   time.Time `json:"createdAt" bson:"created_at"`
    CompletedAt *time.Time `json:"completedAt,omitempty" bson:"completed_at,omitempty"`
}
```

### SensitiveOperation

```go
type SensitiveOperation struct {
    ID          string    `json:"id" bson:"_id"`
    UserID      string    `json:"userId" bson:"user_id"`
    Operation   string    `json:"operation" bson:"operation"`
    TargetType  string    `json:"targetType" bson:"target_type"`
    TargetID    string    `json:"targetId" bson:"target_id"`
    Reason      string    `json:"reason" bson:"reason"`
    Approved    bool      `json:"approved" bson:"approved"`
    CreatedAt   time.Time `json:"createdAt" bson:"created_at"`
}
```

## 技术特点

### 1. 数据分析
- 多维度统计分析
- 时间序列分析
- 趋势预测
- 实时数据处理

### 2. 数据导出
- 多格式支持（CSV/Excel）
- 大数据量处理
- 异步导出
- 导出历史记录

### 3. 审计追踪
- 完整的操作记录
- 敏感操作追踪
- 审计日志查询
- 审计数据导出

### 4. 性能优化
- 数据聚合缓存
- 分页查询
- 异步处理
- 批量操作

## 测试覆盖

主要测试文件：
- `analytics_service_test.go` - 统计分析服务测试
- `audit_log_service_test.go` - 审计日志服务测试
- `export_service_test.go` - 导出服务测试
- `sensitive_operation_service_test.go` - 敏感操作服务测试

## 依赖关系

### 依赖的模块
- `models/admin` - 管理员数据模型
- `models/audit` - 审计日志模型
- `repository` - 数据仓储层
- `pkg/metrics` - 指标收集

### 被依赖的模块
- `api/v1/admin` - 管理员API层

## 使用示例

### 统计分析示例

```go
// 获取用户增长趋势
analytics := admin.NewAnalyticsService(repo, metrics)
data, err := analytics.GetUserGrowthTrend(ctx, "2026-01-01", "2026-01-31", "daily")
if err != nil {
    return err
}
```

### 审计日志示例

```go
// 记录审计日志
auditLog := admin.NewAuditLogService(repo)
log := &admin.AuditLog{
    UserID:       userID,
    Action:       "user.update",
    ResourceType: "user",
    ResourceID:   targetUserID,
    Details:      "更新用户状态",
}
err := auditLog.CreateAuditLog(ctx, log)
```

### 导出服务示例

```go
// 导出书籍数据
exportService := admin.NewExportService(repo, storage)
data, err := exportService.ExportBooksToCSV(ctx, &admin.BookFilter{
    Status: &[]string{"published"}[0],
})
```

## 相关文档

- [Admin API](../../api/v1/admin/README.md)
- [审计系统设计](../../docs/design/modules/07-admin/README.md)
- [统计分析API文档](../../docs/api/analytics.md)

---

**版本**: v1.0
**更新日期**: 2026-02-27
**维护者**: Admin Service模块开发组
