# Router 层设计说明

> **版本**: v1.0
> **创建日期**: 2026-03-19
> **状态**: 进行中

---

## 1. 职责边界与依赖关系

### 1.1 职责定义

Router 层是后端系统的**路由分发层**，负责：

1. **路由注册**：将 URL 路径映射到 API 处理器
2. **路由分组**：按业务模块组织路由
3. **中间件应用**：为路由组配置中间件
4. **服务注入**：将 Service 注入到 API 处理器

### 1.2 与上下层级的交互边界

```
┌─────────────────────────────────────────────────────────┐
│                   Middleware 层                         │
│              (请求预处理)                               │
└─────────────────────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────┐
│                    Router 层                            │
│  ┌─────────────────────────────────────────────────────┤
│  │ 职责：                                              │
│  │ - URL 路径匹配                                      │
│  │ - 路由分组管理                                      │
│  │ - 中间件应用                                        │
│  │ - Service 注入                                      │
│  │ - 路由注册入口                                      │
│  └─────────────────────────────────────────────────────┤
│  禁止：                                                  │
│  │ - 编写业务逻辑                                      │
│  │ - 直接操作数据库                                    │
│  │ - 处理请求参数                                      │
└─────────────────────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────┐
│                     API 层                              │
│              (请求处理)                                 │
└─────────────────────────────────────────────────────────┘
```

### 1.3 依赖关系

```go
// Router 层允许的依赖
import (
    "github.com/gin-gonic/gin"
    "Qingyu_backend/api/v1/xxx"         // API 处理器
    "Qingyu_backend/service/xxx"        // Service 接口
    "Qingyu_backend/service/container"  // 服务容器
    "Qingyu_backend/pkg/logger"         // 日志
)

// Router 层禁止的依赖
import (
    "Qingyu_backend/repository/xxx"     // ❌ 禁止直接依赖 Repository
    "go.mongodb.org/mongo-driver"       // ❌ 禁止直接操作数据库
)
```

---

## 2. 命名与代码规范

### 2.1 文件命名规范

| 类型 | 规范 | 示例 |
|------|------|------|
| 路由注册 | `{模块}_router.go` | `bookstore_router.go`, `user_router.go` |
| 主入口 | `enter.go` | `enter.go` |
| 测试文件 | `{文件名}_test.go` | `bookstore_router_test.go` |

### 2.2 函数命名规范

| 类型 | 规范 | 示例 |
|------|------|------|
| 路由注册 | `Register{模块}Routes` | `RegisterBookstoreRoutes`, `RegisterUserRoutes` |
| 初始化 | `Init{模块}Router` | `InitBookstoreRouter`, `InitAIRouter` |
| 主入口 | `RegisterRoutes` | `RegisterRoutes` |

### 2.3 目录组织规范

```
router/
├── enter.go                # 主入口，注册所有路由
├── admin/                  # 管理后台路由
│   └── admin_router.go
├── ai/                     # AI服务路由
│   ├── ai_router.go
│   └── creative.go
├── bookstore/              # 书城路由
│   └── bookstore_router.go
├── finance/                # 财务路由
│   └── finance_router.go
├── reader/                 # 阅读器路由
│   └── reader_router.go
├── recommendation/         # 推荐系统路由
│   └── recommendation_router.go
├── shared/                 # 共享服务路由
│   ├── shared_router.go
│   ├── storage_router.go
│   └── health.go
├── social/                 # 社交路由
│   └── social_router.go
├── system/                 # 系统监控路由
│   └── system_router.go
├── user/                   # 用户路由
│   └── user_router.go
└── writer/                 # 作家路由
    ├── writer_main_router.go
    ├── project_router.go
    ├── document_router.go
    └── publish_router.go
```

---

## 3. 设计模式与最佳实践

### 3.1 路由注册模式

