package ai

import (
	"fmt"

	aiService "Qingyu_backend/service/ai"

	"Qingyu_backend/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// StoryWriteApi 故事上下文写作 API
type StoryWriteApi struct {
	contextEngine *aiService.StoryContextEngine
	quotaService  *aiService.QuotaService
}

// NewStoryWriteApi 创建故事上下文写作 API 实例
func NewStoryWriteApi(
	contextEngine *aiService.StoryContextEngine,
	quotaService *aiService.QuotaService,
) *StoryWriteApi {
	return &StoryWriteApi{
		contextEngine: contextEngine,
		quotaService:  quotaService,
	}
}

// StoryGenerateRequest 故事生成请求
type StoryGenerateRequest struct {
	ProjectID    string `json:"projectId" binding:"required"`
	DocumentID   string `json:"documentId" binding:"required"`
	Mode         string `json:"mode" binding:"required,oneof=continue rewrite suggest"`
	Instruction  string `json:"instruction,omitempty"`
	SelectedText string `json:"selectedText,omitempty"`
}

// Generate 统一 AI 生成入口
// @Summary 故事生成
// @Description 基于三层上下文进行 AI 故事生成
// @Tags AI写作
// @Accept json
// @Produce json
// @Param request body StoryGenerateRequest true "故事生成请求"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/ai/story/generate [post]
func (api *StoryWriteApi) Generate(c *gin.Context) {
	var req StoryGenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	// 生成请求ID
	requestID := uuid.New().String()
	c.Set("requestID", requestID)
	c.Set("aiService", "story_write")

	// 组装三层上下文
	sc, err := api.contextEngine.BuildStoryContext(
		c.Request.Context(),
		req.ProjectID,
		req.DocumentID,
		req.Mode,
		req.Instruction,
		req.SelectedText,
	)
	if err != nil {
		response.InternalError(c, fmt.Errorf("上下文组装失败: %w", err))
		return
	}

	// 构建 prompt（包级函数）
	prompt := aiService.BuildPrompt(sc)

	// TODO: Phase 2 - 通过 gRPC 发送给 AI Service
	// 目前返回上下文预览，用于调试验证
	response.SuccessWithMessage(c, "上下文组装成功", gin.H{
		"prompt":       prompt,
		"contextStats": sc,
	})
}

// ContextPreview 上下文预览（调试用）
// @Summary 上下文预览
// @Description 预览三层上下文组装结果，用于调试
// @Tags AI写作
// @Produce json
// @Param projectId query string true "项目ID"
// @Param documentId query string true "文档ID"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/ai/story/context-preview [get]
func (api *StoryWriteApi) ContextPreview(c *gin.Context) {
	projectID := c.Query("projectId")
	documentID := c.Query("documentId")
	if projectID == "" || documentID == "" {
		response.BadRequest(c, "参数错误", "需要 projectId 和 documentId")
		return
	}

	sc, err := api.contextEngine.BuildStoryContext(
		c.Request.Context(),
		projectID,
		documentID,
		"continue",
		"",
		"",
	)
	if err != nil {
		response.InternalError(c, fmt.Errorf("上下文组装失败: %w", err))
		return
	}

	prompt := aiService.BuildPrompt(sc)

	response.Success(c, gin.H{
		"storyContext": sc,
		"prompt":       prompt,
		"promptLength": len(prompt),
	})
}

// UpdateSceneStateRequest 更新场景状态请求
type UpdateSceneStateRequest struct {
	SceneGoal      string `json:"sceneGoal,omitempty"`
	ActiveConflict string `json:"activeConflict,omitempty"`
}

// UpdateSceneState 更新场景状态
// @Summary 更新场景状态
// @Description 更新文档的场景目标和活跃冲突
// @Tags AI写作
// @Accept json
// @Produce json
// @Param id path string true "文档ID"
// @Param request body UpdateSceneStateRequest true "场景状态请求"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/ai/story/scene-state/{id} [put]
func (api *StoryWriteApi) UpdateSceneState(c *gin.Context) {
	documentID := c.Param("id")
	var req UpdateSceneStateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	// TODO: 调用 documentRepo.Update 更新 SceneGoal/ActiveConflict 字段
	response.SuccessWithMessage(c, "场景状态已更新", gin.H{
		"documentId":     documentID,
		"sceneGoal":      req.SceneGoal,
		"activeConflict": req.ActiveConflict,
	})
}
