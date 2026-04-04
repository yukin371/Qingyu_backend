package reader

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/pkg/response"
	readerservice "Qingyu_backend/service/reader"
	socialservice "Qingyu_backend/service/social"
)

const (
	readerStatisticsDefaultDays = 30
	readerStatisticsHoursPerDay = 24
	readerStatisticsDaysPerWeek = 7
)

// ReaderStatisticsAPI 提供读者聚合统计接口。
type ReaderStatisticsAPI struct {
	readerService         *readerservice.ReaderService
	bookmarkService       readerservice.BookmarkService
	likeService           *socialservice.LikeService
	collectionService     *socialservice.CollectionService
	readingHistoryService *readerservice.ReadingHistoryService
}

type ReaderStatisticsOverviewCore struct {
	TotalReadingTime int64 `json:"totalReadingTime"`
	UnfinishedCount  int   `json:"unfinishedCount"`
	FinishedCount    int   `json:"finishedCount"`
}

type ReaderStatisticsBookmarkSummary struct {
	Total        int64 `json:"total"`
	PublicCount  int64 `json:"publicCount"`
	PrivateCount int64 `json:"privateCount"`
	ThisWeek     int64 `json:"thisWeek"`
	ThisMonth    int64 `json:"thisMonth"`
}

type ReaderStatisticsHistorySummary struct {
	TotalBooks       int       `json:"totalBooks"`
	TotalChapters    int       `json:"totalChapters"`
	TotalDuration    int       `json:"totalDuration"`
	AvgDailyDuration int       `json:"avgDailyDuration"`
	LastReadTime     time.Time `json:"lastReadTime"`
}

type ReaderStatisticsOverviewResponse struct {
	TotalReadingTime int                             `json:"totalReadingTime"`
	UnfinishedCount  int                             `json:"unfinishedCount"`
	FinishedCount    int                             `json:"finishedCount"`
	Bookmarks        ReaderStatisticsBookmarkSummary `json:"bookmarks"`
	Likes            map[string]interface{}          `json:"likes"`
	Collections      map[string]interface{}          `json:"collections"`
	History          ReaderStatisticsHistorySummary  `json:"history"`
	Overview         ReaderStatisticsOverviewCore    `json:"overview"`
}

type ReaderStatisticsReadingTimeResponse struct {
	Period           string `json:"period"`
	TotalReadingTime int64  `json:"totalReadingTime"`
	Minutes          int64  `json:"minutes"`
}

type ReaderStatisticsTrendItem struct {
	Date        string `json:"date"`
	ReadingTime int    `json:"readingTime"`
	Minutes     int    `json:"minutes,omitempty"`
	Books       int    `json:"books"`
}

type ReaderStatisticsTrendResponse struct {
	Days             int                         `json:"days"`
	TotalReadingTime int64                       `json:"totalReadingTime,omitempty"`
	Items            []ReaderStatisticsTrendItem `json:"items"`
	List             []ReaderStatisticsTrendItem `json:"list"`
}

// NewReaderStatisticsAPI 创建读者聚合统计 API。
func NewReaderStatisticsAPI(
	readerService *readerservice.ReaderService,
	bookmarkService readerservice.BookmarkService,
	likeService *socialservice.LikeService,
	collectionService *socialservice.CollectionService,
	readingHistoryService *readerservice.ReadingHistoryService,
) *ReaderStatisticsAPI {
	return &ReaderStatisticsAPI{
		readerService:         readerService,
		bookmarkService:       bookmarkService,
		likeService:           likeService,
		collectionService:     collectionService,
		readingHistoryService: readingHistoryService,
	}
}

func (api *ReaderStatisticsAPI) getUserID(c *gin.Context) (string, bool) {
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

	return userID, true
}

