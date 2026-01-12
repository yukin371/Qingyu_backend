package social

import (
	"context"
	"fmt"
	"time"

	"Qingyu_backend/models/social"
	socialRepo "Qingyu_backend/repository/interfaces/social"
	"Qingyu_backend/service/base"
)

// ReviewService 书评服务
type ReviewService struct {
	reviewRepo  socialRepo.ReviewRepository
	eventBus    base.EventBus
	serviceName string
	version     string
}

// NewReviewService 创建书评服务实例
func NewReviewService(
	reviewRepo socialRepo.ReviewRepository,
	eventBus base.EventBus,
) *ReviewService {
	return &ReviewService{
		reviewRepo:  reviewRepo,
		eventBus:    eventBus,
		serviceName: "ReviewService",
		version:     "1.0.0",
	}
}

// BaseService 接口实现
func (s *ReviewService) Initialize(ctx context.Context) error { return nil }
func (s *ReviewService) Health(ctx context.Context) error {
	if err := s.reviewRepo.Health(ctx); err != nil {
		return fmt.Errorf("书评Repository健康检查失败: %w", err)
	}
	return nil
}
func (s *ReviewService) Close(ctx context.Context) error { return nil }
func (s *ReviewService) GetServiceName() string          { return s.serviceName }
func (s *ReviewService) GetVersion() string              { return s.version }

