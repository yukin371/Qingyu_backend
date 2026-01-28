package engine

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/olivere/elastic/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =========================
// 辅助函数
// =========================

// getESClient 获取 ES 客户端，如果环境变量未设置则跳过测试
func getESClient(t *testing.T) *elastic.Client {
	t.Helper()

	esURL := os.Getenv("ELASTICSEARCH_URL")
	if esURL == "" {
		t.Skip("ELASTICSEARCH_URL not set, skipping ES tests")
	}

	client, err := elastic.NewClient(
		elastic.SetURL(esURL),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false),
	)
	require.NoError(t, err, "Failed to create ES client")

	return client
}

// getTestEngine 获取测试用的 ES 引擎
func getTestEngine(t *testing.T) *ElasticsearchEngine {
	t.Helper()

	client := getESClient(t)
	engine, err := NewElasticsearchEngine(client)
	require.NoError(t, err, "Failed to create ES engine")
	return engine
}

// generateTestIndexName 生成唯一的测试索引名
func generateTestIndexName(base string) string {
	return fmt.Sprintf("%s_%d", base, time.Now().UnixNano())
}

// cleanupTestIndex 清理测试索引
func cleanupTestIndex(t *testing.T, ctx context.Context, engine *ElasticsearchEngine, indexName string) {
	t.Helper()

	// 忽略删除错误，可能索引已经不存在
	_ = engine.DeleteIndex(ctx, indexName)
}

// =========================
// 基础测试
// =========================

// TestElasticsearchEngine_NewEngine 测试创建 Elasticsearch 引擎
func TestElasticsearchEngine_NewEngine(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := getESClient(t)

	// 测试正常创建
	engine, err := NewElasticsearchEngine(client)
	assert.NoError(t, err)
	assert.NotNil(t, engine)

	// 通过 Health 方法验证引擎可用
	ctx := context.Background()
	err = engine.Health(ctx)
	assert.NoError(t, err, "Engine health check should pass")

	// 测试 nil 客户端
	_, err = NewElasticsearchEngine(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be nil")
}

// TestElasticsearchEngine_Health 测试健康检查
func TestElasticsearchEngine_Health(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	engine := getTestEngine(t)
	ctx := context.Background()

	// 测试健康检查
	err := engine.Health(ctx)
	assert.NoError(t, err, "Health check should pass")
}

// =========================
// 索引管理测试
// =========================

// TestElasticsearchEngine_CreateIndex 测试创建索引
func TestElasticsearchEngine_CreateIndex(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	engine := getTestEngine(t)
	ctx := context.Background()
	indexName := generateTestIndexName("test_create_index")
	defer cleanupTestIndex(t, ctx, engine, indexName)

	// 测试创建带映射的索引
	mapping := map[string]interface{}{
		"settings": map[string]interface{}{
			"number_of_shards":   1,
			"number_of_replicas": 0,
		},
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"title": map[string]interface{}{
					"type": "text",
				},
				"author": map[string]interface{}{
					"type": "keyword",
				},
			},
		},
	}

	err := engine.CreateIndex(ctx, indexName, mapping)
	assert.NoError(t, err, "Failed to create index")

	// 尝试创建已存在的索引（应该失败）
	err = engine.CreateIndex(ctx, indexName, mapping)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}

// TestElasticsearchEngine_CreateIndexWithDefaultMapping 测试使用默认映射创建索引
func TestElasticsearchEngine_CreateIndexWithDefaultMapping(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	engine := getTestEngine(t)
	ctx := context.Background()
	indexName := generateTestIndexName("test_default_mapping")
	defer cleanupTestIndex(t, ctx, engine, indexName)

	// 使用 nil 映射，应该使用默认映射
	err := engine.CreateIndex(ctx, indexName, nil)
	assert.NoError(t, err, "Failed to create index with default mapping")

	// 验证删除索引功能
	err = engine.DeleteIndex(ctx, indexName)
	assert.NoError(t, err, "Should be able to delete index")

	// 再次删除应该失败
	err = engine.DeleteIndex(ctx, indexName)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}

// =========================
// 文档索引测试
// =========================

