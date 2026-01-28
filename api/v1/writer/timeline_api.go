package writer

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/shared"
	"Qingyu_backend/service/interfaces"
	"Qingyu_backend/pkg/response"
)

// TimelineApi 时间线API处理器
type TimelineApi struct {
	timelineService interfaces.TimelineService
}

// NewTimelineApi 创建TimelineApi实例
func NewTimelineApi(timelineService interfaces.TimelineService) *TimelineApi {
	return &TimelineApi{
		timelineService: timelineService,
	}
}

// CreateTimeline 创建时间线
func (api *TimelineApi) CreateTimeline(c *gin.Context) {
	projectID := c.Param("projectId")
	if projectID == "" {
		response.BadRequest(c,  "项目ID不能为空", "")
		return
	}

	var req interfaces.CreateTimelineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c,  "参数错误", err.Error())
		return
	}

	timeline, err := api.timelineService.CreateTimeline(c.Request.Context(), projectID, &req)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	shared.Success(c, http.StatusCreated, "创建成功", timeline)
}

// GetTimeline 获取时间线详情
func (api *TimelineApi) GetTimeline(c *gin.Context) {
	timelineID := c.Param("timelineId")
	projectID := c.Query("projectId")

	if timelineID == "" || projectID == "" {
		response.BadRequest(c,  "参数错误", "timelineId和projectId不能为空")
		return
	}

	timeline, err := api.timelineService.GetTimeline(c.Request.Context(), timelineID, projectID)
	if err != nil {
		shared.Error(c, http.StatusNotFound, "时间线不存在", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取成功", timeline)
}

// ListTimelines 获取项目时间线列表
func (api *TimelineApi) ListTimelines(c *gin.Context) {
	projectID := c.Param("projectId")
	if projectID == "" {
		response.BadRequest(c,  "项目ID不能为空", "")
		return
	}

	timelines, err := api.timelineService.ListTimelines(c.Request.Context(), projectID)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	shared.Success(c, http.StatusOK, "获取成功", timelines)
}

// DeleteTimeline 删除时间线
func (api *TimelineApi) DeleteTimeline(c *gin.Context) {
	timelineID := c.Param("timelineId")
	projectID := c.Query("projectId")

	if timelineID == "" || projectID == "" {
		response.BadRequest(c,  "参数错误", "timelineId和projectId不能为空")
		return
	}

	err := api.timelineService.DeleteTimeline(c.Request.Context(), timelineID, projectID)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	shared.Success(c, http.StatusOK, "删除成功", nil)
}

// CreateTimelineEvent 创建时间线事件
func (api *TimelineApi) CreateTimelineEvent(c *gin.Context) {
	timelineID := c.Param("timelineId")
	projectID := c.Query("projectId")

	if timelineID == "" || projectID == "" {
		response.BadRequest(c,  "参数错误", "timelineId和projectId不能为空")
		return
	}

	var req interfaces.CreateTimelineEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c,  "参数错误", err.Error())
		return
	}

	// 确保timelineId一致
	req.TimelineID = timelineID

	event, err := api.timelineService.CreateEvent(c.Request.Context(), projectID, &req)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	shared.Success(c, http.StatusCreated, "创建成功", event)
}

// GetTimelineEvent 获取事件详情
func (api *TimelineApi) GetTimelineEvent(c *gin.Context) {
	eventID := c.Param("eventId")
	projectID := c.Query("projectId")

	if eventID == "" || projectID == "" {
		response.BadRequest(c,  "参数错误", "eventId和projectId不能为空")
		return
	}

	event, err := api.timelineService.GetEvent(c.Request.Context(), eventID, projectID)
	if err != nil {
		shared.Error(c, http.StatusNotFound, "事件不存在", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取成功", event)
}

// ListTimelineEvents 获取时间线事件列表
func (api *TimelineApi) ListTimelineEvents(c *gin.Context) {
	timelineID := c.Param("timelineId")
	if timelineID == "" {
		response.BadRequest(c,  "时间线ID不能为空", "")
		return
	}

	events, err := api.timelineService.ListEvents(c.Request.Context(), timelineID)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	shared.Success(c, http.StatusOK, "获取成功", events)
}

// UpdateTimelineEvent 更新时间线事件
func (api *TimelineApi) UpdateTimelineEvent(c *gin.Context) {
	eventID := c.Param("eventId")
	projectID := c.Query("projectId")

	if eventID == "" || projectID == "" {
		response.BadRequest(c,  "参数错误", "eventId和projectId不能为空")
		return
	}

	var req interfaces.UpdateTimelineEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c,  "参数错误", err.Error())
		return
	}

	event, err := api.timelineService.UpdateEvent(c.Request.Context(), eventID, projectID, &req)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	shared.Success(c, http.StatusOK, "更新成功", event)
}

// DeleteTimelineEvent 删除时间线事件
func (api *TimelineApi) DeleteTimelineEvent(c *gin.Context) {
	eventID := c.Param("eventId")
	projectID := c.Query("projectId")

	if eventID == "" || projectID == "" {
		response.BadRequest(c,  "参数错误", "eventId和projectId不能为空")
		return
	}

	err := api.timelineService.DeleteEvent(c.Request.Context(), eventID, projectID)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	shared.Success(c, http.StatusOK, "删除成功", nil)
}

// GetTimelineVisualization 获取时间线可视化数据
func (api *TimelineApi) GetTimelineVisualization(c *gin.Context) {
	timelineID := c.Param("timelineId")
	if timelineID == "" {
		response.BadRequest(c,  "时间线ID不能为空", "")
		return
	}

	visualization, err := api.timelineService.GetTimelineVisualization(c.Request.Context(), timelineID)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	shared.Success(c, http.StatusOK, "获取成功", visualization)
}
