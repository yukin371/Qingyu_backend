package reading

import (
	"Qingyu_backend/models/community"
	"context"
	"fmt"
	"time"

	readingRepo "Qingyu_backend/repository/interfaces/reading"
	"Qingyu_backend/service/base"
)

// LikeService 点赞服务
type LikeService struct {
	likeRepo    readingRepo.LikeRepository
	commentRepo readingRepo.CommentRepository
	eventBus    base.EventBus
	serviceName string
	version     string
}

// NewLikeService 创建点赞服务实例
func NewLikeService(
	likeRepo readingRepo.LikeRepository,
	commentRepo readingRepo.CommentRepository,
	eventBus base.EventBus,
) *LikeService {
	return &LikeService{
		likeRepo:    likeRepo,
		commentRepo: commentRepo,
		eventBus:    eventBus,
		serviceName: "LikeService",
		version:     "1.0.0",
	}
}

// =========================
// BaseService 接口实现
// =========================

// Initialize 初始化服务
func (s *LikeService) Initialize(ctx context.Context) error {
	return nil
}

// Health 健康检查
func (s *LikeService) Health(ctx context.Context) error {
	if err := s.likeRepo.Health(ctx); err != nil {
		return fmt.Errorf("点赞Repository健康检查失败: %w", err)
	}
	return nil
}

// Close 关闭服务
func (s *LikeService) Close(ctx context.Context) error {
	return nil
}

// GetServiceName 获取服务名称
func (s *LikeService) GetServiceName() string {
	return s.serviceName
}

// GetVersion 获取服务版本
func (s *LikeService) GetVersion() string {
	return s.version
}

// =========================
// 书籍点赞方法
// =========================

// LikeBook 点赞书籍
func (s *LikeService) LikeBook(ctx context.Context, userID, bookID string) error {
	if userID == "" {
		return fmt.Errorf("用户ID不能为空")
	}
	if bookID == "" {
		return fmt.Errorf("书籍ID不能为空")
	}

	// 创建点赞记录
	like := &community.Like{
		UserID:     userID,
		TargetType: community.LikeTargetTypeBook,
		TargetID:   bookID,
		CreatedAt:  time.Now(),
	}

	// 添加点赞
	if err := s.likeRepo.AddLike(ctx, like); err != nil {
		if err.Error() == "已经点赞过了" {
			// 幂等性：已点赞不报错
			return nil
		}
		return fmt.Errorf("点赞书籍失败: %w", err)
	}

	// 发布点赞事件
	s.publishLikeEvent(ctx, "like.book.added", userID, community.LikeTargetTypeBook, bookID)

	return nil
}

// UnlikeBook 取消点赞书籍
func (s *LikeService) UnlikeBook(ctx context.Context, userID, bookID string) error {
	if userID == "" {
		return fmt.Errorf("用户ID不能为空")
	}
	if bookID == "" {
		return fmt.Errorf("书籍ID不能为空")
	}

	// 取消点赞
	if err := s.likeRepo.RemoveLike(ctx, userID, community.LikeTargetTypeBook, bookID); err != nil {
		if err.Error() == "点赞记录不存在" {
			// 幂等性：未点赞不报错
			return nil
		}
		return fmt.Errorf("取消点赞书籍失败: %w", err)
	}

	// 发布取消点赞事件
	s.publishLikeEvent(ctx, "like.book.removed", userID, community.LikeTargetTypeBook, bookID)

	return nil
}

// GetBookLikeCount 获取书籍点赞数
func (s *LikeService) GetBookLikeCount(ctx context.Context, bookID string) (int64, error) {
	if bookID == "" {
		return 0, fmt.Errorf("书籍ID不能为空")
	}

	count, err := s.likeRepo.GetLikeCount(ctx, community.LikeTargetTypeBook, bookID)
	if err != nil {
		return 0, fmt.Errorf("获取书籍点赞数失败: %w", err)
	}

	return count, nil
}

// IsBookLiked 检查书籍是否已点赞
func (s *LikeService) IsBookLiked(ctx context.Context, userID, bookID string) (bool, error) {
	if userID == "" || bookID == "" {
		return false, fmt.Errorf("用户ID和书籍ID不能为空")
	}

	liked, err := s.likeRepo.IsLiked(ctx, userID, community.LikeTargetTypeBook, bookID)
	if err != nil {
		return false, fmt.Errorf("检查点赞状态失败: %w", err)
	}

	return liked, nil
}

// =========================
// 评论点赞方法
// =========================

// LikeComment 点赞评论
func (s *LikeService) LikeComment(ctx context.Context, userID, commentID string) error {
	if userID == "" {
		return fmt.Errorf("用户ID不能为空")
	}
	if commentID == "" {
		return fmt.Errorf("评论ID不能为空")
	}

	// 创建点赞记录
	like := &community.Like{
		UserID:     userID,
		TargetType: community.LikeTargetTypeComment,
		TargetID:   commentID,
		CreatedAt:  time.Now(),
	}

	// 添加点赞
	if err := s.likeRepo.AddLike(ctx, like); err != nil {
		if err.Error() == "已经点赞过了" {
			return nil
		}
		return fmt.Errorf("点赞评论失败: %w", err)
	}

	// 更新评论点赞数
	if s.commentRepo != nil {
		if err := s.commentRepo.IncrementLikeCount(ctx, commentID); err != nil {
			fmt.Printf("Warning: Failed to increment comment like count: %v\n", err)
		}
	}

	// 发布点赞事件
	s.publishLikeEvent(ctx, "like.comment.added", userID, community.LikeTargetTypeComment, commentID)

	return nil
}

