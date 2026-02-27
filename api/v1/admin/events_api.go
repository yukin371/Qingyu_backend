package admin

import (
	"time"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/pkg/response"
	eventservice "Qingyu_backend/service/events"
)

// EventsAdminAPI 事件管理API
type EventsAdminAPI struct {
	eventBus eventservice.PersistedEventBusInterface
}

// NewEventsAdminAPI 创建事件管理API实例
func NewEventsAdminAPI(eventBus eventservice.PersistedEventBusInterface) *EventsAdminAPI {
	return &EventsAdminAPI{
		eventBus: eventBus,
	}
}

// ReplayEventsRequest 事件回放请求
type ReplayEventsRequest struct {
	EventType string `json:"event_type" binding:"omitempty"` // 事件类型（可选）
	From      string `json:"from" binding:"omitempty"`       // 开始时间（ISO8601格式，可选）
	To        string `json:"to" binding:"omitempty"`         // 结束时间（ISO8601格式，可选）
	Limit     int    `json:"limit" binding:"omitempty"`      // 限制数量（可选，默认1000，最大10000）
	DryRun    bool   `json:"dry_run" binding:"omitempty"`    // 是否为dry-run模式（可选，默认false）
}

// ReplayEvents 事件回放
//
//	@Summary		事件回放
//	@Description	管理员回放历史事件，支持按类型、时间范围过滤，支持dry-run模式
//	@Tags			Admin-Events
//	@Accept			json
//	@Produce		json
//	@Param			request	body		ReplayEventsRequest	true	"回放请求参数"
//	@Success		200		{object}	response.APIResponse
//	@Failure		400		{object}	response.APIResponse
//	@Failure		401		{object}	response.APIResponse
//	@Failure		403		{object}	response.APIResponse
//	@Failure		500		{object}	response.APIResponse
//	@Router			/api/v1/admin/events/replay [post]
func (api *EventsAdminAPI) ReplayEvents(c *gin.Context) {
	var req ReplayEventsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	// 参数验证
	if req.Limit == 0 {
		req.Limit = 1000 // 默认值
	}
	if req.Limit < 1 || req.Limit > 10000 {
		response.BadRequest(c, "参数错误", "limit必须在1-10000之间")
		return
	}

	// 构建EventFilter
	filter := eventservice.EventFilter{
		EventType: req.EventType,
		Limit:     int64(req.Limit),
		DryRun:    req.DryRun,
	}

	// 解析时间范围
	if req.From != "" {
		fromTime, err := time.Parse(time.RFC3339, req.From)
		if err != nil {
			response.BadRequest(c, "参数错误", "from时间格式无效，请使用ISO8601格式")
			return
		}
		filter.StartTime = &fromTime
	}

	if req.To != "" {
		toTime, err := time.Parse(time.RFC3339, req.To)
		if err != nil {
			response.BadRequest(c, "参数错误", "to时间格式无效，请使用ISO8601格式")
			return
		}
		filter.EndTime = &toTime
	}

	// 调用Replay
	// 使用nil handler，因为我们只是统计事件，不需要实际处理
	result, err := api.eventBus.Replay(c.Request.Context(), nil, filter)
	if err != nil {
		c.Error(err)
		return
	}

	// 构建响应
	data := gin.H{
		"replayed_count": result.ReplayedCount,
		"failed_count":   result.FailedCount,
		"skipped_count":  result.SkippedCount,
		"duration_ms":    result.Duration.Milliseconds(),
	}

	// 如果是dry-run模式，添加提示信息
	if req.DryRun {
		data["message"] = "dry-run模式：未实际执行事件处理"
	}

	response.Success(c, data)
}
