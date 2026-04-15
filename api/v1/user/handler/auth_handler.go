package handler

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/shared"
	"Qingyu_backend/api/v1/user/dto"
	"Qingyu_backend/pkg/response"
	"Qingyu_backend/pkg/utils"
	authsvc "Qingyu_backend/service/auth"
	serviceInterfaces "Qingyu_backend/service/interfaces/base"
)

// AuthHandler 认证处理器
type authHandlerService interface {
	Register(ctx context.Context, req *authsvc.RegisterRequest) (*authsvc.RegisterResponse, error)
	Login(ctx context.Context, req *authsvc.LoginRequest) (*authsvc.LoginResponse, error)
	Logout(ctx context.Context, token string) error
}

type AuthHandler struct {
	authService authHandlerService
}

// NewAuthHandler 创建认证处理器实例
func NewAuthHandler(authService authHandlerService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
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
//	@Success		200		{object}	response.APIResponse{data=dto.RegisterResponse}
//	@Failure		400		{object}	response.APIResponse
//	@Failure		500		{object}	response.APIResponse
//	@Router			/api/v1/user/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if !shared.ValidateRequest(c, &req) {
		return
	}

	// 调用Service层
	serviceReq := &authsvc.RegisterRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	resp, err := h.authService.Register(c.Request.Context(), serviceReq)
	if err != nil {
		// 根据错误类型返回不同的HTTP状态码
		if serviceErr, ok := err.(*serviceInterfaces.ServiceError); ok {
			switch serviceErr.Type {
			case serviceInterfaces.ErrorTypeValidation:
				response.BadRequest(c, "注册失败", serviceErr.Message)
			case serviceInterfaces.ErrorTypeBusiness:
				// 根据错误消息返回具体的错误码
				if serviceErr.Message == "用户名已存在" {
					c.JSON(http.StatusConflict, response.APIResponse{
						Code:      2003,
						Message:   "用户名已被注册",
						Timestamp: time.Now().UnixMilli(),
					})
				} else if serviceErr.Message == "邮箱已存在" {
					c.JSON(http.StatusConflict, response.APIResponse{
						Code:      2004,
						Message:   "邮箱已被注册",
						Timestamp: time.Now().UnixMilli(),
					})
				} else {
					response.BadRequest(c, "注册失败", serviceErr.Message)
				}
			default:
				response.InternalError(c, err)
			}
			return
		}
		response.InternalError(c, err)
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

	response.Created(c, registerResp)
}

// Login 用户登录
//
//	@Summary		用户登录
//	@Description	用户登录获取Token
//	@Tags			用户管理-认证
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.LoginRequest	true	"登录信息"
//	@Success		200		{object}	response.APIResponse{data=dto.LoginResponse}
//	@Failure		400		{object}	response.APIResponse
//	@Failure		401		{object}	response.APIResponse
//	@Failure		500		{object}	response.APIResponse
//	@Router			/api/v1/user/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if !shared.ValidateRequest(c, &req) {
		return
	}

	// 获取客户端IP
	clientIP := utils.GetClientIP(c)

	// 调用Service层
	serviceReq := &authsvc.LoginRequest{
		Username: req.Username,
		Password: req.Password,
		ClientIP: clientIP,
	}

	resp, err := h.authService.Login(c.Request.Context(), serviceReq)
	if err != nil {
		if serviceErr, ok := err.(*serviceInterfaces.ServiceError); ok {
			switch serviceErr.Type {
			case serviceInterfaces.ErrorTypeNotFound:
				response.Unauthorized(c, "用户名或密码错误")
			case serviceInterfaces.ErrorTypeUnauthorized:
				response.Unauthorized(c, "用户名或密码错误")
			case serviceInterfaces.ErrorTypeValidation:
				response.BadRequest(c, "登录失败", serviceErr.Message)
			default:
				response.InternalError(c, err)
			}
			return
		}
		response.InternalError(c, err)
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

	response.Success(c, loginResp)
}

// Logout 用户登出
//
//	@Summary		用户登出
//	@Description	用户登出，清除服务端会话/Token
//	@Tags			用户管理-认证
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success		200	{object}	response.APIResponse
//	@Failure		401	{object}	response.APIResponse
//	@Failure		500	{object}	response.APIResponse
//	@Router			/api/v1/user/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// 获取Token（从Authorization header中）
	token := c.GetHeader("Authorization")
	if token == "" {
		// 即使没有token也返回成功，因为登出应该是幂等的
		response.Success(c, gin.H{
			"message": "Logged out successfully",
		})
		return
	}

	token = strings.TrimSpace(strings.TrimPrefix(token, "Bearer "))

	if err := h.authService.Logout(c.Request.Context(), token); err != nil {
		response.InternalError(c, err)
		return
	}

	// 返回成功响应
	response.Success(c, gin.H{
		"message": "Logged out successfully",
		"success": true,
	})
}
