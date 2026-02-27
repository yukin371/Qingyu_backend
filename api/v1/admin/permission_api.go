package admin

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/shared"
	"Qingyu_backend/models/auth"
	"Qingyu_backend/pkg/response"
	sharedService "Qingyu_backend/service/shared"
)

// PermissionAPI 权限管理API
type PermissionAPI struct {
	permissionService sharedService.PermissionService
}

// NewPermissionAPI 创建权限管理API实例
func NewPermissionAPI(permissionService sharedService.PermissionService) *PermissionAPI {
	return &PermissionAPI{
		permissionService: permissionService,
	}
}

// ==================== 权限管理 ====================

// GetAllPermissions 获取所有权限
//
//	@Summary		获取所有权限
//	@Description	管理员获取所有权限列表
//	@Tags			Admin-Permission
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	shared.APIResponse
//	@Failure		401	{object}	shared.APIResponse
//	@Failure		403	{object}	shared.APIResponse
//	@Router			/api/v1/admin/permissions [get]
func (api *PermissionAPI) GetAllPermissions(c *gin.Context) {
	permissions, err := api.permissionService.GetAllPermissions(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, permissions)
}

// GetPermission 获取权限详情
//
//	@Summary		获取权限详情
//	@Description	管理员获取权限详细信息
//	@Tags			Admin-Permission
//	@Accept			json
//	@Produce		json
//	@Param			code	path	string	true	"权限代码"
//	@Success		200	{object}	shared.APIResponse
//	@Failure		400	{object}	shared.APIResponse
//	@Failure		404	{object}	shared.APIResponse
//	@Router			/api/v1/admin/permissions/{code} [get]
func (api *PermissionAPI) GetPermission(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		response.BadRequest(c, "参数错误", "权限代码不能为空")
		return
	}

	permission, err := api.permissionService.GetPermissionByCode(c.Request.Context(), code)
	if err != nil {
		response.NotFound(c, "权限不存在")
		return
	}

	response.Success(c, permission)
}

// CreatePermission 创建权限
//
//	@Summary		创建权限
//	@Description	管理员创建新权限
//	@Tags			Admin-Permission
//	@Accept			json
//	@Produce		json
//	@Param			request	body		auth.Permission	true	"权限信息"
//	@Success		201		{object}	shared.APIResponse
//	@Failure		400		{object}	shared.APIResponse
//	@Failure		401		{object}	shared.APIResponse
//	@Failure		403		{object}	shared.APIResponse
//	@Router			/api/v1/admin/permissions [post]
func (api *PermissionAPI) CreatePermission(c *gin.Context) {
	var req auth.Permission
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	if err := api.permissionService.CreatePermission(c.Request.Context(), &req); err != nil {
		c.Error(err)
		return
	}

	response.Created(c, nil)
}

// UpdatePermission 更新权限
//
//	@Summary		更新权限
//	@Description	管理员更新权限信息
//	@Tags			Admin-Permission
//	@Accept			json
//	@Produce		json
//	@Param			code		path		string			true	"权限代码"
//	@Param			request	body		auth.Permission	true	"权限信息"
//	@Success		200		{object}	shared.APIResponse
//	@Failure		400		{object}	shared.APIResponse
//	@Failure		404		{object}	shared.APIResponse
//	@Router			/api/v1/admin/permissions/{code} [put]
func (api *PermissionAPI) UpdatePermission(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		response.BadRequest(c, "参数错误", "权限代码不能为空")
		return
	}

	var req auth.Permission
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	req.Code = code
	if err := api.permissionService.UpdatePermission(c.Request.Context(), &req); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, nil)
}

// DeletePermission 删除权限
//
//	@Summary		删除权限
//	@Description	管理员删除权限
//	@Tags			Admin-Permission
//	@Accept			json
//	@Produce		json
//	@Param			code	path	string	true	"权限代码"
//	@Success		200	{object}	shared.APIResponse
//	@Failure		400	{object}	shared.APIResponse
//	@Failure		404	{object}	shared.APIResponse
//	@Router			/api/v1/admin/permissions/{code} [delete]
func (api *PermissionAPI) DeletePermission(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		response.BadRequest(c, "参数错误", "权限代码不能为空")
		return
	}

	if err := api.permissionService.DeletePermission(c.Request.Context(), code); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, nil)
}

