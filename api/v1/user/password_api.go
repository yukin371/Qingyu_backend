package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

	shared "Qingyu_backend/api/v1/shared"
	"Qingyu_backend/api/v1/user/dto"
	userService "Qingyu_backend/service/user"
)

// PasswordAPI 密码API处理器
type PasswordAPI struct {
	passwordService *userService.PasswordService
}

// NewPasswordAPI 创建密码API实例
func NewPasswordAPI(
	passwordService *userService.PasswordService,
) *PasswordAPI {
	return &PasswordAPI{
		passwordService: passwordService,
	}
}

// SendPasswordResetCode 发送密码重置验证码
//
//	@Summary		发送密码重置验证码
//	@Description	向用户邮箱发送密码重置验证码
//	@Tags			User Password
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.SendPasswordResetRequest	true	"发送重置验证码请求"
//	@Success 200 {object} shared.APIResponse{data=dto.SendPasswordResetResponse}
//	@Failure		400		{object}	shared.APIResponse	"参数错误"
//	@Router			/api/v1/users/password/reset/send [post]
func (api *PasswordAPI) SendPasswordResetCode(c *gin.Context) {
	var req dto.SendPasswordResetRequest
	if !shared.ValidateRequest(c, &req) {
		return
	}

	// 发送重置验证码
	err := api.passwordService.SendResetCode(c.Request.Context(), req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, shared.ErrorResponseWithCode(
			500,
			"发送失败",
			err,
		))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(
		dto.SendPasswordResetResponse{
			ExpiresIn: 300,
			Message:   "重置验证码已发送到您的邮箱",
		},
		"success",
	))
}

// ResetPassword 重置密码
//
//	@Summary		重置密码
//	@Description	验证重置码并重置用户密码
//	@Tags			User Password
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.ResetPasswordRequest	true	"重置密码请求"
//	@Success 200 {object} shared.APIResponse{data=dto.ResetPasswordResponse}
//	@Failure		400		{object}	shared.APIResponse	"参数错误"
//	@Failure		400		{object}	shared.APIResponse	"验证码无效"
//	@Router			/api/v1/users/password/reset/verify [post]
func (api *PasswordAPI) ResetPassword(c *gin.Context) {
	var req dto.ResetPasswordRequest
	if !shared.ValidateRequest(c, &req) {
		return
	}

	// 重置密码
	err := api.passwordService.ResetPassword(
		c.Request.Context(),
		req.Email,
		req.Code,
		req.NewPassword,
	)
	if err != nil {
		if err == userService.ErrInvalidCode {
			c.JSON(http.StatusBadRequest, shared.ErrorResponseWithCode(
				400,
				"验证码无效或已过期",
				err,
			))
		} else {
			c.JSON(http.StatusInternalServerError, shared.ErrorResponseWithCode(
				500,
				"重置失败",
				err,
			))
		}
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(
		dto.ResetPasswordResponse{
			Success: true,
			Message:  "密码重置成功",
		},
		"success",
	))
}

// UpdatePassword 修改密码
//
//	@Summary		修改密码
//	@Description	修改当前用户密码（需要验证旧密码）
//	@Tags			User Password
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		dto.UpdatePasswordRequest	true	"修改密码请求"
//	@Success 200 {object} shared.APIResponse
//	@Failure		400		{object}	shared.APIResponse	"参数错误"
//	@Failure		401		{object}	shared.APIResponse	"旧密码错误"
//	@Router			/api/v1/users/password [put]
func (api *PasswordAPI) UpdatePassword(c *gin.Context) {
	userID := c.GetString("userID")

	var req dto.UpdatePasswordRequest
	if !shared.ValidateRequest(c, &req) {
		return
	}

	// 修改密码
	err := api.passwordService.UpdatePassword(
		c.Request.Context(),
		userID,
		req.OldPassword,
		req.NewPassword,
	)
	if err != nil {
		if err == userService.ErrOldPasswordMismatch {
			c.JSON(http.StatusUnauthorized, shared.ErrorResponseWithCode(
				401,
				"旧密码错误",
				err,
			))
		} else {
			c.JSON(http.StatusInternalServerError, shared.ErrorResponseWithCode(
				500,
				"修改失败",
				err,
			))
		}
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(nil, "success"))
}
