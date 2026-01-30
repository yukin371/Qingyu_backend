package mongodb

import (
	bookstore2 "Qingyu_backend/models/bookstore"
	base "Qingyu_backend/repository/interfaces/infrastructure"
	"context"

	BookstoreInterface "Qingyu_backend/repository/interfaces/bookstore"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoBookListQueryRepository MongoDB书籍列表查询仓储实现
// 通过组合MongoBookRepository实现BookListQueryRepository接口
type MongoBookListQueryRepository struct {
	*MongoBookRepository
}

// NewMongoBookListQueryRepository 创建MongoDB书籍列表查询仓储实例
func NewMongoBookListQueryRepository(client *mongo.Client, database string) BookstoreInterface.BookListQueryRepository {
	baseRepo := NewMongoBookRepository(client, database)
	return &MongoBookListQueryRepository{
		MongoBookRepository: baseRepo.(*MongoBookRepository),
	}
}

// 确保实现了接口
var _ BookstoreInterface.BookListQueryRepository = (*MongoBookListQueryRepository)(nil)

// 以下方法直接委托给MongoBookRepository

// GetByID 根据ID获取书籍
func (r *MongoBookListQueryRepository) GetByID(ctx context.Context, id string) (*bookstore2.Book, error) {
	return r.MongoBookRepository.GetByID(ctx, id)
}

// List 列出书籍
func (r *MongoBookListQueryRepository) List(ctx context.Context, filter base.Filter) ([]*bookstore2.Book, error) {
	return r.MongoBookRepository.List(ctx, filter)
}

// Count 统计书籍数量
func (r *MongoBookListQueryRepository) Count(ctx context.Context, filter base.Filter) (int64, error) {
	return r.MongoBookRepository.Count(ctx, filter)
}

// Exists 判断书籍是否存在
func (r *MongoBookListQueryRepository) Exists(ctx context.Context, id string) (bool, error) {
	return r.MongoBookRepository.Exists(ctx, id)
}

// Health 健康检查
func (r *MongoBookListQueryRepository) Health(ctx context.Context) error {
	return r.MongoBookRepository.Health(ctx)
}

// GetByCategory 根据分类获取书籍列表
func (r *MongoBookListQueryRepository) GetByCategory(ctx context.Context, categoryID string, limit, offset int) ([]*bookstore2.Book, error) {
	return r.MongoBookRepository.GetByCategory(ctx, categoryID, limit, offset)
}

// GetByAuthor 根据作者获取书籍列表
func (r *MongoBookListQueryRepository) GetByAuthor(ctx context.Context, author string, limit, offset int) ([]*bookstore2.Book, error) {
	return r.MongoBookRepository.GetByAuthor(ctx, author, limit, offset)
}

// GetByAuthorID 根据作者ID获取书籍列表
func (r *MongoBookListQueryRepository) GetByAuthorID(ctx context.Context, authorID string, limit, offset int) ([]*bookstore2.Book, error) {
	return r.MongoBookRepository.GetByAuthorID(ctx, authorID, limit, offset)
}

// GetByStatus 根据状态获取书籍列表
func (r *MongoBookListQueryRepository) GetByStatus(ctx context.Context, status bookstore2.BookStatus, limit, offset int) ([]*bookstore2.Book, error) {
	return r.MongoBookRepository.GetByStatus(ctx, status, limit, offset)
}

// GetRecommended 获取推荐书籍
func (r *MongoBookListQueryRepository) GetRecommended(ctx context.Context, limit, offset int) ([]*bookstore2.Book, error) {
	return r.MongoBookRepository.GetRecommended(ctx, limit, offset)
}

// GetFeatured 获取精选书籍
func (r *MongoBookListQueryRepository) GetFeatured(ctx context.Context, limit, offset int) ([]*bookstore2.Book, error) {
	return r.MongoBookRepository.GetFeatured(ctx, limit, offset)
}

// GetHotBooks 获取热门书籍
func (r *MongoBookListQueryRepository) GetHotBooks(ctx context.Context, limit, offset int) ([]*bookstore2.Book, error) {
	return r.MongoBookRepository.GetHotBooks(ctx, limit, offset)
}

// GetNewReleases 获取新上架书籍
func (r *MongoBookListQueryRepository) GetNewReleases(ctx context.Context, limit, offset int) ([]*bookstore2.Book, error) {
	return r.MongoBookRepository.GetNewReleases(ctx, limit, offset)
}

// GetFreeBooks 获取免费书籍
func (r *MongoBookListQueryRepository) GetFreeBooks(ctx context.Context, limit, offset int) ([]*bookstore2.Book, error) {
	return r.MongoBookRepository.GetFreeBooks(ctx, limit, offset)
}
