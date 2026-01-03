package ai

import (
	"net/http"

	"Qingyu_backend/api/v1/shared"
	aiService "Qingyu_backend/service/ai"
	"Qingyu_backend/service/ai/dto"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// WritingAssistantApi 写作辅助API
// 提供内容总结、校对、敏感词检测等AI辅助功能
type WritingAssistantApi struct {
	summarizeService     *aiService.SummarizeService
	proofreadService     *aiService.ProofreadService
	sensitiveWordsService *aiService.SensitiveWordsService
}

// NewWritingAssistantApi 创建写作辅助API实例
func NewWritingAssistantApi(
	summarizeService *aiService.SummarizeService,
	proofreadService *aiService.ProofreadService,
	sensitiveWordsService *aiService.SensitiveWordsService,
) *WritingAssistantApi {
	return &WritingAssistantApi{
		summarizeService:      summarizeService,
		proofreadService:      proofreadService,
		sensitiveWordsService: sensitiveWordsService,
	}
}

// ===========================
// 内容总结 API
// ===========================

// SummarizeContent 总结文档内容
// @Summary 总结文档内容
// @Description 使用AI总结文档内容，支持简短摘要、详细摘要、关键点提取等多种模式
// @Tags AI写作辅助
// @Accept json
// @Produce json
// @Param request body dto.SummarizeRequest true "总结请求"
// @Success 200 {object} shared.APIResponse{data=dto.SummarizeResponse}
// @Failure 400 {object} shared.ErrorResponse
// @Failure 500 {object} shared.ErrorResponse
// @Router /api/v1/ai/writing/summarize [post]
func (api *WritingAssistantApi) SummarizeContent(c *gin.Context) {
	var req dto.SummarizeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	// 生成请求ID
	requestID := uuid.New().String()
	c.Set("requestID", requestID)
	c.Set("aiService", "summarize")

	// 调用服务
	result, err := api.summarizeService.SummarizeContent(c.Request.Context(), &req)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "总结失败", err.Error())
		return
	}

	// 设置Token使用信息
	c.Set("tokensUsed", result.TokensUsed)
	c.Set("aiModel", result.Model)

	shared.Success(c, http.StatusOK, "总结成功", result)
}

// SummarizeChapter 总结章节内容
// @Summary 总结章节内容
// @Description 自动提取章节要点、情节大纲、涉及角色等详细信息
// @Tags AI写作辅助
// @Accept json
// @Produce json
// @Param request body dto.ChapterSummaryRequest true "章节总结请求"
// @Success 200 {object} shared.APIResponse{data=dto.ChapterSummaryResponse}
// @Failure 400 {object} shared.ErrorResponse
// @Failure 500 {object} shared.ErrorResponse
// @Router /api/v1/ai/writing/summarize-chapter [post]
func (api *WritingAssistantApi) SummarizeChapter(c *gin.Context) {
	var req dto.ChapterSummaryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	// 生成请求ID
	requestID := uuid.New().String()
	c.Set("requestID", requestID)
	c.Set("aiService", "summarize_chapter")

	// 调用服务
	result, err := api.summarizeService.SummarizeChapter(c.Request.Context(), &req)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "章节总结失败", err.Error())
		return
	}

	// 设置Token使用信息
	c.Set("tokensUsed", result.TokensUsed)
	c.Set("aiModel", "summarize_service")

	shared.Success(c, http.StatusOK, "章节总结成功", result)
}

// ===========================
// 文本校对 API
// ===========================

