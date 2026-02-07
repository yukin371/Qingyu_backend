package mongodb

import (
	"context"
	"fmt"

	"Qingyu_backend/models/bookstore"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// PerformanceOptimizer 性能优化器
type PerformanceOptimizer struct {
	collection *mongo.Collection
}

// NewPerformanceOptimizer 创建性能优化器
func NewPerformanceOptimizer(collection *mongo.Collection) *PerformanceOptimizer {
	return &PerformanceOptimizer{
		collection: collection,
	}
}

// ============================================
// Index Optimization
// ============================================

// EnsureIndexes 确保必要的索引存在
func (p *PerformanceOptimizer) EnsureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		// 文本搜索索引
		{
			Keys: bson.D{
				{Key: "title", Value: "text"},
				{Key: "author", Value: "text"},
				{Key: "introduction", Value: "text"},
				{Key: "tags", Value: "text"},
			},
			Options: options.Index().SetName("text_search_index"),
		},
		// 分类索引
		{
			Keys: bson.D{{Key: "category_ids", Value: 1}},
			Options: options.Index().SetName("category_index"),
		},
		// 作者索引
		{
			Keys: bson.D{{Key: "author", Value: 1}},
			Options: options.Index().SetName("author_index"),
		},
		// 状态索引
		{
			Keys: bson.D{{Key: "status", Value: 1}},
			Options: options.Index().SetName("status_index"),
		},
		// 创建时间索引（用于时间戳游标）
		{
			Keys: bson.D{{Key: "created_at", Value: -1}},
			Options: options.Index().SetName("created_at_index"),
		},
		// 更新时间索引
		{
			Keys: bson.D{{Key: "updated_at", Value: -1}},
			Options: options.Index().SetName("updated_at_index"),
		},
		// 复合索引：状态 + 创建时间（常用查询组合）
		{
			Keys: bson.D{
				{Key: "status", Value: 1},
				{Key: "created_at", Value: -1},
			},
			Options: options.Index().SetName("status_created_at_index"),
		},
		// 复合索引：分类 + 状态 + 创建时间
		{
			Keys: bson.D{
				{Key: "category_ids", Value: 1},
				{Key: "status", Value: 1},
				{Key: "created_at", Value: -1},
			},
			Options: options.Index().SetName("category_status_created_at_index"),
		},
		// 推荐书籍索引
		{
			Keys: bson.D{{Key: "is_recommended", Value: 1}},
			Options: options.Index().SetName("recommended_index"),
		},
		// 热门书籍索引
		{
			Keys: bson.D{{Key: "is_hot", Value: 1}},
			Options: options.Index().SetName("hot_index"),
		},
		// 精选书籍索引
		{
			Keys: bson.D{{Key: "is_featured", Value: 1}},
			Options: options.Index().SetName("featured_index"),
		},
		// 免费书籍索引
		{
			Keys: bson.D{{Key: "is_free", Value: 1}},
			Options: options.Index().SetName("free_index"),
		},
		// 标签索引
		{
			Keys: bson.D{{Key: "tags", Value: 1}},
			Options: options.Index().SetName("tags_index"),
		},
	}

	_, err := p.collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return fmt.Errorf("create indexes failed: %w", err)
	}

	return nil
}

// ============================================
// Projection Optimization
// ============================================

// GetProjectionForList 获取列表视图的投影（减少数据传输）
func GetProjectionForList() bson.M {
	return bson.M{
		"title":          1,
		"author":         1,
		"cover":          1,
		"introduction":   1,
		"rating":         1,
		"view_count":     1,
		"word_count":     1,
		"status":         1,
		"created_at":     1,
		"updated_at":     1,
		"tags":           1,
		"is_hot":         1,
		"is_free":        1,
		"is_recommended": 1,
		"is_featured":    1,
	}
}

// GetProjectionForDetail 获取详情视图的投影
func GetProjectionForDetail() bson.M {
	return bson.M{
		"chapter_count": 1,
		"last_chapter":  1,
		"category_ids":  1,
		"price":         1,
	}
}

// ============================================
// Query Optimization
// ============================================

// OptimizeFindOptions 优化查询选项
func OptimizeFindOptions(filter *bookstore.BookFilter) *options.FindOptions {
	opts := options.Find()

	// 设置投影以减少数据传输
	if filter.Limit > 0 && filter.Limit <= 100 {
		// 对于列表查询，只返回必要字段
		opts.SetProjection(GetProjectionForList())
	}

	// 设置批量大小
	batchSize := int32(20)
	if filter.Limit > 0 && filter.Limit < 20 {
		batchSize = int32(filter.Limit)
	}
	opts.SetBatchSize(batchSize)

	// 设置排序
	sortField := filter.SortBy
	if sortField == "" {
		sortField = "created_at"
	}

	sortOrder := -1 // 默认降序
	if filter.SortOrder == "asc" {
		sortOrder = 1
	}

	opts.SetSort(bson.D{{Key: sortField, Value: sortOrder}})

	// 设置限制
	if filter.Limit > 0 {
		opts.SetLimit(int64(filter.Limit))
	} else {
		opts.SetLimit(20) // 默认限制
	}

	// 设置偏移（仅用于offset游标类型）
	if filter.Offset > 0 {
		opts.SetSkip(int64(filter.Offset))
	}

	return opts
}

