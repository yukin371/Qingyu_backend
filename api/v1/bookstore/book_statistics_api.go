package bookstore

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/api/v1/shared"
	bookstoreService "Qingyu_backend/service/bookstore"
)

// BookStatisticsAPI 图书统计API控制器
type BookStatisticsAPI struct {
	BookStatisticsService bookstoreService.BookStatisticsService
}

// NewBookStatisticsAPI 创建新的图书统计API实例
func NewBookStatisticsAPI(bookStatisticsService bookstoreService.BookStatisticsService) *BookStatisticsAPI {
	return &BookStatisticsAPI{
		BookStatisticsService: bookStatisticsService,
	}
}

// GetBookStatistics 获取图书统计信息
// @Summary 获取图书统计信息
// @Description 根据图书ID获取统计信息
// @Tags 图书统计
// @Accept json
// @Produce json
// @Param book_id path string true "图书ID"
// @Success 200 {object} APIResponse "成功"
// @Failure 400 {object} APIResponse "请求参数错误"
// @Failure 404 {object} APIResponse "统计信息不存在"
// @Failure 500 {object} APIResponse "服务器内部错误"
// @Router /api/v1/reading/books/{book_id}/statistics [get]
func (api *BookStatisticsAPI) GetBookStatistics(c *gin.Context) {
	bookIDStr := c.Param("book_id")
	if bookIDStr == "" {
		shared.BadRequest(c, "参数错误", "图书ID不能为空")
		return
	}

	bookID, err := primitive.ObjectIDFromHex(bookIDStr)
	if err != nil {
		shared.BadRequest(c, "参数错误", "无效的图书ID格式")
		return
	}

	statistics, err := api.BookStatisticsService.GetStatisticsByBookID(c.Request.Context(), bookID.Hex())
	if err != nil {
		shared.NotFound(c, "统计信息不存在")
		return
	}

	shared.SuccessData(c, statistics)
}

// GetTopViewedBooks 获取最多浏览的图书
// @Summary 获取最多浏览的图书
// @Description 获取浏览量最高的图书列表
// @Tags 图书统计
// @Accept json
// @Produce json
// @Param limit query int false "返回数量" default(10)
// @Success 200 {object} APIResponse "成功"
// @Failure 400 {object} APIResponse "请求参数错误"
// @Failure 500 {object} APIResponse "服务器内部错误"
// @Router /api/v1/reading/statistics/top-viewed [get]
func (api *BookStatisticsAPI) GetTopViewedBooks(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit < 1 || limit > 100 {
		limit = 10
	}

	statistics, err := api.BookStatisticsService.GetTopViewedBooks(c.Request.Context(), limit)
	if err != nil {
		c.Error(err)
		return
	}

	shared.SuccessData(c, statistics)
}

// GetTopFavoritedBooks 获取最多收藏的图书
// @Summary 获取最多收藏的图书
// @Description 获取收藏量最高的图书列表
// @Tags 图书统计
// @Accept json
// @Produce json
// @Param limit query int false "返回数量" default(10)
// @Success 200 {object} APIResponse "成功"
// @Failure 400 {object} APIResponse "请求参数错误"
// @Failure 500 {object} APIResponse "服务器内部错误"
// @Router /api/v1/reading/statistics/top-favorited [get]
func (api *BookStatisticsAPI) GetTopFavoritedBooks(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit < 1 || limit > 100 {
		limit = 10
	}

	statistics, err := api.BookStatisticsService.GetTopFavoritedBooks(c.Request.Context(), limit)
	if err != nil {
		c.Error(err)
		return
	}

	shared.SuccessData(c, statistics)
}

// GetTopRatedBooks 获取最高评分的图书
// @Summary 获取最高评分的图书
// @Description 获取评分最高的图书列表
// @Tags 图书统计
// @Accept json
// @Produce json
// @Param limit query int false "返回数量" default(10)
// @Success 200 {object} APIResponse "成功"
// @Failure 400 {object} APIResponse "请求参数错误"
// @Failure 500 {object} APIResponse "服务器内部错误"
// @Router /api/v1/reading/statistics/top-rated [get]
func (api *BookStatisticsAPI) GetTopRatedBooks(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit < 1 || limit > 100 {
		limit = 10
	}

	statistics, err := api.BookStatisticsService.GetTopRatedBooks(c.Request.Context(), limit)
	if err != nil {
		c.Error(err)
		return
	}

	shared.SuccessData(c, statistics)
}

// GetHottestBooks 获取最热门的图书
// @Summary 获取最热门的图书
// @Description 获取热度最高的图书列表
// @Tags 图书统计
// @Accept json
// @Produce json
// @Param limit query int false "返回数量" default(10)
// @Success 200 {object} APIResponse "成功"
// @Failure 400 {object} APIResponse "请求参数错误"
// @Failure 500 {object} APIResponse "服务器内部错误"
// @Router /api/v1/reading/statistics/hottest [get]
func (api *BookStatisticsAPI) GetHottestBooks(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit < 1 || limit > 100 {
		limit = 10
	}

	statistics, err := api.BookStatisticsService.GetHottestBooks(c.Request.Context(), limit)
	if err != nil {
		c.Error(err)
		return
	}

	shared.SuccessData(c, statistics)
}

