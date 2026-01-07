# 路由层TODO管理与完善计划

## 📋 执行摘要

**总TODO项**: 37项
**优先级高**: 6项
**优先级中**: 15项
**优先级低**: 16项

---

## 🔴 优先级高 - 立即处理 (6项)

### 1. Router层 - 关键TODO

#### 1.1 审核服务实现
**位置**: `router/enter.go:255`
**问题**: TODO: 获取审核服务实例（需要实现）
**影响**: 管理员审核功能无法使用
**优先级**: 🔴 高

**当前代码**:
```go
// TODO: 获取审核服务实例（需要实现）
// auditSvc := serviceContainer.GetAuditService()
adminRouter.RegisterAdminRoutes(v1, userSvc, quotaService, nil, adminSvc, configSvc)
```

**建议方案**:
1. 在ServiceContainer中添加GetAuditService()方法
2. 实现AuditService（如果还没有）
3. 在router/enter.go中取消注释并传递auditSvc

---

### 2. BookStore路由 - 类型定义TODO

#### 2.1 BookDetailService类型定义
**位置**: `router/bookstore/bookstore_router.go:48`
**问题**: TODO: 改为具体类型 (当前为interface{})
**影响**: 类型检查不完全，运行时可能出错
**优先级**: 🔴 高

**当前代码**:
```go
func InitBookstoreRouter(
	r *gin.RouterGroup,
	bookstoreService bookstore.BookstoreService,
	bookDetailService interface{}, // TODO: 改为具体类型
	...
)
```

**建议方案**:
- 定义具体的BookDetailService接口
- 更改函数签名: `bookDetailService bookstore.BookDetailService`
- 更新调用处的参数传递

#### 2.2 RatingService类型定义
**位置**: `router/bookstore/bookstore_router.go:49`
**问题**: TODO: 改为具体类型
**优先级**: 🔴 高

#### 2.3 StatisticsService类型定义
**位置**: `router/bookstore/bookstore_router.go:50`
**问题**: TODO: 改为具体类型
**优先级**: 🔴 高

---

### 3. Writer路由 - 权限中间件TODO

#### 3.1 审核路由权限检查
**位置**: `router/writer/audit.go:44`
**问题**: TODO: 添加管理员权限中间件
**影响**: 审核端点可能被普通用户访问
**优先级**: 🔴 高

**当前代码**:
```go
// adminGroup.Use(middleware.AdminPermission()) // TODO: 添加管理员权限中间件
```

**建议方案**:
- 启用管理员权限检查
- 确保审核功能只有管理员可以访问

---

## 🟡 优先级中 - 近期处理 (15项)

### 4. BookStore API - 功能实现TODO

#### 4.1 BookDetail API
**位置**: `router/bookstore/bookstore_router.go:99-101`
**问题**: TODO: 当BookDetailAPI实现后添加
**涉及端点**:
- GET /books/:id/detail
- GET /books/:id/similar
- GET /books/:id/statistics

**建议方案**:
1. 检查api/v1/bookstore/book_detail_api.go是否已实现
2. 实现缺失的方法
3. 在bookstore_router.go中注册这些路由

#### 4.2 Chapter API
**位置**: `router/bookstore/bookstore_router.go:104-106`
**问题**: TODO: 当ChapterAPI实现后添加
**涉及端点**:
- GET /chapters/:id
- GET /chapters/book/:id

#### 4.3 Rating API
**位置**: `router/bookstore/bookstore_router.go:116-121`
**问题**: TODO: 当RatingAPI实现后添加
**涉及端点**:
- GET /books/:id/rating
- POST /books/:id/rating
- PUT /books/:id/rating
- DELETE /books/:id/rating
- GET /ratings/user/:id

---

### 5. Admin API - 功能实现TODO (8项)

#### 5.1 系统统计API
**位置**: `api/v1/admin/system_admin_api.go:157`
**问题**: TODO: 实现系统统计功能
**优先级**: 🟡 中

#### 5.2 系统配置API
**位置**: `api/v1/admin/system_admin_api.go:190`
**问题**: TODO: 实现获取系统配置功能
**优先级**: 🟡 中

#### 5.3 更新系统配置API
**位置**: `api/v1/admin/system_admin_api.go:229`
**问题**: TODO: 实现更新系统配置功能
**优先级**: 🟡 中

#### 5.4 公告管理API
**位置**: `api/v1/admin/system_admin_api.go:263, 288`
**问题**: TODO: 实现发布公告功能 / 获取公告列表功能
**优先级**: 🟡 中

#### 5.5 审核统计API
**位置**: `api/v1/admin/audit_admin_api.go:236`
**问题**: TODO: 实现审核统计功能
**优先级**: 🟡 中

