package engine

import (
	"context"
	"fmt"
)

// MongoEngine MongoDB 搜索引擎（兼容过渡）
type MongoEngine struct {
	// TODO: 注入 MongoDB client
}

// NewMongoEngine 创建 MongoDB 引擎
func NewMongoEngine() (*MongoEngine, error) {
	// TODO: 初始化 MongoDB client
	return &MongoEngine{}, nil
}

// Search 执行 MongoDB 搜索
func (m *MongoEngine) Search(ctx context.Context, index string, query interface{}, opts *SearchOptions) (*SearchResult, error) {
	// TODO: 实现 MongoDB 搜索
	return nil, fmt.Errorf("mongodb search not implemented")
}

// Index 批量索引文档（MongoDB 不需要，返回空实现）
func (m *MongoEngine) Index(ctx context.Context, index string, documents []Document) error {
	// MongoDB 作为数据源，不需要此方法
	return nil
}

// Update 更新文档（MongoDB 不需要，返回空实现）
func (m *MongoEngine) Update(ctx context.Context, index string, id string, document Document) error {
	// MongoDB 作为数据源，不需要此方法
	return nil
}

// Delete 删除文档（MongoDB 不需要，返回空实现）
func (m *MongoEngine) Delete(ctx context.Context, index string, id string) error {
	// MongoDB 作为数据源，不需要此方法
	return nil
}

// CreateIndex 创建索引（MongoDB 不需要，返回空实现）
func (m *MongoEngine) CreateIndex(ctx context.Context, index string, mapping interface{}) error {
	// MongoDB 作为数据源，不需要此方法
	return nil
}

// Health 健康检查
func (m *MongoEngine) Health(ctx context.Context) error {
	// TODO: 实现 MongoDB 健康检查
	return fmt.Errorf("mongodb health check not implemented")
}
