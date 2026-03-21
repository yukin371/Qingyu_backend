# Router 层快速参考

## 职责

路由分发层，负责路由注册、路由分组、中间件应用、服务注入。

## 目录结构

```
router/
├── enter.go                # 主入口，注册所有路由
├── admin/                  # 管理后台路由
├── ai/                     # AI服务路由
├── bookstore/              # 书城路由
├── finance/                # 财务路由
├── reader/                 # 阅读器路由
├── recommendation/         # 推荐系统路由
├── shared/                 # 共享服务路由
├── social/                 # 社交路由
├── system/                 # 系统监控路由
├── user/                   # 用户路由
└── writer/                 # 作家路由
```

## 命名规范

| 类型 | 规范 | 示例 |
|------|------|------|
| 文件 | `{模块}_router.go` | `bookstore_router.go` |
| 函数 | `Register{模块}Routes` | `RegisterBookstoreRoutes` |
| 初始化 | `Init{模块}Router` | `InitBookstoreRouter` |

## 快速示例

```go
// 初始化路由
func InitBookstoreRouter(
    v1 *gin.RouterGroup,
    bookstoreSvc bookstore.BookstoreService,
    logger *zap.Logger,
) {
    api := bookstoreAPI.NewBookstoreAPI(bookstoreSvc, logger)

    bookstoreGroup := v1.Group("/bookstore")
    {
        bookstoreGroup.GET("/homepage", api.GetHomepage)

        books := bookstoreGroup.Group("/books")
        {
            books.GET("", api.GetBooks)
            books.GET("/:id", api.GetBookByID)
            books.POST("", api.CreateBook)
        }
    }
}
```

## 路由分组模式

```go
// 公开路由组
public := v1.Group("")
{
    public.POST("/auth/login", authAPI.Login)
}

// 认证路由组
authenticated := v1.Group("")
authenticated.Use(authMiddleware)
{
    authenticated.GET("/profile", userAPI.GetProfile)
}

// 管理员路由组
admin := v1.Group("/admin")
admin.Use(authMiddleware, adminMiddleware)
{
    admin.GET("/users", adminAPI.ListUsers)
}
```

## 渐进式注册

```go
// 按可用服务逐个注册
if authErr == nil && authSvc != nil {
    sharedRouter.RegisterAuthRoutes(sharedGroup, authSvc, oauthSvc, logger)
    logger.Info("✓ 认证服务路由已注册")
} else {
    logger.Warn("⚠ AuthService未配置，跳过认证路由注册")
}
```

## RESTful 规范

| HTTP 方法 | 用途 | 示例 |
|-----------|------|------|
| GET | 获取资源 | `GET /books`, `GET /books/:id` |
| POST | 创建/操作 | `POST /books`, `POST /books/:id/publish` |
| PUT | 完整更新 | `PUT /books/:id` |
| DELETE | 删除 | `DELETE /books/:id` |

## URL 路径规范

```go
// ✅ 推荐：复数名词 + kebab-case
v1.GET("/books", ...)
v1.GET("/reading-stats", ...)

// ❌ 避免：动词
v1.GET("/getBooks", ...)  // 不推荐
```

## 禁止事项

- ❌ 在路由中编写业务逻辑
- ❌ 直接操作数据库
- ❌ 硬编码服务实例

## 详见

完整设计文档: [docs/standards/layer-router.md](../docs/standards/layer-router.md)
