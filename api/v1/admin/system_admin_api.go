package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/shared"
	adminService "Qingyu_backend/service/shared/admin"
)

// SystemAdminAPI 系统管理API处理器（管理员）
type SystemAdminAPI struct {
	adminService adminService.AdminService
}

// NewSystemAdminAPI 创建系统管理API实例
func NewSystemAdminAPI(adminSvc adminService.AdminService) *SystemAdminAPI {
	return &SystemAdminAPI{
		adminService: adminSvc,
	}
}

// ReviewWithdraw 审核提现（管理员）
//
//	@Summary		审核提现
//	@Description	管理员审核用户提现申请（批准或拒绝）
//	@Tags			管理员-系统管理
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			request	body		ReviewWithdrawRequest	true	"审核信息"
//	@Success		200		{object}	shared.APIResponse
//	@Failure		400		{object}	shared.ErrorResponse
//	@Failure		401		{object}	shared.ErrorResponse
//	@Failure		403		{object}	shared.ErrorResponse
//	@Failure		500		{object}	shared.ErrorResponse
//	@Router			/api/v1/admin/withdraw/review [post]
func (api *SystemAdminAPI) ReviewWithdraw(c *gin.Context) {
	var req ReviewWithdrawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	// 从context获取管理员ID
	adminID, exists := c.Get("user_id")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "无法获取管理员信息")
		return
	}

	// 调用Service层
	err := api.adminService.ReviewWithdraw(c.Request.Context(), req.WithdrawID, adminID.(string), req.Approved, req.Reason)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "审核提现失败", err.Error())
		return
	}

	message := "拒绝提现成功"
	if req.Approved {
		message = "批准提现成功"
	}

	shared.Success(c, http.StatusOK, message, nil)
}

// GetUserStatistics 获取用户统计（管理员）
//
//	@Summary		获取用户统计
//	@Description	管理员获取指定用户的统计信息
//	@Tags			管理员-系统管理
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		string	true	"用户ID"
//	@Success		200		{object}	shared.APIResponse
//	@Failure		400		{object}	shared.ErrorResponse
//	@Failure		401		{object}	shared.ErrorResponse
//	@Failure		403		{object}	shared.ErrorResponse
//	@Failure		500		{object}	shared.ErrorResponse
//	@Router			/api/v1/admin/users/{id}/statistics [get]
func (api *SystemAdminAPI) GetUserStatistics(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "缺少用户ID")
		return
	}

	// 调用Service层
	stats, err := api.adminService.GetUserStatistics(c.Request.Context(), userID)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "获取用户统计失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取用户统计成功", stats)
}

// GetOperationLogs 获取操作日志（管理员）
//
//	@Summary		获取操作日志
//	@Description	管理员获取操作日志
//	@Tags			管理员-系统管理
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			page		query		int		false	"页码"			default(1)
//	@Param			page_size	query		int		false	"每页数量"		default(20)
//	@Param			admin_id	query		string	false	"管理员ID"
//	@Param			operation	query		string	false	"操作类型"
//	@Success		200			{object}	shared.PaginatedResponse
//	@Failure		401			{object}	shared.ErrorResponse
//	@Failure		403			{object}	shared.ErrorResponse
//	@Failure		500			{object}	shared.ErrorResponse
//	@Router			/api/v1/admin/operation-logs [get]
func (api *SystemAdminAPI) GetOperationLogs(c *gin.Context) {
	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	adminID := c.Query("admin_id")
	operation := c.Query("operation")

	// 构建请求
	req := &adminService.GetLogsRequest{
		AdminID:   adminID,
		Operation: operation,
		Page:      page,
		PageSize:  pageSize,
	}

	// 调用Service层
	logs, err := api.adminService.GetOperationLogs(c.Request.Context(), req)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "获取操作日志失败", err.Error())
		return
	}

	shared.Paginated(c, logs, int64(len(logs)), page, pageSize, "获取操作日志成功")
}