// TestElasticsearchEngine_Index 测试索引文档
func TestElasticsearchEngine_Index(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	engine := getTestEngine(t)
	ctx := context.Background()
	indexName := generateTestIndexName("test_index_docs")
	defer cleanupTestIndex(t, ctx, engine, indexName)

	// 先创建索引
	mapping := map[string]interface{}{
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"title":  map[string]interface{}{"type": "text"},
				"author": map[string]interface{}{"type": "keyword"},
				"price":  map[string]interface{}{"type": "float"},
			},
		},
	}
	err := engine.CreateIndex(ctx, indexName, mapping)
	require.NoError(t, err)

	// 等待索引就绪
	time.Sleep(500 * time.Millisecond)

	// 测试索引单个文档
	docs := []Document{
		{
			ID: "1",
			Source: map[string]interface{}{
				"title":  "测试书籍1",
				"author": "作者A",
				"price":  29.99,
			},
		},
		{
			ID: "2",
			Source: map[string]interface{}{
				"title":  "测试书籍2",
				"author": "作者B",
				"price":  39.99,
			},
		},
	}

	err = engine.Index(ctx, indexName, docs)
	assert.NoError(t, err, "Failed to index documents")

	// 等待 ES 索引刷新
	time.Sleep(1 * time.Second)

	// 通过搜索验证文档已索引
	opts := &SearchOptions{
		From: 0,
		Size: 10,
	}
	result, err := engine.Search(ctx, indexName, "测试书籍1", opts)
	assert.NoError(t, err, "Search should succeed")
	assert.Greater(t, result.Total, int64(0), "Should find indexed document")
}

// TestElasticsearchEngine_IndexEmpty 测试索引空文档列表
func TestElasticsearchEngine_IndexEmpty(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	engine := getTestEngine(t)
	ctx := context.Background()
	indexName := generateTestIndexName("test_index_empty")
	defer cleanupTestIndex(t, ctx, engine, indexName)

	// 创建索引
	mapping := map[string]interface{}{"mappings": map[string]interface{}{}}
	err := engine.CreateIndex(ctx, indexName, mapping)
	require.NoError(t, err)

	// 索引空列表应该不报错
	err = engine.Index(ctx, indexName, []Document{})
	assert.NoError(t, err)
}

// =========================
// 文档更新测试
// =========================

// TestElasticsearchEngine_Update 测试更新文档
func TestElasticsearchEngine_Update(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	engine := getTestEngine(t)
	ctx := context.Background()
	indexName := generateTestIndexName("test_update")
	defer cleanupTestIndex(t, ctx, engine, indexName)

	// 创建索引并添加文档
	mapping := map[string]interface{}{
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"title": map[string]interface{}{"type": "text"},
				"price": map[string]interface{}{"type": "float"},
			},
		},
	}
	err := engine.CreateIndex(ctx, indexName, mapping)
	require.NoError(t, err)

	docs := []Document{
		{
			ID: "1",
			Source: map[string]interface{}{
				"title": "原标题",
				"price": 10.0,
			},
		},
	}
	err = engine.Index(ctx, indexName, docs)
	require.NoError(t, err)

	// 等待索引
	time.Sleep(1 * time.Second)

	// 更新文档
	updateDoc := Document{
		Source: map[string]interface{}{
			"title": "新标题",
			"price": 20.0,
		},
	}
	err = engine.Update(ctx, indexName, "1", updateDoc)
	assert.NoError(t, err, "Failed to update document")

	// 等待索引更新
	time.Sleep(1 * time.Second)

	// 通过搜索验证更新
	opts := &SearchOptions{
		From: 0,
		Size: 10,
	}
	result, err := engine.Search(ctx, indexName, "新标题", opts)
	assert.NoError(t, err)
	assert.Greater(t, result.Total, int64(0), "Should find updated document")

	// 验证内容
	if len(result.Hits) > 0 {
		assert.Equal(t, "新标题", result.Hits[0].Source["title"])
	}
}

// TestElasticsearchEngine_UpdateNotFound 测试更新不存在的文档
func TestElasticsearchEngine_UpdateNotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	engine := getTestEngine(t)
	ctx := context.Background()
	indexName := generateTestIndexName("test_update_not_found")
	defer cleanupTestIndex(t, ctx, engine, indexName)

	// 创建索引
	mapping := map[string]interface{}{"mappings": map[string]interface{}{}}
	err := engine.CreateIndex(ctx, indexName, mapping)
	require.NoError(t, err)

	// 尝试更新不存在的文档
	updateDoc := Document{
		Source: map[string]interface{}{
			"title": "新标题",
		},
	}
	err = engine.Update(ctx, indexName, "nonexistent", updateDoc)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

