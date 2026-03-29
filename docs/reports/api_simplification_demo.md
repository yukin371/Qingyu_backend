# API层简化演示：使用现有错误系统

## 当前状态 vs 简化后

### 1. 错误处理中间件（项目已有）

项目 `pkg/errors/middleware_funcs.go` 已经实现了：
- `ErrorMiddleware()` - 统一错误处理中间件
- `UnifiedError` - 完整的错误结构
- `ErrorFactory` - 错误工厂模式

### 2. 代码对比演示

#### 当前实现（冗长）

```go
func (api *ChapterAPI) GetChapterByNumber(c *gin.Context) {
	// 参数绑定
	bookID, ok := shared.GetRequiredParam(c, "bookId", "书籍ID")
	if !ok {
		return
	}
	chapterNum := shared.GetIntParam(c, "chapterNum", false, 0, 1, 0)
	if chapterNum == 0 {
		response.BadRequest(c, "参数错误", "无效的章节号")
		return
	}
	userID := shared.GetUserIDOptional(c)

	// 调用Service层
	content, err := api.chapterService.GetChapterByNumber(c.Request.Context(), userID, bookID, chapterNum)
	if err != nil {
		// 冗长的错误处理
		if err == readerservice.ErrChapterNotFound {
			response.NotFound(c, "章节不存在")
			return
		}
		if err == readerservice.ErrAccessDenied && content != nil {
			response.Forbidden(c, "无权访问")
			return
		}
		response.InternalError(c, err)
		return
	}

	// 响应封装
	response.Success(c, content)
}
```

#### 简化后（使用现有错误系统）

```go
// 1. 首先确保路由注册了错误中间件
// router.Use(errors.ErrorMiddleware("reader"))

func (api *ChapterAPI) GetChapterByNumber(c *gin.Context) {
	// 参数绑定
	bookID, ok := shared.GetRequiredParam(c, "bookId", "书籍ID")
	if !ok {
		return
	}
	chapterNum := shared.GetIntParam(c, "chapterNum", false, 1, 1, 0)
	userID := shared.GetUserIDOptional(c)

	// 调用Service层
	content, err := api.chapterService.GetChapterByNumber(c.Request.Context(), userID, bookID, chapterNum)

	// 使用c.Error()将错误交给中间件处理
	if err != nil {
		c.Error(err)
		return
	}

	// 响应封装
	c.JSON(http.StatusOK, content)
}
```

### 3. 更简洁的版本（创建辅助函数）

使用辅助函数 `api/v1/shared/helpers.go`:

```go
package shared

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

// BindParams 绑定路径参数到结构体
func BindParams(c *gin.Context, params interface{}) bool {
	if err := c.ShouldBindUri(params); err != nil {
		c.Error(err)
		return false
	}
	return true
}

// BindJSON 绑定JSON请求体
func BindJSON(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		c.Error(err)
		return false
	}
	return true
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    data,
	})
}

// Created 创建成功响应
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, gin.H{
		"code":    201,
		"message": "created",
		"data":    data,
	})
}
```

### 4. 最简洁的API实现

```go
package reader

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/helpers"
	"Qingyu_backend/service/interfaces"
)

type ChapterAPI struct {
	chapterService interfaces.ReaderChapterService
}

// GetChapterByNumberParams 请求参数结构
type GetChapterByNumberParams struct {
	BookID     string `uri:"bookId" binding:"required"`
	ChapterNum int    `uri:"chapterNum" binding:"required,min=1"`
}

func (api *ChapterAPI) GetChapterByNumber(c *gin.Context) {
	// 1. 参数绑定（自动验证）
	var params GetChapterByNumberParams
	if !helpers.BindParams(c, &params) {
		return
	}

	// 2. 获取用户ID
	userID := helpers.GetUserIDOptional(c)

	// 3. 调用Service层（错误交给中间件）
	content, err := api.chapterService.GetChapterByNumber(c.Request.Context(), userID, params.BookID, params.ChapterNum)
	if err != nil {
		c.Error(err)
		return
	}

	// 4. 成功响应
	helpers.Success(c, content)
}
```

### 5. 代码量对比

| 版本 | 行数 | 说明 |
|------|------|------|
| 当前实现 | ~40行 | 大量重复的错误处理 |
| 使用c.Error | ~20行 | 减少约50% |
| 使用辅助函数 | ~15行 | 减少约60% |

### 6. 实施步骤

#### Step 1: 确保路由注册了错误中间件
```go
// 在 router.go 中
import "Qingyu_backend/pkg/errors"

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// 注册错误中间件
	router.Use(errors.ErrorMiddleware("reader"))

	return router
}
```

#### Step 2: Service层返回UnifiedError
```go
// Service层示例
func (s *ReaderChapterService) GetChapterByNumber(ctx context.Context, userID, bookID string, chapterNum int) (*ChapterDTO, error) {
	// 业务逻辑检查
	if chapter == nil {
		return nil, errors.NewNotFound("Chapter")
	}
	if !canAccess {
		return nil, errors.NewForbidden("无权访问该章节")
	}

	return chapter, nil
}
```

#### Step 3: 创建辅助函数（可选）
将常用的绑定和响应操作封装为辅助函数

#### Step 4: 重构API层
逐个重构API函数，使用c.Error()替代手动错误处理

## 结论

项目**已有完善的错误处理系统**，只需要：
1. ✅ 确保中间件已注册
2. ✅ Service层返回`UnifiedError`
3. ✅ API层使用`c.Error(err)`

代码即可减少40-60%，而且更加一致和可维护喵~