// ProofreadContent 文本校对
// @Summary 文本校对
// @Description 检查拼写、语法、标点错误，返回修改建议列表和整体评分
// @Tags AI写作辅助
// @Accept json
// @Produce json
// @Param request body dto.ProofreadRequest true "校对请求"
// @Success 200 {object} shared.APIResponse{data=dto.ProofreadResponse}
// @Failure 400 {object} shared.ErrorResponse
// @Failure 500 {object} shared.ErrorResponse
// @Router /api/v1/ai/writing/proofread [post]
func (api *WritingAssistantApi) ProofreadContent(c *gin.Context) {
	var req dto.ProofreadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	// 生成请求ID
	requestID := uuid.New().String()
	c.Set("requestID", requestID)
	c.Set("aiService", "proofread")

	// 调用服务
	result, err := api.proofreadService.ProofreadContent(c.Request.Context(), &req)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "校对失败", err.Error())
		return
	}

	// 设置Token使用信息
	c.Set("tokensUsed", result.TokensUsed)
	c.Set("aiModel", result.Model)

	shared.Success(c, http.StatusOK, "校对完成", result)
}

// GetProofreadSuggestion 获取校对建议详情
// @Summary 获取校对建议详情
// @Description 根据建议ID获取详细的修改建议和说明
// @Tags AI写作辅助
// @Accept json
// @Produce json
// @Param id path string true "建议ID"
// @Success 200 {object} shared.APIResponse{data=dto.ProofreadSuggestion}
// @Failure 400 {object} shared.ErrorResponse
// @Failure 404 {object} shared.ErrorResponse
// @Failure 500 {object} shared.ErrorResponse
// @Router /api/v1/ai/writing/suggestions/{id} [get]
func (api *WritingAssistantApi) GetProofreadSuggestion(c *gin.Context) {
	suggestionID := c.Param("id")
	if suggestionID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "建议ID不能为空")
		return
	}

	// 调用服务
	result, err := api.proofreadService.GetProofreadSuggestion(c.Request.Context(), suggestionID)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "获取建议失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取建议成功", result)
}

// ===========================
// 敏感词检测 API
// ===========================

// CheckSensitiveWords 检测敏感词
// @Summary 检测敏感词
// @Description 检测文本中的敏感词，返回敏感词列表、位置和修改建议
// @Tags AI内容审核
// @Accept json
// @Produce json
// @Param request body dto.SensitiveWordsCheckRequest true "敏感词检测请求"
// @Success 200 {object} shared.APIResponse{data=dto.SensitiveWordsCheckResponse}
// @Failure 400 {object} shared.ErrorResponse
// @Failure 500 {object} shared.ErrorResponse
// @Router /api/v1/ai/audit/sensitive-words [post]
func (api *WritingAssistantApi) CheckSensitiveWords(c *gin.Context) {
	var req dto.SensitiveWordsCheckRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	// 生成请求ID
	requestID := uuid.New().String()
	c.Set("requestID", requestID)
	c.Set("aiService", "sensitive_words_check")

	// 调用服务
	result, err := api.sensitiveWordsService.CheckSensitiveWords(c.Request.Context(), &req)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "敏感词检测失败", err.Error())
		return
	}

	// 设置Token使用信息
	if result.TokensUsed > 0 {
		c.Set("tokensUsed", result.TokensUsed)
		c.Set("aiModel", "sensitive_words_detector")
	}

	shared.Success(c, http.StatusOK, "检测完成", result)
}

// GetSensitiveWordsDetail 获取敏感词检测结果
// @Summary 获取敏感词检测结果
// @Description 根据检测ID获取详细的敏感词检测结果
// @Tags AI内容审核
// @Accept json
// @Produce json
// @Param id path string true "检测ID"
// @Success 200 {object} shared.APIResponse{data=dto.SensitiveWordsDetail}
// @Failure 400 {object} shared.ErrorResponse
// @Failure 404 {object} shared.ErrorResponse
// @Failure 500 {object} shared.ErrorResponse
// @Router /api/v1/ai/audit/sensitive-words/{id} [get]
func (api *WritingAssistantApi) GetSensitiveWordsDetail(c *gin.Context) {
	checkID := c.Param("id")
	if checkID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "检测ID不能为空")
		return
	}

	// 调用服务
	result, err := api.sensitiveWordsService.GetSensitiveWordsDetail(c.Request.Context(), checkID)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "获取检测结果失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取检测结果成功", result)
}
