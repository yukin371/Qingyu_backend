package writer

import (
	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/shared"
	"Qingyu_backend/models/writer"
	"Qingyu_backend/pkg/response"
	"Qingyu_backend/service/writer/storyharness"
)

// ChangeRequestApi 变更建议 API 处理器
type ChangeRequestApi struct {
	crSvc *storyharness.ChangeRequestService
}

// NewChangeRequestApi 创建 ChangeRequestApi 实例
func NewChangeRequestApi(crSvc *storyharness.ChangeRequestService) *ChangeRequestApi {
	return &ChangeRequestApi{crSvc: crSvc}
}

// ListChangeRequests 获取章节建议列表
// @Summary 获取章节建议列表
// @Description 获取指定章节的变更建议列表
// @Tags 变更建议
// @Accept json
// @Produce json
// @Param id path string true "项目ID"
// @Param chapterId path string true "章节ID"
// @Param status query string false "状态筛选 (pending/accepted/ignored/deferred)"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/writer/projects/{id}/chapters/{chapterId}/change-requests [get]
func (api *ChangeRequestApi) ListChangeRequests(c *gin.Context) {
	projectID := c.Param("id")
	chapterID := c.Param("chapterId")

	if projectID == "" || chapterID == "" {
		response.BadRequest(c, "参数错误", "projectID 和 chapterID 不能为空")
		return
	}

	items, err := api.crSvc.ListByChapter(c.Request.Context(), projectID, chapterID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, gin.H{
		"items": items,
		"total": len(items),
	})
}

// GetChangeRequest 获取建议详情
// @Summary 获取建议详情
// @Description 根据ID获取变更建议详情
// @Tags 变更建议
// @Accept json
// @Produce json
// @Param requestId path string true "建议ID"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/writer/change-requests/{requestId} [get]
func (api *ChangeRequestApi) GetChangeRequest(c *gin.Context) {
	requestID := c.Param("requestId")
	if requestID == "" {
		response.BadRequest(c, "参数错误", "requestId 不能为空")
		return
	}

	cr, err := api.crSvc.GetByID(c.Request.Context(), requestID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, cr)
}

// ProcessChangeRequest 处理建议
// @Summary 处理建议
// @Description 接受/忽略/延后一条变更建议
// @Tags 变更建议
// @Accept json
// @Produce json
// @Param requestId path string true "建议ID"
// @Param request body object true "处理请求"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/writer/change-requests/{requestId}/status [put]
func (api *ChangeRequestApi) ProcessChangeRequest(c *gin.Context) {
	requestID := c.Param("requestId")
	if requestID == "" {
		response.BadRequest(c, "参数错误", "requestId 不能为空")
		return
	}

	var req struct {
		Status writer.ChangeRequestStatus `json:"status" validate:"required"`
	}
	if !shared.BindJSON(c, &req) {
		return
	}

	userID := shared.GetUserIDOptional(c)

	err := api.crSvc.Process(c.Request.Context(), requestID, req.Status, userID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, nil)
}
