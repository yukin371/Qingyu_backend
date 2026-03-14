package writer

import (
	"strconv"

	"github.com/gin-gonic/gin"

	statsModels "Qingyu_backend/models/stats"
	"Qingyu_backend/pkg/response"
	bookstoreRepo "Qingyu_backend/repository/interfaces/bookstore"
	readingStats "Qingyu_backend/service/reader/stats"
)

// WriterStatsAggregateAPI 提供面向作者工作台的聚合统计接口。
type WriterStatsAggregateAPI struct {
	statsService *readingStats.ReadingStatsService
	bookRepo     bookstoreRepo.BookRepository
}

type WriterStatsOverviewSummary struct {
	Views       int64 `json:"views"`
	Subscribers int64 `json:"subscribers"`
	Bookmarks   int64 `json:"bookmarks"`
	Comments    int64 `json:"comments"`
}

type WriterStatsTodaySummary struct {
	Views       int64 `json:"views"`
	Subscribers int64 `json:"subscribers"`
	Words       int64 `json:"words"`
}

type WriterStatsOverviewResponse struct {
	BookID        string                     `json:"bookId"`
	Title         string                     `json:"title"`
	TotalViews    int64                      `json:"totalViews"`
	Subscribers   int64                      `json:"subscribers"`
	Favorites     int64                      `json:"favorites"`
	Comments      int64                      `json:"comments"`
	TodayViews    int64                      `json:"todayViews"`
	MonthViews    int64                      `json:"monthViews"`
	TodayWords    int64                      `json:"todayWords"`
	WordCount     int64                      `json:"wordCount"`
	TotalRevenue  float64                    `json:"totalRevenue"`
	RetentionRate float64                    `json:"retentionRate"`
	Overview      WriterStatsOverviewSummary `json:"overview"`
	Today         WriterStatsTodaySummary    `json:"today"`
}

type WriterStatsTrendItem struct {
	Date           interface{} `json:"date"`
	Views          int64       `json:"views,omitempty"`
	DailyViews     int64       `json:"dailyViews,omitempty"`
	NewSubscribers int64       `json:"newSubscribers,omitempty"`
	Subscribers    int64       `json:"subscribers,omitempty"`
	Count          int64       `json:"count,omitempty"`
	ChapterID      string      `json:"chapterId,omitempty"`
	ChapterTitle   string      `json:"chapterTitle,omitempty"`
	Title          string      `json:"title,omitempty"`
	Reads          int64       `json:"reads,omitempty"`
	CompletionRate float64     `json:"completionRate,omitempty"`
	DropOffRate    float64     `json:"dropOffRate,omitempty"`
	Revenue        float64     `json:"revenue,omitempty"`
}

type WriterStatsTrendResponse struct {
	BookID string                 `json:"bookId"`
	Days   int                    `json:"days,omitempty"`
	Page   int                    `json:"page,omitempty"`
	Size   int                    `json:"size,omitempty"`
	Total  int64                  `json:"total"`
	Items  []WriterStatsTrendItem `json:"items"`
	List   []WriterStatsTrendItem `json:"list"`
}

type WriterStatsTodayResponse struct {
	BookID           string `json:"bookId"`
	TodayViews       int64  `json:"todayViews"`
	TodaySubscribers int64  `json:"todaySubscribers"`
	TodayWords       int64  `json:"todayWords"`
	Views            int64  `json:"views"`
	Subscribers      int64  `json:"subscribers"`
	Words            int64  `json:"words"`
}

// NewWriterStatsAggregateAPI 创建作者聚合统计 API。
func NewWriterStatsAggregateAPI(statsService *readingStats.ReadingStatsService, bookRepo bookstoreRepo.BookRepository) *WriterStatsAggregateAPI {
	return &WriterStatsAggregateAPI{
		statsService: statsService,
		bookRepo:     bookRepo,
	}
}

