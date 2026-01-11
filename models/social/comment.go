package social

import (
	"time"
)

// Comment 评论模型
type Comment struct {
	IdentifiedEntity     `bson:",inline"`
	Timestamps           `bson:",inline"`
	ThreadedConversation `bson:",inline"`
	Likable              `bson:",inline"`

	// 作者信息
	AuthorID string `bson:"author_id" json:"authorId" validate:"required"` // 评论者ID

	// 目标信息
	TargetType CommentTargetType `bson:"target_type" json:"targetType" validate:"required"` // 目标类型
	TargetID   string            `bson:"target_id" json:"targetId" validate:"required"`     // 目标ID

	// 旧字段兼容（向后兼容）
	BookID    string `bson:"book_id,omitempty" json:"book_id,omitempty"`             // 书籍ID（兼容旧版本）
	ChapterID string `bson:"chapter_id,omitempty" json:"chapter_id,omitempty"`       // 章节ID（兼容旧版本）

	// 评论内容
	Content     string      `bson:"content" json:"content" validate:"required,min=1,max=5000"` // 评论内容
	RichContent interface{} `bson:"rich_content,omitempty" json:"richContent,omitempty"`     // 富文本内容（JSON）
	Rating      int         `bson:"rating" json:"rating" validate:"min=0,max=5"`              // 评分（0-5星，0表示无评分）

	// 状态
	State       CommentState `bson:"state" json:"state" validate:"required,oneof=normal hidden deleted rejected"` // 评论状态
	RejectReason string      `bson:"reject_reason,omitempty" json:"rejectReason,omitempty"`                       // 拒绝原因

	// 作者快照
	AuthorSnapshot *CommentAuthorSnapshot `bson:"author_snapshot,omitempty" json:"authorSnapshot,omitempty"`

	// 回复目标
	ReplyToUserID      *string                `bson:"reply_to_user_id,omitempty" json:"replyToUserId,omitempty"`       // 被回复的用户ID
	ReplyToUserSnapshot *CommentAuthorSnapshot `bson:"reply_to_user_snapshot,omitempty" json:"replyToUserSnapshot,omitempty"` // 被回复用户快照
	ReplyToCommentID   *string                `bson:"reply_to_comment_id,omitempty" json:"replyToCommentId,omitempty"` // 被回复的评论ID
	ReplyToContent     *string                `bson:"reply_to_content,omitempty" json:"replyToContent,omitempty"`     // 被回复的内容摘要

	// 管理信息
	IsPinned      bool       `bson:"is_pinned" json:"isPinned"`                      // 是否置顶
	PinnedAt      *time.Time `bson:"pinned_at,omitempty" json:"pinnedAt,omitempty"` // 置顶时间
	IsAuthorReply bool       `bson:"is_author_reply" json:"isAuthorReply"`           // 是否是作者回复
	IsFeatured    bool       `bson:"is_featured" json:"isFeatured"`                  // 是否精选评论
}

// CommentTargetType 评论目标类型
type CommentTargetType string

const (
	CommentTargetTypeBook        CommentTargetType = "book"        // 书籍评论
	CommentTargetTypeChapter     CommentTargetType = "chapter"     // 章节评论
	CommentTargetTypeArticle     CommentTargetType = "article"     // 文章评论
	CommentTargetTypeAnnouncement CommentTargetType = "announcement" // 公告评论
	CommentTargetTypeProject     CommentTargetType = "project"     // 项目评论
)

// CommentState 评论状态
type CommentState string

const (
	CommentStateNormal  CommentState = "normal"  // 正常
	CommentStateHidden  CommentState = "hidden"  // 已隐藏
	CommentStateDeleted CommentState = "deleted" // 已删除
	CommentStateRejected CommentState = "rejected" // 已拒绝（审核未通过）
)

// CommentStatus 评论状态常量（向后兼容）
const (
	CommentStatusPending  = "pending"  // 待审核（保留但不使用）
	CommentStatusApproved = "normal"   // 已通过（映射到normal）
	CommentStatusRejected = "rejected" // 已拒绝
)

// CommentSortBy 评论排序方式
const (
	CommentSortByLatest = "latest" // 最新
	CommentSortByHot    = "hot"    // 最热（点赞数）
)

// CommentAuthorSnapshot 评论作者快照
type CommentAuthorSnapshot struct {
	ID       string `bson:"id" json:"id"`
	Username string `bson:"username" json:"username"`
	Avatar   string `bson:"avatar,omitempty" json:"avatar,omitempty"`
}

