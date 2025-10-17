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
