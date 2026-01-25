package search

import "time"

// SearchType 搜索类型
type SearchType string

const (
	SearchTypeBooks     SearchType = "books"
	SearchTypeProjects  SearchType = "projects"
	SearchTypeDocuments SearchType = "documents"
	SearchTypeUsers     SearchType = "users"
	SearchTypeVector    SearchType = "vector"
)

// SortField 排序字段
type SortField struct {
	Field     string `json:"field"`     // 排序字段名
	Ascending bool   `json:"ascending"` // 是否升序
}

// SearchItem 搜索结果项
type SearchItem struct {
	ID        string                 `json:"id"`                  // 文档 ID
	Score     float64                `json:"score"`               // 相关性评分
	Data      map[string]interface{} `json:"data"`                // 文档数据
	Highlight map[string][]string    `json:"highlight,omitempty"` // 高亮片段
}

// HighlightConfig 高亮配置
type HighlightConfig struct {
	Fields       []string `json:"fields"`        // 需要高亮的字段
	PreTags      []string `json:"pre_tags"`      // 前置标签
	PostTags     []string `json:"post_tags"`     // 后置标签
	FragmentSize int      `json:"fragment_size"` // 片段大小
}

// SearchOptions 搜索选项
type SearchOptions struct {
	Page       int              `json:"page"`                 // 页码
	PageSize   int              `json:"page_size"`            // 每页数量
	Sort       []SortField      `json:"sort"`                 // 排序
	Filter     *Filter          `json:"filter"`               // 过滤条件
	Highlight  *HighlightConfig `json:"highlight,omitempty"`  // 高亮配置
	Options    map[string]any   `json:"options,omitempty"`    // 额外选项
	Explain    bool             `json:"explain,omitempty"`    // 是否返回评分说明
	Source     []string         `json:"source,omitempty"`     // 返回字段
	TrackScore bool             `json:"track_score"`          // 是否跟踪评分
}

// Pagination 分页信息
type Pagination struct {
	Page     int   `json:"page"`      // 当前页
	PageSize int   `json:"page_size"` // 每页数量
	Total    int64 `json:"total"`     // 总数
}

// SearchMetrics 搜索指标
type SearchMetrics struct {
	Took        time.Duration `json:"took"`          // 搜索耗时
	CacheHit    bool          `json:"cache_hit"`     // 是否命中缓存
	TotalDocs   int64         `json:"total_docs"`    // 总文档数
	MatchedDocs int64         `json:"matched_docs"`  // 匹配文档数
	Shards      *ShardInfo    `json:"shards,omitempty"` // 分片信息
}

// ShardInfo 分片信息
type ShardInfo struct {
	Total      int `json:"total"`       // 总分片数
	Successful int `json:"successful"`  // 成功分片数
	Failed     int `json:"failed"`      // 失败分片数
}
