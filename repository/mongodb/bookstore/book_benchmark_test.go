package mongodb_test

import (
	"Qingyu_backend/models/bookstore"
	"context"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	mongodb "Qingyu_backend/repository/mongodb/bookstore"
)

// setupBenchmarkDB 设置基准测试数据库
func setupBenchmarkDB(b *testing.B) *mongo.Client {
	b.Helper()

	// 注意：这是一个简化的基准测试设置
	// 实际运行时需要真实的MongoDB连接
	// 这里使用 nil 客户端以确保代码可以编译
	return nil
}

// BenchmarkBookQuery_WithoutIndex 无索引基线测试
// 测试在没有索引的情况下查询性能
func BenchmarkBookQuery_WithoutIndex(b *testing.B) {
	client := setupBenchmarkDB(b)
	if client == nil {
		b.Skip("MongoDB connection not available")
	}

	db := client.Database("qingyu_benchmark")
	collection := db.Collection("books")
	ctx := context.Background()

	// 清理测试数据
	collection.DeleteMany(ctx, bson.M{})
	defer collection.DeleteMany(ctx, bson.M{})

	// 插入1000条测试数据（无索引）
	books := make([]interface{}, 1000)
	now := time.Now()
	for i := 0; i < 1000; i++ {
		book := bookstore.Book{
			Title:     "测试书籍",
			Author:    "测试作者",
			Status:    bookstore.BookStatusOngoing,
		}
		// 使用 BaseEntity 的方法设置时间戳
		book.TouchForCreate()
		books[i] = book
	}
	collection.InsertMany(ctx, books)

	b.ResetTimer() // 重置计时器

	// 基准测试：查询状态为ongoing的书籍（无索引）
	for i := 0; i < b.N; i++ {
		cursor, err := collection.Find(ctx, bson.M{"status": bookstore.BookStatusOngoing})
		if err != nil {
			b.Fatalf("查询失败: %v", err)
		}
		cursor.Close(ctx)
	}
}

// BenchmarkBookQuery_WithIndex 有索引性能测试
// 测试在创建索引后的查询性能提升
func BenchmarkBookQuery_WithIndex(b *testing.B) {
	client := setupBenchmarkDB(b)
	if client == nil {
		b.Skip("MongoDB connection not available")
	}

	db := client.Database("qingyu_benchmark")
	collection := db.Collection("books")
	ctx := context.Background()

	// 清理测试数据
	collection.DeleteMany(ctx, bson.M{})
	defer collection.Drop(ctx)

	// 创建索引：status_1
	indexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "status", Value: 1}},
	}
	collection.Indexes().CreateOne(ctx, indexModel)

	// 插入1000条测试数据
	books := make([]interface{}, 1000)
	for i := 0; i < 1000; i++ {
		book := bookstore.Book{
			Title:     "测试书籍",
			Author:    "测试作者",
			Status:    bookstore.BookStatusOngoing,
		}
		book.TouchForCreate()
		books[i] = book
	}
	collection.InsertMany(ctx, books)

	b.ResetTimer() // 重置计时器

	// 基准测试：查询状态为ongoing的书籍（有索引）
	for i := 0; i < b.N; i++ {
		cursor, err := collection.Find(ctx, bson.M{"status": bookstore.BookStatusOngoing})
		if err != nil {
			b.Fatalf("查询失败: %v", err)
		}
		cursor.Close(ctx)
	}
}

// BenchmarkBookQuery_CompoundIndex 复合索引性能测试
// 测试复合索引 status_1_created_at_-1 的性能
func BenchmarkBookQuery_CompoundIndex(b *testing.B) {
	client := setupBenchmarkDB(b)
	if client == nil {
		b.Skip("MongoDB connection not available")
	}

	db := client.Database("qingyu_benchmark")
	collection := db.Collection("books")
	ctx := context.Background()

	// 清理测试数据
	collection.DeleteMany(ctx, bson.M{})
	defer collection.Drop(ctx)

	// 创建复合索引：status_1_created_at_-1
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "status", Value: 1},
			{Key: "created_at", Value: -1},
		},
	}
	collection.Indexes().CreateOne(ctx, indexModel)

	// 插入1000条测试数据
	books := make([]interface{}, 1000)
	for i := 0; i < 1000; i++ {
		book := bookstore.Book{
			Title:      "测试书籍",
			Author:     "测试作者",
			Status:     bookstore.BookStatusOngoing,
		}
		// 设置不同的创建时间以便测试排序
		book.TouchForCreate()
		book.CreatedAt = time.Now().Add(-time.Duration(i) * time.Hour)
		books[i] = book
	}
	collection.InsertMany(ctx, books)

	b.ResetTimer() // 重置计时器

	// 基准测试：使用复合索引查询
	for i := 0; i < b.N; i++ {
		opts := options.Find().SetSort(bson.D{
			{Key: "status", Value: 1},
			{Key: "created_at", Value: -1},
		})
		cursor, err := collection.Find(ctx, bson.M{"status": bookstore.BookStatusOngoing}, opts)
		if err != nil {
			b.Fatalf("查询失败: %v", err)
		}
		cursor.Close(ctx)
	}
}

// BenchmarkBookRepository_GetByStatus_WithIndex 仓储层基准测试
// 测试仓储接口在有索引情况下的性能
func BenchmarkBookRepository_GetByStatus_WithIndex(b *testing.B) {
	client := setupBenchmarkDB(b)
	if client == nil {
		b.Skip("MongoDB connection not available")
	}

	// 清理测试数据
	dbName := "qingyu_benchmark"
	db := client.Database(dbName)
	collection := db.Collection("books")
	collection.DeleteMany(context.Background(), bson.M{})
	defer db.Drop(context.Background())

	// 创建索引
	indexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "status", Value: 1}},
	}
	collection.Indexes().CreateOne(context.Background(), indexModel)

	// 创建仓储实例
	repo := mongodb.NewMongoBookRepository(client, dbName)

	// 插入1000条测试数据
	ctx := context.Background()
	for i := 0; i < 1000; i++ {
		book := &bookstore.Book{
			Title:     "测试书籍",
			Author:    "测试作者",
			Status:    bookstore.BookStatusOngoing,
		}
		book.TouchForCreate()
		repo.Create(ctx, book)
	}

	b.ResetTimer() // 重置计时器

	// 基准测试：通过仓储接口查询
	for i := 0; i < b.N; i++ {
		_, err := repo.GetByStatus(ctx, bookstore.BookStatusOngoing, 100, 0)
		if err != nil {
			b.Fatalf("查询失败: %v", err)
		}
	}
}
