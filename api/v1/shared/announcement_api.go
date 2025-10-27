package shared

import (
	"net/http"
	"strconv"

	sharedService "Qingyu_backend/service/shared"

	"github.com/gin-gonic/gin"
)

// AnnouncementPublicAPI 公告公开API
type AnnouncementPublicAPI struct {
	announcementService sharedService.AnnouncementService
}

// NewAnnouncementPublicAPI 创建公告公开API实例
func NewAnnouncementPublicAPI(announcementService sharedService.AnnouncementService) *AnnouncementPublicAPI {
	return &AnnouncementPublicAPI{
		announcementService: announcementService,
	}
}

// GetEffectiveAnnouncements 获取当前有效的公告
// @Summary 获取当前有效的公告
// @Description 获取当前时间有效的公告列表
// @Tags 公告
// @Accept json
// @Produce json
// @Param targetUsers query string false "目标用户(all/reader/writer)" default(all)
// @Param limit query int false "数量限制" default(10)
// @Success 200 {object} APIResponse{data=[]models.Announcement}
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/announcements/effective [get]
func (api *AnnouncementPublicAPI) GetEffectiveAnnouncements(c *gin.Context) {
	// 解析参数
	targetUsers := c.DefaultQuery("targetUsers", "all")
	limit := 10

	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	// 从上下文获取用户角色（如果已登录）
	if role, exists := c.Get("user_role"); exists {
		switch role.(string) {
		case "reader":
			targetUsers = "reader"
		case "writer":
			targetUsers = "writer"
		case "admin":
			targetUsers = "admin"
		}
	}

	// 调用Service层
	announcements, err := api.announcementService.GetEffectiveAnnouncements(c.Request.Context(), targetUsers, limit)
	if err != nil {
		InternalError(c, "获取公告失败", err)
		return
	}

	Success(c, http.StatusOK, "获取公告成功", announcements)
}

// IncrementViewCount 增加公告查看次数
// @Summary 增加公告查看次数
// @Description 记录公告被查看
// @Tags 公告
// @Accept json
// @Produce json
// @Param id path string true "公告ID"
// @Success 200 {object} APIResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/announcements/{id}/view [post]
func (api *AnnouncementPublicAPI) IncrementViewCount(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		BadRequest(c, "公告ID不能为空", "")
		return
	}

	if err := api.announcementService.IncrementViewCount(c.Request.Context(), id); err != nil {
		HandleServiceError(c, err)
		return
	}

	Success(c, http.StatusOK, "记录成功", nil)
}
