package document

import (
	"time"

	"Qingyu_backend/models/document"
)

// CreateDocumentRequest 创建文档请求
type CreateDocumentRequest struct {
	ProjectID    string   `json:"projectId" validate:"required"`
	ParentID     string   `json:"parentId,omitempty"`
	Title        string   `json:"title" validate:"required,min=1,max=200"`
	Type         string   `json:"type" validate:"required"`
	CharacterIDs []string `json:"characterIds,omitempty"`
	LocationIDs  []string `json:"locationIds,omitempty"`
	TimelineIDs  []string `json:"timelineIds,omitempty"`
	Tags         []string `json:"tags,omitempty"`
	Notes        string   `json:"notes,omitempty"`
	Order        int      `json:"order,omitempty"`
}

// CreateDocumentResponse 创建文档响应
type CreateDocumentResponse struct {
	DocumentID string    `json:"documentId"`
	Title      string    `json:"title"`
	Type       string    `json:"type"`
	CreatedAt  time.Time `json:"createdAt"`
}

// DocumentTreeNode 文档树节点
type DocumentTreeNode struct {
	*document.Document
	Children []*DocumentTreeNode `json:"children"`
}

// DocumentTreeResponse 文档树响应
type DocumentTreeResponse struct {
	ProjectID string              `json:"projectId"`
	Documents []*DocumentTreeNode `json:"documents"`
}

// MoveDocumentRequest 移动文档请求
type MoveDocumentRequest struct {
	DocumentID  string `json:"documentId" validate:"required"`
	NewParentID string `json:"newParentId"`
	Order       int    `json:"order"`
}

// ReorderDocumentsRequest 重新排序请求
type ReorderDocumentsRequest struct {
	ProjectID string         `json:"projectId" validate:"required"`
	ParentID  string         `json:"parentId"`
	Orders    map[string]int `json:"orders" validate:"required"` // documentID -> order
}

// UpdateDocumentRequest 更新文档请求
type UpdateDocumentRequest struct {
	Title        string   `json:"title,omitempty"`
	Notes        string   `json:"notes,omitempty"`
	Status       string   `json:"status,omitempty"`
	CharacterIDs []string `json:"characterIds,omitempty"`
	LocationIDs  []string `json:"locationIds,omitempty"`
	TimelineIDs  []string `json:"timelineIds,omitempty"`
	Tags         []string `json:"tags,omitempty"`
}

// ListDocumentsRequest 文档列表请求
type ListDocumentsRequest struct {
	ProjectID string `json:"projectId" validate:"required"`
	Page      string `json:"page"`
	PageSize  string `json:"pageSize"`
}

// ListDocumentsResponse 文档列表响应
type ListDocumentsResponse struct {
	Documents []*document.Document `json:"documents"`
	Total     int                  `json:"total"`
	Page      int                  `json:"page"`
	PageSize  int                  `json:"pageSize"`
}

// AutoSaveRequest 自动保存请求
type AutoSaveRequest struct {
	DocumentID     string `json:"documentId" validate:"required"`
	Content        string `json:"content" validate:"required"`
	CurrentVersion int    `json:"currentVersion"` // 客户端当前版本号
	SaveType       string `json:"saveType"`       // auto|manual
}

// AutoSaveResponse 自动保存响应
type AutoSaveResponse struct {
	Saved       bool      `json:"saved"`
	NewVersion  int       `json:"newVersion"`
	WordCount   int       `json:"wordCount"`
	SavedAt     time.Time `json:"savedAt"`
	HasConflict bool      `json:"hasConflict"`
}

// SaveStatusResponse 保存状态响应
type SaveStatusResponse struct {
	DocumentID     string    `json:"documentId"`
	LastSavedAt    time.Time `json:"lastSavedAt"`
	CurrentVersion int       `json:"currentVersion"`
	IsSaving       bool      `json:"isSaving"`
	WordCount      int       `json:"wordCount"`
}

// DocumentContentResponse 文档内容响应
type DocumentContentResponse struct {
	DocumentID string    `json:"documentId"`
	Content    string    `json:"content"`
	Version    int       `json:"version"`
	WordCount  int       `json:"wordCount"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

// UpdateContentRequest 更新内容请求
type UpdateContentRequest struct {
	DocumentID string `json:"documentId" validate:"required"`
	Content    string `json:"content" validate:"required"`
	Version    int    `json:"version"` // 用于版本冲突检测
}
