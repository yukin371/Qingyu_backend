package sync

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	searchengine "Qingyu_backend/service/search/engine"
)

// ConsistencyChecker 一致性校验器
type ConsistencyChecker struct {
	mongoClient *mongo.Client
	mongoDB     *mongo.Database
	esEngine    searchengine.Engine
	redisClient *redis.Client
	zapLogger   *zap.Logger
	ctx         context.Context
	cancel      context.CancelFunc
	config      *CheckerConfig
}

// CheckerConfig 校验器配置
type CheckerConfig struct {
	CheckInterval   time.Duration // 检查间隔（用于定时任务）
	BatchSize       int           // 批量查询大小
	MaxMissingDocs  int           // 最大缺失文档数（超过则标记为异常）
	AutoRepair      bool          // 是否自动修复
}

// DefaultCheckerConfig 默认配置
func DefaultCheckerConfig() *CheckerConfig {
	return &CheckerConfig{
		CheckInterval:  24 * time.Hour, // 每天检查一次
		BatchSize:      1000,
		MaxMissingDocs: 100,
		AutoRepair:     false, // 默认不自动修复
	}
}

// NewConsistencyChecker 创建一致性校验器
func NewConsistencyChecker(
	mongoClient *mongo.Client,
	mongoDB *mongo.Database,
	esEngine searchengine.Engine,
	redisClient *redis.Client,
	zapLogger *zap.Logger,
	config *CheckerConfig,
) *ConsistencyChecker {
	if config == nil {
		config = DefaultCheckerConfig()
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &ConsistencyChecker{
		mongoClient: mongoClient,
		mongoDB:     mongoDB,
		esEngine:    esEngine,
		redisClient: redisClient,
		zapLogger:   zapLogger,
		ctx:         ctx,
		cancel:      cancel,
		config:      config,
	}
}

// CheckConsistency 检查指定集合的一致性
func (cc *ConsistencyChecker) CheckConsistency(ctx context.Context, collection string) (*ConsistencyReport, error) {
	report := &ConsistencyReport{
		ID:         fmt.Sprintf("consistency_%s_%d", collection, time.Now().Unix()),
		Collection: collection,
		CheckedAt:  time.Now(),
		Details:    make(map[string]interface{}),
	}

	// 1. 获取 MongoDB 文档计数
	mongoCount, err := cc.getMongoDocumentCount(ctx, collection)
	if err != nil {
		report.Status = "error"
		report.Details["mongo_error"] = err.Error()
		return report, fmt.Errorf("failed to get MongoDB count: %w", err)
	}
	report.MongoCount = mongoCount

	// 2. 获取 ES 文档计数
	indexName := cc.getElasticsearchIndex(collection)
	esCount, err := cc.getESDocumentCount(ctx, indexName)
	if err != nil {
		report.Status = "error"
		report.Details["es_error"] = err.Error()
		return report, fmt.Errorf("failed to get ES count: %w", err)
	}
	report.ESCount = esCount

	// 3. 判断是否一致
	if mongoCount == esCount {
		report.Status = "consistent"
		if cc.zapLogger != nil {
			cc.zapLogger.Info("Collection is consistent",
				zap.String("collection", collection),
				zap.Int64("count", mongoCount),
			)
		}
		return report, nil
	}

	// 4. 不一致，找出差异
	report.Status = "inconsistent"

	// 找出缺失的文档
	missingDocs, extraDocs, err := cc.FindMissingDocuments(ctx, collection)
	if err != nil {
		report.Status = "error"
		report.Details["find_error"] = err.Error()
		return report, fmt.Errorf("failed to find missing documents: %w", err)
	}

	report.MissingDocs = missingDocs
	report.ExtraDocs = extraDocs
	report.Details["missing_count"] = len(missingDocs)
	report.Details["extra_count"] = len(extraDocs)

	if cc.zapLogger != nil {
		cc.zapLogger.Warn("Collection is inconsistent",
			zap.String("collection", collection),
			zap.Int64("mongo_count", mongoCount),
			zap.Int64("es_count", esCount),
			zap.Int("missing_docs", len(missingDocs)),
			zap.Int("extra_docs", len(extraDocs)),
		)
	}

	return report, nil
}

// getMongoDocumentCount 获取 MongoDB 文档计数
func (cc *ConsistencyChecker) getMongoDocumentCount(ctx context.Context, collection string) (int64, error) {
	coll := cc.mongoDB.Collection(collection)
	return coll.CountDocuments(ctx, bson.M{})
}

// getESDocumentCount 获取 ES 文档计数
func (cc *ConsistencyChecker) getESDocumentCount(ctx context.Context, indexName string) (int64, error) {
	// ES Engine 没有 Count 方法，使用 Search 并限制返回 0 条结果来获取总数
	result, err := cc.esEngine.Search(ctx, indexName, map[string]interface{}{
		"size": 0,
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
	}, nil)
	if err != nil {
		return 0, err
	}

	if result == nil {
		return 0, fmt.Errorf("search result is nil")
	}

	return result.Total, nil
}

// FindMissingDocuments 找出 MongoDB 和 ES 之间的文档差异
func (cc *ConsistencyChecker) FindMissingDocuments(ctx context.Context, collection string) ([]string, []string, error) {
	indexName := cc.getElasticsearchIndex(collection)

	// 1. 获取 MongoDB 所有文档 ID
	mongoIDs, err := cc.getMongoDocumentIDs(ctx, collection)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get MongoDB IDs: %w", err)
	}

	// 2. 获取 ES 所有文档 ID
	esIDs, err := cc.getESDocumentIDs(ctx, indexName)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get ES IDs: %w", err)
	}

	// 3. 找出差异
	mongoIDSet := make(map[string]bool)
	for _, id := range mongoIDs {
		mongoIDSet[id] = true
	}

	esIDSet := make(map[string]bool)
	for _, id := range esIDs {
		esIDSet[id] = true
	}

	// MongoDB 有但 ES 没有的
	var missingDocs []string
	for _, id := range mongoIDs {
		if !esIDSet[id] {
			missingDocs = append(missingDocs, id)
		}
	}

	// ES 有但 MongoDB 没有的
	var extraDocs []string
	for _, id := range esIDs {
		if !mongoIDSet[id] {
			extraDocs = append(extraDocs, id)
		}
	}

	return missingDocs, extraDocs, nil
}

