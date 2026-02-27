package announcements

import (
	messagingModel "Qingyu_backend/models/messaging"
	"Qingyu_backend/api/v1/shared"
	"Qingyu_backend/pkg/response"
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
//	@Success		200				{object}	response.APIResponse
//	@Failure		500				{object}	response.ErrorResponse
//	@Router			/api/v1/announcements/effective [get]
func (api *AnnouncementPublicAPI) GetEffectiveAnnouncements(c *gin.Context) {
	// 1. 获取参数
	targetRole := c.DefaultQuery("targetRole", "all")
	limit := shared.GetIntParam(c, "limit", true, 10, 1, 50)

	// 2. 获取有效公告
	announcements, err := api.announcementService.GetEffectiveAnnouncements(c.Request.Context(), targetRole, limit)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	// 3. 返回结果
	response.Success(c, announcements)
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
	id, ok := shared.GetRequiredParam(c, "id", "公告ID")
	if !ok {
		return
	}

	// 2. 增加查看次数
	err := api.announcementService.IncrementViewCount(c.Request.Context(), id)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	// 3. 返回成功
	response.Success(c, nil)
}

// GetAnnouncementByID 获取公告详情
//
//	@Summary		获取公告详情
//	@Description	根据ID获取公告详情（公开访问）
//	@Tags			公告
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"公告ID"
//	@Success		200		{object}	shared.APIResponse
//	@Failure		400		{object}	shared.ErrorResponse
//	@Failure		404		{object}	shared.ErrorResponse
//	@Failure		500		{object}	shared.ErrorResponse
//	@Router			/api/v1/announcements/{id} [get]
func (api *AnnouncementPublicAPI) GetAnnouncementByID(c *gin.Context) {
	// 1. 获取公告ID
	id, ok := shared.GetRequiredParam(c, "id", "公告ID")
	if !ok {
		return
	}

	// 2. 获取公告详情
	announcement, err := api.announcementService.GetAnnouncementByID(c.Request.Context(), id)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	if announcement == nil {
		response.NotFound(c, "公告不存在")
		return
	}

	// 3. 返回结果
	response.Success(c, announcement)
}

var _ = messagingModel.Announcement{}
