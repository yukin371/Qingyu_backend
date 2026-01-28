package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/olivere/elastic/v7"
	"go.uber.org/zap"

	"Qingyu_backend/pkg/logger"
)

// ElasticsearchEngine Elasticsearch 搜索引擎实现
type ElasticsearchEngine struct {
	client *elastic.Client
	logger *logger.Logger
}

// NewElasticsearchEngine 创建 Elasticsearch 引擎
func NewElasticsearchEngine(client *elastic.Client) (*ElasticsearchEngine, error) {
	if client == nil {
		return nil, fmt.Errorf("Elasticsearch client cannot be nil")
	}

	return &ElasticsearchEngine{
		client: client,
		logger: logger.Get().WithModule("search-engine"),
	}, nil
}

// Search 执行 Elasticsearch 搜索
func (e *ElasticsearchEngine) Search(ctx context.Context, index string, query interface{}, opts *SearchOptions) (*SearchResult, error) {
	startTime := time.Now()

	// 默认选项
	if opts == nil {
		opts = &SearchOptions{
			From: 0,
			Size: 20,
		}
	}

	// 构建布尔查询
	boolQuery := elastic.NewBoolQuery()

	// 处理关键词搜索
	if queryStr, ok := query.(string); ok && queryStr != "" {
		// 使用 multi_match 查询多个字段
		multiMatchQuery := elastic.NewMultiMatchQuery(queryStr, "title^2", "content", "description", "introduction", "author").
			Type("best_fields").
			Fuzziness("AUTO")
		boolQuery.Must(multiMatchQuery)
	}

	// 处理复杂的查询结构
	if queryMap, ok := query.(map[string]interface{}); ok {
		// 处理 must 条件
		if mustConditions, ok := queryMap["must"].([]interface{}); ok {
			for _, condition := range mustConditions {
				boolQuery.Must(e.buildQuery(condition))
			}
		}

		// 处理 should 条件
		if shouldConditions, ok := queryMap["should"].([]interface{}); ok {
			for _, condition := range shouldConditions {
				boolQuery.Should(e.buildQuery(condition))
			}
		}

		// 处理 filter 条件
		if filterConditions, ok := queryMap["filter"].([]interface{}); ok {
			for _, condition := range filterConditions {
				boolQuery.Filter(e.buildQuery(condition))
			}
		}

		// 处理 must_not 条件
		if mustNotConditions, ok := queryMap["must_not"].([]interface{}); ok {
			for _, condition := range mustNotConditions {
				boolQuery.MustNot(e.buildQuery(condition))
			}
		}
	}

	// 添加过滤条件
	if opts.Filter != nil {
		for key, value := range opts.Filter {
			boolQuery.Filter(elastic.NewTermQuery(key, value))
		}
	}

	// 构建搜索服务
	searchService := e.client.Search().Index(index).Query(boolQuery)

	// 添加分页
	searchService.From(opts.From).Size(opts.Size)

	// 添加排序
	if len(opts.Sort) > 0 {
		for _, sortField := range opts.Sort {
			if sortField.Ascending {
				searchService.Sort(sortField.Field, true)
			} else {
				searchService.Sort(sortField.Field, false)
			}
		}
	} else {
		// 默认按评分降序排序
		searchService.Sort("_score", false)
	}

	// 添加高亮
	if opts.Highlight != nil && len(opts.Highlight.Fields) > 0 {
		highlight := elastic.NewHighlight()
		for _, field := range opts.Highlight.Fields {
			highlight = highlight.Field(field)
		}
		if len(opts.Highlight.PreTags) > 0 {
			highlight = highlight.PreTags(opts.Highlight.PreTags...)
		}
		if len(opts.Highlight.PostTags) > 0 {
			highlight = highlight.PostTags(opts.Highlight.PostTags...)
		}
		if opts.Highlight.FragmentSize > 0 {
			highlight = highlight.FragmentSize(opts.Highlight.FragmentSize)
		}
		searchService.Highlight(highlight)
	}

	// 执行搜索
	result, err := searchService.Do(ctx)
	if err != nil {
		e.logger.Error("Elasticsearch search failed",
			zap.Error(err),
			zap.String("index", index),
			zap.Any("query", query),
			zap.Int("from", opts.From),
			zap.Int("size", opts.Size),
		)
		return nil, fmt.Errorf("search failed: %w", err)
	}

	// 转换结果
	hits := make([]Hit, 0, len(result.Hits.Hits))
	for _, hit := range result.Hits.Hits {
		// 将 json.RawMessage 转换为 map[string]interface{}
		var sourceMap map[string]interface{}
		if len(hit.Source) > 0 {
			_ = json.Unmarshal(hit.Source, &sourceMap)
		}

		hits = append(hits, Hit{
			ID:     hit.Id,
			Score:  *hit.Score,
			Source: sourceMap,
			Highlight: func() map[string][]string {
				if hit.Highlight != nil {
					return hit.Highlight
				}
				return nil
			}(),
		})
	}

	took := time.Since(startTime)

	// 记录搜索日志
	e.logger.Info("Elasticsearch search completed",
		zap.String("index", index),
		zap.Int64("total", result.TotalHits()),
		zap.Int("returned", len(hits)),
		zap.Duration("took", took),
		zap.Int64("es_took_ms", result.TookInMillis),
	)

	return &SearchResult{
		Total: result.TotalHits(),
		Hits:  hits,
		Took:  took,
	}, nil
}

