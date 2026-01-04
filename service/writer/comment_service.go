package writer

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/models/writer"
	writerrepo "Qingyu_backend/repository/interfaces/writer"
)

var (
	// ErrCommentNotFound 批注不存在
	ErrCommentNotFound = fmt.Errorf("comment not found")
	// ErrCommentAlreadyResolved 批注已解决
	ErrCommentAlreadyResolved = fmt.Errorf("comment already resolved")
	// ErrInvalidPosition 位置信息无效
	ErrInvalidPosition = fmt.Errorf("invalid position")
)

// CommentService 批注服务接口
type CommentService interface {
	// CreateComment 创建批注
	CreateComment(ctx context.Context, comment *writer.DocumentComment) (*writer.DocumentComment, error)

	// GetComment 获取批注详情
	GetComment(ctx context.Context, commentID string) (*writer.DocumentComment, error)

	// UpdateComment 更新批注
	UpdateComment(ctx context.Context, commentID string, comment *writer.DocumentComment) error

	// DeleteComment 删除批注
	DeleteComment(ctx context.Context, commentID string) error

	// ListComments 查询批注列表
	ListComments(ctx context.Context, filter *writer.CommentFilter, page, pageSize int) ([]*writer.DocumentComment, int64, error)

	// GetDocumentComments 获取文档的所有批注
	GetDocumentComments(ctx context.Context, documentID string, includeResolved bool) ([]*writer.DocumentComment, error)

	// GetChapterComments 获取章节的所有批注
	GetChapterComments(ctx context.Context, chapterID string, includeResolved bool) ([]*writer.DocumentComment, error)

	// ResolveComment 标记批注为已解决
	ResolveComment(ctx context.Context, commentID, userID string) error

	// UnresolveComment 标记批注为未解决
	UnresolveComment(ctx context.Context, commentID string) error

	// ReplyComment 回复批注
	ReplyComment(ctx context.Context, parentID, content, userID, userName string) (*writer.DocumentComment, error)

	// GetCommentThread 获取批注线程
	GetCommentThread(ctx context.Context, threadID string) (*writer.CommentThread, error)

	// GetCommentReplies 获取批注的回复
	GetCommentReplies(ctx context.Context, parentID string) ([]*writer.DocumentComment, error)

	// GetCommentStats 获取批注统计
	GetCommentStats(ctx context.Context, documentID string) (*writer.CommentStats, error)

	// BatchDeleteComments 批量删除批注
	BatchDeleteComments(ctx context.Context, commentIDs []string) error

	// SearchComments 搜索批注
	SearchComments(ctx context.Context, keyword, documentID string, page, pageSize int) ([]*writer.DocumentComment, int64, error)
}

// CommentServiceImpl 批注服务实现
type CommentServiceImpl struct {
	commentRepo writerrepo.CommentRepository
}

// NewCommentService 创建批注服务
func NewCommentService(commentRepo writerrepo.CommentRepository) CommentService {
	return &CommentServiceImpl{
		commentRepo: commentRepo,
	}
}

// CreateComment 创建批注
func (s *CommentServiceImpl) CreateComment(ctx context.Context, comment *writer.DocumentComment) (*writer.DocumentComment, error) {
	// 验证位置信息
	if err := s.validatePosition(comment.Position); err != nil {
		return nil, err
	}

	// 验证类型
	if comment.Type == "" {
		comment.Type = writer.CommentTypeComment // 默认为普通评论
	}

	// 如果是回复，验证父评论是否存在
	if comment.ParentID != nil {
		parent, err := s.commentRepo.GetByID(ctx, *comment.ParentID)
		if err != nil {
			return nil, fmt.Errorf("parent comment not found: %w", err)
		}
		// 继承线程ID
		comment.ThreadID = parent.ThreadID
	}

	if err := s.commentRepo.Create(ctx, comment); err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	return comment, nil
}

// GetComment 获取批注详情
func (s *CommentServiceImpl) GetComment(ctx context.Context, commentID string) (*writer.DocumentComment, error) {
	id, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		return nil, fmt.Errorf("invalid comment ID: %w", err)
	}

	comment, err := s.commentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrCommentNotFound
	}

	return comment, nil
}

// UpdateComment 更新批注
func (s *CommentServiceImpl) UpdateComment(ctx context.Context, commentID string, comment *writer.DocumentComment) error {
	id, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		return fmt.Errorf("invalid comment ID: %w", err)
	}

	// 只允许更新内容和元数据
	existing, err := s.commentRepo.GetByID(ctx, id)
	if err != nil {
		return ErrCommentNotFound
	}

	// 更新允许的字段
	existing.Content = comment.Content
	existing.Type = comment.Type
	existing.Position = comment.Position
	existing.Metadata = comment.Metadata
	existing.UpdatedAt = time.Now()

	return s.commentRepo.Update(ctx, id, existing)
}

// DeleteComment 删除批注
func (s *CommentServiceImpl) DeleteComment(ctx context.Context, commentID string) error {
	id, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		return fmt.Errorf("invalid comment ID: %w", err)
	}

	return s.commentRepo.Delete(ctx, id)
}

// ListComments 查询批注列表
func (s *CommentServiceImpl) ListComments(ctx context.Context, filter *writer.CommentFilter, page, pageSize int) ([]*writer.DocumentComment, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	return s.commentRepo.List(ctx, filter, page, pageSize)
}

