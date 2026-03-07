package admin

import (
	"strconv"

	"Qingyu_backend/pkg/response"
	serviceInterfaces "Qingyu_backend/service/interfaces"

	"github.com/gin-gonic/gin"
)

type PublicationAdminAPI struct {
	publishService serviceInterfaces.PublishService
}

func NewPublicationAdminAPI(publishService serviceInterfaces.PublishService) *PublicationAdminAPI {
	return &PublicationAdminAPI{publishService: publishService}
}

type ReviewPublicationRequest struct {
	Action string `json:"action" binding:"required,oneof=approve reject"`
	Note   string `json:"note"`
}

func (api *PublicationAdminAPI) GetPendingPublications(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	records, total, err := api.publishService.GetPendingPublicationRecords(c.Request.Context(), page, pageSize)
	if err != nil {
		c.Error(err)
		return
	}

	response.Paginated(c, records, total, page, pageSize, "获取成功")
}

func (api *PublicationAdminAPI) ReviewPublication(c *gin.Context) {
	recordID := c.Param("id")
	if recordID == "" {
		response.BadRequest(c, "参数错误", "记录ID不能为空")
		return
	}

	var req ReviewPublicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	reviewerID := ""
	if value, exists := c.Get("user_id"); exists {
		if id, ok := value.(string); ok {
			reviewerID = id
		}
	}

	record, err := api.publishService.ReviewPublication(c.Request.Context(), recordID, reviewerID, req.Action == "approve", req.Note)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, record)
}