```go
// bookstore_router.go
package bookstore

import (
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"

    bookstoreAPI "Qingyu_backend/api/v1/bookstore"
    "Qingyu_backend/service/bookstore"
)

// InitBookstoreRouter 初始化书城路由
func InitBookstoreRouter(
    v1 *gin.RouterGroup,
    bookstoreSvc bookstore.BookstoreService,
    bookDetailSvc bookstore.BookDetailService,
    ratingSvc bookstore.BookRatingService,
    statisticsSvc bookstore.BookStatisticsService,
    chapterSvc bookstore.ChapterService,
    chapterPurchaseSvc bookstore.ChapterPurchaseService,
    searchSvc *searchService.SearchService,
    logger *zap.Logger,
) {
    // 创建 API 处理器
    bookstoreAPI := bookstoreAPI.NewBookstoreAPI(bookstoreSvc, searchSvc, logger)

    // 创建路由组
    bookstoreGroup := v1.Group("/bookstore")
    {
        // 首页
        bookstoreGroup.GET("/homepage", bookstoreAPI.GetHomepage)

        // 书籍
        books := bookstoreGroup.Group("/books")
        {
            books.GET("", bookstoreAPI.GetBooks)
            books.GET("/search", bookstoreAPI.SearchBooks)
            books.GET("/recommended", bookstoreAPI.GetRecommendedBooks)
            books.GET("/featured", bookstoreAPI.GetFeaturedBooks)
            books.GET("/:id", bookstoreAPI.GetBookByID)
            books.GET("/:id/similar", bookstoreAPI.GetSimilarBooks)
            books.POST("/:id/view", bookstoreAPI.IncrementBookView)
        }

        // 分类
        categories := bookstoreGroup.Group("/categories")
        {
            categories.GET("/tree", bookstoreAPI.GetCategoryTree)
            categories.GET("/:id", bookstoreAPI.GetCategoryByID)
            categories.GET("/:id/books", bookstoreAPI.GetBooksByCategory)
        }

        // 榜单
        rankings := bookstoreGroup.Group("/rankings")
        {
            rankings.GET("/realtime", bookstoreAPI.GetRealtimeRanking)
            rankings.GET("/weekly", bookstoreAPI.GetWeeklyRanking)
            rankings.GET("/monthly", bookstoreAPI.GetMonthlyRanking)
            rankings.GET("/:type", bookstoreAPI.GetRankingByType)
        }
    }

    logger.Info("✓ 书店路由已注册到: /api/v1/bookstore/")
}
```

### 3.2 路由分组模式

```go
// 按业务模块分组
v1 := r.Group("/api/v1")

// 公开路由组（无需认证）
public := v1.Group("")
{
    public.POST("/auth/login", authAPI.Login)
    public.POST("/auth/register", authAPI.Register)
    public.GET("/books", bookAPI.ListBooks)
}

// 认证路由组（需要登录）
authenticated := v1.Group("")
authenticated.Use(authMiddleware)
{
    authenticated.GET("/profile", userAPI.GetProfile)
    authenticated.PUT("/profile", userAPI.UpdateProfile)
}

// 管理员路由组（需要管理员权限）
admin := v1.Group("/admin")
admin.Use(authMiddleware, adminMiddleware)
{
    admin.GET("/users", adminAPI.ListUsers)
    admin.DELETE("/users/:id", adminAPI.DeleteUser)
}
```

### 3.3 服务注入模式

```go
// 方式1：通过参数注入（推荐）
func RegisterBookstoreRoutes(
    v1 *gin.RouterGroup,
    bookstoreSvc bookstore.BookstoreService,
    logger *zap.Logger,
) {
    api := bookstoreAPI.NewBookstoreAPI(bookstoreSvc, logger)
    // 注册路由...
}

// 方式2：通过服务容器获取
func RegisterRoutes(r *gin.Engine) {
    serviceContainer := service.GetServiceContainer()

    bookstoreSvc, err := serviceContainer.GetBookstoreService()
    if err != nil {
        logger.Fatal("获取书店服务失败", zap.Error(err))
    }

    bookstoreRouter.InitBookstoreRouter(v1, bookstoreSvc, logger)
}
```

### 3.4 渐进式注册模式

```go
// 按可用服务逐个注册（推荐）
registeredCount := 0

// 1. 注册认证服务路由
if authErr == nil && authSvc != nil {
    sharedRouter.RegisterAuthRoutes(sharedGroup, authSvc, oauthSvc, logger)
    logger.Info("✓ 认证服务路由已注册")
    registeredCount++
} else {
    logger.Warn("⚠ AuthService未配置，跳过认证路由注册")
}

// 2. 注册存储服务路由
if storageErr == nil && storageSvc != nil {
    sharedRouter.RegisterStorageRoutes(sharedGroup, storageSvc)
    logger.Info("✓ 存储服务路由已注册")
    registeredCount++
} else {
    logger.Warn("⚠ StorageService未配置，跳过存储路由注册")
}

// 总结
if registeredCount > 0 {
    logger.Info(fmt.Sprintf("✓ 已注册 %d 个服务模块", registeredCount))
}
```