// ==================== 角色管理 ====================

// GetAllRoles 获取所有角色
//
//	@Summary		获取所有角色
//	@Description	管理员获取所有角色列表
//	@Tags			Admin-Role
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	shared.APIResponse
//	@Failure		401	{object}	shared.APIResponse
//	@Failure		403	{object}	shared.APIResponse
//	@Router			/api/v1/admin/roles [get]
func (api *PermissionAPI) GetAllRoles(c *gin.Context) {
	roles, err := api.permissionService.GetAllRoles(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, roles)
}

// GetRole 获取角色详情
//
//	@Summary		获取角色详情
//	@Description	管理员获取角色详细信息
//	@Tags			Admin-Role
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"角色ID"
//	@Success		200	{object}	shared.APIResponse
//	@Failure		400	{object}	shared.APIResponse
//	@Failure		404	{object}	shared.APIResponse
//	@Router			/api/v1/admin/roles/{id} [get]
func (api *PermissionAPI) GetRole(c *gin.Context) {
	roleID := c.Param("id")
	if roleID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "角色ID不能为空")
		return
	}

	role, err := api.permissionService.GetRoleByID(c.Request.Context(), roleID)
	if err != nil {
		response.NotFound(c, "角色不存在")
		return
	}

	response.Success(c, role)
}

// CreateRole 创建角色
//
//	@Summary		创建角色
//	@Description	管理员创建新角色
//	@Tags			Admin-Role
//	@Accept			json
//	@Produce		json
//	@Param			request	body		object	true	"角色信息"
//	@Success		201		{object}	shared.APIResponse
//	@Failure		400		{object}	shared.APIResponse
//	@Failure		401		{object}	shared.APIResponse
//	@Failure		403		{object}	shared.APIResponse
//	@Router			/api/v1/admin/roles [post]
func (api *PermissionAPI) CreateRole(c *gin.Context) {
	var req auth.Role
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	if err := api.permissionService.CreateRole(c.Request.Context(), &req); err != nil {
		c.Error(err)
		return
	}

	response.Created(c, nil)
}

// UpdateRole 更新角色
//
//	@Summary		更新角色
//	@Description	管理员更新角色信息
//	@Tags			Admin-Role
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string		true	"角色ID"
//	@Param			request	body		object	true	"角色信息"
//	@Success		200		{object}	shared.APIResponse
//	@Failure		400		{object}	shared.APIResponse
//	@Failure		404		{object}	shared.APIResponse
//	@Router			/api/v1/admin/roles/{id} [put]
func (api *PermissionAPI) UpdateRole(c *gin.Context) {
	roleID := c.Param("id")
	if roleID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "角色ID不能为空")
		return
	}

	var req auth.Role
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	req.ID = roleID
	if err := api.permissionService.UpdateRole(c.Request.Context(), &req); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, nil)
}

// DeleteRole 删除角色
//
//	@Summary		删除角色
//	@Description	管理员删除角色
//	@Tags			Admin-Role
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"角色ID"
//	@Success		200	{object}	shared.APIResponse
//	@Failure		400	{object}	shared.APIResponse
//	@Failure		404	{object}	shared.APIResponse
//	@Router			/api/v1/admin/roles/{id} [delete]
func (api *PermissionAPI) DeleteRole(c *gin.Context) {
	roleID := c.Param("id")
	if roleID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "角色ID不能为空")
		return
	}

	if err := api.permissionService.DeleteRole(c.Request.Context(), roleID); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, nil)
}