// UnlikeComment 取消点赞评论
func (s *LikeService) UnlikeComment(ctx context.Context, userID, commentID string) error {
	if userID == "" {
		return fmt.Errorf("用户ID不能为空")
	}
	if commentID == "" {
		return fmt.Errorf("评论ID不能为空")
	}

	// 取消点赞
	if err := s.likeRepo.RemoveLike(ctx, userID, community.LikeTargetTypeComment, commentID); err != nil {
		if err.Error() == "点赞记录不存在" {
			return nil
		}
		return fmt.Errorf("取消点赞评论失败: %w", err)
	}

	// 更新评论点赞数
	if s.commentRepo != nil {
		if err := s.commentRepo.DecrementLikeCount(ctx, commentID); err != nil {
			fmt.Printf("Warning: Failed to decrement comment like count: %v\n", err)
		}
	}

	// 发布取消点赞事件
	s.publishLikeEvent(ctx, "like.comment.removed", userID, community.LikeTargetTypeComment, commentID)

	return nil
}

// =========================
// 用户点赞列表
// =========================

// GetUserLikedBooks 获取用户点赞的书籍列表
func (s *LikeService) GetUserLikedBooks(ctx context.Context, userID string, page, size int) ([]*community.Like, int64, error) {
	if userID == "" {
		return nil, 0, fmt.Errorf("用户ID不能为空")
	}

	// 参数验证和默认值
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	likes, total, err := s.likeRepo.GetUserLikes(ctx, userID, community.LikeTargetTypeBook, page, size)
	if err != nil {
		return nil, 0, fmt.Errorf("获取用户点赞书籍列表失败: %w", err)
	}

	return likes, total, nil
}

// GetUserLikedComments 获取用户点赞的评论列表
func (s *LikeService) GetUserLikedComments(ctx context.Context, userID string, page, size int) ([]*community.Like, int64, error) {
	if userID == "" {
		return nil, 0, fmt.Errorf("用户ID不能为空")
	}

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	likes, total, err := s.likeRepo.GetUserLikes(ctx, userID, community.LikeTargetTypeComment, page, size)
	if err != nil {
		return nil, 0, fmt.Errorf("获取用户点赞评论列表失败: %w", err)
	}

	return likes, total, nil
}

// =========================
// 批量查询方法
// =========================

// GetBooksLikeCount 批量获取书籍点赞数
func (s *LikeService) GetBooksLikeCount(ctx context.Context, bookIDs []string) (map[string]int64, error) {
	if len(bookIDs) == 0 {
		return make(map[string]int64), nil
	}

	counts, err := s.likeRepo.GetLikesCountBatch(ctx, community.LikeTargetTypeBook, bookIDs)
	if err != nil {
		return nil, fmt.Errorf("批量获取书籍点赞数失败: %w", err)
	}

	return counts, nil
}

// GetUserLikeStatus 批量检查用户点赞状态
func (s *LikeService) GetUserLikeStatus(ctx context.Context, userID string, bookIDs []string) (map[string]bool, error) {
	if userID == "" || len(bookIDs) == 0 {
		return make(map[string]bool), nil
	}

	status, err := s.likeRepo.GetUserLikeStatusBatch(ctx, userID, community.LikeTargetTypeBook, bookIDs)
	if err != nil {
		return nil, fmt.Errorf("批量检查点赞状态失败: %w", err)
	}

	return status, nil
}

// =========================
// 统计方法
// =========================

// GetUserLikeStats 获取用户点赞统计
func (s *LikeService) GetUserLikeStats(ctx context.Context, userID string) (map[string]interface{}, error) {
	if userID == "" {
		return nil, fmt.Errorf("用户ID不能为空")
	}

	totalCount, err := s.likeRepo.CountUserLikes(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("获取用户点赞统计失败: %w", err)
	}

	stats := map[string]interface{}{
		"total_likes": totalCount,
	}

	return stats, nil
}

// =========================
// 私有辅助方法
// =========================

// publishLikeEvent 发布点赞事件
func (s *LikeService) publishLikeEvent(ctx context.Context, eventType string, userID, targetType, targetID string) {
	if s.eventBus == nil {
		return
	}

	// 获取点赞数（用于事件数据）
	likeCount, _ := s.likeRepo.GetLikeCount(ctx, targetType, targetID)

	event := &base.BaseEvent{
		EventType: eventType,
		EventData: map[string]interface{}{
			"user_id":     userID,
			"target_type": targetType,
			"target_id":   targetID,
			"like_count":  likeCount,
		},
		Timestamp: time.Now(),
		Source:    s.serviceName,
	}

	s.eventBus.PublishAsync(ctx, event)
}
