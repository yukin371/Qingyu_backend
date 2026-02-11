package impl

import (
	"context"

	"Qingyu_backend/models/writer"
	serviceWriter "Qingyu_backend/service/interfaces/writer"
	writerservice "Qingyu_backend/service/writer"
)

// CollaborationImpl 协作批注端口实现
type CollaborationImpl struct {
	commentService writerservice.CommentService
	serviceName    string
	version        string
}

// NewCollaborationImpl 创建协作批注端口实现
func NewCollaborationImpl(commentService writerservice.CommentService) serviceWriter.CollaborationPort {
	return &CollaborationImpl{
		commentService: commentService,
		serviceName:    "CollaborationPort",
		version:        "1.0.0",
	}
}

// ============================================================================
// BaseService 生命周期方法实现
// ============================================================================

func (c *CollaborationImpl) Initialize(ctx context.Context) error {
	// CommentService 没有初始化方法，返回 nil
	return nil
}

func (c *CollaborationImpl) Health(ctx context.Context) error {
	// CommentService 没有健康检查方法，返回 nil 表示健康
	return nil
}

func (c *CollaborationImpl) Close(ctx context.Context) error {
	// CommentService 没有需要清理的资源
	return nil
}

func (c *CollaborationImpl) GetServiceName() string {
	return c.serviceName
}

func (c *CollaborationImpl) GetVersion() string {
	return c.version
}

// ============================================================================
// CollaborationPort 方法实现
// ============================================================================

// CreateComment 创建批注
func (c *CollaborationImpl) CreateComment(ctx context.Context, comment *writer.DocumentComment) (*writer.DocumentComment, error) {
	return c.commentService.CreateComment(ctx, comment)
}

// GetComment 获取批注详情
func (c *CollaborationImpl) GetComment(ctx context.Context, commentID string) (*writer.DocumentComment, error) {
	return c.commentService.GetComment(ctx, commentID)
}

// UpdateComment 更新批注
func (c *CollaborationImpl) UpdateComment(ctx context.Context, commentID string, comment *writer.DocumentComment) error {
	return c.commentService.UpdateComment(ctx, commentID, comment)
}

// DeleteComment 删除批注
func (c *CollaborationImpl) DeleteComment(ctx context.Context, commentID string) error {
	return c.commentService.DeleteComment(ctx, commentID)
}

// ListComments 查询批注列表
func (c *CollaborationImpl) ListComments(ctx context.Context, filter *writer.CommentFilter, page, pageSize int) ([]*writer.DocumentComment, int64, error) {
	return c.commentService.ListComments(ctx, filter, page, pageSize)
}

// GetDocumentComments 获取文档的所有批注
func (c *CollaborationImpl) GetDocumentComments(ctx context.Context, documentID string, includeResolved bool) ([]*writer.DocumentComment, error) {
	return c.commentService.GetDocumentComments(ctx, documentID, includeResolved)
}

// GetChapterComments 获取章节的所有批注
func (c *CollaborationImpl) GetChapterComments(ctx context.Context, chapterID string, includeResolved bool) ([]*writer.DocumentComment, error) {
	return c.commentService.GetChapterComments(ctx, chapterID, includeResolved)
}

// ResolveComment 标记批注为已解决
func (c *CollaborationImpl) ResolveComment(ctx context.Context, commentID, userID string) error {
	return c.commentService.ResolveComment(ctx, commentID, userID)
}

// UnresolveComment 标记批注为未解决
func (c *CollaborationImpl) UnresolveComment(ctx context.Context, commentID string) error {
	return c.commentService.UnresolveComment(ctx, commentID)
}

// ReplyComment 回复批注
func (c *CollaborationImpl) ReplyComment(ctx context.Context, parentID, content, userID, userName string) (*writer.DocumentComment, error) {
	return c.commentService.ReplyComment(ctx, parentID, content, userID, userName)
}

// GetCommentThread 获取批注线程
func (c *CollaborationImpl) GetCommentThread(ctx context.Context, threadID string) (*writer.CommentThread, error) {
	return c.commentService.GetCommentThread(ctx, threadID)
}

// GetCommentReplies 获取批注的回复
func (c *CollaborationImpl) GetCommentReplies(ctx context.Context, parentID string) ([]*writer.DocumentComment, error) {
	return c.commentService.GetCommentReplies(ctx, parentID)
}

// GetCommentStats 获取批注统计
func (c *CollaborationImpl) GetCommentStats(ctx context.Context, documentID string) (*writer.CommentStats, error) {
	return c.commentService.GetCommentStats(ctx, documentID)
}

// BatchDeleteComments 批量删除批注
func (c *CollaborationImpl) BatchDeleteComments(ctx context.Context, commentIDs []string) error {
	return c.commentService.BatchDeleteComments(ctx, commentIDs)
}

// SearchComments 搜索批注
func (c *CollaborationImpl) SearchComments(ctx context.Context, keyword, documentID string, page, pageSize int) ([]*writer.DocumentComment, int64, error) {
	return c.commentService.SearchComments(ctx, keyword, documentID, page, pageSize)
}
