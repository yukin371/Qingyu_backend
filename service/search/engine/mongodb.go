package engine

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	"Qingyu_backend/pkg/logger"
)

// MongoEngine MongoDB 搜索引擎实现
type MongoEngine struct {
	client   *mongo.Client
	database *mongo.Database
	logger   *logger.Logger
}

// NewMongoEngine 创建 MongoDB 引擎
func NewMongoEngine(client *mongo.Client, database *mongo.Database) (*MongoEngine, error) {
	if client == nil {
		return nil, fmt.Errorf("MongoDB client cannot be nil")
	}
	if database == nil {
		return nil, fmt.Errorf("MongoDB database cannot be nil")
	}

	return &MongoEngine{
		client:   client,
		database: database,
		logger:   logger.Get().WithModule("search-engine"),
	}, nil
}

// Search 执行 MongoDB 搜索
func (m *MongoEngine) Search(ctx context.Context, index string, query interface{}, opts *SearchOptions) (*SearchResult, error) {
	startTime := time.Now()

	// 默认选项
	if opts == nil {
		opts = &SearchOptions{
			From: 0,
			Size: 20,
		}
	}

	// 构建查询过滤器
	filter := m.buildSearchQuery(query, opts)

	// 构建排序选项
	sortOpts := m.buildSortOptions(opts.Sort)

	// 计算跳过的文档数
	skip := int64(opts.From)
	limit := int64(opts.Size)

	// 创建查找选项
	findOpts := options.Find()
	findOpts.SetSkip(skip)
	findOpts.SetLimit(limit)
	if len(sortOpts) > 0 {
		findOpts.SetSort(sortOpts)
	}

	// 获取集合
	collection := m.database.Collection(index)

	// 执行查询
	cursor, err := collection.Find(ctx, filter, findOpts)
	if err != nil {
		m.logger.Error("MongoDB search failed",
			zap.Error(err),
			zap.String("index", index),
			zap.Any("query", query),
			zap.Int("from", opts.From),
			zap.Int("size", opts.Size),
		)
		return nil, fmt.Errorf("search failed: %w", err)
	}
	defer cursor.Close(ctx)

	// 解码结果
	var rawDocs []bson.D
	if err = cursor.All(ctx, &rawDocs); err != nil {
		m.logger.Error("Failed to decode search results",
			zap.Error(err),
			zap.String("index", index),
		)
		return nil, fmt.Errorf("decode results failed: %w", err)
	}

	// 统计总数
	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		m.logger.Error("Failed to count documents",
			zap.Error(err),
			zap.String("index", index),
		)
		return nil, fmt.Errorf("count documents failed: %w", err)
	}

	// 转换为搜索结果
	hits := make([]Hit, 0, len(rawDocs))
	for _, doc := range rawDocs {
		hit := m.convertDocumentToHit(doc, opts.Highlight)
		hits = append(hits, hit)
	}

	took := time.Since(startTime)

	// 记录搜索日志
	m.logger.Info("MongoDB search completed",
		zap.String("index", index),
		zap.Int64("total", total),
		zap.Int("returned", len(hits)),
		zap.Duration("took", took),
	)

	return &SearchResult{
		Total: total,
		Hits:  hits,
		Took:  took,
	}, nil
}

// buildSearchQuery 构建搜索查询
func (m *MongoEngine) buildSearchQuery(query interface{}, opts *SearchOptions) bson.M {
	filter := bson.M{}

	// 处理关键词搜索
	if queryStr, ok := query.(string); ok && queryStr != "" {
		// 使用 $or 进行多字段搜索
		searchConditions := []bson.M{
			{"title": bson.M{"$regex": queryStr, "$options": "i"}},
			{"author": bson.M{"$regex": queryStr, "$options": "i"}},
			{"introduction": bson.M{"$regex": queryStr, "$options": "i"}},
			{"description": bson.M{"$regex": queryStr, "$options": "i"}},
			{"content": bson.M{"$regex": queryStr, "$options": "i"}},
		}
		filter["$or"] = searchConditions
	}

	// 处理过滤条件
	if opts != nil && opts.Filter != nil {
		for key, value := range opts.Filter {
			filter[key] = value
		}
	}

	return filter
}

