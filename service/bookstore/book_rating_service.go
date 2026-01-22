package bookstore

import (
	"Qingyu_backend/models/bookstore"
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	BookstoreRepo "Qingyu_backend/repository/interfaces/bookstore"
)

// BookRatingService 书籍评分服务接口
type BookRatingService interface {
	// 评分基础操作
	CreateRating(ctx context.Context, rating *bookstore.BookRating) error
	GetRatingByID(ctx context.Context, id primitive.ObjectID) (*bookstore.BookRating, error)
	UpdateRating(ctx context.Context, rating *bookstore.BookRating) error
	DeleteRating(ctx context.Context, id primitive.ObjectID) error

	// 评分查询
	GetRatingsByBookID(ctx context.Context, bookID primitive.ObjectID, page, pageSize int) ([]*bookstore.BookRating, int64, error)
	GetRatingsByUserID(ctx context.Context, userID primitive.ObjectID, page, pageSize int) ([]*bookstore.BookRating, int64, error)
	GetRatingByBookIDAndUserID(ctx context.Context, bookID, userID primitive.ObjectID) (*bookstore.BookRating, error)
	GetRatingsByRating(ctx context.Context, rating float64, page, pageSize int) ([]*bookstore.BookRating, int64, error)
	GetRatingsByTags(ctx context.Context, tags []string, page, pageSize int) ([]*bookstore.BookRating, int64, error)

	// 评分统计
	GetAverageRating(ctx context.Context, bookID primitive.ObjectID) (float64, error)
	GetRatingDistribution(ctx context.Context, bookID primitive.ObjectID) (map[string]int64, error)
	GetRatingCount(ctx context.Context, bookID primitive.ObjectID) (int64, error)
	GetRatingStats(ctx context.Context, bookID primitive.ObjectID) (map[string]interface{}, error)
	GetTopRatedBooks(ctx context.Context, limit int) ([]*bookstore.BookRating, error)

	// 评分互动
	LikeRating(ctx context.Context, ratingID primitive.ObjectID, userID primitive.ObjectID) error
	UnlikeRating(ctx context.Context, ratingID primitive.ObjectID, userID primitive.ObjectID) error
	GetRatingLikes(ctx context.Context, ratingID primitive.ObjectID) (int64, error)

	// 用户评分管理
	HasUserRated(ctx context.Context, bookID, userID primitive.ObjectID) (bool, error)
	GetUserRatingForBook(ctx context.Context, bookID, userID primitive.ObjectID) (*bookstore.BookRating, error)
	UpdateUserRating(ctx context.Context, bookID, userID primitive.ObjectID, rating float64, comment string, tags []string) error
	DeleteUserRating(ctx context.Context, bookID, userID primitive.ObjectID) error

	// 评分批量操作
	BatchUpdateRatingTags(ctx context.Context, ratingIDs []primitive.ObjectID, tags []string) error
	BatchDeleteRatings(ctx context.Context, ratingIDs []primitive.ObjectID) error
	BatchDeleteRatingsByBookID(ctx context.Context, bookID primitive.ObjectID) error
	BatchDeleteRatingsByUserID(ctx context.Context, userID primitive.ObjectID) error

	// 评分搜索和过滤
	SearchRatings(ctx context.Context, keyword string, page, pageSize int) ([]*bookstore.BookRating, int64, error)
	GetRatingsWithComments(ctx context.Context, bookID primitive.ObjectID, page, pageSize int) ([]*bookstore.BookRating, int64, error)
	GetHighRatedComments(ctx context.Context, bookID primitive.ObjectID, minRating float64, page, pageSize int) ([]*bookstore.BookRating, int64, error)
}

// BookRatingServiceImpl 书籍评分服务实现
type BookRatingServiceImpl struct {
	ratingRepo   BookstoreRepo.BookRatingRepository
	cacheService CacheService
}

// NewBookRatingService 创建书籍评分服务实例
func NewBookRatingService(ratingRepo BookstoreRepo.BookRatingRepository, cacheService CacheService) BookRatingService {
	return &BookRatingServiceImpl{
		ratingRepo:   ratingRepo,
		cacheService: cacheService,
	}
}

