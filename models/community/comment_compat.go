package community

import socialmodels "Qingyu_backend/models/social"

// 向后兼容的类型别名 - 指向 social 模块

// Comment 评论模型
type Comment = socialmodels.Comment

// CommentTargetType 评论目标类型
type CommentTargetType = socialmodels.CommentTargetType

// CommentState 评论状态
type CommentState = socialmodels.CommentState

// CommentAuthorSnapshot 评论作者快照
type CommentAuthorSnapshot = socialmodels.CommentAuthorSnapshot

// CommentFilter 评论查询过滤器
type CommentFilter = socialmodels.CommentFilter

// CommentThread 评论线程
type CommentThread = socialmodels.CommentThread

// ValidationError 验证错误
type ValidationError = socialmodels.ValidationError

// 常量导出
const (
	CommentTargetTypeBook        = socialmodels.CommentTargetTypeBook
	CommentTargetTypeChapter     = socialmodels.CommentTargetTypeChapter
	CommentTargetTypeArticle     = socialmodels.CommentTargetTypeArticle
	CommentTargetTypeAnnouncement = socialmodels.CommentTargetTypeAnnouncement
	CommentTargetTypeProject     = socialmodels.CommentTargetTypeProject

	CommentStateNormal  = socialmodels.CommentStateNormal
	CommentStateHidden  = socialmodels.CommentStateHidden
	CommentStateDeleted = socialmodels.CommentStateDeleted
	CommentStateRejected = socialmodels.CommentStateRejected

	CommentStatusPending  = socialmodels.CommentStatusPending
	CommentStatusApproved = socialmodels.CommentStatusApproved
	CommentStatusRejected = socialmodels.CommentStatusRejected

	CommentSortByLatest = socialmodels.CommentSortByLatest
	CommentSortByHot    = socialmodels.CommentSortByHot
)
