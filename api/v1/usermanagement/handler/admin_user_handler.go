package handler

import (
	serviceInterfaces "Qingyu_backend/service/interfaces/base"
	userServiceInterface "Qingyu_backend/service/interfaces/user"
	"net/http"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/shared"
	"Qingyu_backend/api/v1/usermanagement/dto"
)

// AdminUserHandler 管理员用户管理处理器
type AdminUserHandler struct {
	userService userServiceInterface.UserService
}

// NewAdminUserHandler 创建管理员用户管理处理器实例
func NewAdminUserHandler(userService userServiceInterface.UserService) *AdminUserHandler {
	return &AdminUserHandler{
		userService: userService,
	}
}

// ListUsers 获取用户列表（管理员）
//
//	@Summary		获取用户列表
//	@Description	管理员获取用户列表（支持分页和筛选）
//	@Tags			用户管理-管理员
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			page		query		int		false	"页码"				default(1)
//	@Param			page_size	query		int		false	"每页数量"			default(10)
//	@Param			username	query		string	false	"用户名筛选"
//	@Param			email		query		string	false	"邮箱筛选"
//	@Param			role		query		string	false	"角色筛选"
//	@Param			status		query		string	false	"状态筛选"
//	@Success		200			{object}	shared.PaginatedResponse
//	@Failure		400			{object}	shared.ErrorResponse
//	@Failure		401			{object}	shared.ErrorResponse
//	@Failure		403			{object}	shared.ErrorResponse
//	@Failure		500			{object}	shared.ErrorResponse
//	@Router			/api/v1/user-management/users [get]
func (h *AdminUserHandler) ListUsers(c *gin.Context) {
	var req dto.ListUsersRequest
	if !shared.ValidateQueryParams(c, &req) {
		return
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 调用Service层
	serviceReq := &userServiceInterface.ListUsersRequest{
		Username: req.Username,
		Email:    req.Email,
		Status:   req.Status,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	resp, err := h.userService.ListUsers(c.Request.Context(), serviceReq)
	if err != nil {
		shared.InternalError(c, "获取用户列表失败", err)
		return
	}

	// 构建响应
	var users []dto.UserProfileResponse
	for _, user := range resp.Users {
		users = append(users, dto.UserProfileResponse{
			UserID:        user.ID,
			Username:      user.Username,
			Email:         user.Email,
			Phone:         user.Phone,
			Role:          user.Role,
			Status:        string(user.Status),
			Avatar:        user.Avatar,
			Nickname:      user.Nickname,
			Bio:           user.Bio,
			EmailVerified: user.EmailVerified,
			PhoneVerified: user.PhoneVerified,
			LastLoginAt:   user.LastLoginAt,
			LastLoginIP:   user.LastLoginIP,
			CreatedAt:     user.CreatedAt,
			UpdatedAt:     user.UpdatedAt,
		})
	}

	shared.Paginated(c, users, resp.Total, resp.Page, resp.PageSize, "获取成功")
}

// GetUser 获取指定用户信息（管理员）
//
//	@Summary		获取指定用户信息
//	@Description	管理员获取指定用户的详细信息
//	@Tags			用户管理-管理员
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		string	true	"用户ID"
//	@Success		200		{object}	shared.APIResponse{data=dto.UserProfileResponse}
//	@Failure		400		{object}	shared.ErrorResponse
//	@Failure		401		{object}	shared.ErrorResponse
//	@Failure		403		{object}	shared.ErrorResponse
//	@Failure		404		{object}	shared.ErrorResponse
//	@Failure		500		{object}	shared.ErrorResponse
//	@Router			/api/v1/user-management/users/{id} [get]
func (h *AdminUserHandler) GetUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		shared.BadRequest(c, "用户ID不能为空", "")
		return
	}

	// 调用Service层
	serviceReq := &userServiceInterface.GetUserRequest{
		ID: userID,
	}

	resp, err := h.userService.GetUser(c.Request.Context(), serviceReq)
	if err != nil {
		if serviceErr, ok := err.(*serviceInterfaces.ServiceError); ok {
			switch serviceErr.Type {
			case serviceInterfaces.ErrorTypeNotFound:
				shared.NotFound(c, "用户不存在")
			default:
				shared.InternalError(c, "获取用户信息失败", err)
			}
			return
		}
		shared.InternalError(c, "获取用户信息失败", err)
		return
	}

	// 构建响应
	profileResp := dto.UserProfileResponse{
		UserID:        resp.User.ID,
		Username:      resp.User.Username,
		Email:         resp.User.Email,
		Phone:         resp.User.Phone,
		Role:          resp.User.Role,
		Status:        string(resp.User.Status),
		Avatar:        resp.User.Avatar,
		Nickname:      resp.User.Nickname,
		Bio:           resp.User.Bio,
		EmailVerified: resp.User.EmailVerified,
		PhoneVerified: resp.User.PhoneVerified,
		LastLoginAt:   resp.User.LastLoginAt,
		LastLoginIP:   resp.User.LastLoginIP,
		CreatedAt:     resp.User.CreatedAt,
		UpdatedAt:     resp.User.UpdatedAt,
	}

	shared.Success(c, http.StatusOK, "获取成功", profileResp)
}

