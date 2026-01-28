package writer

import (
	"Qingyu_backend/api/v1/shared"
	"Qingyu_backend/models/stats"
	readingStats "Qingyu_backend/service/reader/stats"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"Qingyu_backend/pkg/response"
)

// StatsApi 阅读/书店统计API
// 职责：处理作品、章节的阅读统计
type StatsApi struct {
	statsService *readingStats.ReadingStatsService
}

// NewStatsApi 创建统计API
func NewStatsApi(statsService *readingStats.ReadingStatsService) *StatsApi {
	return &StatsApi{
		statsService: statsService,
	}
}

// GetBookStats 获取作品统计
// @Summary 获取作品统计数据
// @Description 获取作品的完整统计信息，包括阅读、收入、互动等数据
// @Tags Stats
// @Accept json
// @Produce json
// @Param book_id path string true "作品ID"
// @Success 200 {object} shared.APIResponse
// @Failure 400 {object} shared.APIResponse
// @Failure 500 {object} shared.APIResponse
// @Router /api/v1/writer/books/{book_id}/stats [get]
func (api *StatsApi) GetBookStats(c *gin.Context) {
	bookID := c.Param("book_id")
	if bookID == "" {
		response.BadRequest(c,  "参数错误", "作品ID不能为空")
		return
	}

	// 获取作品统计
	bookStats, err := api.statsService.CalculateBookStats(c.Request.Context(), bookID)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	if bookStats == nil {
		shared.Error(c, http.StatusNotFound, "作品统计不存在", "")
		return
	}

	shared.Success(c, http.StatusOK, "获取成功", bookStats)
}

// GetChapterStats 获取章节统计
// @Summary 获取章节统计数据
// @Description 获取单个章节的统计信息
// @Tags Stats
// @Accept json
// @Produce json
// @Param chapter_id path string true "章节ID"
// @Success 200 {object} shared.APIResponse
// @Failure 400 {object} shared.APIResponse
// @Failure 500 {object} shared.APIResponse
// @Router /api/v1/writer/chapters/{chapter_id}/stats [get]
func (api *StatsApi) GetChapterStats(c *gin.Context) {
	chapterID := c.Param("chapter_id")
	if chapterID == "" {
		response.BadRequest(c,  "参数错误", "章节ID不能为空")
		return
	}

	// 获取章节统计
	chapterStats, err := api.statsService.CalculateChapterStats(c.Request.Context(), chapterID)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	if chapterStats == nil {
		shared.Error(c, http.StatusNotFound, "章节统计不存在", "")
		return
	}

	shared.Success(c, http.StatusOK, "获取成功", chapterStats)
}

// GetBookHeatmap 获取阅读热力图
// @Summary 获取作品阅读热力图
// @Description 获取作品各章节的阅读热度分布
// @Tags Stats
// @Accept json
// @Produce json
// @Param book_id path string true "作品ID"
// @Success 200 {object} shared.APIResponse
// @Failure 400 {object} shared.APIResponse
// @Failure 500 {object} shared.APIResponse
// @Router /api/v1/writer/books/{book_id}/heatmap [get]
func (api *StatsApi) GetBookHeatmap(c *gin.Context) {
	bookID := c.Param("book_id")
	if bookID == "" {
		response.BadRequest(c,  "参数错误", "作品ID不能为空")
		return
	}

	// 生成热力图
	heatmap, err := api.statsService.GenerateHeatmap(c.Request.Context(), bookID)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	shared.Success(c, http.StatusOK, "获取成功", heatmap)
}

// GetBookRevenue 获取收入统计
// @Summary 获取作品收入统计
// @Description 获取作品的收入细分数据
// @Tags Stats
// @Accept json
// @Produce json
// @Param book_id path string true "作品ID"
// @Param start_date query string false "开始日期 (YYYY-MM-DD)"
// @Param end_date query string false "结束日期 (YYYY-MM-DD)"
// @Success 200 {object} shared.APIResponse
// @Failure 400 {object} shared.APIResponse
// @Failure 500 {object} shared.APIResponse
// @Router /api/v1/writer/books/{book_id}/revenue [get]
func (api *StatsApi) GetBookRevenue(c *gin.Context) {
	bookID := c.Param("book_id")
	if bookID == "" {
		response.BadRequest(c,  "参数错误", "作品ID不能为空")
		return
	}

	// 解析日期参数
	startDateStr := c.DefaultQuery("start_date", "")
	endDateStr := c.DefaultQuery("end_date", "")

	var startDate, endDate time.Time
	var err error

	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			response.BadRequest(c,  "参数错误", "开始日期格式错误")
			return
		}
	} else {
		// 默认最近30天
		startDate = time.Now().AddDate(0, 0, -30)
	}

	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			response.BadRequest(c,  "参数错误", "结束日期格式错误")
			return
		}
	} else {
		endDate = time.Now()
	}

	// 获取收入细分
	revenueBreakdown, err := api.statsService.GetRevenueBreakdown(c.Request.Context(), bookID, startDate, endDate)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	shared.Success(c, http.StatusOK, "获取成功", revenueBreakdown)
}

