package project

import (
	"time"
)

// VersionHistoryResponse 版本历史响应
type VersionHistoryResponse struct {
	Versions []*VersionInfo `json:"versions"`
	Total    int            `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"pageSize"`
}

// VersionInfo 版本信息
type VersionInfo struct {
	VersionID string    `json:"versionId"`
	Version   int       `json:"version"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"createdAt"`
	CreatedBy string    `json:"createdBy"`
	WordCount int       `json:"wordCount"`
}

// VersionDetail 版本详情
type VersionDetail struct {
	VersionID  string    `json:"versionId"`
	DocumentID string    `json:"documentId"`
	Version    int       `json:"version"`
	Content    string    `json:"content"`
	Message    string    `json:"message"`
	CreatedAt  time.Time `json:"createdAt"`
	CreatedBy  string    `json:"createdBy"`
	WordCount  int       `json:"wordCount"`
}

// VersionDiff 版本差异
type VersionDiff struct {
	FromVersion  string       `json:"fromVersion"`
	ToVersion    string       `json:"toVersion"`
	Changes      []ChangeItem `json:"changes"`
	AddedLines   int          `json:"addedLines"`
	DeletedLines int          `json:"deletedLines"`
}

// ChangeItem 变更项
type ChangeItem struct {
	Type    string `json:"type"` // added, deleted, modified
	Line    int    `json:"line"`
	Content string `json:"content"`
}