func (api *ReaderStatisticsAPI) resolvePeriodRange(period string) (time.Time, time.Time, bool) {
	now := time.Now()
	switch period {
	case "today":
		start := now.Truncate(readerStatisticsHoursPerDay * time.Hour)
		return start, start.Add(readerStatisticsHoursPerDay * time.Hour), true
	case "week":
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = readerStatisticsDaysPerWeek
		}
		start := now.AddDate(0, 0, -(weekday - 1)).Truncate(readerStatisticsHoursPerDay * time.Hour)
		return start, start.AddDate(0, 0, readerStatisticsDaysPerWeek), true
	case "month":
		start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		return start, start.AddDate(0, 1, 0), true
	default:
		return time.Time{}, time.Time{}, false
	}
}

// GetOverview 获取读者统计概览。
// @Summary 获取读者统计概览
// @Description 获取当前登录读者的聚合统计概览
// @Tags Reader-Statistics
// @Accept json
// @Produce json
// @Success 200 {object} response.APIResponse{data=ReaderStatisticsOverviewResponse}
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/reader/statistics [get]
// @Router /api/v1/reader/statistics/overview [get]
// GetOverview 获取读者统计概览。
func (api *ReaderStatisticsAPI) GetOverview(c *gin.Context) {
	userID, ok := api.getUserID(c)
	if !ok {
		return
	}

	totalReadingTime, err := api.readerService.GetTotalReadingTime(c.Request.Context(), userID)
	if err != nil {
		c.Error(err)
		return
	}

	unfinishedBooks, err := api.readerService.GetUnfinishedBooks(c.Request.Context(), userID)
	if err != nil {
		c.Error(err)
		return
	}

	finishedBooks, err := api.readerService.GetFinishedBooks(c.Request.Context(), userID)
	if err != nil {
		c.Error(err)
		return
	}

	bookmarks := gin.H{}
	if api.bookmarkService != nil {
		bookmarkStats, err := api.bookmarkService.GetBookmarkStats(c.Request.Context(), userID)
		if err != nil {
			c.Error(err)
			return
		}
		bookmarks = gin.H{
			"total":        bookmarkStats.TotalCount,
			"publicCount":  bookmarkStats.PublicCount,
			"privateCount": bookmarkStats.PrivateCount,
			"thisWeek":     bookmarkStats.ThisWeekCount,
			"thisMonth":    bookmarkStats.ThisMonthCount,
		}
	}

	likes := map[string]interface{}{}
	if api.likeService != nil {
		likeStats, err := api.likeService.GetUserLikeStats(c.Request.Context(), userID)
		if err != nil {
			c.Error(err)
			return
		}
		likes = likeStats
	}

	collections := map[string]interface{}{}
	if api.collectionService != nil {
		collectionStats, err := api.collectionService.GetUserCollectionStats(c.Request.Context(), userID)
		if err != nil {
			c.Error(err)
			return
		}
		collections = collectionStats
	}

	history := gin.H{}
	if api.readingHistoryService != nil {
		historyStats, err := api.readingHistoryService.GetUserReadingStats(c.Request.Context(), userID)
		if err != nil {
			c.Error(err)
			return
		}
		if historyStats != nil {
			history = gin.H{
				"totalBooks":       historyStats.TotalBooks,
				"totalChapters":    historyStats.TotalChapters,
				"totalDuration":    historyStats.TotalDuration,
				"avgDailyDuration": historyStats.AvgDailyDuration,
				"lastReadTime":     historyStats.LastReadTime,
			}
		}
	}

	response.Success(c, gin.H{
		"totalReadingTime": totalReadingTime,
		"unfinishedCount":  len(unfinishedBooks),
		"finishedCount":    len(finishedBooks),
		"bookmarks":        bookmarks,
		"likes":            likes,
		"collections":      collections,
		"history":          history,
		"overview": gin.H{
			"totalReadingTime": totalReadingTime,
			"unfinishedCount":  len(unfinishedBooks),
			"finishedCount":    len(finishedBooks),
		},
	})
}