// GetSystemStats 获取系统统计（管理员）
//
//	@Summary		获取系统统计
//	@Description	管理员获取系统整体统计数据
//	@Tags			管理员-系统管理
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success		200	{object}	shared.APIResponse{data=SystemStatsResponse}
//	@Failure		401	{object}	shared.ErrorResponse
//	@Failure		403	{object}	shared.ErrorResponse
//	@Failure		500	{object}	shared.ErrorResponse
//	@Router			/api/v1/admin/stats [get]
func (api *SystemAdminAPI) GetSystemStats(c *gin.Context) {
	// TODO: 实现系统统计功能
	// 统计内容包括：
	// - 总用户数
	// - 活跃用户数
	// - 总书籍数
	// - 总收入
	// - 待审核数量

	stats := SystemStatsResponse{
		TotalUsers:    0,
		ActiveUsers:   0,
		TotalBooks:    0,
		TotalRevenue:  0.0,
		PendingAudits: 0,
	}

	shared.Success(c, http.StatusOK, "获取成功", stats)
}

// GetSystemConfig 获取系统配置（管理员）
//
//	@Summary		获取系统配置
//	@Description	管理员获取系统配置信息
//	@Tags			管理员-系统管理
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success		200	{object}	shared.APIResponse
//	@Failure		401	{object}	shared.ErrorResponse
//	@Failure		403	{object}	shared.ErrorResponse
//	@Failure		500	{object}	shared.ErrorResponse
//	@Router			/api/v1/admin/config [get]
func (api *SystemAdminAPI) GetSystemConfig(c *gin.Context) {
	// TODO: 实现获取系统配置功能
	// 配置内容包括：
	// - 是否允许注册
	// - 是否需要邮箱验证
	// - 最大上传文件大小
	// - 是否启用审核

	config := map[string]interface{}{
		"allowRegistration":        true,
		"requireEmailVerification": true,
		"maxUploadSize":            10485760,
		"enableAudit":              true,
	}

	shared.Success(c, http.StatusOK, "获取成功", config)
}

// UpdateSystemConfig 更新系统配置（管理员）
//
//	@Summary		更新系统配置
//	@Description	管理员更新系统配置
//	@Tags			管理员-系统管理
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			request	body		SystemConfigRequest	true	"配置信息"
//	@Success		200		{object}	shared.APIResponse
//	@Failure		400		{object}	shared.ErrorResponse
//	@Failure		401		{object}	shared.ErrorResponse
//	@Failure		403		{object}	shared.ErrorResponse
//	@Failure		500		{object}	shared.ErrorResponse
//	@Router			/api/v1/admin/config [put]
func (api *SystemAdminAPI) UpdateSystemConfig(c *gin.Context) {
	var req SystemConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	// TODO: 实现更新系统配置功能
	// 这里需要实现配置持久化和生效逻辑

	shared.Success(c, http.StatusOK, "更新成功", nil)
}

// CreateAnnouncement 发布公告（管理员）
//
//	@Summary		发布公告
//	@Description	管理员发布系统公告
//	@Tags			管理员-系统管理
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			request	body		object	true	"公告信息"
//	@Success		200		{object}	shared.APIResponse
//	@Failure		400		{object}	shared.ErrorResponse
//	@Failure		401		{object}	shared.ErrorResponse
//	@Failure		403		{object}	shared.ErrorResponse
//	@Failure		500		{object}	shared.ErrorResponse
//	@Router			/api/v1/admin/announcements [post]
func (api *SystemAdminAPI) CreateAnnouncement(c *gin.Context) {
	var req struct {
		Title    string `json:"title" binding:"required"`
		Content  string `json:"content" binding:"required"`
		Type     string `json:"type" binding:"required,oneof=system maintenance event"`
		Priority string `json:"priority" binding:"required,oneof=low normal high urgent"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	// TODO: 实现发布公告功能
	// 公告应该存储到数据库，并通知所有在线用户

	shared.Success(c, http.StatusCreated, "公告发布成功", nil)
}

// GetAnnouncements 获取公告列表（管理员）
//
//	@Summary		获取公告列表
//	@Description	管理员获取公告列表
//	@Tags			管理员-系统管理
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			page		query		int	false	"页码"		default(1)
//	@Param			page_size	query		int	false	"每页数量"	default(20)
//	@Success		200			{object}	shared.PaginatedResponse
//	@Failure		401			{object}	shared.ErrorResponse
//	@Failure		403			{object}	shared.ErrorResponse
//	@Failure		500			{object}	shared.ErrorResponse
//	@Router			/api/v1/admin/announcements [get]
func (api *SystemAdminAPI) GetAnnouncements(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	// TODO: 实现获取公告列表功能

	announcements := []interface{}{}
	shared.Paginated(c, announcements, 0, page, pageSize, "获取成功")
}
