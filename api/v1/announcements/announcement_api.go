package announcements

import (
	"strconv"

	"Qingyu_backend/api/v1/shared"
	messagingModel "Qingyu_backend/models/messaging" // Imported for Swagger annotations
	messagingService "Qingyu_backend/service/messaging"

	"github.com/gin-gonic/gin"
)

// AnnouncementPublicAPI 公告公开API处理器
type AnnouncementPublicAPI struct {
	announcementService messagingService.AnnouncementService
}

// NewAnnouncementPublicAPI 创建公告公开API实例
func NewAnnouncementPublicAPI(announcementService messagingService.AnnouncementService) *AnnouncementPublicAPI {
	return &AnnouncementPublicAPI{
		announcementService: announcementService,
	}
}

// ============ 公告公开API ============

// GetEffectiveAnnouncements 获取有效公告列表
//
//	@Summary		获取有效公告
//	@Description	获取当前有效的公告列表（公开访问，无需登录）
//	@Tags			公告
//	@Accept			json
//	@Produce		json
//	@Param			targetRole	query		string	false	"目标角色(all/reader/writer/admin)"	default(all)
//	@Param			limit			query		int		false	"限制数量"	default(10)
//	@Success		200				{object}	shared.APIResponse{data=[]models.Announcement}
//	@Failure		500				{object}	shared.ErrorResponse
//	@Router			/api/v1/announcements/effective [get]
func (api *AnnouncementPublicAPI) GetEffectiveAnnouncements(c *gin.Context) {
	// 1. 获取参数
	targetRole := c.DefaultQuery("targetRole", "all")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// 2. 限制最大数量
	if limit > 50 {
		limit = 50
	}
	if limit <= 0 {
		limit = 10
	}

	// 3. 获取有效公告
	announcements, err := api.announcementService.GetEffectiveAnnouncements(c.Request.Context(), targetRole, limit)
	if err != nil {
		shared.Error(c, 500, "获取公告失败", err.Error())
		return
	}

	// 4. 返回结果
	shared.Success(c, 200, "获取成功", announcements)
}

// IncrementViewCount 增加公告查看次数
//
//	@Summary		增加查看次数
//	@Description	增加指定公告的查看次数
//	@Tags			公告
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"公告ID"
//	@Success		200		{object}	shared.APIResponse
//	@Failure		400		{object}	shared.ErrorResponse
//	@Failure		500		{object}	shared.ErrorResponse
//	@Router			/api/v1/announcements/{id}/view [post]
func (api *AnnouncementPublicAPI) IncrementViewCount(c *gin.Context) {
	// 1. 获取公告ID
	id := c.Param("id")
	if id == "" {
		shared.BadRequest(c, "参数错误", "公告ID不能为空")
		return
	}

	// 2. 增加查看次数
	err := api.announcementService.IncrementViewCount(c.Request.Context(), id)
	if err != nil {
		shared.Error(c, 500, "操作失败", err.Error())
		return
	}

	// 3. 返回成功
	shared.Success(c, 200, "操作成功", nil)
}

// GetAnnouncementByID 获取公告详情
//
//	@Summary		获取公告详情
//	@Description	根据ID获取公告详情（公开访问）
//	@Tags			公告
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"公告ID"
//	@Success		200		{object}	shared.APIResponse{data=models.Announcement}
//	@Failure		400		{object}	shared.ErrorResponse
//	@Failure		404		{object}	shared.ErrorResponse
//	@Failure		500		{object}	shared.ErrorResponse
//	@Router			/api/v1/announcements/{id} [get]
func (api *AnnouncementPublicAPI) GetAnnouncementByID(c *gin.Context) {
	// 1. 获取公告ID
	id := c.Param("id")
	if id == "" {
		shared.BadRequest(c, "参数错误", "公告ID不能为空")
		return
	}

	// 2. 获取公告详情
	announcement, err := api.announcementService.GetAnnouncementByID(c.Request.Context(), id)
	if err != nil {
		// 检查是否为404错误
		if len(err.Error()) > 0 && (err.Error()[0:4] == "not_" || err.Error()[0:4] == "NOT_") {
			shared.NotFound(c, "公告不存在")
			return
		}
		shared.Error(c, 500, "获取公告失败", err.Error())
		return
	}

	if announcement == nil {
		shared.NotFound(c, "公告不存在")
		return
	}

	// 3. 返回结果
	shared.Success(c, 200, "获取成功", announcement)
}

var _ = messagingModel.Announcement{}
