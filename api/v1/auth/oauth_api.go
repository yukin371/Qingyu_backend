package auth

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"Qingyu_backend/api/v1/shared"
	authModel "Qingyu_backend/models/auth"
	"Qingyu_backend/pkg/response"
	"Qingyu_backend/service/shared/auth"
)

// OAuthAPI OAuth认证API处理器
type OAuthAPI struct {
	oauthService auth.OAuthServiceInterface
	authService  auth.AuthService
	logger       *zap.Logger
}

// NewOAuthAPI 创建OAuth API实例
func NewOAuthAPI(oauthService auth.OAuthServiceInterface, authService auth.AuthService, logger *zap.Logger) *OAuthAPI {
	return &OAuthAPI{
		oauthService: oauthService,
		authService:  authService,
		logger:       logger,
	}
}

// GetAuthorizeURL 获取OAuth授权URL
//
//	@Summary		获取OAuth授权URL
//	@Description	获取第三方登录授权URL，用户需要访问此URL进行授权
//	@Tags			OAuth认证
//	@Accept			json
//	@Produce		json
//	@Param			provider	path		string				true	"OAuth提供商 (google/github/qq)"
//	@Param			request		body		OAuthAuthorizeRequest	true	"授权请求"
//	@Success		200			{object} shared.APIResponse
//	@Failure		400			{object} shared.APIResponse
//	@Failure		500			{object} shared.APIResponse
//	@Router			/api/v1/shared/oauth/{provider}/authorize [post]
func (api *OAuthAPI) GetAuthorizeURL(c *gin.Context) {
	provider := authModel.OAuthProvider(c.Param("provider"))

	var req OAuthAuthorizeRequest
	if !shared.ValidateRequest(c, &req) {
		return
	}

	// 检查是否为绑定模式（已登录用户绑定OAuth账号）
	linkMode := false
	userID := ""
	if userIDVal, exists := c.Get("user_id"); exists {
		linkMode = true
		userID = userIDVal.(string)
	}

	authURL, err := api.oauthService.GetAuthURL(c.Request.Context(), provider, req.RedirectURI, req.State, linkMode, userID)
	if err != nil {
		response.BadRequest(c, "获取授权URL失败: "+err.Error(), nil)
		return
	}

	response.SuccessWithMessage(c, "获取授权URL成功", gin.H{
		"authorize_url": authURL,
		"provider":      provider,
	})
}

// HandleCallback 处理OAuth回调
//
//	@Summary		处理OAuth回调
//	@Description	处理第三方登录回调，完成登录或注册
//	@Tags			OAuth认证
//	@Accept			json
//	@Produce		json
//	@Param			provider	path		string				true	"OAuth提供商 (google/github/qq)"
//	@Param			request		body		OAuthCallbackRequest	true	"回调请求"
//	@Success		200			{object} shared.APIResponse
//	@Failure		400			{object} shared.APIResponse
//	@Failure		401			{object} shared.APIResponse
//	@Failure		500			{object} shared.APIResponse
//	@Router			/api/v1/shared/oauth/{provider}/callback [post]
func (api *OAuthAPI) HandleCallback(c *gin.Context) {
	provider := authModel.OAuthProvider(c.Param("provider"))

	var req OAuthCallbackRequest
	if !shared.ValidateRequest(c, &req) {
		return
	}

	// 交换授权码获取Token
	token, session, err := api.oauthService.ExchangeCode(c.Request.Context(), provider, req.Code, req.State)
	if err != nil {
		response.BadRequest(c, "授权码交换失败: "+err.Error(), nil)
		return
	}

	// 获取用户信息
	identity, err := api.oauthService.GetUserInfo(c.Request.Context(), provider, token)
	if err != nil {
		response.BadRequest(c, "获取用户信息失败: "+err.Error(), nil)
		return
	}

	// 检查是否为绑定模式
	if session.LinkMode && session.UserID != "" {
		// 绑定模式：将OAuth账号绑定到已登录用户
		account, err := api.oauthService.LinkAccount(c.Request.Context(), session.UserID, provider, token, identity)
		if err != nil {
			response.BadRequest(c, "绑定账号失败: "+err.Error(), nil)
			return
		}

		response.SuccessWithMessage(c, "账号绑定成功", gin.H{
			"account":  account,
			"provider": provider,
		})
		return
	}

	// 登录/注册模式：使用AuthService的OAuthLogin方法
	oauthLoginReq := &auth.OAuthLoginRequest{
		Provider:   provider,
		ProviderID: identity.ProviderID,
		Email:      identity.Email,
		Name:       identity.Name,
		Avatar:     identity.Avatar,
		Username:   identity.Username,
	}

	loginResp, err := api.authService.OAuthLogin(c.Request.Context(), oauthLoginReq)
	if err != nil {
		response.BadRequest(c, "OAuth登录失败: "+err.Error(), nil)
		return
	}

	response.SuccessWithMessage(c, "OAuth登录成功", loginResp)
}