#### 5.6 用户信息扩展
**位置**: `api/v1/admin/user_admin_api.go:324`
**问题**: TODO: 可以扩展添加ban_reason, ban_until等字段
**建议**: 这是可选优化，可后续处理
**优先级**: 🟡 中

---

### 6. Reader API - 功能缺失TODO

#### 6.1 删除阅读进度方法
**位置**: `api/v1/reader/books_api.go:122`
**问题**: TODO: 实现删除阅读进度的方法
**优先级**: 🟡 中

---

### 7. AI系统API - 功能实现TODO

#### 7.1 获取提供商列表
**位置**: `api/v1/ai/system_api.go:52`
**问题**: TODO: 实现获取提供商列表的逻辑
**优先级**: 🟡 中

#### 7.2 获取模型列表
**位置**: `api/v1/ai/system_api.go:80`
**问题**: TODO: 实现获取模型列表的逻辑
**优先级**: 🟡 中

---

### 8. Audit API - 权限检查TODO

#### 8.1 权限检查未实现
**位置**: `api/v1/writer/audit_api.go:279, 318`
**问题**: TODO: 检查是否为管理员
**优先级**: 🟡 中

**当前代码**:
```go
// TODO: 检查是否为管理员
```

**建议方案**:
- 添加权限检查逻辑
- 使用middleware.RequireRole("admin")或在API中检查

---

## 🟢 优先级低 - 后续处理 (16项)

### 9. Phase3功能 - 后续实现 (4项)

#### 9.1 统计API (Phase3)
**位置**: `api/v1/shared/stats_api.go:240`
**优先级**: 🟢 低

#### 9.2 推送通知API (Phase3)
**位置**: `api/v1/shared/notification_api.go:241, 247`
**优先级**: 🟢 低

#### 9.3 搜索历史和热门搜索API (Phase3)
**位置**: `api/v1/shared/search_api.go:166, 178`
**优先级**: 🟢 低

#### 9.4 相似推荐API
**位置**: `api/v1/recommendation/similar.go:11`
**优先级**: 🟢 低

#### 9.5 个人推荐API
**位置**: `api/v1/recommendation/personal.go:11`
**优先级**: 🟢 低

---

### 10. 数据转换和实现TODO (10项)

#### 10.1 Audit API转换逻辑
**位置**: `api/v1/writer/audit_api.go:378, 384, 389`
**问题**: TODO: 实现完整的转换逻辑
**优先级**: 🟢 低

#### 10.2 审核过滤逻辑
**位置**: `api/v1/admin/audit_admin_api.go:66`
**问题**: TODO: 这里应该在Service层实现过滤逻辑
**优先级**: 🟢 低

#### 10.3 Rating API搜索方法
**位置**: `api/v1/bookstore/book_rating_api.go:652`
**问题**: TODO: SearchByKeyword方法尚未在Service层实现
**优先级**: 🟢 低

#### 10.4 User API作品查询
**位置**: `api/v1/user/user_api.go:413`
**问题**: TODO: 调用BookService查询用户的已发布作品
**优先级**: 🟢 低

#### 10.5 分类推荐逻辑
**位置**: `api/v1/recommendation/recommendation_api.go:254`
**问题**: TODO: 后续可以基于category参数实现真正的分类推荐
**优先级**: 🟢 低

#### 10.6 操作日志记录
**位置**: `api/v1/admin/README.md:326`
**问题**: TODO: 记录管理员操作日志
**优先级**: 🟢 低

---

## 📊 优先级分布

```
🔴 高优先级 (6项) - 立即处理
  ├─ 审核服务实现 (1项)
  └─ BookStore类型定义 (3项)
  └─ Writer权限中间件 (1项)
  └─ 其他关键功能 (1项)

🟡 中优先级 (15项) - 近期处理
  ├─ BookStore API功能 (3项)
  ├─ Admin API功能 (8项)
  ├─ Reader API功能 (1项)
  └─ AI系统API功能 (2项)
  └─ 权限检查 (1项)

🟢 低优先级 (16项) - 后续处理
  ├─ Phase3功能 (4项)
  └─ 数据转换优化 (10项)
  └─ 其他优化 (2项)
```

---

## ✅ 完善方案

### 第一阶段 - 关键修复 (1-2周)

#### 1.1 审核服务实现
```go
// 1. 在ServiceContainer中添加
func (c *ServiceContainer) GetAuditService() (audit.ContentAuditService, error) {
    if c.auditService == nil {
        return nil, fmt.Errorf("AuditService未初始化")
    }
    return c.auditService, nil
}

// 2. 在router/enter.go中
auditSvc, _ := serviceContainer.GetAuditService()
adminRouter.RegisterAdminRoutes(v1, userSvc, quotaService, auditSvc, adminSvc, configSvc)
```

