package reader

import (
	readerModels "Qingyu_backend/models/reader"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/shared"
	"Qingyu_backend/service/interfaces"
)

const (
	hoursPerDay = 24
	daysPerWeek = 7
)

// ProgressAPI 阅读进度API
type ProgressAPI struct {
	readerService interfaces.ReaderService
}

// NewProgressAPI 创建阅读进度API实例
func NewProgressAPI(readerService interfaces.ReaderService) *ProgressAPI {
	return &ProgressAPI{
		readerService: readerService,
	}
}

// SaveProgressRequest 保存进度请求
type SaveProgressRequest struct {
	BookID    string  `json:"bookId" binding:"required"`
	ChapterID string  `json:"chapterId" binding:"required"`
	Progress  float64 `json:"progress" binding:"required,min=0,max=1"`
}

// UpdateReadingTimeRequest 更新阅读时长请求
type UpdateReadingTimeRequest struct {
	BookID   string `json:"bookId" binding:"required"`
	Duration int64  `json:"duration" binding:"required,min=1"`
}

// GetReadingProgress 获取阅读进度
//
//	@Summary	获取阅读进度
//	@Tags		阅读器
//	@Param		bookId	path		string	true	"书籍ID"
//	@Success	200		{object}	shared.APIResponse
//	@Router		/api/v1/reader/progress/{bookId} [get]
func (api *ProgressAPI) GetReadingProgress(c *gin.Context) {
	bookID := c.Param("bookId")

	// 获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	progress, err := api.readerService.GetReadingProgress(c.Request.Context(), userID.(string), bookID)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "获取阅读进度失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取成功", progress)
}

// SaveReadingProgress 保存阅读进度
//
//	@Summary	保存阅读进度
//	@Tags		阅读器
//	@Param		request	body		SaveProgressRequest	true	"保存进度请求"
//	@Success	200		{object}	shared.APIResponse
//	@Router		/api/v1/reader/progress [post]
func (api *ProgressAPI) SaveReadingProgress(c *gin.Context) {
	var req SaveProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.ValidationError(c, err)
		return
	}

	// 获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	err := api.readerService.SaveReadingProgress(c.Request.Context(), userID.(string), req.BookID, req.ChapterID, req.Progress)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "保存阅读进度失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "保存成功", nil)
}

// UpdateReadingTime 更新阅读时长
//
//	@Summary	更新阅读时长
//	@Tags		阅读器
//	@Param		request	body		UpdateReadingTimeRequest	true	"更新时长请求"
//	@Success	200		{object}	shared.APIResponse
//	@Router		/api/v1/reader/progress/reading-time [put]
func (api *ProgressAPI) UpdateReadingTime(c *gin.Context) {
	var req UpdateReadingTimeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.ValidationError(c, err)
		return
	}

	// 获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	err := api.readerService.UpdateReadingTime(c.Request.Context(), userID.(string), req.BookID, req.Duration)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "更新阅读时长失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "更新成功", nil)
}

// GetRecentReading 获取最近阅读记录
//
//	@Summary	获取最近阅读记录
//	@Tags		阅读器
//	@Param		limit	query		int	false	"数量限制"	default(20)
//	@Success	200		{object}	shared.APIResponse
//	@Router		/api/v1/reader/progress/recent [get]
func (api *ProgressAPI) GetRecentReading(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	progresses, err := api.readerService.GetRecentReading(c.Request.Context(), userID.(string), limit)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "获取最近阅读记录失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取成功", progresses)
}

// GetReadingHistory 获取阅读历史
//
//	@Summary	获取阅读历史
//	@Tags		阅读器
//	@Param		page	query		int	false	"页码"	default(1)
//	@Param		size	query		int	false	"每页数量"	default(20)
//	@Success	200		{object}	shared.APIResponse
//	@Router		/api/v1/reader/progress/history [get]
func (api *ProgressAPI) GetReadingHistory(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	progresses, total, err := api.readerService.GetReadingHistory(c.Request.Context(), userID.(string), page, size)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "获取阅读历史失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取成功", gin.H{
		"progresses": progresses,
		"total":      total,
		"page":       page,
		"size":       size,
	})
}

// GetReadingStats 获取阅读统计
//
//	@Summary	获取阅读统计
//	@Tags		阅读器
//	@Param		period	query		string	false	"统计周期"	default("all")
//	@Success	200		{object}	shared.APIResponse
//	@Router		/api/v1/reader/progress/stats [get]
func (api *ProgressAPI) GetReadingStats(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	period := c.DefaultQuery("period", "all")

	var totalTime int64
	var err error

	switch period {
	case "today":
		// 今天
		start := time.Now().Truncate(hoursPerDay * time.Hour)
		end := start.Add(hoursPerDay * time.Hour)
		totalTime, err = api.readerService.GetReadingTimeByPeriod(c.Request.Context(), userID.(string), start, end)
	case "week":
		// 本周
		now := time.Now()
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = daysPerWeek
		}
		start := now.AddDate(0, 0, -(weekday - 1)).Truncate(hoursPerDay * time.Hour)
		end := start.AddDate(0, 0, daysPerWeek)
		totalTime, err = api.readerService.GetReadingTimeByPeriod(c.Request.Context(), userID.(string), start, end)
	case "month":
		// 本月
		now := time.Now()
		start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		end := start.AddDate(0, 1, 0)
		totalTime, err = api.readerService.GetReadingTimeByPeriod(c.Request.Context(), userID.(string), start, end)
	default:
		// 总计
		totalTime, err = api.readerService.GetTotalReadingTime(c.Request.Context(), userID.(string))
	}

	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "获取阅读统计失败", err.Error())
		return
	}

	// 获取未读完和已读完的书籍
	unfinished, errUnfinished := api.readerService.GetUnfinishedBooks(c.Request.Context(), userID.(string))
	if errUnfinished != nil {
		unfinished = []*readerModels.ReadingProgress{} // 返回空列表而非失败
	}

	finished, errFinished := api.readerService.GetFinishedBooks(c.Request.Context(), userID.(string))
	if errFinished != nil {
		finished = []*readerModels.ReadingProgress{} // 返回空列表而非失败
	}

	shared.Success(c, http.StatusOK, "获取成功", gin.H{
		"totalReadingTime": totalTime,
		"unfinishedCount":  len(unfinished),
		"finishedCount":    len(finished),
		"period":           period,
	})
}

// GetUnfinishedBooks 获取未读完的书籍
//
//	@Summary	获取未读完的书籍
//	@Tags		阅读器
//	@Success	200	{object}	shared.APIResponse
//	@Router		/api/v1/reader/progress/unfinished [get]
func (api *ProgressAPI) GetUnfinishedBooks(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	progresses, err := api.readerService.GetUnfinishedBooks(c.Request.Context(), userID.(string))
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "获取未读完书籍失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取成功", progresses)
}

// GetFinishedBooks 获取已读完的书籍
//
//	@Summary	获取已读完的书籍
//	@Tags		阅读器
//	@Success	200	{object}	shared.APIResponse
//	@Router		/api/v1/reader/progress/finished [get]
func (api *ProgressAPI) GetFinishedBooks(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	progresses, err := api.readerService.GetFinishedBooks(c.Request.Context(), userID.(string))
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "获取已读完书籍失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取成功", progresses)
}
