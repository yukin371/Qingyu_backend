package mongodb

import (
	bookstore2 "Qingyu_backend/models/bookstore"
	"context"

	BookstoreInterface "Qingyu_backend/repository/interfaces/bookstore"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoBookSearchRepository MongoDB书籍搜索仓储实现
// 通过组合MongoBookRepository实现BookSearchRepository接口
type MongoBookSearchRepository struct {
	*MongoBookRepository
}

// NewMongoBookSearchRepository 创建MongoDB书籍搜索仓储实例
func NewMongoBookSearchRepository(client *mongo.Client, database string) BookstoreInterface.BookSearchRepository {
	baseRepo := NewMongoBookRepository(client, database)
	return &MongoBookSearchRepository{
		MongoBookRepository: baseRepo.(*MongoBookRepository),
	}
}

// 确保实现了接口
var _ BookstoreInterface.BookSearchRepository = (*MongoBookSearchRepository)(nil)

// Search 搜索书籍
func (r *MongoBookSearchRepository) Search(ctx context.Context, keyword string, limit, offset int) ([]*bookstore2.Book, error) {
	return r.MongoBookRepository.Search(ctx, keyword, limit, offset)
}

// SearchWithFilter 使用过滤器搜索书籍
func (r *MongoBookSearchRepository) SearchWithFilter(ctx context.Context, filter *bookstore2.BookFilter) ([]*bookstore2.Book, error) {
	return r.MongoBookRepository.SearchWithFilter(ctx, filter)
}

// GetByPriceRange 按价格区间获取书籍
func (r *MongoBookSearchRepository) GetByPriceRange(ctx context.Context, minPrice, maxPrice float64, limit, offset int) ([]*bookstore2.Book, error) {
	return r.MongoBookRepository.GetByPriceRange(ctx, minPrice, maxPrice, limit, offset)
}