// CommentFilter 评论查询过滤器
type CommentFilter struct {
	TargetType    *CommentTargetType `json:"targetType,omitempty"`
	TargetID      *string            `json:"targetId,omitempty"`
	BookID        *string            `json:"bookId,omitempty"`         // 兼容旧版本
	ChapterID     *string            `json:"chapterId,omitempty"`      // 兼容旧版本
	AuthorID      *string            `json:"authorId,omitempty"`
	State         *CommentState      `json:"state,omitempty"`
	ParentID      *string            `json:"parentId,omitempty"`       // 只查询顶级评论或指定评论的回复
	IsPinned      *bool              `json:"isPinned,omitempty"`
	IsFeatured    *bool              `json:"isFeatured,omitempty"`
	IsAuthorReply *bool              `json:"isAuthorReply,omitempty"`
	HasRating     *bool              `json:"hasRating,omitempty"`      // 是否有评分
	StartTime     *time.Time         `json:"startTime,omitempty"`
	EndTime       *time.Time         `json:"endTime,omitempty"`
	SortBy        string             `json:"sortBy,omitempty"`
	SortOrder     string             `json:"sortOrder,omitempty"`
	Limit         int                `json:"limit,omitempty"`
	Offset        int                `json:"offset,omitempty"`
}

// GetConditions 获取查询条件
func (f *CommentFilter) GetConditions() map[string]interface{} {
	conditions := make(map[string]interface{})

	if f.TargetType != nil {
		conditions["target_type"] = *f.TargetType
	} else if f.BookID != nil {
		// 兼容旧版本：book_id映射到target_type
		conditions["target_type"] = CommentTargetTypeBook
	}

	if f.TargetID != nil {
		conditions["target_id"] = *f.TargetID
	} else if f.BookID != nil {
		// 兼容旧版本：book_id映射到target_id
		conditions["target_id"] = *f.BookID
	}

	if f.ChapterID != nil {
		conditions["chapter_id"] = *f.ChapterID
	}

	if f.AuthorID != nil {
		conditions["author_id"] = *f.AuthorID
	}
	if f.State != nil {
		conditions["state"] = *f.State
	}
	if f.IsPinned != nil {
		conditions["is_pinned"] = *f.IsPinned
	}
	if f.IsFeatured != nil {
		conditions["is_featured"] = *f.IsFeatured
	}
	if f.IsAuthorReply != nil {
		conditions["is_author_reply"] = *f.IsAuthorReply
	}
	if f.HasRating != nil {
		if *f.HasRating {
			conditions["rating"] = map[string]interface{}{"$gt": 0}
		} else {
			conditions["rating"] = 0
		}
	}

	// 父评论筛选
	if f.ParentID != nil {
		if *f.ParentID == "" {
			// 空字符串表示查询顶级评论
			conditions["parent_id"] = nil
		} else {
			conditions["parent_id"] = *f.ParentID
		}
	}

	// 时间范围
	if f.StartTime != nil || f.EndTime != nil {
		timeCondition := make(map[string]interface{})
		if f.StartTime != nil {
			timeCondition["$gte"] = *f.StartTime
		}
		if f.EndTime != nil {
			timeCondition["$lte"] = *f.EndTime
		}
		conditions["created_at"] = timeCondition
	}

	// 默认不显示已删除的评论
	if f.State == nil {
		conditions["state"] = map[string]interface{}{"$ne": CommentStateDeleted}
	}

	return conditions
}

// GetSort 获取排序
func (f *CommentFilter) GetSort() map[string]int {
	sort := make(map[string]int)

	sortValue := -1 // 默认降序
	if f.SortOrder == "asc" {
		sortValue = 1
	}

	switch f.SortBy {
	case "created_at", "latest":
		sort["created_at"] = sortValue
	case "like_count", "hot":
		sort["like_count"] = sortValue
	case "rating":
		sort["rating"] = sortValue
	case "is_pinned":
		sort["is_pinned"] = -1 // 置顶始终在前
	default:
		// 默认排序：置顶优先，然后按点赞数，最后按时间
		sort["is_pinned"] = -1
		sort["like_count"] = -1
		sort["created_at"] = -1
	}

	return sort
}

// GetLimit 获取限制
func (f *CommentFilter) GetLimit() int {
	return f.Limit
}