### 3.5 RESTful 路由规范

```go
// RESTful 资源路由
books := v1.Group("/books")
{
    // 列表
    books.GET("", bookAPI.List)           // GET /api/v1/books

    // 创建
    books.POST("", bookAPI.Create)        // POST /api/v1/books

    // 单个资源
    books.GET("/:id", bookAPI.GetByID)    // GET /api/v1/books/:id
    books.PUT("/:id", bookAPI.Update)     // PUT /api/v1/books/:id
    books.DELETE("/:id", bookAPI.Delete)  // DELETE /api/v1/books/:id

    // 子资源
    books.GET("/:id/chapters", bookAPI.ListChapters)  // GET /api/v1/books/:id/chapters

    // 自定义操作（使用动词）
    books.POST("/:id/publish", bookAPI.Publish)      // POST /api/v1/books/:id/publish
    books.POST("/:id/archive", bookAPI.Archive)      // POST /api/v1/books/:id/archive
}
```

### 3.6 反模式警示

❌ **禁止事项**：

```go
// ❌ 禁止：在路由中编写业务逻辑
func RegisterBookRoutes(v1 *gin.RouterGroup) {
    v1.GET("/books/:id", func(c *gin.Context) {
        // 业务逻辑应该在 API 层
        book, _ := db.Collection("books").FindOne(...)
        c.JSON(200, book)
    })
}

// ❌ 禁止：直接操作数据库
func RegisterRoutes(r *gin.Engine) {
    mongoClient := mongo.Connect(...)  // 禁止在路由中初始化
}

// ❌ 禁止：硬编码服务实例
func RegisterRoutes(r *gin.Engine) {
    bookstoreSvc := bookstore.NewBookstoreService(...)  // 应该通过容器获取
}
```

---

## 4. 路由规范

### 4.1 URL 路径规范

| 类型 | 规范 | 示例 |
|------|------|------|
| 资源列表 | `/{资源}` | `/books`, `/users` |
| 单个资源 | `/{资源}/:id` | `/books/:id` |
| 子资源 | `/{资源}/:id/{子资源}` | `/books/:id/chapters` |
| 自定义操作 | `/{资源}/:id/{动词}` | `/books/:id/publish` |
| 搜索/过滤 | `/{资源}/search` | `/books/search` |

### 4.2 HTTP 方法规范

| HTTP 方法 | 用途 | 示例 |
|-----------|------|------|
| GET | 获取资源 | `GET /books`, `GET /books/:id` |
| POST | 创建资源/执行操作 | `POST /books`, `POST /books/:id/publish` |
| PUT | 完整更新 | `PUT /books/:id` |
| PATCH | 部分更新 | `PATCH /books/:id` |
| DELETE | 删除资源 | `DELETE /books/:id` |

### 4.3 路由命名规范

```go
// ✅ 推荐：使用复数名词
v1.GET("/books", ...)           // 书籍列表
v1.GET("/users", ...)           // 用户列表
v1.GET("/chapters", ...)        // 章节列表

// ✅ 推荐：使用 kebab-case
v1.GET("/reading-stats", ...)   // 阅读统计
v1.GET("/user-management", ...) // 用户管理

// ❌ 避免：使用动词
v1.GET("/getBooks", ...)        // 不推荐
v1.POST("/createBook", ...)     // 不推荐
```

---

## 5. 测试策略

### 5.1 路由测试编写指南

```go
// bookstore_router_test.go
package bookstore

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

func TestBookstoreRouter_Routes(t *testing.T) {
    // 1. 准备
    gin.SetMode(gin.TestMode)

    mockService := new(MockBookstoreService)
    mockService.On("GetHomepageData", mock.Anything).
        Return(&bookstore.HomepageData{}, nil)

    // 2. 创建路由
    r := gin.New()
    v1 := r.Group("/api/v1")

    api := bookstoreAPI.NewBookstoreAPI(mockService, nil, nil)
    bookstoreGroup := v1.Group("/bookstore")
    {
        bookstoreGroup.GET("/homepage", api.GetHomepage)
    }

    // 3. 测试路由
    w := httptest.NewRecorder()
    req := httptest.NewRequest("GET", "/api/v1/bookstore/homepage", nil)
    r.ServeHTTP(w, req)

    // 4. 验证
    assert.Equal(t, http.StatusOK, w.Code)
}

func TestBookstoreRouter_RouteNotFound(t *testing.T) {
    gin.SetMode(gin.TestMode)

    r := gin.New()
    v1 := r.Group("/api/v1")

    // 注册路由...

    w := httptest.NewRecorder()
    req := httptest.NewRequest("GET", "/api/v1/bookstore/nonexistent", nil)
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusNotFound, w.Code)
}
```

