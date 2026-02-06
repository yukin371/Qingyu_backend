package handler

import (
	serviceInterfaces "Qingyu_backend/service/interfaces/base"
	userServiceInterface "Qingyu_backend/service/interfaces/user"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/shared"
	"Qingyu_backend/api/v1/user/dto"
	"Qingyu_backend/pkg/utils"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	userService userServiceInterface.UserService
}

// NewAuthHandler 创建认证处理器实例
func NewAuthHandler(userService userServiceInterface.UserService) *AuthHandler {
	return &AuthHandler{
		userService: userService,
	}
}

// Register 用户注册
//
//	@Summary		用户注册
//	@Description	注册新用户账户
//	@Tags			用户管理-认证
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.RegisterRequest	true	"注册信息"
//	@Success		200		{object}	shared.APIResponse{data=dto.RegisterResponse}
//	@Failure		400		{object}	shared.ErrorResponse
//	@Failure		500		{object}	shared.ErrorResponse
//	@Router			/api/v1/user/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if !shared.ValidateRequest(c, &req) {
		return
	}

	// 调用Service层
	serviceReq := &userServiceInterface.RegisterUserRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	resp, err := h.userService.RegisterUser(c.Request.Context(), serviceReq)
	if err != nil {
		// 根据错误类型返回不同的HTTP状态码
		if serviceErr, ok := err.(*serviceInterfaces.ServiceError); ok {
			switch serviceErr.Type {
			case serviceInterfaces.ErrorTypeValidation:
				shared.BadRequest(c, "注册失败", serviceErr.Message)
			case serviceInterfaces.ErrorTypeBusiness:
				// 根据错误消息返回具体的错误码
				if serviceErr.Message == "用户名已存在" {
					c.JSON(http.StatusConflict, shared.ErrorResponse{
						Code:      2003,
						Message:   "用户名已被注册",
						Timestamp: time.Now().UnixMilli(),
					})
				} else if serviceErr.Message == "邮箱已存在" {
					c.JSON(http.StatusConflict, shared.ErrorResponse{
						Code:      2004,
						Message:   "邮箱已被注册",
						Timestamp: time.Now().UnixMilli(),
					})
				} else {
					shared.BadRequest(c, "注册失败", serviceErr.Message)
				}
			default:
				shared.InternalError(c, "注册失败", err)
			}
			return
		}
		shared.InternalError(c, "注册失败", err)
		return
	}

	// 构建响应
	role := ""
	if len(resp.User.Roles) > 0 {
		role = resp.User.Roles[0]
	}
	registerResp := dto.RegisterResponse{
		UserID:   resp.User.ID,
		Username: resp.User.Username,
		Email:    resp.User.Email,
		Role:     role,
		Roles:    resp.User.Roles, // 返回完整的角色列表
		Status:   resp.User.Status,
		Token:    resp.Token,
	}

	shared.Success(c, http.StatusCreated, "注册成功", registerResp)
}

// Login 用户登录
//
//	@Summary		用户登录
//	@Description	用户登录获取Token
//	@Tags			用户管理-认证
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.LoginRequest	true	"登录信息"
//	@Success		200		{object}	shared.APIResponse{data=dto.LoginResponse}
//	@Failure		400		{object}	shared.ErrorResponse
//	@Failure		401		{object}	shared.ErrorResponse
//	@Failure		500		{object}	shared.ErrorResponse
//	@Router			/api/v1/user/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
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

	resp, err := h.userService.LoginUser(c.Request.Context(), serviceReq)
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
	role := ""
	if len(resp.User.Roles) > 0 {
		role = resp.User.Roles[0]
	}
	loginResp := dto.LoginResponse{
		Token: resp.Token,
		User: dto.UserBasicInfo{
			UserID:   resp.User.ID,
			Username: resp.User.Username,
			Email:    resp.User.Email,
			Role:     role,
			Roles:    resp.User.Roles, // 返回完整的角色列表
		},
		Roles: resp.User.Roles, // 顶层也返回roles，方便前端访问
	}

	shared.Success(c, http.StatusOK, "登录成功", loginResp)
}

// Logout 用户登出
//
//	@Summary		用户登出
//	@Description	用户登出，清除服务端会话/Token
//	@Tags			用户管理-认证
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success		200	{object}	shared.APIResponse
//	@Failure		401	{object}	shared.ErrorResponse
//	@Failure		500	{object}	shared.ErrorResponse
//	@Router			/api/v1/user/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// 获取Token（从Authorization header中）
	token := c.GetHeader("Authorization")
	if token == "" {
		// 即使没有token也返回成功，因为登出应该是幂等的
		shared.Success(c, http.StatusOK, "登出成功", gin.H{
			"message": "Logged out successfully",
		})
		return
	}

	// 去除 "Bearer " 前缀
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	// TODO: 将Token加入黑名单或清除服务端会话
	// 当前JWT是无状态的，主要依赖客户端删除token
	// 如果需要服务端控制，可以实现token黑名单机制

	// 返回成功响应
	shared.Success(c, http.StatusOK, "登出成功", gin.H{
		"message": "Logged out successfully",
	})
}