// GetOffset 获取偏移
func (f *CommentFilter) GetOffset() int {
	return f.Offset
}

// GetFields 获取字段
func (f *CommentFilter) GetFields() []string {
	return []string{}
}

// Validate 验证
func (f *CommentFilter) Validate() error {
	if f.Limit < 0 {
		return &ValidationError{Message: "limit不能为负数"}
	}
	if f.Offset < 0 {
		return &ValidationError{Message: "offset不能为负数"}
	}
	return nil
}

// IsTopLevel 判断是否为顶级评论
func (c *Comment) IsTopLevel() bool {
	return c.ParentID == nil
}

// IsReply 判断是否为回复
func (c *Comment) IsReply() bool {
	return c.ParentID != nil
}

// IsEditable 检查评论是否可编辑（30分钟内）
func (c *Comment) IsEditable() bool {
	return time.Since(c.CreatedAt) <= 30*time.Minute && c.State == CommentStateNormal
}

// HasRating 判断是否包含评分
func (c *Comment) HasRating() bool {
	return c.Rating > 0
}

// IsVisible 判断评论是否可见
func (c *Comment) IsVisible() bool {
	return c.State == CommentStateNormal
}

// CanEdit 判断是否可以编辑
func (c *Comment) CanEdit(userID string) bool {
	// 只有作者可以编辑，且只能编辑正常状态的评论
	return c.AuthorID == userID && c.State == CommentStateNormal
}

// CanDelete 判断是否可以删除
func (c *Comment) CanDelete(userID string, isAdmin bool) bool {
	// 作者或管理员可以删除
	return c.AuthorID == userID || isAdmin
}

// Hide 隐藏评论
func (c *Comment) Hide() {
	c.State = CommentStateHidden
	c.Touch()
}

// Show 显示评论
func (c *Comment) Show() {
	c.State = CommentStateNormal
	c.Touch()
}

// Delete 删除评论
func (c *Comment) Delete() {
	c.State = CommentStateDeleted
	c.Touch()
}

// Reject 拒绝评论（审核未通过）
func (c *Comment) Reject(reason string) {
	c.State = CommentStateRejected
	c.RejectReason = reason
	c.Touch()
}

// Pin 置顶评论
func (c *Comment) Pin() {
	c.IsPinned = true
	now := time.Now()
	c.PinnedAt = &now
	c.Touch()
}

// Unpin 取消置顶
func (c *Comment) Unpin() {
	c.IsPinned = false
	c.PinnedAt = nil
	c.Touch()
}

// Feature 精选评论
func (c *Comment) Feature() {
	c.IsFeatured = true
	c.Touch()
}

// Unfeature 取消精选
func (c *Comment) Unfeature() {
	c.IsFeatured = false
	c.Touch()
}

// MarkAsAuthorReply 标记为作者回复
func (c *Comment) MarkAsAuthorReply() {
	c.IsAuthorReply = true
	c.Touch()
}

// SetReplyTarget 设置回复目标
func (c *Comment) SetReplyTarget(parentComment *Comment, replyToUserSnapshot *CommentAuthorSnapshot) {
	c.ParentID = &parentComment.ID
	if parentComment.RootID != nil {
		c.RootID = parentComment.RootID
	} else {
		rootID := parentComment.ID
		c.RootID = &rootID
	}

	c.ReplyToCommentID = &parentComment.ID
	c.ReplyToUserID = &parentComment.AuthorID
	c.ReplyToUserSnapshot = replyToUserSnapshot

	// 截取被回复内容作为摘要（最多100字符）
	content := parentComment.Content
	if len(content) > 100 {
		content = content[:100] + "..."
	}
	c.ReplyToContent = &content

	parentComment.AddReply()
	parentComment.Touch()
}

// CanBeReplied 判断是否可以被回复
func (c *Comment) CanBeReplied() bool {
	// 已删除、隐藏或拒绝的评论不能被回复
	return c.State != CommentStateDeleted && c.State != CommentStateRejected
}

// CommentThread 评论线程（树状结构）
type CommentThread struct {
	Comment *Comment         `json:"comment"`  // 主评论
	Replies []*CommentThread `json:"replies"`  // 回复列表（支持嵌套）
	Total   int64            `json:"total"`    // 总回复数
	HasMore bool             `json:"has_more"` // 是否还有更多回复
}

// ValidationError 验证错误
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}
