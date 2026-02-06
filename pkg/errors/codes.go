package errors

import "net/http"

// ErrorCode 错误码类型
type ErrorCode int

// 常用错误码定义
// 格式：4位数字，前1位表示错误类别，后3位为具体错误编号
// 0: 成功
// 1xxx: 通用客户端错误（参数、验证等）
// 2xxx: 用户相关错误（认证、授权等）
// 3xxx: 业务逻辑错误
// 4xxx: 限流和配额错误
// 5xxx: 服务器内部错误
const (
	// 成功
	Success ErrorCode = 0

	// ========== 通用客户端错误 (1000-1099) ==========

	// 参数验证错误 (1000-1019)
	InvalidParams    ErrorCode = 1001 // 无效参数
	MissingParam     ErrorCode = 1008 // 缺少必填参数
	InvalidFormat    ErrorCode = 1009 // 参数格式无效
	InvalidLength    ErrorCode = 1010 // 参数长度无效
	InvalidType      ErrorCode = 1011 // 参数类型无效
	OutOfRange       ErrorCode = 1012 // 参数超出范围
	DuplicateField   ErrorCode = 1013 // 字段重复
	UnknownField     ErrorCode = 1014 // 未知字段
	ValidationFailed ErrorCode = 1015 // 验证失败
	InvalidValue     ErrorCode = 1016 // 值无效
	InvalidOperation ErrorCode = 1017 // 无效操作

	// 认证授权错误 (1020-1039)
	Unauthorized ErrorCode = 1002 // 未授权
	Forbidden    ErrorCode = 1003 // 禁止访问

	// 资源相关错误 (1040-1059)
	NotFound      ErrorCode = 1004 // 资源不存在 (HTTP 404)
	AlreadyExists ErrorCode = 1005 // 资源已存在
	Conflict      ErrorCode = 1006 // 冲突
	ResourceGone  ErrorCode = 1018 // 资源已删除

	// 请求相关错误 (1060-1079)
	MethodNotAllowed ErrorCode = 1019 // 方法不允许
	RequestTimeout   ErrorCode = 1020 // 请求超时

	// ========== 用户相关错误 (2000-2199) ==========

	// 用户认证 (2000-2019)
	UserNotFound         ErrorCode = 2001 // 用户不存在
	InvalidCredentials   ErrorCode = 2002 // 无效凭证（用户名或密码错误）
	UsernameAlreadyUsed  ErrorCode = 2003 // 用户名已被使用
	EmailAlreadyUsed     ErrorCode = 2004 // 邮箱已被使用
	EmailSendFailed      ErrorCode = 2005 // 邮件发送失败
	InvalidCode          ErrorCode = 2006 // 验证码无效
	CodeExpired          ErrorCode = 2007 // 验证码过期
	TokenExpired         ErrorCode = 2008 // Token过期
	TokenInvalid         ErrorCode = 2009 // Token无效
	TokenFormatError     ErrorCode = 2010 // Token格式错误
	TokenMissing         ErrorCode = 2011 // Token缺失
	RefreshTokenExpired  ErrorCode = 2012 // Refresh Token过期
	RefreshTokenInvalid  ErrorCode = 2013 // Refresh Token无效
	PasswordTooWeak      ErrorCode = 2014 // 密码强度不足
	AccountLocked        ErrorCode = 2015 // 账户已锁定
	AccountDisabled      ErrorCode = 2016 // 账户已禁用
	TokenRevoked         ErrorCode = 2017 // Token已被撤销
	SessionExpired       ErrorCode = 2018 // 会话过期
	TooManyAttempts      ErrorCode = 2019 // 尝试次数过多
	AccountNotVerified   ErrorCode = 2020 // 账户未验证

	// 邮箱和手机 (2020-2039)
	PhoneAlreadyUsed     ErrorCode = 2021 // 手机号已使用
	InvalidPhoneFormat   ErrorCode = 2022 // 手机号格式无效
	EmailNotVerified     ErrorCode = 2023 // 邮箱未验证
	PhoneNotVerified     ErrorCode = 2024 // 手机号未验证
	SmsSendFailed        ErrorCode = 2025 // 短信发送失败
	EmailSendTooFrequent ErrorCode = 2026 // 邮件发送过于频繁
	SmsSendTooFrequent   ErrorCode = 2027 // 短信发送过于频繁

	// 评分相关 (2500-2599)
	RatingNotFound       ErrorCode = 2501 // 评分不存在
	RatingInvalid        ErrorCode = 2502 // 评分值无效（不在1-5范围）
	RatingAlreadyExists  ErrorCode = 2503 // 用户已评分
	RatingUnauthorized   ErrorCode = 2504 // 无权操作此评分
	RatingTargetNotFound ErrorCode = 2505 // 评分目标不存在

	// ========== 业务逻辑错误 (3000-3199) ==========

	// 书籍相关 (3000-3039)
	BookNotFound      ErrorCode = 3001 // 书籍不存在
	BookAlreadyExists ErrorCode = 3004 // 书籍已存在
	InvalidBookStatus ErrorCode = 3005 // 书籍状态无效
	BookDeleted       ErrorCode = 3006 // 书籍已删除

	// 章节相关 (3040-3069)
	ChapterNotFound      ErrorCode = 3002 // 章节不存在
	ChapterAlreadyExists ErrorCode = 3007 // 章节已存在
	InvalidChapterStatus ErrorCode = 3008 // 章节状态无效
	ChapterDeleted       ErrorCode = 3009 // 章节已删除

	// 财务相关 (3070-3099)
	InsufficientBalance ErrorCode = 3003 // 余额不足
	InsufficientQuota   ErrorCode = 3010 // 配额不足
	WalletFrozen        ErrorCode = 3011 // 钱包已冻结
	TransactionFailed   ErrorCode = 3017 // 交易失败

	// 内容相关 (3100-3129)
	ContentNotPublished ErrorCode = 3012 // 内容未发布
	ChapterLocked       ErrorCode = 3013 // 章节已锁定
	ContentLocked       ErrorCode = 3018 // 内容已锁定
	ContentDeleted      ErrorCode = 3019 // 内容已删除

	// 内容审核 (3130-3159)
	ContentPendingReview ErrorCode = 3014 // 内容待审核
	ContentRejected      ErrorCode = 3015 // 内容被拒绝
	ContentViolation     ErrorCode = 3016 // 内容违规

	// 角色相关 (3160-3179)
	CharacterNotFound    ErrorCode = 3020 // 角色不存在
	InvalidCharacterData ErrorCode = 3021 // 角色数据无效

	// 评论相关 (3180-3199)
	ReviewNotFound ErrorCode = 3022 // 评论不存在

	// 收藏和关注 (3200-3219)
	CollectionNotFound ErrorCode = 3023 // 收藏不存在
	AlreadyCollected   ErrorCode = 3024 // 已收藏
	AlreadyFollowed    ErrorCode = 3025 // 已关注

	// ========== 限流和配额错误 (4000-4099) ==========

	RateLimitExceeded     ErrorCode = 4000 // 频率限制超出
	DailyLimitExceeded    ErrorCode = 4001 // 每日限制超出
	HourlyLimitExceeded   ErrorCode = 4002 // 每小时限制超出
	MinuteLimitExceeded   ErrorCode = 4003 // 每分钟限制超出
	UploadLimitExceeded   ErrorCode = 4004 // 上传限制超出
	StorageLimitExceeded  ErrorCode = 4005 // 存储限制超出
	ApiQuotaExceeded      ErrorCode = 4006 // API配额超出
	ConcurrentLimitExceeded ErrorCode = 4007 // 并发限制超出
	RateLimitLogin        ErrorCode = 4008 // 登录频率限制
	RateLimitEmailSend    ErrorCode = 4009 // 邮件发送频率限制
	RateLimitSmsSend      ErrorCode = 4010 // 短信发送频率限制
	HourlyLimitExceededOld ErrorCode = 4291 // 小时级限制超出（旧版，保留兼容）

	// ========== 服务器内部错误 (5000-5099) ==========

	// 系统错误 (5000-5019)
	InternalError ErrorCode = 5000 // 内部错误
	DatabaseError ErrorCode = 5001 // 数据库错误
	RedisError    ErrorCode = 5003 // Redis错误

	// 外部服务 (5020-5039)
	ExternalAPIError  ErrorCode = 5004 // 外部API错误
	ServiceUnavailable ErrorCode = 5002 // 服务不可用
	CacheError        ErrorCode = 5005 // 缓存错误
	QueueError        ErrorCode = 5006 // 队列错误
	StorageError      ErrorCode = 5007 // 存储错误
	NetworkError      ErrorCode = 5008 // 网络错误
	ConfigurationError ErrorCode = 5009 // 配置错误
	DependencyError   ErrorCode = 5010 // 依赖错误
	TimeoutError      ErrorCode = 5011 // 超时错误
	OverloadedError   ErrorCode = 5012 // 过载错误
	MaintenanceError  ErrorCode = 5013 // 维护中

	// 数据库详细错误 (5040-5059)
	DatabaseConnectionFailed ErrorCode = 5014 // 数据库连接失败
	DatabaseQueryTimeout     ErrorCode = 5015 // 数据库查询超时
	DatabaseTransactionFailed ErrorCode = 5016 // 数据库事务失败
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
	// 成功
	Success: "成功",

	// 通用客户端错误 (1000-1099)
	InvalidParams:     "请求参数无效",
	MissingParam:      "缺少必填参数",
	InvalidFormat:     "参数格式无效",
	InvalidLength:     "参数长度无效",
	InvalidType:       "参数类型无效",
	OutOfRange:        "参数超出范围",
	DuplicateField:    "字段重复",
	UnknownField:      "未知字段",
	ValidationFailed:  "验证失败",
	InvalidValue:      "值无效",
	InvalidOperation:  "无效操作",
	Unauthorized:      "未授权访问",
	Forbidden:         "禁止访问",
	NotFound:          "资源不存在",
	AlreadyExists:     "资源已存在",
	Conflict:          "请求冲突",
	ResourceGone:      "资源已删除",
	MethodNotAllowed:  "方法不允许",
	RequestTimeout:    "请求超时",

	// 用户相关错误 (2000-2199)
	UserNotFound:         "用户不存在",
	InvalidCredentials:   "用户名或密码错误",
	UsernameAlreadyUsed:  "用户名已被使用",
	EmailAlreadyUsed:     "邮箱已被使用",
	EmailSendFailed:      "邮件发送失败",
	InvalidCode:          "验证码无效",
	CodeExpired:          "验证码过期",
	TokenExpired:         "登录已过期，请重新登录",
	TokenInvalid:         "无效的登录凭证",
	TokenFormatError:     "Token格式错误",
	TokenMissing:         "用户未认证",
	RefreshTokenExpired:  "Refresh Token已过期",
	RefreshTokenInvalid:  "无效的Refresh Token",
	PasswordTooWeak:      "密码强度不足",
	AccountLocked:        "账户已被锁定",
	AccountDisabled:      "账户已被禁用",
	TokenRevoked:         "Token已被撤销",
	SessionExpired:       "会话已过期",
	TooManyAttempts:      "尝试次数过多",
	AccountNotVerified:   "账户未验证",
	PhoneAlreadyUsed:     "手机号已被使用",
	InvalidPhoneFormat:   "手机号格式无效",
	EmailNotVerified:     "邮箱未验证",
	PhoneNotVerified:     "手机号未验证",
	SmsSendFailed:        "短信发送失败",
	EmailSendTooFrequent: "邮件发送过于频繁",
	SmsSendTooFrequent:   "短信发送过于频繁",

	// 评分相关 (2500-2599)
	RatingNotFound:       "评分不存在",
	RatingInvalid:        "评分值无效（必须在1-5之间）",
	RatingAlreadyExists:  "您已经评分过了",
	RatingUnauthorized:   "无权操作此评分",
	RatingTargetNotFound: "评分目标不存在",

	// 业务逻辑错误 (3000-3199)
	BookNotFound:         "书籍不存在",
	BookAlreadyExists:    "书籍已存在",
	InvalidBookStatus:    "书籍状态无效",
	BookDeleted:          "书籍已删除",
	ChapterNotFound:      "章节不存在",
	ChapterAlreadyExists: "章节已存在",
	InvalidChapterStatus: "章节状态无效",
	ChapterDeleted:       "章节已删除",
	InsufficientBalance:  "余额不足",
	InsufficientQuota:    "配额不足",
	WalletFrozen:         "钱包已被冻结",
	TransactionFailed:    "交易失败",
	ContentNotPublished:  "内容尚未发布",
	ChapterLocked:        "章节已锁定",
	ContentLocked:        "内容已锁定",
	ContentDeleted:       "内容已删除",
	ContentPendingReview: "内容待审核",
	ContentRejected:      "内容审核未通过",
	ContentViolation:     "内容违规",
	CharacterNotFound:    "角色不存在",
	InvalidCharacterData: "角色数据无效",
	ReviewNotFound:       "评论不存在",
	CollectionNotFound:   "收藏不存在",
	AlreadyCollected:     "已收藏",
	AlreadyFollowed:      "已关注",

	// 限流和配额错误 (4000-4099)
	RateLimitExceeded:      "请求过于频繁",
	DailyLimitExceeded:     "每日请求次数已达上限",
	HourlyLimitExceeded:    "每小时请求次数已达上限",
	MinuteLimitExceeded:    "每分钟请求次数已达上限",
	UploadLimitExceeded:    "上传次数已达上限",
	StorageLimitExceeded:   "存储空间已达上限",
	ApiQuotaExceeded:       "API调用配额已达上限",
	ConcurrentLimitExceeded: "并发请求数已达上限",
	RateLimitLogin:         "登录过于频繁，请稍后再试",
	RateLimitEmailSend:     "邮件发送过于频繁",
	RateLimitSmsSend:       "短信发送过于频繁",
	HourlyLimitExceededOld: "小时级限制超出",

	// 服务器内部错误 (5000-5099)
	InternalError:           "服务器内部错误",
	DatabaseError:           "数据库错误",
	RedisError:              "缓存服务错误",
	ExternalAPIError:        "外部服务错误",
	ServiceUnavailable:      "服务暂时不可用",
	CacheError:              "缓存错误",
	QueueError:              "队列错误",
	StorageError:            "存储错误",
	NetworkError:            "网络错误",
	ConfigurationError:      "配置错误",
	DependencyError:         "依赖服务错误",
	TimeoutError:            "操作超时",
	OverloadedError:         "服务过载",
	MaintenanceError:        "系统维护中",
	DatabaseConnectionFailed: "数据库连接失败",
	DatabaseQueryTimeout:    "数据库查询超时",
	DatabaseTransactionFailed: "数据库事务失败",
}

