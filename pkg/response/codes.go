package response

// 业务错误码常量定义
// 错误码格式：4位数字
// 分类规则：
//   0       - 成功
//   1xxx    - 通用客户端错误 (1000-1999)
//   2xxx    - 用户相关错误 (2000-2999)
//   3xxx    - 业务逻辑错误 (3000-3999)
//   4xxx    - 频率限制错误 (4000-4999)
//   5xxx    - 服务端错误 (5000-5999)
const (
	// 成功
	CodeSuccess = 0 // 成功

	// 通用客户端错误 (1000-1999)
	CodeParamError       = 1001 // 参数错误
	CodeUnauthorized     = 1002 // 未授权
	CodeForbidden        = 1003 // 禁止访问
	CodeNotFound         = 1004 // 资源不存在（通用）
	CodeAlreadyExists    = 1005 // 资源已存在
	CodeConflict         = 1006 // 资源冲突
	CodeInvalidOperation = 1007 // 无效操作

	// 用户相关错误 (2000-2999)
	CodeUserNotFound      = 2001 // 用户不存在
	CodeInvalidCredential = 2002 // 用户名或密码错误
	CodeEmailAlreadyUsed  = 2003 // 邮箱已被使用
	CodeEmailSendFailed   = 2004 // 邮件发送失败
	CodeInvalidCode       = 2005 // 验证码无效
	CodeCodeExpired       = 2006 // 验证码过期
	CodeTokenExpired      = 2007 // Token过期
	CodeTokenInvalid      = 2008 // Token无效
	CodePasswordTooWeak   = 2009 // 密码强度不足
	CodeAccountLocked     = 2010 // 账户已锁定
	CodeAccountDisabled   = 2011 // 账户已禁用

	// 业务逻辑错误 (3000-3999)
	CodeBookNotFound        = 3001 // 书籍不存在
	CodeChapterNotFound     = 3002 // 章节不存在
	CodeInsufficientBalance = 3003 // 余额不足
	CodeInsufficientQuota   = 3010 // 配额不足
	CodeWalletFrozen        = 3011 // 钱包已冻结
	CodeContentNotPublished = 3012 // 内容未发布
	CodeChapterLocked       = 3013 // 章节已锁定
	CodeContentPendingReview = 3014 // 内容待审核
	CodeContentRejected     = 3015 // 内容被拒绝
	CodeContentViolation    = 3016 // 内容违规

	// 评分相关错误 (2500-2599)
	CodeRatingNotFound       = 2501 // 评分不存在
	CodeRatingInvalid        = 2502 // 评分值无效（不在1-5范围）
	CodeRatingAlreadyExists  = 2503 // 用户已评分
	CodeRatingUnauthorized   = 2504 // 无权操作此评分
	CodeRatingTargetNotFound = 2505 // 评分目标不存在

	// 频率限制错误 (4000-4999)
	CodeRateLimitExceeded   = 4290 // 频率限制超出
	CodeHourlyLimitExceeded = 4291 // 小时级限制超出

	// 服务端错误 (5000-5999)
	CodeInternalError      = 5000 // 内部错误
	CodeDatabaseError      = 5001 // 数据库错误
	CodeServiceUnavailable = 5002 // 服务不可用
	CodeRedisError         = 5003 // Redis错误
	CodeExternalAPIError   = 5004 // 外部API错误
)