// buildQuery 构建查询条件
func (e *ElasticsearchEngine) buildQuery(condition interface{}) elastic.Query {
	if conditionMap, ok := condition.(map[string]interface{}); ok {
		// term 查询
		if term, ok := conditionMap["term"].(map[string]interface{}); ok {
			for field, value := range term {
				return elastic.NewTermQuery(field, value)
			}
		}

		// match 查询
		if match, ok := conditionMap["match"].(map[string]interface{}); ok {
			for field, value := range match {
				if valueStr, ok := value.(string); ok {
					return elastic.NewMatchQuery(field, valueStr)
				}
			}
		}

		// range 查询
		if rangeQuery, ok := conditionMap["range"].(map[string]interface{}); ok {
			for field, value := range rangeQuery {
				if valueMap, ok := value.(map[string]interface{}); ok {
					rangeQ := elastic.NewRangeQuery(field)
					if gte, ok := valueMap["gte"]; ok {
						rangeQ.Gte(gte)
					}
					if gt, ok := valueMap["gt"]; ok {
						rangeQ.Gt(gt)
					}
					if lte, ok := valueMap["lte"]; ok {
						rangeQ.Lte(lte)
					}
					if lt, ok := valueMap["lt"]; ok {
						rangeQ.Lt(lt)
					}
					return rangeQ
				}
			}
		}

		// bool 查询（递归）
		if boolQuery, ok := conditionMap["bool"].(map[string]interface{}); ok {
			bq := elastic.NewBoolQuery()
			if must, ok := boolQuery["must"].([]interface{}); ok {
				for _, m := range must {
					bq.Must(e.buildQuery(m))
				}
			}
			if should, ok := boolQuery["should"].([]interface{}); ok {
				for _, s := range should {
					bq.Should(e.buildQuery(s))
				}
			}
			if filter, ok := boolQuery["filter"].([]interface{}); ok {
				for _, f := range filter {
					bq.Filter(e.buildQuery(f))
				}
			}
			if mustNot, ok := boolQuery["must_not"].([]interface{}); ok {
				for _, mn := range mustNot {
					bq.MustNot(e.buildQuery(mn))
				}
			}
			return bq
		}
	}

	// 默认返回 match_all 查询
	return elastic.NewMatchAllQuery()
}

// Index 批量索引文档
func (e *ElasticsearchEngine) Index(ctx context.Context, index string, documents []Document) error {
	if len(documents) == 0 {
		return nil
	}

	startTime := time.Now()

	// 创建批量请求
	bulkRequest := e.client.Bulk()

	// 添加文档到批量请求
	for _, doc := range documents {
		// 准备索引请求
		indexReq := elastic.NewBulkIndexRequest().Index(index)

		// 如果指定了 ID，使用指定的 ID
		if doc.ID != "" {
			indexReq.Id(doc.ID)
		}

		// 设置文档内容
		indexReq.Doc(doc.Source)

		bulkRequest.Add(indexReq)
	}

	// 执行批量请求
	result, err := bulkRequest.Do(ctx)
	if err != nil {
		e.logger.Error("Elasticsearch bulk index failed",
			zap.Error(err),
			zap.String("index", index),
			zap.Int("count", len(documents)),
		)
		return fmt.Errorf("bulk index failed: %w", err)
	}

	// 检查是否有失败的文档
	if result.Errors {
		failed := 0
		for _, indexed := range result.Indexed() {
			if indexed.Error != nil {
				failed++
			}
		}
		e.logger.Warn("Elasticsearch bulk index completed with errors",
			zap.String("index", index),
			zap.Int("succeeded", len(result.Indexed())-failed),
			zap.Int("failed", failed),
		)
	}

	took := time.Since(startTime)

	e.logger.Info("Elasticsearch bulk index completed",
		zap.String("index", index),
		zap.Int("indexed_count", len(result.Indexed())),
		zap.Duration("took", took),
	)

	return nil
}

