package handler

import (
	serviceInterfaces "Qingyu_backend/service/interfaces/base"
	userServiceInterface "Qingyu_backend/service/interfaces/user"
	"net/http"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/shared"
	"Qingyu_backend/api/v1/usermanagement/dto"
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
//	@Description	注册新用户账号
//	@Tags			用户管理-认证
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.RegisterRequest	true	"注册信息"
//	@Success		200		{object}	shared.APIResponse{data=dto.RegisterResponse}
//	@Failure		400		{object}	shared.ErrorResponse
//	@Failure		500		{object}	shared.ErrorResponse
//	@Router			/api/v1/user-management/auth/register [post]
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
	registerResp := dto.RegisterResponse{
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
//	@Tags			用户管理-认证
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.LoginRequest	true	"登录信息"
//	@Success		200		{object}	shared.APIResponse{data=dto.LoginResponse}
//	@Failure		400		{object}	shared.ErrorResponse
//	@Failure		401		{object}	shared.ErrorResponse
//	@Failure		500		{object}	shared.ErrorResponse
//	@Router			/api/v1/user-management/auth/login [post]
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
	loginResp := dto.LoginResponse{
		Token: resp.Token,
		User: dto.UserBasicInfo{
			UserID:   resp.User.ID,
			Username: resp.User.Username,
			Email:    resp.User.Email,
			Role:     resp.User.Role,
		},
	}

	shared.Success(c, http.StatusOK, "登录成功", loginResp)
}
