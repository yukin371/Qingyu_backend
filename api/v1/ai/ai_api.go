package ai

import (
	"net/http"
	"time"

	aiService "Qingyu_backend/service/ai"
	"github.com/gin-gonic/gin"
)

// AIApi AI相关API控制器
type AIApi struct {
	service *aiService.Service
}

// NewAIApi 创建AI API控制器
func NewAIApi(service *aiService.Service) *AIApi {
	return &AIApi{
		service: service,
	}
}

// GenerateContent 生成内容
// POST /api/v1/ai/generate
func (a *AIApi) GenerateContent(c *gin.Context) {
	var req aiService.GenerateContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      400,
			"message":   "请求参数错误: " + err.Error(),
			"timestamp": getTimestamp(),
		})
		return
	}

	// 验证必填字段
	if req.ProjectID == "" || req.Prompt == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      400,
			"message":   "项目ID和提示词不能为空",
			"timestamp": getTimestamp(),
		})
		return
	}

	// 调用服务生成内容
	response, err := a.service.GenerateContent(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":      500,
			"message":   "生成内容失败: " + err.Error(),
			"timestamp": getTimestamp(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":      0,
		"message":   "success",
		"data":      response,
		"timestamp": getTimestamp(),
	})
}

// AnalyzeContent 分析内容
// POST /api/v1/ai/analyze
func (a *AIApi) AnalyzeContent(c *gin.Context) {
	var req aiService.AnalyzeContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      400,
			"message":   "请求参数错误: " + err.Error(),
			"timestamp": getTimestamp(),
		})
		return
	}

	// 验证必填字段
	if req.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      400,
			"message":   "分析内容不能为空",
			"timestamp": getTimestamp(),
		})
		return
	}

	// 设置默认分析类型
	if req.AnalysisType == "" {
		req.AnalysisType = "general"
	}

	// 调用服务分析内容
	response, err := a.service.AnalyzeContent(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":      500,
			"message":   "分析内容失败: " + err.Error(),
			"timestamp": getTimestamp(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":      0,
		"message":   "success",
		"data":      response,
		"timestamp": getTimestamp(),
	})
}

// ContinueWriting 续写内容
// POST /api/v1/ai/continue
func (a *AIApi) ContinueWriting(c *gin.Context) {
	var req aiService.ContinueWritingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      400,
			"message":   "请求参数错误: " + err.Error(),
			"timestamp": getTimestamp(),
		})
		return
	}

	// 验证必填字段
	if req.ProjectID == "" || req.ChapterID == "" || req.CurrentText == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      400,
			"message":   "项目ID、章节ID和当前文本不能为空",
			"timestamp": getTimestamp(),
		})
		return
	}

	// 调用服务续写内容
	response, err := a.service.ContinueWriting(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":      500,
			"message":   "续写内容失败: " + err.Error(),
			"timestamp": getTimestamp(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":      0,
		"message":   "success",
		"data":      response,
		"timestamp": getTimestamp(),
	})
}

// OptimizeText 优化文本
// POST /api/v1/ai/optimize
func (a *AIApi) OptimizeText(c *gin.Context) {
	var req aiService.OptimizeTextRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      400,
			"message":   "请求参数错误: " + err.Error(),
			"timestamp": getTimestamp(),
		})
		return
	}

	// 验证必填字段
	if req.ProjectID == "" || req.OriginalText == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      400,
			"message":   "项目ID和原始文本不能为空",
			"timestamp": getTimestamp(),
		})
		return
	}

	// 设置默认优化类型
	if req.OptimizeType == "" {
		req.OptimizeType = "style"
	}

	// 调用服务优化文本
	response, err := a.service.OptimizeText(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":      500,
			"message":   "优化文本失败: " + err.Error(),
			"timestamp": getTimestamp(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":      0,
		"message":   "success",
		"data":      response,
		"timestamp": getTimestamp(),
	})
}

// GenerateOutline 生成大纲
// POST /api/v1/ai/outline
func (a *AIApi) GenerateOutline(c *gin.Context) {
	var req aiService.GenerateOutlineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      400,
			"message":   "请求参数错误: " + err.Error(),
			"timestamp": getTimestamp(),
		})
		return
	}

	// 验证必填字段
	if req.ProjectID == "" || req.Theme == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      400,
			"message":   "项目ID和主题不能为空",
			"timestamp": getTimestamp(),
		})
		return
	}

	// 调用服务生成大纲
	response, err := a.service.GenerateOutline(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":      500,
			"message":   "生成大纲失败: " + err.Error(),
			"timestamp": getTimestamp(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":      0,
		"message":   "success",
		"data":      response,
		"timestamp": getTimestamp(),
	})
}

// GetContextInfo 获取上下文信息
// GET /api/v1/ai/context/:projectId/:chapterId
func (a *AIApi) GetContextInfo(c *gin.Context) {
	projectID := c.Param("projectId")
	chapterID := c.Param("chapterId")

	if projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      400,
			"message":   "项目ID不能为空",
			"timestamp": getTimestamp(),
		})
		return
	}

	// 调用服务获取上下文信息
	context, err := a.service.GetContextInfo(c.Request.Context(), projectID, chapterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":      500,
			"message":   "获取上下文信息失败: " + err.Error(),
			"timestamp": getTimestamp(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":      0,
		"message":   "success",
		"data":      context,
		"timestamp": getTimestamp(),
	})
}

// UpdateContextFeedback 更新上下文反馈
// POST /api/v1/ai/context/feedback
func (a *AIApi) UpdateContextFeedback(c *gin.Context) {
	type FeedbackRequest struct {
		ProjectID string `json:"projectId" binding:"required"`
		ChapterID string `json:"chapterId"`
		Feedback  string `json:"feedback" binding:"required"`
	}

	var req FeedbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      400,
			"message":   "请求参数错误: " + err.Error(),
			"timestamp": getTimestamp(),
		})
		return
	}

	// 调用服务更新上下文反馈
	err := a.service.UpdateContextWithFeedback(c.Request.Context(), req.ProjectID, req.ChapterID, req.Feedback)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":      500,
			"message":   "更新上下文反馈失败: " + err.Error(),
			"timestamp": getTimestamp(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":      0,
		"message":   "success",
		"timestamp": getTimestamp(),
	})
}

// TextGeneration 文本生成请求
func (a *AIApi) TextGeneration(c *gin.Context) {
	var req aiService.GenerateContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      400,
			"message":   "请求参数错误: " + err.Error(),
			"timestamp": getTimestamp(),
		})
		return
	}

	// 调用生成内容服务
	response, err := a.service.GenerateContent(c.Request.Context(), &req)
	if err != nil {
		// 处理不同类型的错误
		statusCode := http.StatusInternalServerError
		errorCode := 10009

		c.JSON(statusCode, gin.H{
			"code":      errorCode,
			"message":   err.Error(),
			"timestamp": getTimestamp(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":      0,
		"message":   "success",
		"data":      response,
		"timestamp": getTimestamp(),
	})
}

// getTimestamp 获取当前时间戳
func getTimestamp() int64 {
	return time.Now().Unix()
}