// UpdateUser 更新指定用户信息（管理员）
//
//	@Summary		更新指定用户信息
//	@Description	管理员更新指定用户的信息
//	@Tags			用户管理-管理员
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id		path		string						true	"用户ID"
//	@Param			request	body		dto.AdminUpdateUserRequest	true	"更新信息"
//	@Success		200		{object}	shared.APIResponse
//	@Failure		400		{object}	shared.ErrorResponse
//	@Failure		401		{object}	shared.ErrorResponse
//	@Failure		403		{object}	shared.ErrorResponse
//	@Failure		404		{object}	shared.ErrorResponse
//	@Failure		500		{object}	shared.ErrorResponse
//	@Router			/api/v1/user-management/users/{id} [put]
func (h *AdminUserHandler) UpdateUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		shared.BadRequest(c, "用户ID不能为空", "")
		return
	}

	var req dto.AdminUpdateUserRequest
	if !shared.ValidateRequest(c, &req) {
		return
	}

	// 构建更新数据
	updates := make(map[string]interface{})
	if req.Nickname != nil {
		updates["nickname"] = *req.Nickname
	}
	if req.Bio != nil {
		updates["bio"] = *req.Bio
	}
	if req.Avatar != nil {
		updates["avatar"] = *req.Avatar
	}
	if req.Phone != nil {
		updates["phone"] = *req.Phone
	}
	if req.Role != nil {
		updates["role"] = *req.Role
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.EmailVerified != nil {
		updates["email_verified"] = *req.EmailVerified
	}
	if req.PhoneVerified != nil {
		updates["phone_verified"] = *req.PhoneVerified
	}

	// 调用Service层
	serviceReq := &userServiceInterface.UpdateUserRequest{
		ID:      userID,
		Updates: updates,
	}

	_, err := h.userService.UpdateUser(c.Request.Context(), serviceReq)
	if err != nil {
		if serviceErr, ok := err.(*serviceInterfaces.ServiceError); ok {
			switch serviceErr.Type {
			case serviceInterfaces.ErrorTypeNotFound:
				shared.NotFound(c, "用户不存在")
			case serviceInterfaces.ErrorTypeValidation:
				shared.BadRequest(c, "更新失败", serviceErr.Message)
			default:
				shared.InternalError(c, "更新失败", err)
			}
			return
		}
		shared.InternalError(c, "更新失败", err)
		return
	}

	shared.Success(c, http.StatusOK, "更新成功", nil)
}

// DeleteUser 删除用户（管理员）
//
//	@Summary		删除用户
//	@Description	管理员删除指定用户
//	@Tags			用户管理-管理员
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		string	true	"用户ID"
//	@Success		200	{object}	shared.APIResponse
//	@Failure		400	{object}	shared.ErrorResponse
//	@Failure		401	{object}	shared.ErrorResponse
//	@Failure		403	{object}	shared.ErrorResponse
//	@Failure		404		{object}	shared.ErrorResponse
//	@Failure		500		{object}	shared.ErrorResponse
//	@Router			/api/v1/user-management/users/{id} [delete]
func (h *AdminUserHandler) DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		shared.BadRequest(c, "用户ID不能为空", "")
		return
	}

	// 调用Service层
	serviceReq := &userServiceInterface.DeleteUserRequest{
		ID: userID,
	}

	_, err := h.userService.DeleteUser(c.Request.Context(), serviceReq)
	if err != nil {
		if serviceErr, ok := err.(*serviceInterfaces.ServiceError); ok {
			switch serviceErr.Type {
			case serviceInterfaces.ErrorTypeNotFound:
				shared.NotFound(c, "用户不存在")
			default:
				shared.InternalError(c, "删除失败", err)
			}
			return
		}
		shared.InternalError(c, "删除失败", err)
		return
	}

	shared.Success(c, http.StatusOK, "删除成功", nil)
}

