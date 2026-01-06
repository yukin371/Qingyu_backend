package errors

import "net/http"

// ErrorCode 错误码类型
type ErrorCode int

// 常用错误码定义
const (
	// 成功
	Success ErrorCode = 0

	// 客户端错误 1000-1999
	InvalidParams   ErrorCode = 1001 // 无效参数
	Unauthorized    ErrorCode = 1002 // 未授权
	Forbidden       ErrorCode = 1003 // 禁止访问
	NotFound        ErrorCode = 1004 // 资源不存在
	AlreadyExists   ErrorCode = 1005 // 资源已存在
	Conflict        ErrorCode = 1006 // 冲突
	RateLimitExceeded ErrorCode = 1007 // 超出频率限制

	// 认证授权错误 1100-1199
	InvalidCredentials ErrorCode = 1101 // 无效凭证
	TokenExpired       ErrorCode = 1102 // Token过期
	TokenInvalid       ErrorCode = 1103 // Token无效
	PasswordTooWeak    ErrorCode = 1104 // 密码强度不足
	AccountLocked      ErrorCode = 1105 // 账户已锁定
	AccountDisabled    ErrorCode = 1106 // 账户已禁用

	// 业务逻辑错误 1200-1299
	InsufficientBalance ErrorCode = 1201 // 余额不足
	InsufficientQuota   ErrorCode = 1202 // 配额不足
	WalletFrozen        ErrorCode = 1203 // 钱包已冻结
	ContentNotPublished ErrorCode = 1204 // 内容未发布
	ChapterLocked       ErrorCode = 1205 // 章节已锁定

	// 内容审核错误 1300-1399
	ContentPendingReview ErrorCode = 1301 // 内容待审核
	ContentRejected       ErrorCode = 1302 // 内容被拒绝
	ContentViolation      ErrorCode = 1303 // 内容违规

	// 服务器错误 5000-5999
	InternalError   ErrorCode = 5000 // 内部错误
	DatabaseError   ErrorCode = 5001 // 数据库错误
	RedisError      ErrorCode = 5002 // Redis错误
	ExternalAPIError ErrorCode = 5003 // 外部API错误
)

// ErrorInfo 错误信息
type ErrorInfo struct {
	Code       ErrorCode `json:"code"`
	Message    string    `json:"message"`
	HTTPStatus int       `json:"-"`
	Details    interface{} `json:"details,omitempty"`
}

// DefaultErrorMessages 默认错误消息
var DefaultErrorMessages = map[ErrorCode]string{
	InvalidParams:    "请求参数无效",
	Unauthorized:     "未授权访问",
	Forbidden:        "禁止访问",
	NotFound:         "资源不存在",
	AlreadyExists:    "资源已存在",
	Conflict:         "请求冲突",
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

	InternalError:   "服务器内部错误",
	DatabaseError:   "数据库错误",
	RedisError:      "缓存服务错误",
	ExternalAPIError: "外部服务错误",
}

// DefaultHTTPStatus 默认HTTP状态码
var DefaultHTTPStatus = map[ErrorCode]int{
	InvalidParams:    http.StatusBadRequest,
	Unauthorized:     http.StatusUnauthorized,
	Forbidden:        http.StatusForbidden,
	NotFound:         http.StatusNotFound,
	AlreadyExists:    http.StatusConflict,
	Conflict:         http.StatusConflict,
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

	InternalError:   http.StatusInternalServerError,
	DatabaseError:   http.StatusInternalServerError,
	RedisError:      http.StatusInternalServerError,
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
