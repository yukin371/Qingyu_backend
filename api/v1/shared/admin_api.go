package shared

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/service/shared/admin"
)

// AdminAPI 管理服务API处理器
type AdminAPI struct {
	adminService admin.AdminService
}

// NewAdminAPI 创建管理API实例
func NewAdminAPI(adminService admin.AdminService) *AdminAPI {
	return &AdminAPI{
		adminService: adminService,
	}
}

// GetPendingReviews 获取待审核内容
//
//	@Summary		获取待审核内容
//	@Description	获取待审核内容列表
//	@Tags			管理
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			content_type	query		string	false	"内容类型"
//	@Success		200				{object}	APIResponse
//	@Failure		401				{object}	APIResponse
//	@Failure		403				{object}	APIResponse
//	@Failure		500				{object}	APIResponse
//	@Router			/api/v1/shared/admin/reviews/pending [get]
func (api *AdminAPI) GetPendingReviews(c *gin.Context) {
	// 检查管理员权限
	role, exists := c.Get("user_role")
	if !exists || role != "admin" {
		c.JSON(http.StatusForbidden, APIResponse{
			Code:    403,
			Message: "无权限访问",
		})
		return
	}

	contentType := c.Query("content_type")

	reviews, err := api.adminService.GetPendingReviews(c.Request.Context(), contentType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "获取待审核内容失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "获取待审核内容成功",
		Data:    reviews,
	})
}

// ReviewContentRequest 审核内容请求
type ReviewContentRequest struct {
	ContentID   string `json:"content_id" binding:"required,min=1"`
	ContentType string `json:"content_type" binding:"required" validate:"content_type"`
	Action      string `json:"action" binding:"required,oneof=approve reject"`
	Reason      string `json:"reason" validate:"omitempty,max=500"`
}

// ReviewContent 审核内容
//
//	@Summary		审核内容
//	@Description	审核用户提交的内容
//	@Tags			管理
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			request	body		ReviewContentRequest	true	"审核信息"
//	@Success		200		{object}	APIResponse
//	@Failure		400		{object}	APIResponse
//	@Failure		401		{object}	APIResponse
//	@Failure		403		{object}	APIResponse
//	@Failure		500		{object}	APIResponse
//	@Router			/api/v1/shared/admin/reviews [post]
func (api *AdminAPI) ReviewContent(c *gin.Context) {
	// 检查管理员权限
	role, exists := c.Get("user_role")
	if !exists || role != "admin" {
		c.JSON(http.StatusForbidden, APIResponse{
			Code:    403,
			Message: "无权限访问",
		})
		return
	}

	adminID, _ := c.Get("user_id")

	var reqBody ReviewContentRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	req := &admin.ReviewContentRequest{
		ContentID:   reqBody.ContentID,
		ContentType: reqBody.ContentType,
		Action:      reqBody.Action,
		Reason:      reqBody.Reason,
		ReviewerID:  adminID.(string),
	}

	err := api.adminService.ReviewContent(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "审核失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "审核成功",
	})
}

// ReviewWithdrawRequest 审核提现请求
type ReviewWithdrawRequest struct {
	WithdrawID string `json:"withdraw_id" binding:"required"`
	Approved   bool   `json:"approved"`
	Reason     string `json:"reason"`
}

// ReviewWithdraw 审核提现
//
//	@Summary		审核提现
//	@Description	审核用户提现申请（批准或拒绝）
//	@Tags			管理
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			request	body		ReviewWithdrawRequest	true	"审核信息"
//	@Success		200		{object}	APIResponse
//	@Failure		400		{object}	APIResponse
//	@Failure		401		{object}	APIResponse
//	@Failure		403		{object}	APIResponse
//	@Failure		500		{object}	APIResponse
//	@Router			/api/v1/shared/admin/withdraw/review [post]
func (api *AdminAPI) ReviewWithdraw(c *gin.Context) {
	// 检查管理员权限
	role, exists := c.Get("user_role")
	if !exists || role != "admin" {
		c.JSON(http.StatusForbidden, APIResponse{
			Code:    403,
			Message: "无权限访问",
		})
		return
	}

	adminID, _ := c.Get("user_id")

	var req ReviewWithdrawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	err := api.adminService.ReviewWithdraw(c.Request.Context(), req.WithdrawID, adminID.(string), req.Approved, req.Reason)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "审核提现失败: " + err.Error(),
		})
		return
	}

	message := "拒绝提现成功"
	if req.Approved {
		message = "批准提现成功"
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: message,
	})
}

// GetUserStatistics 获取用户统计
//
//	@Summary		获取用户统计
//	@Description	获取指定用户的统计信息
//	@Tags			管理
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			user_id	path		string	true	"用户ID"
//	@Success 200 {object} APIResponse
//	@Failure		401		{object}	APIResponse
//	@Failure		403		{object}	APIResponse
//	@Failure		500		{object}	APIResponse
//	@Router			/api/v1/shared/admin/users/{user_id}/statistics [get]
func (api *AdminAPI) GetUserStatistics(c *gin.Context) {
	// 检查管理员权限
	role, exists := c.Get("user_role")
	if !exists || role != "admin" {
		c.JSON(http.StatusForbidden, APIResponse{
			Code:    403,
			Message: "无权限访问",
		})
		return
	}

	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "缺少用户ID",
		})
		return
	}

	stats, err := api.adminService.GetUserStatistics(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "获取用户统计失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "获取用户统计成功",
		Data:    stats,
	})
}

// GetOperationLogs 获取操作日志
//
//	@Summary		获取操作日志
//	@Description	获取管理员操作日志
//	@Tags			管理
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			page		query		int		false	"页码"	default(1)
//	@Param			page_size	query		int		false	"每页数量"	default(20)
//	@Param			admin_id	query		string	false	"管理员ID"
//	@Param			operation	query		string	false	"操作类型"
//	@Success		200			{object}	PaginatedResponse
//	@Failure		401			{object}	APIResponse
//	@Failure		403			{object}	APIResponse
//	@Failure		500			{object}	APIResponse
//	@Router			/api/v1/shared/admin/operation-logs [get]
func (api *AdminAPI) GetOperationLogs(c *gin.Context) {
	// 检查管理员权限
	role, exists := c.Get("user_role")
	if !exists || role != "admin" {
		c.JSON(http.StatusForbidden, APIResponse{
			Code:    403,
			Message: "无权限访问",
		})
		return
	}

	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	adminID := c.Query("admin_id")
	operation := c.Query("operation")

	req := &admin.GetLogsRequest{
		AdminID:   adminID,
		Operation: operation,
		Page:      page,
		PageSize:  pageSize,
	}

	logs, err := api.adminService.GetOperationLogs(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "获取操作日志失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, PaginatedResponseHelper(
		logs,
		int64(len(logs)),
		page,
		pageSize,
		"获取操作日志成功",
	))
}
