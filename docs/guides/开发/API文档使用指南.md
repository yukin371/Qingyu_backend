# API文档使用指南

> **文档版本**: v1.0
> **创建日期**: 2026-01-06
> **适用系统**: 青羽写作平台后端

## 📋 概述

青羽平台使用Swagger/OpenAPI规范进行API文档管理。通过swaggo工具自动生成API文档，提供在线的API浏览和测试界面。

## 🎯 访问API文档

### 1. 本地开发环境

启动服务后，访问以下地址查看API文档：

```
http://localhost:8080/swagger/index.html
```

### 2. 生产环境

```
https://api.qingyu.com/swagger/index.html
```

### 3. 查看JSON格式文档

```
http://localhost:8080/swagger/doc.json
```

## 📝 Swagger注解规范

### 1. 基础API注解

```go
// Summary 简短描述
// Description 详细描述
// Tags 标签（用于分组）
// Accept 接受的内容类型
// Produce 返回的内容类型
// Param 参数说明
// Success 成功响应
// Failure 失败响应
// Router 路由路径
// Security 安全认证

// Example: 获取用户信息
// @Summary      获取用户信息
// @Description  根据用户ID获取用户详细信息
// @Tags         用户
// @Accept       json
// @Produce      json
// @Param        user_id   path      string  true  "用户ID"
// @Param        fields    query      string  false  "返回字段（逗号分隔）"  Extensions(id=fields,example=id,username,email)
// @Success      200       {object}  models.User  "成功返回用户信息"
// @Failure      400       {object}  responses.ErrorResponse  "参数错误"
// @Failure      404       {object}  responses.ErrorResponse  "用户不存在"
// @Failure      500       {object}  responses.ErrorResponse  "服务器错误"
// @Router       /api/v1/users/{user_id} [get]
// @Security     Bearer
func (api *UserAPI) GetUser(c *gin.Context) {
    // 实现代码
}
```

### 2. 通用注解格式

| 注解 | 说明 | 示例 |
|------|------|------|
| `@Summary` | API简短描述 | `@Summary 获取用户列表` |
| `@Description` | API详细描述 | `@Description 分页获取用户列表，支持搜索和过滤` |
| `@Tags` | API分组标签 | `@Tags 用户` |
| `@Accept` | 接受的请求类型 | `@Accept json` |
| `@Produce` | 返回的响应类型 | `@Produce json` |
| `@Param` | 参数说明 | 见下方详细说明 |
| `@Success` | 成功响应 | `@Success 200 {object} models.User` |
| `@Failure` | 失败响应 | `@Failure 400 {object} ErrorResponse` |
| `@Router` | 路由路径和方法 | `@Router /api/v1/users [get]` |
| `@Security` | 安全认证 | `@Security Bearer` |

### 3. 参数注解格式

```go
// @Param name param_type data_type required description
//         ├─ name:        参数名
//         ├─ param_type:  参数位置 (path/query/header/body/form)
//         ├─ data_type:   数据类型 (string/int/bool/array/object等)
//         ├─ required:    是否必填 (true/false)
//         └─ description: 参数描述

// Path参数示例
// @Param user_id path string true "用户ID"

// Query参数示例
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10) minimum(1) maximum(100)
// @Param keyword query string false "搜索关键词"

// Header参数示例
// @Param Authorization header string true "Bearer Token"

// Body参数示例
// @Param body body models.CreateUserRequest true "用户信息"
```

### 4. 响应模型定义

```go
// models/response.go
package models

// User 用户模型
type User struct {
    ID       string `json:"id" example:"123"`
    Username string `json:"username" example:"john_doe"`
    Email    string `json:"email" example:"john@example.com"`
    Role     string `json:"role" example:"user"`
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
    Username string `json:"username" binding:"required" example:"john_doe"`
    Email    string `json:"email" binding:"required,email" example:"john@example.com"`
    Password string `json:"password" binding:"required,min=6" example:"password123"`
}

// PageResponse 分页响应
type PageResponse struct {
    Total int64       `json:"total" example:"100"`
    Page  int         `json:"page" example:"1"`
    Size  int         `json:"size" example:"10"`
    Data  interface{} `json:"data"`
}

// ErrorResponse 错误响应
type ErrorResponse struct {
    Code    int    `json:"code" example:"400"`
    Message string `json:"message" example:"参数错误"`
}
```

## 🏷️ API标签定义

主文件中定义的标签（`cmd/server/main.go`）：

