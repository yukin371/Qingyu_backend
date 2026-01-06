package reader

import "time"

// ChapterComment 章节评论（扩展自通用评论）
type ChapterComment struct {
	ID          string    `bson:"_id,omitempty" json:"id"`
	ChapterID   string    `bson:"chapter_id" json:"chapterId"`       // 章节ID
	BookID      string    `bson:"book_id" json:"bookId"`             // 书籍ID
	UserID      string    `bson:"user_id" json:"userId"`             // 用户ID
	Content     string    `bson:"content" json:"content"`            // 评论内容
	Rating      int       `bson:"rating" json:"rating"`              // 评分（1-5，0表示无评分）

	// 段落级评论
	ParagraphIndex *int     `bson:"paragraph_index,omitempty" json:"paragraphIndex,omitempty"` // 段落索引（从0开始）
	ParagraphText   *string  `bson:"paragraph_text,omitempty" json:"paragraphText,omitempty"`   // 段落文本摘要
	CharStart       *int     `bson:"char_start,omitempty" json:"charStart,omitempty"`           // 字符起始位置
	CharEnd         *int     `bson:"char_end,omitempty" json:"charEnd,omitempty"`                 // 字符结束位置

	// 回复
	ParentID   *string    `bson:"parent_id,omitempty" json:"parentId,omitempty"`       // 父评论ID
	RootID     *string    `bson:"root_id,omitempty" json:"rootId,omitempty"`           // 根评论ID（用于 threaded）
	ReplyCount int        `bson:"reply_count" json:"replyCount"`                       // 回复数量

	// 点赞
	LikeCount  int        `bson:"like_count" json:"likeCount"`                         // 点赞数

	// 状态
	IsVisible  bool       `bson:"is_visible" json:"isVisible"`                         // 是否可见（可由管理员隐藏）
	IsDeleted  bool       `bson:"is_deleted" json:"isDeleted"`                         // 是否已删除

	// 时间
	CreatedAt  time.Time  `bson:"created_at" json:"createdAt"`
	UpdatedAt  time.Time  `bson:"updated_at" json:"updatedAt"`

	// 用户快照（避免联表查询）
	UserSnapshot *CommentUserSnapshot `bson:"user_snapshot,omitempty" json:"userSnapshot,omitempty"`
}

// CommentUserSnapshot 评论用户快照
type CommentUserSnapshot struct {
	ID       string `bson:"id" json:"id"`
	Username string `bson:"username" json:"username"`
	Avatar   string `bson:"avatar,omitempty" json:"avatar,omitempty"`
}

// IsParagraphComment 判断是否为段落级评论
func (c *ChapterComment) IsParagraphComment() bool {
	return c.ParagraphIndex != nil
}

// IsTopLevel 判断是否为顶级评论
func (c *ChapterComment) IsTopLevel() bool {
	return c.ParentID == nil
}

// CanEdit 判断是否可编辑（发布后30分钟内）
func (c *ChapterComment) CanEdit() bool {
	return time.Since(c.CreatedAt) <= 30*time.Minute && !c.IsDeleted
}

// ChapterCommentFilter 章节评论过滤器
type ChapterCommentFilter struct {
	BookID        *string    `json:"bookId,omitempty"`
	ChapterID     *string    `json:"chapterId,omitempty"`
	UserID        *string    `json:"userId,omitempty"`
	ParentID      *string    `json:"parentId,omitempty"`       // 空字符串表示查询顶级评论
	ParagraphIndex *int      `json:"paragraphIndex,omitempty"` // 查询特定段落的评论
	MinRating     *int       `json:"minRating,omitempty"`      // 最低评分
	HasRating     *bool      `json:"hasRating,omitempty"`      // 是否有评分
	SortBy        string     `json:"sortBy,omitempty"`         // 排序字段：created_at/like_count/rating
	SortOrder     string     `json:"sortOrder,omitempty"`      // 排序方向：asc/desc
	Page          int        `json:"page,omitempty"`
	PageSize      int        `json:"pageSize,omitempty"`
}

// CreateChapterCommentRequest 创建章节评论请求
type CreateChapterCommentRequest struct {
	ChapterID      string  `json:"chapterId" validate:"required"`
	BookID         string  `json:"bookId" validate:"required"`
	Content        string  `json:"content" validate:"required,min=1,max=5000"`
	Rating         int     `json:"rating" validate:"min=0,max=5"`
	ParentID       *string `json:"parentId,omitempty"`
	ParagraphIndex *int    `json:"paragraphIndex,omitempty"`
	CharStart      *int    `json:"charStart,omitempty"`
	CharEnd        *int    `json:"charEnd,omitempty"`
}

// UpdateChapterCommentRequest 更新章节评论请求
type UpdateChapterCommentRequest struct {
	Content *string `json:"content" validate:"omitempty,min=1,max=5000"`
	Rating  *int    `json:"rating" validate:"omitempty,min=0,max=5"`
}

// ChapterCommentListResponse 章节评论列表响应
type ChapterCommentListResponse struct {
	Comments    []*ChapterComment `json:"comments"`
	Total       int64             `json:"total"`
	Page        int               `json:"page"`
	PageSize    int               `json:"pageSize"`
	TotalPages  int               `json:"totalPages"`
	AvgRating   float64           `json:"avgRating"`   // 平均评分
	RatingCount int               `json:"ratingCount"` // 评分数量
}

// ParagraphCommentResponse 段落评论响应
type ParagraphCommentResponse struct {
	ParagraphIndex int               `json:"paragraphIndex"`
	ParagraphText  string            `json:"paragraphText"`
	CommentCount   int               `json:"commentCount"`
	Comments       []*ChapterComment `json:"comments"`
}