// GetDocumentComments 获取文档的所有批注
func (s *CommentServiceImpl) GetDocumentComments(ctx context.Context, documentID string, includeResolved bool) ([]*writer.DocumentComment, error) {
	id, err := primitive.ObjectIDFromHex(documentID)
	if err != nil {
		return nil, fmt.Errorf("invalid document ID: %w", err)
	}

	return s.commentRepo.GetByDocument(ctx, id, includeResolved)
}

// GetChapterComments 获取章节的所有批注
func (s *CommentServiceImpl) GetChapterComments(ctx context.Context, chapterID string, includeResolved bool) ([]*writer.DocumentComment, error) {
	id, err := primitive.ObjectIDFromHex(chapterID)
	if err != nil {
		return nil, fmt.Errorf("invalid chapter ID: %w", err)
	}

	return s.commentRepo.GetByChapter(ctx, id, includeResolved)
}

// ResolveComment 标记批注为已解决
func (s *CommentServiceImpl) ResolveComment(ctx context.Context, commentID, userID string) error {
	commentIDObj, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		return fmt.Errorf("invalid comment ID: %w", err)
	}

	userIDObj, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	comment, err := s.commentRepo.GetByID(ctx, commentIDObj)
	if err != nil {
		return ErrCommentNotFound
	}

	if comment.Resolved {
		return ErrCommentAlreadyResolved
	}

	return s.commentRepo.MarkAsResolved(ctx, commentIDObj, userIDObj)
}

// UnresolveComment 标记批注为未解决
func (s *CommentServiceImpl) UnresolveComment(ctx context.Context, commentID string) error {
	id, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		return fmt.Errorf("invalid comment ID: %w", err)
	}

	return s.commentRepo.MarkAsUnresolved(ctx, id)
}

// ReplyComment 回复批注
func (s *CommentServiceImpl) ReplyComment(ctx context.Context, parentID, content, userID, userName string) (*writer.DocumentComment, error) {
	parentIDObj, err := primitive.ObjectIDFromHex(parentID)
	if err != nil {
		return nil, fmt.Errorf("invalid parent comment ID: %w", err)
	}

	userIDObj, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	// 获取父评论
	parent, err := s.commentRepo.GetByID(ctx, parentIDObj)
	if err != nil {
		return nil, fmt.Errorf("parent comment not found: %w", err)
	}

	// 创建回复
	reply := &writer.DocumentComment{
		DocumentID: parent.DocumentID,
		ChapterID:  parent.ChapterID,
		UserID:     userIDObj,
		UserName:   userName,
		Content:    content,
		Type:       writer.CommentTypeComment,
		Position:   parent.Position, // 继承父评论的位置
		ParentID:   &parentIDObj,
		ReplyTo:    &parentIDObj,
		ThreadID:   parent.ThreadID,
		Resolved:   false,
		Metadata:   writer.CommentMetadata{},
	}

	if err := s.commentRepo.Create(ctx, reply); err != nil {
		return nil, fmt.Errorf("failed to create reply: %w", err)
	}

	return reply, nil
}

// GetCommentThread 获取批注线程
func (s *CommentServiceImpl) GetCommentThread(ctx context.Context, threadID string) (*writer.CommentThread, error) {
	id, err := primitive.ObjectIDFromHex(threadID)
	if err != nil {
		return nil, fmt.Errorf("invalid thread ID: %w", err)
	}

	return s.commentRepo.GetThread(ctx, id)
}

// GetCommentReplies 获取批注的回复
func (s *CommentServiceImpl) GetCommentReplies(ctx context.Context, parentID string) ([]*writer.DocumentComment, error) {
	id, err := primitive.ObjectIDFromHex(parentID)
	if err != nil {
		return nil, fmt.Errorf("invalid parent comment ID: %w", err)
	}

	return s.commentRepo.GetReplies(ctx, id)
}

// GetCommentStats 获取批注统计
func (s *CommentServiceImpl) GetCommentStats(ctx context.Context, documentID string) (*writer.CommentStats, error) {
	id, err := primitive.ObjectIDFromHex(documentID)
	if err != nil {
		return nil, fmt.Errorf("invalid document ID: %w", err)
	}

	return s.commentRepo.GetStats(ctx, id)
}

// BatchDeleteComments 批量删除批注
func (s *CommentServiceImpl) BatchDeleteComments(ctx context.Context, commentIDs []string) error {
	ids := make([]primitive.ObjectID, 0, len(commentIDs))
	for _, idStr := range commentIDs {
		id, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			continue // 跳过无效ID
		}
		ids = append(ids, id)
	}

	if len(ids) == 0 {
		return fmt.Errorf("no valid comment IDs provided")
	}

	return s.commentRepo.BatchDelete(ctx, ids)
}

// SearchComments 搜索批注
func (s *CommentServiceImpl) SearchComments(ctx context.Context, keyword, documentID string, page, pageSize int) ([]*writer.DocumentComment, int64, error) {
	id, err := primitive.ObjectIDFromHex(documentID)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid document ID: %w", err)
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	return s.commentRepo.Search(ctx, keyword, id, page, pageSize)
}

// validatePosition 验证位置信息
func (s *CommentServiceImpl) validatePosition(pos writer.CommentPosition) error {
	if pos.ChapterID.IsZero() {
		return fmt.Errorf("%w: chapter ID is required", ErrInvalidPosition)
	}
	if pos.Paragraph < 0 {
		return fmt.Errorf("%w: paragraph index cannot be negative", ErrInvalidPosition)
	}
	if pos.Offset < 0 {
		return fmt.Errorf("%w: offset cannot be negative", ErrInvalidPosition)
	}
	if pos.Length <= 0 {
		return fmt.Errorf("%w: length must be positive", ErrInvalidPosition)
	}
	return nil
}