// buildSortOptions 构建排序选项
func (m *MongoEngine) buildSortOptions(sortFields []SortField) bson.D {
	if len(sortFields) == 0 {
		// 默认按更新时间降序
		return bson.D{{Key: "updated_at", Value: -1}}
	}

	sortOpts := bson.D{}
	for _, sf := range sortFields {
		order := 1 // 默认升序
		if !sf.Ascending {
			order = -1 // 降序
		}
		sortOpts = append(sortOpts, bson.E{Key: sf.Field, Value: order})
	}

	return sortOpts
}

// convertDocumentToHit 转换文档为搜索命中项
func (m *MongoEngine) convertDocumentToHit(doc bson.D, highlightCfg *HighlightConfig) Hit {
	hit := Hit{
		Source: make(map[string]interface{}),
	}

	for _, elem := range doc {
		key := elem.Key
		value := elem.Value

		// 提取文档 ID
		if key == "_id" {
			if oid, ok := value.(primitive.ObjectID); ok {
				hit.ID = oid.Hex()
			} else {
				hit.ID = fmt.Sprintf("%v", value)
			}
			continue
		}

		// 提取评分（如果有）
		if key == "score" {
			if score, ok := value.(float64); ok {
				hit.Score = score
			}
			continue
		}

		// 其他字段放入 Source
		hit.Source[key] = value
	}

	// 处理高亮（MongoDB 不支持原生高亮，这里做简单实现）
	if highlightCfg != nil && len(highlightCfg.Fields) > 0 {
		hit.Highlight = m.buildHighlight(hit.Source, highlightCfg)
	}

	return hit
}

// buildHighlight 构建高亮片段
func (m *MongoEngine) buildHighlight(source map[string]interface{}, cfg *HighlightConfig) map[string][]string {
	highlight := make(map[string][]string)

	preTag := ""
	postTag := ""
	if len(cfg.PreTags) > 0 {
		preTag = cfg.PreTags[0]
	}
	if len(cfg.PostTags) > 0 {
		postTag = cfg.PostTags[0]
	}

	for _, field := range cfg.Fields {
		if value, exists := source[field]; exists {
			if strValue, ok := value.(string); ok {
				// 简单高亮：在字段值前后添加标签
				highlight[field] = []string{preTag + strValue + postTag}
			}
		}
	}

	return highlight
}

// Index 批量索引文档
func (m *MongoEngine) Index(ctx context.Context, index string, documents []Document) error {
	if len(documents) == 0 {
		return nil
	}

	startTime := time.Now()
	collection := m.database.Collection(index)

	// 准备批量插入的文档
	docsToInsert := make([]interface{}, len(documents))
	for i, doc := range documents {
		// 构建文档，包含 _id
		mongoDoc := bson.M{}
		for key, value := range doc.Source {
			mongoDoc[key] = value
		}

		// 设置文档 ID
		if doc.ID != "" {
			if objectID, err := primitive.ObjectIDFromHex(doc.ID); err == nil {
				mongoDoc["_id"] = objectID
			} else {
				mongoDoc["_id"] = doc.ID
			}
		}

		docsToInsert[i] = mongoDoc
	}

	// 执行批量插入
	result, err := collection.InsertMany(ctx, docsToInsert)
	if err != nil {
		m.logger.Error("MongoDB bulk insert failed",
			zap.Error(err),
			zap.String("index", index),
			zap.Int("count", len(documents)),
		)
		return fmt.Errorf("bulk insert failed: %w", err)
	}

	took := time.Since(startTime)

	m.logger.Info("MongoDB bulk insert completed",
		zap.String("index", index),
		zap.Int("inserted_count", len(result.InsertedIDs)),
		zap.Duration("took", took),
	)

	return nil
}