// getMongoDocumentIDs 获取 MongoDB 所有文档 ID
func (cc *ConsistencyChecker) getMongoDocumentIDs(ctx context.Context, collection string) ([]string, error) {
	coll := cc.mongoDB.Collection(collection)

	cursor, err := coll.Find(ctx, bson.M{}, nil)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var ids []string
	for cursor.Next(ctx) {
		var result struct {
			ID string `bson:"_id"`
		}
		if err := cursor.Decode(&result); err != nil {
			continue
		}
		ids = append(ids, result.ID)
	}

	return ids, nil
}

// getESDocumentIDs 获取 ES 所有文档 ID
func (cc *ConsistencyChecker) getESDocumentIDs(ctx context.Context, indexName string) ([]string, error) {
	// 使用 ES Engine 的 Search 方法获取所有 ID
	// 注意：这里需要实现一个批量获取 ID 的方法
	// 暂时返回空列表，实际使用时需要根据 ES Engine 的能力实现

	// TODO: 实现批量获取 ES 文档 ID 的逻辑
	// 可以使用 scroll API 或者 search after
	return []string{}, nil
}

// RepairMissingDocuments 修复不一致的文档
func (cc *ConsistencyChecker) RepairMissingDocuments(ctx context.Context, collection string, missingDocs, extraDocs []string) error {
	if cc.zapLogger != nil {
		cc.zapLogger.Info("Starting repair",
			zap.String("collection", collection),
			zap.Int("missing_docs", len(missingDocs)),
			zap.Int("extra_docs", len(extraDocs)),
		)
	}

	// 1. 补充缺失的文档（从 MongoDB 同步到 ES）
	for _, docID := range missingDocs {
		if err := cc.syncDocumentToES(ctx, collection, docID); err != nil {
			if cc.zapLogger != nil {
				cc.zapLogger.Error("Failed to sync document to ES",
					zap.String("document_id", docID),
					zap.Error(err),
				)
			}
			continue
		}

		if cc.zapLogger != nil {
			cc.zapLogger.Debug("Document synced to ES",
				zap.String("document_id", docID),
			)
		}
	}

	// 2. 删除多余的文档（从 ES 删除）
	indexName := cc.getElasticsearchIndex(collection)
	for _, docID := range extraDocs {
		if err := cc.esEngine.Delete(ctx, indexName, docID); err != nil {
			if cc.zapLogger != nil {
				cc.zapLogger.Error("Failed to delete extra document from ES",
					zap.String("document_id", docID),
					zap.Error(err),
				)
			}
			continue
		}

		if cc.zapLogger != nil {
			cc.zapLogger.Debug("Extra document deleted from ES",
				zap.String("document_id", docID),
			)
		}
	}

	return nil
}

