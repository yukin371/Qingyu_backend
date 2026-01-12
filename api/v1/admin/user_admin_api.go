package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/shared"
	"Qingyu_backend/models/users"
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
		shared.Error(c, http.StatusInternalServerError, "获取失败", err.Error())
		return
	}

	result := gin.H{
		"users": usersList,
		"total": total,
		"page":  page,
		"size":  size,
	}

	shared.Success(c, http.StatusOK, "获取成功", result)
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
		shared.Error(c, http.StatusBadRequest, "参数错误", "用户ID不能为空")
		return
	}

	user, err := api.userAdminService.GetUserDetail(c.Request.Context(), userID)
	if err != nil {
		if err == adminservice.ErrUserNotFound {
			shared.Error(c, http.StatusNotFound, "用户不存在", err.Error())
			return
		}
		shared.Error(c, http.StatusInternalServerError, "获取失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取成功", user)
}

// UpdateUserStatus 更新用户状态
//
//	@Summary		更新用户状态
//	@Description	管理员更新用户状态
//	@Tags			Admin-User
//	@Accept			json
//	@Produce		json
//	@Param			_id		path		string					true	"用户ID"
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
		shared.Error(c, http.StatusBadRequest, "参数错误", "用户ID不能为空")
		return
	}

	var req UpdateUserStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	status := users.UserStatus(req.Status)
	if err := api.userAdminService.UpdateUserStatus(c.Request.Context(), userID, status); err != nil {
		if err == adminservice.ErrUserNotFound {
			shared.Error(c, http.StatusNotFound, "用户不存在", err.Error())
			return
		}
		if err == adminservice.ErrCannotModifySuperAdmin {
			shared.Error(c, http.StatusForbidden, "权限不足", err.Error())
			return
		}
		shared.Error(c, http.StatusInternalServerError, "更新失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "更新成功", nil)
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
		shared.Error(c, http.StatusBadRequest, "参数错误", "用户ID不能为空")
		return
	}

	var req UpdateUserRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	if err := api.userAdminService.UpdateUserRole(c.Request.Context(), userID, req.Role); err != nil {
		if err == adminservice.ErrUserNotFound {
			shared.Error(c, http.StatusNotFound, "用户不存在", err.Error())
			return
		}
		if err == adminservice.ErrCannotModifySuperAdmin {
			shared.Error(c, http.StatusForbidden, "权限不足", err.Error())
			return
		}
		if err == adminservice.ErrInvalidRole {
			shared.Error(c, http.StatusBadRequest, "无效的角色", err.Error())
			return
		}
		shared.Error(c, http.StatusInternalServerError, "更新失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "更新成功", nil)
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
		shared.Error(c, http.StatusBadRequest, "参数错误", "用户ID不能为空")
		return
	}

	if err := api.userAdminService.DeleteUser(c.Request.Context(), userID); err != nil {
		if err == adminservice.ErrUserNotFound {
			shared.Error(c, http.StatusNotFound, "用户不存在", err.Error())
			return
		}
		if err == adminservice.ErrCannotModifySuperAdmin {
			shared.Error(c, http.StatusForbidden, "权限不足", err.Error())
			return
		}
		shared.Error(c, http.StatusInternalServerError, "删除失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "删除成功", nil)
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
		shared.Error(c, http.StatusBadRequest, "参数错误", "用户ID不能为空")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	activities, total, err := api.userAdminService.GetUserActivities(c.Request.Context(), userID, page, size)
	if err != nil {
		if err == adminservice.ErrInvalidUserID {
			shared.Error(c, http.StatusBadRequest, "无效的用户ID", err.Error())
			return
		}
		shared.Error(c, http.StatusInternalServerError, "获取失败", err.Error())
		return
	}

	result := gin.H{
		"activities": activities,
		"total":      total,
		"page":       page,
		"size":       size,
	}

	shared.Success(c, http.StatusOK, "获取成功", result)
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
		shared.Error(c, http.StatusBadRequest, "参数错误", "用户ID不能为空")
		return
	}

	stats, err := api.userAdminService.GetUserStatistics(c.Request.Context(), userID)
	if err != nil {
		if err == adminservice.ErrInvalidUserID {
			shared.Error(c, http.StatusBadRequest, "无效的用户ID", err.Error())
			return
		}
		shared.Error(c, http.StatusInternalServerError, "获取失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取成功", stats)
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
		shared.Error(c, http.StatusBadRequest, "参数错误", "用户ID不能为空")
		return
	}

	newPassword, err := api.userAdminService.ResetUserPassword(c.Request.Context(), userID)
	if err != nil {
		if err == adminservice.ErrUserNotFound {
			shared.Error(c, http.StatusNotFound, "用户不存在", err.Error())
			return
		}
		shared.Error(c, http.StatusInternalServerError, "重置失败", err.Error())
		return
	}

	result := gin.H{
		"newPassword": newPassword,
		"message":     "请将新密码安全地发送给用户",
	}

	shared.Success(c, http.StatusOK, "重置成功", result)
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
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	if len(req.UserIds) == 0 {
		shared.Error(c, http.StatusBadRequest, "参数错误", "用户ID列表不能为空")
		return
	}

	status := users.UserStatus(req.Status)
	if err := api.userAdminService.BatchUpdateStatus(c.Request.Context(), req.UserIds, status); err != nil {
		shared.Error(c, http.StatusInternalServerError, "批量更新失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "批量更新成功", nil)
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
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	if len(req.UserIds) == 0 {
		shared.Error(c, http.StatusBadRequest, "参数错误", "用户ID列表不能为空")
		return
	}

	if err := api.userAdminService.BatchDeleteUsers(c.Request.Context(), req.UserIds); err != nil {
		shared.Error(c, http.StatusInternalServerError, "批量删除失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "批量删除成功", nil)
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
		shared.Error(c, http.StatusBadRequest, "参数错误", "关键词不能为空")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	usersList, total, err := api.userAdminService.SearchUsers(c.Request.Context(), keyword, page, size)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "搜索失败", err.Error())
		return
	}

	result := gin.H{
		"users": usersList,
		"total": total,
		"page":  page,
		"size":  size,
	}

	shared.Success(c, http.StatusOK, "搜索成功", result)
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
		shared.Error(c, http.StatusInternalServerError, "统计失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "统计成功", counts)
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
