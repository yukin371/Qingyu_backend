package stats

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/shared"
	readingStatsService "Qingyu_backend/service/reader/stats"
	"Qingyu_backend/pkg/response"
)

// ReadingStatsAPI 阅读统计API处理器
type ReadingStatsAPI struct {
	statsService *readingStatsService.ReadingStatsService
}

// NewReadingStatsAPI 创建阅读统计API实例
func NewReadingStatsAPI(statsService *readingStatsService.ReadingStatsService) *ReadingStatsAPI {
	return &ReadingStatsAPI{
		statsService: statsService,
	}
}

// GetMyStats 获取我的阅读统计
// @Summary 获取我的阅读统计
// @Tags 阅读统计
// @Accept json
// @Produce json
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/reading-stats/my/stats [get]
// @Security Bearer
func (api *ReadingStatsAPI) GetMyStats(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	// 获取过去30天的统计数据
	days := 30
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)

	// 获取时间范围统计
	stats, err := api.statsService.GetTimeRangeStats(c.Request.Context(), "", startDate, endDate)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	shared.Success(c, http.StatusOK, "获取统计成功", stats)
}

// GetMyDailyStats 获取我的每日阅读统计
// @Summary 获取我的每日阅读统计
// @Tags 阅读统计
// @Accept json
// @Produce json
// @Param days query int false "天数" default(30)
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/reading-stats/my/daily [get]
// @Security Bearer
func (api *ReadingStatsAPI) GetMyDailyStats(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	daysStr := c.DefaultQuery("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days <= 0 || days > 365 {
		days = 30
	}

	// 获取每日统计
	dailyStats, err := api.statsService.GetDailyStats(c.Request.Context(), "", days)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	shared.Success(c, http.StatusOK, "获取每日统计成功", gin.H{
		"days":  days,
		"stats": dailyStats,
	})
}

// GetMyRanking 获取我的阅读排名
// @Summary 获取我的阅读排名
// @Tags 阅读统计
// @Accept json
// @Produce json
// @Param period query string false "周期(week/month/all)" default(week)
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/reading-stats/my/ranking [get]
// @Security Bearer
func (api *ReadingStatsAPI) GetMyRanking(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	period := c.DefaultQuery("period", "week")

	// 这里调用排名服务（需要实现）
	// 暂时返回模拟数据
	ranking := gin.H{
		"period":  period,
		"rank":    0,
		"total":   0,
		"message": "排名功能待实现",
	}

	shared.Success(c, http.StatusOK, "获取排名成功", ranking)
}

// GetMyReadingTime 获取我的阅读时长统计
// @Summary 获取我的阅读时长统计
// @Tags 阅读统计
// @Accept json
// @Produce json
// @Param days query int false "天数" default(7)
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/reading-stats/my/reading-time [get]
// @Security Bearer
func (api *ReadingStatsAPI) GetMyReadingTime(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	daysStr := c.DefaultQuery("days", "7")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days <= 0 || days > 365 {
		days = 7
	}

	// 计算时间范围
	_ = time.Now()
	_ = time.Now().AddDate(0, 0, -days)

	// 这里调用阅读时长服务（需要从 reader 服务获取）
	// 暂时返回模拟数据
	readingTime := gin.H{
		"days":         days,
		"totalSeconds": 0,
		"totalHours":   0,
		"avgDaily":     0,
		"message":      "功能待实现",
	}

	shared.Success(c, http.StatusOK, "获取阅读时长成功", readingTime)
}

// GetMyHistory 获取我的阅读历史统计
// @Summary 获取我的阅读历史统计
// @Tags 阅读统计
// @Accept json
// @Produce json
// @Param days query int false "天数" default(30)
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/reading-stats/my/history [get]
// @Security Bearer
func (api *ReadingStatsAPI) GetMyHistory(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	daysStr := c.DefaultQuery("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days <= 0 || days > 365 {
		days = 30
	}

	// 这里调用阅读历史服务
	// 暂时返回模拟数据
	history := gin.H{
		"days":         days,
		"booksRead":    0,
		"chaptersRead": 0,
		"totalWords":   0,
		"totalTime":    0,
		"message":      "功能待实现",
	}

	shared.Success(c, http.StatusOK, "获取阅读历史成功", history)
}

// GetRecommendations 获取阅读推荐（基于统计）
// @Summary 获取阅读推荐
// @Tags 阅读统计
// @Accept json
// @Produce json
// @Param type query string false "推荐类型(hot/new/personal)" default(personal)
// @Param limit query int false "数量限制" default(10)
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/reading-stats/recommendations [get]
// @Security Bearer
func (api *ReadingStatsAPI) GetRecommendations(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	recommendType := c.DefaultQuery("type", "personal")
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 50 {
		limit = 10
	}

	// 这里调用推荐服务
	// 暂时返回模拟数据
	recommendations := gin.H{
		"type":    recommendType,
		"limit":   limit,
		"books":   []interface{}{},
		"message": "推荐功能待实现",
	}

	shared.Success(c, http.StatusOK, "获取推荐成功", recommendations)
}
