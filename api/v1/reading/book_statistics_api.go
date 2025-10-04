package reading

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/models/reading/bookstore"
	"Qingyu_backend/service/bookstore"
)

// BookStatisticsAPI 图书统计API控制器
type BookStatisticsAPI struct {
	BookStatisticsService bookstore.BookStatisticsService
}

// NewBookStatisticsAPI 创建新的图书统计API实例
func NewBookStatisticsAPI(bookStatisticsService bookstore.BookStatisticsService) *BookStatisticsAPI {
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
// @Success 200 {object} APIResponse{data=bookstore.BookStatistics} "成功"
// @Failure 400 {object} APIResponse "请求参数错误"
// @Failure 404 {object} APIResponse "统计信息不存在"
// @Failure 500 {object} APIResponse "服务器内部错误"
// @Router /api/v1/reading/books/{book_id}/statistics [get]
func (api *BookStatisticsAPI) GetBookStatistics(c *gin.Context) {
	bookIDStr := c.Param("book_id")
	if bookIDStr == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "图书ID不能为空",
			Data:    nil,
		})
		return
	}

	bookID, err := primitive.ObjectIDFromHex(bookIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "无效的图书ID格式",
			Data:    nil,
		})
		return
	}

	statistics, err := api.BookStatisticsService.GetByBookID(c.Request.Context(), bookID)
	if err != nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Code:    404,
			Message: "统计信息不存在",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "获取成功",
		Data:    statistics,
	})
}

// GetTopViewedBooks 获取最多浏览的图书
// @Summary 获取最多浏览的图书
// @Description 获取浏览量最高的图书列表
// @Tags 图书统计
// @Accept json
// @Produce json
// @Param limit query int false "返回数量" default(10)
// @Success 200 {object} APIResponse{data=[]bookstore.BookStatistics} "成功"
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
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "获取最多浏览图书失败",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "获取成功",
		Data:    statistics,
	})
}

// GetTopFavoritedBooks 获取最多收藏的图书
// @Summary 获取最多收藏的图书
// @Description 获取收藏量最高的图书列表
// @Tags 图书统计
// @Accept json
// @Produce json
// @Param limit query int false "返回数量" default(10)
// @Success 200 {object} APIResponse{data=[]bookstore.BookStatistics} "成功"
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
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "获取最多收藏图书失败",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "获取成功",
		Data:    statistics,
	})
}

// GetTopRatedBooks 获取最高评分的图书
// @Summary 获取最高评分的图书
// @Description 获取评分最高的图书列表
// @Tags 图书统计
// @Accept json
// @Produce json
// @Param limit query int false "返回数量" default(10)
// @Success 200 {object} APIResponse{data=[]bookstore.BookStatistics} "成功"
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
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "获取最高评分图书失败",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "获取成功",
		Data:    statistics,
	})
}

// GetHottestBooks 获取最热门的图书
// @Summary 获取最热门的图书
// @Description 获取热度最高的图书列表
// @Tags 图书统计
// @Accept json
// @Produce json
// @Param limit query int false "返回数量" default(10)
// @Success 200 {object} APIResponse{data=[]bookstore.BookStatistics} "成功"
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
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "获取最热门图书失败",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "获取成功",
		Data:    statistics,
	})
}

// GetTrendingBooks 获取趋势图书
// @Summary 获取趋势图书
// @Description 获取当前趋势图书列表
// @Tags 图书统计
// @Accept json
// @Produce json
// @Param limit query int false "返回数量" default(10)
// @Success 200 {object} APIResponse{data=[]bookstore.BookStatistics} "成功"
// @Failure 400 {object} APIResponse "请求参数错误"
// @Failure 500 {object} APIResponse "服务器内部错误"
// @Router /api/v1/reading/statistics/trending [get]
func (api *BookStatisticsAPI) GetTrendingBooks(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit < 1 || limit > 100 {
		limit = 10
	}

	statistics, err := api.BookStatisticsService.GetTrendingBooks(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "获取趋势图书失败",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "获取成功",
		Data:    statistics,
	})
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
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "图书ID不能为空",
			Data:    nil,
		})
		return
	}

	bookID, err := primitive.ObjectIDFromHex(bookIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "无效的图书ID格式",
			Data:    nil,
		})
		return
	}

	err = api.BookStatisticsService.IncrementViewCount(c.Request.Context(), bookID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "增加浏览量失败",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "增加浏览量成功",
		Data:    nil,
	})
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
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "图书ID不能为空",
			Data:    nil,
		})
		return
	}

	bookID, err := primitive.ObjectIDFromHex(bookIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "无效的图书ID格式",
			Data:    nil,
		})
		return
	}

	err = api.BookStatisticsService.IncrementFavoriteCount(c.Request.Context(), bookID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "增加收藏量失败",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "增加收藏量成功",
		Data:    nil,
	})
}

