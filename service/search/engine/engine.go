package engine

import (
	"context"
	"time"
)

// EngineType 搜索引擎类型
type EngineType string

const (
	EngineElasticsearch EngineType = "elasticsearch"
	EngineMilvus        EngineType = "milvus"
	EngineMongoDB       EngineType = "mongodb"
)

// Engine 搜索引擎接口
type Engine interface {
	// Search 执行搜索
	Search(ctx context.Context, index string, query interface{}, opts *SearchOptions) (*SearchResult, error)

	// Index 批量索引文档
	Index(ctx context.Context, index string, documents []Document) error

	// Update 更新文档
	Update(ctx context.Context, index string, id string, document Document) error

	// Delete 删除文档
	Delete(ctx context.Context, index string, id string) error

	// CreateIndex 创建索引
	CreateIndex(ctx context.Context, index string, mapping interface{}) error

	// Health 健康检查
	Health(ctx context.Context) error
}

// SearchOptions 搜索选项
type SearchOptions struct {
	From     int                    // 分页起始
	Size     int                    // 返回数量
	Sort     []SortField            // 排序
	Filter   map[string]interface{} // 过滤条件
	Highlight *HighlightConfig      // 高亮配置
}

// SortField 排序字段
type SortField struct {
	Field     string
	Ascending bool
}

// HighlightConfig 高亮配置
type HighlightConfig struct {
	Fields       []string
	PreTags      []string
	PostTags     []string
	FragmentSize int
}

// SearchResult 搜索结果
type SearchResult struct {
	Total int64                  `json:"total"`
	Hits  []Hit                  `json:"hits"`
	Aggs  map[string]interface{} `json:"aggs,omitempty"`
	Took  time.Duration          `json:"took"`
}

// Hit 搜索命中项
type Hit struct {
	ID       string                 `json:"id"`
	Score    float64                `json:"score"`
	Source   map[string]interface{} `json:"source"`
	Highlight map[string][]string   `json:"highlight,omitempty"`
}

// Document 文档
type Document struct {
	ID     string                 `json:"id"`
	Source map[string]interface{} `json:"source"`
}
