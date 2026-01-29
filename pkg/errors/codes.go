package errors

import "net/http"

// ErrorCode 错误码类型
type ErrorCode int

// 常用错误码定义
// 格式：4位数字，前1位表示错误类别，后3位为具体错误编号
// 1xxx: 客户端错误 - 参数/请求相关
// 1xxx: 客户端错误 - 认证授权相关
// 5xxx: 服务器错误 - 系统内部错误
// 5xxx: 服务器错误 - 外部服务错误
const (
	// 成功
	Success ErrorCode = 0

	// 客户端错误 - 参数相关 (1xxx)
	InvalidParams ErrorCode = 1001 // 无效参数

	// 客户端错误 - 认证授权相关 (1xxx)
	Unauthorized ErrorCode = 1002 // 未授权
	Forbidden    ErrorCode = 1003 // 禁止访问

	// 客户端错误 - 资源相关 (1xxx)
	NotFound      ErrorCode = 1004 // 资源不存在 (HTTP 404)
	AlreadyExists ErrorCode = 1005 // 资源已存在
	Conflict      ErrorCode = 1006 // 冲突

	// 认证授权错误 (1xxx)
	InvalidCredentials ErrorCode = 2002 // 无效凭证
	TokenExpired       ErrorCode = 2007 // Token过期
	TokenInvalid       ErrorCode = 2008 // Token无效
	PasswordTooWeak    ErrorCode = 2009 // 密码强度不足
	AccountLocked      ErrorCode = 2010 // 账户已锁定
	AccountDisabled    ErrorCode = 2011 // 账户已禁用

	// 业务逻辑错误 (1xxx)
	InsufficientBalance ErrorCode = 3003 // 余额不足
	InsufficientQuota   ErrorCode = 3010 // 配额不足
	WalletFrozen        ErrorCode = 3011 // 钱包已冻结
	ContentNotPublished ErrorCode = 3012 // 内容未发布
	ChapterLocked       ErrorCode = 3013 // 章节已锁定

	// 内容审核错误 (1xxx)
	ContentPendingReview ErrorCode = 3014 // 内容待审核
	ContentRejected      ErrorCode = 3015 // 内容被拒绝
	ContentViolation     ErrorCode = 3016 // 内容违规

	// 服务器错误 - 系统错误 (5xxx)
	InternalError ErrorCode = 5000 // 内部错误
	DatabaseError ErrorCode = 5001 // 数据库错误
	RedisError    ErrorCode = 5003 // Redis错误

	// 服务器错误 - 外部服务 (5xxx)
	ExternalAPIError  ErrorCode = 5004 // 外部API错误
	RateLimitExceeded ErrorCode = 4290 // 超出频率限制
)

// ErrorInfo 错误信息
type ErrorInfo struct {
	Code       ErrorCode   `json:"code"`
	Message    string      `json:"message"`
	HTTPStatus int         `json:"-"`
	Details    interface{} `json:"details,omitempty"`
}

// DefaultErrorMessages 默认错误消息
var DefaultErrorMessages = map[ErrorCode]string{
	InvalidParams:     "请求参数无效",
	Unauthorized:      "未授权访问",
	Forbidden:         "禁止访问",
	NotFound:          "资源不存在",
	AlreadyExists:     "资源已存在",
	Conflict:          "请求冲突",
	RateLimitExceeded: "请求过于频繁",

	InvalidCredentials: "用户名或密码错误",
	TokenExpired:       "登录已过期，请重新登录",
	TokenInvalid:       "无效的登录凭证",
	PasswordTooWeak:    "密码强度不足",
	AccountLocked:      "账户已被锁定",
	AccountDisabled:    "账户已被禁用",

	InsufficientBalance: "余额不足",
	InsufficientQuota:   "配额不足",
	WalletFrozen:        "钱包已被冻结",
	ContentNotPublished: "内容尚未发布",
	ChapterLocked:       "章节已锁定",

	ContentPendingReview: "内容待审核",
	ContentRejected:      "内容审核未通过",
	ContentViolation:     "内容违规",

	InternalError:    "服务器内部错误",
	DatabaseError:    "数据库错误",
	RedisError:       "缓存服务错误",
	ExternalAPIError: "外部服务错误",
}

// DefaultHTTPStatus 默认HTTP状态码
var DefaultHTTPStatus = map[ErrorCode]int{
	InvalidParams:     http.StatusBadRequest,
	Unauthorized:      http.StatusUnauthorized,
	Forbidden:         http.StatusForbidden,
	NotFound:          http.StatusNotFound,
	AlreadyExists:     http.StatusConflict,
	Conflict:          http.StatusConflict,
	RateLimitExceeded: http.StatusTooManyRequests,

	InvalidCredentials: http.StatusUnauthorized,
	TokenExpired:       http.StatusUnauthorized,
	TokenInvalid:       http.StatusUnauthorized,
	PasswordTooWeak:    http.StatusBadRequest,
	AccountLocked:      http.StatusForbidden,
	AccountDisabled:    http.StatusForbidden,

	InsufficientBalance: http.StatusBadRequest,
	InsufficientQuota:   http.StatusBadRequest,
	WalletFrozen:        http.StatusForbidden,
	ContentNotPublished: http.StatusForbidden,
	ChapterLocked:       http.StatusForbidden,

	ContentPendingReview: http.StatusAccepted,
	ContentRejected:      http.StatusForbidden,
	ContentViolation:     http.StatusForbidden,

	InternalError:    http.StatusInternalServerError,
	DatabaseError:    http.StatusInternalServerError,
	RedisError:       http.StatusInternalServerError,
	ExternalAPIError: http.StatusBadGateway,
}

// GetHTTPStatus 获取错误码对应的HTTP状态码
func GetHTTPStatus(code ErrorCode) int {
	if status, ok := DefaultHTTPStatus[code]; ok {
		return status
	}
	return http.StatusInternalServerError
}

// GetDefaultMessage 获取默认错误消息
func GetDefaultMessage(code ErrorCode) string {
	if msg, ok := DefaultErrorMessages[code]; ok {
		return msg
	}
	return "未知错误"
}
