package shared

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/service/shared/auth"
)

// AuthAPI 认证服务API处理器
type AuthAPI struct {
	authService auth.AuthService
}

// NewAuthAPI 创建认证API实例
func NewAuthAPI(authService auth.AuthService) *AuthAPI {
	return &AuthAPI{
		authService: authService,
	}
}

// Register 用户注册
// @Summary 用户注册
// @Description 注册新用户账号
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body auth.RegisterRequest true "注册信息"
// @Success 200 {object} APIResponse{data=auth.RegisterResponse}
// @Failure 400 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /api/v1/shared/auth/register [post]
func (api *AuthAPI) Register(c *gin.Context) {
	var req auth.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	resp, err := api.authService.Register(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "注册失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "注册成功",
		Data:    resp,
	})
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录获取Token
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body auth.LoginRequest true "登录信息"
// @Success 200 {object} APIResponse{data=auth.LoginResponse}
// @Failure 400 {object} APIResponse
// @Failure 401 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /api/v1/shared/auth/login [post]
func (api *AuthAPI) Login(c *gin.Context) {
	var req auth.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	resp, err := api.authService.Login(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Code:    401,
			Message: "登录失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "登录成功",
		Data:    resp,
	})
}

// Logout 用户登出
// @Summary 用户登出
// @Description 用户登出，使Token失效
// @Tags 认证
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} APIResponse
// @Failure 401 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /api/v1/shared/auth/logout [post]
func (api *AuthAPI) Logout(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Code:    401,
			Message: "未提供Token",
		})
		return
	}

	// 去除 "Bearer " 前缀
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	err := api.authService.Logout(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "登出失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "登出成功",
	})
}

// RefreshToken 刷新Token
// @Summary 刷新Token
// @Description 使用当前Token获取新Token
// @Tags 认证
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} APIResponse{data=string}
// @Failure 401 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /api/v1/shared/auth/refresh [post]
func (api *AuthAPI) RefreshToken(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Code:    401,
			Message: "未提供Token",
		})
		return
	}

	// 去除 "Bearer " 前缀
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	newToken, err := api.authService.RefreshToken(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Code:    401,
			Message: "Token刷新失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "Token刷新成功",
		Data:    map[string]string{"token": newToken},
	})
}

// GetUserPermissions 获取用户权限
// @Summary 获取用户权限
// @Description 获取当前用户的权限列表
// @Tags 认证
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} APIResponse{data=[]string}
// @Failure 401 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /api/v1/shared/auth/permissions [get]
func (api *AuthAPI) GetUserPermissions(c *gin.Context) {
	// 从Context中获取当前用户ID（由中间件设置）
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Code:    401,
			Message: "未认证",
		})
		return
	}

	permissions, err := api.authService.GetUserPermissions(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "获取权限失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "获取权限成功",
		Data:    permissions,
	})
}

// GetUserRoles 获取用户角色
// @Summary 获取用户角色
// @Description 获取当前用户的角色列表
// @Tags 认证
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} APIResponse{data=[]string}
// @Failure 401 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /api/v1/shared/auth/roles [get]
func (api *AuthAPI) GetUserRoles(c *gin.Context) {
	// 从Context中获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Code:    401,
			Message: "未认证",
		})
		return
	}

	roles, err := api.authService.GetUserRoles(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "获取角色失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "获取角色成功",
		Data:    roles,
	})
}
