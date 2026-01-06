package admin

import (
	"net/http"
	"strconv"

	"Qingyu_backend/api/v1/shared"
	messagingService "Qingyu_backend/service/messaging"

	"github.com/gin-gonic/gin"
)

// AnnouncementAPI 公告管理API
type AnnouncementAPI struct {
	announcementService messagingService.AnnouncementService
}

// NewAnnouncementAPI 创建公告管理API实例
func NewAnnouncementAPI(announcementService messagingService.AnnouncementService) *AnnouncementAPI {
	return &AnnouncementAPI{
		announcementService: announcementService,
	}
}

// GetAnnouncements 获取公告列表
// @Summary 获取公告列表
// @Description 获取公告列表，支持筛选和分页
// @Tags 管理员-公告管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param isActive query boolean false "是否激活"
// @Param type query string false "类型(info/warning/notice)"
// @Param targetUsers query string false "目标用户(all/reader/writer/admin)"
// @Param limit query int false "每页数量" default(20)
// @Param offset query int false "偏移量" default(0)
// @Param sortBy query string false "排序字段(priority/created_at/view_count)" default(priority)
// @Param sortOrder query string false "排序方向(asc/desc)" default(desc)
// @Success 200 {object} shared.APIResponse{data=messagingService.GetAnnouncementsResponse}
// @Failure 401 {object} shared.ErrorResponse
// @Failure 500 {object} shared.ErrorResponse
// @Router /api/v1/admin/announcements [get]
func (api *AnnouncementAPI) GetAnnouncements(c *gin.Context) {
	// 解析查询参数
	req := &messagingService.GetAnnouncementsRequest{
		Limit:     20,
		Offset:    0,
		SortBy:    "priority",
		SortOrder: "desc",
	}

	// 解析isActive
	if isActiveStr := c.Query("isActive"); isActiveStr != "" {
		isActive, err := strconv.ParseBool(isActiveStr)
		if err == nil {
			req.IsActive = &isActive
		}
	}

	// 解析type
	if announcementType := c.Query("type"); announcementType != "" {
		req.Type = &announcementType
	}

	// 解析targetUsers
	if targetUsers := c.Query("targetUsers"); targetUsers != "" {
		req.TargetRole = &targetUsers
	}

	// 解析limit
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			req.Limit = limit
		}
	}

	// 解析offset
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			req.Offset = offset
		}
	}

	// 解析sortBy
	if sortBy := c.Query("sortBy"); sortBy != "" {
		req.SortBy = sortBy
	}

	// 解析sortOrder
	if sortOrder := c.Query("sortOrder"); sortOrder != "" {
		req.SortOrder = sortOrder
	}

	// 调用Service层
	resp, err := api.announcementService.GetAnnouncements(c.Request.Context(), req)
	if err != nil {
		shared.InternalError(c, "获取公告列表失败", err)
		return
	}

	shared.Success(c, http.StatusOK, "获取公告列表成功", resp)
}

// GetAnnouncementByID 获取公告详情
// @Summary 获取公告详情
// @Description 根据ID获取公告详情
// @Tags 管理员-公告管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "公告ID"
// @Success 200 {object} shared.APIResponse{data=models.Announcement}
// @Failure 400 {object} shared.ErrorResponse
// @Failure 404 {object} shared.ErrorResponse
// @Failure 500 {object} shared.ErrorResponse
// @Router /api/v1/admin/announcements/{id} [get]
func (api *AnnouncementAPI) GetAnnouncementByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		shared.BadRequest(c, "公告ID不能为空", "")
		return
	}

	announcement, err := api.announcementService.GetAnnouncementByID(c.Request.Context(), id)
	if err != nil {
		shared.HandleServiceError(c, err)
		return
	}

	shared.Success(c, http.StatusOK, "获取公告成功", announcement)
}

// CreateAnnouncement 创建公告
// @Summary 创建公告
// @Description 创建新的公告
// @Tags 管理员-公告管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body messagingService.CreateAnnouncementRequest true "创建公告请求"
// @Success 201 {object} shared.APIResponse{data=models.Announcement}
// @Failure 400 {object} shared.ErrorResponse
// @Failure 401 {object} shared.ErrorResponse
// @Failure 500 {object} shared.ErrorResponse
// @Router /api/v1/admin/announcements [post]
func (api *AnnouncementAPI) CreateAnnouncement(c *gin.Context) {
	var req messagingService.CreateAnnouncementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.BadRequest(c, "参数错误", err.Error())
		return
	}

	// 从上下文获取当前用户ID
	if userID, exists := c.Get("user_id"); exists {
		req.CreatedBy = userID.(string)
	}

	announcement, err := api.announcementService.CreateAnnouncement(c.Request.Context(), &req)
	if err != nil {
		shared.HandleServiceError(c, err)
		return
	}

	shared.Success(c, http.StatusCreated, "创建公告成功", announcement)
}