```go
// @tag.name 书城
// @tag.description 书城相关接口，包括首页、书籍列表、分类等

// @tag.name 书籍
// @tag.description 书籍详情、搜索、评分等功能

// @tag.name 用户
// @tag.description 用户注册、登录、个人信息管理

// @tag.name 项目
// @tag.description 写作项目管理

// @tag.name 文档
// @tag.description 文档编辑、版本控制

// @tag.name AI辅助
// @tag.description AI写作辅助功能

// @tag.name 钱包
// @tag.description 钱包、充值、提现功能
```

## 📚 完整API示例

### 1. 用户注册

```go
// @Summary      用户注册
// @Description  创建新用户账号
// @Tags         用户
// @Accept       json
// @Produce      json
// @Param        body  body  models.RegisterRequest  true  "注册信息"
// @Success      201   {object}  responses.AuthResponse  "注册成功"
// @Failure      400   {object}  responses.ErrorResponse  "参数错误"
// @Failure      409   {object}  responses.ErrorResponse  "用户已存在"
// @Router       /api/v1/auth/register [post]
func (api *AuthAPI) Register(c *gin.Context) {
    // 实现代码
}
```

### 2. 获取书籍列表

```go
// @Summary      获取书籍列表
// @Description  分页获取书籍列表，支持分类、状态、搜索等过滤条件
// @Tags         书城
// @Accept       json
// @Produce      json
// @Param        page      query     int     false  "页码"        default(1)  minimum(1)
// @Param        size      query     int     false  "每页数量"    default(10) minimum(1) maximum(100)
// @Param        category  query     string  false  "分类ID"
// @Param        status    query     string  false  "状态"        Enums(published, draft, completed)
// @Param        keyword   query     string  false  "搜索关键词"
// @Param        sort_by   query     string  false  "排序字段"    Enums(created_at, updated_at, rating, read_count)
// @Param        order     query     string  false  "排序方向"    Enums(asc, desc)  default(desc)
// @Success      200       {object}  responses.BookListResponse  "成功"
// @Failure      400       {object}  responses.ErrorResponse     "参数错误"
// @Failure      500       {object}  responses.ErrorResponse     "服务器错误"
// @Router       /api/v1/bookstore/books [get]
// @Security     Bearer
func (api *BookStoreAPI) GetBooks(c *gin.Context) {
    // 实现代码
}
```

### 3. 创建书籍

```go
// @Summary      创建书籍
// @Description  创建新书籍
// @Tags         书籍
// @Accept       json
// @Produce      json
// @Param        body  body  models.CreateBookRequest  true  "书籍信息"
// @Success      201   {object}  responses.BookResponse  "创建成功"
// @Failure      400   {object}  responses.ErrorResponse "参数错误"
// @Failure      401   {object}  responses.ErrorResponse "未授权"
// @Failure      500   {object}  responses.ErrorResponse "服务器错误"
// @Router       /api/v1/books [post]
// @Security     Bearer
func (api *BookAPI) CreateBook(c *gin.Context) {
    // 实现代码
}
```

### 4. 获取章节内容

```go
// @Summary      获取章节内容
// @Description  获取指定章节的详细内容
// @Tags         书籍
// @Accept       json
// @Produce      json
// @Param        book_id         path      string  true  "书籍ID"
// @Param        chapter_number  path      int     true  "章节号"
// @Param        fields          query     string  false  "返回字段（逗号分隔）"
// @Success      200             {object}  responses.ChapterResponse  "成功"
// @Failure      400             {object}  responses.ErrorResponse    "参数错误"
// @Failure      404             {object}  responses.ErrorResponse    "章节不存在"
// @Failure      500             {object}  responses.ErrorResponse    "服务器错误"
// @Router       /api/v1/books/{book_id}/chapters/{chapter_number} [get]
// @Security     Bearer
func (api *BookAPI) GetChapter(c *gin.Context) {
    // 实现代码
}
```

## 🔄 生成Swagger文档

### 1. 安装swag工具

```bash
# 安装最新版
go install github.com/swaggo/swag/cmd/swag@latest

# 验证安装
swag --version
```

### 2. 生成文档

在项目根目录执行：

```bash
# 生成文档
swag init

# 或指定main.go路径
swag init -g cmd/server/main.go

# 或指定输出目录
swag init -g cmd/server/main.go -o docs

# 解析依赖包
swag init --parseDependency --parseInternal
```

### 3. 文件说明

生成的文件位于 `docs/` 目录：

```
docs/
├── docs.go         # 主文档文件（自动生成）
├── swagger.json    # JSON格式文档
└── swagger.yaml    # YAML格式文档
```

### 4. 自动化脚本

创建 `scripts/swagger/generate.sh`：

