package admin

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/models/users"
	"Qingyu_backend/pkg/response"
	adminrepo "Qingyu_backend/repository/interfaces/admin"
	adminservice "Qingyu_backend/service/admin"
)

// UserAdminAPI 用户管理API
type UserAdminAPI struct {
	userAdminService adminservice.UserAdminService
}

// NewUserAdminAPI 创建用户管理API实例
func NewUserAdminAPI(userAdminService adminservice.UserAdminService) *UserAdminAPI {
	return &UserAdminAPI{
		userAdminService: userAdminService,
	}
}

// ListUsers 获取用户列表
//
//	@Summary		获取用户列表
//	@Description	管理员获取用户列表，支持分页和筛选
//	@Tags			Admin-User
//	@Accept			json
//	@Produce		json
//	@Param			keyword	query		string				false	"搜索关键词"
//	@Param			status		query		string				false	"状态筛选"
//	@Param			role		query		string				false	"角色筛选"
//	@Param			page		query		int					false	"页码"	default(1)
//	@Param			size		query		int					false	"每页数量"	default(20)
//	@Success		200			{object}	shared.APIResponse
//	@Failure		400			{object}	shared.APIResponse
//	@Failure		401			{object}	shared.APIResponse
//	@Failure		403			{object}	shared.APIResponse
//	@Router			/api/v1/admin/users [get]
func (api *UserAdminAPI) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	filter := &adminrepo.UserFilter{
		Keyword: c.Query("keyword"),
		Role:    c.Query("role"),
	}
	// Parse status if provided
	if statusStr := c.Query("status"); statusStr != "" {
		filter.Status = users.UserStatus(statusStr)
	}

	usersList, total, err := api.userAdminService.GetUserList(c.Request.Context(), filter, page, size)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	result := gin.H{
		"users": usersList,
		"total": total,
		"page":  page,
		"size":  size,
	}

	response.Success(c, result)
}

// GetUserDetail 获取用户详情
//
//	@Summary		获取用户详情
//	@Description	管理员获取用户详细信息
//	@Tags			Admin-User
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"用户ID"
//	@Success		200		{object}	shared.APIResponse
//	@Failure		400		{object}	shared.APIResponse
//	@Failure		401		{object}	shared.APIResponse
//	@Failure		403		{object}	shared.APIResponse
//	@Failure		404		{object}	shared.APIResponse
//	@Router			/api/v1/admin/users/{id} [get]
func (api *UserAdminAPI) GetUserDetail(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		response.BadRequest(c, "参数错误", "用户ID不能为空")
		return
	}

	user, err := api.userAdminService.GetUserDetail(c.Request.Context(), userID)
	if err != nil {
		if err == adminservice.ErrUserNotFound {
			response.NotFound(c, "用户不存在")
			return
		}
		response.InternalError(c, err)
		return
	}

	response.Success(c, user)
}

// UpdateUserStatus 更新用户状态
//
//	@Summary		更新用户状态
//	@Description	管理员更新用户状态
//	@Tags			Admin-User
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"用户ID"
//	@Param			request	body		UpdateUserStatusRequest	true	"状态信息"
//	@Success		200		{object}	shared.APIResponse
//	@Failure		400		{object}	shared.APIResponse
//	@Failure		401		{object}	shared.APIResponse
//	@Failure		403		{object}	shared.APIResponse
//	@Failure		404		{object}	shared.APIResponse
//	@Router			/api/v1/admin/users/{id}/status [put]
func (api *UserAdminAPI) UpdateUserStatus(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		response.BadRequest(c, "参数错误", "用户ID不能为空")
		return
	}

	var req UpdateUserStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	status := users.UserStatus(req.Status)
	if err := api.userAdminService.UpdateUserStatus(c.Request.Context(), userID, status); err != nil {
		if err == adminservice.ErrUserNotFound {
			response.NotFound(c, "用户不存在")
			return
		}
		if err == adminservice.ErrCannotModifySuperAdmin {
			response.Forbidden(c, "权限不足")
			return
		}
		response.InternalError(c, err)
		return
	}

	response.Success(c, nil)
}