// CreateRating 创建评分
func (s *BookRatingServiceImpl) CreateRating(ctx context.Context, rating *bookstore.BookRating) error {
	if rating == nil {
		return errors.New("rating cannot be nil")
	}

	// 验证必填字段
	if rating.BookID.IsZero() {
		return errors.New("book ID is required")
	}
	if rating.UserID.IsZero() {
		return errors.New("user ID is required")
	}
	if !rating.IsValidRating() {
		return errors.New("rating must be between 1 and 5")
	}

	// 检查用户是否已经评分
	existingRating, err := s.ratingRepo.GetByBookIDAndUserID(ctx, rating.BookID, rating.UserID)
	if err != nil {
		return fmt.Errorf("failed to check existing rating: %w", err)
	}
	if existingRating != nil {
		return errors.New("user has already rated this book")
	}

	// 创建评分
	if err := s.ratingRepo.Create(ctx, rating); err != nil {
		return fmt.Errorf("failed to create rating: %w", err)
	}

	// 清除相关缓存
	s.invalidateRelatedCache(ctx, rating)

	return nil
}

// GetRatingByID 根据ID获取评分
func (s *BookRatingServiceImpl) GetRatingByID(ctx context.Context, id primitive.ObjectID) (*bookstore.BookRating, error) {
	// 先尝试从缓存获取
	if s.cacheService != nil {
		if cachedRating, err := s.cacheService.GetBookRating(ctx, id.Hex()); err == nil && cachedRating != nil {
			return cachedRating, nil
		}
	}

	// 从数据库获取
	rating, err := s.ratingRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get rating: %w", err)
	}
	if rating == nil {
		return nil, errors.New("rating not found")
	}

	// 缓存结果
	if s.cacheService != nil {
		s.cacheService.SetBookRating(ctx, id.Hex(), rating, 30*time.Minute)
	}

	return rating, nil
}

// UpdateRating 更新评分
func (s *BookRatingServiceImpl) UpdateRating(ctx context.Context, rating *bookstore.BookRating) error {
	if rating == nil {
		return errors.New("rating cannot be nil")
	}

	// 验证必填字段
	if rating.BookID.IsZero() {
		return errors.New("book ID is required")
	}
	if rating.UserID.IsZero() {
		return errors.New("user ID is required")
	}
	if !rating.IsValidRating() {
		return errors.New("rating must be between 1 and 5")
	}

	// 更新评分
	updates := map[string]interface{}{
		"rating":     rating.Rating,
		"comment":    rating.Comment,
		"tags":       rating.Tags,
		"updated_at": rating.UpdatedAt,
	}
	if err := s.ratingRepo.Update(ctx, rating.ID, updates); err != nil {
		return fmt.Errorf("failed to update rating: %w", err)
	}

	// 清除相关缓存
	s.invalidateRelatedCache(ctx, rating)

	return nil
}

// DeleteRating 删除评分
func (s *BookRatingServiceImpl) DeleteRating(ctx context.Context, id primitive.ObjectID) error {
	// 先获取评分信息用于清除缓存
	rating, err := s.ratingRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get rating for deletion: %w", err)
	}
	if rating == nil {
		return errors.New("rating not found")
	}

	// 删除评分
	if err := s.ratingRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete rating: %w", err)
	}

	// 清除相关缓存
	s.invalidateRelatedCache(ctx, rating)

	return nil
}

