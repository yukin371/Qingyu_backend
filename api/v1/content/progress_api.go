package content

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/models/dto"
	"Qingyu_backend/api/v1/shared"
	contentService "Qingyu_backend/service/interfaces/content"
)

// ProgressAPI 阅读进度API
type ProgressAPI struct {
	progressService contentService.ReadingProgressServicePort
}

// NewProgressAPI 创建进度API实例
func NewProgressAPI(progressService contentService.ReadingProgressServicePort) *ProgressAPI {
	return &ProgressAPI{
		progressService: progressService,
	}
}

// GetProgress 获取阅读进度
//
//	@Summary		获取阅读进度
//	@Description	获取用户对指定书籍的阅读进度
//	@Tags			内容管理-进度
//	@Accept			json
//	@Produce		json
//	@Param			bookId	path		string	true	"书籍ID"
//	@Success		200		{object}	shared.APIResponse
//	@Failure		400		{object}	shared.APIResponse
//	@Failure		401		{object}	shared.APIResponse
//	@Failure		500		{object}	shared.APIResponse
//	@Router			/api/v1/content/progress/{bookId} [get]
func (api *ProgressAPI) GetProgress(c *gin.Context) {
	bookID := c.Param("bookId")
	if bookID == "" {
		shared.BadRequest(c, "参数错误", "书籍ID不能为空")
		return
	}

	// 获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		shared.Unauthorized(c, "请先登录")
		return
	}

	result, err := api.progressService.GetProgress(c.Request.Context(), userID.(string), bookID)
	if err != nil {
		c.Error(err)
		return
	}

	shared.Success(c, 200, "获取成功", result)
}

// SaveProgress 保存阅读进度
//
//	@Summary		保存阅读进度
//	@Description	保存用户的阅读位置和进度百分比
//	@Tags			内容管理-进度
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.SaveProgressRequest	true	"保存进度请求"
//	@Success		200		{object}	shared.APIResponse
//	@Failure		400		{object}	shared.APIResponse
//	@Failure		401		{object}	shared.APIResponse
//	@Failure		500		{object}	shared.APIResponse
//	@Router			/api/v1/content/progress [post]
func (api *ProgressAPI) SaveProgress(c *gin.Context) {
	var req dto.SaveProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.BadRequest(c, "参数错误", err.Error())
		return
	}

	// 获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		shared.Unauthorized(c, "请先登录")
		return
	}

	req.UserID = userID.(string)

	err := api.progressService.SaveProgress(c.Request.Context(), &req)
	if err != nil {
		c.Error(err)
		return
	}

	shared.Success(c, 200, "保存成功", nil)
}

// UpdateReadingTime 更新阅读时长
//
//	@Summary		更新阅读时长
//	@Description	累加用户的阅读时长
//	@Tags			内容管理-进度
//	@Accept			json
//	@Produce		json
//	@Param			request	body		object{bookId=string,duration=int}	true	"更新时长请求"
//	@Success		200		{object}	shared.APIResponse
//	@Failure		400		{object}	shared.APIResponse
//	@Failure		401		{object}	shared.APIResponse
//	@Failure		500		{object}	shared.APIResponse
//	@Router			/api/v1/content/progress/reading-time [put]
func (api *ProgressAPI) UpdateReadingTime(c *gin.Context) {
	var req struct {
		BookID   string `json:"bookId" binding:"required"`
		Duration int    `json:"duration" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.BadRequest(c, "参数错误", err.Error())
		return
	}

	// 获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		shared.Unauthorized(c, "请先登录")
		return
	}

	err := api.progressService.UpdateReadingTime(c.Request.Context(), userID.(string), req.BookID, req.Duration)
	if err != nil {
		c.Error(err)
		return
	}

	shared.Success(c, 200, "更新成功", nil)
}

// GetRecentBooks 获取最近阅读的书籍
//
//	@Summary		获取最近阅读的书籍
//	@Description	返回用户最近阅读的书籍列表
//	@Tags			内容管理-进度
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		int	false	"数量限制"	default(20)
//	@Success		200		{object}	shared.APIResponse
//	@Failure		401	{object}	shared.APIResponse
//	@Failure		500	{object}	shared.APIResponse
//	@Router			/api/v1/content/progress/recent [get]
func (api *ProgressAPI) GetRecentBooks(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		shared.Unauthorized(c, "请先登录")
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if limit < 1 || limit > 100 {
		limit = 20
	}

	result, err := api.progressService.GetRecentBooks(c.Request.Context(), userID.(string), limit)
	if err != nil {
		c.Error(err)
		return
	}

	shared.Success(c, 200, "获取成功", result)
}

// GetReadingStats 获取阅读统计
//
//	@Summary		获取阅读统计
//	@Description	返回用户的阅读统计数据
//	@Tags			内容管理-进度
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	shared.APIResponse
//	@Failure		401	{object}	shared.APIResponse
//	@Failure		500	{object}	shared.APIResponse
//	@Router			/api/v1/content/progress/stats [get]
func (api *ProgressAPI) GetReadingStats(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		shared.Unauthorized(c, "请先登录")
		return
	}

	result, err := api.progressService.GetReadingStats(c.Request.Context(), userID.(string))
	if err != nil {
		c.Error(err)
		return
	}

	shared.Success(c, 200, "获取成功", result)
}

// GetReadingHistory 获取阅读历史
//
//	@Summary		获取阅读历史
//	@Description	返回用户的阅读历史记录，支持分页
//	@Tags			内容管理-进度
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int	false	"页码"	default(1)
//	@Param			pageSize	query		int		false	"每页数量"	default(20)
//	@Success		200			{object}	shared.APIResponse
//	@Failure		400			{object}	shared.APIResponse
//	@Failure		401			{object}	shared.APIResponse
//	@Failure		500			{object}	shared.APIResponse
//	@Router			/api/v1/content/progress/history [get]
func (api *ProgressAPI) GetReadingHistory(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		shared.Unauthorized(c, "请先登录")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	result, err := api.progressService.GetReadingHistory(c.Request.Context(), userID.(string), page, pageSize)
	if err != nil {
		c.Error(err)
		return
	}

	shared.Paginated(c, result.Progresses, result.Total, page, pageSize, "获取成功")
}

// GetUnfinishedBooks 获取未读完的书籍
//
//	@Summary		获取未读完的书籍
//	@Description	返回所有未读完的书籍列表
//	@Tags			内容管理-进度
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	shared.APIResponse
//	@Failure		401	{object}	shared.APIResponse
//	@Failure		500	{object}	shared.APIResponse
//	@Router			/api/v1/content/progress/unfinished [get]
func (api *ProgressAPI) GetUnfinishedBooks(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		shared.Unauthorized(c, "请先登录")
		return
	}

	result, err := api.progressService.GetUnfinishedBooks(c.Request.Context(), userID.(string))
	if err != nil {
		c.Error(err)
		return
	}

	shared.Success(c, 200, "获取成功", result)
}

// GetFinishedBooks 获取已读完的书籍
//
//	@Summary		获取已读完的书籍
//	@Description	返回所有已读完的书籍列表
//	@Tags			内容管理-进度
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	shared.APIResponse
//	@Failure		401	{object}	shared.APIResponse
//	@Failure		500	{object}	shared.APIResponse
//	@Router			/api/v1/content/progress/finished [get]
func (api *ProgressAPI) GetFinishedBooks(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		shared.Unauthorized(c, "请先登录")
		return
	}

	result, err := api.progressService.GetFinishedBooks(c.Request.Context(), userID.(string))
	if err != nil {
		c.Error(err)
		return
	}

	shared.Success(c, 200, "获取成功", result)
}