// Update 更新文档
func (e *ElasticsearchEngine) Update(ctx context.Context, index string, id string, document Document) error {
	startTime := time.Now()

	// 执行更新
	result, err := e.client.Update().Index(index).Id(id).Doc(document.Source).Do(ctx)
	if err != nil {
		// 检查是否是文档不存在的错误
		if elastic.IsNotFound(err) {
			return fmt.Errorf("document not found: %s", id)
		}
		e.logger.Error("Elasticsearch update failed",
			zap.Error(err),
			zap.String("index", index),
			zap.String("id", id),
		)
		return fmt.Errorf("update failed: %w", err)
	}

	if result.Result == "noop" {
		e.logger.Warn("Elasticsearch update noop",
			zap.String("index", index),
			zap.String("id", id),
		)
	}

	took := time.Since(startTime)

	e.logger.Info("Elasticsearch update completed",
		zap.String("index", index),
		zap.String("id", id),
		zap.String("result", result.Result),
		zap.Duration("took", took),
	)

	return nil
}

// Delete 删除文档
func (e *ElasticsearchEngine) Delete(ctx context.Context, index string, id string) error {
	startTime := time.Now()

	// 执行删除
	result, err := e.client.Delete().Index(index).Id(id).Do(ctx)
	if err != nil {
		// 检查是否是文档不存在的错误
		if elastic.IsNotFound(err) {
			return fmt.Errorf("document not found: %s", id)
		}
		e.logger.Error("Elasticsearch delete failed",
			zap.Error(err),
			zap.String("index", index),
			zap.String("id", id),
		)
		return fmt.Errorf("delete failed: %w", err)
	}

	if result.Result == "not_found" {
		return fmt.Errorf("document not found: %s", id)
	}

	took := time.Since(startTime)

	e.logger.Info("Elasticsearch delete completed",
		zap.String("index", index),
		zap.String("id", id),
		zap.String("result", result.Result),
		zap.Duration("took", took),
	)

	return nil
}

// CreateIndex 创建索引
func (e *ElasticsearchEngine) CreateIndex(ctx context.Context, index string, mapping interface{}) error {
	startTime := time.Now()

	// 检查索引是否已存在
	exists, err := e.client.IndexExists(index).Do(ctx)
	if err != nil {
		e.logger.Error("Elasticsearch check index existence failed",
			zap.Error(err),
			zap.String("index", index),
		)
		return fmt.Errorf("check index existence failed: %w", err)
	}

	if exists {
		e.logger.Warn("Elasticsearch index already exists",
			zap.String("index", index),
		)
		return fmt.Errorf("index already exists: %s", index)
	}

	// 构建索引创建请求
	createIndexRequest := e.client.CreateIndex(index)

	// 解析映射配置
	if mappingMap, ok := mapping.(map[string]interface{}); ok {
		// 检查是否有 settings 配置
		if settings, ok := mappingMap["settings"].(map[string]interface{}); ok {
			createIndexRequest.BodyJson(map[string]interface{}{
				"settings": settings,
			})
		}

		// 检查是否有 mappings 配置
		if mappings, ok := mappingMap["mappings"].(map[string]interface{}); ok {
			// 合并 settings 和 mappings
			body := map[string]interface{}{}

			if settings, ok := mappingMap["settings"].(map[string]interface{}); ok {
				body["settings"] = settings
			}

			body["mappings"] = mappings

			createIndexRequest.BodyJson(body)
		}
	} else {
		// 使用默认映射配置
		defaultMapping := e.getDefaultMapping()
		createIndexRequest.BodyJson(defaultMapping)
	}

	// 创建索引
	result, err := createIndexRequest.Do(ctx)
	if err != nil {
		e.logger.Error("Elasticsearch create index failed",
			zap.Error(err),
			zap.String("index", index),
			zap.Any("mapping", mapping),
		)
		return fmt.Errorf("create index failed: %w", err)
	}

	if !result.Acknowledged {
		return fmt.Errorf("create index not acknowledged: %s", index)
	}

	took := time.Since(startTime)

	e.logger.Info("Elasticsearch index created",
		zap.String("index", index),
		zap.Bool("acknowledged", result.Acknowledged),
		zap.Bool("shards_acknowledged", result.ShardsAcknowledged),
		zap.Duration("took", took),
	)

	return nil
}

