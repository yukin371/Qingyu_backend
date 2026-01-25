package search

import "time"

// SearchResponse 统一搜索响应
type SearchResponse struct {
	Success bool          `json:"success"`           // 是否成功
	Data    *SearchData   `json:"data,omitempty"`    // 搜索数据
	Error   *ErrorInfo    `json:"error,omitempty"`   // 错误信息
	Meta    *MetaInfo     `json:"meta,omitempty"`    // 元信息
}

// SearchData 搜索数据
type SearchData struct {
	Type         SearchType     `json:"type"`                   // 搜索类型
	Total        int64          `json:"total"`                  // 总数
	Page         int            `json:"page"`                   // 当前页
	PageSize     int            `json:"page_size"`              // 每页数量
	Results      []SearchItem   `json:"results"`                // 搜索结果
	Aggregations map[string]any `json:"aggregations,omitempty"` // 聚合结果
	Took         time.Duration  `json:"took"`                   // 耗时
}

// ErrorInfo 错误信息
type ErrorInfo struct {
	Code    string `json:"code"`              // 错误代码
	Message string `json:"message"`           // 错误消息
	Details string `json:"details,omitempty"` // 详细信息
}

// MetaInfo 元信息
type MetaInfo struct {
	RequestID string `json:"request_id"` // 请求 ID
	TookMs    int64  `json:"took_ms"`    // 耗时（毫秒）
}

// BookSearchItem 书籍搜索结果
type BookSearchItem struct {
	ID           string  `json:"id"`
	Title        string  `json:"title"`
	Author       string  `json:"author"`
	Introduction string  `json:"introduction"`
	CoverURL     string  `json:"cover_url"`
	ViewCount    int64   `json:"view_count"`
	LikeCount    int64   `json:"like_count"`
	Rating       float64 `json:"rating"`
	WordCount    int64   `json:"word_count"`
	Status       string  `json:"status"`
}

// ProjectSearchItem 项目搜索结果
type ProjectSearchItem struct {
	ID                 string `json:"id"`
	Title              string `json:"title"`
	Synopsis           string `json:"synopsis"`
	Genre              string `json:"genre"`
	TargetWordCount    int    `json:"target_word_count"`
	CurrentWordCount   int    `json:"current_word_count"`
	Status             string `json:"status"`
}

// DocumentSearchItem 文档搜索结果
type DocumentSearchItem struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Type      string `json:"type"`
	WordCount int    `json:"word_count"`
}

// UserSearchItem 用户搜索结果
type UserSearchItem struct {
	ID         string `json:"id"`
	Username   string `json:"username"`
	Nickname   string `json:"nickname"`
	Bio        string `json:"bio"`
	AvatarURL  string `json:"avatar_url"`
	Role       string `json:"role"`
	IsVerified bool   `json:"is_verified"`
}