// GetTrendingBooks 获取趋势图书
// @Summary 获取趋势图书
// @Description 获取当前趋势图书列表
// @Tags 图书统计
// @Accept json
// @Produce json
// @Param limit query int false "返回数量" default(10)
// @Success 200 {object} APIResponse "成功"
// @Failure 400 {object} APIResponse "请求参数错误"
// @Failure 500 {object} APIResponse "服务器内部错误"
// @Router /api/v1/reading/statistics/trending [get]
func (api *BookStatisticsAPI) GetTrendingBooks(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit < 1 || limit > 100 {
		limit = 10
	}

	statistics, err := api.BookStatisticsService.GetTrendingBooks(c.Request.Context(), 0, limit)
	if err != nil {
		c.Error(err)
		return
	}

	shared.SuccessData(c, statistics)
}

// IncrementViewCount 增加浏览量
// @Summary 增加浏览量
// @Description 为指定图书增加浏览量
// @Tags 图书统计
// @Accept json
// @Produce json
// @Param book_id path string true "图书ID"
// @Success 200 {object} APIResponse "增加成功"
// @Failure 400 {object} APIResponse "请求参数错误"
// @Failure 500 {object} APIResponse "服务器内部错误"
// @Router /api/v1/reading/books/{book_id}/view [post]
func (api *BookStatisticsAPI) IncrementViewCount(c *gin.Context) {
	bookIDStr := c.Param("book_id")
	if bookIDStr == "" {
		shared.BadRequest(c, "参数错误", "图书ID不能为空")
		return
	}

	bookID, err := primitive.ObjectIDFromHex(bookIDStr)
	if err != nil {
		shared.BadRequest(c, "参数错误", "无效的图书ID格式")
		return
	}

	err = api.BookStatisticsService.IncrementViewCount(c.Request.Context(), bookID.Hex())
	if err != nil {
		c.Error(err)
		return
	}

	shared.SuccessData(c, nil)
}

// IncrementFavoriteCount 增加收藏量
// @Summary 增加收藏量
// @Description 为指定图书增加收藏量
// @Tags 图书统计
// @Accept json
// @Produce json
// @Param book_id path string true "图书ID"
// @Success 200 {object} APIResponse "增加成功"
// @Failure 400 {object} APIResponse "请求参数错误"
// @Failure 500 {object} APIResponse "服务器内部错误"
// @Router /api/v1/reading/books/{book_id}/favorite [post]
func (api *BookStatisticsAPI) IncrementFavoriteCount(c *gin.Context) {
	bookIDStr := c.Param("book_id")
	if bookIDStr == "" {
		shared.BadRequest(c, "参数错误", "图书ID不能为空")
		return
	}

	bookID, err := primitive.ObjectIDFromHex(bookIDStr)
	if err != nil {
		shared.BadRequest(c, "参数错误", "无效的图书ID格式")
		return
	}

	err = api.BookStatisticsService.IncrementFavoriteCount(c.Request.Context(), bookID.Hex())
	if err != nil {
		c.Error(err)
		return
	}

	shared.SuccessData(c, nil)
}