func (api *WriterStatsAggregateAPI) resolveAggregateBookID(c *gin.Context) (string, bool) {
	if bookID := c.Query("bookId"); bookID != "" {
		return bookID, true
	}

	if projectID := c.Query("projectId"); projectID != "" {
		if api.bookRepo == nil {
			return projectID, true
		}

		book, err := api.bookRepo.GetByProjectID(c.Request.Context(), projectID)
		if err != nil {
			c.Error(err)
			return "", false
		}
		if book != nil {
			return book.ID.Hex(), true
		}

		return projectID, true
	}

	userIDValue, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未登录")
		return "", false
	}

	userID, ok := userIDValue.(string)
	if !ok || userID == "" {
		response.Unauthorized(c, "未登录")
		return "", false
	}

	if api.bookRepo == nil {
		response.BadRequest(c, "参数错误", "缺少bookId或projectId")
		return "", false
	}

	books, err := api.bookRepo.GetByAuthorID(c.Request.Context(), userID, 1, 0)
	if err != nil {
		c.Error(err)
		return "", false
	}
	if len(books) == 0 {
		response.Success(c, gin.H{
			"bookId":        "",
			"title":         "",
			"totalViews":    0,
			"subscribers":   0,
			"favorites":     0,
			"comments":      0,
			"todayViews":    0,
			"monthViews":    0,
			"wordCount":     0,
			"todayWords":    0,
			"totalRevenue":  0,
			"retentionRate": 0,
		})
		return "", false
	}

	return books[0].ID.Hex(), true
}

func sumViewMetrics(items []*statsModels.BookStatsDaily) (int64, int64) {
	var views int64
	var subscribers int64
	for _, item := range items {
		views += item.DailyViews
		subscribers += item.DailySubscribers
	}

	return views, subscribers
}

// GetOverview 获取作者工作台概览统计。
// @Summary 获取作者统计概览
// @Description 获取作者工作台概览统计，支持通过 bookId 或 projectId 指定作品
// @Tags Writer-Stats
// @Accept json
// @Produce json
// @Param bookId query string false "作品ID"
// @Param projectId query string false "项目ID"
// @Success 200 {object} response.APIResponse{data=WriterStatsOverviewResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/writer/stats/overview [get]
// GetOverview 获取作者工作台概览统计。
func (api *WriterStatsAggregateAPI) GetOverview(c *gin.Context) {
	bookID, ok := api.resolveAggregateBookID(c)
	if !ok {
		return
	}

	bookStats, err := api.statsService.CalculateBookStats(c.Request.Context(), bookID)
	if err != nil {
		c.Error(err)
		return
	}
	if bookStats == nil {
		response.Success(c, gin.H{})
		return
	}

	todayStats, err := api.statsService.GetDailyStats(c.Request.Context(), bookID, 1)
	if err != nil {
		c.Error(err)
		return
	}
	monthStats, err := api.statsService.GetDailyStats(c.Request.Context(), bookID, 30)
	if err != nil {
		c.Error(err)
		return
	}
	todayViews, todaySubscribers := sumViewMetrics(todayStats)
	monthViews, _ := sumViewMetrics(monthStats)

	retentionRate, err := api.statsService.CalculateRetention(c.Request.Context(), bookID, 7)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, gin.H{
		"bookId":        bookID,
		"title":         bookStats.Title,
		"totalViews":    bookStats.TotalViews,
		"subscribers":   bookStats.TotalSubscribers,
		"favorites":     bookStats.TotalBookmarks,
		"comments":      bookStats.TotalComments,
		"todayViews":    todayViews,
		"monthViews":    monthViews,
		"todayWords":    0,
		"wordCount":     bookStats.TotalWords,
		"totalRevenue":  bookStats.TotalRevenue,
		"retentionRate": retentionRate,
		"overview": gin.H{
			"views":       bookStats.TotalViews,
			"subscribers": bookStats.TotalSubscribers,
			"bookmarks":   bookStats.TotalBookmarks,
			"comments":    bookStats.TotalComments,
		},
		"today": gin.H{
			"views":       todayViews,
			"subscribers": todaySubscribers,
			"words":       0,
		},
	})
}

// GetViews 获取阅读量趋势统计。
// @Summary 获取作者阅读量趋势
// @Description 获取作者作品最近 N 天的阅读量趋势
// @Tags Writer-Stats
// @Accept json
// @Produce json
// @Param bookId query string false "作品ID"
// @Param projectId query string false "项目ID"
// @Param days query int false "天数" default(30)
// @Success 200 {object} response.APIResponse{data=WriterStatsTrendResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/writer/stats/views [get]
// GetViews 获取阅读量趋势统计。
func (api *WriterStatsAggregateAPI) GetViews(c *gin.Context) {
	bookID, ok := api.resolveAggregateBookID(c)
	if !ok {
		return
	}

	days, err := strconv.Atoi(c.DefaultQuery("days", "30"))
	if err != nil || days < 1 || days > 365 {
		response.BadRequest(c, "参数错误", "天数必须在1-365之间")
		return
	}

	dailyStats, err := api.statsService.GetDailyStats(c.Request.Context(), bookID, days)
	if err != nil {
		c.Error(err)
		return
	}

	items := make([]gin.H, 0, len(dailyStats))
	var total int64
	for _, item := range dailyStats {
		total += item.DailyViews
		items = append(items, gin.H{
			"date":       item.Date,
			"views":      item.DailyViews,
			"dailyViews": item.DailyViews,
		})
	}

	response.Success(c, gin.H{
		"bookId": bookID,
		"days":   days,
		"total":  total,
		"items":  items,
		"list":   items,
	})
}