// UpdateUserRole 更新用户角色
//
//	@Summary		更新用户角色
//	@Description	管理员更新用户角色
//	@Tags			Admin-User
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string					true	"用户ID"
//	@Param			request	body		UpdateUserRoleRequest	true	"角色信息"
//	@Success		200		{object}	shared.APIResponse
//	@Failure		400		{object}	shared.APIResponse
//	@Failure		401		{object}	shared.APIResponse
//	@Failure		403		{object}	shared.APIResponse
//	@Failure		404		{object}	shared.APIResponse
//	@Router			/api/v1/admin/users/{id}/role [put]
func (api *UserAdminAPI) UpdateUserRole(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		response.BadRequest(c, "参数错误", "用户ID不能为空")
		return
	}

	var req UpdateUserRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	if err := api.userAdminService.UpdateUserRole(c.Request.Context(), userID, req.Role); err != nil {
		if err == adminservice.ErrUserNotFound {
			response.NotFound(c, "用户不存在")
			return
		}
		if err == adminservice.ErrCannotModifySuperAdmin {
			response.Forbidden(c, "权限不足")
			return
		}
		if err == adminservice.ErrInvalidRole {
			response.BadRequest(c, "无效的角色", err.Error())
			return
		}
		response.InternalError(c, err)
		return
	}

	response.Success(c, nil)
}

// DeleteUser 删除用户
//
//	@Summary		删除用户
//	@Description	管理员删除用户
//	@Tags			Admin-User
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"用户ID"
//	@Success		200		{object}	shared.APIResponse
//	@Failure		400		{object}	shared.APIResponse
//	@Failure		401		{object}	shared.APIResponse
//	@Failure		403		{object}	shared.APIResponse
//	@Failure		404		{object}	shared.APIResponse
//	@Router			/api/v1/admin/users/{id} [delete]
func (api *UserAdminAPI) DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		response.BadRequest(c, "参数错误", "用户ID不能为空")
		return
	}

	if err := api.userAdminService.DeleteUser(c.Request.Context(), userID); err != nil {
		if err == adminservice.ErrUserNotFound {
			response.NotFound(c, "用户不存在")
			return
		}
		if err == adminservice.ErrCannotModifySuperAdmin {
			response.Forbidden(c, "权限不足")
			return
		}
		response.InternalError(c, err)
		return
	}

	response.Success(c, nil)
}

// GetUserActivities 获取用户活动记录
//
//	@Summary		获取用户活动记录
//	@Description	管理员获取用户活动记录
//	@Tags			Admin-User
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"用户ID"
//	@Param			page	query		int		false	"页码"	default(1)
//	@Param			size	query		int		false	"每页数量"	default(20)
//	@Success		200		{object}	shared.APIResponse
//	@Failure		400		{object}	shared.APIResponse
//	@Failure		401		{object}	shared.APIResponse
//	@Failure		404		{object}	shared.APIResponse
//	@Router			/api/v1/admin/users/{id}/activities [get]
func (api *UserAdminAPI) GetUserActivities(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		response.BadRequest(c, "参数错误", "用户ID不能为空")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	activities, total, err := api.userAdminService.GetUserActivities(c.Request.Context(), userID, page, size)
	if err != nil {
		if err == adminservice.ErrInvalidUserID {
			response.BadRequest(c, "无效的用户ID", err.Error())
			return
		}
		response.InternalError(c, err)
		return
	}

	result := gin.H{
		"activities": activities,
		"total":      total,
		"page":       page,
		"size":       size,
	}

	response.Success(c, result)
}

// GetUserStatistics 获取用户统计信息
//
//	@Summary		获取用户统计信息
//	@Description	管理员获取用户统计信息
//	@Tags			Admin-User
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"用户ID"
//	@Success		200		{object}	shared.APIResponse
//	@Failure		400		{object}	shared.APIResponse
//	@Failure		401		{object}	shared.APIResponse
//	@Failure		404		{object}	shared.APIResponse
//	@Router			/api/v1/admin/users/{id}/statistics [get]
func (api *UserAdminAPI) GetUserStatistics(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		response.BadRequest(c, "参数错误", "用户ID不能为空")
		return
	}

	stats, err := api.userAdminService.GetUserStatistics(c.Request.Context(), userID)
	if err != nil {
		if err == adminservice.ErrInvalidUserID {
			response.BadRequest(c, "无效的用户ID", err.Error())
			return
		}
		response.InternalError(c, err)
		return
	}

	response.Success(c, stats)
}