#### 1.2 BookStore类型定义
```go
// 修改函数签名
func InitBookstoreRouter(
    r *gin.RouterGroup,
    bookstoreService bookstore.BookstoreService,
    bookDetailService bookstore.BookDetailService,
    ratingService bookstore.RatingService,
    statisticsService bookstore.StatisticsService,
)

// 或使用指针类型
func InitBookstoreRouter(
    r *gin.RouterGroup,
    bookstoreService bookstore.BookstoreService,
    bookDetailService *bookstore.BookDetailService,
    ...
)
```

#### 1.3 Writer权限检查
```go
// 在audit.go中启用权限检查
adminGroup := r.Group("/audit")
adminGroup.Use(middleware.JWTAuth())
adminGroup.Use(middleware.RequireRole("admin"))
{
    // 路由定义
}
```

---

### 第二阶段 - 功能完善 (2-4周)

#### 2.1 BookStore API补全
```go
// bookstore_router.go中取消注释并完成

if bookDetailService != nil {
    bookDetailApiHandler := bookstoreApi.NewBookDetailAPI(bookDetailService)
    public.GET("/books/:id/detail", bookDetailApiHandler.GetBookDetail)
    public.GET("/books/:id/similar", bookDetailApiHandler.GetSimilarBooks)
    public.GET("/books/:id/statistics", bookDetailApiHandler.GetBookStatistics)
}

if chapterApiHandler != nil {
    public.GET("/chapters/:id", chapterApiHandler.GetChapter)
    public.GET("/chapters/book/:id", chapterApiHandler.GetChaptersByBookID)
}

if ratingApiHandler != nil {
    authenticated.GET("/books/:id/rating", ratingApiHandler.GetBookRating)
    authenticated.POST("/books/:id/rating", ratingApiHandler.CreateRating)
    authenticated.PUT("/books/:id/rating", ratingApiHandler.UpdateRating)
    authenticated.DELETE("/books/:id/rating", ratingApiHandler.DeleteRating)
    authenticated.GET("/ratings/user/:id", ratingApiHandler.GetRatingsByUserID)
}
```

#### 2.2 Admin API实现
- 实现系统统计API
- 实现系统配置API
- 实现公告管理API
- 实现审核统计API

#### 2.3 Reader API补全
- 实现删除阅读进度方法

---

### 第三阶段 - 优化完善 (后续)

#### 3.1 数据转换优化
- 完整的转换逻辑实现
- 过滤逻辑优化

#### 3.2 Phase3功能
- 高级统计API
- 推送通知API
- 搜索历史API

#### 3.3 日志和监控
- 操作日志记录
- 性能监控

---

## 📋 实施检查清单

### 立即行动 (本周)
- [ ] 实现审核服务获取方法
- [ ] 修改BookStore路由类型定义
- [ ] 启用Writer权限检查
- [ ] 编译验证无错误

### 近期行动 (2周内)
- [ ] 完成BookDetail API注册
- [ ] 完成Chapter API注册
- [ ] 完成Rating API注册
- [ ] 实现Admin系统统计
- [ ] 实现Admin配置管理
- [ ] 集成测试验证

### 后续完善 (持续)
- [ ] 优化数据转换逻辑
- [ ] 实现Phase3功能
- [ ] 性能测试
- [ ] 文档更新

---

## 📝 文件修改清单

### 需要修改的关键文件

| 文件 | 优先级 | 操作 | 说明 |
|------|--------|------|------|
| router/enter.go | 🔴 高 | 修改 | 启用审核服务 |
| router/bookstore/bookstore_router.go | 🔴 高 | 修改 | 类型定义、路由注册 |
| router/writer/audit.go | 🔴 高 | 修改 | 添加权限检查 |
| api/v1/admin/system_admin_api.go | 🟡 中 | 修改 | 实现系统管理功能 |
| api/v1/admin/audit_admin_api.go | 🟡 中 | 修改 | 实现审核统计功能 |
| api/v1/ai/system_api.go | 🟡 中 | 修改 | 实现提供商和模型列表 |
| api/v1/writer/audit_api.go | 🟡 中 | 修改 | 完整转换逻辑、权限检查 |
| api/v1/reader/books_api.go | 🟡 中 | 修改 | 实现删除阅读进度 |

---

## 🎯 成功标准

- ✅ 所有高优先级TODO完成
- ✅ 编译通过，无警告
- ✅ 所有类型定义正确
- ✅ 中优先级TODO > 50% 完成
- ✅ 路由注册完整
- ✅ 集成测试通过

---

## 相关文档

- 架构设计规范: `doc/architecture/`
- API设计规范: `doc/api/API设计规范.md`
- 工程规范: `doc/engineering/软件工程规范_v2.0.md`

---

**最后更新**: 2025-10-31
**状态**: 待执行
