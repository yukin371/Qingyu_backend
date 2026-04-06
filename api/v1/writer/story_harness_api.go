package writer

import (
	"github.com/gin-gonic/gin"

	"Qingyu_backend/pkg/response"
	"Qingyu_backend/service/writer/storyharness"
)

// StoryHarnessApi Story Harness API 处理器
// 负责章节上下文（Context Lens）数据的查询
type StoryHarnessApi struct {
	contextSvc *storyharness.ContextService
	indexerSvc *storyharness.IndexerService
	crSvc      *storyharness.ChangeRequestService
}

// NewStoryHarnessApi 创建 StoryHarnessApi 实例
func NewStoryHarnessApi(
	contextSvc *storyharness.ContextService,
	indexerSvc *storyharness.IndexerService,
	crSvc *storyharness.ChangeRequestService,
) *StoryHarnessApi {
	return &StoryHarnessApi{
		contextSvc: contextSvc,
		indexerSvc: indexerSvc,
		crSvc:      crSvc,
	}
}

// GetChapterContext 获取章节上下文
// @Summary 获取章节上下文
// @Description 获取当前章节的角色快照、关系、作用域等 Context Lens 数据
// @Tags Story Harness
// @Accept json
// @Produce json
// @Param id path string true "项目ID"
// @Param chapterId path string true "章节ID"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /api/v1/writer/projects/{id}/chapters/{chapterId}/context [get]
func (api *StoryHarnessApi) GetChapterContext(c *gin.Context) {
	projectID := c.Param("id")
	chapterID := c.Param("chapterId")

	if projectID == "" || chapterID == "" {
		response.BadRequest(c, "参数错误", "projectID 和 chapterID 不能为空")
		return
	}

	data, err := api.contextSvc.GetChapterContext(c.Request.Context(), projectID, chapterID)
	if err != nil {
		c.Error(err)
		return
	}

	// 组装响应 DTO
	result := gin.H{
		"characters": data.Characters,
		"relations":  data.Relations,
		"pendingCRs": data.PendingCRs,
	}

	response.Success(c, result)
}

// TriggerChapterIndex 手动触发章节索引
// @Summary 触发章节索引
// @Description 对当前章节运行最小规则引擎，生成变更建议批次
// @Tags Story Harness
// @Accept json
// @Produce json
// @Param id path string true "项目ID"
// @Param chapterId path string true "章节ID"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /api/v1/writer/projects/{id}/chapters/{chapterId}/trigger-index [post]
func (api *StoryHarnessApi) TriggerChapterIndex(c *gin.Context) {
	projectID := c.Param("id")
	chapterID := c.Param("chapterId")

	if projectID == "" || chapterID == "" {
		response.BadRequest(c, "参数错误", "projectID 和 chapterID 不能为空")
		return
	}
	if api.indexerSvc == nil {
		response.BadRequest(c, "索引服务不可用", "indexer service is nil")
		return
	}

	result, err := api.indexerSvc.TriggerChapterIndex(c.Request.Context(), projectID, chapterID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, result)
}

// RebuildChapterProjection 手动重建章节投影
// @Summary 重建章节投影
// @Description 基于当前章节已 accepted 的建议重新构建 projection，用于修复或回放
// @Tags Story Harness
// @Accept json
// @Produce json
// @Param id path string true "项目ID"
// @Param chapterId path string true "章节ID"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /api/v1/writer/projects/{id}/chapters/{chapterId}/rebuild-projection [post]
func (api *StoryHarnessApi) RebuildChapterProjection(c *gin.Context) {
	projectID := c.Param("id")
	chapterID := c.Param("chapterId")

	if projectID == "" || chapterID == "" {
		response.BadRequest(c, "参数错误", "projectID 和 chapterID 不能为空")
		return
	}
	if api.crSvc == nil {
		response.BadRequest(c, "建议服务不可用", "change request service is nil")
		return
	}

	result, err := api.crSvc.RebuildProjection(c.Request.Context(), projectID, chapterID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, result)
}