// DefaultHTTPStatus 默认HTTP状态码
var DefaultHTTPStatus = map[ErrorCode]int{
	// 成功
	Success: http.StatusOK,

	// 通用客户端错误 (1000-1099)
	InvalidParams:     http.StatusBadRequest,
	MissingParam:      http.StatusBadRequest,
	InvalidFormat:     http.StatusBadRequest,
	InvalidLength:     http.StatusBadRequest,
	InvalidType:       http.StatusBadRequest,
	OutOfRange:        http.StatusBadRequest,
	DuplicateField:    http.StatusConflict,
	UnknownField:      http.StatusBadRequest,
	ValidationFailed:  http.StatusBadRequest,
	InvalidValue:      http.StatusBadRequest,
	InvalidOperation:  http.StatusBadRequest,
	Unauthorized:      http.StatusUnauthorized,
	Forbidden:         http.StatusForbidden,
	NotFound:          http.StatusNotFound,
	AlreadyExists:     http.StatusConflict,
	Conflict:          http.StatusConflict,
	ResourceGone:      http.StatusGone,
	MethodNotAllowed:  http.StatusMethodNotAllowed,
	RequestTimeout:    http.StatusRequestTimeout,

	// 用户相关错误 (2000-2199)
	UserNotFound:         http.StatusNotFound,
	InvalidCredentials:   http.StatusUnauthorized,
	UsernameAlreadyUsed:  http.StatusConflict,
	EmailAlreadyUsed:     http.StatusConflict,
	EmailSendFailed:      http.StatusInternalServerError,
	InvalidCode:          http.StatusBadRequest,
	CodeExpired:          http.StatusBadRequest,
	TokenExpired:         http.StatusUnauthorized,
	TokenInvalid:         http.StatusUnauthorized,
	TokenFormatError:     http.StatusBadRequest,
	TokenMissing:         http.StatusUnauthorized,
	RefreshTokenExpired:  http.StatusUnauthorized,
	RefreshTokenInvalid:  http.StatusUnauthorized,
	PasswordTooWeak:      http.StatusBadRequest,
	AccountLocked:        http.StatusForbidden,
	AccountDisabled:      http.StatusForbidden,
	TokenRevoked:         http.StatusUnauthorized,
	SessionExpired:       http.StatusUnauthorized,
	TooManyAttempts:      http.StatusTooManyRequests,
	AccountNotVerified:   http.StatusForbidden,
	PhoneAlreadyUsed:     http.StatusConflict,
	InvalidPhoneFormat:   http.StatusBadRequest,
	EmailNotVerified:     http.StatusForbidden,
	PhoneNotVerified:     http.StatusForbidden,
	SmsSendFailed:        http.StatusInternalServerError,
	EmailSendTooFrequent: http.StatusTooManyRequests,
	SmsSendTooFrequent:   http.StatusTooManyRequests,

	// 评分相关 (2500-2599)
	RatingNotFound:       http.StatusNotFound,
	RatingInvalid:        http.StatusBadRequest,
	RatingAlreadyExists:  http.StatusConflict,
	RatingUnauthorized:   http.StatusForbidden,
	RatingTargetNotFound: http.StatusNotFound,

	// 业务逻辑错误 (3000-3199)
	BookNotFound:         http.StatusNotFound,
	BookAlreadyExists:    http.StatusConflict,
	InvalidBookStatus:    http.StatusBadRequest,
	BookDeleted:          http.StatusGone,
	ChapterNotFound:      http.StatusNotFound,
	ChapterAlreadyExists: http.StatusConflict,
	InvalidChapterStatus: http.StatusBadRequest,
	ChapterDeleted:       http.StatusGone,
	InsufficientBalance:  http.StatusBadRequest,
	InsufficientQuota:    http.StatusBadRequest,
	WalletFrozen:         http.StatusForbidden,
	TransactionFailed:    http.StatusInternalServerError,
	ContentNotPublished:  http.StatusForbidden,
	ChapterLocked:        http.StatusForbidden,
	ContentLocked:        http.StatusForbidden,
	ContentDeleted:       http.StatusGone,
	ContentPendingReview: http.StatusAccepted,
	ContentRejected:      http.StatusForbidden,
	ContentViolation:     http.StatusForbidden,
	CharacterNotFound:    http.StatusNotFound,
	InvalidCharacterData: http.StatusBadRequest,
	ReviewNotFound:       http.StatusNotFound,
	CollectionNotFound:   http.StatusNotFound,
	AlreadyCollected:     http.StatusConflict,
	AlreadyFollowed:      http.StatusConflict,

	// 限流和配额错误 (4000-4099)
	RateLimitExceeded:       http.StatusTooManyRequests,
	DailyLimitExceeded:      http.StatusTooManyRequests,
	HourlyLimitExceeded:     http.StatusTooManyRequests,
	MinuteLimitExceeded:     http.StatusTooManyRequests,
	UploadLimitExceeded:     http.StatusTooManyRequests,
	StorageLimitExceeded:    http.StatusInsufficientStorage,
	ApiQuotaExceeded:        http.StatusTooManyRequests,
	ConcurrentLimitExceeded: http.StatusTooManyRequests,
	RateLimitLogin:          http.StatusTooManyRequests,
	RateLimitEmailSend:      http.StatusTooManyRequests,
	RateLimitSmsSend:        http.StatusTooManyRequests,
	HourlyLimitExceededOld:  http.StatusTooManyRequests,

	// 服务器内部错误 (5000-5099)
	InternalError:             http.StatusInternalServerError,
	DatabaseError:             http.StatusInternalServerError,
	RedisError:                http.StatusInternalServerError,
	ExternalAPIError:          http.StatusBadGateway,
	ServiceUnavailable:        http.StatusServiceUnavailable,
	CacheError:                http.StatusInternalServerError,
	QueueError:                http.StatusInternalServerError,
	StorageError:              http.StatusInternalServerError,
	NetworkError:              http.StatusInternalServerError,
	ConfigurationError:        http.StatusInternalServerError,
	DependencyError:           http.StatusBadGateway,
	TimeoutError:              http.StatusGatewayTimeout,
	OverloadedError:           http.StatusServiceUnavailable,
	MaintenanceError:          http.StatusServiceUnavailable,
	DatabaseConnectionFailed:  http.StatusInternalServerError,
	DatabaseQueryTimeout:      http.StatusInternalServerError,
	DatabaseTransactionFailed: http.StatusInternalServerError,
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
