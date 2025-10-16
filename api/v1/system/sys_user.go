package system

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/shared"
	serviceInterfaces "Qingyu_backend/service/interfaces"
)

// UserAPI 用户管理API处理器
type UserAPI struct {
	userService serviceInterfaces.UserService
}

// NewUserAPI 创建用户API实例
func NewUserAPI(userService serviceInterfaces.UserService) *UserAPI {
	return &UserAPI{
		userService: userService,
	}
}

// Register 用户注册
//
//	@Summary		用户注册
//	@Description	注册新用户账号
//	@Tags			用户
//	@Accept			json
//	@Produce		json
//	@Param			request	body		RegisterRequest	true	"注册信息"
//	@Success		200		{object}	shared.APIResponse{data=RegisterResponse}
//	@Failure		400		{object}	shared.ErrorResponse
//	@Failure		500		{object}	shared.ErrorResponse
//	@Router			/api/v1/register [post]
func (api *UserAPI) Register(c *gin.Context) {
	var req RegisterRequest
	if !shared.ValidateRequest(c, &req) {
		return
	}

	// 调用Service层
	serviceReq := &serviceInterfaces.RegisterUserRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	resp, err := api.userService.RegisterUser(c.Request.Context(), serviceReq)
	if err != nil {
		// 根据错误类型返回不同的HTTP状态码
		if serviceErr, ok := err.(*serviceInterfaces.ServiceError); ok {
			switch serviceErr.Type {
			case serviceInterfaces.ErrorTypeValidation:
				shared.BadRequest(c, "注册失败", serviceErr.Message)
			case serviceInterfaces.ErrorTypeBusiness:
				shared.BadRequest(c, "注册失败", serviceErr.Message)
			default:
				shared.InternalError(c, "注册失败", err)
			}
			return
		}
		shared.InternalError(c, "注册失败", err)
		return
	}

	// 构建响应
	registerResp := RegisterResponse{
		UserID:   resp.User.ID,
		Username: resp.User.Username,
		Email:    resp.User.Email,
		Role:     resp.User.Role,
		Status:   string(resp.User.Status),
		Token:    resp.Token,
	}

	shared.Success(c, http.StatusCreated, "注册成功", registerResp)
}

// Login 用户登录
//
//	@Summary		用户登录
//	@Description	用户登录获取Token
//	@Tags			用户
//	@Accept			json
//	@Produce		json
//	@Param			request	body		LoginRequest	true	"登录信息"
//	@Success		200		{object}	shared.APIResponse{data=LoginResponse}
//	@Failure		400		{object}	shared.ErrorResponse
//	@Failure		401		{object}	shared.ErrorResponse
//	@Failure		500		{object}	shared.ErrorResponse
//	@Router			/api/v1/login [post]
func (api *UserAPI) Login(c *gin.Context) {
	var req LoginRequest
	if !shared.ValidateRequest(c, &req) {
		return
	}

	// 获取客户端IP
	clientIP := c.ClientIP()

	// 调用Service层
	serviceReq := &serviceInterfaces.LoginUserRequest{
		Username: req.Username,
		Password: req.Password,
	}

	// TODO: 将IP通过context传递给Service层
	_ = clientIP

	resp, err := api.userService.LoginUser(c.Request.Context(), serviceReq)
	if err != nil {
		if serviceErr, ok := err.(*serviceInterfaces.ServiceError); ok {
			switch serviceErr.Type {
			case serviceInterfaces.ErrorTypeNotFound:
				shared.Unauthorized(c, "用户名或密码错误")
			case serviceInterfaces.ErrorTypeUnauthorized:
				shared.Unauthorized(c, "用户名或密码错误")
			case serviceInterfaces.ErrorTypeValidation:
				shared.BadRequest(c, "登录失败", serviceErr.Message)
			default:
				shared.InternalError(c, "登录失败", err)
			}
			return
		}
		shared.InternalError(c, "登录失败", err)
		return
	}

	// 构建响应
	loginResp := LoginResponse{
		UserID:   resp.User.ID,
		Username: resp.User.Username,
		Email:    resp.User.Email,
		Token:    resp.Token,
	}

	shared.Success(c, http.StatusOK, "登录成功", loginResp)
}