```bash
#!/bin/bash
# Swagger文档生成脚本

echo "生成Swagger文档..."

# 进入项目根目录
cd "$(dirname "$0")/../.."

# 生成文档
swag init -g cmd/server/main.go -o docs --parseDependency --parseInternal

if [ $? -eq 0 ]; then
    echo "✓ Swagger文档生成成功"
    echo "文档地址: http://localhost:8080/swagger/index.html"
else
    echo "✗ Swagger文档生成失败"
    exit 1
fi
```

### 5. Git Hooks自动生成

在 `.git/hooks/pre-commit` 中添加：

```bash
#!/bin/bash
# Pre-commit hook: 自动生成Swagger文档

echo "检查Swagger注解..."

# 检查是否有API文件修改
if git diff --cached --name-only | grep -E "api/v1/.*\.go$"; then
    echo "检测到API文件修改，重新生成Swagger文档..."
    swag init -g cmd/server/main.go -o docs --parseDependency --parseInternal

    if [ $? -eq 0 ]; then
        git add docs/
        echo "✓ Swagger文档已更新"
    else
        echo "✗ Swagger文档生成失败"
        exit 1
    fi
fi
```

## 🎨 Swagger UI配置

### 1. 自定义Swagger UI

在 `cmd/server/main.go` 中添加配置：

```go
// @title           青羽写作平台 API
// @version         1.0
// @description     青羽写作平台后端服务API文档
// @termsOfService  http://qingyu.com/terms/

// @contact.name   API Support
// @contact.url    http://qingyu.com/support
// @contact.email  support@qingyu.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// X-Total-Count header for pagination
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Key
```

### 2. 环境变量配置

根据不同环境使用不同的host：

```go
// 开发环境
// @host localhost:8080

// 测试环境
// @host test-api.qingyu.com

// 生产环境
// @host api.qingyu.com
```

## 📋 注解检查清单

在添加新的API时，确保包含以下注解：

- [ ] `@Summary` - 简短描述
- [ ] `@Description` - 详细描述（可选）
- [ ] `@Tags` - 分组标签
- [ ] `@Accept` - 接受的请求类型
- [ ] `@Produce` - 返回的响应类型
- [ ] `@Param` - 所有参数说明
- [ ] `@Success` - 成功响应（至少一个）
- [ ] `@Failure` - 失败响应（至少包含400, 401, 500）
- [ ] `@Router` - 路由路径和方法
- [ ] `@Security` - 安全认证（如需要）

## 🔧 常见问题

### 1. 文档不更新

**问题**: 修改注解后文档不更新

**解决**:
```bash
# 重新生成文档
swag init -g cmd/server/main.go -o docs

# 重启服务
go run cmd/server/main.go

# 清除浏览器缓存
```

### 2. 模型不显示

**问题**: 响应模型在Swagger UI中不显示

**解决**: 确保模型导出（首字母大写）：
```go
type User struct {  // 正确
    ID string
}

type user struct {  // 错误 - 小写不会导出
    ID string
}
```

### 3. 参数验证规则不显示

**问题**: 参数的验证规则（如required, min, max）不显示

**解决**: 在模型中使用binding tag：
```go
type CreateUserRequest struct {
    Username string `json:"username" binding:"required" example:"john"`
    Email    string `json:"email" binding:"required,email" example:"john@example.com"`
    Age      int    `json:"age" binding:"min=18,max=100" example:"25"`
}
```

## 📖 最佳实践

### 1. 注解组织

```go
// 1. 基本信息
// @Summary 获取用户列表
// @Description 分页获取用户列表，支持搜索和过滤
// @Tags 用户

// 2. 请求/响应类型
// @Accept json
// @Produce json

// 3. 参数说明（按path -> query -> body顺序）
// @Param user_id path string true "用户ID"
// @Param page query int false "页码" default(1)

// 4. 响应说明
// @Success 200 {object} responses.UserListResponse
// @Failure 400 {object} responses.ErrorResponse

// 5. 路由和安全
// @Router /api/v1/users [get]
// @Security Bearer
```

### 2. 模型组织

创建 `models/` 和 `responses/` 目录：

```
models/
├── user.go          # 用户模型
├── book.go          # 书籍模型
└── chapter.go       # 章节模型

responses/
├── user_response.go # 用户响应
├── book_response.go # 书籍响应
└── error.go         # 错误响应
```

### 3. 示例值

为所有字段添加example tag：

```go
type User struct {
    ID        string    `json:"id" example:"123"`
    Username  string    `json:"username" example:"john_doe"`
    Email     string    `json:"email" example:"john@example.com"`
    CreatedAt time.Time `json:"created_at" example:"2026-01-06T10:00:00Z"`
}
```

### 4. 枚举值

使用Enums注解：

```go
// @Param status query string false "状态" Enums(published, draft, completed)
```

---

**文档维护**: 青羽后端架构团队
**最后更新**: 2026-01-06
