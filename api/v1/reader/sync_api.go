package reader

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	gorilla_websocket "github.com/gorilla/websocket"

	"Qingyu_backend/api/v1/shared"
	progressSync "Qingyu_backend/pkg/sync"
	ws "Qingyu_backend/pkg/websocket"
	"Qingyu_backend/service/interfaces"
)

// SyncAPI 阅读进度同步API
type SyncAPI struct {
	syncService interfaces.ProgressSyncService
}

// NewSyncAPI 创建同步API实例
func NewSyncAPI(syncService interfaces.ProgressSyncService) *SyncAPI {
	return &SyncAPI{
		syncService: syncService,
	}
}

// WebSocketUpgrader WebSocket升级器
var WebSocketUpgrader = gorilla_websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // 生产环境应该验证Origin
	},
}

// SyncWebSocket WebSocket同步连接
//
//	@Summary		WebSocket进度同步
//	@Description	建立WebSocket连接进行实时进度同步
//	@Tags			Reader-Sync
//	@Accept			json
//	@Produce		json
//	@Param			token		header	string	true	"JWT Token"
//	@Success		101			{string}	string	"Switching Protocols"
//	@Failure		400			{object}	shared.APIResponse
//	@Failure		401			{object}	shared.APIResponse
//	@Router			/api/v1/reader/progress/ws [get]
func (api *SyncAPI) SyncWebSocket(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "需要登录")
		return
	}

	userIDStr := userID.(string)

	// 获取设备ID
	deviceID := c.GetHeader("X-Device-ID")
	if deviceID == "" {
		deviceID = "unknown"
	}

	// 升级到WebSocket
	conn, err := WebSocketUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		shared.Error(c, http.StatusBadRequest, "WebSocket升级失败", err.Error())
		return
	}

	// 创建客户端
	hub := api.syncService.GetHub()
	client := &ws.Client{
		UserID:   userIDStr,
		DeviceID: deviceID,
		Hub:      hub,
		Conn:     conn,
		Send:     make(chan *ws.ProgressMessage, 256),
	}

	// 注册客户端
	hub.Register <- client

	// 启动读写协程
	go client.WritePump()
	go client.ReadPump()
}

// SyncProgress 同步进度
//
//	@Summary		同步阅读进度
//	@Description	同步阅读进度到服务器并推送到其他设备
//	@Tags			Reader-Sync
//	@Accept			json
//	@Produce		json
//	@Param			request	body		SyncProgressRequest	true	"同步请求"
//	@Success		200		{object}	shared.APIResponse
//	@Failure		400		{object}	shared.APIResponse
//	@Failure		401		{object}	shared.APIResponse
//	@Router			/api/v1/reader/progress/sync [post]
func (api *SyncAPI) SyncProgress(c *gin.Context) {
	var req SyncProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	// 获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "需要登录")
		return
	}

	userIDStr := userID.(string)

	// 同步进度
	if err := api.syncService.SyncProgress(c.Request.Context(), userIDStr, req.BookID, req.ChapterID, req.DeviceID, req.Progress); err != nil {
		shared.Error(c, http.StatusInternalServerError, "同步失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "同步成功", nil)
}

// MergeOfflineProgresses 合并离线进度
//
//	@Summary		合并离线进度
//	@Description	批量合并离线阅读进度
//	@Tags			Reader-Sync
//	@Accept			json
//	@Produce		json
//	@Param			request	body		MergeProgressRequest	true	"合并请求"
//	@Success		200		{object}	shared.APIResponse
//	@Failure		400		{object}	shared.APIResponse
//	@Failure		401		{object}	shared.APIResponse
//	@Router			/api/v1/reader/progress/merge [post]
func (api *SyncAPI) MergeOfflineProgresses(c *gin.Context) {
	var req MergeProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	// 获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "需要登录")
		return
	}

	userIDStr := userID.(string)

	// 转换进度
	progresses := make([]progressSync.OfflineProgress, len(req.Progresses))
	for i, p := range req.Progresses {
		// 解析时间戳
		timestamp, err := time.Parse(time.RFC3339, p.Timestamp)
		if err != nil {
			shared.Error(c, http.StatusBadRequest, "时间戳格式错误", err.Error())
			return
		}

		progresses[i] = progressSync.OfflineProgress{
			UserID:    userIDStr,
			BookID:    p.BookID,
			ChapterID: p.ChapterID,
			Progress:  p.Progress,
			Timestamp: timestamp,
			DeviceID:  p.DeviceID,
		}
	}

	// 合并进度
	if err := api.syncService.MergeOfflineProgresses(c.Request.Context(), userIDStr, progresses); err != nil {
		shared.Error(c, http.StatusInternalServerError, "合并失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "合并成功", nil)
}

// GetSyncStatus 获取同步状态
//
//	@Summary		获取同步状态
//	@Description	获取当前用户的同步状态
//	@Tags			Reader-Sync
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	shared.APIResponse
//	@Failure		401	{object}	shared.APIResponse
//	@Router			/api/v1/reader/progress/sync-status [get]
func (api *SyncAPI) GetSyncStatus(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "需要登录")
		return
	}

	userIDStr := userID.(string)

	// 获取同步状态
	status := api.syncService.GetSyncStatus(userIDStr)

	shared.Success(c, http.StatusOK, "获取成功", status)
}

// SyncProgressRequest 同步进度请求
type SyncProgressRequest struct {
	BookID    string  `json:"bookId" binding:"required"`
	ChapterID string  `json:"chapterId" binding:"required"`
	Progress  float64 `json:"progress" binding:"required,min=0,max=1"`
	DeviceID  string  `json:"deviceId" binding:"required"`
}

// MergeProgressRequest 合并进度请求
type MergeProgressRequest struct {
	Progresses []OfflineProgressItem `json:"progresses" binding:"required"`
}

// OfflineProgressItem 离线进度项
type OfflineProgressItem struct {
	BookID    string  `json:"bookId" binding:"required"`
	ChapterID string  `json:"chapterId" binding:"required"`
	Progress  float64 `json:"progress" binding:"required,min=0,max=1"`
	Timestamp string  `json:"timestamp" binding:"required"`
	DeviceID  string  `json:"deviceId" binding:"required"`
}