// Update 更新文档
func (m *MongoEngine) Update(ctx context.Context, index string, id string, document Document) error {
	startTime := time.Now()
	collection := m.database.Collection(index)

	// 构建 ID 过滤器
	var filter bson.M
	if objectID, err := primitive.ObjectIDFromHex(id); err == nil {
		filter = bson.M{"_id": objectID}
	} else {
		filter = bson.M{"_id": id}
	}

	// 构建更新文档
	update := bson.M{"$set": bson.M{}}
	for key, value := range document.Source {
		update["$set"].(bson.M)[key] = value
	}

	// 执行更新
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		m.logger.Error("MongoDB update failed",
			zap.Error(err),
			zap.String("index", index),
			zap.String("id", id),
		)
		return fmt.Errorf("update failed: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("document not found: %s", id)
	}

	took := time.Since(startTime)

	m.logger.Info("MongoDB update completed",
		zap.String("index", index),
		zap.String("id", id),
		zap.Int64("matched_count", result.MatchedCount),
		zap.Int64("modified_count", result.ModifiedCount),
		zap.Duration("took", took),
	)

	return nil
}

// Delete 删除文档
func (m *MongoEngine) Delete(ctx context.Context, index string, id string) error {
	startTime := time.Now()
	collection := m.database.Collection(index)

	// 构建 ID 过滤器
	var filter bson.M
	if objectID, err := primitive.ObjectIDFromHex(id); err == nil {
		filter = bson.M{"_id": objectID}
	} else {
		filter = bson.M{"_id": id}
	}

	// 执行删除
	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		m.logger.Error("MongoDB delete failed",
			zap.Error(err),
			zap.String("index", index),
			zap.String("id", id),
		)
		return fmt.Errorf("delete failed: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("document not found: %s", id)
	}

	took := time.Since(startTime)

	m.logger.Info("MongoDB delete completed",
		zap.String("index", index),
		zap.String("id", id),
		zap.Int64("deleted_count", result.DeletedCount),
		zap.Duration("took", took),
	)

	return nil
}

// CreateIndex 创建索引
func (m *MongoEngine) CreateIndex(ctx context.Context, index string, mapping interface{}) error {
	startTime := time.Now()
	collection := m.database.Collection(index)

	// 解析映射配置，创建文本索引
	indexModel := mongo.IndexModel{
		Keys: bson.D{},
	}

	// 如果 mapping 是 map，尝试提取字段
	if mappingMap, ok := mapping.(map[string]interface{}); ok {
		for field, config := range mappingMap {
			if field == "_id" {
				continue
			}

			// 如果配置为 "text"，则创建文本索引
			if configStr, ok := config.(string); ok && configStr == "text" {
				indexModel.Keys = append(indexModel.Keys.(bson.D), bson.E{Key: field, Value: "text"})
			} else {
				// 否则创建普通索引（降序）
				indexModel.Keys = append(indexModel.Keys.(bson.D), bson.E{Key: field, Value: -1})
			}
		}
	} else {
		// 默认在常用字段上创建文本索引
		indexModel.Keys = bson.D{
			{Key: "title", Value: "text"},
			{Key: "content", Value: "text"},
			{Key: "description", Value: "text"},
		}
	}

	// 创建索引
	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		m.logger.Error("MongoDB create index failed",
			zap.Error(err),
			zap.String("index", index),
			zap.Any("mapping", mapping),
		)
		return fmt.Errorf("create index failed: %w", err)
	}

	took := time.Since(startTime)

	m.logger.Info("MongoDB index created",
		zap.String("index", index),
		zap.Any("mapping", mapping),
		zap.Duration("took", took),
	)

	return nil
}

// Health 健康检查
func (m *MongoEngine) Health(ctx context.Context) error {
	if m.client == nil {
		return fmt.Errorf("MongoDB client is not initialized")
	}

	// 设置超时
	healthCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 执行 Ping 检查
	if err := m.client.Ping(healthCtx, nil); err != nil {
		m.logger.Error("MongoDB health check failed",
			zap.Error(err),
		)
		return fmt.Errorf("MongoDB ping failed: %w", err)
	}

	return nil
}