// GetLinkedAccounts 获取用户绑定的OAuth账号列表
//
//	@Summary		获取绑定的OAuth账号
//	@Description	获取当前用户绑定的所有第三方账号
//	@Tags			OAuth认证
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success		200	{object} shared.APIResponse
//	@Failure		401	{object} shared.APIResponse
//	@Failure		500	{object} shared.APIResponse
//	@Router			/api/v1/shared/oauth/accounts [get]
func (api *OAuthAPI) GetLinkedAccounts(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未授权")
		return
	}

	accounts, err := api.oauthService.GetLinkedAccounts(c.Request.Context(), userID.(string))
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.SuccessWithMessage(c, "获取成功", accounts)
}

// UnlinkAccount 解绑OAuth账号
//
//	@Summary		解绑OAuth账号
//	@Description	解绑指定的第三方账号
//	@Tags			OAuth认证
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			accountID	path		string	true	"OAuth账号ID"
//	@Success		200			{object} shared.APIResponse
//	@Failure		400			{object} shared.APIResponse
//	@Failure		401			{object} shared.APIResponse
//	@Failure		500			{object} shared.APIResponse
//	@Router			/api/v1/shared/oauth/accounts/{accountID} [delete]
func (api *OAuthAPI) UnlinkAccount(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未授权")
		return
	}

	accountID := c.Param("accountID")
	if accountID == "" {
		response.BadRequest(c, "账号ID不能为空", nil)
		return
	}

	err := api.oauthService.UnlinkAccount(c.Request.Context(), userID.(string), accountID)
	if err != nil {
		response.BadRequest(c, "解绑账号失败: "+err.Error(), nil)
		return
	}

	response.SuccessWithMessage(c, "解绑成功", nil)
}

// SetPrimaryAccount 设置主账号
//
//	@Summary		设置主账号
//	@Description	将指定的OAuth账号设为主账号
//	@Tags			OAuth认证
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			accountID	path		string	true	"OAuth账号ID"
//	@Success		200			{object} shared.APIResponse
//	@Failure		400			{object} shared.APIResponse
//	@Failure		401			{object} shared.APIResponse
//	@Failure		500			{object} shared.APIResponse
//	@Router			/api/v1/shared/oauth/accounts/{accountID}/primary [put]
func (api *OAuthAPI) SetPrimaryAccount(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未授权")
		return
	}

	accountID := c.Param("accountID")
	if accountID == "" {
		response.BadRequest(c, "账号ID不能为空", nil)
		return
	}

	err := api.oauthService.SetPrimaryAccount(c.Request.Context(), userID.(string), accountID)
	if err != nil {
		response.BadRequest(c, "设置主账号失败: "+err.Error(), nil)
		return
	}

	response.SuccessWithMessage(c, "设置成功", nil)
}

// ==================== 请求和响应结构 ====================

// OAuthAuthorizeRequest OAuth授权请求
type OAuthAuthorizeRequest struct {
	RedirectURI string `json:"redirect_uri" binding:"required"` // 回调地址
	State       string `json:"state"`                           // 状态参数，用于防止CSRF攻击
}

// OAuthCallbackRequest OAuth回调请求
type OAuthCallbackRequest struct {
	Code  string `json:"code" binding:"required"`  // 授权码
	State string `json:"state" binding:"required"` // 状态参数
}