// ============================================
// Aggregation Optimization
// ============================================

// BuildOptimizedAggregationPipeline 构建优化的聚合管道
func (p *PerformanceOptimizer) BuildOptimizedAggregationPipeline(
	filter *bookstore.BookFilter,
	requiresCount bool,
) []bson.D {
	pipeline := []bson.D{}

	// 1. Match阶段（尽早过滤数据）
	matchStage := p.buildMatchStage(filter)
	if len(matchStage) > 0 {
		pipeline = append(pipeline, matchStage)
	}

	// 2. Sort阶段（在Limit之前排序）
	sortStage := p.buildSortStage(filter)
	if len(sortStage) > 0 {
		pipeline = append(pipeline, sortStage)
	}

	// 3. Limit阶段（尽早限制结果集）
	if filter.Limit > 0 {
		limitStage := bson.D{{Key: "$limit", Value: filter.Limit}}
		pipeline = append(pipeline, limitStage)
	}

	// 4. Projection阶段（减少数据传输）
	projectionStage := bson.D{{Key: "$project", Value: GetProjectionForList()}}
	pipeline = append(pipeline, projectionStage)

	return pipeline
}

// buildMatchStage 构建Match阶段
func (p *PerformanceOptimizer) buildMatchStage(filter *bookstore.BookFilter) bson.D {
	matchConditions := bson.D{}

	// 状态过滤
	if filter.Status != nil {
		matchConditions = append(matchConditions, bson.E{Key: "status", Value: *filter.Status})
	}

	// 分类过滤
	if filter.CategoryID != nil {
		matchConditions = append(matchConditions, bson.E{Key: "category_ids", Value: *filter.CategoryID})
	}

	// 作者过滤
	if filter.Author != nil {
		matchConditions = append(matchConditions, bson.E{Key: "author", Value: *filter.Author})
	}

	// 标签过滤
	if len(filter.Tags) > 0 {
		matchConditions = append(matchConditions, bson.E{Key: "tags", Value: bson.D{{Key: "$in", Value: filter.Tags}}})
	}

	// 关键词搜索
	if filter.Keyword != nil && *filter.Keyword != "" {
		searchConditions := bson.D{
			{Key: "$or", Value: []bson.D{
				{{Key: "title", Value: bson.D{{Key: "$regex", Value: *filter.Keyword}, {Key: "$options", Value: "i"}}}},
				{{Key: "author", Value: bson.D{{Key: "$regex", Value: *filter.Keyword}, {Key: "$options", Value: "i"}}}},
				{{Key: "introduction", Value: bson.D{{Key: "$regex", Value: *filter.Keyword}, {Key: "$options", Value: "i"}}}},
			}},
		}
		matchConditions = append(matchConditions, searchConditions...)
	}

	if len(matchConditions) > 0 {
		return bson.D{{Key: "$match", Value: matchConditions}}
	}

	return bson.D{}
}

// buildSortStage 构建Sort阶段
func (p *PerformanceOptimizer) buildSortStage(filter *bookstore.BookFilter) bson.D {
	sortField := filter.SortBy
	if sortField == "" {
		sortField = "created_at"
	}

	sortOrder := -1
	if filter.SortOrder == "asc" {
		sortOrder = 1
	}

	return bson.D{{Key: "$sort", Value: bson.D{{Key: sortField, Value: sortOrder}}}}
}

// ============================================
// Cursor Optimization
// ============================================

// OptimizeCursorBatchSize 优化游标批量大小
func (p *PerformanceOptimizer) OptimizeCursorBatchSize(filter *bookstore.BookFilter) int32 {
	// 根据不同的场景返回不同的批量大小
	batchSize := int32(20)

	// 如果是移动端，减少批量大小以节省内存
	if filter.Limit > 0 && filter.Limit < 20 {
		batchSize = int32(filter.Limit)
	}

	return batchSize
}

// ============================================
// Connection Pool Optimization
// ============================================

// OptimizeConnectionPool 优化连接池配置
func (p *PerformanceOptimizer) OptimizeConnectionPool() *options.ClientOptions {
	return options.Client().
		SetMaxPoolSize(100).          // 最大连接数
		SetMinPoolSize(10).           // 最小连接数
		SetMaxConnIdleTime(30 * 60)   // 连接最大空闲时间（秒）
}

// ============================================
// Bulk Operations
// ============================================

// BulkWriteOptimized 优化的批量写入
func (p *PerformanceOptimizer) BulkWriteOptimized(
	ctx context.Context,
	models []mongo.WriteModel,
	ordered bool,
) (*mongo.BulkWriteResult, error) {
	// 使用无序写入以提高性能
	opts := options.BulkWrite().
		SetOrdered(ordered).
		SetBypassDocumentValidation(false)

	return p.collection.BulkWrite(ctx, models, opts)
}

// ============================================
// Performance Monitoring
// ============================================

// ExplainQuery 解释查询计划（用于性能分析）
func (p *PerformanceOptimizer) ExplainQuery(
	ctx context.Context,
	filter bson.M,
	opts *options.FindOptions,
) (map[string]interface{}, error) {
	// 获取查询计划
	cursor, err := p.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// 执行explain
	// 注意：这需要MongoDB的支持
	return map[string]interface{}{
		"message": "Query explanation requires MongoDB profiling",
	}, nil
}