// GetProfile 获取当前用户信息
//
//	@Summary		获取当前用户信息
//	@Description	获取当前登录用户的详细信息
//	@Tags			用户
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success		200	{object}	shared.APIResponse{data=UserProfileResponse}
//	@Failure		401	{object}	shared.ErrorResponse
//	@Failure		404	{object}	shared.ErrorResponse
//	@Failure		500	{object}	shared.ErrorResponse
//	@Router			/api/v1/users/profile [get]
func (api *UserAPI) GetProfile(c *gin.Context) {
	// 从Context中获取当前用户ID（由JWT中间件设置）
	userID, exists := c.Get("user_id")
	if !exists {
		shared.Unauthorized(c, "未认证")
		return
	}

	// 调用Service层
	serviceReq := &serviceInterfaces.GetUserRequest{
		ID: userID.(string),
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

// UpdateProfile 更新当前用户信息
//
//	@Summary		更新当前用户信息
//	@Description	更新当前登录用户的个人信息
//	@Tags			用户
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			request	body		UpdateProfileRequest	true	"更新信息"
//	@Success		200		{object}	shared.APIResponse
//	@Failure		400		{object}	shared.ErrorResponse
//	@Failure		401		{object}	shared.ErrorResponse
//	@Failure		500		{object}	shared.ErrorResponse
//	@Router			/api/v1/users/profile [put]
func (api *UserAPI) UpdateProfile(c *gin.Context) {
	// 从Context中获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		shared.Unauthorized(c, "未认证")
		return
	}

	var req UpdateProfileRequest
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

	// 调用Service层
	serviceReq := &serviceInterfaces.UpdateUserRequest{
		ID:      userID.(string),
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

// ChangePassword 修改密码
//
//	@Summary		修改密码
//	@Description	修改当前用户的登录密码
//	@Tags			用户
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			request	body		ChangePasswordRequest	true	"密码信息"
//	@Success		200		{object}	shared.APIResponse
//	@Failure		400		{object}	shared.ErrorResponse
//	@Failure		401		{object}	shared.ErrorResponse
//	@Failure		500		{object}	shared.ErrorResponse
//	@Router			/api/v1/users/password [put]
func (api *UserAPI) ChangePassword(c *gin.Context) {
	// 从Context中获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		shared.Unauthorized(c, "未认证")
		return
	}

	var req ChangePasswordRequest
	if !shared.ValidateRequest(c, &req) {
		return
	}

	// 调用Service层
	serviceReq := &serviceInterfaces.UpdatePasswordRequest{
		ID:          userID.(string),
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	}

	_, err := api.userService.UpdatePassword(c.Request.Context(), serviceReq)
	if err != nil {
		if serviceErr, ok := err.(*serviceInterfaces.ServiceError); ok {
			switch serviceErr.Type {
			case serviceInterfaces.ErrorTypeNotFound:
				shared.NotFound(c, "用户不存在")
			case serviceInterfaces.ErrorTypeUnauthorized:
				shared.Unauthorized(c, "旧密码错误")
			case serviceInterfaces.ErrorTypeValidation:
				shared.BadRequest(c, "修改密码失败", serviceErr.Message)
			default:
				shared.InternalError(c, "修改密码失败", err)
			}
			return
		}
		shared.InternalError(c, "修改密码失败", err)
		return
	}

	shared.Success(c, http.StatusOK, "密码修改成功", nil)
}

// GetUser 获取指定用户信息（管理员）
//
//	@Summary		获取指定用户信息
//	@Description	管理员获取指定用户的详细信息
//	@Tags			用户-管理
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
func (api *UserAPI) GetUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		shared.BadRequest(c, "用户ID不能为空", "")
		return
	}

	// 调用Service层
	serviceReq := &serviceInterfaces.GetUserRequest{
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
//	@Tags			用户-管理
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			page		query		int						false	"页码"				default(1)
//	@Param			page_size	query		int						false	"每页数量"			default(10)
//	@Param			username	query		string					false	"用户名筛选"
//	@Param			email		query		string					false	"邮箱筛选"
//	@Param			role		query		string					false	"角色筛选"
//	@Param			status		query		usersModel.UserStatus	false	"状态筛选"
//	@Success		200			{object}	shared.PaginatedResponse{data=[]UserProfileResponse}
//	@Failure		400			{object}	shared.ErrorResponse
//	@Failure		401			{object}	shared.ErrorResponse
//	@Failure		403			{object}	shared.ErrorResponse
//	@Failure		500			{object}	shared.ErrorResponse
//	@Router			/api/v1/admin/users [get]
func (api *UserAPI) ListUsers(c *gin.Context) {
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
	serviceReq := &serviceInterfaces.ListUsersRequest{
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
//	@Tags			用户-管理
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
func (api *UserAPI) UpdateUser(c *gin.Context) {
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
	serviceReq := &serviceInterfaces.UpdateUserRequest{
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
//	@Tags			用户-管理
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
func (api *UserAPI) DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		shared.BadRequest(c, "用户ID不能为空", "")
		return
	}

	// 调用Service层
	serviceReq := &serviceInterfaces.DeleteUserRequest{
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