// =========================
// 文档删除测试
// =========================

// TestElasticsearchEngine_Delete 测试删除文档
func TestElasticsearchEngine_Delete(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	engine := getTestEngine(t)
	ctx := context.Background()
	indexName := generateTestIndexName("test_delete")
	defer cleanupTestIndex(t, ctx, engine, indexName)

	// 创建索引并添加文档
	mapping := map[string]interface{}{
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"title": map[string]interface{}{"type": "text"},
			},
		},
	}
	err := engine.CreateIndex(ctx, indexName, mapping)
	require.NoError(t, err)

	docs := []Document{
		{
			ID: "1",
			Source: map[string]interface{}{
				"title": "要删除的文档",
			},
		},
	}
	err = engine.Index(ctx, indexName, docs)
	require.NoError(t, err)

	// 等待索引
	time.Sleep(1 * time.Second)

	// 验证文档存在
	opts := &SearchOptions{From: 0, Size: 10}
	result, err := engine.Search(ctx, indexName, "要删除的文档", opts)
	require.NoError(t, err)
	assert.Greater(t, result.Total, int64(0), "Document should exist before deletion")

	// 删除文档
	err = engine.Delete(ctx, indexName, "1")
	assert.NoError(t, err, "Failed to delete document")

	// 等待索引更新
	time.Sleep(1 * time.Second)

	// 验证删除 - 搜索应该找不到文档
	result, err = engine.Search(ctx, indexName, "要删除的文档", opts)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), result.Total, "Document should be deleted")
}

// TestElasticsearchEngine_DeleteNotFound 测试删除不存在的文档
func TestElasticsearchEngine_DeleteNotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	engine := getTestEngine(t)
	ctx := context.Background()
	indexName := generateTestIndexName("test_delete_not_found")
	defer cleanupTestIndex(t, ctx, engine, indexName)

	// 创建索引
	mapping := map[string]interface{}{"mappings": map[string]interface{}{}}
	err := engine.CreateIndex(ctx, indexName, mapping)
	require.NoError(t, err)

	// 尝试删除不存在的文档
	err = engine.Delete(ctx, indexName, "nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

// =========================
// 搜索功能测试
// =========================

// TestElasticsearchEngine_SearchKeyword 测试关键词搜索
func TestElasticsearchEngine_SearchKeyword(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	engine := getTestEngine(t)
	ctx := context.Background()
	indexName := generateTestIndexName("test_search_keyword")
	defer cleanupTestIndex(t, ctx, engine, indexName)

	// 创建索引并添加测试数据
	mapping := map[string]interface{}{
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"title":       map[string]interface{}{"type": "text"},
				"content":     map[string]interface{}{"type": "text"},
				"description": map[string]interface{}{"type": "text"},
			},
		},
	}
	err := engine.CreateIndex(ctx, indexName, mapping)
	require.NoError(t, err)

	docs := []Document{
		{
			ID: "1",
			Source: map[string]interface{}{
				"title":   "Go 语言编程",
				"content": "这是一本关于 Go 语言的编程书籍",
			},
		},
		{
			ID: "2",
			Source: map[string]interface{}{
				"title":   "Python 入门",
				"content": "Python 是一门易学的编程语言",
			},
		},
		{
			ID: "3",
			Source: map[string]interface{}{
				"title":   "Java 高级特性",
				"content": "深入学习 Java 的高级特性",
			},
		},
	}
	err = engine.Index(ctx, indexName, docs)
	require.NoError(t, err)

	// 等待索引刷新
	time.Sleep(1 * time.Second)

	// 测试关键词搜索
	opts := &SearchOptions{
		From: 0,
		Size: 10,
	}
	result, err := engine.Search(ctx, indexName, "Go", opts)
	assert.NoError(t, err, "Search should succeed")
	assert.Greater(t, result.Total, int64(0), "Should find results")
	assert.Len(t, result.Hits, 2, "Should find 2 documents containing 'Go'")
}

