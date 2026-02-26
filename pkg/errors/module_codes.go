package errors

// ============================================================================
// 模块专属错误码
// 
// 本文件定义了各模块的专属错误码，补充 codes.go 中的通用错误码
// 错误码分配规则：
// - Writer模块: 3300-3399
// - Reader模块: 3400-3499
// - AIService模块: 3500-3599
// - Social模块: 3600-3699
// - Messaging模块: 3700-3799
// - Admin模块: 3800-3899
// ============================================================================

const (
	// ========== Writer模块错误 (3300-3399) ==========

	// Writer项目相关 (3300-3319)
	WriterProjectNotFound     ErrorCode = 3301 // 项目不存在
	WriterDocumentNotFound    ErrorCode = 3302 // 文档不存在
	WriterCommentNotFound     ErrorCode = 3303 // 批注不存在
	WriterCharacterNotFound   ErrorCode = 3304 // 角色不存在
	WriterLocationNotFound    ErrorCode = 3305 // 地点不存在
	WriterTimelineNotFound    ErrorCode = 3306 // 时间线不存在
	WriterVersionNotFound     ErrorCode = 3307 // 版本不存在
	WriterPublicationNotFound ErrorCode = 3308 // 发布记录不存在
	WriterExportTaskNotFound  ErrorCode = 3309 // 导出任务不存在

	// Writer验证错误 (3320-3339)
	WriterInvalidProjectID    ErrorCode = 3320 // 无效的项目ID
	WriterInvalidDocumentID   ErrorCode = 3321 // 无效的文档ID
	WriterInvalidContent      ErrorCode = 3322 // 内容无效
	WriterInvalidVersion      ErrorCode = 3323 // 无效的版本号
	WriterInvalidExportFormat ErrorCode = 3324 // 无效的导出格式
	WriterInvalidRelationType ErrorCode = 3325 // 无效的关系类型

	// Writer冲突错误 (3340-3359)
	WriterProjectAlreadyExists  ErrorCode = 3340 // 项目已存在
	WriterDocumentAlreadyExists ErrorCode = 3341 // 文档已存在
	WriterNameAlreadyExists     ErrorCode = 3342 // 名称已存在
	WriterVersionConflict       ErrorCode = 3343 // 版本冲突
	WriterEditConflict          ErrorCode = 3344 // 编辑冲突

	// Writer授权错误 (3360-3379)
	WriterUnauthorized ErrorCode = 3360 // 未授权
	WriterForbidden    ErrorCode = 3361 // 禁止访问
	WriterNoPermission ErrorCode = 3362 // 无权限

	// Writer系统错误 (3380-3399)
	WriterPublishFailed      ErrorCode = 3380 // 发布失败
	WriterExportFailed       ErrorCode = 3381 // 导出失败
	WriterStorageError       ErrorCode = 3382 // 存储错误
	WriterExternalServiceErr ErrorCode = 3383 // 外部服务错误

	// ========== Reader模块错误 (3400-3499) ==========

	// Reader资源相关 (3400-3419)
	ReaderProgressNotFound   ErrorCode = 3401 // 阅读进度不存在
	ReaderAnnotationNotFound ErrorCode = 3402 // 标注不存在
	ReaderSettingsNotFound   ErrorCode = 3403 // 阅读设置不存在
	ReaderChapterNotFound    ErrorCode = 3404 // 章节不存在
	ReaderBookNotFound       ErrorCode = 3405 // 书籍不存在

	// Reader验证错误 (3420-3439)
	ReaderInvalidProgress      ErrorCode = 3420 // 无效的阅读进度
	ReaderInvalidAnnotation    ErrorCode = 3421 // 无效的标注
	ReaderInvalidSettings      ErrorCode = 3422 // 无效的阅读设置
	ReaderInvalidStatus        ErrorCode = 3423 // 无效的书籍状态
	ReaderInvalidChapterNumber ErrorCode = 3424 // 无效的章节号

	// Reader授权错误 (3440-3459)
	ReaderAccessDenied ErrorCode = 3440 // 访问被拒绝
	ReaderForbidden    ErrorCode = 3441 // 禁止访问

	// Reader系统错误 (3460-3479)
	ReaderSyncFailed ErrorCode = 3460 // 同步失败

	// ========== AIService模块错误 (3500-3599) ==========

	// AI服务可用性 (3500-3509)
	AIServiceUnavailable ErrorCode = 3501 // AI服务不可用
	AIServiceTimeout     ErrorCode = 3502 // AI服务超时
	AIServiceRateLimit   ErrorCode = 3503 // AI服务频率限制

	// AI配额相关 (3510-3519)
	AIQuotaExhausted ErrorCode = 3510 // AI配额不足
	AIQuotaCheckFail ErrorCode = 3511 // AI配额检查失败

	// AI请求相关 (3520-3529)
	AIInvalidRequest  ErrorCode = 3520 // AI请求参数无效
	AIRequestTooLarge ErrorCode = 3521 // AI请求过大
	AIContextExceeded ErrorCode = 3522 // AI上下文超限

	// AI模型相关 (3530-3539)
	AIModelNotFound     ErrorCode = 3530 // AI模型不存在
	AIModelNotAvailable ErrorCode = 3531 // AI模型不可用

	// ========== Social模块错误 (3600-3699) ==========

	// Social资源相关 (3600-3619)
	SocialCommentNotFound ErrorCode = 3601 // 评论不存在
	SocialLikeNotFound    ErrorCode = 3602 // 点赞不存在
	SocialFollowNotFound  ErrorCode = 3603 // 关注不存在

	// Social验证错误 (3620-3639)
	SocialInvalidComment   ErrorCode = 3620 // 无效的评论
	SocialInvalidRating    ErrorCode = 3621 // 无效的评分
	SocialInvalidRelation  ErrorCode = 3622 // 无效的关系

	// Social冲突错误 (3640-3659)
	SocialAlreadyLiked  ErrorCode = 3640 // 已点赞
	SocialAlreadyFollowed ErrorCode = 3641 // 已关注

	// ========== Messaging模块错误 (3700-3799) ==========

	// Messaging资源相关 (3700-3719)
	MessagingConversationNotFound ErrorCode = 3701 // 会话不存在
	MessagingMessageNotFound      ErrorCode = 3702 // 消息不存在

	// Messaging验证错误 (3720-3739)
	MessagingInvalidRecipient ErrorCode = 3720 // 无效的接收者
	MessagingInvalidContent   ErrorCode = 3721 // 无效的消息内容
	MessagingEmptyMessage     ErrorCode = 3722 // 消息为空

	// Messaging系统错误 (3740-3759)
	MessagingSendFailed ErrorCode = 3740 // 消息发送失败

	// ========== Admin模块错误 (3800-3899) ==========

	// Admin资源相关 (3800-3819)
	AdminAuditLogNotFound ErrorCode = 3801 // 审计日志不存在
	AdminReportNotFound   ErrorCode = 3802 // 举报记录不存在

	// Admin验证错误 (3820-3839)
	AdminInvalidAction    ErrorCode = 3820 // 无效的管理操作
	AdminInvalidStatus    ErrorCode = 3821 // 无效的状态
	AdminInvalidPermission ErrorCode = 3822 // 无效的权限设置

	// Admin授权错误 (3840-3859)
	AdminUnauthorized     ErrorCode = 3840 // 未授权的管理操作
	AdminInsufficientRole ErrorCode = 3841 // 角色权限不足
)

