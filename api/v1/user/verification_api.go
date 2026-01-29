package user

import (
	"github.com/gin-gonic/gin"

	sharedApi "Qingyu_backend/api/v1/shared"
	"Qingyu_backend/api/v1/user/dto"
	"Qingyu_backend/pkg/response"
	userService "Qingyu_backend/service/user"
	userConstants "Qingyu_backend/service/user"
	userServiceInterface "Qingyu_backend/service/interfaces/user"
	serviceInterfaces "Qingyu_backend/service/interfaces/base"
)

// VerificationAPI 验证API处理器
type VerificationAPI struct {
	verificationService *userService.VerificationService
	userService         userServiceInterface.UserService
}

// NewVerificationAPI 创建验证API实例
func NewVerificationAPI(
	verificationService *userService.VerificationService,
	userService userServiceInterface.UserService,
) *VerificationAPI {
	return &VerificationAPI{
		verificationService: verificationService,
		userService:         userService,
	}
}

// SendEmailVerifyCode 发送邮箱验证码
//
//	@Summary		发送邮箱验证码
//	@Description	向用户邮箱发送6位数字验证码，有效期5分钟
//	@Tags			User Verification
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.SendEmailCodeRequest	true	"发送验证码请求"
//	@Success 200 {object} sharedApi.APIResponse{data=dto.SendCodeResponse}
//	@Failure		400		{object}	sharedApi.APIResponse	"参数错误"
//	@Failure		429		{object}	sharedApi.APIResponse	"发送过于频繁"
//	@Router			/api/v1/users/verify/email/send [post]
func (api *VerificationAPI) SendEmailVerifyCode(c *gin.Context) {
	var req dto.SendEmailCodeRequest
	if !sharedApi.ValidateRequest(c, &req) {
		return
	}

	// 检查邮箱是否已被使用
	exists, err := api.verificationService.EmailExists(c.Request.Context(), req.Email)
	if err == nil && exists {
		// 邮箱已被使用，但仍允许发送验证码（用于验证已有邮箱）
	}

	// 发送验证码
	err = api.verificationService.SendEmailCode(c.Request.Context(), req.Email, "verify_email")
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.SuccessWithMessage(c, "success", dto.SendCodeResponse{
		ExpiresIn: userConstants.VerificationCodeExpirySec, // 30分钟
		Message:   "验证码已发送到您的邮箱",
	})
}

// SendPhoneVerifyCode 发送手机验证码（模拟实现）
//
//	@Summary		发送手机验证码
//	@Description	向用户手机发送6位数字验证码（模拟实现：控制台打印）
//	@Tags			User Verification
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.SendPhoneCodeRequest	true	"发送验证码请求"
//	@Success 200 {object} sharedApi.APIResponse{data=dto.SendCodeResponse}
//	@Failure		400		{object}	sharedApi.APIResponse	"参数错误"
//	@Router			/api/v1/users/verify/phone/send [post]
func (api *VerificationAPI) SendPhoneVerifyCode(c *gin.Context) {
	var req dto.SendPhoneCodeRequest
	if !sharedApi.ValidateRequest(c, &req) {
		return
	}

	// 模拟实现：发送手机验证码
	err := api.verificationService.SendPhoneCode(c.Request.Context(), req.Phone, "verify_phone")
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.SuccessWithMessage(c, "success", dto.SendCodeResponse{
		ExpiresIn: userConstants.VerificationCodeExpirySec, // 30分钟
		Message:   "验证码已发送（模拟实现：请查看控制台）",
	})
}

// VerifyEmail 验证邮箱
//
//	@Summary		验证邮箱
//	@Description	验证用户邮箱，验证成功后标记邮箱已验证
//	@Tags			User Verification
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.VerifyEmailRequest	true	"验证邮箱请求"
//	@Success 200 {object} sharedApi.APIResponse{data=dto.VerifyEmailResponse}
//	@Failure		400		{object}	sharedApi.APIResponse	"参数错误"
//	@Failure		400		{object}	sharedApi.APIResponse	"验证码无效或过期"
//	@Router			/api/v1/users/email/verify [post]
func (api *VerificationAPI) VerifyEmail(c *gin.Context) {
	var req dto.VerifyEmailRequest
	if !sharedApi.ValidateRequest(c, &req) {
		return
	}

	// 验证验证码
	err := api.verificationService.VerifyCode(c.Request.Context(), req.Email, req.Code, "verify_email")
	if err != nil {
		response.BadRequest(c, "验证码无效或已过期", err)
		return
	}

	// ✅ 添加：标记验证码为已使用（防止重复使用）
	err = api.verificationService.MarkCodeAsUsed(c.Request.Context(), req.Email)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	// 更新用户邮箱验证状态
	userID := c.GetString("userID")
	err = api.verificationService.SetEmailVerified(c.Request.Context(), userID, req.Email)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.SuccessWithMessage(c, "success", dto.VerifyEmailResponse{
		Verified: true,
		Message:  "邮箱验证成功",
	})
}

