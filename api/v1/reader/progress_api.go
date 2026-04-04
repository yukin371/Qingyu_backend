package reader

import (
	"context"
	"strconv"
	"strings"
	"time"

	"Qingyu_backend/api/v1/shared"
	"Qingyu_backend/models"
	readerModels "Qingyu_backend/models/reader"
	"Qingyu_backend/pkg/response"
	readerRepo "Qingyu_backend/repository/interfaces/reader"
	"Qingyu_backend/service/interfaces"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	hoursPerDay = 24
	daysPerWeek = 7
)

// ProgressAPI 阅读进度API
type ProgressAPI struct {
	readerService interfaces.ReaderService
	deviceRepo   readerRepo.DeviceRepository
}

// NewProgressAPI 创建阅读进度API实例
func NewProgressAPI(readerService interfaces.ReaderService, deviceRepo readerRepo.DeviceRepository) *ProgressAPI {
	return &ProgressAPI{
		readerService: readerService,
		deviceRepo:   deviceRepo,
	}
}

// SaveProgressRequest 保存进度请求
type SaveProgressRequest struct {
	BookID    string  `json:"bookId" binding:"required"`
	ChapterID string  `json:"chapterId" binding:"required"`
	Progress  float64 `json:"progress" binding:"min=0,max=1"`
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
//	@Success	200		{object}	response.APIResponse
//	@Router		/api/v1/reader/progress/{bookId} [get]
func (api *ProgressAPI) GetReadingProgress(c *gin.Context) {
	bookID := c.Param("bookId")

	// 获取用户ID
	userID, ok := shared.GetUserID(c)
	if !ok {
		return
	}

	progress, err := api.readerService.GetReadingProgress(c.Request.Context(), userID, bookID)
	if err != nil {
		c.Error(err)
		return
	}

	// 转换为 DTO
	progressDTO := ToReadingProgressDTO(progress)
	response.Success(c, progressDTO)
}

// SaveReadingProgress 保存阅读进度
//
//	@Summary	保存阅读进度
//	@Tags		阅读器
//	@Param		request	body		SaveProgressRequest	true	"保存进度请求"
//	@Success	200		{object}	response.APIResponse
//	@Router		/api/v1/reader/progress [post]
func (api *ProgressAPI) SaveReadingProgress(c *gin.Context) {
	var req SaveProgressRequest
	if !shared.BindAndValidate(c, &req) {
		return
	}

	// 获取用户ID
	userID, ok := shared.GetUserID(c)
	if !ok {
		return
	}

	err := api.readerService.SaveReadingProgress(c.Request.Context(), userID, req.BookID, req.ChapterID, req.Progress)
	if err != nil {
		c.Error(err)
		return
	}

	// 异步追踪设备信息
	go api.trackDevice(c, userID)

	response.Success(c, nil)
}

// UpdateReadingTime 更新阅读时长
//
//	@Summary	更新阅读时长
//	@Tags		阅读器
//	@Param		request	body		UpdateReadingTimeRequest	true	"更新时长请求"
//	@Success	200		{object}	response.APIResponse
//	@Router		/api/v1/reader/progress/reading-time [put]
func (api *ProgressAPI) UpdateReadingTime(c *gin.Context) {
	var req UpdateReadingTimeRequest
	if !shared.BindAndValidate(c, &req) {
		return
	}

	// 获取用户ID
	userID, ok := shared.GetUserID(c)
	if !ok {
		return
	}

	err := api.readerService.UpdateReadingTime(c.Request.Context(), userID, req.BookID, req.Duration)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, nil)
}

// GetRecentReading 获取最近阅读记录
//
//	@Summary	获取最近阅读记录
//	@Tags		阅读器
//	@Param		limit	query		int	false	"数量限制"	default(20)
//	@Success	200		{object}	response.APIResponse
//	@Router		/api/v1/reader/progress/recent [get]
func (api *ProgressAPI) GetRecentReading(c *gin.Context) {
	// 获取用户ID
	userID, ok := shared.GetUserID(c)
	if !ok {
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	progresses, err := api.readerService.GetRecentReading(c.Request.Context(), userID, limit)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, progresses)
}