// ============================================================================
// 模块专属错误消息
// ============================================================================

// ModuleErrorMessages 模块专属错误消息映射
var ModuleErrorMessages = map[ErrorCode]string{
	// Writer模块
	WriterProjectNotFound:     "项目不存在",
	WriterDocumentNotFound:    "文档不存在",
	WriterCommentNotFound:     "批注不存在",
	WriterCharacterNotFound:   "角色不存在",
	WriterLocationNotFound:    "地点不存在",
	WriterTimelineNotFound:    "时间线不存在",
	WriterVersionNotFound:     "版本不存在",
	WriterPublicationNotFound: "发布记录不存在",
	WriterExportTaskNotFound:  "导出任务不存在",
	WriterInvalidProjectID:    "无效的项目ID",
	WriterInvalidDocumentID:   "无效的文档ID",
	WriterInvalidContent:      "内容无效",
	WriterInvalidVersion:      "无效的版本号",
	WriterInvalidExportFormat: "无效的导出格式",
	WriterInvalidRelationType: "无效的关系类型",
	WriterProjectAlreadyExists: "项目已存在",
	WriterDocumentAlreadyExists: "文档已存在",
	WriterNameAlreadyExists:    "名称已存在",
	WriterVersionConflict:      "版本冲突",
	WriterEditConflict:         "编辑冲突",
	WriterUnauthorized:         "未授权",
	WriterForbidden:            "禁止访问",
	WriterNoPermission:         "无权限",
	WriterPublishFailed:        "发布失败",
	WriterExportFailed:         "导出失败",
	WriterStorageError:         "存储错误",
	WriterExternalServiceErr:   "外部服务错误",

	// Reader模块
	ReaderProgressNotFound:     "阅读进度不存在",
	ReaderAnnotationNotFound:   "标注不存在",
	ReaderSettingsNotFound:     "阅读设置不存在",
	ReaderChapterNotFound:      "章节不存在",
	ReaderBookNotFound:         "书籍不存在",
	ReaderInvalidProgress:      "无效的阅读进度",
	ReaderInvalidAnnotation:    "无效的标注",
	ReaderInvalidSettings:      "无效的阅读设置",
	ReaderInvalidStatus:        "无效的书籍状态",
	ReaderInvalidChapterNumber: "无效的章节号",
	ReaderAccessDenied:         "访问被拒绝",
	ReaderForbidden:            "禁止访问",
	ReaderSyncFailed:           "同步失败",

	// AIService模块
	AIServiceUnavailable:  "AI服务不可用",
	AIServiceTimeout:      "AI服务超时",
	AIServiceRateLimit:    "AI服务频率限制",
	AIQuotaExhausted:      "AI配额不足",
	AIQuotaCheckFail:      "AI配额检查失败",
	AIInvalidRequest:      "AI请求参数无效",
	AIRequestTooLarge:     "AI请求过大",
	AIContextExceeded:     "AI上下文超限",
	AIModelNotFound:       "AI模型不存在",
	AIModelNotAvailable:   "AI模型不可用",

	// Social模块
	SocialCommentNotFound:  "评论不存在",
	SocialLikeNotFound:     "点赞不存在",
	SocialFollowNotFound:   "关注不存在",
	SocialInvalidComment:   "无效的评论",
	SocialInvalidRating:    "无效的评分",
	SocialInvalidRelation:  "无效的关系",
	SocialAlreadyLiked:     "已点赞",
	SocialAlreadyFollowed:  "已关注",

	// Messaging模块
	MessagingConversationNotFound: "会话不存在",
	MessagingMessageNotFound:      "消息不存在",
	MessagingInvalidRecipient:     "无效的接收者",
	MessagingInvalidContent:       "无效的消息内容",
	MessagingEmptyMessage:         "消息为空",
	MessagingSendFailed:           "消息发送失败",

	// Admin模块
	AdminAuditLogNotFound:    "审计日志不存在",
	AdminReportNotFound:      "举报记录不存在",
	AdminInvalidAction:       "无效的管理操作",
	AdminInvalidStatus:       "无效的状态",
	AdminInvalidPermission:   "无效的权限设置",
	AdminUnauthorized:        "未授权的管理操作",
	AdminInsufficientRole:    "角色权限不足",
}

// GetModuleMessage 获取模块专属错误消息
// 如果没有找到模块消息，返回通用错误消息
func GetModuleMessage(code ErrorCode) string {
	if msg, ok := ModuleErrorMessages[code]; ok {
		return msg
	}
	return GetDefaultMessage(code)
}

// init 初始化模块错误消息到默认消息映射
func init() {
	// 将模块错误消息合并到默认消息映射
	for code, msg := range ModuleErrorMessages {
		if _, exists := DefaultErrorMessages[code]; !exists {
			DefaultErrorMessages[code] = msg
		}
	}
}