// getDefaultMapping 获取默认映射配置
func (e *ElasticsearchEngine) getDefaultMapping() map[string]interface{} {
	return map[string]interface{}{
		"settings": map[string]interface{}{
			"number_of_shards":   1,
			"number_of_replicas": 1,
			"analysis": map[string]interface{}{
				"analyzer": map[string]interface{}{
					"ik_max_word": map[string]interface{}{
						"type":      "custom",
						"tokenizer": "ik_max_word",
					},
					"ik_smart": map[string]interface{}{
						"type":      "custom",
						"tokenizer": "ik_smart",
					},
				},
			},
		},
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"title": map[string]interface{}{
					"type":            "text",
					"analyzer":        "ik_max_word",
					"search_analyzer": "ik_smart",
					"fields": map[string]interface{}{
						"keyword": map[string]interface{}{
							"type": "keyword",
						},
					},
				},
				"content": map[string]interface{}{
					"type":            "text",
					"analyzer":        "ik_max_word",
					"search_analyzer": "ik_smart",
				},
				"description": map[string]interface{}{
					"type":            "text",
					"analyzer":        "ik_max_word",
					"search_analyzer": "ik_smart",
				},
				"introduction": map[string]interface{}{
					"type":            "text",
					"analyzer":        "ik_max_word",
					"search_analyzer": "ik_smart",
				},
				"author": map[string]interface{}{
					"type":            "text",
					"analyzer":        "ik_max_word",
					"search_analyzer": "ik_smart",
					"fields": map[string]interface{}{
						"keyword": map[string]interface{}{
							"type": "keyword",
						},
					},
				},
				"created_at": map[string]interface{}{
					"type": "date",
				},
				"updated_at": map[string]interface{}{
					"type": "date",
				},
			},
		},
	}
}

// Health 健康检查
func (e *ElasticsearchEngine) Health(ctx context.Context) error {
	if e.client == nil {
		return fmt.Errorf("Elasticsearch client is not initialized")
	}

	// 设置超时
	healthCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 获取集群健康状态
	health, err := e.client.ClusterHealth().Do(healthCtx)
	if err != nil {
		e.logger.Error("Elasticsearch health check failed",
			zap.Error(err),
		)
		return fmt.Errorf("cluster health check failed: %w", err)
	}

	// 检查健康状态
	if health.Status == "red" {
		e.logger.Error("Elasticsearch cluster health is red",
			zap.String("status", health.Status),
			zap.Int("active_shards", health.ActiveShards),
			zap.Int("relocating_shards", health.RelocatingShards),
			zap.Int("unassigned_shards", health.UnassignedShards),
		)
		return fmt.Errorf("cluster health is red")
	}

	// yellow 状态是警告，但不是错误
	if health.Status == "yellow" {
		e.logger.Warn("Elasticsearch cluster health is yellow",
			zap.String("status", health.Status),
			zap.Int("active_shards", health.ActiveShards),
			zap.Int("relocating_shards", health.RelocatingShards),
			zap.Int("unassigned_shards", health.UnassignedShards),
		)
	}

	e.logger.Info("Elasticsearch health check passed",
		zap.String("status", health.Status),
		zap.String("cluster_name", health.ClusterName),
		zap.Int("number_of_nodes", health.NumberOfNodes),
		zap.Int("active_shards", health.ActiveShards),
	)

	return nil
}

// DeleteIndex 删除索引
func (e *ElasticsearchEngine) DeleteIndex(ctx context.Context, index string) error {
	startTime := time.Now()

	// 检查索引是否存在
	exists, err := e.client.IndexExists(index).Do(ctx)
	if err != nil {
		e.logger.Error("Elasticsearch check index existence failed",
			zap.Error(err),
			zap.String("index", index),
		)
		return fmt.Errorf("check index existence failed: %w", err)
	}

	if !exists {
		e.logger.Warn("Elasticsearch index does not exist",
			zap.String("index", index),
		)
		return fmt.Errorf("index does not exist: %s", index)
	}

	// 删除索引
	result, err := e.client.DeleteIndex(index).Do(ctx)
	if err != nil {
		e.logger.Error("Elasticsearch delete index failed",
			zap.Error(err),
			zap.String("index", index),
		)
		return fmt.Errorf("delete index failed: %w", err)
	}

	if !result.Acknowledged {
		return fmt.Errorf("delete index not acknowledged: %s", index)
	}

	took := time.Since(startTime)

	e.logger.Info("Elasticsearch index deleted",
		zap.String("index", index),
		zap.Bool("acknowledged", result.Acknowledged),
		zap.Duration("took", took),
	)

	return nil
}
