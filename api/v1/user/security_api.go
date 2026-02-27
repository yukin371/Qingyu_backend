package user

import (
	"github.com/gin-gonic/gin"

	shared "Qingyu_backend/api/v1/shared"
	"Qingyu_backend/pkg/response"
	user2 "Qingyu_backend/service/interfaces/user"
)

// SecurityAPI 用户安全API处理器
type SecurityAPI struct {
	userService user2.UserService
}

// NewSecurityAPI 创建安全API实例
func NewSecurityAPI(userService user2.UserService) *SecurityAPI {
	return &SecurityAPI{
		userService: userService,
	}
}

// ==================== 邮箱验证相关API ====================

// SendEmailVerification 发送邮箱验证码
//
//	@Summary		发送邮箱验证码
//	@Description	向用户邮箱发送6位数字验证码，用于验证邮箱
//	@Tags			用户安全
//	@Accept			json
//	@Produce		json
//	@Param			request	body		object	true	"发送验证码请求"
//	@Success 200 {object} shared.APIResponse
//	@Failure		400		{object}	shared.APIResponse
//	@Failure		404		{object}	shared.APIResponse
//	@Failure		500		{object}	shared.APIResponse
//	@Router			/api/v1/user/email/send-code [post]
func (api *SecurityAPI) SendEmailVerification(c *gin.Context) {
	var req user2.SendEmailVerificationRequest
	if !shared.ValidateRequest(c, &req) {
		return
	}

	resp, err := api.userService.SendEmailVerification(c.Request.Context(), &req)
	if err != nil {
		c.Error(err)
		return
	}

	response.SuccessWithMessage(c, resp.Message, resp)
}

// VerifyEmail 验证邮箱
//
//	@Summary		验证邮箱
//	@Description	使用验证码验证用户邮箱
//	@Tags			用户安全
//	@Accept			json
//	@Produce		json
//	@Param			request	body		object	true	"验证邮箱请求"
//	@Success 200 {object} shared.APIResponse
//	@Failure		400		{object}	shared.APIResponse
//	@Failure		404		{object}	shared.APIResponse
//	@Failure		500		{object}	shared.APIResponse
//	@Router			/api/v1/user/email/verify [post]
func (api *SecurityAPI) VerifyEmail(c *gin.Context) {
	var req user2.VerifyEmailRequest
	if !shared.ValidateRequest(c, &req) {
		return
	}

	resp, err := api.userService.VerifyEmail(c.Request.Context(), &req)
	if err != nil {
		response.BadRequest(c, "验证失败: "+err.Error(), nil)
		return
	}

	response.SuccessWithMessage(c, resp.Message, resp)
}

// ==================== 密码重置相关API ====================

// RequestPasswordReset 请求密码重置
//
//	@Summary		请求密码重置
//	@Description	向用户邮箱发送密码重置链接（包含Token）
//	@Tags			用户安全
//	@Accept			json
//	@Produce		json
//	@Param			request	body		object	true	"请求重置密码"
//	@Success 200 {object} shared.APIResponse
//	@Failure		400		{object}	shared.APIResponse
//	@Failure		500		{object}	shared.APIResponse
//	@Router			/api/v1/user/password/reset-request [post]
func (api *SecurityAPI) RequestPasswordReset(c *gin.Context) {
	var req user2.RequestPasswordResetRequest
	if !shared.ValidateRequest(c, &req) {
		return
	}

	resp, err := api.userService.RequestPasswordReset(c.Request.Context(), &req)
	if err != nil {
		c.Error(err)
		return
	}

	response.SuccessWithMessage(c, resp.Message, resp)
}

// ConfirmPasswordReset 确认密码重置
//
//	@Summary		确认密码重置
//	@Description	使用Token和新密码完成密码重置
//	@Tags			用户安全
//	@Accept			json
//	@Produce		json
//	@Param			request	body		object	true	"确认重置密码"
//	@Success 200 {object} shared.APIResponse
//	@Failure		400		{object}	shared.APIResponse
//	@Failure		404		{object}	shared.APIResponse
//	@Failure		500		{object}	shared.APIResponse
//	@Router			/api/v1/user/password/reset [post]
func (api *SecurityAPI) ConfirmPasswordReset(c *gin.Context) {
	var req user2.ConfirmPasswordResetRequest
	if !shared.ValidateRequest(c, &req) {
		return
	}

	resp, err := api.userService.ConfirmPasswordReset(c.Request.Context(), &req)
	if err != nil {
		response.BadRequest(c, "重置失败: "+err.Error(), nil)
		return
	}

	response.SuccessWithMessage(c, resp.Message, resp)
}