// ResetUserPassword 重置用户密码
//
//	@Summary		重置用户密码
//	@Description	管理员重置用户密码
//	@Tags			Admin-User
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"用户ID"
//	@Success		200		{object}	shared.APIResponse
//	@Failure		400		{object}	shared.APIResponse
//	@Failure		401		{object}	shared.APIResponse
//	@Failure		403		{object}	shared.APIResponse
//	@Failure		404		{object}	shared.APIResponse
//	@Router			/api/v1/admin/users/{id}/reset-password [post]
func (api *UserAdminAPI) ResetUserPassword(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		response.BadRequest(c, "参数错误", "用户ID不能为空")
		return
	}

	newPassword, err := api.userAdminService.ResetUserPassword(c.Request.Context(), userID)
	if err != nil {
		if err == adminservice.ErrUserNotFound {
			response.NotFound(c, "用户不存在")
			return
		}
		response.InternalError(c, err)
		return
	}

	result := gin.H{
		"newPassword": newPassword,
		"message":     "请将新密码安全地发送给用户",
	}

	response.Success(c, result)
}

// BatchUpdateStatus 批量更新用户状态
//
//	@Summary		批量更新用户状态
//	@Description	管理员批量更新用户状态
//	@Tags			Admin-User
//	@Accept			json
//	@Produce		json
//	@Param			request	body		BatchUpdateStatusRequest	true	"批量更新请求"
//	@Success		200		{object}	shared.APIResponse
//	@Failure		400		{object}	shared.APIResponse
//	@Failure		401		{object}	shared.APIResponse
//	@Failure		403		{object}	shared.APIResponse
//	@Router			/api/v1/admin/users/batch-update-status [post]
func (api *UserAdminAPI) BatchUpdateStatus(c *gin.Context) {
	var req BatchUpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	if len(req.UserIds) == 0 {
		response.BadRequest(c, "参数错误", "用户ID列表不能为空")
		return
	}

	status := users.UserStatus(req.Status)
	if err := api.userAdminService.BatchUpdateStatus(c.Request.Context(), req.UserIds, status); err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, nil)
}

// BatchDeleteUsers 批量删除用户
//
//	@Summary		批量删除用户
//	@Description	管理员批量删除用户
//	@Tags			Admin-User
//	@Accept			json
//	@Produce		json
//	@Param			request	body		BatchDeleteUsersRequest	true	"批量删除请求"
//	@Success		200		{object}	shared.APIResponse
//	@Failure		400		{object}	shared.APIResponse
//	@Failure		401		{object}	shared.APIResponse
//	@Failure		403		{object}	shared.APIResponse
//	@Router			/api/v1/admin/users/batch-delete [post]
func (api *UserAdminAPI) BatchDeleteUsers(c *gin.Context) {
	var req BatchDeleteUsersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	if len(req.UserIds) == 0 {
		response.BadRequest(c, "参数错误", "用户ID列表不能为空")
		return
	}

	if err := api.userAdminService.BatchDeleteUsers(c.Request.Context(), req.UserIds); err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, nil)
}

// SearchUsers 搜索用户
//
//	@Summary		搜索用户
//	@Description	管理员搜索用户
//	@Tags			Admin-User
//	@Accept			json
//	@Produce		json
//	@Param			keyword	query		string	true	"搜索关键词"
//	@Param			page		query		int		false	"页码"	default(1)
//	@Param			size		query		int		false	"每页数量"	default(20)
//	@Success		200		{object}	shared.APIResponse
//	@Failure		400		{object}	shared.APIResponse
//	@Failure		401		{object}	shared.APIResponse
//	@Router			/api/v1/admin/users/search [get]
func (api *UserAdminAPI) SearchUsers(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		response.BadRequest(c, "参数错误", "关键词不能为空")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	usersList, total, err := api.userAdminService.SearchUsers(c.Request.Context(), keyword, page, size)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	result := gin.H{
		"users": usersList,
		"total": total,
		"page":  page,
		"size":  size,
	}

	response.Success(c, result)
}