// syncDocumentToES 同步单个文档到 ES
func (cc *ConsistencyChecker) syncDocumentToES(ctx context.Context, collection, docID string) error {
	coll := cc.mongoDB.Collection(collection)

	var doc bson.Raw
	err := coll.FindOne(ctx, bson.M{"_id": docID}).Decode(&doc)
	if err != nil {
		return fmt.Errorf("failed to find document: %w", err)
	}

	// 转换为 map
	docMap := make(map[string]interface{})
	if err := bson.Unmarshal(doc, &docMap); err != nil {
		return fmt.Errorf("failed to unmarshal document: %w", err)
	}

	// 删除 _id 字段
	delete(docMap, "_id")

	indexName := cc.getElasticsearchIndex(collection)

	// 索引到 Elasticsearch
	documents := []searchengine.Document{
		{
			ID:     docID,
			Source: docMap,
		},
	}

	return cc.esEngine.Index(ctx, indexName, documents)
}

// StartScheduledCheck 启动定时检查任务
func (cc *ConsistencyChecker) StartScheduledCheck(collection string, reportHandler func(*ConsistencyReport)) {
	go func() {
		ticker := time.NewTicker(cc.config.CheckInterval)
		defer ticker.Stop()

		// 首次立即执行一次
		cc.performCheck(collection, reportHandler)

		for {
			select {
			case <-cc.ctx.Done():
				if cc.zapLogger != nil {
					cc.zapLogger.Info("Consistency checker stopped")
				}
				return
			case <-ticker.C:
				cc.performCheck(collection, reportHandler)
			}
		}
	}()
}

// performCheck 执行检查
func (cc *ConsistencyChecker) performCheck(collection string, reportHandler func(*ConsistencyReport)) {
	ctx, cancel := context.WithTimeout(cc.ctx, 5*time.Minute)
	defer cancel()

	report, err := cc.CheckConsistency(ctx, collection)
	if err != nil {
		if cc.zapLogger != nil {
			cc.zapLogger.Error("Consistency check failed",
				zap.String("collection", collection),
				zap.Error(err),
			)
		}
		return
	}

	// 调用报告处理器
	if reportHandler != nil {
		reportHandler(report)
	}

	// 如果配置了自动修复，则自动修复不一致
	if cc.config.AutoRepair && report.Status == "inconsistent" {
		if len(report.MissingDocs) > 0 || len(report.ExtraDocs) > 0 {
			if err := cc.RepairMissingDocuments(ctx, collection, report.MissingDocs, report.ExtraDocs); err != nil {
				if cc.zapLogger != nil {
					cc.zapLogger.Error("Auto repair failed",
						zap.String("collection", collection),
						zap.Error(err),
					)
				}
			}
		}
	}
}

// getElasticsearchIndex 获取 ES 索引名称
func (cc *ConsistencyChecker) getElasticsearchIndex(collection string) string {
	indexMap := map[string]string{
		"books":     "books_search",
		"projects":  "projects_search",
		"documents": "documents_search",
		"users":     "users_search",
	}

	if indexName, ok := indexMap[collection]; ok {
		return indexName
	}

	return collection + "_search"
}

// Stop 停止校验器
func (cc *ConsistencyChecker) Stop() {
	cc.cancel()
	if cc.zapLogger != nil {
		cc.zapLogger.Info("Consistency checker stopped")
	}
}

// CheckAllCollections 检查所有集合的一致性
func (cc *ConsistencyChecker) CheckAllCollections(ctx context.Context) ([]*ConsistencyReport, error) {
	collections := []string{"books", "projects", "documents", "users"}

	var reports []*ConsistencyReport
	for _, collection := range collections {
		report, err := cc.CheckConsistency(ctx, collection)
		if err != nil {
			if cc.zapLogger != nil {
				cc.zapLogger.Error("Failed to check collection",
					zap.String("collection", collection),
					zap.Error(err),
				)
			}
			continue
		}
		reports = append(reports, report)
	}

	return reports, nil
}
