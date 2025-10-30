package user

import (
	serviceInterfaces "Qingyu_backend/service/interfaces/base"
	userServiceInterface "Qingyu_backend/service/interfaces/user"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/shared"
	"Qingyu_backend/pkg/utils"
)

// UserAPI 用户管理API处理器
type UserAPI struct {
	userService userServiceInterface.UserService
}

// NewUserAPI 创建用户API实例
func NewUserAPI(userService userServiceInterface.UserService) *UserAPI {
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
	serviceReq := &userServiceInterface.RegisterUserRequest{
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
	clientIP := utils.GetClientIP(c)

	// 调用Service层
	serviceReq := &userServiceInterface.LoginUserRequest{
		Username: req.Username,
		Password: req.Password,
		ClientIP: clientIP,
	}

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
		Token: resp.Token,
		User: UserBasicInfo{
			UserID:   resp.User.ID,
			Username: resp.User.Username,
			Email:    resp.User.Email,
			Role:     resp.User.Role,
		},
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
	serviceReq := &userServiceInterface.GetUserRequest{
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
	serviceReq := &userServiceInterface.UpdateUserRequest{
		ID:      userID.(string),
		Updates: updates,
	}

	updatedUser, err := api.userService.UpdateUser(c.Request.Context(), serviceReq)
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

	// 返回更新后的用户信息
	shared.Success(c, http.StatusOK, "更新成功", updatedUser)
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
	serviceReq := &userServiceInterface.UpdatePasswordRequest{
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

// GetUserProfile 获取用户公开信息（用于用户主页）
//
//	@Summary		获取用户公开信息
//	@Description	获取指定用户的公开信息，用于展示用户主页
//	@Tags			用户
//	@Accept			json
//	@Produce		json
//	@Param			userId	path		string	true	"用户ID"
//	@Success		200		{object}	shared.APIResponse{data=PublicUserProfileResponse}
//	@Failure		404		{object}	shared.ErrorResponse
//	@Failure		500		{object}	shared.ErrorResponse
//	@Router			/api/v1/users/{userId}/profile [get]
func (api *UserAPI) GetUserProfile(c *gin.Context) {
	userID := c.Param("userId")
	if userID == "" {
		shared.BadRequest(c, "参数错误", "用户ID不能为空")
		return
	}

	// 调用Service层获取用户信息
	serviceReq := &userServiceInterface.GetUserRequest{
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

	// 构建公开信息响应（不包含敏感信息）
	publicProfile := PublicUserProfileResponse{
		UserID:    resp.User.ID,
		Username:  resp.User.Username,
		Avatar:    resp.User.Avatar,
		Nickname:  resp.User.Nickname,
		Bio:       resp.User.Bio,
		Role:      resp.User.Role,
		CreatedAt: resp.User.CreatedAt,
	}

	shared.Success(c, http.StatusOK, "获取成功", publicProfile)
}

// GetUserBooks 获取用户的作品列表
//
//	@Summary		获取用户作品列表
//	@Description	获取指定用户的已发布作品列表
//	@Tags			用户
//	@Accept			json
//	@Produce		json
//	@Param			userId	path		string	true	"用户ID"
//	@Param			page	query		int		false	"页码"		default(1)
//	@Param			size	query		int		false	"每页数量"	default(20)
//	@Param			status	query		string	false	"状态筛选"	Enums(published, completed)
//	@Success		200		{object}	shared.APIResponse{data=UserBooksResponse}
//	@Failure		404		{object}	shared.ErrorResponse
//	@Failure		500		{object}	shared.ErrorResponse
//	@Router			/api/v1/users/{userId}/books [get]
func (api *UserAPI) GetUserBooks(c *gin.Context) {
	userID := c.Param("userId")
	if userID == "" {
		shared.BadRequest(c, "参数错误", "用户ID不能为空")
		return
	}

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	_ = c.Query("status") // status 参数预留，待 BookService 实现时使用

	// TODO: 调用BookService查询用户的已发布作品
	// 这里简化实现，返回模拟数据
	// 实际应该调用 bookService.GetBooksByAuthor(ctx, userID, page, size, status)

	books := []map[string]interface{}{
		{
			"book_id":     "book123",
			"title":       "示例作品",
			"cover":       "",
			"description": "这是一部精彩的作品",
			"category":    "玄幻",
			"status":      "published",
			"word_count":  100000,
			"created_at":  "2024-01-01T00:00:00Z",
		},
	}

	response := UserBooksResponse{
		Books: books,
		Total: len(books),
		Page:  page,
		Size:  size,
	}

	shared.Success(c, http.StatusOK, "获取成功", response)
}
