package admin

import (
	serviceInterfaces "Qingyu_backend/service/interfaces/base"
	"Qingyu_backend/service/interfaces/user"
	"net/http"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/shared"
)

// UserAdminAPI 用户管理API处理器（管理员）
type UserAdminAPI struct {
	userService user.UserService
}

// NewUserAdminAPI 创建用户管理API实例
func NewUserAdminAPI(userService user.UserService) *UserAdminAPI {
	return &UserAdminAPI{
		userService: userService,
	}
}

// GetUser 获取指定用户信息（管理员）
//
//	@Summary		获取指定用户信息
//	@Description	管理员获取指定用户的详细信息
//	@Tags			管理员-用户管理
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		string	true	"用户ID"
//	@Success		200	{object}	shared.APIResponse{data=UserProfileResponse}
//	@Failure		400	{object}	shared.ErrorResponse
//	@Failure		401	{object}	shared.ErrorResponse
//	@Failure		403	{object}	shared.ErrorResponse
//	@Failure		404	{object}	shared.ErrorResponse
//	@Failure		500	{object}	shared.ErrorResponse
//	@Router			/api/v1/admin/users/{id} [get]
func (api *UserAdminAPI) GetUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		shared.BadRequest(c, "用户ID不能为空", "")
		return
	}

	// 调用Service层
	serviceReq := &user.GetUserRequest{
		ID: userID,
	}

	resp, err := api.userService.GetUser(c.Request.Context(), serviceReq)
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
	profileResp := UserProfileResponse{
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

// ListUsers 获取用户列表（管理员）
//
//	@Summary		获取用户列表
//	@Description	管理员获取用户列表（支持分页和筛选）
//	@Tags			管理员-用户管理
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
//	@Router			/api/v1/admin/users [get]
func (api *UserAdminAPI) ListUsers(c *gin.Context) {
	var req ListUsersRequest
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
	serviceReq := &user.ListUsersRequest{
		Username: req.Username,
		Email:    req.Email,
		Status:   string(req.Status),
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	resp, err := api.userService.ListUsers(c.Request.Context(), serviceReq)
	if err != nil {
		shared.InternalError(c, "获取用户列表失败", err)
		return
	}

	// 构建响应
	var users []UserProfileResponse
	for _, user := range resp.Users {
		users = append(users, UserProfileResponse{
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

// UpdateUser 更新指定用户信息（管理员）
//
//	@Summary		更新指定用户信息
//	@Description	管理员更新指定用户的信息
//	@Tags			管理员-用户管理
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id		path		string					true	"用户ID"
//	@Param			request	body		AdminUpdateUserRequest	true	"更新信息"
//	@Success		200		{object}	shared.APIResponse
//	@Failure		400		{object}	shared.ErrorResponse
//	@Failure		401		{object}	shared.ErrorResponse
//	@Failure		403		{object}	shared.ErrorResponse
//	@Failure		404		{object}	shared.ErrorResponse
//	@Failure		500		{object}	shared.ErrorResponse
//	@Router			/api/v1/admin/users/{id} [put]
func (api *UserAdminAPI) UpdateUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		shared.BadRequest(c, "用户ID不能为空", "")
		return
	}

	var req AdminUpdateUserRequest
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
	serviceReq := &user.UpdateUserRequest{
		ID:      userID,
		Updates: updates,
	}

	_, err := api.userService.UpdateUser(c.Request.Context(), serviceReq)
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
//	@Tags			管理员-用户管理
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		string	true	"用户ID"
//	@Success		200	{object}	shared.APIResponse
//	@Failure		400	{object}	shared.ErrorResponse
//	@Failure		401	{object}	shared.ErrorResponse
//	@Failure		403	{object}	shared.ErrorResponse
//	@Failure		404	{object}	shared.ErrorResponse
//	@Failure		500	{object}	shared.ErrorResponse
//	@Router			/api/v1/admin/users/{id} [delete]
func (api *UserAdminAPI) DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		shared.BadRequest(c, "用户ID不能为空", "")
		return
	}

	// 调用Service层
	serviceReq := &user.DeleteUserRequest{
		ID: userID,
	}

	_, err := api.userService.DeleteUser(c.Request.Context(), serviceReq)
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
//	@Tags			管理员-用户管理
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id		path		string			true	"用户ID"
//	@Param			request	body		BanUserRequest	true	"封禁信息"
//	@Success		200		{object}	shared.APIResponse
//	@Failure		400		{object}	shared.ErrorResponse
//	@Failure		401		{object}	shared.ErrorResponse
//	@Failure		403		{object}	shared.ErrorResponse
//	@Failure		404		{object}	shared.ErrorResponse
//	@Failure		500		{object}	shared.ErrorResponse
//	@Router			/api/v1/admin/users/{id}/ban [post]
func (api *UserAdminAPI) BanUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		shared.BadRequest(c, "用户ID不能为空", "")
		return
	}

	var req BanUserRequest
	if !shared.ValidateRequest(c, &req) {
		return
	}

	// 构建更新数据 - 设置为banned状态
	updates := make(map[string]interface{})
	updates["status"] = "banned"
	// TODO: 可以扩展添加ban_reason, ban_until等字段

	// 调用Service层
	serviceReq := &user.UpdateUserRequest{
		ID:      userID,
		Updates: updates,
	}

	_, err := api.userService.UpdateUser(c.Request.Context(), serviceReq)
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
//	@Tags			管理员-用户管理
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		string	true	"用户ID"
//	@Success		200	{object}	shared.APIResponse
//	@Failure		400	{object}	shared.ErrorResponse
//	@Failure		401	{object}	shared.ErrorResponse
//	@Failure		403	{object}	shared.ErrorResponse
//	@Failure		404	{object}	shared.ErrorResponse
//	@Failure		500	{object}	shared.ErrorResponse
//	@Router			/api/v1/admin/users/{id}/unban [post]
func (api *UserAdminAPI) UnbanUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		shared.BadRequest(c, "用户ID不能为空", "")
		return
	}

	// 构建更新数据 - 恢复为active状态
	updates := make(map[string]interface{})
	updates["status"] = "active"

	// 调用Service层
	serviceReq := &user.UpdateUserRequest{
		ID:      userID,
		Updates: updates,
	}

	_, err := api.userService.UpdateUser(c.Request.Context(), serviceReq)
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
