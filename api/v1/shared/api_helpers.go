package shared

import (
	"bytes"
	"context"
	"io"
	"strconv"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/pkg/response"
)

// ============ 参数获取辅助函数 ============

// GetUserID 从上下文获取用户ID，如果不存在或类型错误则返回未授权响应
// 返回: (userID, ok) - ok为false表示已发送错误响应
func GetUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "请先登录")
		return "", false
	}

	uid, ok := userID.(string)
	if !ok {
		response.Unauthorized(c, "用户信息格式错误")
		return "", false
	}

	return uid, true
}

// GetUserIDOptional 从上下文获取用户ID（可选），不存在时返回空字符串
// 返回: userID字符串，如果不存在或类型错误返回空字符串
func GetUserIDOptional(c *gin.Context) string {
	userID, exists := c.Get("user_id")
	if !exists {
		return ""
	}

	if uid, ok := userID.(string); ok {
		return uid
	}
	return ""
}

// ============ 路径参数辅助函数 ============

// GetRequiredParam 获取必需的路径参数，如果为空则返回错误响应
// 返回: (param, ok) - ok为false表示已发送错误响应
func GetRequiredParam(c *gin.Context, key, displayName string) (string, bool) {
	param := c.Param(key)
	if param == "" {
		response.BadRequest(c, "参数错误", displayName+"不能为空")
		return "", false
	}
	return param, true
}

// GetRequiredQuery 获取必需的查询参数，如果为空则返回错误响应
// 返回: (param, ok) - ok为false表示已发送错误响应
func GetRequiredQuery(c *gin.Context, key, displayName string) (string, bool) {
	param := c.Query(key)
	if param == "" {
		response.BadRequest(c, "参数错误", displayName+"不能为空")
		return "", false
	}
	return param, true
}

// ============ 分页参数辅助函数 ============

// PaginationParams 分页参数
type PaginationParams struct {
	Page     int
	PageSize int
	Limit    int
	Offset   int
}

// GetPaginationParams 获取分页参数并验证
// defaultPage: 默认页码
// defaultPageSize: 默认每页数量
// maxPageSize: 最大每页数量（0表示不限制）
func GetPaginationParams(c *gin.Context, defaultPage, defaultPageSize, maxPageSize int) PaginationParams {
	page, _ := strconv.Atoi(c.DefaultQuery("page", strconv.Itoa(defaultPage)))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("size", strconv.Itoa(defaultPageSize)))

	// 验证页码
	if page < 1 {
		page = defaultPage
	}

	// 验证每页数量
	if pageSize < 1 {
		pageSize = defaultPageSize
	}
	if maxPageSize > 0 && pageSize > maxPageSize {
		pageSize = maxPageSize
	}

	return PaginationParams{
		Page:     page,
		PageSize: pageSize,
		Limit:    pageSize,
		Offset:   (page - 1) * pageSize,
	}
}

// GetPaginationParamsStandard 获取标准分页参数（默认1页，每页20条，最大100条）
func GetPaginationParamsStandard(c *gin.Context) PaginationParams {
	return GetPaginationParams(c, 1, 20, 100)
}

// GetPaginationParamsLarge 获取大容量分页参数（默认1页，每页50条，最大200条）
func GetPaginationParamsLarge(c *gin.Context) PaginationParams {
	return GetPaginationParams(c, 1, 50, 200)
}

// GetPaginationParamsSmall 获取小容量分页参数（默认1页，每页10条，最大50条）
func GetPaginationParamsSmall(c *gin.Context) PaginationParams {
	return GetPaginationParams(c, 1, 10, 50)
}

// ============ 数字参数辅助函数 ============

//GetIntParam 获取整数参数（路径或查询），支持默认值和范围验证
// key: 参数名
// isQuery: true表示查询参数，false表示路径参数
// defaultValue: 默认值
// min, max: 最小值和最大值（0表示不限制）
func GetIntParam(c *gin.Context, key string, isQuery bool, defaultValue, min, max int) int {
	var valueStr string
	if isQuery {
		valueStr = c.Query(key)
	} else {
		valueStr = c.Param(key)
	}

	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}

	if min > 0 && value < min {
		return min
	}
	if max > 0 && value > max {
		return max
	}

	return value
}

// ============ JSON绑定辅助函数 ============

// BindAndValidate 绑定并验证JSON请求体
// 返回: (ok) - false表示已发送错误响应
// 注意：ShouldBindJSON会自动验证结构体的binding标签
func BindAndValidate(c *gin.Context, req interface{}) bool {
	// 读取原始请求体用于错误响应
	bodyBytes, _ := c.GetRawData()

	// 重新设置请求体供后续使用
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// 绑定并验证（ShouldBindJSON会自动验证binding标签）
	if err := c.ShouldBindJSON(req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return false
	}
	return true
}

// BindJSON 仅绑定JSON请求体（不验证）
// 返回: (ok) - false表示已发送错误响应
func BindJSON(c *gin.Context, req interface{}) bool {
	if err := c.ShouldBindJSON(req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return false
	}
	return true
}

// BindParams 绑定URI路径参数和Query查询参数（自动验证）
// 使用结构体标签进行参数绑定和验证
// 返回: (ok) - false表示已发送错误响应
//
// 使用示例：
//	var params struct {
//	    BookID string `uri:"bookId" binding:"required"`
//	    Page   int    `form:"page" binding:"min=1"`
//	}
//	if !BindParams(c, &params) { return }
func BindParams(c *gin.Context, req interface{}) bool {
	// 先绑定URI参数
	if err := c.ShouldBindUri(req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return false
	}
	// 再绑定Query参数
	if err := c.ShouldBindQuery(req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return false
	}
	return true
}

// ============ 分页响应辅助函数 ============

// RespondWithPaginated 响应分页数据
func RespondWithPaginated(c *gin.Context, data interface{}, total, page, size int, message string) {
	response.Success(c, gin.H{
		"list":  data,
		"total": total,
		"page":  page,
		"size":  size,
	})
}

// ============ 上下文辅助函数 ============

// AddUserIDToContext 将用户ID添加到context.Context（用于传递给service层）
func AddUserIDToContext(c *gin.Context) context.Context {
	ctx := c.Request.Context()
	if userID, exists := c.Get("user_id"); exists {
		ctx = context.WithValue(ctx, "userId", userID.(string))
	}
	return ctx
}

// ContextWithUserID 创建带有用户ID的context（wrapper版本，返回gin.Context）
func ContextWithUserID(c *gin.Context) bool {
	userID, ok := GetUserID(c)
	if !ok {
		return false
	}

	ctx := context.WithValue(c.Request.Context(), "userId", userID)
	c.Request = c.Request.WithContext(ctx)

	return true
}

// ============ 批量操作辅助函数 ============

// ValidateBatchIDs 验证批量操作的ID列表
func ValidateBatchIDs(c *gin.Context, ids []string, displayName string) bool {
	if len(ids) == 0 {
		response.BadRequest(c, "参数错误", displayName+"不能为空")
		return false
	}

	if len(ids) > 1000 {
		response.BadRequest(c, "参数错误", displayName+"数量不能超过1000个")
		return false
	}

	return true
}
