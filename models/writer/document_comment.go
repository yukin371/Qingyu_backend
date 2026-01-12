package writer

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CommentType 评论类型
type CommentType string

const (
	CommentTypeComment    CommentType = "comment"    // 普通评论
	CommentTypeSuggestion CommentType = "suggestion" // 建议性修改
	CommentTypeQuestion   CommentType = "question"   // 疑问
	CommentTypeApproval   CommentType = "approval"   // 认可
)

// DocumentComment 文档批注
type DocumentComment struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	DocumentID primitive.ObjectID `bson:"document_id" json:"documentId"`
	ChapterID  primitive.ObjectID `bson:"chapter_id,omitempty" json:"chapterId,omitempty"` // 可选，用于章节级批注
	UserID     primitive.ObjectID `bson:"user_id" json:"userId"`
	UserName   string             `bson:"user_name" json:"userName"`
	UserAvatar string             `bson:"user_avatar,omitempty" json:"userAvatar,omitempty"`

	// 批注内容
	Content  string          `bson:"content" json:"content"`
	Type     CommentType     `bson:"type" json:"type"`
	Position CommentPosition `bson:"position" json:"position"` // 位置信息

	// 状态
	Resolved   bool               `bson:"resolved" json:"resolved"`
	ResolvedBy primitive.ObjectID `bson:"resolved_by,omitempty" json:"resolvedBy,omitempty"`
	ResolvedAt *time.Time         `bson:"resolved_at,omitempty" json:"resolvedAt,omitempty"`

	// 线程化
	ParentID *primitive.ObjectID `bson:"parent_id,omitempty" json:"parentId,omitempty"` // 父评论ID
	ReplyTo  *primitive.ObjectID `bson:"reply_to,omitempty" json:"replyTo,omitempty"`   // 回复的评论ID
	ThreadID *primitive.ObjectID `bson:"thread_id,omitempty" json:"threadId,omitempty"` // 线程ID

	// 元数据
	Metadata CommentMetadata `bson:"metadata,omitempty" json:"metadata,omitempty"`

	// 时间戳
	CreatedAt time.Time  `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time  `bson:"updated_at" json:"updatedAt"`
	DeletedAt *time.Time `bson:"deleted_at,omitempty" json:"deletedAt,omitempty"`
}

// CommentPosition 批注位置
type CommentPosition struct {
	ChapterID    primitive.ObjectID `bson:"chapter_id" json:"chapterId"`                           // 章节ID
	Paragraph    int                `bson:"paragraph" json:"paragraph"`                            // 段落索引（从0开始）
	Offset       int                `bson:"offset" json:"offset"`                                  // 字符偏移（从段落开始）
	Length       int                `bson:"length" json:"length"`                                  // 选中长度
	SelectedText string             `bson:"selected_text,omitempty" json:"selectedText,omitempty"` // 选中的文本
	Line         int                `bson:"line,omitempty" json:"line,omitempty"`                  // 行号（可选）
}

// CommentMetadata 批注元数据
type CommentMetadata struct {
	Labels      []string `bson:"labels,omitempty" json:"labels,omitempty"`           // 标签
	Priority    string   `bson:"priority,omitempty" json:"priority,omitempty"`       // 优先级：low/medium/high
	Attachments []string `bson:"attachments,omitempty" json:"attachments,omitempty"` // 附件ID列表
}

// CommentThread 批注线程
type CommentThread struct {
	ThreadID    primitive.ObjectID `bson:"thread_id" json:"threadId"`
	RootComment *DocumentComment   `bson:"root_comment" json:"rootComment"`
	Replies     []DocumentComment  `bson:"replies" json:"replies"`
	ReplyCount  int                `bson:"reply_count" json:"replyCount"`
	Unresolved  int                `bson:"unresolved" json:"unresolved"` // 未解决的回复数
}

// CommentFilter 批注筛选条件
type CommentFilter struct {
	DocumentID *primitive.ObjectID `bson:"document_id,omitempty"`
	ChapterID  *primitive.ObjectID `bson:"chapter_id,omitempty"`
	UserID     *primitive.ObjectID `bson:"user_id,omitempty"`
	Type       CommentType         `bson:"type,omitempty"`
	Resolved   *bool               `bson:"resolved,omitempty"`
	ParentID   *primitive.ObjectID `bson:"parent_id,omitempty"` // nil表示顶级评论
	ThreadID   *primitive.ObjectID `bson:"thread_id,omitempty"`
	StartDate  *time.Time          `bson:"start_date,omitempty"`
	EndDate    *time.Time          `bson:"end_date,omitempty"`
	Keyword    string              `bson:"-"` // 内容关键词搜索
	Labels     []string            `bson:"-"` // 标签筛选
}

// CommentStats 批注统计
type CommentStats struct {
	TotalCount      int            `json:"totalCount"`
	ResolvedCount   int            `json:"resolvedCount"`
	UnresolvedCount int            `json:"unresolvedCount"`
	ByType          map[string]int `json:"byType"`
	ByUser          map[string]int `json:"byUser"` // userID -> count
}
