package user

// 验证相关常量
const (
	// 验证码相关
	VerificationCodeLength     = 6           // 验证码长度
	VerificationCodeExpiry     = 30 * 60     // 验证码有效期：30分钟（秒）
	VerificationCodeExpiryMin  = 30          // 验证码有效期：30分钟
	VerificationCodeExpirySec  = 1800        // 验证码有效期：1800秒

	// 频率限制相关
	VerificationRateLimitCount   = 3         // 验证码发送频率限制：每分钟最多3次
	VerificationRateLimitWindow  = 60        // 频率限制时间窗口：60秒

	// 密码重置相关
	PasswordResetTokenExpiry     = 3600       // 密码重置Token有效期：1小时（秒）
	PasswordResetTokenExpiryHour = 1          // 密码重置Token有效期：1小时

	// Token相关
	TokenDefaultExpiry           = 3600 * 24 * 7 // Token默认有效期：7天
)

// 验证目的常量
const (
	VerificationPurposeEmail   = "verify_email"   // 邮箱验证
	VerificationPurposePhone   = "verify_phone"   // 手机验证
	VerificationPurposeReset   = "reset_password" // 密码重置
)