// GetRatingsByBookID 根据书籍ID获取评分列表
func (s *BookRatingServiceImpl) GetRatingsByBookID(ctx context.Context, bookID primitive.ObjectID, page, pageSize int) ([]*bookstore.BookRating, int64, error) {
	if bookID.IsZero() {
		return nil, 0, errors.New("book ID cannot be empty")
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	// 获取评分列表
	ratings, err := s.ratingRepo.GetByBookID(ctx, bookID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get ratings by book ID: %w", err)
	}

	// 获取总数
	total, err := s.ratingRepo.CountByBookID(ctx, bookID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count ratings by book ID: %w", err)
	}

	return ratings, total, nil
}

// GetRatingsByUserID 根据用户ID获取评分列表
func (s *BookRatingServiceImpl) GetRatingsByUserID(ctx context.Context, userID primitive.ObjectID, page, pageSize int) ([]*bookstore.BookRating, int64, error) {
	if userID.IsZero() {
		return nil, 0, errors.New("user ID cannot be empty")
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	// 获取评分列表
	ratings, err := s.ratingRepo.GetByUserID(ctx, userID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get ratings by user ID: %w", err)
	}

	// 获取总数
	total, err := s.ratingRepo.CountByUserID(ctx, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count ratings by user ID: %w", err)
	}

	return ratings, total, nil
}

// GetRatingByBookIDAndUserID 根据书籍ID和用户ID获取评分
func (s *BookRatingServiceImpl) GetRatingByBookIDAndUserID(ctx context.Context, bookID, userID primitive.ObjectID) (*bookstore.BookRating, error) {
	if bookID.IsZero() {
		return nil, errors.New("book ID cannot be empty")
	}
	if userID.IsZero() {
		return nil, errors.New("user ID cannot be empty")
	}

	rating, err := s.ratingRepo.GetByBookIDAndUserID(ctx, bookID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get rating by book ID and user ID: %w", err)
	}

	return rating, nil
}

// GetRatingsByRating 根据评分值获取评分列表
func (s *BookRatingServiceImpl) GetRatingsByRating(ctx context.Context, rating float64, page, pageSize int) ([]*bookstore.BookRating, int64, error) {
	if rating < 1 || rating > 5 {
		return nil, 0, errors.New("rating must be between 1 and 5")
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	// 获取评分列表
	ratings, err := s.ratingRepo.GetByRating(ctx, int(rating), pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get ratings by rating value: %w", err)
	}

	// 对于获取所有特定评分的评分，我们需要计算总数
	// 由于CountByRating需要bookID参数，我们使用返回的ratings长度作为近似
	// 或者实现一个新的Count方法
	total := int64(len(ratings))
	if len(ratings) == pageSize {
		// 如果返回的数量等于pageSize，可能还有更多数据
		// 这里简化处理，实际应该有专门的CountByRatingValue方法
		total = int64(pageSize * page) // 估算值
	}

	return ratings, total, nil
}

// GetRatingsByTags 根据标签获取评分列表
func (s *BookRatingServiceImpl) GetRatingsByTags(ctx context.Context, tags []string, page, pageSize int) ([]*bookstore.BookRating, int64, error) {
	if len(tags) == 0 {
		return nil, 0, errors.New("tags cannot be empty")
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	// 获取评分列表
	ratings, err := s.ratingRepo.GetByTags(ctx, tags, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get ratings by tags: %w", err)
	}

	// 这里简化处理，实际应该有专门的计数方法
	total := int64(len(ratings))
	if len(ratings) == pageSize {
		total = int64((page + 1) * pageSize)
	}

	return ratings, total, nil
}

// GetAverageRating 获取书籍平均评分
func (s *BookRatingServiceImpl) GetAverageRating(ctx context.Context, bookID primitive.ObjectID) (float64, error) {
	if bookID.IsZero() {
		return 0, errors.New("book ID cannot be empty")
	}

	// 先尝试从缓存获取
	if s.cacheService != nil {
		if avgRating, err := s.cacheService.GetBookAverageRating(ctx, bookID.Hex()); err == nil && avgRating > 0 {
			return avgRating, nil
		}
	}

	// 从数据库获取
	avgRating, err := s.ratingRepo.GetAverageRating(ctx, bookID)
	if err != nil {
		return 0, fmt.Errorf("failed to get average rating: %w", err)
	}

	// 缓存结果
	if s.cacheService != nil {
		s.cacheService.SetBookAverageRating(ctx, bookID.Hex(), avgRating, 10*time.Minute)
	}

	return avgRating, nil
}

// GetRatingDistribution 获取评分分布
func (s *BookRatingServiceImpl) GetRatingDistribution(ctx context.Context, bookID primitive.ObjectID) (map[string]int64, error) {
	if bookID.IsZero() {
		return nil, errors.New("book ID cannot be empty")
	}

	distribution, err := s.ratingRepo.GetRatingDistribution(ctx, bookID)
	if err != nil {
		return nil, fmt.Errorf("failed to get rating distribution: %w", err)
	}

	return distribution, nil
}

// GetRatingCount 获取评分总数
func (s *BookRatingServiceImpl) GetRatingCount(ctx context.Context, bookID primitive.ObjectID) (int64, error) {
	if bookID.IsZero() {
		return 0, errors.New("book ID cannot be empty")
	}

	count, err := s.ratingRepo.CountByBookID(ctx, bookID)
	if err != nil {
		return 0, fmt.Errorf("failed to get rating count: %w", err)
	}

	return count, nil
}

// GetRatingStats 获取评分统计信息
func (s *BookRatingServiceImpl) GetRatingStats(ctx context.Context, bookID primitive.ObjectID) (map[string]interface{}, error) {
	if bookID.IsZero() {
		return nil, errors.New("book ID cannot be empty")
	}

	stats := make(map[string]interface{})

	// 平均评分
	avgRating, err := s.ratingRepo.GetAverageRating(ctx, bookID)
	if err != nil {
		return nil, fmt.Errorf("failed to get average rating: %w", err)
	}
	stats["average_rating"] = avgRating

	// 评分总数
	totalCount, err := s.ratingRepo.CountByBookID(ctx, bookID)
	if err != nil {
		return nil, fmt.Errorf("failed to get rating count: %w", err)
	}
	stats["total_ratings"] = totalCount

	// 评分分布
	distribution, err := s.ratingRepo.GetRatingDistribution(ctx, bookID)
	if err != nil {
		return nil, fmt.Errorf("failed to get rating distribution: %w", err)
	}
	stats["rating_distribution"] = distribution

	return stats, nil
}

// GetTopRatedBooks 获取高评分书籍
func (s *BookRatingServiceImpl) GetTopRatedBooks(ctx context.Context, limit int) ([]*bookstore.BookRating, error) {
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// 由于BookRatingRepository接口中没有GetTopRatedBooks方法，
	// 我们使用GetTopRated方法，但需要一个bookID参数
	// 这里暂时返回空结果，实际应该重新设计这个方法
	return []*bookstore.BookRating{}, nil
}

// LikeRating 点赞评分
func (s *BookRatingServiceImpl) LikeRating(ctx context.Context, ratingID primitive.ObjectID, userID primitive.ObjectID) error {
	if ratingID.IsZero() {
		return errors.New("rating ID cannot be empty")
	}
	if userID.IsZero() {
		return errors.New("user ID cannot be empty")
	}

	// TODO: 这里应该检查用户是否已经点赞过，避免重复点赞
	// 简化处理，直接增加点赞数
	if err := s.ratingRepo.IncrementLikes(ctx, ratingID); err != nil {
		return fmt.Errorf("failed to like rating: %w", err)
	}

	// 清除相关缓存
	if s.cacheService != nil {
		s.cacheService.InvalidateBookRatingCache(ctx, ratingID.Hex())
	}

	return nil
}

// UnlikeRating 取消点赞评分
func (s *BookRatingServiceImpl) UnlikeRating(ctx context.Context, ratingID primitive.ObjectID, userID primitive.ObjectID) error {
	if ratingID.IsZero() {
		return errors.New("rating ID cannot be empty")
	}
	if userID.IsZero() {
		return errors.New("user ID cannot be empty")
	}

	// TODO: 这里应该检查用户是否已经点赞过
	// 简化处理，直接减少点赞数
	if err := s.ratingRepo.DecrementLikes(ctx, ratingID); err != nil {
		return fmt.Errorf("failed to unlike rating: %w", err)
	}

	// 清除相关缓存
	if s.cacheService != nil {
		s.cacheService.InvalidateBookRatingCache(ctx, ratingID.Hex())
	}

	return nil
}

// GetRatingLikes 获取评分点赞数
func (s *BookRatingServiceImpl) GetRatingLikes(ctx context.Context, ratingID primitive.ObjectID) (int64, error) {
	if ratingID.IsZero() {
		return 0, errors.New("rating ID cannot be empty")
	}

	rating, err := s.ratingRepo.GetByID(ctx, ratingID)
	if err != nil {
		return 0, fmt.Errorf("failed to get rating: %w", err)
	}
	if rating == nil {
		return 0, errors.New("rating not found")
	}

	return int64(rating.Likes), nil
}

// HasUserRated 检查用户是否已评分
func (s *BookRatingServiceImpl) HasUserRated(ctx context.Context, bookID, userID primitive.ObjectID) (bool, error) {
	if bookID.IsZero() {
		return false, errors.New("book ID cannot be empty")
	}
	if userID.IsZero() {
		return false, errors.New("user ID cannot be empty")
	}

	rating, err := s.ratingRepo.GetByBookIDAndUserID(ctx, bookID, userID)
	if err != nil {
		return false, fmt.Errorf("failed to check user rating: %w", err)
	}

	return rating != nil, nil
}

// GetUserRatingForBook 获取用户对书籍的评分
func (s *BookRatingServiceImpl) GetUserRatingForBook(ctx context.Context, bookID, userID primitive.ObjectID) (*bookstore.BookRating, error) {
	if bookID.IsZero() {
		return nil, errors.New("book ID cannot be empty")
	}
	if userID.IsZero() {
		return nil, errors.New("user ID cannot be empty")
	}

	rating, err := s.ratingRepo.GetByBookIDAndUserID(ctx, bookID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user rating: %w", err)
	}

	return rating, nil
}

// UpdateUserRating 更新用户评分
func (s *BookRatingServiceImpl) UpdateUserRating(ctx context.Context, bookID, userID primitive.ObjectID, rating float64, comment string, tags []string) error {
	if bookID.IsZero() {
		return errors.New("book ID cannot be empty")
	}
	if userID.IsZero() {
		return errors.New("user ID cannot be empty")
	}
	if rating < 1 || rating > 5 {
		return errors.New("rating must be between 1 and 5")
	}

	// 获取现有评分
	existingRating, err := s.ratingRepo.GetByBookIDAndUserID(ctx, bookID, userID)
	if err != nil {
		return fmt.Errorf("failed to get existing rating: %w", err)
	}

	if existingRating == nil {
		// 创建新评分
		newRating := &bookstore.BookRating{
			BookID:  bookID,
			UserID:  userID,
			Rating:  int(rating),
			Comment: comment,
			Tags:    tags,
		}
		return s.CreateRating(ctx, newRating)
	} else {
		// 更新现有评分
		existingRating.Rating = int(rating)
		existingRating.Comment = comment
		existingRating.Tags = tags
		return s.UpdateRating(ctx, existingRating)
	}
}

// DeleteUserRating 删除用户评分
func (s *BookRatingServiceImpl) DeleteUserRating(ctx context.Context, bookID, userID primitive.ObjectID) error {
	if bookID.IsZero() {
		return errors.New("book ID cannot be empty")
	}
	if userID.IsZero() {
		return errors.New("user ID cannot be empty")
	}

	// 获取现有评分
	existingRating, err := s.ratingRepo.GetByBookIDAndUserID(ctx, bookID, userID)
	if err != nil {
		return fmt.Errorf("failed to get existing rating: %w", err)
	}
	if existingRating == nil {
		return errors.New("rating not found")
	}

	// 删除评分
	return s.DeleteRating(ctx, existingRating.ID)
}

// BatchUpdateRatingTags 批量更新评分标签
func (s *BookRatingServiceImpl) BatchUpdateRatingTags(ctx context.Context, ratingIDs []primitive.ObjectID, tags []string) error {
	if len(ratingIDs) == 0 {
		return errors.New("rating IDs cannot be empty")
	}
	if len(tags) == 0 {
		return errors.New("tags cannot be empty")
	}

	if err := s.ratingRepo.BatchUpdateTags(ctx, ratingIDs, tags); err != nil {
		return fmt.Errorf("failed to batch update rating tags: %w", err)
	}

	// 清除相关缓存
	for _, ratingID := range ratingIDs {
		if s.cacheService != nil {
			s.cacheService.InvalidateBookRatingCache(ctx, ratingID.Hex())
		}
	}

	return nil
}

// BatchDeleteRatings 批量删除评分
func (s *BookRatingServiceImpl) BatchDeleteRatings(ctx context.Context, ratingIDs []primitive.ObjectID) error {
	if len(ratingIDs) == 0 {
		return errors.New("rating IDs cannot be empty")
	}

	if err := s.ratingRepo.BatchDelete(ctx, ratingIDs); err != nil {
		return fmt.Errorf("failed to batch delete ratings: %w", err)
	}

	// 清除相关缓存
	for _, ratingID := range ratingIDs {
		if s.cacheService != nil {
			s.cacheService.InvalidateBookRatingCache(ctx, ratingID.Hex())
		}
	}

	return nil
}

// BatchDeleteRatingsByBookID 批量删除书籍的所有评分
func (s *BookRatingServiceImpl) BatchDeleteRatingsByBookID(ctx context.Context, bookID primitive.ObjectID) error {
	if bookID.IsZero() {
		return errors.New("book ID cannot be empty")
	}

	if err := s.ratingRepo.DeleteByBookID(ctx, bookID); err != nil {
		return fmt.Errorf("failed to batch delete ratings by book ID: %w", err)
	}

	// 清除相关缓存
	if s.cacheService != nil {
		s.cacheService.InvalidateBookRatingsCache(ctx, bookID.Hex())
	}

	return nil
}

// BatchDeleteRatingsByUserID 批量删除用户的所有评分
func (s *BookRatingServiceImpl) BatchDeleteRatingsByUserID(ctx context.Context, userID primitive.ObjectID) error {
	if userID.IsZero() {
		return errors.New("user ID cannot be empty")
	}

	if err := s.ratingRepo.DeleteByUserID(ctx, userID); err != nil {
		return fmt.Errorf("failed to batch delete ratings by user ID: %w", err)
	}

	// 清除相关缓存
	if s.cacheService != nil {
		s.cacheService.InvalidateUserRatingsCache(ctx, userID.Hex())
	}

	return nil
}

// SearchRatings 搜索评分
func (s *BookRatingServiceImpl) SearchRatings(ctx context.Context, keyword string, page, pageSize int) ([]*bookstore.BookRating, int64, error) {
	if keyword == "" {
		return nil, 0, errors.New("keyword cannot be empty")
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	// 搜索评分
	ratings, err := s.ratingRepo.Search(ctx, keyword, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search ratings: %w", err)
	}

	// 这里简化处理，实际应该有专门的搜索计数方法
	total := int64(len(ratings))
	if len(ratings) == pageSize {
		total = int64((page + 1) * pageSize)
	}

	return ratings, total, nil
}

// GetRatingsWithComments 获取有评论的评分
func (s *BookRatingServiceImpl) GetRatingsWithComments(ctx context.Context, bookID primitive.ObjectID, page, pageSize int) ([]*bookstore.BookRating, int64, error) {
	if bookID.IsZero() {
		return nil, 0, errors.New("book ID cannot be empty")
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	// 获取所有评分，然后过滤有评论的
	allRatings, err := s.ratingRepo.GetByBookID(ctx, bookID, pageSize*2, offset) // 获取更多数据以便过滤
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get ratings: %w", err)
	}

	// 过滤有评论的评分
	var ratingsWithComments []*bookstore.BookRating
	for _, rating := range allRatings {
		if rating.HasComment() {
			ratingsWithComments = append(ratingsWithComments, rating)
			if len(ratingsWithComments) >= pageSize {
				break
			}
		}
	}

	// 简化处理总数计算
	total := int64(len(ratingsWithComments))
	if len(ratingsWithComments) == pageSize {
		total = int64((page + 1) * pageSize)
	}

	return ratingsWithComments, total, nil
}

// GetHighRatedComments 获取高评分评论
func (s *BookRatingServiceImpl) GetHighRatedComments(ctx context.Context, bookID primitive.ObjectID, minRating float64, page, pageSize int) ([]*bookstore.BookRating, int64, error) {
	if bookID.IsZero() {
		return nil, 0, errors.New("book ID cannot be empty")
	}
	if minRating < 1 || minRating > 5 {
		return nil, 0, errors.New("minimum rating must be between 1 and 5")
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	// 获取所有评分，然后过滤高评分且有评论的
	allRatings, err := s.ratingRepo.GetByBookID(ctx, bookID, pageSize*2, offset) // 获取更多数据以便过滤
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get ratings: %w", err)
	}

	// 过滤高评分且有评论的评分
	var highRatedComments []*bookstore.BookRating
	for _, rating := range allRatings {
		if rating.Rating >= int(minRating) && rating.HasComment() {
			highRatedComments = append(highRatedComments, rating)
			if len(highRatedComments) >= pageSize {
				break
			}
		}
	}

	// 简化处理总数计算
	total := int64(len(highRatedComments))
	if len(highRatedComments) == pageSize {
		total = int64((page + 1) * pageSize)
	}

	return highRatedComments, total, nil
}

// invalidateRelatedCache 清除相关缓存
func (s *BookRatingServiceImpl) invalidateRelatedCache(ctx context.Context, rating *bookstore.BookRating) {
	if s.cacheService == nil {
		return
	}

	// 清除评分缓存
	s.cacheService.InvalidateBookRatingCache(ctx, rating.ID.Hex())

	// 清除书籍评分相关缓存
	s.cacheService.InvalidateBookRatingsCache(ctx, rating.BookID.Hex())
	s.cacheService.InvalidateBookAverageRatingCache(ctx, rating.BookID.Hex())

	// 清除用户评分缓存
	s.cacheService.InvalidateUserRatingsCache(ctx, rating.UserID.Hex())

	// 清除书籍详情缓存（因为平均评分可能变化）
	s.cacheService.InvalidateBookDetailCache(ctx, rating.BookID.Hex())
}
