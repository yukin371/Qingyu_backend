package shared

import (
	"github.com/gin-gonic/gin"

	"Qingyu_backend/pkg/emailcode"
	"Qingyu_backend/pkg/response"
	"Qingyu_backend/service/auth"
)

// AuthAPI 认证服务API处理器
type AuthAPI struct {
	authService auth.AuthService
	codeManager *emailcode.Manager
}

// NewAuthAPI 创建认证API实例
func NewAuthAPI(authService auth.AuthService) *AuthAPI {
	return &AuthAPI{
		authService: authService,
		codeManager: emailcode.NewManager(),
	}
}

// Register 用户注册
//
//	@Summary		用户注册
//	@Description	注册新用户账号
//	@Tags			认证
//	@Accept			json
//	@Produce		json
//	@Param			request	body		auth.RegisterRequest	true	"注册信息"
//	@Success 200 {object} APIResponse
//	@Failure		400		{object}	APIResponse
//	@Failure		500		{object}	APIResponse
//	@Router			/api/v1/shared/auth/register [post]
func (api *AuthAPI) Register(c *gin.Context) {
	var req auth.RegisterRequest
	if !ValidateRequest(c, &req) {
		return
	}

	if api.codeManager.Enabled() {
		if req.VerificationCode == "" {
			response.BadRequest(c, "请先填写邮箱验证码", nil)
			return
		}
		if err := api.codeManager.VerifyRegisterCode(req.Email, req.VerificationCode); err != nil {
			response.BadRequest(c, "邮箱验证码校验失败: "+err.Error(), nil)
			return
		}
	}

	resp, err := api.authService.Register(c.Request.Context(), &req)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.SuccessWithMessage(c, "注册成功", resp)
}

// Login 用户登录
//
//	@Summary		用户登录
//	@Description	用户登录获取Token
//	@Tags			认证
//	@Accept			json
//	@Produce		json
//	@Param			request	body		auth.LoginRequest	true	"登录信息"
//	@Success 200 {object} APIResponse
//	@Failure		400		{object}	APIResponse
//	@Failure		401		{object}	APIResponse
//	@Failure		500		{object}	APIResponse
//	@Router			/api/v1/shared/auth/login [post]
func (api *AuthAPI) Login(c *gin.Context) {
	var req auth.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error(), nil)
		return
	}

	resp, err := api.authService.Login(c.Request.Context(), &req)
	if err != nil {
		response.Unauthorized(c, "登录失败: "+err.Error())
		return
	}

	response.SuccessWithMessage(c, "登录成功", resp)
}

// Logout 用户登出
//
//	@Summary		用户登出
//	@Description	用户登出，使Token失效
//	@Tags			认证
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success		200	{object}	APIResponse
//	@Failure		401	{object}	APIResponse
//	@Failure		500	{object}	APIResponse
//	@Router			/api/v1/shared/auth/logout [post]
func (api *AuthAPI) Logout(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		response.Unauthorized(c, "未提供Token")
		return
	}

	// 去除 "Bearer " 前缀
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	err := api.authService.Logout(c.Request.Context(), token)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.SuccessWithMessage(c, "登出成功", nil)
}

// RefreshToken 刷新Token
//
//	@Summary		刷新Token
//	@Description	使用当前Token获取新Token
//	@Tags			认证
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success 200 {object} APIResponse
//	@Failure		401	{object}	APIResponse
//	@Failure		500	{object}	APIResponse
//	@Router			/api/v1/shared/auth/refresh [post]
func (api *AuthAPI) RefreshToken(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		response.Unauthorized(c, "未提供Token")
		return
	}

	// 去除 "Bearer " 前缀
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	newToken, err := api.authService.RefreshToken(c.Request.Context(), token)
	if err != nil {
		response.Unauthorized(c, "Token刷新失败: "+err.Error())
		return
	}

	response.SuccessWithMessage(c, "Token刷新成功", map[string]string{"token": newToken})
}

// GetUserPermissions 获取用户权限
//
//	@Summary		获取用户权限
//	@Description	获取当前用户的权限列表
//	@Tags			认证
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success 200 {object} APIResponse
//	@Failure		401	{object}	APIResponse
//	@Failure		500	{object}	APIResponse
//	@Router			/api/v1/shared/auth/permissions [get]
func (api *AuthAPI) GetUserPermissions(c *gin.Context) {
	// 从Context中获取当前用户ID（由中间件设置）
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未认证")
		return
	}

	permissions, err := api.authService.GetUserPermissions(c.Request.Context(), userID.(string))
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.SuccessWithMessage(c, "获取权限成功", permissions)
}

// GetUserRoles 获取用户角色
//
//	@Summary		获取用户角色
//	@Description	获取当前用户的角色列表
//	@Tags			认证
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success 200 {object} APIResponse
//	@Failure		401	{object}	APIResponse
//	@Failure		500	{object}	APIResponse
//	@Router			/api/v1/shared/auth/roles [get]
func (api *AuthAPI) GetUserRoles(c *gin.Context) {
	// 从Context中获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未认证")
		return
	}

	roles, err := api.authService.GetUserRoles(c.Request.Context(), userID.(string))
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.SuccessWithMessage(c, "获取角色成功", roles)
}

type sendVerificationCodeRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// SendVerificationCode 发送注册邮箱验证码
//
//	@Summary		发送邮箱验证码
//	@Description	发送用户注册验证码到邮箱
//	@Tags			认证
//	@Accept			json
//	@Produce		json
//	@Param			request	body		sendVerificationCodeRequest	true	"邮箱信息"
//	@Success 200 {object} APIResponse
//	@Failure		400		{object}	APIResponse
//	@Failure		500		{object}	APIResponse
//	@Router			/api/v1/shared/auth/send-verification-code [post]
func (api *AuthAPI) SendVerificationCode(c *gin.Context) {
	var req sendVerificationCodeRequest
	if !ValidateRequest(c, &req) {
		return
	}

	if err := api.codeManager.SendRegisterCode(c.Request.Context(), req.Email); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.SuccessWithMessage(c, "验证码发送成功", map[string]interface{}{
		"expires_in_seconds":  600,
		"cooldown_in_seconds": 60,
	})
}
