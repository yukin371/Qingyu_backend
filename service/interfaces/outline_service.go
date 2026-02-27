package interfaces

import (
	"Qingyu_backend/models/writer"
	"context"
)

// OutlineService 大纲服务接口
type OutlineService interface {
	// 基础CRUD
	Create(ctx context.Context, projectID, userID string, req *CreateOutlineRequest) (*writer.OutlineNode, error)
	GetByID(ctx context.Context, outlineID, projectID string) (*writer.OutlineNode, error)
	List(ctx context.Context, projectID string) ([]*writer.OutlineNode, error)
	Update(ctx context.Context, outlineID, projectID string, req *UpdateOutlineRequest) (*writer.OutlineNode, error)
	Delete(ctx context.Context, outlineID, projectID string) error

	// 树形结构操作
	GetTree(ctx context.Context, projectID string) ([]*OutlineTreeNode, error)
	GetChildren(ctx context.Context, projectID, parentID string) ([]*writer.OutlineNode, error)
}

// CreateOutlineRequest 创建大纲请求
type CreateOutlineRequest struct {
	Title      string   `json:"title" validate:"required"`
	ParentID   string   `json:"parentId"`
	Summary    string   `json:"summary"`
	Type       string   `json:"type"`
	Tension    int      `json:"tension"`
	ChapterID  string   `json:"chapterId"`
	Characters []string `json:"characters"`
	Items      []string `json:"items"`
	Order      int      `json:"order"`
}

// UpdateOutlineRequest 更新大纲请求
type UpdateOutlineRequest struct {
	Title      *string   `json:"title"`
	ParentID   *string   `json:"parentId"`
	Summary    *string   `json:"summary"`
	Type       *string   `json:"type"`
	Tension    *int      `json:"tension"`
	ChapterID  *string   `json:"chapterId"`
	Characters *[]string `json:"characters"`
	Items      *[]string `json:"items"`
	Order      *int      `json:"order"`
}

// OutlineTreeNode 大纲树节点
type OutlineTreeNode struct {
	*writer.OutlineNode
	Children []*OutlineTreeNode `json:"children,omitempty"`
}
