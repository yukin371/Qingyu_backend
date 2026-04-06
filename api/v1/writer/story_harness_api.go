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
}

// NewStoryHarnessApi 创建 StoryHarnessApi 实例
func NewStoryHarnessApi(contextSvc *storyharness.ContextService) *StoryHarnessApi {
	return &StoryHarnessApi{contextSvc: contextSvc}
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
