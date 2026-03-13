package writer

import (
	"math"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/models/stats"
	"Qingyu_backend/pkg/response"
	bookstoreRepo "Qingyu_backend/repository/interfaces/bookstore"
	readingStats "Qingyu_backend/service/reader/stats"
)

type compareBooksRequest struct {
	BookIDs   []string `json:"bookIds"`
	Metrics   []string `json:"metrics"`
	StartDate *string  `json:"startDate"`
	EndDate   *string  `json:"endDate"`
}

// StatsApi 阅读/书店统计API
// 职责：处理作品、章节的阅读统计
type StatsApi struct {
	statsService *readingStats.ReadingStatsService
	bookRepo     bookstoreRepo.BookRepository
}

// NewStatsApi 创建统计API
func NewStatsApi(statsService *readingStats.ReadingStatsService, bookRepo bookstoreRepo.BookRepository) *StatsApi {
	return &StatsApi{
		statsService: statsService,
		bookRepo:     bookRepo,
	}
}

func (api *StatsApi) resolveBookID(c *gin.Context) (string, bool) {
	rawID := c.Param("book_id")
	if rawID == "" {
		response.BadRequest(c, "参数错误", "作品ID不能为空")
		return "", false
	}

	if api.bookRepo == nil {
		return rawID, true
	}

	book, err := api.bookRepo.GetByID(c.Request.Context(), rawID)
	if err == nil && book != nil {
		return rawID, true
	}

	book, err = api.bookRepo.GetByProjectID(c.Request.Context(), rawID)
	if err != nil {
		c.Error(err)
		return "", false
	}
	if book != nil {
		return book.ID.Hex(), true
	}

	return rawID, true
}

// GetBookStats 获取作品统计
// @Summary 获取作品统计数据
// @Description 获取作品的完整统计信息，包括阅读、收入、互动等数据
// @Tags Stats
// @Accept json
// @Produce json
// @Param book_id path string true "作品ID"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/writer/books/{book_id}/stats [get]
func (api *StatsApi) GetBookStats(c *gin.Context) {
	bookID, ok := api.resolveBookID(c)
	if !ok {
		return
	}

	// 获取作品统计
	bookStats, err := api.statsService.CalculateBookStats(c.Request.Context(), bookID)
	if err != nil {
		c.Error(err)
		return
	}

	if bookStats == nil {
		response.Success(c, gin.H{
			"totalViews":  0,
			"subscribers": 0,
			"favorites":   0,
			"comments":    0,
			"todayViews":  0,
			"monthViews":  0,
			"wordCount":   0,
			"total_views": 0,
			"word_count":  0,
		})
		return
	}

	todayStats, _ := api.statsService.GetDailyStats(c.Request.Context(), bookID, 1)
	monthStats, _ := api.statsService.GetDailyStats(c.Request.Context(), bookID, 30)
	var todayViews int64
	var monthViews int64
	for _, item := range todayStats {
		todayViews += item.DailyViews
	}
	for _, item := range monthStats {
		monthViews += item.DailyViews
	}

	response.Success(c, gin.H{
		"bookId":       bookID,
		"title":        bookStats.Title,
		"totalViews":   bookStats.TotalViews,
		"subscribers":  bookStats.TotalSubscribers,
		"favorites":    bookStats.TotalBookmarks,
		"comments":     bookStats.TotalComments,
		"todayViews":   todayViews,
		"monthViews":   monthViews,
		"wordCount":    bookStats.TotalWords,
		"rating":       0,
		"total_views":  bookStats.TotalViews,
		"word_count":   bookStats.TotalWords,
		"totalRevenue": bookStats.TotalRevenue,
	})
}

// GetChapterStats 获取章节统计
// @Summary 获取章节统计数据
// @Description 获取单个章节的统计信息
// @Tags Stats
// @Accept json
// @Produce json
// @Param chapter_id path string true "章节ID"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/writer/chapters/{chapter_id}/stats [get]
func (api *StatsApi) GetChapterStats(c *gin.Context) {
	chapterID := c.Param("chapter_id")
	if chapterID == "" {
		response.BadRequest(c, "参数错误", "章节ID不能为空")
		return
	}

	// 获取章节统计
	chapterStats, err := api.statsService.CalculateChapterStats(c.Request.Context(), chapterID)
	if err != nil {
		c.Error(err)
		return
	}

	if chapterStats == nil {
		response.NotFound(c, "章节统计不存在")
		return
	}

	response.Success(c, chapterStats)
}

