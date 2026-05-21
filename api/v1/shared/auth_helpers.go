package shared

import (
	"context"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/pkg/emailcode"
	"Qingyu_backend/pkg/response"
	authservice "Qingyu_backend/service/auth"
)

// ValidateRegisterVerificationCode 校验注册邮箱验证码。
// 返回 false 表示已写入错误响应。
func ValidateRegisterVerificationCode(c *gin.Context, codeManager *emailcode.Manager, email, verificationCode string) bool {
	if codeManager == nil || !codeManager.Enabled() {
		return true
	}

	if verificationCode == "" {
		response.BadRequest(c, "请先填写邮箱验证码", nil)
		return false
	}

	if err := codeManager.VerifyRegisterCode(email, verificationCode); err != nil {
		response.BadRequest(c, "邮箱验证码校验失败: "+err.Error(), nil)
		return false
	}

	return true
}

// SendRegisterVerificationCode 发送注册邮箱验证码。
// 返回 false 表示已写入错误响应。
func SendRegisterVerificationCode(ctx context.Context, c *gin.Context, codeManager *emailcode.Manager, email string) bool {
	if codeManager == nil {
		response.BadRequest(c, "邮箱验证码功能未启用，请先配置 QINGYU_EMAIL_ENABLED 和 SMTP 参数", nil)
		return false
	}

	if err := codeManager.SendRegisterCode(ctx, email); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return false
	}

	return true
}

// HandleAuthLogin 处理标准认证登录请求。
func HandleAuthLogin(c *gin.Context, authService authservice.AuthService) {
	var req authservice.LoginRequest
	if !BindJSONWithMessage(c, &req, "请求参数错误: ") {
		return
	}

	resp, err := authService.Login(c.Request.Context(), &req)
	if err != nil {
		response.Unauthorized(c, "登录失败: "+err.Error())
		return
	}

	response.SuccessWithMessage(c, "登录成功", resp)
}
