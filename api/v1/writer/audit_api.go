package writer

import (
	"Qingyu_backend/service/interfaces/audit"

	"github.com/gin-gonic/gin"

	auditDTO "Qingyu_backend/service/audit"
	"Qingyu_backend/pkg/response"
)

// AuditApi 审核API
type AuditApi struct {
	auditService audit.ContentAuditService
}

// NewAuditApi 创建审核API
func NewAuditApi(auditService audit.ContentAuditService) *AuditApi {
	return &AuditApi{
		auditService: auditService,
	}
}

// CheckContent 实时检测内容
// @Summary 实时检测内容
// @Description 快速检测内容是否包含违规信息（不创建审核记录）
// @Tags 审核
// @Accept json
// @Produce json
// @Param request body object true "检测请求"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /api/v1/audit/check [post]
func (api *AuditApi) CheckContent(c *gin.Context) {
	var req auditDTO.CheckContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c,  "参数错误", err.Error())
		return
	}

	result, err := api.auditService.CheckContent(c.Request.Context(), req.Content)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, result)
}

// AuditDocument 全文审核文档
// @Summary 全文审核文档
// @Description 对文档进行全文审核并创建审核记录
// @Tags 审核
// @Accept json
// @Produce json
// @Param id path string true "文档ID"
// @Param request body object true "审核请求"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /api/v1/documents/{id}/audit [post]
func (api *AuditApi) AuditDocument(c *gin.Context) {
	documentID := c.Param("id")

	var req auditDTO.AuditDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c,  "参数错误", err.Error())
		return
	}

	// 从context获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "无法获取用户信息")
		return
	}

	req.DocumentID = documentID

	record, err := api.auditService.AuditDocument(c.Request.Context(), req.DocumentID, req.Content, userID.(string))
	if err != nil {
		c.Error(err)
		return
	}

	// 转换为响应DTO
	resp := convertAuditRecordToResponse(record)

	response.Success(c, resp)
}

// GetAuditResult 获取审核结果
// @Summary 获取审核结果
// @Description 根据文档ID获取审核结果
// @Tags 审核
// @Accept json
// @Produce json
// @Param id path string true "文档ID"
// @Param targetType query string false "目标类型" default(document)
// @Success 200 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /api/v1/documents/{id}/audit-result [get]
func (api *AuditApi) GetAuditResult(c *gin.Context) {
	documentID := c.Param("id")
	targetType := c.DefaultQuery("targetType", "document")

	record, err := api.auditService.GetAuditResult(c.Request.Context(), targetType, documentID)
	if err != nil {
			response.NotFound(c, err.Error())
		return
	}

	auditResp := convertAuditRecordToResponse(record)

	response.Success(c, auditResp)
}

// SubmitAppeal 提交申诉
// @Summary 提交申诉
// @Description 对审核结果提交申诉
// @Tags 审核
// @Accept json
// @Produce json
// @Param id path string true "审核记录ID"
// @Param request body object true "申诉请求"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /api/v1/audit/{id}/appeal [post]
func (api *AuditApi) SubmitAppeal(c *gin.Context) {
	auditID := c.Param("id")

	var req auditDTO.SubmitAppealRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c,  "参数错误", err.Error())
		return
	}

	// 从context获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "无法获取用户信息")
		return
	}

	err := api.auditService.SubmitAppeal(c.Request.Context(), auditID, userID.(string), req.Reason)
	if err != nil {
		response.BadRequest(c,  "提交申诉失败", err.Error())
		return
	}

	response.Success(c, nil)
}

// GetPendingReviews 获取待复核列表
// @Summary 获取待复核列表
// @Description 获取需要人工复核的审核记录（管理员接口）
// @Tags 审核-管理
// @Accept json
// @Produce json
// @Param limit query int false "数量限制" default(50)
// @Success 200 {object} response.APIResponse
// @Router /api/v1/admin/audit/pending [get]
func (api *AuditApi) GetPendingReviews(c *gin.Context) {
	limit := 50
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := c.GetQuery("limit"); err {
			limit = int(l[0])
		}
	}

	records, err := api.auditService.GetPendingReviews(c.Request.Context(), limit)
	if err != nil {
		c.Error(err)
		return
	}

	responses := make([]auditDTO.AuditRecordResponse, len(records))
	for i, record := range records {
		responses[i] = convertAuditRecordToResponse(record)
	}

	response.Success(c, responses)
}

// ReviewAudit 复核审核结果
// @Summary 复核审核结果
// @Description 人工复核审核结果（管理员接口）
// @Tags 审核-管理
// @Accept json
// @Produce json
// @Param id path string true "审核记录ID"
// @Param request body object true "复核请求"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/admin/audit/{id}/review [post]
func (api *AuditApi) ReviewAudit(c *gin.Context) {
	auditID := c.Param("id")

	var req auditDTO.ReviewAuditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c,  "参数错误", err.Error())
		return
	}

	// 从context获取管理员ID
	reviewerID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "无法获取管理员信息")
		return
	}

	err := api.auditService.ReviewAudit(c.Request.Context(), auditID, reviewerID.(string), req.Approved, req.Note)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, nil)
}