// GetBookHeatmap 获取阅读热力图
// @Summary 获取作品阅读热力图
// @Description 获取作品各章节的阅读热度分布
// @Tags Stats
// @Accept json
// @Produce json
// @Param book_id path string true "作品ID"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/writer/books/{book_id}/heatmap [get]
func (api *StatsApi) GetBookHeatmap(c *gin.Context) {
	bookID, ok := api.resolveBookID(c)
	if !ok {
		return
	}

	heatmap, err := api.statsService.GenerateReadingTimeHeatmap(c.Request.Context(), bookID, 90)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, heatmap)
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
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/writer/books/{book_id}/revenue [get]
func (api *StatsApi) GetBookRevenue(c *gin.Context) {
	bookID, ok := api.resolveBookID(c)
	if !ok {
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
			response.BadRequest(c, "参数错误", "开始日期格式错误")
			return
		}
	} else {
		// 默认最近30天
		startDate = time.Now().AddDate(0, 0, -30)
	}

	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			response.BadRequest(c, "参数错误", "结束日期格式错误")
			return
		}
	} else {
		endDate = time.Now()
	}

	// 获取收入细分
	revenueBreakdown, err := api.statsService.GetRevenueBreakdown(c.Request.Context(), bookID, startDate, endDate)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, revenueBreakdown)
}

// GetTopChapters 获取热门章节
// @Summary 获取热门章节列表
// @Description 获取作品的热门章节统计（阅读量最高、收入最高、完读率最低、跳出率最高）
// @Tags Stats
// @Accept json
// @Produce json
// @Param book_id path string true "作品ID"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/writer/books/{book_id}/top-chapters [get]
func (api *StatsApi) GetTopChapters(c *gin.Context) {
	bookID, ok := api.resolveBookID(c)
	if !ok {
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	topChapters, total, err := api.statsService.GetChapterRankings(c.Request.Context(), bookID, page, size)
	if err != nil {
		c.Error(err)
		return
	}

	result := make([]gin.H, 0, len(topChapters))
	for _, item := range topChapters {
		result = append(result, gin.H{
			"chapterId":      item.ChapterID,
			"chapterTitle":   item.Title,
			"title":          item.Title,
			"views":          item.ViewCount,
			"reads":          item.ViewCount,
			"comments":       0,
			"avgReadTime":    0,
			"retention":      math.Round(item.CompletionRate*10000) / 100,
			"completionRate": item.CompletionRate,
			"dropOffRate":    item.DropOffRate,
			"revenue":        item.Revenue,
		})
	}

	response.Success(c, gin.H{
		"items": result,
		"total": total,
		"list":  result,
	})
}

// GetDailyStats 获取每日统计
// @Summary 获取作品每日统计
// @Description 获取作品最近N天的每日统计数据
// @Tags Stats
// @Accept json
// @Produce json
// @Param book_id path string true "作品ID"
// @Param days query int false "天数" default(7)
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/writer/books/{book_id}/daily-stats [get]
func (api *StatsApi) GetDailyStats(c *gin.Context) {
	bookID, ok := api.resolveBookID(c)
	if !ok {
		return
	}

	// 解析天数参数
	daysStr := c.DefaultQuery("days", "7")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days < 1 || days > 365 {
		response.BadRequest(c, "参数错误", "天数必须在1-365之间")
		return
	}

	dailyStats, err := api.statsService.GetDailyStats(c.Request.Context(), bookID, days)
	if err != nil {
		c.Error(err)
		return
	}

	result := make([]gin.H, 0, len(dailyStats))
	for _, item := range dailyStats {
		result = append(result, gin.H{
			"date":             item.Date,
			"views":            item.DailyViews,
			"dailyViews":       item.DailyViews,
			"newSubscribers":   item.DailySubscribers,
			"dailySubscribers": item.DailySubscribers,
			"newFavorites":     int64(0),
			"comments":         int64(0),
		})
	}

	response.Success(c, result)
}

// GetDropOffPoints 获取跳出点分析
// @Summary 获取作品跳出点分析
// @Description 获取跳出率最高的章节列表
// @Tags Stats
// @Accept json
// @Produce json
// @Param book_id path string true "作品ID"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/writer/books/{book_id}/drop-off-points [get]
func (api *StatsApi) GetDropOffPoints(c *gin.Context) {
	bookID, ok := api.resolveBookID(c)
	if !ok {
		return
	}

	// 获取跳出点
	dropOffPoints, err := api.statsService.CalculateDropOffPoints(c.Request.Context(), bookID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, dropOffPoints)
}

// RecordBehavior 记录读者行为
// @Summary 记录读者行为
// @Description 记录读者的阅读行为（浏览、完读、跳出等）
// @Tags Stats
// @Accept json
// @Produce json
// @Param request body object true "读者行为数据"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/reader/behavior [post]
func (api *StatsApi) RecordBehavior(c *gin.Context) {
	var behavior stats.ReaderBehavior
	if err := c.ShouldBindJSON(&behavior); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	// 从context获取用户ID
	userID, exists := c.Get("user_id")
	if exists {
		behavior.UserID = userID.(string)
	}

	// 记录行为
	err := api.statsService.RecordReaderBehavior(c.Request.Context(), &behavior)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, nil)
}