// BanUser 封禁用户（管理员）
//
//	@Summary		封禁用户
//	@Description	管理员封禁指定用户
//	@Tags			用户管理-管理员
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id		path		string				true	"用户ID"
//	@Param			request	body		dto.BanUserRequest	true	"封禁信息"
//	@Success		200		{object}	shared.APIResponse
//	@Failure		400		{object}	shared.ErrorResponse
//	@Failure		401		{object}	shared.ErrorResponse
//	@Failure		403		{object}	shared.ErrorResponse
//	@Failure		404		{object}	shared.ErrorResponse
//	@Failure		500		{object}	shared.ErrorResponse
//	@Router			/api/v1/user-management/users/{id}/ban [post]
func (h *AdminUserHandler) BanUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		shared.BadRequest(c, "用户ID不能为空", "")
		return
	}

	var req dto.BanUserRequest
	if !shared.ValidateRequest(c, &req) {
		return
	}

	// 构建更新数据 - 设置为banned状态
	updates := make(map[string]interface{})
	updates["status"] = "banned"
	// 添加ban_reason和ban_until字段以完善用户信息扩展
	if req.Reason != "" {
		updates["ban_reason"] = req.Reason
	}
	if req.BanUntil != nil {
		updates["ban_until"] = *req.BanUntil
	}

	// 调用Service层
	serviceReq := &userServiceInterface.UpdateUserRequest{
		ID:      userID,
		Updates: updates,
	}

	_, err := h.userService.UpdateUser(c.Request.Context(), serviceReq)
	if err != nil {
		if serviceErr, ok := err.(*serviceInterfaces.ServiceError); ok {
			switch serviceErr.Type {
			case serviceInterfaces.ErrorTypeNotFound:
				shared.NotFound(c, "用户不存在")
			default:
				shared.InternalError(c, "封禁失败", err)
			}
			return
		}
		shared.InternalError(c, "封禁失败", err)
		return
	}

	shared.Success(c, http.StatusOK, "封禁成功", nil)
}

// UnbanUser 解除封禁（管理员）
//
//	@Summary		解除封禁
//	@Description	管理员解除用户封禁
//	@Tags			用户管理-管理员
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		string	true	"用户ID"
//	@Success		200	{object}	shared.APIResponse
//	@Failure		400	{object}	shared.ErrorResponse
//	@Failure		401		{object}	shared.ErrorResponse
//	@Failure		403		{object}	shared.ErrorResponse
//	@Failure		404		{object}	shared.ErrorResponse
//	@Failure		500		{object}	shared.ErrorResponse
//	@Router			/api/v1/user-management/users/{id}/unban [post]
func (h *AdminUserHandler) UnbanUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		shared.BadRequest(c, "用户ID不能为空", "")
		return
	}

	// 构建更新数据 - 恢复为active状态
	updates := make(map[string]interface{})
	updates["status"] = "active"

	// 调用Service层
	serviceReq := &userServiceInterface.UpdateUserRequest{
		ID:      userID,
		Updates: updates,
	}

	_, err := h.userService.UpdateUser(c.Request.Context(), serviceReq)
	if err != nil {
		if serviceErr, ok := err.(*serviceInterfaces.ServiceError); ok {
			switch serviceErr.Type {
			case serviceInterfaces.ErrorTypeNotFound:
				shared.NotFound(c, "用户不存在")
			default:
				shared.InternalError(c, "解除封禁失败", err)
			}
			return
		}
		shared.InternalError(c, "解除封禁失败", err)
		return
	}

	shared.Success(c, http.StatusOK, "解除封禁成功", nil)
}