// AssignPermissionToRole 为角色分配权限
//
//	@Summary		为角色分配权限
//	@Description	管理员为角色分配权限
//	@Tags			Admin-Role
//	@Accept			json
//	@Produce		json
//	@Param			id			path		string	true	"角色ID"
//	@Param			permissionCode	path		string	true	"权限代码"
//	@Success		200			{object}	shared.APIResponse
//	@Failure		400			{object}	shared.APIResponse
//	@Failure		404			{object}	shared.APIResponse
//	@Router			/api/v1/admin/roles/{id}/permissions/{permissionCode} [post]
func (api *PermissionAPI) AssignPermissionToRole(c *gin.Context) {
	roleID := c.Param("id")
	permissionCode := c.Param("permissionCode")

	if roleID == "" || permissionCode == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "角色ID和权限代码不能为空")
		return
	}

	if err := api.permissionService.AssignPermissionToRole(c.Request.Context(), roleID, permissionCode); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, nil)
}

// RemovePermissionFromRole 移除角色权限
//
//	@Summary		移除角色权限
//	@Description	管理员移除角色权限
//	@Tags			Admin-Role
//	@Accept			json
//	@Produce		json
//	@Param			id			path		string	true	"角色ID"
//	@Param			permissionCode	path		string	true	"权限代码"
//	@Success		200			{object}	shared.APIResponse
//	@Failure		400			{object}	shared.APIResponse
//	@Failure		404			{object}	shared.APIResponse
//	@Router			/api/v1/admin/roles/{id}/permissions/{permissionCode} [delete]
func (api *PermissionAPI) RemovePermissionFromRole(c *gin.Context) {
	roleID := c.Param("id")
	permissionCode := c.Param("permissionCode")

	if roleID == "" || permissionCode == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "角色ID和权限代码不能为空")
		return
	}

	if err := api.permissionService.RemovePermissionFromRole(c.Request.Context(), roleID, permissionCode); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, nil)
}

// GetRolePermissions 获取角色的所有权限
//
//	@Summary		获取角色权限
//	@Description	管理员获取角色的所有权限
//	@Tags			Admin-Role
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"角色ID"
//	@Success		200	{object}	shared.APIResponse
//	@Failure		400	{object}	shared.APIResponse
//	@Failure		404	{object}	shared.APIResponse
//	@Router			/api/v1/admin/roles/{id}/permissions [get]
func (api *PermissionAPI) GetRolePermissions(c *gin.Context) {
	roleID := c.Param("id")
	if roleID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "角色ID不能为空")
		return
	}

	permissions, err := api.permissionService.GetRolePermissions(c.Request.Context(), roleID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, permissions)
}

// ==================== 用户角色管理 ====================

// GetUserRoles 获取用户的所有角色
//
//	@Summary		获取用户角色
//	@Description	管理员获取用户的所有角色
//	@Tags			Admin-User
//	@Accept			json
//	@Produce		json
//	@Param			userId	path		string	true	"用户ID"
//	@Success		200		{object}	shared.APIResponse
//	@Failure		400		{object}	shared.APIResponse
//	@Router			/api/v1/admin/users/{userId}/roles [get]
func (api *PermissionAPI) GetUserRoles(c *gin.Context) {
	userID := c.Param("userId")
	if userID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "用户ID不能为空")
		return
	}

	roles, err := api.permissionService.GetUserRoles(c.Request.Context(), userID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, roles)
}

// AssignRoleToUser 为用户分配角色
//
//	@Summary		为用户分配角色
//	@Description	管理员为用户分配角色
//	@Tags			Admin-User
//	@Accept			json
//	@Produce		json
//	@Param			userId		path		string					true	"用户ID"
//	@Param			request	body		AssignRoleRequest	true	"角色信息"
//	@Success		200			{object}	shared.APIResponse
//	@Failure		400			{object}	shared.APIResponse
//	@Router			/api/v1/admin/users/{userId}/roles [post]
func (api *PermissionAPI) AssignRoleToUser(c *gin.Context) {
	userID := c.Param("userId")
	if userID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "用户ID不能为空")
		return
	}

	var req AssignRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	if err := api.permissionService.AssignRoleToUser(c.Request.Context(), userID, req.Role); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, nil)
}

