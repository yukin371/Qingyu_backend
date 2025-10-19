package writer

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/shared"
	auditDTO "Qingyu_backend/service/audit"
	"Qingyu_backend/service/interfaces"
)

// AuditApi 审核API
type AuditApi struct {
	auditService interfaces.ContentAuditService
}

// NewAuditApi 创建审核API
func NewAuditApi(auditService interfaces.ContentAuditService) *AuditApi {
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
// @Param request body auditDTO.CheckContentRequest true "检测请求"
// @Success 200 {object} shared.Response{data=interfaces.AuditCheckResult}
// @Failure 400 {object} shared.Response
// @Router /api/v1/audit/check [post]
func (api *AuditApi) CheckContent(c *gin.Context) {
	var req auditDTO.CheckContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	result, err := api.auditService.CheckContent(c.Request.Context(), req.Content)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "检测失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "检测完成", result)
}

// AuditDocument 全文审核文档
// @Summary 全文审核文档
// @Description 对文档进行全文审核并创建审核记录
// @Tags 审核
// @Accept json
// @Produce json
// @Param id path string true "文档ID"
// @Param request body auditDTO.AuditDocumentRequest true "审核请求"
// @Success 200 {object} shared.Response{data=auditDTO.AuditRecordResponse}
// @Failure 400 {object} shared.Response
// @Router /api/v1/documents/{id}/audit [post]
func (api *AuditApi) AuditDocument(c *gin.Context) {
	documentID := c.Param("id")

	var req auditDTO.AuditDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	// 从context获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "无法获取用户信息")
		return
	}

	req.DocumentID = documentID

	record, err := api.auditService.AuditDocument(c.Request.Context(), req.DocumentID, req.Content, userID.(string))
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "审核失败", err.Error())
		return
	}

	// 转换为响应DTO
	response := convertAuditRecordToResponse(record)

	shared.Success(c, http.StatusOK, "审核完成", response)
}

// GetAuditResult 获取审核结果
// @Summary 获取审核结果
// @Description 根据文档ID获取审核结果
// @Tags 审核
// @Accept json
// @Produce json
// @Param id path string true "文档ID"
// @Param targetType query string false "目标类型" default(document)
// @Success 200 {object} shared.Response{data=auditDTO.AuditRecordResponse}
// @Failure 404 {object} shared.Response
// @Router /api/v1/documents/{id}/audit-result [get]
func (api *AuditApi) GetAuditResult(c *gin.Context) {
	documentID := c.Param("id")
	targetType := c.DefaultQuery("targetType", "document")

	record, err := api.auditService.GetAuditResult(c.Request.Context(), targetType, documentID)
	if err != nil {
		shared.Error(c, http.StatusNotFound, "未找到审核记录", err.Error())
		return
	}

	response := convertAuditRecordToResponse(record)

	shared.Success(c, http.StatusOK, "获取成功", response)
}

// SubmitAppeal 提交申诉
// @Summary 提交申诉
// @Description 对审核结果提交申诉
// @Tags 审核
// @Accept json
// @Produce json
// @Param id path string true "审核记录ID"
// @Param request body auditDTO.SubmitAppealRequest true "申诉请求"
// @Success 200 {object} shared.Response
// @Failure 400 {object} shared.Response
// @Router /api/v1/audit/{id}/appeal [post]
func (api *AuditApi) SubmitAppeal(c *gin.Context) {
	auditID := c.Param("id")

	var req auditDTO.SubmitAppealRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	// 从context获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "无法获取用户信息")
		return
	}

	err := api.auditService.SubmitAppeal(c.Request.Context(), auditID, userID.(string), req.Reason)
	if err != nil {
		shared.Error(c, http.StatusBadRequest, "提交申诉失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "申诉已提交，等待复核", nil)
}

// GetPendingReviews 获取待复核列表
// @Summary 获取待复核列表
// @Description 获取需要人工复核的审核记录（管理员接口）
// @Tags 审核-管理
// @Accept json
// @Produce json
// @Param limit query int false "数量限制" default(50)
// @Success 200 {object} shared.Response{data=[]auditDTO.AuditRecordResponse}
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
		shared.Error(c, http.StatusInternalServerError, "查询失败", err.Error())
		return
	}

	responses := make([]auditDTO.AuditRecordResponse, len(records))
	for i, record := range records {
		responses[i] = convertAuditRecordToResponse(record)
	}

	shared.Success(c, http.StatusOK, "获取成功", responses)
}