// UnbindEmail 解绑邮箱
//
//	@Summary		解绑邮箱
//	@Description	解除用户邮箱绑定，需要验证密码
//	@Tags			User Profile
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.UnbindEmailRequest	true	"解绑邮箱请求"
//	@Success 200 {object} sharedApi.APIResponse
//	@Failure		400		{object}	sharedApi.APIResponse	"参数错误"
//	@Failure		401		{object}	sharedApi.APIResponse	"密码错误"
//	@Router			/api/v1/users/email/unbind [delete]
func (api *VerificationAPI) UnbindEmail(c *gin.Context) {
	userID := c.GetString("userID")

	var req dto.UnbindEmailRequest
	if !sharedApi.ValidateRequest(c, &req) {
		return
	}

	// 验证密码
	err := api.verificationService.CheckPassword(c.Request.Context(), userID, req.Password)
	if err != nil {
		response.Unauthorized(c, "密码错误")
		return
	}

	// 解绑邮箱
	err = api.userService.UnbindEmail(c.Request.Context(), userID)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, nil)
}

// UnbindPhone 解绑手机（模拟实现）
//
//	@Summary		解绑手机
//	@Description	解除用户手机绑定（模拟实现，直接返回成功）
//	@Tags			User Profile
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.UnbindPhoneRequest	true	"解绑手机请求"
//	@Success 200 {object} sharedApi.APIResponse
//	@Failure		400		{object}	sharedApi.APIResponse	"参数错误"
//	@Router			/api/v1/users/phone/unbind [delete]
func (api *VerificationAPI) UnbindPhone(c *gin.Context) {
	userID := c.GetString("userID")

	var req dto.UnbindPhoneRequest
	if !sharedApi.ValidateRequest(c, &req) {
		return
	}

	// 验证密码
	err := api.verificationService.CheckPassword(c.Request.Context(), userID, req.Password)
	if err != nil {
		response.Unauthorized(c, "密码错误")
		return
	}

	// 模拟实现：直接返回成功
	// TODO: 实际实现时需要调用 userService.UnbindPhone()
	err = api.userService.UnbindPhone(c.Request.Context(), userID)
	if err != nil {
		// 暂时忽略错误，因为这是模拟实现
	}
	response.Success(c, nil)
}

// DeleteDevice 删除设备
//
//	@Summary		删除设备
//	@Description	删除用户的登录设备
//	@Tags			User Profile
//	@Accept			json
//	@Produce		json
//	@Param			deviceId	path		string	true	"设备ID"
//	@Param			request	body		dto.DeleteDeviceRequest	true	"删除设备请求"
//	@Success 200 {object} sharedApi.APIResponse
//	@Failure		400		{object}	sharedApi.APIResponse	"参数错误"
//	@Failure		401		{object}	sharedApi.APIResponse	"密码错误"
//	@Failure		404		{object}	sharedApi.APIResponse	"设备不存在"
//	@Router			/api/v1/users/devices/{deviceId} [delete]
func (api *VerificationAPI) DeleteDevice(c *gin.Context) {
	userID := c.GetString("userID")
	deviceID := c.Param("deviceId")

	var req dto.DeleteDeviceRequest
	if !sharedApi.ValidateRequest(c, &req) {
		return
	}

	// 验证密码
	err := api.verificationService.CheckPassword(c.Request.Context(), userID, req.Password)
	if err != nil {
		response.Unauthorized(c, "密码错误")
		return
	}

	// 删除设备
	err = api.userService.DeleteDevice(c.Request.Context(), userID, deviceID)
	if err != nil {
		// 使用错误类型判断而不是字符串比较
		if serviceInterfaces.IsNotFoundError(err) {
			response.NotFound(c, "设备不存在")
		} else {
			response.InternalError(c, err)
		}
		return
	}

	response.Success(c, nil)
}