### 5.2 测试场景清单

- [ ] 路由路径正确
- [ ] HTTP 方法正确
- [ ] 路由分组正确
- [ ] 中间件应用正确
- [ ] 404 路由处理

---

## 6. 完整代码示例

### 6.1 完整路由注册示例

```go
// enter.go
package router

import (
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"

    "Qingyu_backend/service"
    "Qingyu_backend/service/container"
    adminRouter "Qingyu_backend/router/admin"
    bookstoreRouter "Qingyu_backend/router/bookstore"
    userRouter "Qingyu_backend/router/user"
)

// RegisterRoutes 注册所有路由
func RegisterRoutes(r *gin.Engine) {
    // 初始化日志器
    logger := initRouterLogger()

    // API 版本组
    v1 := r.Group("/api/v1")

    // 获取全局服务容器
    serviceContainer := service.GetServiceContainer()
    if serviceContainer == nil {
        logger.Fatal("服务容器未初始化")
    }

    logger.Info("✓ 服务容器已初始化，开始注册路由...")

    // ============ 注册书店路由 ============
    registerBookstoreRoutes(v1, serviceContainer, logger)

    // ============ 注册用户路由 ============
    registerUserRoutes(v1, serviceContainer, logger)

    // ============ 注册管理员路由 ============
    registerAdminRoutes(v1, serviceContainer, logger)

    // ============ 健康检查 ============
    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "pong"})
    })

    logger.Info("✓ 所有路由注册完成!")
}

// registerBookstoreRoutes 注册书店路由
func registerBookstoreRoutes(v1 *gin.RouterGroup, container *container.ServiceContainer, logger *zap.Logger) {
    bookstoreSvc, err := container.GetBookstoreService()
    if err != nil {
        logger.Warn("获取书店服务失败", zap.Error(err))
        return
    }

    // 获取其他依赖服务
    bookDetailSvc, _ := container.GetBookDetailService()
    ratingSvc, _ := container.GetBookRatingService()
    statisticsSvc, _ := container.GetBookStatisticsService()
    chapterSvc, _ := container.GetChapterService()

    // 注册路由
    bookstoreRouter.InitBookstoreRouter(
        v1,
        bookstoreSvc,
        bookDetailSvc,
        ratingSvc,
        statisticsSvc,
        chapterSvc,
        nil,  // chapterPurchaseSvc
        nil,  // searchSvc
        logger,
    )

    logger.Info("✓ 书店路由已注册到: /api/v1/bookstore/")
}

// registerUserRoutes 注册用户路由
func registerUserRoutes(v1 *gin.RouterGroup, container *container.ServiceContainer, logger *zap.Logger) {
    userSvc, err := container.GetUserService()
    if err != nil {
        logger.Warn("获取用户服务失败", zap.Error(err))
        return
    }

    userRouter.RegisterUserRoutes(v1, userSvc, nil, nil, nil)
    logger.Info("✓ 用户路由已注册到: /api/v1/user/")
}

// registerAdminRoutes 注册管理员路由
func registerAdminRoutes(v1 *gin.RouterGroup, container *container.ServiceContainer, logger *zap.Logger) {
    userSvc, err := container.GetUserService()
    if err != nil {
        logger.Warn("获取用户服务失败", zap.Error(err))
        return
    }

    quotaSvc, _ := container.GetQuotaService()
    auditSvc, _ := container.GetAuditService()

    adminRouter.RegisterAdminRoutes(
        v1,
        userSvc,
        quotaSvc,
        auditSvc,
        nil,  // adminSvc
        nil,  // configSvc
        nil,  // announcementSvc
        nil,  // userAdminSvc
        nil,  // permissionSvc
        container.GetPersistedEventBus(),
        nil,  // categorySvc
        nil,  // publicationSvc
        nil,  // bannerSvc
    )

    logger.Info("✓ 管理员路由已注册到: /api/v1/admin/")
}
```

---

## 7. 参考资料

- [Router 层快速参考](../router/README.md)
- [API 层设计说明](./layer-api.md)
- [Middleware 层设计说明](./layer-middleware.md)
- [Service 层设计说明](./layer-service.md)

---

*最后更新：2026-03-19*