// TestElasticsearchEngine_SearchBooleanQuery 测试布尔查询
func TestElasticsearchEngine_SearchBooleanQuery(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	engine := getTestEngine(t)
	ctx := context.Background()
	indexName := generateTestIndexName("test_search_boolean")
	defer cleanupTestIndex(t, ctx, engine, indexName)

	// 创建索引并添加测试数据
	mapping := map[string]interface{}{
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"title":  map[string]interface{}{"type": "text"},
				"author": map[string]interface{}{"type": "keyword"},
				"price":  map[string]interface{}{"type": "integer"},
			},
		},
	}
	err := engine.CreateIndex(ctx, indexName, mapping)
	require.NoError(t, err)

	docs := []Document{
		{
			ID: "1",
			Source: map[string]interface{}{
				"title":  "Go 语言编程",
				"author": "张三",
				"price":  50,
			},
		},
		{
			ID: "2",
			Source: map[string]interface{}{
				"title":  "Python 入门",
				"author": "李四",
				"price":  30,
			},
		},
		{
			ID: "3",
			Source: map[string]interface{}{
				"title":  "Go 实战",
				"author": "张三",
				"price":  60,
			},
		},
	}
	err = engine.Index(ctx, indexName, docs)
	require.NoError(t, err)

	// 等待索引刷新
	time.Sleep(1 * time.Second)

	// 测试布尔查询
	query := map[string]interface{}{
		"bool": map[string]interface{}{
			"must": []interface{}{
				map[string]interface{}{
					"match": map[string]interface{}{
						"title": "Go",
					},
				},
			},
			"filter": []interface{}{
				map[string]interface{}{
					"term": map[string]interface{}{
						"author": "张三",
					},
				},
			},
		},
	}

	opts := &SearchOptions{
		From: 0,
		Size: 10,
	}
	result, err := engine.Search(ctx, indexName, query, opts)
	assert.NoError(t, err)
	assert.Greater(t, result.Total, int64(0), "Should find results")
}

// TestElasticsearchEngine_SearchWithFilter 测试带过滤条件的搜索
func TestElasticsearchEngine_SearchWithFilter(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	engine := getTestEngine(t)
	ctx := context.Background()
	indexName := generateTestIndexName("test_search_filter")
	defer cleanupTestIndex(t, ctx, engine, indexName)

	// 创建索引并添加测试数据
	mapping := map[string]interface{}{
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"title":  map[string]interface{}{"type": "text"},
				"author": map[string]interface{}{"type": "keyword"},
				"status": map[string]interface{}{"type": "keyword"},
			},
		},
	}
	err := engine.CreateIndex(ctx, indexName, mapping)
	require.NoError(t, err)

	docs := []Document{
		{
			ID: "1",
			Source: map[string]interface{}{
				"title":  "Go 语言编程",
				"author": "张三",
				"status": "published",
			},
		},
		{
			ID: "2",
			Source: map[string]interface{}{
				"title":  "Python 入门",
				"author": "李四",
				"status": "draft",
			},
		},
		{
			ID: "3",
			Source: map[string]interface{}{
				"title":  "Java 实战",
				"author": "王五",
				"status": "published",
			},
		},
	}
	err = engine.Index(ctx, indexName, docs)
	require.NoError(t, err)

	// 等待索引刷新
	time.Sleep(1 * time.Second)

	// 测试带过滤条件的搜索
	opts := &SearchOptions{
		From:   0,
		Size:   10,
		Filter: map[string]interface{}{"status": "published"},
	}
	result, err := engine.Search(ctx, indexName, "编程", opts)
	assert.NoError(t, err)
	assert.Greater(t, result.Total, int64(0), "Should find published results")
}

// TestElasticsearchEngine_SearchWithSort 测试带排序的搜索
func TestElasticsearchEngine_SearchWithSort(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	engine := getTestEngine(t)
	ctx := context.Background()
	indexName := generateTestIndexName("test_search_sort")
	defer cleanupTestIndex(t, ctx, engine, indexName)

	// 创建索引并添加测试数据
	mapping := map[string]interface{}{
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"title": map[string]interface{}{"type": "text"},
				"price": map[string]interface{}{"type": "integer"},
			},
		},
	}
	err := engine.CreateIndex(ctx, indexName, mapping)
	require.NoError(t, err)

	docs := []Document{
		{
			ID: "1",
			Source: map[string]interface{}{
				"title": "书籍1",
				"price": 30,
			},
		},
		{
			ID: "2",
			Source: map[string]interface{}{
				"title": "书籍2",
				"price": 50,
			},
		},
		{
			ID: "3",
			Source: map[string]interface{}{
				"title": "书籍3",
				"price": 20,
			},
		},
	}
	err = engine.Index(ctx, indexName, docs)
	require.NoError(t, err)

	// 等待索引刷新
	time.Sleep(1 * time.Second)

	// 测试按价格降序排序
	opts := &SearchOptions{
		From: 0,
		Size: 10,
		Sort: []SortField{
			{Field: "price", Ascending: false},
		},
	}
	result, err := engine.Search(ctx, indexName, "书籍", opts)
	assert.NoError(t, err)
	assert.Greater(t, result.Total, int64(0), "Should find results")

	// 验证排序结果
	if len(result.Hits) >= 2 {
		firstPrice := result.Hits[0].Source["price"].(int)
		secondPrice := result.Hits[1].Source["price"].(int)
		assert.Greater(t, firstPrice, secondPrice, "Results should be sorted by price descending")
	}
}