// GetTopChapters 获取热门章节
// @Summary 获取热门章节列表
// @Description 获取作品的热门章节统计（阅读量最高、收入最高、完读率最低、跳出率最高）
// @Tags Stats
// @Accept json
// @Produce json
// @Param book_id path string true "作品ID"
// @Success 200 {object} shared.APIResponse
// @Failure 400 {object} shared.APIResponse
// @Failure 500 {object} shared.APIResponse
// @Router /api/v1/writer/books/{book_id}/top-chapters [get]
func (api *StatsApi) GetTopChapters(c *gin.Context) {
	bookID := c.Param("book_id")
	if bookID == "" {
		response.BadRequest(c,  "参数错误", "作品ID不能为空")
		return
	}

	// 获取热门章节
	topChapters, err := api.statsService.GetTopChapters(c.Request.Context(), bookID)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	shared.Success(c, http.StatusOK, "获取成功", topChapters)
}

// GetDailyStats 获取每日统计
// @Summary 获取作品每日统计
// @Description 获取作品最近N天的每日统计数据
// @Tags Stats
// @Accept json
// @Produce json
// @Param book_id path string true "作品ID"
// @Param days query int false "天数" default(7)
// @Success 200 {object} shared.APIResponse
// @Failure 400 {object} shared.APIResponse
// @Failure 500 {object} shared.APIResponse
// @Router /api/v1/writer/books/{book_id}/daily-stats [get]
func (api *StatsApi) GetDailyStats(c *gin.Context) {
	bookID := c.Param("book_id")
	if bookID == "" {
		response.BadRequest(c,  "参数错误", "作品ID不能为空")
		return
	}

	// 解析天数参数
	daysStr := c.DefaultQuery("days", "7")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days < 1 || days > 365 {
		response.BadRequest(c,  "参数错误", "天数必须在1-365之间")
		return
	}

	// 获取每日统计
	dailyStats, err := api.statsService.GetDailyStats(c.Request.Context(), bookID, days)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	shared.Success(c, http.StatusOK, "获取成功", dailyStats)
}

// GetDropOffPoints 获取跳出点分析
// @Summary 获取作品跳出点分析
// @Description 获取跳出率最高的章节列表
// @Tags Stats
// @Accept json
// @Produce json
// @Param book_id path string true "作品ID"
// @Success 200 {object} shared.APIResponse
// @Failure 400 {object} shared.APIResponse
// @Failure 500 {object} shared.APIResponse
// @Router /api/v1/writer/books/{book_id}/drop-off-points [get]
func (api *StatsApi) GetDropOffPoints(c *gin.Context) {
	bookID := c.Param("book_id")
	if bookID == "" {
		response.BadRequest(c,  "参数错误", "作品ID不能为空")
		return
	}

	// 获取跳出点
	dropOffPoints, err := api.statsService.CalculateDropOffPoints(c.Request.Context(), bookID)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	shared.Success(c, http.StatusOK, "获取成功", dropOffPoints)
}

// RecordBehavior 记录读者行为
// @Summary 记录读者行为
// @Description 记录读者的阅读行为（浏览、完读、跳出等）
// @Tags Stats
// @Accept json
// @Produce json
// @Param request body object true "读者行为数据"
// @Success 200 {object} shared.APIResponse
// @Failure 400 {object} shared.APIResponse
// @Failure 500 {object} shared.APIResponse
// @Router /api/v1/reader/behavior [post]
func (api *StatsApi) RecordBehavior(c *gin.Context) {
	var behavior stats.ReaderBehavior
	if err := c.ShouldBindJSON(&behavior); err != nil {
		response.BadRequest(c,  "参数错误", err.Error())
		return
	}

	// 从context获取用户ID
	userID, exists := c.Get("userId")
	if exists {
		behavior.UserID = userID.(string)
	}

	// 记录行为
	err := api.statsService.RecordReaderBehavior(c.Request.Context(), &behavior)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	shared.Success(c, http.StatusOK, "记录成功", nil)
}

// GetRetentionRate 获取留存率
// @Summary 获取作品留存率
// @Description 获取作品的N日留存率
// @Tags Stats
// @Accept json
// @Produce json
// @Param book_id path string true "作品ID"
// @Param days query int false "天数" default(7)
// @Success 200 {object} shared.APIResponse
// @Failure 400 {object} shared.APIResponse
// @Failure 500 {object} shared.APIResponse
// @Router /api/v1/writer/books/{book_id}/retention [get]
func (api *StatsApi) GetRetentionRate(c *gin.Context) {
	bookID := c.Param("book_id")
	if bookID == "" {
		response.BadRequest(c,  "参数错误", "作品ID不能为空")
		return
	}

	// 解析天数参数
	daysStr := c.DefaultQuery("days", "7")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days < 1 || days > 90 {
		response.BadRequest(c,  "参数错误", "天数必须在1-90之间")
		return
	}

	// 计算留存率
	retentionRate, err := api.statsService.CalculateRetention(c.Request.Context(), bookID, days)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	result := map[string]interface{}{
		"book_id":        bookID,
		"days":           days,
		"retention_rate": retentionRate,
	}

	shared.Success(c, http.StatusOK, "获取成功", result)
}