// GetRetentionRate 获取留存率
// @Summary 获取作品留存率
// @Description 获取作品的N日留存率
// @Tags Stats
// @Accept json
// @Produce json
// @Param book_id path string true "作品ID"
// @Param days query int false "天数" default(7)
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/writer/books/{book_id}/retention [get]
func (api *StatsApi) GetRetentionRate(c *gin.Context) {
	bookID, ok := api.resolveBookID(c)
	if !ok {
		return
	}

	// 解析天数参数
	daysStr := c.DefaultQuery("days", "7")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days < 1 || days > 90 {
		response.BadRequest(c, "参数错误", "天数必须在1-90之间")
		return
	}

	// 计算留存率
	retentionRate, err := api.statsService.CalculateRetention(c.Request.Context(), bookID, days)
	if err != nil {
		c.Error(err)
		return
	}

	result := map[string]interface{}{
		"book_id":        bookID,
		"days":           days,
		"retention_rate": retentionRate,
	}

	response.Success(c, result)
}

// GetSubscribersTrend 获取订阅增长趋势
func (api *StatsApi) GetSubscribersTrend(c *gin.Context) {
	bookID, ok := api.resolveBookID(c)
	if !ok {
		return
	}

	days, err := strconv.Atoi(c.DefaultQuery("days", "30"))
	if err != nil || days < 1 || days > 365 {
		response.BadRequest(c, "参数错误", "天数必须在1-365之间")
		return
	}

	dailyStats, err := api.statsService.GetSubscribersTrend(c.Request.Context(), bookID, days)
	if err != nil {
		c.Error(err)
		return
	}

	result := make([]gin.H, 0, len(dailyStats))
	for _, item := range dailyStats {
		result = append(result, gin.H{
			"date":           item.Date,
			"newSubscribers": item.DailySubscribers,
			"subscribers":    item.DailySubscribers,
			"count":          item.DailySubscribers,
		})
	}

	response.Success(c, result)
}

// GetReaderActivity 获取读者活跃度分布
func (api *StatsApi) GetReaderActivity(c *gin.Context) {
	bookID, ok := api.resolveBookID(c)
	if !ok {
		return
	}

	activity, err := api.statsService.GetReaderActivity(c.Request.Context(), bookID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, activity)
}

// CompareBooks 对比多个作品统计
func (api *StatsApi) CompareBooks(c *gin.Context) {
	var req compareBooksRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}
	if len(req.BookIDs) == 0 {
		response.BadRequest(c, "参数错误", "bookIds不能为空")
		return
	}
	if len(req.Metrics) == 0 {
		req.Metrics = []string{"views", "subscribers", "favorites", "comments", "revenue"}
	}

	var startDate *time.Time
	var endDate *time.Time
	if req.StartDate != nil && *req.StartDate != "" {
		parsed, err := time.Parse(time.RFC3339, *req.StartDate)
		if err != nil {
			if parsedDate, parseErr := time.Parse("2006-01-02", *req.StartDate); parseErr == nil {
				parsed = parsedDate
			} else {
				response.BadRequest(c, "参数错误", "startDate格式错误")
				return
			}
		}
		startDate = &parsed
	}
	if req.EndDate != nil && *req.EndDate != "" {
		parsed, err := time.Parse(time.RFC3339, *req.EndDate)
		if err != nil {
			if parsedDate, parseErr := time.Parse("2006-01-02", *req.EndDate); parseErr == nil {
				parsed = parsedDate
			} else {
				response.BadRequest(c, "参数错误", "endDate格式错误")
				return
			}
		}
		endDate = &parsed
	}

	result, err := api.statsService.CompareBooks(c.Request.Context(), req.BookIDs, req.Metrics, startDate, endDate)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, gin.H{
		"items":   result,
		"metrics": req.Metrics,
	})
}
