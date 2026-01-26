package search

// SearchRequest 统一搜索请求
type SearchRequest struct {
	Type     SearchType            `json:"type" binding:"required"`     // 搜索类型：books, projects, documents, users, vector
	Query    string                `json:"query" binding:"required"`    // 搜索关键词
	Filter   map[string]interface{} `json:"filter"`                      // 过滤条件
	Sort     []SortField           `json:"sort"`                        // 排序字段
	Page     int                   `json:"page"`                        // 页码，默认 1
	PageSize int                   `json:"page_size"`                   // 每页数量，默认 20
	Options  map[string]interface{} `json:"options"`                     // 额外选项
}

// BookSearchFilter 书籍搜索过滤条件
type BookSearchFilter struct {
	CategoryID    string   `json:"category_id"`    // 分类 ID
	Author        string   `json:"author"`         // 作者
	Tags          []string `json:"tags"`           // 标签
	Status        string   `json:"status"`         // 状态
	WordCountMin  int      `json:"word_count_min"` // 最小字数
	WordCountMax  int      `json:"word_count_max"` // 最大字数
	RatingMin     float64  `json:"rating_min"`     // 最小评分
}

// DocumentSearchFilter 文档搜索过滤条件
type DocumentSearchFilter struct {
	ProjectID string `json:"project_id"` // 项目 ID
	Type      string `json:"type"`       // 文档类型
	Status    string `json:"status"`     // 状态
}

// UserSearchFilter 用户搜索过滤条件
type UserSearchFilter struct {
	Role       string `json:"role"`       // 角色
	IsVerified bool   `json:"is_verified"` // 是否认证
	Status     string `json:"status"`     // 状态
}
