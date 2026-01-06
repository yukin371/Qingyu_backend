package search

import (
	"context"
	"fmt"

	bookstoreRepo "Qingyu_backend/repository/interfaces/bookstore"
	writerRepo "Qingyu_backend/repository/interfaces/writer"
)

// TODO(Phase3): 集成Elasticsearch实现高级搜索
// 当前使用MongoDB全文索引，性能满足MVP需求
// 后续可切换到Elasticsearch以支持:
// - 更复杂的查询语法
// - 更好的中文分词
// - 更高的查询性能
// - 搜索结果高亮
// 优先级: P1
// 预计工时: 3天

// SearchService 搜索服务接口
type SearchService interface {
	// 书籍搜索
	SearchBooks(ctx context.Context, req *SearchRequest) (*SearchResult, error)

	// 文档搜索
	SearchDocuments(ctx context.Context, req *SearchRequest) (*SearchResult, error)

	// 搜索建议（自动补全）
	GetSuggestions(ctx context.Context, keyword string, limit int) ([]string, error)

	// TODO(Phase3): 搜索历史记录
	// SaveSearchHistory(ctx context.Context, userID, keyword string) error
	// GetSearchHistory(ctx context.Context, userID string, limit int) ([]string, error)

	// TODO(Phase3): 热门搜索
	// GetHotSearches(ctx context.Context, limit int) ([]string, error)

	// BaseService接口
	Initialize(ctx context.Context) error
	Health(ctx context.Context) error
	Close(ctx context.Context) error
	GetServiceName() string
	GetVersion() string
}

// SearchRequest 搜索请求
type SearchRequest struct {
	Keyword   string                 `json:"keyword"`              // 搜索关键词
	Category  string                 `json:"category,omitempty"`   // 分类过滤
	Author    string                 `json:"author,omitempty"`     // 作者过滤
	Tags      []string               `json:"tags,omitempty"`       // 标签过滤
	SortBy    string                 `json:"sort_by,omitempty"`    // 排序字段: relevance, time, popularity
	Page      int                    `json:"page"`                 // 页码
	PageSize  int                    `json:"page_size"`            // 每页数量
	Filters   map[string]interface{} `json:"filters,omitempty"`    // 额外过滤条件
	ProjectID string                 `json:"project_id,omitempty"` // 项目ID（文档搜索用）
}

// SearchResult 搜索结果
type SearchResult struct {
	Items      []interface{} `json:"items"`       // 搜索结果项
	Total      int64         `json:"total"`       // 总数
	Page       int           `json:"page"`        // 当前页
	PageSize   int           `json:"page_size"`   // 每页数量
	TotalPages int           `json:"total_pages"` // 总页数
	Keyword    string        `json:"keyword"`     // 搜索关键词
}

// SearchServiceImpl 搜索服务实现（MongoDB）
type SearchServiceImpl struct {
	bookRepo     bookstoreRepo.BookRepository
	documentRepo writerRepo.DocumentRepository
	initialized  bool
}

// NewSearchService 创建搜索服务
func NewSearchService(
	bookRepo bookstoreRepo.BookRepository,
	documentRepo writerRepo.DocumentRepository,
) SearchService {
	return &SearchServiceImpl{
		bookRepo:     bookRepo,
		documentRepo: documentRepo,
	}
}

// ============ 搜索功能 ============

// SearchBooks 搜索书籍
func (s *SearchServiceImpl) SearchBooks(ctx context.Context, req *SearchRequest) (*SearchResult, error) {
	if req == nil || req.Keyword == "" {
		return nil, fmt.Errorf("search keyword is required")
	}

	// 设置默认值
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}

	// TODO: 实现MongoDB全文搜索
	// 这里需要调用BookRepository的搜索方法
	// 目前返回空结果作为示例

	result := &SearchResult{
		Items:      []interface{}{},
		Total:      0,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: 0,
		Keyword:    req.Keyword,
	}

	// TODO(Phase2): 实现实际的MongoDB全文搜索
	// filter := bson.M{
	// 	"$text": bson.M{"$search": req.Keyword},
	// }
	// if req.Category != "" {
	// 	filter["category"] = req.Category
	// }
	// books, total, err := s.bookRepo.SearchBooks(ctx, filter, req.Page, req.PageSize)

	return result, nil
}

// SearchDocuments 搜索文档
func (s *SearchServiceImpl) SearchDocuments(ctx context.Context, req *SearchRequest) (*SearchResult, error) {
	if req == nil || req.Keyword == "" {
		return nil, fmt.Errorf("search keyword is required")
	}

	// 设置默认值
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}

	// TODO: 实现MongoDB全文搜索
	result := &SearchResult{
		Items:      []interface{}{},
		Total:      0,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: 0,
		Keyword:    req.Keyword,
	}

	// TODO(Phase2): 实现实际的文档搜索
	// 需要DocumentRepository支持全文搜索

	return result, nil
}

// GetSuggestions 获取搜索建议
func (s *SearchServiceImpl) GetSuggestions(ctx context.Context, keyword string, limit int) ([]string, error) {
	if keyword == "" {
		return []string{}, nil
	}

	if limit < 1 || limit > 20 {
		limit = 10
	}

	// TODO(Phase2): 实现搜索建议
	// 可以基于:
	// 1. 热门搜索词
	// 2. 书籍标题/作者前缀匹配
	// 3. 用户搜索历史

	// 暂时返回空列表
	suggestions := []string{}

	// TODO(Phase3): 实现智能搜索建议
	// - 拼音搜索
	// - 同义词扩展
	// - 纠错建议

	return suggestions, nil
}

// TODO(Phase3): 搜索历史功能
// func (s *SearchServiceImpl) SaveSearchHistory(ctx context.Context, userID, keyword string) error {
//     // 保存用户搜索历史到Redis或MongoDB
//     // 限制每个用户保存最近50条
//     return errors.New("not implemented yet")
// }

// TODO(Phase3): 热门搜索功能
// func (s *SearchServiceImpl) GetHotSearches(ctx context.Context, limit int) ([]string, error) {
//     // 从Redis或MongoDB获取热门搜索词
//     // 基于搜索频率排序
//     return nil, errors.New("not implemented yet")
// }

// ============ BaseService接口实现 ============

// Initialize 初始化服务
func (s *SearchServiceImpl) Initialize(ctx context.Context) error {
	if s.initialized {
		return nil
	}

	// TODO(Phase2): 创建全文索引
	// 需要在MongoDB中创建text索引：
	// db.books.createIndex({
	//     "title": "text",
	//     "author": "text",
	//     "description": "text",
	//     "tags": "text"
	// })

	s.initialized = true
	return nil
}

// Health 健康检查
func (s *SearchServiceImpl) Health(ctx context.Context) error {
	// 检查依赖的Repository是否健康
	if s.bookRepo == nil {
		return fmt.Errorf("book repository is nil")
	}

	// TODO: 调用Repository的Health方法

	return nil
}

// Close 关闭服务
func (s *SearchServiceImpl) Close(ctx context.Context) error {
	s.initialized = false
	return nil
}

// GetServiceName 获取服务名称
func (s *SearchServiceImpl) GetServiceName() string {
	return "SearchService"
}

// GetVersion 获取服务版本
func (s *SearchServiceImpl) GetVersion() string {
	return "1.0.0"
}