// TestElasticsearchEngine_SearchWithPagination 测试分页搜索
func TestElasticsearchEngine_SearchWithPagination(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	engine := getTestEngine(t)
	ctx := context.Background()
	indexName := generateTestIndexName("test_search_pagination")
	defer cleanupTestIndex(t, ctx, engine, indexName)

	// 创建索引并添加测试数据
	mapping := map[string]interface{}{
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"title": map[string]interface{}{"type": "text"},
			},
		},
	}
	err := engine.CreateIndex(ctx, indexName, mapping)
	require.NoError(t, err)

	// 添加 15 个文档
	docs := make([]Document, 15)
	for i := 0; i < 15; i++ {
		docs[i] = Document{
			ID: fmt.Sprintf("%d", i+1),
			Source: map[string]interface{}{
				"title": fmt.Sprintf("测试文档%d", i+1),
			},
		}
	}
	err = engine.Index(ctx, indexName, docs)
	require.NoError(t, err)

	// 等待索引刷新
	time.Sleep(1 * time.Second)

	// 测试第一页
	opts := &SearchOptions{
		From: 0,
		Size: 5,
	}
	result, err := engine.Search(ctx, indexName, "测试", opts)
	assert.NoError(t, err)
	assert.Equal(t, int64(15), result.Total, "Should have 15 total results")
	assert.Len(t, result.Hits, 5, "Should return 5 results on first page")

	// 测试第二页
	opts.From = 5
	result, err = engine.Search(ctx, indexName, "测试", opts)
	assert.NoError(t, err)
	assert.Len(t, result.Hits, 5, "Should return 5 results on second page")
}

// TestElasticsearchEngine_SearchWithHighlight 测试高亮搜索
func TestElasticsearchEngine_SearchWithHighlight(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	engine := getTestEngine(t)
	ctx := context.Background()
	indexName := generateTestIndexName("test_search_highlight")
	defer cleanupTestIndex(t, ctx, engine, indexName)

	// 创建索引并添加测试数据
	mapping := map[string]interface{}{
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"title":   map[string]interface{}{"type": "text"},
				"content": map[string]interface{}{"type": "text"},
			},
		},
	}
	err := engine.CreateIndex(ctx, indexName, mapping)
	require.NoError(t, err)

	docs := []Document{
		{
			ID: "1",
			Source: map[string]interface{}{
				"title":   "Go 语言编程",
				"content": "Go 是一门强大的编程语言",
			},
		},
	}
	err = engine.Index(ctx, indexName, docs)
	require.NoError(t, err)

	// 等待索引刷新
	time.Sleep(1 * time.Second)

	// 测试高亮搜索
	opts := &SearchOptions{
		From: 0,
		Size: 10,
		Highlight: &HighlightConfig{
			Fields:       []string{"title", "content"},
			PreTags:      []string{"<em>"},
			PostTags:     []string{"</em>"},
			FragmentSize: 100,
		},
	}
	result, err := engine.Search(ctx, indexName, "Go", opts)
	assert.NoError(t, err)
	assert.Greater(t, result.Total, int64(0), "Should find results")

	// 验证高亮结果
	if len(result.Hits) > 0 && result.Hits[0].Highlight != nil {
		assert.NotEmpty(t, result.Hits[0].Highlight, "Should have highlight results")
	}
}

// =========================
// 范围查询测试
// =========================