// UpdateAnnouncement 更新公告
// @Summary 更新公告
// @Description 更新公告信息
// @Tags 管理员-公告管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "公告ID"
// @Param request body messagingService.UpdateAnnouncementRequest true "更新公告请求"
// @Success 200 {object} shared.APIResponse
// @Failure 400 {object} shared.ErrorResponse
// @Failure 401 {object} shared.ErrorResponse
// @Failure 404 {object} shared.ErrorResponse
// @Failure 500 {object} shared.ErrorResponse
// @Router /api/v1/admin/announcements/{id} [put]
func (api *AnnouncementAPI) UpdateAnnouncement(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		shared.BadRequest(c, "公告ID不能为空", "")
		return
	}

	var req messagingService.UpdateAnnouncementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.BadRequest(c, "参数错误", err.Error())
		return
	}

	if err := api.announcementService.UpdateAnnouncement(c.Request.Context(), id, &req); err != nil {
		shared.HandleServiceError(c, err)
		return
	}

	shared.Success(c, http.StatusOK, "更新公告成功", nil)
}

// DeleteAnnouncement 删除公告
// @Summary 删除公告
// @Description 删除指定的公告
// @Tags 管理员-公告管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "公告ID"
// @Success 200 {object} shared.APIResponse
// @Failure 400 {object} shared.ErrorResponse
// @Failure 401 {object} shared.ErrorResponse
// @Failure 404 {object} shared.ErrorResponse
// @Failure 500 {object} shared.ErrorResponse
// @Router /api/v1/admin/announcements/{id} [delete]
func (api *AnnouncementAPI) DeleteAnnouncement(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		shared.BadRequest(c, "公告ID不能为空", "")
		return
	}

	if err := api.announcementService.DeleteAnnouncement(c.Request.Context(), id); err != nil {
		shared.HandleServiceError(c, err)
		return
	}

	shared.Success(c, http.StatusOK, "删除公告成功", nil)
}

// BatchUpdateStatus 批量更新状态
// @Summary 批量更新公告状态
// @Description 批量启用或禁用公告
// @Tags 管理员-公告管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body messagingService.BatchUpdateAnnouncementStatusRequest true "批量更新状态请求"
// @Success 200 {object} shared.APIResponse
// @Failure 400 {object} shared.ErrorResponse
// @Failure 401 {object} shared.ErrorResponse
// @Failure 500 {object} shared.ErrorResponse
// @Router /api/v1/admin/announcements/batch-status [put]
func (api *AnnouncementAPI) BatchUpdateStatus(c *gin.Context) {
	var req messagingService.BatchUpdateAnnouncementStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.BadRequest(c, "参数错误", err.Error())
		return
	}

	if err := api.announcementService.BatchUpdateStatus(c.Request.Context(), &req); err != nil {
		shared.HandleServiceError(c, err)
		return
	}

	shared.Success(c, http.StatusOK, "批量更新状态成功", nil)
}

// BatchDeleteRequest 批量删除请求
type BatchDeleteRequest struct {
	IDs []string `json:"ids" validate:"required,min=1"`
}

// BatchDelete 批量删除公告
// @Summary 批量删除公告
// @Description 批量删除公告
// @Tags 管理员-公告管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body BatchDeleteRequest true "批量删除请求"
// @Success 200 {object} shared.APIResponse
// @Failure 400 {object} shared.ErrorResponse
// @Failure 401 {object} shared.ErrorResponse
// @Failure 500 {object} shared.ErrorResponse
// @Router /api/v1/admin/announcements/batch-delete [delete]
func (api *AnnouncementAPI) BatchDelete(c *gin.Context) {
	var req BatchDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.BadRequest(c, "参数错误", err.Error())
		return
	}

	if err := api.announcementService.BatchDelete(c.Request.Context(), req.IDs); err != nil {
		shared.HandleServiceError(c, err)
		return
	}

	shared.Success(c, http.StatusOK, "批量删除成功", nil)
}