// GetAggregatedStatistics 获取聚合统计信息
// @Summary 获取聚合统计信息
// @Description 获取系统整体的聚合统计信息
// @Tags 图书统计
// @Accept json
// @Produce json
// @Success		200	{object}	APIResponse
// @Failure 500 {object} APIResponse "服务器内部错误"
// @Router /api/v1/reading/statistics/aggregated [get]
func (api *BookStatisticsAPI) GetAggregatedStatistics(c *gin.Context) {
	statistics, err := api.BookStatisticsService.GetAggregatedStatistics(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	shared.SuccessData(c, statistics)
}

// GetStatisticsByTimeRange 获取时间范围内的统计信息
// @Summary 获取时间范围内的统计信息
// @Description 获取指定时间范围内的统计信息
// @Tags 图书统计
// @Accept json
// @Produce json
// @Param start_time query string true "开始时间" format(date-time)
// @Param end_time query string true "结束时间" format(date-time)
// @Param page query int false "页码" default(1)
// @Param limit query int false "每页数量" default(10)
// @Success 200 {object} APIResponse "成功"
// @Failure 400 {object} APIResponse "请求参数错误"
// @Failure 500 {object} APIResponse "服务器内部错误"
// @Router /api/v1/reading/statistics/time-range [get]
func (api *BookStatisticsAPI) GetStatisticsByTimeRange(c *gin.Context) {
	startTimeStr := c.Query("start_time")
	endTimeStr := c.Query("end_time")

	if startTimeStr == "" || endTimeStr == "" {
		shared.BadRequest(c, "参数错误", "开始时间和结束时间不能为空")
		return
	}

	startTime, err := time.Parse(time.RFC3339, startTimeStr)
	if err != nil {
		shared.BadRequest(c, "参数错误", "开始时间格式错误")
		return
	}

	endTime, err := time.Parse(time.RFC3339, endTimeStr)
	if err != nil {
		shared.BadRequest(c, "参数错误", "结束时间格式错误")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	statistics, err := api.BookStatisticsService.GetStatisticsByTimeRange(c.Request.Context(), startTime, endTime)
	if err != nil {
		c.Error(err)
		return
	}

	total := int64(len(statistics))

	shared.Paginated(c, statistics, total, page, limit, "获取成功")
}

// GetDailyStatisticsReport 获取日统计报告
// @Summary 获取日统计报告
// @Description 获取指定日期的统计报告
// @Tags 图书统计
// @Accept json
// @Produce json
// @Param date query string true "日期" format(date)
// @Success		200	{object}	APIResponse
// @Failure 400 {object} APIResponse "请求参数错误"
// @Failure 500 {object} APIResponse "服务器内部错误"
// @Router /api/v1/reading/statistics/daily-report [get]
func (api *BookStatisticsAPI) GetDailyStatisticsReport(c *gin.Context) {
	dateStr := c.Query("date")
	if dateStr == "" {
		shared.BadRequest(c, "参数错误", "日期不能为空")
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		shared.BadRequest(c, "参数错误", "日期格式错误")
		return
	}

	report, err := api.BookStatisticsService.GenerateDailyReport(c.Request.Context(), date)
	if err != nil {
		c.Error(err)
		return
	}

	shared.SuccessData(c, report)
}

// GetWeeklyStatisticsReport 获取周统计报告
// @Summary 获取周统计报告
// @Description 获取指定周的统计报告
// @Tags 图书统计
// @Accept json
// @Produce json
// @Param year query int true "年份"
// @Param week query int true "周数"
// @Success		200	{object}	APIResponse
// @Failure 400 {object} APIResponse "请求参数错误"
// @Failure 500 {object} APIResponse "服务器内部错误"
// @Router /api/v1/reading/statistics/weekly-report [get]
func (api *BookStatisticsAPI) GetWeeklyStatisticsReport(c *gin.Context) {
	yearStr := c.Query("year")
	weekStr := c.Query("week")

	if yearStr == "" || weekStr == "" {
		shared.BadRequest(c, "参数错误", "年份和周数不能为空")
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		shared.BadRequest(c, "参数错误", "年份格式错误")
		return
	}

	week, err := strconv.Atoi(weekStr)
	if err != nil {
		shared.BadRequest(c, "参数错误", "周数格式错误")
		return
	}

	// 计算周的起始日期
	startDate := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	// 找到第一个星期一
	for startDate.Weekday() != time.Monday {
		startDate = startDate.AddDate(0, 0, 1)
	}
	// 加上周数
	startDate = startDate.AddDate(0, 0, (week-1)*7)

	report, err := api.BookStatisticsService.GenerateWeeklyReport(c.Request.Context(), startDate)
	if err != nil {
		c.Error(err)
		return
	}

	shared.SuccessData(c, report)
}

// GetMonthlyStatisticsReport 获取月统计报告
// @Summary 获取月统计报告
// @Description 获取指定月份的统计报告
// @Tags 图书统计
// @Accept json
// @Produce json
// @Param year query int true "年份"
// @Param month query int true "月份"
// @Success		200	{object}	APIResponse
// @Failure 400 {object} APIResponse "请求参数错误"
// @Failure 500 {object} APIResponse "服务器内部错误"
// @Router /api/v1/reading/statistics/monthly-report [get]
func (api *BookStatisticsAPI) GetMonthlyStatisticsReport(c *gin.Context) {
	yearStr := c.Query("year")
	monthStr := c.Query("month")

	if yearStr == "" || monthStr == "" {
		shared.BadRequest(c, "参数错误", "年份和月份不能为空")
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		shared.BadRequest(c, "参数错误", "年份格式错误")
		return
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		shared.BadRequest(c, "参数错误", "月份格式错误")
		return
	}

	report, err := api.BookStatisticsService.GenerateMonthlyReport(c.Request.Context(), year, month)
	if err != nil {
		c.Error(err)
		return
	}

	shared.SuccessData(c, report)
}

// SearchStatistics 搜索统计信息
// @Summary 搜索统计信息
// @Description 根据关键词搜索统计信息
// @Tags 图书统计
// @Accept json
// @Produce json
// @Param keyword query string true "搜索关键词"
// @Param page query int false "页码" default(1)
// @Param limit query int false "每页数量" default(10)
// @Success 200 {object} APIResponse "成功"
// @Failure 400 {object} APIResponse "请求参数错误"
// @Failure 500 {object} APIResponse "服务器内部错误"
// @Router /api/v1/reading/statistics/search [get]
func (api *BookStatisticsAPI) SearchStatistics(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		shared.BadRequest(c, "参数错误", "搜索关键词不能为空")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	statistics, total, err := api.BookStatisticsService.SearchStatistics(c.Request.Context(), keyword, page, limit)
	if err != nil {
		c.Error(err)
		return
	}

	shared.Paginated(c, statistics, total, page, limit, "搜索成功")
}
