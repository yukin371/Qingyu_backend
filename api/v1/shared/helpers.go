package shared

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/pkg/errors"
)

// === 辅助函数：简化API层代码 ===

// BindParams 绑定路径参数和查询参数（自动验证）
// 使用示例：
//
//	var params struct {
//	    BookID string `uri:"bookId" binding:"required"`
//	    Page   int    `form:"page" binding:"min=1"`
//	}
//	if !BindParams(c, &params) { return }
func BindParams(c *gin.Context, params interface{}) bool {
	// 先绑定URI参数
	if err := c.ShouldBindUri(params); err != nil {
		c.Error(errors.UserServiceFactory.ValidationError("INVALID_PARAMS", "参数错误", err.Error()))
		return false
	}
	// 再绑定Query参数
	if err := c.ShouldBindQuery(params); err != nil {
		c.Error(errors.UserServiceFactory.ValidationError("INVALID_PARAMS", "参数错误", err.Error()))
		return false
	}
	return true
}

// BindJSON 绑定JSON请求体（自动验证）
// 使用示例：
//
//	var req CreateChapterRequest
//	if !BindJSON(c, &req) { return }
func BindJSON(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		c.Error(errors.UserServiceFactory.ValidationError("INVALID_JSON", "请求体格式错误", err.Error()))
		return false
	}
	return true
}

// Success 成功响应（200）
// 使用示例：
//
//	Success(c, userData)
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    data,
	})
}

// SuccessWithMessage 成功响应（带消息）
// 使用示例：
//
//	SuccessWithMessage(c, "创建成功", newData)
func SuccessWithMessage(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": msg,
		"data":    data,
	})
}

// Created 创建成功响应（201）
// 使用示例：
//
//	Created(c, createdResource)
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, gin.H{
		"code":    201,
		"message": "created",
		"data":    data,
	})
}

// NoContent 无内容响应（204）
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// === 简化前后的对比 ===

// 简化前：
/*
func (api *ChapterAPI) GetChapter(c *gin.Context) {
	bookID, ok := shared.GetRequiredParam(c, "bookId", "书籍ID")
	if !ok {
		return
	}
	chapterID, ok := shared.GetRequiredParam(c, "chapterId", "章节ID")
	if !ok {
		return
	}
	userID := shared.GetUserIDOptional(c)

	content, err := api.chapterService.GetChapterContent(c.Request.Context(), userID, bookID, chapterID)
	if err != nil {
		if err == readerservice.ErrChapterNotFound {
			response.NotFound(c, "章节不存在")
			return
		}
		if err == readerservice.ErrAccessDenied {
			response.Forbidden(c, "无权访问")
			return
		}
		response.InternalError(c, err)
		return
	}

	response.Success(c, content)
}
*/

// 简化后：
/*
func (api *ChapterAPI) GetChapter(c *gin.Context) {
	// 1. 参数绑定（结构体+验证）
	var params struct {
		BookID    string `uri:"bookId" binding:"required"`
		ChapterID string `uri:"chapterId" binding:"required"`
	}
	if !BindParams(c, &params) {
		return
	}

	// 2. 获取用户ID
	userID := GetUserIDOptional(c)

	// 3. 调用Service层（错误交给中间件）
	content, err := api.chapterService.GetChapterContent(c.Request.Context(), userID, params.BookID, params.ChapterID)
	if err != nil {
		c.Error(err)
		return
	}

	// 4. 成功响应
	Success(c, content)
}
*/

// === 行数对比 ===
// 简化前: ~35行
// 简化后: ~18行（减少约50%）