// GetSubscribers 获取订阅趋势统计。
// @Summary 获取作者订阅趋势
// @Description 获取作者作品最近 N 天的订阅增长趋势
// @Tags Writer-Stats
// @Accept json
// @Produce json
// @Param bookId query string false "作品ID"
// @Param projectId query string false "项目ID"
// @Param days query int false "天数" default(30)
// @Success 200 {object} response.APIResponse{data=WriterStatsTrendResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/writer/stats/subscribers [get]
// GetSubscribers 获取订阅趋势统计。
func (api *WriterStatsAggregateAPI) GetSubscribers(c *gin.Context) {
	bookID, ok := api.resolveAggregateBookID(c)
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

	items := make([]gin.H, 0, len(dailyStats))
	var total int64
	for _, item := range dailyStats {
		total += item.DailySubscribers
		items = append(items, gin.H{
			"date":           item.Date,
			"newSubscribers": item.DailySubscribers,
			"subscribers":    item.DailySubscribers,
			"count":          item.DailySubscribers,
		})
	}

	response.Success(c, gin.H{
		"bookId": bookID,
		"days":   days,
		"total":  total,
		"items":  items,
		"list":   items,
	})
}

// GetChapters 获取章节表现统计。
// @Summary 获取作者章节统计
// @Description 获取作者作品的章节表现排行
// @Tags Writer-Stats
// @Accept json
// @Produce json
// @Param bookId query string false "作品ID"
// @Param projectId query string false "项目ID"
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(20)
// @Success 200 {object} response.APIResponse{data=WriterStatsTrendResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/writer/stats/chapters [get]
// GetChapters 获取章节表现统计。
func (api *WriterStatsAggregateAPI) GetChapters(c *gin.Context) {
	bookID, ok := api.resolveAggregateBookID(c)
	if !ok {
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 20
	}

	items, total, err := api.statsService.GetChapterRankings(c.Request.Context(), bookID, page, size)
	if err != nil {
		c.Error(err)
		return
	}

	result := make([]gin.H, 0, len(items))
	for _, item := range items {
		result = append(result, gin.H{
			"chapterId":      item.ChapterID,
			"chapterTitle":   item.Title,
			"title":          item.Title,
			"views":          item.ViewCount,
			"reads":          item.ViewCount,
			"completionRate": item.CompletionRate,
			"dropOffRate":    item.DropOffRate,
			"revenue":        item.Revenue,
		})
	}

	response.Success(c, gin.H{
		"bookId": bookID,
		"page":   page,
		"size":   size,
		"total":  total,
		"items":  result,
		"list":   result,
	})
}

// GetToday 获取当日统计。
// @Summary 获取作者今日统计
// @Description 获取作者作品当日阅读与订阅统计
// @Tags Writer-Stats
// @Accept json
// @Produce json
// @Param bookId query string false "作品ID"
// @Param projectId query string false "项目ID"
// @Success 200 {object} response.APIResponse{data=WriterStatsTodayResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/writer/stats/today [get]
// GetToday 获取当日统计。
func (api *WriterStatsAggregateAPI) GetToday(c *gin.Context) {
	bookID, ok := api.resolveAggregateBookID(c)
	if !ok {
		return
	}

	dailyStats, err := api.statsService.GetDailyStats(c.Request.Context(), bookID, 1)
	if err != nil {
		c.Error(err)
		return
	}

	todayViews, todaySubscribers := sumViewMetrics(dailyStats)

	response.Success(c, gin.H{
		"bookId":           bookID,
		"todayViews":       todayViews,
		"todaySubscribers": todaySubscribers,
		"todayWords":       0,
		"views":            todayViews,
		"subscribers":      todaySubscribers,
		"words":            0,
	})
}
