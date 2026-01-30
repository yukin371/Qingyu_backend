package mongodb

import (
	bookstore2 "Qingyu_backend/models/bookstore"
	"context"

	BookstoreInterface "Qingyu_backend/repository/interfaces/bookstore"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoBookDataMutationRepository MongoDB书籍数据变更仓储实现
// 通过组合MongoBookRepository实现BookDataMutationRepository接口
type MongoBookDataMutationRepository struct {
	*MongoBookRepository
}

// NewMongoBookDataMutationRepository 创建MongoDB书籍数据变更仓储实例
func NewMongoBookDataMutationRepository(client *mongo.Client, database string) BookstoreInterface.BookDataMutationRepository {
	baseRepo := NewMongoBookRepository(client, database)
	return &MongoBookDataMutationRepository{
		MongoBookRepository: baseRepo.(*MongoBookRepository),
	}
}

// 确保实现了接口
var _ BookstoreInterface.BookDataMutationRepository = (*MongoBookDataMutationRepository)(nil)

// Create 创建书籍
func (r *MongoBookDataMutationRepository) Create(ctx context.Context, book *bookstore2.Book) error {
	return r.MongoBookRepository.Create(ctx, book)
}

// Update 更新书籍
func (r *MongoBookDataMutationRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	return r.MongoBookRepository.Update(ctx, id, updates)
}

// Delete 删除书籍
func (r *MongoBookDataMutationRepository) Delete(ctx context.Context, id string) error {
	return r.MongoBookRepository.Delete(ctx, id)
}

// Transaction 执行事务
func (r *MongoBookDataMutationRepository) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return r.MongoBookRepository.Transaction(ctx, fn)
}