// ReviewAppeal 复核申诉
// @Summary 复核申诉
// @Description 人工复核申诉（管理员接口）
// @Tags 审核-管理
// @Accept json
// @Produce json
// @Param id path string true "审核记录ID"
// @Param request body object true "复核申诉请求"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/admin/audit/{id}/appeal/review [post]
func (api *AuditApi) ReviewAppeal(c *gin.Context) {
	auditID := c.Param("id")

	var req auditDTO.ReviewAppealRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c,  "参数错误", err.Error())
		return
	}

	// 从context获取管理员ID
	reviewerID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "无法获取管理员信息")
		return
	}

	err := api.auditService.ReviewAppeal(c.Request.Context(), auditID, reviewerID.(string), req.Approved, req.Note)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, nil)
}

// GetUserViolations 获取用户违规记录
// @Summary 获取用户违规记录
// @Description 获取指定用户的所有违规记录
// @Tags 审核
// @Accept json
// @Produce json
// @Param userId path string true "用户ID"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/users/{userId}/violations [get]
func (api *AuditApi) GetUserViolations(c *gin.Context) {
	userID := c.Param("userId")

	// 验证权限：只能查看自己的违规记录，或管理员可以查看所有
	currentUserID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "无法获取用户信息")
		return
	}

	// 检查是否为管理员或用户本身
	currentRole, roleExists := c.Get("role")
	isAdmin := roleExists && currentRole.(string) == "admin"

	if !isAdmin && currentUserID.(string) != userID {
		response.Forbidden(c, "只能查看自己的违规记录")
		return
	}

	violations, err := api.auditService.GetUserViolations(c.Request.Context(), userID)
	if err != nil {
		c.Error(err)
		return
	}

	responses := make([]auditDTO.ViolationRecordResponse, len(violations))
	for i, v := range violations {
		responses[i] = convertViolationRecordToResponse(v)
	}

	response.Success(c, responses)
}

// GetUserViolationSummary 获取用户违规统计
// @Summary 获取用户违规统计
// @Description 获取指定用户的违规统计信息
// @Tags 审核
// @Accept json
// @Produce json
// @Param userId path string true "用户ID"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/users/{userId}/violation-summary [get]
func (api *AuditApi) GetUserViolationSummary(c *gin.Context) {
	userID := c.Param("userId")

	// 验证权限
	currentUserID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "无法获取用户信息")
		return
	}

	// 检查是否为管理员或用户本身
	currentRole, roleExists := c.Get("role")
	isAdmin := roleExists && currentRole.(string) == "admin"

	if !isAdmin && currentUserID.(string) != userID {
		response.Forbidden(c, "只能查看自己的违规统计")
		return
	}

	summary, err := api.auditService.GetUserViolationSummary(c.Request.Context(), userID)
	if err != nil {
		c.Error(err)
		return
	}

	summaryResp := convertUserViolationSummaryToResponse(summary)

	response.Success(c, summaryResp)
}

// GetHighRiskAudits 获取高风险审核记录
// @Summary 获取高风险审核记录
// @Description 获取高风险审核记录（管理员接口）
// @Tags 审核-管理
// @Accept json
// @Produce json
// @Param minRiskLevel query int false "最低风险等级" default(3)
// @Param limit query int false "数量限制" default(50)
// @Success 200 {object} response.APIResponse
// @Router /api/v1/admin/audit/high-risk [get]
func (api *AuditApi) GetHighRiskAudits(c *gin.Context) {
	minRiskLevel := 3
	limit := 50

	if level := c.Query("minRiskLevel"); level != "" {
		if l, err := c.GetQuery("minRiskLevel"); err {
			minRiskLevel = int(l[0])
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := c.GetQuery("limit"); err {
			limit = int(l[0])
		}
	}

	records, err := api.auditService.GetHighRiskAudits(c.Request.Context(), minRiskLevel, limit)
	if err != nil {
		c.Error(err)
		return
	}

	responses := make([]auditDTO.AuditRecordResponse, len(records))
	for i, record := range records {
		responses[i] = convertAuditRecordToResponse(record)
	}

	response.Success(c, responses)
}

// 辅助转换函数

func convertAuditRecordToResponse(record interface{}) auditDTO.AuditRecordResponse {
	// TODO: 实现完整的转换逻辑
	// 这里简化处理，实际应该使用类型断言或反射
	return auditDTO.AuditRecordResponse{}
}

func convertViolationRecordToResponse(violation interface{}) auditDTO.ViolationRecordResponse {
	// TODO: 实现完整的转换逻辑
	return auditDTO.ViolationRecordResponse{}
}

func convertUserViolationSummaryToResponse(summary interface{}) auditDTO.UserViolationSummaryResponse {
	// TODO: 实现完整的转换逻辑
	return auditDTO.UserViolationSummaryResponse{}
}
