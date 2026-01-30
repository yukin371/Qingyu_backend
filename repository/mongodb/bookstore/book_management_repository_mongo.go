package mongodb

import (
	bookstore2 "Qingyu_backend/models/bookstore"
	"context"

	BookstoreInterface "Qingyu_backend/repository/interfaces/bookstore"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoBookManagementRepository MongoDB书籍管理仓储实现
// 通过组合MongoBookRepository实现BookManagementRepository接口
type MongoBookManagementRepository struct {
	*MongoBookRepository
}

// NewMongoBookManagementRepository 创建MongoDB书籍管理仓储实例
func NewMongoBookManagementRepository(client *mongo.Client, database string) BookstoreInterface.BookManagementRepository {
	baseRepo := NewMongoBookRepository(client, database)
	return &MongoBookManagementRepository{
		MongoBookRepository: baseRepo.(*MongoBookRepository),
	}
}

// 确保实现了接口
var _ BookstoreInterface.BookManagementRepository = (*MongoBookManagementRepository)(nil)

// BatchUpdateStatus 批量更新书籍状态
func (r *MongoBookManagementRepository) BatchUpdateStatus(ctx context.Context, bookIDs []string, status bookstore2.BookStatus) error {
	return r.MongoBookRepository.BatchUpdateStatus(ctx, bookIDs, status)
}

// BatchUpdateCategory 批量更新书籍分类
func (r *MongoBookManagementRepository) BatchUpdateCategory(ctx context.Context, bookIDs []string, categoryIDs []string) error {
	return r.MongoBookRepository.BatchUpdateCategory(ctx, bookIDs, categoryIDs)
}

// BatchUpdateRecommended 批量更新推荐状态
func (r *MongoBookManagementRepository) BatchUpdateRecommended(ctx context.Context, bookIDs []string, isRecommended bool) error {
	return r.MongoBookRepository.BatchUpdateRecommended(ctx, bookIDs, isRecommended)
}

// BatchUpdateFeatured 批量更新精选状态
func (r *MongoBookManagementRepository) BatchUpdateFeatured(ctx context.Context, bookIDs []string, isFeatured bool) error {
	return r.MongoBookRepository.BatchUpdateFeatured(ctx, bookIDs, isFeatured)
}

// GetYears 获取所有书籍的发布年份列表
func (r *MongoBookManagementRepository) GetYears(ctx context.Context) ([]int, error) {
	return r.MongoBookRepository.GetYears(ctx)
}

// GetTags 获取所有标签列表
func (r *MongoBookManagementRepository) GetTags(ctx context.Context, categoryID *string) ([]string, error) {
	return r.MongoBookRepository.GetTags(ctx, categoryID)
}