// TestElasticsearchEngine_SearchRange 测试范围查询
func TestElasticsearchEngine_SearchRange(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	engine := getTestEngine(t)
	ctx := context.Background()
	indexName := generateTestIndexName("test_search_range")
	defer cleanupTestIndex(t, ctx, engine, indexName)

	// 创建索引并添加测试数据
	mapping := map[string]interface{}{
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"title": map[string]interface{}{"type": "text"},
				"price": map[string]interface{}{"type": "integer"},
			},
		},
	}
	err := engine.CreateIndex(ctx, indexName, mapping)
	require.NoError(t, err)

	docs := []Document{
		{
			ID: "1",
			Source: map[string]interface{}{
				"title": "书籍1",
				"price": 10,
			},
		},
		{
			ID: "2",
			Source: map[string]interface{}{
				"title": "书籍2",
				"price": 30,
			},
		},
		{
			ID: "3",
			Source: map[string]interface{}{
				"title": "书籍3",
				"price": 50,
			},
		},
	}
	err = engine.Index(ctx, indexName, docs)
	require.NoError(t, err)

	// 等待索引刷新
	time.Sleep(1 * time.Second)

	// 测试范围查询：价格在 20 到 60 之间
	query := map[string]interface{}{
		"bool": map[string]interface{}{
			"filter": []interface{}{
				map[string]interface{}{
					"range": map[string]interface{}{
						"price": map[string]interface{}{
							"gte": 20,
							"lte": 60,
						},
					},
				},
			},
		},
	}

	opts := &SearchOptions{
		From: 0,
		Size: 10,
	}
	result, err := engine.Search(ctx, indexName, query, opts)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, result.Total, int64(2), "Should find at least 2 documents in range")
}

// =========================
// 完整工作流测试
// =========================

// TestElasticsearchEngine_FullWorkflow 测试完整的索引和搜索工作流
func TestElasticsearchEngine_FullWorkflow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	engine := getTestEngine(t)
	ctx := context.Background()
	indexName := generateTestIndexName("test_full_workflow")
	defer cleanupTestIndex(t, ctx, engine, indexName)

	// 1. 创建索引
	mapping := map[string]interface{}{
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"title":       map[string]interface{}{"type": "text"},
				"author":      map[string]interface{}{"type": "keyword"},
				"price":       map[string]interface{}{"type": "integer"},
				"description": map[string]interface{}{"type": "text"},
			},
		},
	}
	err := engine.CreateIndex(ctx, indexName, mapping)
	require.NoError(t, err)

	// 2. 索引文档
	docs := []Document{
		{
			ID: "1",
			Source: map[string]interface{}{
				"title":       "Go 语言实战",
				"author":      "张三",
				"price":       59,
				"description": "Go 语言实战教程",
			},
		},
		{
			ID: "2",
			Source: map[string]interface{}{
				"title":       "Python 基础",
				"author":      "李四",
				"price":       39,
				"description": "Python 入门教程",
			},
		},
	}
	err = engine.Index(ctx, indexName, docs)
	require.NoError(t, err)

	// 等待索引刷新
	time.Sleep(1 * time.Second)

	// 3. 搜索文档
	opts := &SearchOptions{
		From: 0,
		Size: 10,
	}
	result, err := engine.Search(ctx, indexName, "Go", opts)
	require.NoError(t, err)
	assert.Greater(t, result.Total, int64(0), "Should find results")

	// 4. 更新文档
	updateDoc := Document{
		Source: map[string]interface{}{
			"price": 69,
		},
	}
	err = engine.Update(ctx, indexName, "1", updateDoc)
	require.NoError(t, err)

	// 等待索引刷新
	time.Sleep(1 * time.Second)

	// 5. 验证更新 - 通过搜索验证
	result, err = engine.Search(ctx, indexName, "Go", opts)
	require.NoError(t, err)
	assert.Greater(t, result.Total, int64(0), "Should find updated document")
	if len(result.Hits) > 0 {
		price := result.Hits[0].Source["price"]
		assert.Equal(t, 69, int(price.(float64)), "Price should be updated")
	}

	// 6. 删除文档
	err = engine.Delete(ctx, indexName, "2")
	require.NoError(t, err)

	// 等待索引刷新
	time.Sleep(1 * time.Second)

	// 7. 验证删除 - 搜索应该找不到文档
	result, err = engine.Search(ctx, indexName, "Python", opts)
	require.NoError(t, err)
	assert.Equal(t, int64(0), result.Total, "Document should be deleted")
}
