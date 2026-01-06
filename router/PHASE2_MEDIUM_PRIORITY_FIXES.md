# 第二阶段：中优先级TODO修复计划

**开始时间**: 2025-10-31
**目标**: 完成15个中优先级TODO
**预计耗时**: 2-4周
**状态**: 开始执行

---

## 📋 修复优先级排序

### 第1组 - BookStore API补全 (3项) - 优先级: 高
**目标**: 启用Rating和Chapter API，完善BookStore功能

| 序号 | 任务 | 文件 | 估时 | 状态 |
|------|------|------|------|------|
| 1 | 启用Rating API | router/bookstore/bookstore_router.go | 20分钟 | ⏳ 待做 |
| 2 | 启用Chapter API | router/bookstore/bookstore_router.go | 15分钟 | ⏳ 待做 |
| 3 | BookStore类型完善 | 完成type定义 | 10分钟 | ⏳ 待做 |

### 第2组 - Admin API实现 (6项) - 优先级: 中
**目标**: 实现系统管理功能

| 序号 | 任务 | 文件 | 估时 | 状态 |
|------|------|------|------|------|
| 4 | 系统统计API | api/v1/admin/system_admin_api.go | 30分钟 | ⏳ 待做 |
| 5 | 系统配置API读取 | api/v1/admin/system_admin_api.go | 20分钟 | ⏳ 待做 |
| 6 | 系统配置API更新 | api/v1/admin/system_admin_api.go | 20分钟 | ⏳ 待做 |
| 7 | 公告管理API | api/v1/admin/system_admin_api.go | 25分钟 | ⏳ 待做 |
| 8 | 审核统计API | api/v1/admin/audit_admin_api.go | 20分钟 | ⏳ 待做 |
| 9 | 用户信息扩展 | api/v1/admin/user_admin_api.go | 15分钟 | ⏳ 待做 |

### 第3组 - 其他API功能 (6项) - 优先级: 中
**目标**: 完善Reader、AI系统等模块

| 序号 | 任务 | 文件 | 估时 | 状态 |
|------|------|------|------|------|
| 10 | 删除阅读进度 | api/v1/reader/books_api.go | 15分钟 | ⏳ 待做 |
| 11 | AI提供商列表 | api/v1/ai/system_api.go | 20分钟 | ⏳ 待做 |
| 12 | AI模型列表 | api/v1/ai/system_api.go | 20分钟 | ⏳ 待做 |
| 13 | Audit权限检查 | api/v1/writer/audit_api.go | 15分钟 | ⏳ 待做 |
| 14 | 权限逻辑 | api/v1/writer/audit_api.go | 15分钟 | ⏳ 待做 |
| 15 | 操作日志记录 | 相关API文件 | 20分钟 | ⏳ 待做 |

---

## 🎯 今日目标：完成第1组 (BookStore API补全)

### 任务1: 启用Rating API

**文件**: `router/bookstore/bookstore_router.go:60-62`

**当前状态**: 被注释
```go
// if ratingService != nil {
// 	ratingApiHandler := bookstoreApi.NewBookRatingAPI(ratingService.(bookstore.RatingService))
// }
```

**修复方案**:
1. 初始化RatingAPI处理器
2. 启用Rating相关路由
3. 验证编译通过

**预计时间**: 20分钟

---

### 任务2: 启用Chapter API

**文件**: `router/bookstore/bookstore_router.go:66`

**当前状态**: 被注释
```go
// chapterApiHandler := bookstoreApi.NewChapterAPI(...)
```

**修复方案**:
1. 实现Chapter API初始化
2. 启用Chapter相关路由
3. 验证编译通过

**预计时间**: 15分钟

---

### 任务3: 完善BookStore类型定义

**文件**: `router/bookstore/bookstore_router.go:49-50`

**当前状态**: 仍为interface{}
```go
ratingService interface{},
statisticsService interface{},
```

**修复方案**:
1. 更改为具体类型
2. 删除类型断言
3. 验证编译通过

**预计时间**: 10分钟

---

## 💡 实施步骤

### Step 1: 启用Rating API处理器

```go
// 初始化Rating API处理器
var ratingApiHandler *bookstoreApi.BookRatingAPI
if ratingService != nil {
	ratingApiHandler = bookstoreApi.NewBookRatingAPI(ratingService)
}
```

### Step 2: 启用Rating相关路由

```go
// Rating API路由（需要认证）
if ratingApiHandler != nil {
	authenticated.GET("/books/:id/rating", ratingApiHandler.GetBookRating)
	authenticated.POST("/books/:id/rating", ratingApiHandler.CreateRating)
	authenticated.PUT("/books/:id/rating", ratingApiHandler.UpdateRating)
	authenticated.DELETE("/books/:id/rating", ratingApiHandler.DeleteRating)
	authenticated.GET("/ratings/user/:id", ratingApiHandler.GetRatingsByUserID)
}
```

### Step 3: 完善Chapter API

类似的步骤初始化和注册Chapter API

### Step 4: 修改类型定义

```go
ratingService bookstore.RatingService,
statisticsService bookstore.StatisticsService,
```

---

## ✅ 验证清单

- [ ] 启用Rating API处理器
- [ ] 启用Rating路由注册
- [ ] 启用Chapter API处理器
- [ ] 启用Chapter路由注册
- [ ] 修改类型定义（ratingService）
- [ ] 修改类型定义（statisticsService）
- [ ] 编译验证通过
- [ ] 单元测试通过

---

## 📊 预计进度

**本周目标**: 完成第1-2组 (9项) 
- BookStore API (3项)
- Admin API主要功能 (6项)

**预计完成**: 2-3天内

**下周目标**: 完成第3组 (6项)
- Reader、AI等模块功能

---

## 🚀 开始第一个修复

准备开始修复 **Task 1: 启用Rating API**

预计耗时: 20分钟
目标: 修复 router/bookstore/bookstore_router.go

---

**计划更新时间**: 2025-10-31
**状态**: 准备执行第一个修复