// CountByStatus 按状态统计用户数量
//
//	@Summary		按状态统计用户数量
//	@Description	管理员按状态统计用户数量
//	@Tags			Admin-User
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	shared.APIResponse
//	@Failure		401		{object}	shared.APIResponse
//	@Router			/api/v1/admin/users/count-by-status [get]
func (api *UserAdminAPI) CountByStatus(c *gin.Context) {
	counts, err := api.userAdminService.CountByStatus(c.Request.Context())
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, counts)
}

// ============ 请求体结构 ============

// UpdateUserStatusRequest 更新用户状态请求
type UpdateUserStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// UpdateUserRoleRequest 更新用户角色请求
type UpdateUserRoleRequest struct {
	Role string `json:"role" binding:"required"`
}

// BatchUpdateStatusRequest 批量更新状态请求
type BatchUpdateStatusRequest struct {
	UserIds []string `json:"userIds" binding:"required"`
	Status  string   `json:"status" binding:"required"`
}

// BatchDeleteUsersRequest 批量删除用户请求
type BatchDeleteUsersRequest struct {
	UserIds []string `json:"userIds" binding:"required"`
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	Role     string `json:"role" binding:"required,oneof=reader author admin"`
	Status   string `json:"status"`
	Bio      string `json:"bio"`
}

// BatchCreateUserRequest 批量创建用户请求
type BatchCreateUserRequest struct {
	Count  int    `json:"count" binding:"required,min=1,max=100"`
	Prefix string `json:"prefix"`
	Role   string `json:"role" binding:"required,oneof=reader author admin"`
	Status string `json:"status"`
}

// CreateUser 创建用户
//
//	@Summary		创建用户
//	@Description	管理员创建新用户
//	@Tags			Admin-User
//	@Accept			json
//	@Produce		json
//	@Param			request	body		CreateUserRequest	true	"用户信息"
//	@Success		200		{object}	shared.APIResponse
//	@Failure		400		{object}	shared.APIResponse
//	@Failure		401		{object}	shared.APIResponse
//	@Failure		403		{object}	shared.APIResponse
//	@Router			/api/v1/admin/users [post]
func (api *UserAdminAPI) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	serviceReq := &adminservice.CreateUserRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Nickname: req.Nickname,
		Role:     req.Role,
		Status:   users.UserStatus(req.Status),
		Bio:      req.Bio,
	}

	user, err := api.userAdminService.CreateUser(c.Request.Context(), serviceReq)
	if err != nil {
		if err == adminservice.ErrUserAlreadyExists {
			response.BadRequest(c, "用户已存在", err.Error())
			return
		}
		if err == adminservice.ErrInvalidRole {
			response.BadRequest(c, "无效的角色", err.Error())
			return
		}
		response.InternalError(c, err)
		return
	}

	response.Success(c, user)
}

// BatchCreateUsers 批量创建用户
//
//	@Summary		批量创建用户
//	@Description	管理员批量创建用户
//	@Tags			Admin-User
//	@Accept			json
//	@Produce		json
//	@Param			request	body		BatchCreateUserRequest	true	"批量创建请求"
//	@Success		200		{object}	shared.APIResponse
//	@Failure		400		{object}	shared.APIResponse
//	@Failure		401		{object}	shared.APIResponse
//	@Failure		403		{object}	shared.APIResponse
//	@Router			/api/v1/admin/users/batch-create [post]
func (api *UserAdminAPI) BatchCreateUsers(c *gin.Context) {
	var req BatchCreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	serviceReq := &adminservice.BatchCreateUserRequest{
		Count:  req.Count,
		Prefix: req.Prefix,
		Role:   req.Role,
		Status: users.UserStatus(req.Status),
	}

	usersList, err := api.userAdminService.BatchCreateUsers(c.Request.Context(), serviceReq)
	if err != nil {
		if err == adminservice.ErrInvalidRole {
			response.BadRequest(c, "无效的角色", err.Error())
			return
		}
		if err == adminservice.ErrInvalidBatchCount {
			response.BadRequest(c, "无效的批量创建数量", err.Error())
			return
		}
		response.InternalError(c, err)
		return
	}

	result := gin.H{
		"users": usersList,
		"count": len(usersList),
	}

	response.Success(c, result)
}