// RemoveRoleFromUser 移除用户角色
//
//	@Summary		移除用户角色
//	@Description	管理员移除用户角色
//	@Tags			Admin-User
//	@Accept			json
//	@Produce		json
//	@Param			userId	path		string	true	"用户ID"
//	@Param			role	query		string	true	"角色名称"
//	@Success		200		{object}	shared.APIResponse
//	@Failure		400	{object}	shared.APIResponse
//	@Router			/api/v1/admin/users/{userId}/roles [delete]
func (api *PermissionAPI) RemoveRoleFromUser(c *gin.Context) {
	userID := c.Param("userId")
	role := c.Query("role")

	if userID == "" || role == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "用户ID和角色不能为空")
		return
	}

	if err := api.permissionService.RemoveRoleFromUser(c.Request.Context(), userID, role); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, nil)
}

// GetUserPermissions 获取用户的所有权限
//
//	@Summary		获取用户权限
//	@Description	管理员获取用户的所有权限
//	@Tags			Admin-User
//	@Accept			json
//	@Produce		json
//	@Param			userId	path		string	true	"用户ID"
//	@Success		200		{object}	shared.APIResponse
//	@Failure		400		{object}	shared.APIResponse
//	@Router			/api/v1/admin/users/{userId}/permissions [get]
func (api *PermissionAPI) GetUserPermissions(c *gin.Context) {
	userID := c.Param("userId")
	if userID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "用户ID不能为空")
		return
	}

	permissions, err := api.permissionService.GetUserPermissions(c.Request.Context(), userID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, permissions)
}

// ==================== 请求体结构 ====================

// AssignRoleRequest 分配角色请求
type AssignRoleRequest struct {
	Role string `json:"role" binding:"required"`
}

// UpdateRolePermissionsRequest 更新角色权限请求
type UpdateRolePermissionsRequest struct {
	Permissions []string `json:"permissions" binding:"required"`
}

// ==================== 批量角色操作 ====================

// batchOperationFunc 批量操作函数类型
type batchOperationFunc func(ctx context.Context, userID, role string) error

// executeBatchRoleOperation 执行批量角色操作的通用方法
func (api *PermissionAPI) executeBatchRoleOperation(c *gin.Context, req BatchRoleOperationRequest, operation batchOperationFunc) {
	result := BatchOperationResult{
		Total:  len(req.UserIDs),
		Errors: []string{},
	}

	for _, userID := range req.UserIDs {
		if err := operation(c.Request.Context(), userID, req.Role); err != nil {
			result.Failed++
			result.Errors = append(result.Errors, userID+": "+err.Error())
		} else {
			result.Success++
		}
	}

	response.Success(c, result)
}

// BatchAssignRoleToUsers 批量分配角色
//
//	@Summary		批量分配角色
//	@Description	管理员批量为用户分配角色
//	@Tags			Admin-Role
//	@Accept			json
//	@Produce		json
//	@Param			request	body		BatchRoleOperationRequest	true	"批量角色操作请求"
//	@Success		200		{object}	shared.APIResponse
//	@Failure		400		{object}	shared.APIResponse
//	@Router			/api/v1/admin/users/batch-assign-role [post]
func (api *PermissionAPI) BatchAssignRoleToUsers(c *gin.Context) {
	var req BatchRoleOperationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	api.executeBatchRoleOperation(c, req, api.permissionService.AssignRoleToUser)
}

// BatchRevokeRoleFromUsers 批量撤销角色
//
//	@Summary		批量撤销角色
//	@Description	管理员批量撤销用户角色
//	@Tags			Admin-Role
//	@Accept			json
//	@Produce		json
//	@Param			request	body		BatchRoleOperationRequest	true	"批量角色操作请求"
//	@Success		200		{object}	shared.APIResponse
//	@Failure		400		{object}	shared.APIResponse
//	@Router			/api/v1/admin/users/batch-revoke-role [post]
func (api *PermissionAPI) BatchRevokeRoleFromUsers(c *gin.Context) {
	var req BatchRoleOperationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	api.executeBatchRoleOperation(c, req, api.permissionService.RemoveRoleFromUser)
}
