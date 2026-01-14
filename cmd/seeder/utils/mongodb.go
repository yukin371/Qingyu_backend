// Package utils 提供 MongoDB 工具函数
package utils

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Database MongoDB 数据库包装器
type Database struct {
	*mongo.Database
}

// Collection 获取集合
func (d *Database) Collection(name string) *mongo.Collection {
	return d.Database.Collection(name)
}

// Disconnect 断开数据库连接
func (d *Database) Disconnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return d.Client().Disconnect(ctx)
}

// ConnectDB 连接 MongoDB 数据库
func ConnectDB(uri, database string) (*Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("创建 MongoDB 客户端失败: %w", err)
	}

	// 测试连接
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("连接 MongoDB 失败: %w", err)
	}

	db := client.Database(database)
	return &Database{db}, nil
}

// BulkInserter MongoDB 批量插入器
type BulkInserter struct {
	collection *mongo.Collection
	batchSize  int
}

// NewBulkInserter 创建批量插入器
func NewBulkInserter(collection *mongo.Collection, batchSize int) *BulkInserter {
	return &BulkInserter{
		collection: collection,
		batchSize:  batchSize,
	}
}

// InsertMany 批量插入文档
func (b *BulkInserter) InsertMany(ctx context.Context, documents interface{}) error {
	// 将切片转换为 []interface{}
	docsSlice := toInterfaceSlice(documents)

	for i := 0; i < len(docsSlice); i += b.batchSize {
		end := i + b.batchSize
		if end > len(docsSlice) {
			end = len(docsSlice)
		}

		_, err := b.collection.InsertMany(ctx, docsSlice[i:end])
		if err != nil {
			return fmt.Errorf("批量插入失败，批次 %d: %w", i/b.batchSize, err)
		}
	}

	return nil
}

// toInterfaceSlice 将任意切片转换为 []interface{}
func toInterfaceSlice(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		return nil
	}

	ret := make([]interface{}, s.Len())
	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}
	return ret
}

// CountDocuments 统计文档数量
func (b *BulkInserter) CountDocuments(ctx context.Context, filter interface{}) (int64, error) {
	return b.collection.CountDocuments(ctx, filter)
}

// DeleteMany 批量删除文档
func (b *BulkInserter) DeleteMany(ctx context.Context, filter interface{}) (int64, error) {
	result, err := b.collection.DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}
