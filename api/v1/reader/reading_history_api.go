package reader

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/pkg/response"
	"Qingyu_backend/service/interfaces"
)

// ReadingHistoryAPI 阅读历史API处理器
type ReadingHistoryAPI struct {
	historyService interfaces.ReadingHistoryService
}

// NewReadingHistoryAPI 创建API处理器
func NewReadingHistoryAPI(historyService interfaces.ReadingHistoryService) *ReadingHistoryAPI {
	return &ReadingHistoryAPI{
		historyService: historyService,
	}
}

// RecordReading 记录阅读
// @Summary 记录阅读历史
// @Tags 阅读历史
// @Accept json
// @Produce json
// @Param body body RecordReadingRequest true "阅读记录参数"
// @Success 200 {object} pkg.Response
// @Failure 400 {object} pkg.Response
// @Router /api/v1/reader/reading-history [post]
func (api *ReadingHistoryAPI) RecordReading(c *gin.Context) {
	var req RecordReadingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未认证")
		return
	}

	// 记录阅读
	err := api.historyService.RecordReading(
		c.Request.Context(),
		userID.(string),
		req.BookID,
		req.ChapterID,
		req.StartTime,
		req.EndTime,
		req.Progress,
		req.DeviceType,
		req.DeviceID,
	)
	if err != nil {
		c.Error(err)
		return
	}

	response.Created(c, gin.H{})
}

// GetReadingHistories 获取阅读历史列表
// @Summary 获取阅读历史
// @Tags 阅读历史
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param book_id query string false "书籍ID"
// @Success 200 {object} pkg.Response
// @Failure 400 {object} pkg.Response
// @Router /api/v1/reader/reading-history [get]
func (api *ReadingHistoryAPI) GetReadingHistories(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未认证")
		return
	}

	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	bookID := c.Query("book_id")

	var histories interface{}
	var pagination interface{}
	var err error

	// 根据是否有book_id参数决定调用哪个方法
	if bookID != "" {
		histories, pagination, err = api.historyService.GetUserHistoriesByBook(
			c.Request.Context(),
			userID.(string),
			bookID,
			page,
			pageSize,
		)
	} else {
		histories, pagination, err = api.historyService.GetUserHistories(
			c.Request.Context(),
			userID.(string),
			page,
			pageSize,
		)
	}

	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, gin.H{
		"histories":  histories,
		"pagination": pagination,
	})
}

// GetReadingStats 获取阅读统计
// @Summary 获取阅读统计
// @Tags 阅读历史
// @Accept json
// @Produce json
// @Param days query int false "统计天数" default(30)
// @Success 200 {object} pkg.Response
// @Failure 400 {object} pkg.Response
// @Router /api/v1/reader/reading-history/stats [get]
func (api *ReadingHistoryAPI) GetReadingStats(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未认证")
		return
	}

	// 解析天数参数
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))

	// 获取总体统计
	stats, err := api.historyService.GetUserReadingStats(
		c.Request.Context(),
		userID.(string),
	)
	if err != nil {
		c.Error(err)
		return
	}

	// 获取每日统计
	dailyStats, err := api.historyService.GetUserDailyReadingStats(
		c.Request.Context(),
		userID.(string),
		days,
	)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, gin.H{
		"summary":     stats,
		"daily_stats": dailyStats,
	})
}

// DeleteHistory 删除单条历史记录
// @Summary 删除历史记录
// @Tags 阅读历史
// @Accept json
// @Produce json
// @Param id path string true "历史记录ID"
// @Success 200 {object} pkg.Response
// @Failure 400 {object} pkg.Response
// @Router /api/v1/reader/reading-history/{id} [delete]
func (api *ReadingHistoryAPI) DeleteHistory(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未认证")
		return
	}

	historyID := c.Param("id")
	if historyID == "" {
		response.BadRequest(c, "参数错误", "历史记录ID不能为空")
		return
	}

	// 删除历史记录
	err := api.historyService.DeleteHistory(
		c.Request.Context(),
		userID.(string),
		historyID,
	)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, gin.H{})
}

// ClearHistories 清空历史记录
// @Summary 清空历史记录
// @Tags 阅读历史
// @Accept json
// @Produce json
// @Success 200 {object} pkg.Response
// @Failure 400 {object} pkg.Response
// @Router /api/v1/reader/reading-history [delete]
func (api *ReadingHistoryAPI) ClearHistories(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未认证")
		return
	}

	// 清空历史记录
	err := api.historyService.ClearUserHistories(
		c.Request.Context(),
		userID.(string),
	)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, gin.H{})
}

// =======================
// 请求/响应结构
// =======================

// RecordReadingRequest 记录阅读请求
type RecordReadingRequest struct {
	BookID     string    `json:"book_id" binding:"required"`
	ChapterID  string    `json:"chapter_id" binding:"required"`
	StartTime  time.Time `json:"start_time" binding:"required"`
	EndTime    time.Time `json:"end_time" binding:"required"`
	Progress   float64   `json:"progress" binding:"min=0,max=100"`
	DeviceType string    `json:"device_type"`
	DeviceID   string    `json:"device_id"`
}

// RecordReadingResponse 记录阅读响应
type RecordReadingResponse struct {
	Message string `json:"message"`
}

// HistoryListResponse 历史列表响应
type HistoryListResponse struct {
	Histories  interface{} `json:"histories"`
	Pagination interface{} `json:"pagination"`
}

// ReadingStatsResponse 阅读统计响应
type ReadingStatsResponse struct {
	Summary    interface{} `json:"summary"`
	DailyStats interface{} `json:"daily_stats"`
}