// GetReadingTime 获取阅读时长统计。
// @Summary 获取阅读时长统计
// @Description 按 today/week/month/all 返回当前读者的阅读时长
// @Tags Reader-Statistics
// @Accept json
// @Produce json
// @Param period query string false "统计周期" Enums(today,week,month,all) default(all)
// @Success 200 {object} response.APIResponse{data=ReaderStatisticsReadingTimeResponse}
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/reader/statistics/reading-time [get]
// GetReadingTime 获取阅读时长统计。
func (api *ReaderStatisticsAPI) GetReadingTime(c *gin.Context) {
	userID, ok := api.getUserID(c)
	if !ok {
		return
	}

	period := c.DefaultQuery("period", "all")
	var totalTime int64
	var err error

	start, end, hasRange := api.resolvePeriodRange(period)
	if hasRange {
		totalTime, err = api.readerService.GetReadingTimeByPeriod(c.Request.Context(), userID, start, end)
	} else {
		totalTime, err = api.readerService.GetTotalReadingTime(c.Request.Context(), userID)
	}
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, gin.H{
		"period":           period,
		"totalReadingTime": totalTime,
		"minutes":          totalTime,
	})
}

// GetHeatmap 获取阅读热力图数据。
// @Summary 获取阅读热力图
// @Description 获取当前读者最近 N 天的阅读热力图数据
// @Tags Reader-Statistics
// @Accept json
// @Produce json
// @Param days query int false "天数" default(30)
// @Success 200 {object} response.APIResponse{data=ReaderStatisticsTrendResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/reader/statistics/heatmap [get]
// GetHeatmap 获取阅读热力图数据。
func (api *ReaderStatisticsAPI) GetHeatmap(c *gin.Context) {
	userID, ok := api.getUserID(c)
	if !ok {
		return
	}
	if api.readingHistoryService == nil {
		response.Success(c, gin.H{
			"days":  readerStatisticsDefaultDays,
			"items": []gin.H{},
			"list":  []gin.H{},
		})
		return
	}

	days, err := strconv.Atoi(c.DefaultQuery("days", strconv.Itoa(readerStatisticsDefaultDays)))
	if err != nil || days < 1 || days > 365 {
		response.BadRequest(c, "参数错误", "天数必须在1-365之间")
		return
	}

	dailyStats, err := api.readingHistoryService.GetUserDailyReadingStats(c.Request.Context(), userID, days)
	if err != nil {
		c.Error(err)
		return
	}

	items := make([]gin.H, 0, len(dailyStats))
	for _, item := range dailyStats {
		items = append(items, gin.H{
			"date":        item.Date,
			"readingTime": item.Duration,
			"minutes":     item.Duration,
			"books":       item.Books,
		})
	}

	response.Success(c, gin.H{
		"days":  days,
		"items": items,
		"list":  items,
	})
}

// GetTrends 获取阅读趋势数据。
// @Summary 获取阅读趋势
// @Description 获取当前读者最近 N 天的阅读趋势数据
// @Tags Reader-Statistics
// @Accept json
// @Produce json
// @Param days query int false "天数" default(30)
// @Success 200 {object} response.APIResponse{data=ReaderStatisticsTrendResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/reader/statistics/trends [get]
// GetTrends 获取阅读趋势数据。
func (api *ReaderStatisticsAPI) GetTrends(c *gin.Context) {
	userID, ok := api.getUserID(c)
	if !ok {
		return
	}
	if api.readingHistoryService == nil {
		response.Success(c, gin.H{
			"days":  readerStatisticsDefaultDays,
			"items": []gin.H{},
			"list":  []gin.H{},
		})
		return
	}

	days, err := strconv.Atoi(c.DefaultQuery("days", strconv.Itoa(readerStatisticsDefaultDays)))
	if err != nil || days < 1 || days > 365 {
		response.BadRequest(c, "参数错误", "天数必须在1-365之间")
		return
	}

	dailyStats, err := api.readingHistoryService.GetUserDailyReadingStats(c.Request.Context(), userID, days)
	if err != nil {
		c.Error(err)
		return
	}

	items := make([]gin.H, 0, len(dailyStats))
	var totalTime int64
	for _, item := range dailyStats {
		totalTime += int64(item.Duration)
		items = append(items, gin.H{
			"date":        item.Date,
			"readingTime": item.Duration,
			"books":       item.Books,
		})
	}

	response.Success(c, gin.H{
		"days":             days,
		"totalReadingTime": totalTime,
		"items":            items,
		"list":             items,
	})
}