// ReviewAudit 复核审核结果
// @Summary 复核审核结果
// @Description 人工复核审核结果（管理员接口）
// @Tags 审核-管理
// @Accept json
// @Produce json
// @Param id path string true "审核记录ID"
// @Param request body auditDTO.ReviewAuditRequest true "复核请求"
// @Success 200 {object} shared.Response
// @Router /api/v1/admin/audit/{id}/review [post]
func (api *AuditApi) ReviewAudit(c *gin.Context) {
	auditID := c.Param("id")

	var req auditDTO.ReviewAuditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	// 从context获取管理员ID
	reviewerID, exists := c.Get("userID")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "无法获取管理员信息")
		return
	}

	err := api.auditService.ReviewAudit(c.Request.Context(), auditID, reviewerID.(string), req.Approved, req.Note)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "复核失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "复核完成", nil)
}

// ReviewAppeal 复核申诉
// @Summary 复核申诉
// @Description 人工复核申诉（管理员接口）
// @Tags 审核-管理
// @Accept json
// @Produce json
// @Param id path string true "审核记录ID"
// @Param request body auditDTO.ReviewAppealRequest true "复核申诉请求"
// @Success 200 {object} shared.Response
// @Router /api/v1/admin/audit/{id}/appeal/review [post]
func (api *AuditApi) ReviewAppeal(c *gin.Context) {
	auditID := c.Param("id")

	var req auditDTO.ReviewAppealRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	// 从context获取管理员ID
	reviewerID, exists := c.Get("userID")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "无法获取管理员信息")
		return
	}

	err := api.auditService.ReviewAppeal(c.Request.Context(), auditID, reviewerID.(string), req.Approved, req.Note)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "复核申诉失败", err.Error())
		return
	}

	status := "驳回"
	if req.Approved {
		status = "通过"
	}

	shared.Success(c, http.StatusOK, "申诉已"+status, nil)
}

// GetUserViolations 获取用户违规记录
// @Summary 获取用户违规记录
// @Description 获取指定用户的所有违规记录
// @Tags 审核
// @Accept json
// @Produce json
// @Param userId path string true "用户ID"
// @Success 200 {object} shared.Response{data=[]auditDTO.ViolationRecordResponse}
// @Router /api/v1/users/{userId}/violations [get]
func (api *AuditApi) GetUserViolations(c *gin.Context) {
	userID := c.Param("userId")

	// 验证权限：只能查看自己的违规记录，或管理员可以查看所有
	currentUserID, exists := c.Get("userID")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "无法获取用户信息")
		return
	}

	// TODO: 检查是否为管理员
	if currentUserID.(string) != userID {
		shared.Error(c, http.StatusForbidden, "无权限", "只能查看自己的违规记录")
		return
	}

	violations, err := api.auditService.GetUserViolations(c.Request.Context(), userID)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "查询失败", err.Error())
		return
	}

	responses := make([]auditDTO.ViolationRecordResponse, len(violations))
	for i, v := range violations {
		responses[i] = convertViolationRecordToResponse(v)
	}

	shared.Success(c, http.StatusOK, "获取成功", responses)
}

// GetUserViolationSummary 获取用户违规统计
// @Summary 获取用户违规统计
// @Description 获取指定用户的违规统计信息
// @Tags 审核
// @Accept json
// @Produce json
// @Param userId path string true "用户ID"
// @Success 200 {object} shared.Response{data=auditDTO.UserViolationSummaryResponse}
// @Router /api/v1/users/{userId}/violation-summary [get]
func (api *AuditApi) GetUserViolationSummary(c *gin.Context) {
	userID := c.Param("userId")

	// 验证权限
	currentUserID, exists := c.Get("userID")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "无法获取用户信息")
		return
	}

	// TODO: 检查是否为管理员
	if currentUserID.(string) != userID {
		shared.Error(c, http.StatusForbidden, "无权限", "只能查看自己的违规统计")
		return
	}

	summary, err := api.auditService.GetUserViolationSummary(c.Request.Context(), userID)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "查询失败", err.Error())
		return
	}

	response := convertUserViolationSummaryToResponse(summary)

	shared.Success(c, http.StatusOK, "获取成功", response)
}

// GetHighRiskAudits 获取高风险审核记录
// @Summary 获取高风险审核记录
// @Description 获取高风险审核记录（管理员接口）
// @Tags 审核-管理
// @Accept json
// @Produce json
// @Param minRiskLevel query int false "最低风险等级" default(3)
// @Param limit query int false "数量限制" default(50)
// @Success 200 {object} shared.Response{data=[]auditDTO.AuditRecordResponse}
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
		shared.Error(c, http.StatusInternalServerError, "查询失败", err.Error())
		return
	}

	responses := make([]auditDTO.AuditRecordResponse, len(records))
	for i, record := range records {
		responses[i] = convertAuditRecordToResponse(record)
	}

	shared.Success(c, http.StatusOK, "获取成功", responses)
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