// GetAggregatedStatistics 获取聚合统计信息
// @Summary 获取聚合统计信息
// @Description 获取系统整体的聚合统计信息
// @Tags 图书统计
// @Accept json
// @Produce json
// @Success 200 {object} APIResponse{data=map[string]interface{}} "成功"
// @Failure 500 {object} APIResponse "服务器内部错误"
// @Router /api/v1/reading/statistics/aggregated [get]
func (api *BookStatisticsAPI) GetAggregatedStatistics(c *gin.Context) {
	statistics, err := api.BookStatisticsService.GetAggregatedStatistics(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "获取聚合统计信息失败",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "获取成功",
		Data:    statistics,
	})
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
// @Success 200 {object} PaginatedResponse{data=[]bookstore.BookStatistics} "成功"
// @Failure 400 {object} APIResponse "请求参数错误"
// @Failure 500 {object} APIResponse "服务器内部错误"
// @Router /api/v1/reading/statistics/time-range [get]
func (api *BookStatisticsAPI) GetStatisticsByTimeRange(c *gin.Context) {
	startTimeStr := c.Query("start_time")
	endTimeStr := c.Query("end_time")

	if startTimeStr == "" || endTimeStr == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "开始时间和结束时间不能为空",
			Data:    nil,
		})
		return
	}

	startTime, err := time.Parse(time.RFC3339, startTimeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "开始时间格式错误",
			Data:    nil,
		})
		return
	}

	endTime, err := time.Parse(time.RFC3339, endTimeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "结束时间格式错误",
			Data:    nil,
		})
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

	statistics, err := api.BookStatisticsService.GetStatisticsByTimeRange(c.Request.Context(), startTime, endTime, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "获取时间范围统计信息失败",
			Data:    nil,
		})
		return
	}

	total := int64(len(statistics))

	c.JSON(http.StatusOK, PaginatedResponse{
		Code:    200,
		Message: "获取成功",
		Data:    statistics,
		Total:   total,
		Page:    page,
		Limit:   limit,
	})
}

// GetDailyStatisticsReport 获取日统计报告
// @Summary 获取日统计报告
// @Description 获取指定日期的统计报告
// @Tags 图书统计
// @Accept json
// @Produce json
// @Param date query string true "日期" format(date)
// @Success 200 {object} APIResponse{data=map[string]interface{}} "成功"
// @Failure 400 {object} APIResponse "请求参数错误"
// @Failure 500 {object} APIResponse "服务器内部错误"
// @Router /api/v1/reading/statistics/daily-report [get]
func (api *BookStatisticsAPI) GetDailyStatisticsReport(c *gin.Context) {
	dateStr := c.Query("date")
	if dateStr == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "日期不能为空",
			Data:    nil,
		})
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "日期格式错误",
			Data:    nil,
		})
		return
	}

	report, err := api.BookStatisticsService.GetDailyStatisticsReport(c.Request.Context(), date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "获取日统计报告失败",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "获取成功",
		Data:    report,
	})
}

// GetWeeklyStatisticsReport 获取周统计报告
// @Summary 获取周统计报告
// @Description 获取指定周的统计报告
// @Tags 图书统计
// @Accept json
// @Produce json
// @Param year query int true "年份"
// @Param week query int true "周数"
// @Success 200 {object} APIResponse{data=map[string]interface{}} "成功"
// @Failure 400 {object} APIResponse "请求参数错误"
// @Failure 500 {object} APIResponse "服务器内部错误"
// @Router /api/v1/reading/statistics/weekly-report [get]
func (api *BookStatisticsAPI) GetWeeklyStatisticsReport(c *gin.Context) {
	yearStr := c.Query("year")
	weekStr := c.Query("week")

	if yearStr == "" || weekStr == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "年份和周数不能为空",
			Data:    nil,
		})
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "年份格式错误",
			Data:    nil,
		})
		return
	}

	week, err := strconv.Atoi(weekStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "周数格式错误",
			Data:    nil,
		})
		return
	}

	report, err := api.BookStatisticsService.GetWeeklyStatisticsReport(c.Request.Context(), year, week)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "获取周统计报告失败",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "获取成功",
		Data:    report,
	})
}

// GetMonthlyStatisticsReport 获取月统计报告
// @Summary 获取月统计报告
// @Description 获取指定月份的统计报告
// @Tags 图书统计
// @Accept json
// @Produce json
// @Param year query int true "年份"
// @Param month query int true "月份"
// @Success 200 {object} APIResponse{data=map[string]interface{}} "成功"
// @Failure 400 {object} APIResponse "请求参数错误"
// @Failure 500 {object} APIResponse "服务器内部错误"
// @Router /api/v1/reading/statistics/monthly-report [get]
func (api *BookStatisticsAPI) GetMonthlyStatisticsReport(c *gin.Context) {
	yearStr := c.Query("year")
	monthStr := c.Query("month")

	if yearStr == "" || monthStr == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "年份和月份不能为空",
			Data:    nil,
		})
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "年份格式错误",
			Data:    nil,
		})
		return
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "月份格式错误",
			Data:    nil,
		})
		return
	}

	report, err := api.BookStatisticsService.GetMonthlyStatisticsReport(c.Request.Context(), year, month)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "获取月统计报告失败",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "获取成功",
		Data:    report,
	})
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
// @Success 200 {object} PaginatedResponse{data=[]bookstore.BookStatistics} "成功"
// @Failure 400 {object} APIResponse "请求参数错误"
// @Failure 500 {object} APIResponse "服务器内部错误"
// @Router /api/v1/reading/statistics/search [get]
func (api *BookStatisticsAPI) SearchStatistics(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "搜索关键词不能为空",
			Data:    nil,
		})
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

	statistics, err := api.BookStatisticsService.SearchByKeyword(c.Request.Context(), keyword, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "搜索统计信息失败",
			Data:    nil,
		})
		return
	}

	total := int64(len(statistics))

	c.JSON(http.StatusOK, PaginatedResponse{
		Code:    200,
		Message: "搜索成功",
		Data:    statistics,
		Total:   total,
		Page:    page,
		Limit:   limit,
	})
}