// GetReadingHistory 获取阅读历史
//
//	@Summary	获取阅读历史
//	@Tags		阅读器
//	@Param		page	query		int	false	"页码"	default(1)
//	@Param		size	query		int	false	"每页数量"	default(20)
//	@Success	200		{object}	response.APIResponse
//	@Router		/api/v1/reader/progress/history [get]
func (api *ProgressAPI) GetReadingHistory(c *gin.Context) {
	// 获取用户ID
	userID, ok := shared.GetUserID(c)
	if !ok {
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	progresses, total, err := api.readerService.GetReadingHistory(c.Request.Context(), userID, page, size)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, gin.H{
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
//	@Success	200		{object}	response.APIResponse
//	@Router		/api/v1/reader/progress/stats [get]
func (api *ProgressAPI) GetReadingStats(c *gin.Context) {
	// 获取用户ID
	userID, ok := shared.GetUserID(c)
	if !ok {
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
		totalTime, err = api.readerService.GetReadingTimeByPeriod(c.Request.Context(), userID, start, end)
	case "week":
		// 本周
		now := time.Now()
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = daysPerWeek
		}
		start := now.AddDate(0, 0, -(weekday - 1)).Truncate(hoursPerDay * time.Hour)
		end := start.AddDate(0, 0, daysPerWeek)
		totalTime, err = api.readerService.GetReadingTimeByPeriod(c.Request.Context(), userID, start, end)
	case "month":
		// 本月
		now := time.Now()
		start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		end := start.AddDate(0, 1, 0)
		totalTime, err = api.readerService.GetReadingTimeByPeriod(c.Request.Context(), userID, start, end)
	default:
		// 总计
		totalTime, err = api.readerService.GetTotalReadingTime(c.Request.Context(), userID)
	}

	if err != nil {
		c.Error(err)
		return
	}

	// 获取未读完和已读完的书籍
	unfinished, errUnfinished := api.readerService.GetUnfinishedBooks(c.Request.Context(), userID)
	if errUnfinished != nil {
		unfinished = []*readerModels.ReadingProgress{} // 返回空列表而非失败
	}

	finished, errFinished := api.readerService.GetFinishedBooks(c.Request.Context(), userID)
	if errFinished != nil {
		finished = []*readerModels.ReadingProgress{} // 返回空列表而非失败
	}

	response.Success(c, gin.H{
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
//	@Success	200	{object}	response.APIResponse
//	@Router		/api/v1/reader/progress/unfinished [get]
func (api *ProgressAPI) GetUnfinishedBooks(c *gin.Context) {
	// 获取用户ID
	userID, ok := shared.GetUserID(c)
	if !ok {
		return
	}

	progresses, err := api.readerService.GetUnfinishedBooks(c.Request.Context(), userID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, progresses)
}

// GetFinishedBooks 获取已读完的书籍
//
//	@Summary	获取已读完的书籍
//	@Tags		阅读器
//	@Success	200	{object}	response.APIResponse
//	@Router		/api/v1/reader/progress/finished [get]
func (api *ProgressAPI) GetFinishedBooks(c *gin.Context) {
	// 获取用户ID
	userID, ok := shared.GetUserID(c)
	if !ok {
		return
	}

	progresses, err := api.readerService.GetFinishedBooks(c.Request.Context(), userID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, progresses)
}

// DevicesProgressResponse 跨设备阅读记录响应
type DevicesProgressResponse struct {
	Devices []DeviceProgress `json:"devices"`
	Count   int              `json:"count"`
}

// DeviceProgress 设备阅读进度
type DeviceProgress struct {
	DeviceID    string    `json:"deviceId"`
	DeviceName  string    `json:"deviceName"`
	DeviceType  string    `json:"deviceType"`
	LastSyncAt  time.Time `json:"lastSyncAt"`
	CurrentBook string    `json:"currentBook,omitempty"`
	Progress    float64   `json:"progress,omitempty"`
}

// GetDevicesProgress 获取跨设备阅读记录
//
//	@Summary	获取跨设备阅读记录
//	@Tags		阅读器
//	@Success	200	{object}	response.APIResponse
//	@Router		/api/v1/reader/progress/devices [get]
func (api *ProgressAPI) GetDevicesProgress(c *gin.Context) {
	// 获取用户ID
	userID, ok := shared.GetUserID(c)
	if !ok {
		return
	}

	// 从 deviceRepo 查询用户设备列表
	if api.deviceRepo == nil {
		response.Success(c, DevicesProgressResponse{
			Devices: []DeviceProgress{},
			Count:   0,
		})
		return
	}

	ctx := c.Request.Context()
	devices, err := api.deviceRepo.GetByUserID(ctx, userID)
	if err != nil {
		c.Error(err)
		response.Success(c, DevicesProgressResponse{
			Devices: []DeviceProgress{},
			Count:   0,
		})
		return
	}

	result := make([]DeviceProgress, 0, len(devices))
	for _, d := range devices {
		result = append(result, DeviceProgress{
			DeviceID:   d.ID.Hex(),
			DeviceName: d.Name,
			DeviceType: d.Type,
			LastSyncAt: d.LastSeen,
		})
	}

	response.Success(c, DevicesProgressResponse{
		Devices: result,
		Count:   len(result),
	})
}

// trackDevice 异步追踪用户设备信息（通过 User-Agent 解析）
func (api *ProgressAPI) trackDevice(c *gin.Context, userID string) {
	if api.deviceRepo == nil {
		return
	}

	ua := c.GetHeader("User-Agent")
	if ua == "" {
		return
	}

	name, deviceType := detectDevice(ua)
	ip := c.ClientIP()

	device := &models.Device{
		UserID:    primitive.ObjectID{},
		Name:      name,
		Type:      deviceType,
		UserAgent: ua,
		IP:        ip,
	}
	if uid, err := primitive.ObjectIDFromHex(userID); err == nil {
		device.UserID = uid
	}

	_ = api.deviceRepo.UpsertDevice(context.Background(), device)
}

// detectDevice 从 User-Agent 解析设备名称和类型
func detectDevice(ua string) (name string, deviceType string) {
	deviceType = "desktop"
	name = "Unknown Device"

	uaLower := strings.ToLower(ua)

	if strings.Contains(uaLower, "mobile") || strings.Contains(uaLower, "android") && strings.Contains(uaLower, "mobile") {
		deviceType = "mobile"
	} else if strings.Contains(uaLower, "tablet") || strings.Contains(uaLower, "ipad") {
		deviceType = "tablet"
	}

	// 检测浏览器
	browser := ""
	if strings.Contains(uaLower, "chrome") && !strings.Contains(uaLower, "edg") {
		browser = "Chrome"
	} else if strings.Contains(uaLower, "safari") && !strings.Contains(uaLower, "chrome") {
		browser = "Safari"
	} else if strings.Contains(uaLower, "firefox") {
		browser = "Firefox"
	} else if strings.Contains(uaLower, "edg") {
		browser = "Edge"
	}

	// 检测操作系统
	os := ""
	if strings.Contains(uaLower, "windows") {
		os = "Windows"
	} else if strings.Contains(uaLower, "mac os") || strings.Contains(uaLower, "macos") {
		os = "macOS"
	} else if strings.Contains(uaLower, "linux") && !strings.Contains(uaLower, "android") {
		os = "Linux"
	} else if strings.Contains(uaLower, "android") {
		os = "Android"
	} else if strings.Contains(uaLower, "iphone") || strings.Contains(uaLower, "ipad") {
		os = "iOS"
	}

	if browser != "" && os != "" {
		name = browser + " on " + os
	} else if os != "" {
		name = os + " Device"
	}

	return name, deviceType
}