// CreateReview 创建书评
func (s *ReviewService) CreateReview(ctx context.Context, bookID, userID, userName, userAvatar, title, content string, rating int, isSpoiler, isPublic bool) (*social.Review, error) {
	if bookID == "" || userID == "" {
		return nil, fmt.Errorf("书籍ID和用户ID不能为空")
	}
	if title == "" {
		return nil, fmt.Errorf("标题不能为空")
	}
	if len(title) > 100 {
		return nil, fmt.Errorf("标题最多100字")
	}
	if content == "" {
		return nil, fmt.Errorf("内容不能为空")
	}
	if len(content) > 5000 {
		return nil, fmt.Errorf("内容最多5000字")
	}
	if rating < 1 || rating > 5 {
		return nil, fmt.Errorf("评分必须在1-5之间")
	}

	review := &social.Review{
		BookID:       bookID,
		UserID:       userID,
		UserName:     userName,
		UserAvatar:   userAvatar,
		Title:        title,
		Content:      content,
		Rating:       rating,
		LikeCount:    0,
		CommentCount: 0,
		IsSpoiler:    isSpoiler,
		IsPublic:     isPublic,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.reviewRepo.CreateReview(ctx, review); err != nil {
		return nil, fmt.Errorf("创建书评失败: %w", err)
	}

	s.publishReviewEvent(ctx, "review.created", userID, bookID, review.ID.Hex())
	return review, nil
}

// GetReviews 获取书评列表
func (s *ReviewService) GetReviews(ctx context.Context, bookID string, page, size int) ([]*social.Review, int64, error) {
	if bookID == "" {
		return nil, 0, fmt.Errorf("书籍ID不能为空")
	}

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	reviews, total, err := s.reviewRepo.GetReviewsByBook(ctx, bookID, page, size)
	if err != nil {
		return nil, 0, fmt.Errorf("获取书评列表失败: %w", err)
	}

	return reviews, total, nil
}

// GetReviewByID 获取书评详情
func (s *ReviewService) GetReviewByID(ctx context.Context, reviewID string) (*social.Review, error) {
	if reviewID == "" {
		return nil, fmt.Errorf("书评ID不能为空")
	}

	review, err := s.reviewRepo.GetReviewByID(ctx, reviewID)
	if err != nil {
		return nil, fmt.Errorf("获取书评详情失败: %w", err)
	}

	return review, nil
}

// UpdateReview 更新书评
func (s *ReviewService) UpdateReview(ctx context.Context, userID, reviewID string, updates map[string]interface{}) error {
	if userID == "" || reviewID == "" {
		return fmt.Errorf("用户ID和书评ID不能为空")
	}

	// 获取书评
	review, err := s.reviewRepo.GetReviewByID(ctx, reviewID)
	if err != nil {
		return fmt.Errorf("获取书评失败: %w", err)
	}

	// 权限检查
	if review.UserID != userID {
		return fmt.Errorf("无权更新该书评")
	}

	// 更新
	if err := s.reviewRepo.UpdateReview(ctx, reviewID, updates); err != nil {
		return fmt.Errorf("更新书评失败: %w", err)
	}

	return nil
}

// DeleteReview 删除书评
func (s *ReviewService) DeleteReview(ctx context.Context, userID, reviewID string) error {
	if userID == "" || reviewID == "" {
		return fmt.Errorf("用户ID和书评ID不能为空")
	}

	// 获取书评
	review, err := s.reviewRepo.GetReviewByID(ctx, reviewID)
	if err != nil {
		return fmt.Errorf("获取书评失败: %w", err)
	}

	// 权限检查
	if review.UserID != userID {
		return fmt.Errorf("无权删除该书评")
	}

	// 删除
	if err := s.reviewRepo.DeleteReview(ctx, reviewID); err != nil {
		return fmt.Errorf("删除书评失败: %w", err)
	}

	s.publishReviewEvent(ctx, "review.deleted", userID, review.BookID, reviewID)
	return nil
}

// LikeReview 点赞书评
func (s *ReviewService) LikeReview(ctx context.Context, userID, reviewID string) error {
	if userID == "" || reviewID == "" {
		return fmt.Errorf("用户ID和书评ID不能为空")
	}

	// 检查是否已点赞
	isLiked, err := s.reviewRepo.IsReviewLiked(ctx, reviewID, userID)
	if err != nil {
		return fmt.Errorf("检查点赞状态失败: %w", err)
	}
	if isLiked {
		return fmt.Errorf("已经点赞过该书评")
	}

	// 创建点赞
	reviewLike := &social.ReviewLike{
		ReviewID:  reviewID,
		UserID:    userID,
		CreatedAt: time.Now(),
	}

	if err := s.reviewRepo.CreateReviewLike(ctx, reviewLike); err != nil {
		return fmt.Errorf("点赞失败: %w", err)
	}

	// 增加点赞数
	if err := s.reviewRepo.IncrementReviewLikeCount(ctx, reviewID); err != nil {
		fmt.Printf("Warning: Failed to increment like count: %v\n", err)
	}

	return nil
}

// UnlikeReview 取消点赞书评
func (s *ReviewService) UnlikeReview(ctx context.Context, userID, reviewID string) error {
	if userID == "" || reviewID == "" {
		return fmt.Errorf("用户ID和书评ID不能为空")
	}

	// 检查是否已点赞
	isLiked, err := s.reviewRepo.IsReviewLiked(ctx, reviewID, userID)
	if err != nil {
		return fmt.Errorf("检查点赞状态失败: %w", err)
	}
	if !isLiked {
		return fmt.Errorf("未点赞该书评")
	}

	// 删除点赞
	if err := s.reviewRepo.DeleteReviewLike(ctx, reviewID, userID); err != nil {
		return fmt.Errorf("取消点赞失败: %w", err)
	}

	// 减少点赞数
	if err := s.reviewRepo.DecrementReviewLikeCount(ctx, reviewID); err != nil {
		fmt.Printf("Warning: Failed to decrement like count: %v\n", err)
	}

	return nil
}

func (s *ReviewService) publishReviewEvent(ctx context.Context, eventType, userID, bookID, reviewID string) {
	if s.eventBus == nil {
		return
	}

	event := &base.BaseEvent{
		EventType: eventType,
		EventData: map[string]interface{}{
			"user_id":   userID,
			"book_id":   bookID,
			"review_id": reviewID,
		},
		Timestamp: time.Now(),
		Source:    s.serviceName,
	}

	s.eventBus.PublishAsync(ctx, event)
}
