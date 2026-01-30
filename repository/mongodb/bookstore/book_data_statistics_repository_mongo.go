package mongodb

import (
	bookstore2 "Qingyu_backend/models/bookstore"
	"context"

	BookstoreInterface "Qingyu_backend/repository/interfaces/bookstore"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoBookDataStatisticsRepository MongoDB书籍数据统计仓储实现
// 通过组合MongoBookRepository实现BookDataStatisticsRepository接口
type MongoBookDataStatisticsRepository struct {
	*MongoBookRepository
}

// NewMongoBookDataStatisticsRepository 创建MongoDB书籍数据统计仓储实例
func NewMongoBookDataStatisticsRepository(client *mongo.Client, database string) BookstoreInterface.BookDataStatisticsRepository {
	baseRepo := NewMongoBookRepository(client, database)
	return &MongoBookDataStatisticsRepository{
		MongoBookRepository: baseRepo.(*MongoBookRepository),
	}
}

// 确保实现了接口
var _ BookstoreInterface.BookDataStatisticsRepository = (*MongoBookDataStatisticsRepository)(nil)

// CountByCategory 统计分类下的书籍数量
func (r *MongoBookDataStatisticsRepository) CountByCategory(ctx context.Context, categoryID string) (int64, error) {
	return r.MongoBookRepository.CountByCategory(ctx, categoryID)
}

// CountByAuthor 统计作者的书籍数量
func (r *MongoBookDataStatisticsRepository) CountByAuthor(ctx context.Context, author string) (int64, error) {
	return r.MongoBookRepository.CountByAuthor(ctx, author)
}

// CountByStatus 统计指定状态的书籍数量
func (r *MongoBookDataStatisticsRepository) CountByStatus(ctx context.Context, status bookstore2.BookStatus) (int64, error) {
	return r.MongoBookRepository.CountByStatus(ctx, status)
}

// CountByFilter 根据过滤器统计书籍数量
func (r *MongoBookDataStatisticsRepository) CountByFilter(ctx context.Context, filter *bookstore2.BookFilter) (int64, error) {
	return r.MongoBookRepository.CountByFilter(ctx, filter)
}

// GetStats 获取书籍统计概览信息
func (r *MongoBookDataStatisticsRepository) GetStats(ctx context.Context) (*bookstore2.BookStats, error) {
	return r.MongoBookRepository.GetStats(ctx)
}

// IncrementViewCount 增加浏览计数
func (r *MongoBookDataStatisticsRepository) IncrementViewCount(ctx context.Context, bookID string) error {
	return r.MongoBookRepository.IncrementViewCount(ctx, bookID)
}

// IncrementLikeCount 增加点赞数
func (r *MongoBookDataStatisticsRepository) IncrementLikeCount(ctx context.Context, bookID string) error {
	return r.MongoBookRepository.IncrementLikeCount(ctx, bookID)
}

// IncrementCommentCount 增加评论数
func (r *MongoBookDataStatisticsRepository) IncrementCommentCount(ctx context.Context, bookID string) error {
	return r.MongoBookRepository.IncrementCommentCount(ctx, bookID)
}

// UpdateRating 更新评分
func (r *MongoBookDataStatisticsRepository) UpdateRating(ctx context.Context, bookID string, rating float64) error {
	return r.MongoBookRepository.UpdateRating(ctx, bookID, rating)
}
